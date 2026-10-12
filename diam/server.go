// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server, based on net/http.

package diam

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// The Handler interface allow arbitrary objects to be
// registered to serve particular messages like CER, DWR.
type Handler interface {
	// ServeDIAM should write messages to the Conn and then return.
	// Returning signals that the request is finished.
	//
	// If Server.MaxConcurrentHandlers != 0, ServeDIAM may be invoked
	// concurrently on the same connection from multiple goroutines;
	// handlers that maintain per-connection state must synchronize
	// access themselves. By default (MaxConcurrentHandlers == 0) the
	// server dispatches messages on a connection sequentially.
	ServeDIAM(Conn, *Message)
}

// Conn interface is used by a handler to send diameter messages.
type Conn interface {
	Write(b []byte) (int, error)                    // Writes a msg to the connection
	WriteStream(b []byte, stream uint) (int, error) // Writes a msg to the connection's stream
	Close()                                         // Close the connection
	LocalAddr() net.Addr                            // Returns the local IP
	RemoteAddr() net.Addr                           // Returns the remote IP
	TLS() *tls.ConnectionState                      // TLS or nil when not using TLS
	Dictionary() *dict.Parser                       // Dictionary parser of the connection
	Context() context.Context                       // Returns the internal context
	SetContext(ctx context.Context)                 // Stores a new context
	Connection() net.Conn                           // Returns network connection
}

// The CloseNotifier interface is implemented by Conns which
// allow detecting when the underlying connection has gone away.
//
// This mechanism can be used to detect if a peer has disconnected.
type CloseNotifier interface {
	// CloseNotify returns a channel that is closed
	// when the client connection has gone away.
	CloseNotify() <-chan struct{}
}

// TLSUpgrader is implemented by Conns that can be upgraded to TLS after the
// capabilities exchange, as negotiated with Inband-Security-Id (RFC 6733 §13).
// RFC 6733 retains this only for backward compatibility with RFC 3588 peers
// (§1.1.3, §2.1, §5.6). New deployments should prefer a TLS listener and hand
// the Server an already-encrypted connection, which §5.6 and §13 recommend
// because it also protects the CER and CEA.
//
// CanStartTLS reports whether this particular connection can be upgraded, and
// must be consulted before a success CEA is sent: once both sides have agreed
// on TLS in the CER/CEA exchange, the handshake is expected to follow. Only
// buffered stream connections qualify, and only while their handler runs on
// the serve goroutine (see InbandTLSHandler). A connection that is already
// encrypted cannot be upgraded again, and neither can a MultistreamConn:
// RFC 6733 §2.1 runs DTLS, not TLS, over SCTP.
//
//	if u, ok := c.(diam.TLSUpgrader); ok && u.CanStartTLS() {
//	    err := u.StartTLS(tlsConfig)
//	}
//
// If the handshake fails, the connection must be closed: RFC 6733 §13 requires
// both ends to move to the closed state.
type TLSUpgrader interface {
	// CanStartTLS reports whether StartTLS or StartTLSClient can succeed on
	// this connection. It does not change the connection.
	CanStartTLS() bool

	// StartTLS upgrades the connection to TLS as the server side
	// (tls.Server). It must be called from the handler that is processing
	// the capabilities exchange, before the serve loop reads again.
	StartTLS(cfg *tls.Config) error

	// StartTLSClient upgrades the connection to TLS as the client side
	// (tls.Client), under the same constraint as StartTLS.
	StartTLSClient(cfg *tls.Config) error
}

// InbandTLSHandler is implemented by Handlers that may upgrade a connection
// to TLS from inside the handler for the capabilities exchange. A Server asks
// its Handler once per connection. When the answer is true and
// MaxConcurrentHandlers is non-zero, the capabilities exchange (CER and CEA)
// on a connection that can still be upgraded is dispatched on the serve
// goroutine instead of a new one; this is a behaviour change for such
// Handlers only. A Handler that does not implement this interface, or reports
// false, keeps the dispatch behaviour it had, and its connections report
// CanStartTLS false under concurrent dispatch.
//
// sm.StateMachine implements this, reporting whether Settings.TLSConfig is
// set. A Handler that wraps a StateMachine must forward InbandTLSEnabled for
// in-band TLS to be available under concurrent dispatch.
type InbandTLSHandler interface {
	InbandTLSEnabled() bool
}

