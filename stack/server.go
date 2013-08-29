// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Forked off Go's HTTP server.

// Diameter server.

package stack

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/base"
	"github.com/fiorix/go-diameter/dict"
)

// Objects implementing the Handler interface can be
// registered to serve a particular path or subtree
// in the diameter server.
//
// ServeDiam should write messages to the ResponseWriter and then return.
// Returning signals that the request is finished and that the Diam server
// can move on to the next request on the connection.
type Handler interface {
	ServeDiam(Conn, *base.Message)
}

// The CloseNotifier interface is implemented by ResponseWriters which
// allow detecting when the underlying connection has gone away.
//
// This mechanism can be used to cancel long operations on the server
// if the client has disconnected before the response is ready.
type CloseNotifier interface {
	// CloseNotify returns a channel that receives a single value
	// when the client connection has gone away.
	CloseNotify() <-chan bool
}

// A Conn interface is used by a handler to construct a
// diameter message.
type Conn interface {
	Write(*base.Message) (int, error) // Writes a msg to the connection
	Close()                           // Close the connection
	LocalAddr() net.Addr              // Local IP
	RemoteAddr() net.Addr             // Remote IP
	TLS() *tls.ConnectionState        // or nil when not using TLS
}

// A conn represents the server side of a diameter connection.
type conn struct {
	server   *Server              // the Server on which the connection arrived
	rwc      net.Conn             // i/o connection
	tlsState *tls.ConnectionState // or nil when not using TLS

	mu           sync.Mutex // guards the following
	clientGone   bool       // if client has disconnected mid-request
	closeNotifyc chan bool  // made lazily
}

func (c *conn) closeNotify() <-chan bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc == nil {
		c.closeNotifyc = make(chan bool, 1)
		_, pw := io.Pipe()
		readSource := c.rwc
		go func() {
			// TODO: test!!!!!
			_, err := io.Copy(pw, readSource)
			if err == nil {
				err = io.EOF
			}
			pw.CloseWithError(err)
			c.noteClientGone()
		}()
	}
	return c.closeNotifyc
}

func (c *conn) noteClientGone() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc != nil && !c.clientGone {
		c.closeNotifyc <- true
	}
	c.clientGone = true
}

// A response represents the server side of a diameter response.
type response struct {
	conn *conn
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.server = srv
	c.rwc = rwc
	return c, nil
}

// Read next request from connection.
func (c *conn) readRequest() (*base.Message, *response, error) {
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}
	dp := c.server.Dict
	if dp == nil {
		dp = base.Dict
	}
	m, err := base.ReadMessage(c.rwc, dp)
	if err != nil {
		return nil, nil, err
	}
	return m, &response{conn: c}, nil
}

func (w *response) Write(m *base.Message) (int, error) {
	return w.conn.rwc.Write(m.Bytes())
}

func (w *response) Close() {
	w.conn.rwc.Close()
}

func (w *response) LocalAddr() net.Addr {
	return w.conn.rwc.LocalAddr()
}

func (w *response) RemoteAddr() net.Addr {
	return w.conn.rwc.RemoteAddr()
}

func (w *response) TLS() *tls.ConnectionState {
	return w.conn.tlsState
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("diam: panic serving %v: %v\n%s",
				c.rwc.RemoteAddr().String(), err, buf)
		}
		c.rwc.Close()
	}()
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if d := c.server.ReadTimeout; d != 0 {
			c.rwc.SetReadDeadline(time.Now().Add(d))
		}
		if d := c.server.WriteTimeout; d != 0 {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			return
		}
		c.tlsState = new(tls.ConnectionState)
		*c.tlsState = tlsConn.ConnectionState()
	}
	for {
		m, w, err := c.readRequest()
		if err != nil {
			log.Print("Server error:", err)
			c.rwc.Close()
			break
		}
		// Diameter cannot have multiple simultaneous active requests.
		// Until the server replies to this request, it can't
		// read another, so we might as well run the handler in
		// this goroutine.
		serverHandler{c.server}.ServeDiam(w, m)
	}
}

