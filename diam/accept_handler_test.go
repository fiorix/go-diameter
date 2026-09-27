// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

type countingAcceptHandler struct {
	HandlerFunc
	accepts int32
	closes  int32
}

func (h *countingAcceptHandler) HandleAccept(Conn) func() {
	atomic.AddInt32(&h.accepts, 1)
	return func() { atomic.AddInt32(&h.closes, 1) }
}

func TestAcceptHandlerOnlyForAcceptedConns(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srvH := &countingAcceptHandler{HandlerFunc: func(Conn, *Message) {}}
	srv := &Server{Handler: srvH, Dict: dict.Default}
	go srv.Serve(ln)
	defer srv.Close()

	cliH := &countingAcceptHandler{HandlerFunc: func(Conn, *Message) {}}
	cli, err := Dial(ln.Addr().String(), cliH, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadInt32(&srvH.accepts) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := atomic.LoadInt32(&srvH.accepts); got != 1 {
		t.Errorf("server HandleAccept calls = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&cliH.accepts); got != 0 {
		t.Errorf("dialed connection called HandleAccept %d times, want 0", got)
	}

	cli.Close()
	deadline = time.Now().Add(2 * time.Second)
	for atomic.LoadInt32(&srvH.closes) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := atomic.LoadInt32(&srvH.closes); got != 1 {
		t.Errorf("onClose calls after the peer closed = %d, want 1", got)
	}
}