// inbandTLSEnabled reports whether h may upgrade a connection to TLS from
// within the capabilities-exchange handler.
func inbandTLSEnabled(h Handler) bool {
	if h == nil {
		h = DefaultServeMux
	}
	u, ok := h.(InbandTLSHandler)
	return ok && u.InbandTLSEnabled()
}

// A liveSwitchReader is a switchReader that's safe for concurrent
// reads and switches, if its mutex is held.
type liveSwitchReader struct {
	sync.Mutex
	r         io.Reader
	pr        *io.PipeReader
	pipeCopyF func()
}

func (sr *liveSwitchReader) Read(p []byte) (n int, err error) {
	sr.Lock()
	// Check if closeNotifier was created prior to this Read call & start it
	if sr.pr != nil && sr.pipeCopyF != nil {
		go sr.pipeCopyF()
		sr.r = sr.pr
		sr.pr = nil
		sr.pipeCopyF = nil
	}
	r := sr.r
	sr.Unlock()
	return r.Read(p)
}

// conn represents the server side of a diameter connection.
type conn struct {
	server   *Server              // the Server on which the connection arrived
	rwc      net.Conn             // i/o connection
	sr       liveSwitchReader     // reads from rwc
	buf      *bufio.ReadWriter    // buffered(sr, rwc)
	tlsState *tls.ConnectionState // or nil when not using TLS
	writer   *response            // the diam.Conn exposed to handlers

	hwg       sync.WaitGroup // tracks in-flight handler goroutines
	sem       chan struct{}  // bounds concurrent handlers; nil = unbounded/sequential
	accepted  bool           // accepted by Server.Serve, as opposed to dialed
	inbandTLS bool           // Handler may upgrade this connection to TLS
	inline    atomic.Bool    // a handler is running on the serve goroutine

	// tmu guards rwc and tlsState against an in-band TLS upgrade, which
	// replaces both while other goroutines may be inspecting or closing the
	// connection through its diam.Conn. The serve goroutine, which is the
	// only writer, reads them without it.
	tmu sync.RWMutex

	mu           sync.Mutex // guards the following
	closeNotifyc chan struct{}
	clientGone   bool
}

func (c *conn) closeNotify() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc == nil {
		c.closeNotifyc = make(chan struct{})
		if c.clientGone {
			// The connection went away before anyone asked; deliver that now rather than
			// returning a channel nothing will ever close.
			close(c.closeNotifyc)
			return c.closeNotifyc
		}

		if msc, isMulti := c.netConn().(MultistreamConn); isMulti {
			// MultistreamConn provides it's own error handler
			msc.SetErrorHandler(func(mc MultistreamConn, err error) {
				mc.Close()
				c.notifyClientGone()
			})
		} else {
			pr, pw := io.Pipe()
			c.sr.Lock()
			readSource := c.sr.r
			c.sr.pr = pr
			// Create closeNotifier pipe copy routine, but do not start it here
			// If we start it immediately, pipe Write can block indefinitely if we are already in
			// liveSwitchReader.Read() with original sr.r since Pipe.Write blocks in absence of corresponding
			// pipe reader
			// We should only swap the reader outside of r.Read call
			c.sr.pipeCopyF = func() {
				_, err := io.Copy(pw, readSource)
				if err == nil {
					err = io.EOF
				}
				pw.CloseWithError(err)
				c.notifyClientGone()
			}
			c.sr.Unlock()
		}
	}
	return c.closeNotifyc
}

func (c *conn) notifyClientGone() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.clientGone {
		return
	}
	// Record the close even if nobody has called CloseNotify yet, so a later call can see it.
	c.clientGone = true
	if c.closeNotifyc != nil {
		close(c.closeNotifyc) // unblock readers
	}
}

// netConn returns the current transport, which an in-band TLS upgrade may
// replace. Use it wherever the connection is reached from a goroutine other
// than the serve loop.
func (c *conn) netConn() net.Conn {
	c.tmu.RLock()
	defer c.tmu.RUnlock()
	return c.rwc
}

