// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"crypto/tls"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// fakeMultistreamConn is a MultistreamConn backed by a net.Conn, so that
// newConn takes its multistream branch -- which leaves conn.buf nil -- without
// needing a real SCTP association.
type fakeMultistreamConn struct {
	net.Conn
}

func (c fakeMultistreamConn) ReadAny(b []byte) (int, uint, error) {
	n, err := c.Conn.Read(b)
	return n, 0, err
}

func (c fakeMultistreamConn) ReadStream(b []byte, stream uint) (int, error) {
	return c.Conn.Read(b)
}

func (c fakeMultistreamConn) ReadAtLeast(b []byte, min int, strm uint) (int, uint, error) {
	n, err := c.Conn.Read(b)
	return n, 0, err
}

func (c fakeMultistreamConn) CurrentStream() uint          { return 0 }
func (c fakeMultistreamConn) ResetCurrentStream()          {}
func (c fakeMultistreamConn) SetCurrentStream(s uint) uint { return s }
func (c fakeMultistreamConn) WriteStream(b []byte, s uint) (int, error) {
	return c.Conn.Write(b)
}
func (c fakeMultistreamConn) CurrentWriterStream() uint                  { return 0 }
func (c fakeMultistreamConn) ResetWriterStream()                         {}
func (c fakeMultistreamConn) SetWriterStream(s uint) uint                { return s }
func (c fakeMultistreamConn) SetErrorHandler(MutistreamConnErrorHandler) {}

var _ MultistreamConn = fakeMultistreamConn{}

// A multistream connection has no buffered reader, so it cannot be upgraded.
// It must say so rather than being discovered by a nil dereference inside the
// handshake, after a success CEA has already promised the peer TLS.
func TestCanStartTLSRejectsMultistream(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()

	s := &Server{Handler: NewServeMux(), Dict: dict.Default}
	c, err := s.newConn(fakeMultistreamConn{Conn: srv})
	if err != nil {
		t.Fatal(err)
	}
	if c.buf != nil {
		t.Fatal("precondition: expected a nil buffered reader for a multistream conn")
	}
	w := c.writer
	if w.CanStartTLS() {
		t.Error("CanStartTLS = true for a multistream conn, want false")
	}
	// Must report an error, not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("StartTLS panicked on a multistream conn: %v", r)
		}
	}()
	if err := w.StartTLS(&tls.Config{}); err == nil {
		t.Error("StartTLS = nil for a multistream conn, want an error")
	}
	if err := w.StartTLSClient(&tls.Config{}); err == nil {
		t.Error("StartTLSClient = nil for a multistream conn, want an error")
	}
}

// A connection that already arrived over TLS is not upgradable again.
func TestCanStartTLSRejectsExistingTLS(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()

	s := &Server{Handler: NewServeMux(), Dict: dict.Default}
	c, err := s.newConn(srv)
	if err != nil {
		t.Fatal(err)
	}
	c.inline.Store(true) // as if called from a handler on the serve goroutine
	if !c.writer.CanStartTLS() {
		t.Fatal("precondition: a plain stream conn should be upgradable")
	}
	c.tlsState = &tls.ConnectionState{}
	if c.writer.CanStartTLS() {
		t.Error("CanStartTLS = true for a conn already using TLS, want false")
	}
	if err := c.writer.StartTLS(&tls.Config{}); err == nil {
		t.Error("StartTLS = nil for a conn already using TLS, want an error")
	}
}

// Outside a handler running on the serve goroutine the serve loop may be
// reading, so an upgrade must be refused rather than race it.
func TestCanStartTLSRejectsOffServeGoroutine(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()

	s := &Server{Handler: NewServeMux(), Dict: dict.Default}
	c, err := s.newConn(srv)
	if err != nil {
		t.Fatal(err)
	}
	if c.writer.CanStartTLS() {
		t.Error("CanStartTLS = true outside an inline handler, want false")
	}
	if err := c.writer.StartTLS(&tls.Config{}); err == nil {
		t.Error("StartTLS = nil outside an inline handler, want an error")
	}
}

func TestStartTLSRejectsNilConfig(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()

	s := &Server{Handler: NewServeMux(), Dict: dict.Default}
	c, err := s.newConn(srv)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.writer.StartTLS(nil); err == nil {
		t.Error("StartTLS(nil) = nil, want an error")
	}
}

// inbandHandler is a Handler that opts in to in-band TLS, as sm.StateMachine
// does when Settings.TLSConfig is set.
type inbandHandler struct {
	Handler
	enabled bool
}

func (h inbandHandler) InbandTLSEnabled() bool { return h.enabled }

// Dispatch must be unchanged for a Handler that does not offer in-band TLS,
// including for the capabilities exchange: that is what keeps this feature
// free for everyone who does not use it. A CER is only held on the serve
// goroutine when the Handler opted in and the connection can still upgrade.
func TestDispatchUnaffectedWithoutInbandTLS(t *testing.T) {
	for _, tc := range []struct {
		desc       string
		enabled    *bool // nil: Handler does not implement InbandTLSHandler
		alreadyTLS bool
		wantSync   bool
	}{
		{"handler without the opt-in", nil, false, false},
		{"handler opting out", ptr(false), false, false},
		{"handler opting in", ptr(true), false, true},
		// Nothing to upgrade, so nothing to serialise.
		{"handler opting in, conn already TLS", ptr(true), true, false},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			cli, srv := net.Pipe()
			defer cli.Close()
			defer srv.Close()

			release := make(chan struct{})
			mux := NewServeMux()
			mux.HandleFunc("CER", func(c Conn, m *Message) { <-release })
			var h Handler = mux
			if tc.enabled != nil {
				h = inbandHandler{Handler: mux, enabled: *tc.enabled}
			}
			s := &Server{Handler: h, Dict: dict.Default, MaxConcurrentHandlers: 8}
			c, err := s.newConn(srv)
			if err != nil {
				t.Fatal(err)
			}
			if tc.alreadyTLS {
				c.tlsState = &tls.ConnectionState{}
			}

			// Synchronous dispatch cannot return while the handler is
			// blocked; concurrent dispatch returns straight away.
			returned := make(chan struct{})
			go func() {
				c.dispatch(NewRequest(CapabilitiesExchange, 0, dict.Default))
				close(returned)
			}()
			var gotSync bool
			select {
			case <-returned:
			case <-time.After(200 * time.Millisecond):
				gotSync = true
			}
			close(release)
			<-returned
			c.hwg.Wait()
			if gotSync != tc.wantSync {
				t.Errorf("CER dispatched synchronously = %v, want %v", gotSync, tc.wantSync)
			}
		})
	}
}

func ptr(b bool) *bool { return &b }

// Non-CER commands must never be forced onto the serve goroutine, even with
// in-band TLS enabled.
func TestDispatchConcurrentForNonCER(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	defer srv.Close()

	release := make(chan struct{})
	entered := make(chan struct{})
	mux := NewServeMux()
	mux.HandleFunc("DWR", func(c Conn, m *Message) {
		close(entered)
		<-release
	})
	s := &Server{
		Handler:               inbandHandler{Handler: mux, enabled: true},
		Dict:                  dict.Default,
		MaxConcurrentHandlers: 8,
	}
	c, err := s.newConn(srv)
	if err != nil {
		t.Fatal(err)
	}
	// If DWR were dispatched synchronously this would block forever.
	c.dispatch(NewRequest(DeviceWatchdog, 0, dict.Default))
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("DWR handler did not run concurrently")
	}
	close(release)
	c.hwg.Wait()
}
