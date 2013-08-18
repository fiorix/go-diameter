// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Experimental diameter server that currently does nothing but print
// incoming messages.
package main

import (
	"fmt"
	"net"

	"github.com/fiorix/go-diameter/diameter"
)

// var dict *diameter.Dict

func main() {
	var err error
	// dict, err = diameter.NewDict("./dict/diam_app.xml")
	// diameter.BaseDict.Load("./dict/diam_app.xml")
	if err != nil {
		panic(err)
	}
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
	msg, err := diameter.ReadMessage(conn, diameter.BaseDict)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("message:\n%s\n", msg.Header)
	for _, avp := range msg.AVP {
		fmt.Printf("  %s\n", avp)
	}
	fmt.Println()
}