// tls returns the TLS state, or nil, under the same rules as netConn.
func (c *conn) tls() *tls.ConnectionState {
	c.tmu.RLock()
	defer c.tmu.RUnlock()
	return c.tlsState
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	msc, isMulti := rwc.(MultistreamConn)
	if isMulti {
		c = &conn{
			server: srv,
			rwc:    msc,
		}
	} else {
		c = &conn{
			server: srv,
			rwc:    rwc,
			sr:     liveSwitchReader{r: rwc},
		}
		c.buf = bufio.NewReadWriter(bufio.NewReader(&c.sr), bufio.NewWriter(rwc))
	}
	c.writer = &response{conn: c}
	c.inbandTLS = inbandTLSEnabled(srv.Handler)
	if n := srv.MaxConcurrentHandlers; n > 0 {
		c.sem = make(chan struct{}, n)
	}
	return c, nil
}

// Read next message from connection.
func (c *conn) readMessage() (m *Message, err error) {
	if c.server.ReadTimeout > 0 {
		c.rwc.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
	}
	if msc, isMulti := c.rwc.(MultistreamConn); isMulti {
		// If it's a multi-stream association - reset the stream to "undefined" prior to reading next message
		msc.ResetCurrentStream()
		m, err = ReadMessage(msc, c.dictionary()) // MultistreamConn has it's own buffering
	} else {
		m, err = ReadMessage(c.buf.Reader, c.dictionary())
	}
	return m, err
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("diam: panic serving %v: %v\n%s",
				c.rwc.RemoteAddr().String(), err, buf)
		}
		// Wait for in-flight handler goroutines to finish so they are
		// not writing to a closed connection when we call rwc.Close().
		c.hwg.Wait()
		c.rwc.Close()
		c.notifyClientGone()
	}()
	if c.accepted {
		h := c.server.Handler
		if h == nil {
			h = DefaultServeMux
		}
		if ah, ok := h.(AcceptHandler); ok {
			if onClose := ah.HandleAccept(c.writer); onClose != nil {
				defer onClose()
			}
		}
	}
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			return
		}
		c.tlsState = &tls.ConnectionState{}
		*c.tlsState = tlsConn.ConnectionState()
	}
	if cb := c.server.OnNewConnection; cb != nil {
		cb(c.writer)
	}
	for {
		m, err := c.readMessage()
		if err != nil {
			// Report errors to the channel, except EOF.
			// Connection close is handled by the defer above,
			// after draining in-flight handlers via hwg.Wait().
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				h := c.server.Handler
				if h == nil {
					h = DefaultServeMux
				}
				if er, ok := h.(ErrorReporter); ok {
					er.Error(&ErrorReport{c.writer, m, err})
				}
			}
			break
		}
		c.dispatch(m)
	}
}

// dispatch runs the handler for m, either on the serve goroutine
// (sequential, default) or in a new goroutine (concurrent), depending
// on Server.MaxConcurrentHandlers. A positive value bounds concurrency
// via a per-connection semaphore; a negative value is unbounded.
//
// A capabilities exchange runs on the serve goroutine even under concurrent
// dispatch, but only when the Handler reports in-band TLS as enabled (see
// InbandTLSHandler) and this connection could still be upgraded. The handler
// may call StartTLS, and the upgrade has to finish before the next read, or
// the serve loop would read the peer's TLS ClientHello as a diameter message.
// Everything else is unaffected: a Handler that does not offer in-band
// security, a connection that already arrived over TLS, an SCTP connection,
// and every command other than CER dispatch exactly as before.
func (c *conn) dispatch(m *Message) {
	if c.server.MaxConcurrentHandlers == 0 ||
		(c.inbandTLS && m.Header.CommandCode == CapabilitiesExchange && c.canStartTLS()) {
		// Sequential dispatch preserves the historical Handler contract.
		c.inline.Store(true)
		defer c.inline.Store(false)
		serverHandler{c.server}.ServeDIAM(c.writer, m)
		return
	}
	if c.sem != nil {
		c.sem <- struct{}{} // blocks when MaxConcurrentHandlers reached
	}
	c.hwg.Add(1)
	go func() {
		defer c.hwg.Done()
		defer func() {
			if c.sem != nil {
				<-c.sem
			}
		}()
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 4096)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("diam: panic serving %v: %v\n%s",
					c.netConn().RemoteAddr().String(), err, buf)
			}
		}()
		serverHandler{c.server}.ServeDIAM(c.writer, m)
	}()
}

// dictionary returns the dictionary parser associated to the Server instance
// or dict.Default.
func (c *conn) dictionary() *dict.Parser {
	if c.server.Dict == nil {
		return dict.Default
	}
	return c.server.Dict
}

