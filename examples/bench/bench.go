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
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diamtype"
)

const (
	Identity    = diamtype.DiameterIdentity("client")
	Realm       = diamtype.DiameterIdentity("localhost")
	VendorId    = diamtype.Unsigned32(13)
	ProductName = diamtype.UTF8String("go-diameter")
)

var (
	BenchMessages int
	wg            sync.WaitGroup
)

func main() {
	ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
	c := flag.Int("c", 1, "concurrent clients")
	n := flag.Int("n", 1000, "messages per client")
	flag.Parse()
	BenchMessages = *n
	if len(os.Args) < 2 {
		fmt.Println("Use: bench [options] host:port")
		flag.Usage()
		return
	}
	// Use all CPUs.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Launch clients.
	for n := 0; n < *c; n++ {
		wg.Add(1)
		go NewClient(*ssl)
	}
	wg.Wait()
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
		if counter == BenchMessages {
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
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, diamtype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorId, avp.Mbit, 0, VendorId)
	m.NewAVP(avp.ProductName, avp.Mbit, 0, ProductName)
	m.NewAVP(avp.OriginStateId, avp.Mbit, 0, diamtype.Unsigned32(rand.Uint32()))
	// Send message to the connection
	if _, err := m.WriteTo(c); err != nil {
		log.Println("Write failed:", err)
		return
	}
	// Prepare the ACR that is used for benchmarking.
	m = diam.NewRequest(diam.Accounting, 0, nil)
	m.NewAVP(avp.SessionId, avp.Mbit, 0, diamtype.UTF8String("Hello"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, diamtype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorId, avp.Mbit, 0, VendorId)
	m.NewAVP(avp.ProductName, avp.Mbit, 0, ProductName)
	m.NewAVP(avp.OriginStateId, avp.Mbit, 0, diamtype.Unsigned32(rand.Uint32()))
	log.Println("OK, sending messages")
	var n int
	go func() {
		start := time.Now()
		// Wait until the connection is closed.
		<-c.(diam.CloseNotifier).CloseNotify()
		elapsed := time.Since(start)
		mps := int(float64(n) / elapsed.Seconds())
		log.Printf("%d messages in %s seconds, %d msg/s",
			n, elapsed, mps)
		wg.Done()
	}()
	for ; n < BenchMessages; n++ {
		// Send message to the connection
		if _, err := m.WriteTo(c); err != nil {
			log.Fatal("Write failed:", err)
		}
	}
}
