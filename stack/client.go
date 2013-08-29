// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//

// Diameter client.

package stack

import (
	"crypto/tls"
	"net"

	"github.com/fiorix/go-diameter/dict"
)

func Dial(addr string, handler Handler, dict *dict.Parser) (Conn, error) {
	server := &Server{Addr: addr, Handler: handler, Dict: dict}
	return server.Dial()
}

func (srv *Server) Dial() (Conn, error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":3868"
	}
	rw, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(rw)
	if err != nil {
		return nil, err
	}
	go c.serve()
	return &response{conn: c}, nil
}

func DialTLS(addr, certFile, keyFile string, handler Handler, dict *dict.Parser) (Conn, error) {
	server := &Server{Addr: addr, Handler: handler, Dict: dict}
	return server.DialTLS(certFile, keyFile)
}

func (srv *Server) DialTLS(certFile, keyFile string) (Conn, error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":3868"
	}
	config := &tls.Config{InsecureSkipVerify: true}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if certFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
	}
	rw, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c, err := srv.newConn(tls.Client(rw, config))
	if err != nil {
		return nil, err
	}
	go c.serve()
	return &response{conn: c}, nil
}
