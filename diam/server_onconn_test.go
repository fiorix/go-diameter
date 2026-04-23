// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam_test

import (
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
)

// TestServerOnNewConnection verifies that Server.OnNewConnection fires once
// per accepted connection with a Conn that can be used with CloseNotify to
// detect disconnection. Regression test for #152.
func TestServerOnNewConnection(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	opened := make(chan diam.Conn, 1)
	srv := &diam.Server{
		OnNewConnection: func(c diam.Conn) { opened <- c },
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	defer srv.Close()

	peer, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	var c diam.Conn
	select {
	case c = <-opened:
	case <-time.After(time.Second):
		t.Fatal("OnNewConnection not called")
	}

	if c.RemoteAddr() == nil {
		t.Fatal("conn RemoteAddr is nil")
	}

	cn, ok := c.(diam.CloseNotifier)
	if !ok {
		t.Fatal("Conn does not implement CloseNotifier")
	}
	closed := cn.CloseNotify()
	peer.Close()

	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("CloseNotify did not fire after peer close")
	}
}