func (w *response) CloseNotify() <-chan bool {
	return w.conn.closeNotify()
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as diameter handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(Conn, *base.Message)

// ServeDiam calls f(c, m).
func (f HandlerFunc) ServeDiam(c Conn, m *base.Message) {
	f(c, m)
}

// ServeMux is a diameter request multiplexer.
// It matches the command from the incoming message against a list
// of registered commands and calls the handler.
type ServeMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h   Handler
	cmd string
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// ServeDiam dispatches the request to the handler whose code match
// the incoming message, or close the connection if no handler is found.
func (mux *ServeMux) ServeDiam(c Conn, m *base.Message) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	var cmd string
	if dcmd, err := m.Dict.FindCmd(
		m.Header.ApplicationId,
		m.Header.CommandCode(),
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
		panic("diam: nil handler")
	}
	mux.m[cmd] = muxEntry{h: handler, cmd: cmd}
}

// HandleFunc registers the handler function for the given pattern.
// Special @cmd "ALL" may be used as a catch all.
func (mux *ServeMux) HandleFunc(cmd string, handler func(Conn, *base.Message)) {
	mux.Handle(cmd, HandlerFunc(handler))
}

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func Handle(cmd string, handler Handler) { DefaultServeMux.Handle(cmd, handler) }

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(cmd string, handler func(Conn, *base.Message)) {
	DefaultServeMux.HandleFunc(cmd, handler)
}

// Serve accepts incoming HTTP connections on the listener l,
// creating a new service goroutine for each.  The service goroutines
// read requests and then call handler to reply to them.
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

func (sh serverHandler) ServeDiam(w Conn, m *base.Message) {
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
				log.Printf("diam: Accept error: %v; retrying in %v", e, tempDelay)
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
	return nil
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If @handler is nil, diam.DefaultServeMux is used.
//
// If @dict is nil, dict.Base is used.
func ListenAndServe(addr string, handler Handler, dict *dict.Parser) error {
	server := &Server{Addr: addr, Handler: handler, Dict: dict}
	return server.ListenAndServe()
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

// TimeoutHandler returns a Handler that runs h with the given time limit.
//
// The new Handler calls h.ServeDiam to handle each request, but if a
// call runs for longer than its time limit, the connection is closed.
// After such a timeout its MessageWriter will return ErrHandlerTimeout.
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
	f := func() <-chan time.Time {
		return time.After(dt)
	}
	return &timeoutHandler{h, f, msg}
}

// ErrHandlerTimeout is returned on MessageWriter Write calls
// in handlers which have timed out.
var ErrHandlerTimeout = errors.New("diam: Handler timeout")

type timeoutHandler struct {
	handler Handler
	timeout func() <-chan time.Time // returns channel producing a timeout
	body    string
}

func (h *timeoutHandler) ServeDiam(w Conn, m *base.Message) {
	done := make(chan bool, 1)
	tw := &timeoutConn{w: w}
	go func() {
		h.handler.ServeDiam(tw, m)
		done <- true
	}()
	select {
	case <-done:
		return
	case <-h.timeout():
		tw.mu.Lock()
		defer tw.mu.Unlock()
		if !tw.wroteHeader {
			// TODO: Retransmit? Close?
		}
		tw.timedOut = true
	}
}

type timeoutConn struct {
	w Conn

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
}

func (tw *timeoutConn) Write(m *base.Message) (int, error) {
	return tw.w.Write(m)
}

func (tw *timeoutConn) Close() {
	tw.w.Close()
}

func (tw *timeoutConn) LocalAddr() net.Addr {
	return tw.w.LocalAddr()
}

func (tw *timeoutConn) RemoteAddr() net.Addr {
	return tw.w.RemoteAddr()
}

func (tw *timeoutConn) TLS() *tls.ConnectionState {
	return tw.w.TLS()
}
