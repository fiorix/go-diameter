// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/diamtest"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func settingsWithHandshakeTimeout(d time.Duration) *Settings {
	s := *serverSettings
	s.HandshakeTimeout = d
	return &s
}

// closedWithin reports whether the server closes c within d.
func closedWithin(t *testing.T, c net.Conn, d time.Duration) bool {
	t.Helper()
	c.SetReadDeadline(time.Now().Add(d))
	_, err := c.Read(make([]byte, 1))
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return false
	}
	return err != nil
}

func TestHandshakeTimeoutClosesSilentConn(t *testing.T) {
	srv := diamtest.NewServer(New(settingsWithHandshakeTimeout(200*time.Millisecond)), dict.Default)
	defer srv.Close()
	c, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if !closedWithin(t, c, 2*time.Second) {
		t.Fatal("connection that never sent a CER is still open")
	}
}

// TestHandshakeTimeoutClosesSilentTLSConn covers a peer that connects to a
// TLS listener and never starts the TLS handshake.
func TestHandshakeTimeoutClosesSilentTLSConn(t *testing.T) {
	srv := diamtest.NewUnstartedServer(New(settingsWithHandshakeTimeout(200*time.Millisecond)), dict.Default)
	srv.StartTLS()
	defer srv.Close()
	c, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if !closedWithin(t, c, 2*time.Second) {
		t.Fatal("connection that never started TLS is still open")
	}
}

func TestHandshakeTimeoutKeepsHandshakedConn(t *testing.T) {
	srv := diamtest.NewServer(New(settingsWithHandshakeTimeout(200*time.Millisecond)), dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) { mc <- m })
	mux.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) { mc <- m })
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err = m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code for CEA.\n%s", resp)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("No CEA received")
	}
	time.Sleep(500 * time.Millisecond) // well past the handshake timeout
	m = diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	if _, err = m.WriteTo(cli); err != nil {
		t.Fatalf("handshaked connection closed by the handshake timeout: %v", err)
	}
	select {
	case <-mc:
	case <-time.After(2 * time.Second):
		t.Fatal("No DWA on a handshaked connection after the handshake timeout")
	}
}

func TestHandshakeTimeoutNegativeDisables(t *testing.T) {
	srv := diamtest.NewServer(New(settingsWithHandshakeTimeout(-1)), dict.Default)
	defer srv.Close()
	c, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if closedWithin(t, c, 500*time.Millisecond) {
		t.Fatal("connection closed with the handshake timeout disabled")
	}
}

func TestHandshakeTimeoutDefault(t *testing.T) {
	if got := New(serverSettings).handshakeTimeout(); got != DefaultHandshakeTimeout {
		t.Fatalf("handshakeTimeout() = %v, want DefaultHandshakeTimeout (%v)", got, DefaultHandshakeTimeout)
	}
}

// fakeConn records Close calls; Context always reports no handshake.
type fakeConn struct {
	diam.Conn
	closed int32
}

func (f *fakeConn) Context() context.Context { return context.Background() }
func (f *fakeConn) Close()                   { atomic.AddInt32(&f.closed, 1) }

// TestHandshakeTimeoutStopReleasesTimer checks that the function returned by
// HandleAccept stops the timer: once the connection has closed, the timeout
// must no longer fire.
func TestHandshakeTimeoutStopReleasesTimer(t *testing.T) {
	sm := New(settingsWithHandshakeTimeout(50 * time.Millisecond))

	fired := &fakeConn{}
	sm.HandleAccept(fired)
	stopped := &fakeConn{}
	stop := sm.HandleAccept(stopped)
	stop()

	time.Sleep(200 * time.Millisecond)
	if got := atomic.LoadInt32(&fired.closed); got != 1 {
		t.Errorf("unstopped timer: Close calls = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&stopped.closed); got != 0 {
		t.Errorf("stopped timer: Close calls = %d, want 0", got)
	}
}
