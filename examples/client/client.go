// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/fiorix/go-diameter/diam"
)

const (
	Identity    = "client"
	Realm       = "localhost"
	VendorId    = 13
	ProductName = "go-diameter"
)

func main() {
	ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Use: client [-ssl] host:port")
		return
	}
	// ALL incoming messages are handled here.
	diam.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		log.Printf("Receiving message from %s",
			c.RemoteAddr().String())
		m.PrettyPrint()
	})
	// Connect using the default handler and base.Dict.
	addr := os.Args[len(os.Args)-1]
	log.Println("Connecting to", addr)
	var (
		c   diam.Conn
		err error
	)
	if *ssl {
		c, err = diam.DialTLS(addr, "", "", nil, nil)
	} else {
		c, err = diam.Dial(addr, nil, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
	go NewClient(c)
	// Wait until the server kick us out.
	<-c.(diam.CloseNotifier).CloseNotify()
	log.Println("Server disconnected.")
}

// NewClient sends a CER to the server and then a DWR every 10 seconds.
func NewClient(c diam.Conn) {
	// Build CER
	m := diam.NewMessage(
		257,  // CER
		0x80, // Flags: request
		0,    // Application Id
		0,    // HopByHop: 0 means random
		0,    // EndToEnd: 0 means random
		nil,  // nil means diam.dict.Default
	)
	// Add AVPs
	m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
	m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
	IP, _, _ := net.SplitHostPort(c.LocalAddr().String())
	m.NewAVP("Host-IP-Address", 0x40, 0x0, IP)
	m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
	m.NewAVP("Product-Name", 0x40, 0x0, ProductName)
	m.NewAVP("Origin-State-Id", 0x40, 0x0, rand.Uint32())
	log.Printf("Sending message to %s", c.RemoteAddr().String())
	m.PrettyPrint()
	// Send message to the connection
	if _, err := c.Write(m); err != nil {
		log.Println("Write failed:", err)
		return
	}
	// Send watchdog messages every 5 seconds
	for {
		time.Sleep(5 * time.Second)
		m = diam.NewMessage(280, 0x80, 0, 0, 0, nil)
		m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
		m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
		m.NewAVP("Origin-State-Id", 0x40, 0x00, rand.Uint32())
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		m.PrettyPrint()
		if _, err := c.Write(m); err != nil {
			log.Println("Write failed:", err)
			return
		}
	}
}
