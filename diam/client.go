// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package diam

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam/dict"
)

// Dial connects to the peer pointed to by addr and returns the Conn that
// can be used to send diameter messages. Incoming messages are handled
// by the handler, which is typically nil and DefaultServeMux is used.
// If dict is nil, dict.Default is used.
func Dial(addr string, handler Handler, dp *dict.Parser) (Conn, error) {
	srv := &Server{Addr: addr, Handler: handler, Dict: dp}
	return dial(srv, 0)
}

func DialTimeout(addr string, handler Handler, dp *dict.Parser, timeout time.Duration) (Conn, error) {
	srv := &Server{Addr: addr, Handler: handler, Dict: dp}
	return dial(srv, timeout)
}

func dial(srv *Server, timeout time.Duration) (Conn, error) {
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var rw net.Conn
	var err error
	if timeout == 0 {
		rw, err = net.Dial("tcp", addr)
	} else {
		rw, err = net.DialTimeout("tcp", addr, timeout)
	}
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(rw)
	if err != nil {
		return nil, err
	}
	go c.serve()
	return c.writer, nil
}

// DialTLS is the same as Dial, but for TLS.
func DialTLS(addr, certFile, keyFile string, handler Handler, dp *dict.Parser) (Conn, error) {
	srv := &Server{Addr: addr, Handler: handler, Dict: dp}
	return dialTLS(srv, certFile, keyFile, 0)
}

// DialTLSTimeout is the same as DialTimeout, but for TLS.
func DialTLSTimeout(addr, certFile, keyFile string, handler Handler, dp *dict.Parser, timeout time.Duration) (Conn, error) {
	srv := &Server{Addr: addr, Handler: handler, Dict: dp}
	return dialTLS(srv, certFile, keyFile, timeout)
}

func dialTLS(srv *Server, certFile, keyFile string, timeout time.Duration) (Conn, error) {
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var config *tls.Config
	if srv.TLSConfig == nil {
		config = &tls.Config{InsecureSkipVerify: true}
	} else {
		config = TLSConfigClone(srv.TLSConfig)
	}
	if len(certFile) != 0 {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
	}

	var rw net.Conn
	var err error
	if timeout == 0 {
		rw, err = net.Dial("tcp", addr)
	} else {
		rw, err = net.DialTimeout("tcp", addr, timeout)
	}
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(tls.Client(rw, config))
	if err != nil {
		return nil, err
	}
	go c.serve()
	return c.writer, nil
}
