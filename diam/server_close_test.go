// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
)

// TestServerClose verifies Server.Close unblocks Accept, makes Serve return
// ErrServerClosed, and releases the listening port so it can be rebound.
// Regression test for #183.
func TestServerClose(t *testing.T) {
	for _, network := range []string{"tcp", "sctp"} {
		t.Run(network, func(t *testing.T) {
			ln, err := diam.MultistreamListen(network, "127.0.0.1:0")
			if err != nil {
				t.Fatalf("listen: %v", err)
			}
			addr := ln.Addr().String()

			srv := &diam.Server{Network: network}
			errc := make(chan error, 1)
			go func() { errc <- srv.Serve(ln) }()

			time.Sleep(50 * time.Millisecond)

			if err := srv.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			select {
			case err := <-errc:
				if !errors.Is(err, diam.ErrServerClosed) {
					t.Fatalf("Serve returned %v, want ErrServerClosed", err)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Serve did not return after Close")
			}

			ln2, err := diam.MultistreamListen(network, addr)
			if err != nil {
				t.Fatalf("rebind %s %s: %v", network, addr, err)
			}
			ln2.Close()
		})
	}
}

// TestServerCloseBeforeServe ensures Serve returns ErrServerClosed immediately
// when called after Close, without touching the listener.
func TestServerCloseBeforeServe(t *testing.T) {
	srv := &diam.Server{}
	srv.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	if err := srv.Serve(ln); !errors.Is(err, diam.ErrServerClosed) {
		t.Fatalf("Serve after Close returned %v, want ErrServerClosed", err)
	}
}
