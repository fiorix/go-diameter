// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Experimental diameter server that currently does nothing but print
// incoming messages.
package main

import (
	"fmt"
	"net"

	"github.com/fiorix/go-diameter/base"
	"github.com/fiorix/go-diameter/dict"
)

var GlobalDict *dict.Parser

func main() {
	var err error
	if GlobalDict, err = dict.New("./diam_app.xml"); err != nil {
		panic(err)
	}
	GlobalDict.PrettyPrint()
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
	msg, err := base.ReadMessage(conn, GlobalDict)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	msg.PrettyPrint()
}