// A response represents the server side of a diameter response.
// It implements the Conn and CloseNotifier interfaces.
type response struct {
	mu   sync.Mutex      // guards conn and Write
	conn *conn           // socket, reader and writer
	xmu  sync.Mutex      // guards ctx
	ctx  context.Context // context for this Conn
}

// Write writes the message m to the connection.
func (w *response) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writeLocked(b)
}

// writeLocked performs the write, assuming w.mu is already held.
// Serializes access to the underlying (TCP or SCTP) connection so
// concurrently-dispatched handlers do not interleave bytes.
func (w *response) writeLocked(b []byte) (int, error) {
	if w.conn.server.WriteTimeout > 0 {
		w.conn.rwc.SetWriteDeadline(time.Now().Add(w.conn.server.WriteTimeout))
	}
	msc, isMulti := w.conn.rwc.(MultistreamConn) // Note - SetWriteDeadline is not currently supported for SCTP
	if isMulti {                                 // don't use buffered writer for muti-streamming writes it'll mix up streams
		return msc.Write(b)
	}
	n, err := w.conn.buf.Writer.Write(b)
	if err != nil {
		return 0, err
	}
	if err = w.conn.buf.Writer.Flush(); err != nil {
		return 0, err
	}
	return n, nil
}

// WriteStream of MultistreamWriter interface
func (w *response) WriteStream(b []byte, stream uint) (int, error) {
	// TODO - SetWriteDeadline is not currently supported
	w.mu.Lock()
	defer w.mu.Unlock()
	if msc, isMulti := w.conn.rwc.(MultistreamConn); isMulti {
		// don't use buffered writer for muti-streamming writes it'll mix up streams
		return msc.WriteStream(b, stream)
	}
	return w.writeLocked(b)
}

// CurrentWriterStream of MultistreamWriter interface
func (w *response) CurrentWriterStream() uint {
	if msc, isMulti := w.conn.netConn().(MultistreamConn); isMulti {
		return msc.CurrentWriterStream()
	}
	return 0
}

// ResetWriterStream of MultistreamWriter interface
func (w *response) ResetWriterStream() {
	if msc, isMulti := w.conn.netConn().(MultistreamConn); isMulti {
		msc.CurrentWriterStream()
	}
}

// SetWriterStream of MultistreamWriter interface
func (w *response) SetWriterStream(stream uint) uint {
	if msc, isMulti := w.conn.netConn().(MultistreamConn); isMulti {
		return msc.SetWriterStream(stream)
	}
	return 0
}

// Close closes the connection.
func (w *response) Close() {
	w.conn.netConn().Close()
}

// LocalAddr returns the local address of the connection.
func (w *response) LocalAddr() net.Addr {
	return w.conn.netConn().LocalAddr()
}

// RemoteAddr returns the peer address of the connection.
func (w *response) RemoteAddr() net.Addr {
	return w.conn.netConn().RemoteAddr()
}

// TLS returns the TLS connection state, or nil.
func (w *response) TLS() *tls.ConnectionState {
	return w.conn.tls()
}

// Dictionary returns the dictionary parser associated to this connection.
// If none was provided then it returns the default dictionary.
func (w *response) Dictionary() *dict.Parser {
	return w.conn.dictionary()
}

// CloseNotify implements the CloseNotifier interface.
func (w *response) CloseNotify() <-chan struct{} {
	return w.conn.closeNotify()
}

// Context returns the internal context or a new context.Background.
func (w *response) Context() context.Context {
	w.xmu.Lock()
	defer w.xmu.Unlock()
	if w.ctx == nil {
		w.ctx = context.Background()
	}
	return w.ctx
}

// SetContext replaces the internal context with the given one.
func (w *response) SetContext(ctx context.Context) {
	w.xmu.Lock()
	w.ctx = ctx
	w.xmu.Unlock()
}

func (w *response) Connection() net.Conn {
	return w.conn.netConn()
}

// CanStartTLS reports whether this connection can be upgraded to TLS now. It
// is false unless the calling handler runs on the serve goroutine, because the
// upgrade must complete before the serve loop reads again.
func (w *response) CanStartTLS() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.inline.Load() && w.conn.canStartTLS()
}

// canStartTLS reports whether the transport allows an in-band upgrade.
// Callers hold the response mutex, or are the serve goroutine.
func (c *conn) canStartTLS() bool {
	if c.tlsState != nil {
		return false // Already encrypted.
	}
	if _, isMulti := c.rwc.(MultistreamConn); isMulti {
		// RFC 6733 §2.1 secures SCTP with DTLS, not with this handshake,
		// and newConn leaves a multistream connection without the buffered
		// reader the upgrade needs.
		return false
	}
	return c.buf != nil
}

