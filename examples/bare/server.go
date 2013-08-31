// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Experimental diameter server that currently does nothing but print
// incoming messages.
package main

import (
	"fmt"
	"net"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/dict"
)

func main() {
	dict.Default.LoadFile("diam_app.xml")
	dict.Default.PrettyPrint()
	fmt.Println()
	srv, err := net.Listen("tcp", ":3868")
	if err != nil {
		panic(err)
	}
	for {
		if conn, err := srv.Accept(); err != nil {
			panic(err)
		} else {
			go client(conn)
		}
	}
}

func client(conn net.Conn) {
	defer conn.Close()
	msg, err := diam.ReadMessage(conn, dict.Default)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	msg.PrettyPrint()
}
