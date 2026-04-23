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

// TestServerDialWriteTimeout verifies that Server.WriteTimeout is honored on
// the client side when dialing via (*Server).Dial. Regression test for #218:
// client writes to a stalled peer should unblock with a timeout error instead
// of piling up on the write mutex.
func TestServerDialWriteTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	accepted := make(chan net.Conn, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		if tc, ok := c.(*net.TCPConn); ok {
			tc.SetReadBuffer(1024)
		}
		accepted <- c
	}()

	srv := &diam.Server{
		Network:      "tcp",
		Addr:         ln.Addr().String(),
		WriteTimeout: 100 * time.Millisecond,
	}
	cli, err := srv.Dial(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	var peer net.Conn
	select {
	case peer = <-accepted:
	case <-time.After(time.Second):
		t.Fatal("peer did not accept")
	}
	defer peer.Close()

	if tc, ok := cli.Connection().(*net.TCPConn); ok {
		tc.SetWriteBuffer(1024)
	}

	payload := make([]byte, 64*1024)
	deadline := time.After(5 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("Write never timed out; WriteTimeout not honored on client")
		default:
		}
		_, err := cli.Write(payload)
		if err == nil {
			continue
		}
		ne, ok := err.(net.Error)
		if !ok || !ne.Timeout() {
			t.Fatalf("expected timeout error, got %v", err)
		}
		return
	}
}