// StartTLS upgrades the connection to TLS as the server side.
//
// It must be called from the handler processing the capabilities exchange
// while that handler runs on the serve goroutine, so that no read is in
// flight; otherwise it returns an error. See CanStartTLS.
func (w *response) StartTLS(cfg *tls.Config) error {
	return w.startTLS(cfg, true)
}

// StartTLSClient upgrades the connection to TLS as the client side, under the
// same constraint as StartTLS.
func (w *response) StartTLSClient(cfg *tls.Config) error {
	return w.startTLS(cfg, false)
}

func (w *response) startTLS(cfg *tls.Config, server bool) error {
	if cfg == nil {
		return errors.New("diam: StartTLS: nil tls.Config")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	c := w.conn
	if c.tlsState != nil {
		return errors.New("diam: StartTLS: connection is already using TLS")
	}
	if !c.canStartTLS() {
		return fmt.Errorf("diam: StartTLS: connection to %s cannot be upgraded",
			c.rwc.RemoteAddr())
	}
	if !c.inline.Load() {
		return errors.New("diam: StartTLS: not called from a handler running on the serve goroutine")
	}

	// Take over the connection's current read source rather than reading the
	// socket behind it. When a CloseNotifier has been registered the socket is
	// being drained by a copy goroutine into a pipe, and reading the socket
	// here would compete with it for the same bytes.
	c.sr.Lock()
	if c.sr.pr != nil && c.sr.pipeCopyF != nil {
		// A CloseNotifier was registered but the serve loop has not read
		// since, so the pipe is not running yet. Start it here, exactly as
		// liveSwitchReader.Read would, so that close notification keeps
		// working across the upgrade.
		go c.sr.pipeCopyF()
		c.sr.r = c.sr.pr
		c.sr.pr = nil
		c.sr.pipeCopyF = nil
	}
	readSource := c.sr.r
	c.sr.Unlock()

	// The peer may already have sent its ClientHello, which the buffered
	// reader will have consumed along with the last diameter message. Those
	// bytes belong to the handshake.
	var handshakeReader io.Reader = readSource
	if n := c.buf.Reader.Buffered(); n > 0 {
		pending := make([]byte, n)
		if _, err := io.ReadFull(c.buf.Reader, pending); err != nil {
			// The bytes are consumed either way, so the stream cannot be
			// handed to TLS intact; report it instead of losing them.
			return fmt.Errorf("diam: StartTLS: reading %d buffered bytes: %w", n, err)
		}
		handshakeReader = io.MultiReader(bytes.NewReader(pending), readSource)
	}

	// Writes go straight to the socket; reads come through the chain above.
	raw := &upgradeConn{Conn: c.rwc, r: handshakeReader}
	var tlsConn *tls.Conn
	if server {
		tlsConn = tls.Server(raw, cfg)
	} else {
		tlsConn = tls.Client(raw, cfg)
	}
	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("diam: StartTLS: handshake with %s: %w",
			c.rwc.RemoteAddr(), err)
	}

	// Publish the upgraded connection. sr.r is set to the tls.Conn, which
	// reads through handshakeReader, so it must not be reached via sr again.
	state := tlsConn.ConnectionState()
	c.sr.Lock()
	c.sr.r = tlsConn
	c.sr.Unlock()
	c.buf = bufio.NewReadWriter(bufio.NewReader(&c.sr), bufio.NewWriter(tlsConn))
	c.tmu.Lock()
	c.rwc = tlsConn
	c.tlsState = &state
	c.tmu.Unlock()
	return nil
}

// upgradeConn presents the connection's existing read chain and its socket as
// one net.Conn, so the TLS handshake reads the bytes that have already been
// buffered or piped instead of racing whoever is draining the socket.
type upgradeConn struct {
	net.Conn
	r io.Reader
}

func (c *upgradeConn) Read(b []byte) (int, error) { return c.r.Read(b) }

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as diameter handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(Conn, *Message)

// ServeDIAM calls f(c, m).
func (f HandlerFunc) ServeDIAM(c Conn, m *Message) {
	f(c, m)
}

