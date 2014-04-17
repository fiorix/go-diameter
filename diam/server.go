// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server, based on net/http.

package diam

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/diam/dict"
)

// Objects implementing the Handler interface can be
// registered to serve particular messages like CER, DWR.
//
// ServeDiam should write messages to the Conn and then return.
// Returning signals that the request is finished and that the server
// can move on to the next request on the connection.
type Handler interface {
	ServeDiam(Conn, *Message)
	ErrorReports() chan ErrorReport
}

// The CloseNotifier interface is implemented by Conns which
// allow detecting when the underlying connection has gone away.
//
// This mechanism can be used to detect if the client has disconnected.
type CloseNotifier interface {
	// CloseNotify returns a channel that receives a single value
	// when the client connection has gone away.
	CloseNotify() <-chan bool
}

// Conn interface is used by a handler to send diameter messages.
type Conn interface {
	Write(b []byte) (int, error) // Writes a msg to the connection
	Close()                      // Close the connection
	LocalAddr() net.Addr         // Local IP
	RemoteAddr() net.Addr        // Remote IP
	TLS() *tls.ConnectionState   // or nil when not using TLS
}

// A liveSwitchReader is a switchReader that's safe for concurrent
// reads and switches, if its mutex is held.
type liveSwitchReader struct {
	sync.Mutex
	r io.Reader
}

func (sr *liveSwitchReader) Read(p []byte) (n int, err error) {
	sr.Lock()
	r := sr.r
	sr.Unlock()
	return r.Read(p)
}

// conn represents the server side of a diameter connection.
type sConn struct {
	server   *Server              // the Server on which the connection arrived
	rwc      net.Conn             // i/o connection
	sr       liveSwitchReader     // reads from rwc
	lr       *io.LimitedReader    // io.LimitedReader(sr)
	buf      *bufio.ReadWriter    // buffered(lr, rwc) reading from bufio->limitReader->sr->rwc
	tlsState *tls.ConnectionState // or nil when not using TLS

	mu           sync.Mutex // guards the following
	closeNotifyc chan bool
}

func (c *sConn) closeNotify() <-chan bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc == nil {
		c.closeNotifyc = make(chan bool, 1)
		pr, pw := io.Pipe()
		readSource := c.sr.r
		c.sr.Lock()
		c.sr.r = pr
		c.sr.Unlock()
		go func() {
			_, err := io.Copy(pw, readSource)
			if err == nil {
				err = io.EOF
			}
			pw.CloseWithError(err)
			c.closeNotifyc <- true
		}()
	}
	return c.closeNotifyc
}

// A response represents the server side of a diameter response.
type response struct {
	sConn *sConn
}

// noLimit is an effective infinite upper bound for io.LimitedReader
const noLimit int64 = (1 << 63) - 1

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *sConn, err error) {
	c = new(sConn)
	c.server = srv
	c.rwc = rwc
	c.sr = liveSwitchReader{r: c.rwc}
	c.lr = io.LimitReader(&c.sr, noLimit).(*io.LimitedReader)
	br := newBufioReader(c.lr)
	bw := newBufioWriterSize(c.rwc, 4<<10)
	c.buf = bufio.NewReadWriter(br, bw)
	return c, nil
}

// TODO: use a sync.Cache instead
var (
	bufioReaderCache   = make(chan *bufio.Reader, 4)
	bufioWriterCache2k = make(chan *bufio.Writer, 4)
	bufioWriterCache4k = make(chan *bufio.Writer, 4)
)

func newBufioReader(r io.Reader) *bufio.Reader {
	select {
	case p := <-bufioReaderCache:
		p.Reset(r)
		return p
	default:
		return bufio.NewReader(r)
	}
}

func bufioWriterCache(size int) chan *bufio.Writer {
	switch size {
	case 2 << 10:
		return bufioWriterCache2k
	case 4 << 10:
		return bufioWriterCache4k
	}
	return nil
}

func newBufioWriterSize(w io.Writer, size int) *bufio.Writer {
	select {
	case p := <-bufioWriterCache(size):
		p.Reset(w)
		return p
	default:
		return bufio.NewWriterSize(w, size)
	}
}

// Read next message from connection.
func (c *sConn) readMessage() (*Message, *response, error) {
	dp := c.server.Dict
	if dp == nil {
		dp = dict.Default
	}
	if c.server.ReadTimeout > 0 {
		c.rwc.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
	}
	c.lr.N = int64(HeaderLength) + 4096 /* bufio slop */
	m, err := ReadMessage(c.buf.Reader, dp)
	if err != nil {
		return nil, nil, err
	}
	c.lr.N = noLimit
	return m, &response{sConn: c}, nil
}

// Write writes the message m to the connection.
func (w *response) Write(b []byte) (int, error) {
	if w.sConn.server.WriteTimeout > 0 {
		w.sConn.rwc.SetWriteDeadline(time.Now().Add(w.sConn.server.WriteTimeout))
	}
	defer w.sConn.buf.Writer.Flush()
	return w.sConn.buf.Writer.Write(b)
}

// Close closes the connection.
func (w *response) Close() {
	w.sConn.rwc.Close()
}

// LocalAddr returns the local address of the connection.
func (w *response) LocalAddr() net.Addr {
	return w.sConn.rwc.LocalAddr()
}

