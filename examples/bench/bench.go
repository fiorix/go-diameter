// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter benchmark client.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

const (
	identity    = datatype.DiameterIdentity("client")
	realm       = datatype.DiameterIdentity("localhost")
	vendorID    = datatype.Unsigned32(13)
	productName = datatype.UTF8String("go-diameter")
)

var (
	benchMessages int
	totalMessages int32
	wg            sync.WaitGroup
)

func main() {
	ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
	c := flag.Int("c", 1, "concurrent clients")
	n := flag.Int("n", 1000, "messages per client")
	t := flag.Int("t", 0, "threads (0 means one per core)")
	flag.Parse()
	benchMessages = *n
	if len(os.Args) < 2 {
		fmt.Println("Use: bench [options] host:port")
		flag.Usage()
		return
	}
	if *t == 0 {
		*t = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*t)
	// Launch clients.
	start := time.Now()
	for n := 0; n < *c; n++ {
		wg.Add(1)
		go NewClient(*ssl)
	}
	wg.Wait()
	elapsed := time.Since(start)
	mps := int(float64(totalMessages) / elapsed.Seconds())
	log.Println("Done")
	log.Printf("Total of %d messages in %s: %d msg/s",
		totalMessages, elapsed, mps)
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
		if counter == benchMessages {
			c.Close()
		}
	})
	// Connect using the our custom mux and default dictionary.
	if ssl {
		c, err = diam.DialTLS(addr, "", "", mux, nil)
	} else {
		c, err = diam.Dial(addr, mux, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
	// Build CER
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	m.NewAVP(avp.ProductName, 0, 0, productName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(rand.Uint32()))
	// Send message to the connection
	if _, err := m.WriteTo(c); err != nil {
		log.Println("Write failed:", err)
		return
	}
	// Prepare the ACR that is used for benchmarking.
	m = diam.NewRequest(diam.Accounting, 0, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("Hello"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	m.NewAVP(avp.ProductName, 0, 0, productName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(rand.Uint32()))
	log.Println("OK, sending messages")
	var n int
	go func() {
		start := time.Now()
		// Wait until the connection is closed.
		<-c.(diam.CloseNotifier).CloseNotify()
		n *= 2 // request+answer
		elapsed := time.Since(start)
		mps := int(float64(n) / elapsed.Seconds())
		log.Printf("%d messages (request+answer) in %s seconds, %d msg/s",
			n, elapsed, mps)
		atomic.AddInt32(&totalMessages, int32(n))
		wg.Done()
	}()
	for ; n < benchMessages; n++ {
		// Send message to the connection
		if _, err := m.WriteTo(c); err != nil {
			log.Fatal("Write failed:", err)
		}
	}
}