// The ErrorReporter interface is implemented by Handlers that
// allow reading errors from the underlying connection, like
// parsing diameter messages or connection errors.
type ErrorReporter interface {
	// Error writes an error to the reporter.
	Error(err *ErrorReport)

	// ErrorReports returns a channel that receives
	// errors from the connection.
	ErrorReports() <-chan *ErrorReport
}

// The AcceptHandler interface is implemented by Handlers that need to know
// when a Server accepts a connection. HandleAccept is called once for each
// connection accepted by Server.Serve, before the TLS handshake and before
// the first read; it is not called for connections created by Dial. It runs
// on the connection's goroutine and must not block. If it returns a non-nil
// function, that function is called once when the connection is closed.
type AcceptHandler interface {
	HandleAccept(c Conn) (onClose func())
}

// ErrorReport is sent out of the server in case it fails to
// read messages due to a bad dictionary or network errors.
type ErrorReport struct {
	Conn    Conn     // Peer that caused the error
	Message *Message // Message that caused the error
	Error   error    // Error message
}

// String returns an error message. It does not render the Message field.
func (er *ErrorReport) String() string {
	if er.Conn == nil {
		return fmt.Sprintf("diameter error: %s", er.Error)
	}
	return fmt.Sprintf("diameter error on %s: %s", er.Conn.RemoteAddr(), er.Error)
}

// ServeMux is a diameter message multiplexer. It matches the
// command from the incoming message against a list of
// registered commands and calls the handler.
type ServeMux struct {
	e      chan *ErrorReport
	mu     sync.RWMutex // Guards m.
	m      map[string]muxEntry
	idxMap map[CommandIndex]muxEntry
}

type muxEntry struct {
	h      Handler
	cmd    string
	cmdIdx CommandIndex
}

type CommandIndex struct {
	AppID   uint32
	Code    uint32
	Request bool
}

var ALL_CMD_INDEX = CommandIndex{^uint32(0), ^uint32(0), false}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		e:      make(chan *ErrorReport, 1),
		m:      make(map[string]muxEntry),
		idxMap: make(map[CommandIndex]muxEntry),
	}
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// Error implements the ErrorReporter interface.
func (mux *ServeMux) Error(err *ErrorReport) {
	select {
	case mux.e <- err:
	default:
	}
}

// ErrorReports implement the ErrorReporter interface.
func (mux *ServeMux) ErrorReports() <-chan *ErrorReport {
	return mux.e
}

// ServeDIAM dispatches the request to the handler that match the code
// in the incoming message. If the special "ALL" handler is registered
// it is used as a catch-all. Otherwise an ErrorReport is sent out.
func (mux *ServeMux) ServeDIAM(c Conn, m *Message) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	dcmd, err := m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode)

	if err != nil {
		// Try the catch-all.
		mux.serveIdx(ALL_CMD_INDEX, c, m)
		return
	}

	idx := CommandIndex{
		m.Header.ApplicationID,
		m.Header.CommandCode,
		m.Header.CommandFlags&RequestFlag == RequestFlag}
	_, ok := mux.idxMap[idx]
	if ok {
		mux.serveIdx(idx, c, m)
		return
	}

	var cmd string
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		cmd = dcmd.Short + "R"
	} else {
		cmd = dcmd.Short + "A"
	}
	mux.serve(cmd, c, m)
}

func (mux *ServeMux) serveIdx(cmd CommandIndex, c Conn, m *Message) {
	entry, ok := mux.idxMap[cmd]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	// Try catch-all.
	entry, ok = mux.idxMap[ALL_CMD_INDEX]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	mux.Error(&ErrorReport{
		Conn:    c,
		Message: m,
		Error:   fmt.Errorf("unhandled message for index: %+v", cmd),
	})
}

func (mux *ServeMux) serve(cmd string, c Conn, m *Message) {
	entry, ok := mux.m[cmd]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	// Try catch-all.
	entry, ok = mux.idxMap[ALL_CMD_INDEX]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	mux.Error(&ErrorReport{
		Conn:    c,
		Message: m,
		Error:   fmt.Errorf("unhandled message for '%s'", cmd),
	})
}

// Handle registers the handler for the given code.
// If a handler already exists for code, Handle panics.
func (mux *ServeMux) Handle(shortCmd string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("DIAM: nil handler")
	}
	if shortCmd == "ALL" {
		mux.idxMap[ALL_CMD_INDEX] = muxEntry{h: handler, cmd: shortCmd}
		return
	}
	mux.m[shortCmd] = muxEntry{h: handler, cmd: shortCmd}
}