// RemoteAddr returns the peer address of the connection.
func (w *response) RemoteAddr() net.Addr {
	return w.sConn.rwc.RemoteAddr()
}

// TLS returns the TLS connection state, or nil.
func (w *response) TLS() *tls.ConnectionState {
	return w.sConn.tlsState
}

// Serve a new connection.
func (c *sConn) serve() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("DIAM: panic serving %v: %v\n%s",
				c.rwc.RemoteAddr().String(), err, buf)
		}
		c.rwc.Close()
	}()
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			return
		}
		c.tlsState = new(tls.ConnectionState)
		*c.tlsState = tlsConn.ConnectionState()
	}
	for {
		m, w, err := c.readMessage()
		if err != nil {
			c.rwc.Close()
			// Report errors to the channel, except EOF.
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				h := c.server.Handler
				if h == nil {
					h = DefaultServeMux
				}
				h.ErrorReports() <- ErrorReport{m, err}
			}
			break
		}
		// Diameter cannot have multiple simultaneous active requests.
		// Until the server replies to this message, it can't
		// read another, so we might as well run the handler in
		// this goroutine.
		serverHandler{c.server}.ServeDiam(w, m)
	}
}

func (w *response) CloseNotify() <-chan bool {
	return w.sConn.closeNotify()
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as diameter handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(Conn, *Message)

// ServeDiam calls f(c, m).
func (f HandlerFunc) ServeDiam(c Conn, m *Message) {
	f(c, m)
}

// ErrorReports calls f.ErrorReports()
func (f HandlerFunc) ErrorReports() chan ErrorReport {
	return f.ErrorReports()
}

// ServeMux is a diameter message multiplexer.
// It matches the command from the incoming message against a list
// of registered commands and calls the handler.
type ServeMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
	e  chan ErrorReport
}

// ErrorReport is sent out of the server in case it fails to read messages
// because of a bad dictionary or network errors.
type ErrorReport struct {
	Message *Message
	Error   error
}

type muxEntry struct {
	h   Handler
	cmd string
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		m: make(map[string]muxEntry),
		e: make(chan ErrorReport),
	}
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// ServeDiam dispatches the request to the handler whose code match
// the incoming message, or close the connection if no handler is found.
func (mux *ServeMux) ServeDiam(c Conn, m *Message) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	var cmd string
	if dcmd, err := m.Dictionary.FindCMD(
		m.Header.ApplicationId,
		m.Header.CommandCode,
	); err != nil {
		cmd = "ALL"
	} else {
		cmd = dcmd.Short
		if m.Header.CommandFlags&0x80 > 0 {
			cmd += "R"
		} else {
			cmd += "A"
		}
	}
	if me, ok := mux.m[cmd]; ok {
		me.h.ServeDiam(c, m)
	} else if me, ok = mux.m["ALL"]; ok {
		me.h.ServeDiam(c, m)
	} else {
		// This is not HTTP. Ignore messages that are not bound
		// instead of closing the connection.
		//c.Close()
	}
}

// Handle registers the handler for the given code.
// If a handler already exists for code, Handle panics.
func (mux *ServeMux) Handle(cmd string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("DIAM: nil handler")
	}
	mux.m[cmd] = muxEntry{h: handler, cmd: cmd}
}

// HandleFunc registers the handler function for the given command.
// Special cmd "ALL" may be used as a catch all.
func (mux *ServeMux) HandleFunc(cmd string, handler func(Conn, *Message)) {
	mux.Handle(cmd, HandlerFunc(handler))
}

// ErrorReports returns the ErrorReport channel of the handler.
func (mux *ServeMux) ErrorReports() chan ErrorReport { return mux.e }

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
func Handle(cmd string, handler Handler) {
	DefaultServeMux.Handle(cmd, handler)
}

// HandleFunc registers the handler function for the given command
// in the DefaultServeMux.
func HandleFunc(cmd string, handler func(Conn, *Message)) {
	DefaultServeMux.HandleFunc(cmd, handler)
}

// ErrorReport returns the ErrorReport channel of the DefaultServeMux.
func ErrorReports() chan ErrorReport {
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
	Addr         string        // TCP address to listen on, ":3868" if empty
	Handler      Handler       // handler to invoke, diam.DefaultServeMux if nil
	Dict         *dict.Parser  // diameter dictionaries for this server
	ReadTimeout  time.Duration // maximum duration before timing out read of the request
	WriteTimeout time.Duration // maximum duration before timing out write of the response
	TLSConfig    *tls.Config   // optional TLS config, used by ListenAndServeTLS
}

// serverHandler delegates to either the server's Handler or DefaultServeMux.
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeDiam(w Conn, m *Message) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	handler.ServeDiam(w, m)
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
// srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":3868"
	}
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("DIAM: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		go c.serve()
	}
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If handler is nil, diam.DefaultServeMux is used.
//
// If dict is nil, dict.Default is used.
func ListenAndServe(addr string, handler Handler, dict *dict.Parser) error {
	server := &Server{Addr: addr, Handler: handler, Dict: dict}
	return server.ListenAndServe()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for
// the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// If srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":3868"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(conn, config)
	return srv.Serve(tlsListener)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects SSL connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler, dict *dict.Parser) error {
	server := &Server{Addr: addr, Handler: handler, Dict: dict}
	return server.ListenAndServeTLS(certFile, keyFile)
}
