// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// TestWriteErrorClosesConn checks that a failed write ends the connection.
// The bufio.Writer behind Conn.Write keeps returning its first error, so a
// connection left open after a write timeout can never answer again while
// it keeps reading and processing requests.
func TestWriteErrorClosesConn(t *testing.T) {
	writeErr := make(chan error, 1)
	mux := NewServeMux()
	mux.HandleFunc("CER", func(c Conn, m *Message) {
		a := m.Answer(Success)
		a.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String(string(make([]byte, 60000))))
		var err error
		for i := 0; i < 2000 && err == nil; i++ {
			_, err = a.WriteTo(c)
		}
		writeErr <- err
	})
	srv := &Server{Handler: mux, Dict: dict.Default, WriteTimeout: 200 * time.Millisecond}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go srv.Serve(l)

	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	cer := NewRequest(CapabilitiesExchange, 0, dict.Default)
	cer.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("peer"))
	if _, err := cer.WriteTo(c); err != nil {
		t.Fatal(err)
	}
	// Do not read: the server's writes fill the socket buffers and time out.
	select {
	case err := <-writeErr:
		if err == nil {
			t.Fatal("server writes never failed; raise the write count")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("server write did not time out")
	}
	// Draining must end with EOF: the server closed the connection.
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.Copy(io.Discard, c); err != nil {
		t.Fatalf("connection still open after a failed write: read error %v, want EOF", err)
	}
}