// Handle registers the handler for the given code.
// If a handler already exists for code, Handle panics.
func (mux *ServeMux) HandleIdx(cmd CommandIndex, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("DIAM: nil handler")
	}
	mux.idxMap[cmd] = muxEntry{h: handler, cmdIdx: cmd}
}

// HandleFunc registers the handler function for the given command.
// Special cmd "ALL" may be used as a catch all.
func (mux *ServeMux) HandleFunc(cmd string, handler func(Conn, *Message)) {
	mux.Handle(cmd, HandlerFunc(handler))
}

// Handle registers the handler object for the given command
// in the DefaultServeMux.
func Handle(cmd string, handler Handler) {
	DefaultServeMux.Handle(cmd, handler)
}

// HandleFunc registers the handler function for the given command
// in the DefaultServeMux.
func HandleFunc(cmd string, handler func(Conn, *Message)) {
	DefaultServeMux.HandleFunc(cmd, handler)
}

// ErrorReports returns the ErrorReport channel of the DefaultServeMux.
func ErrorReports() <-chan *ErrorReport {
	return DefaultServeMux.ErrorReports()
}

// Serve accepts incoming diameter connections on the listener l,
// creating a new service goroutine for each.  The service goroutines
// read messages and then call handler to reply to them.
// Handler is typically nil, in which case the DefaultServeMux is used.
func Serve(l net.Listener, handler Handler) error {
	srv := &Server{Handler: handler}
	return srv.Serve(l)
}

// A Server defines parameters for running a diameter server.
type Server struct {
	Network      string        // network of the address - empty string defaults to tcp
	Addr         string        // address to listen on, ":3868" if empty
	Handler      Handler       // handler to invoke, DefaultServeMux if nil
	Dict         *dict.Parser  // diameter dictionaries for this server
	ReadTimeout  time.Duration // maximum duration before timing out read of the request
	WriteTimeout time.Duration // maximum duration before timing out write of the response
	TLSConfig    *tls.Config   // optional TLS config, used by ListenAndServeTLS
	LocalAddr    net.Addr      // optional Local Address to bind dailer's (Dail...) socket to

	// MaxConcurrentHandlers controls per-connection handler dispatch.
	//   0 (default) - sequential dispatch: each message is handled to
	//                 completion before the next is read. Preserves the
	//                 historical Handler contract.
	//   >0          - concurrent dispatch bounded by this many in-flight
	//                 handlers per connection; once reached, the read
	//                 loop blocks until a slot frees.
	//   <0          - unbounded concurrent dispatch (not recommended for
	//                 untrusted peers).
	//
	// Why enable concurrent dispatch:
	// Sequential dispatch caps per-connection throughput at
	// 1 / handler_latency, regardless of network or CPU headroom. Real
	// handlers typically do I/O (DB lookups, auth, logging, remote calls)
	// in the 1-10ms range, which caps a connection at 100-1000 msg/s.
	// Concurrent dispatch lets independent requests proceed in parallel,
	// hiding handler latency and raising throughput to roughly
	// MaxConcurrentHandlers / handler_latency until CPU saturates.
	// Measured locally with a 1ms handler: ~810 msg/s sequential vs
	// ~39,900 msg/s with MaxConcurrentHandlers=256 (~49x).
	//
	// Diameter matches answers to requests by Hop-by-Hop ID, not arrival
	// order, so concurrent dispatch does not change what the peer sees.
	// Handlers that mutate shared per-connection state must synchronize
	// themselves.
	MaxConcurrentHandlers int

	// OnNewConnection, if non-nil, is invoked once per accepted connection
	// after the transport is fully established (TLS handshake complete, if
	// applicable) and before the read loop starts. It runs in the
	// connection's serve goroutine, so long-running work will delay message
	// processing on that connection. Type-assert to CloseNotifier to detect
	// disconnection:
	//
	//	srv.OnNewConnection = func(c diam.Conn) {
	//		log.Printf("up: %s", c.RemoteAddr())
	//		go func() {
	//			<-c.(diam.CloseNotifier).CloseNotify()
	//			log.Printf("down: %s", c.RemoteAddr())
	//		}()
	//	}
	OnNewConnection func(Conn)

	mu        sync.Mutex
	listeners map[net.Listener]struct{}
	closed    bool
}

