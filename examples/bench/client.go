// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter bench client skeleton.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/diam"
)

const (
	Identity    = "client"
	Realm       = "localhost"
	VendorId    = 13
	ProductName = "go-diameter"
)

var (
	BenchMessages int
	WG            sync.WaitGroup
)

func main() {
	ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
	c := flag.Int("c", 1, "concurrent clients")
	n := flag.Int("n", 1000, "messages per client")
	flag.Parse()
	BenchMessages = *n
	if len(os.Args) < 2 {
		fmt.Println("Use: client [options] host:port")
		flag.Usage()
		return
	}
	// Launch clients.
	for n := 0; n < *c; n++ {
		WG.Add(1)
		go NewClient(*ssl)
	}
	WG.Wait()
	log.Println("Done.")
}

// NewClient connects to the server and sends a CER.
// In parallel, a DWR is sent every 10 seconds unless
func NewClient(ssl bool) {
	addr := os.Args[len(os.Args)-1]
	log.Println("Connecting to", addr)
	var (
		c   diam.Conn
		err error
	)
	// Handle ACA incoming messages.
	var counter int
	mux := diam.NewServeMux()
	mux.HandleFunc("ACA", func(c diam.Conn, m *diam.Message) {
		counter++
		OnACA(c, m, counter)
	})
	// Connect using the our custom mux and default base.Dict.
	if ssl {
		c, err = diam.DialTLS(addr, "", "", mux, nil)
	} else {
		c, err = diam.Dial(addr, mux, nil)
	}
	if err != nil {
		log.Fatal(err)
		WG.Done()
	}
	go func() {
		// Wait until the connection is closed.
		<-c.(diam.CloseNotifier).CloseNotify()
		WG.Done()
	}()
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
	// Send message to the connection
	if _, err := c.Write(m); err != nil {
		log.Println("Write failed:", err)
		return
	}
	// Prepare the ACR that is used for benchmarking.
	m = diam.NewMessage(
		271,  // ACR
		0x80, // Flags: request
		0,    // Application Id
		0,    // HopByHop: 0 means random
		0,    // EndToEnd: 0 means random
		nil,  // nil means diam.dict.Default
	)
	// Add AVPs
	SessId, _ := m.NewAVP("Session-Id", 0x40, 0x00, "Hello")
	m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
	m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
	m.NewAVP("Host-IP-Address", 0x40, 0x0, IP)
	m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
	m.NewAVP("Product-Name", 0x40, 0x0, ProductName)
	StateId, _ := m.NewAVP("Origin-State-Id", 0x40, 0x0, rand.Uint32())
	log.Println("OK, sending messages")
	start := time.Now()
	var n int
	for ; n < BenchMessages; n++ {
		// Send message to the connection
		if _, err := c.Write(m); err != nil {
			log.Println("Write failed:", err)
			break
		}
		SessId.Set(fmt.Sprintf("%d", rand.Uint32()))
		StateId.Set(rand.Uint32())
	}
	elapsed := time.Since(start)
	mps := int(float64(n) / elapsed.Seconds())
	log.Printf("%d messages in %s seconds, %d msg/s",
		n, elapsed, mps)
}

// OnACA handles all incoming messages Accounting-Answer messages.
func OnACA(c diam.Conn, m *diam.Message, counter int) {
	if counter == BenchMessages {
		c.Close()
	}
}