// ErrServerClosed is returned by Server.Serve and Server.ListenAndServe(TLS)
// after a call to Server.Close.
var ErrServerClosed = fmt.Errorf("diam: Server closed")

// Close immediately closes all listeners registered with the server. It does
// not affect already-accepted connections or in-flight handlers; those
// continue to run until their read loop exits naturally or their underlying
// connection is closed by the peer. After Close, Server.Serve returns
// ErrServerClosed and no new connections are accepted.
func (srv *Server) Close() error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.closed = true
	var firstErr error
	for l := range srv.listeners {
		if err := l.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	srv.listeners = nil
	return firstErr
}

func (srv *Server) trackListener(l net.Listener, add bool) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if add {
		if srv.closed {
			return false
		}
		if srv.listeners == nil {
			srv.listeners = make(map[net.Listener]struct{})
		}
		srv.listeners[l] = struct{}{}
		return true
	}
	delete(srv.listeners, l)
	return true
}

func (srv *Server) isClosed() bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.closed
}

// onceCloseListener wraps a net.Listener so that Close is idempotent.
// Used internally by ListenAndServe(TLS) so the defer l.Close() and
// Server.Close do not both close the underlying listener, which triggers
// file-descriptor reuse races on SCTP (see 29cbaef).
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(func() { oc.closeErr = oc.Listener.Close() })
	return oc.closeErr
}

// serverHandler delegates to either the server's Handler or DefaultServeMux.
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeDIAM(w Conn, m *Message) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	handler.ServeDIAM(w, m)
}

// ListenAndServe listens on the network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
//
// If srv.Network is blank, "tcp" is used
// If srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServe() error {
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	l, e := MultistreamListen(network, addr)
	if e != nil {
		return e
	}
	l = &onceCloseListener{Listener: l}
	defer l.Close()
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call srv.Handler to reply to them.
// The caller is responsible for closing l when Serve returns.
// Serve returns ErrServerClosed after srv.Close is called.
func (srv *Server) Serve(l net.Listener) error {
	if !srv.trackListener(l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(l, false)
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if srv.isClosed() {
				return ErrServerClosed
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("diam: accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			network := "<nil>"
			address := network
			addr := l.Addr()
			if addr != nil {
				network = addr.Network()
				address = addr.String()
			}
			log.Printf("diam: accept error: %v for %s %s", e, network, address)
			return e
		}
		tempDelay = 0
		if c, err := srv.newConn(rw); err != nil {
			log.Printf("srv.newConn error: %v", err)
			continue
		} else {
			c.accepted = true
			go c.serve()
		}
	}
}

// ListenAndServeNetwork listens on the network & addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If handler is nil, DefaultServeMux is used.
//
// If dict is nil, dict.Default is used.
func ListenAndServeNetwork(network, addr string, handler Handler, dp *dict.Parser) error {
	server := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp}
	return server.ListenAndServe()
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If handler is nil, DefaultServeMux is used.
//
// If dict is nil, dict.Default is used.
func ListenAndServe(addr string, handler Handler, dp *dict.Parser) error {
	return ListenAndServeNetwork("tcp", addr, handler, dp)
}

// ListenAndServeTLS listens on the network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
//
// Either filenames containing a certificate and matching private key
// for the server must be provided either the callback
// srv.TLSConfig.GetCertificate should be filled in advance.
// If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// If srv.Network is blank, "tcp" is used
// If srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var config *tls.Config
	if srv.TLSConfig == nil {
		config = new(tls.Config)
	} else {
		config = TLSConfigClone(srv.TLSConfig)
	}
	var err error
	if config.GetCertificate == nil {
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}
	conn, err := Listen(network, addr)
	if err != nil {
		return err
	}
	var tlsListener net.Listener = tls.NewListener(conn, config)
	tlsListener = &onceCloseListener{Listener: tlsListener}
	defer tlsListener.Close()
	return srv.Serve(tlsListener)
}

// ListenAndServeNetworkTLS acts identically to ListenAndServeNetwork, except that it
// expects SSL connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeNetworkTLS(network, addr string, certFile string, keyFile string, handler Handler, dp *dict.Parser) error {
	server := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp}
	return server.ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects SSL connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler, dp *dict.Parser) error {
	return ListenAndServeNetworkTLS("tcp", addr, certFile, keyFile, handler, dp)
}
