// Copyright 2013-2014 go-diameter authors.  All rights reserved.
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
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/avp/format"
)

const (
	Identity    = format.DiameterIdentity("client")
	Realm       = format.DiameterIdentity("localhost")
	VendorId    = format.Unsigned32(13)
	ProductName = format.UTF8String("go-diameter")
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
		log.Printf("Receiving message from %s", c.RemoteAddr().String())
		log.Println(m)
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
	// Passing nil as the last argument to NewRequest will use the default Parser
	// and load dict.Default. To load other dictionaries you can use your own parser:
	// parser, _ := dict.NewParser()
	// parser.Load(bytes.NewReader(dict.DefaultXML))
	// parser.Load(bytes.NewReader(dict.CreditControlXML))
	// m := diam.NewRequest(diam.CapabilitiesExchange, 0, parser)

	// Alternatively you can load more dictionaries into the default parser. e.g.
	// dict.Default.load(bytes.NewReader(dict.CreditControlXML))

	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, format.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorId, avp.Mbit, 0, VendorId)
	m.NewAVP(avp.ProductName, avp.Mbit, 0, ProductName)
	m.NewAVP(avp.OriginStateId, avp.Mbit, 0, format.Unsigned32(rand.Uint32()))
	m.NewAVP(avp.VendorSpecificApplicationId, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationId, avp.Mbit, 0, format.Unsigned32(4)),
			diam.NewAVP(avp.VendorId, avp.Mbit, 0, format.Unsigned32(10415)),
		},
	})
	log.Printf("Sending message to %s", c.RemoteAddr().String())
	log.Println(m)
	// Send message to the connection
	if _, err := m.WriteTo(c); err != nil {
		log.Fatal("Write failed:", err)
	}
	// Send watchdog messages every 5 seconds
	for {
		time.Sleep(5 * time.Second)
		m = diam.NewRequest(diam.DeviceWatchdog, 0, nil)
		m.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
		m.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
		m.NewAVP(avp.OriginStateId, avp.Mbit, 0, format.Unsigned32(rand.Uint32()))
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		log.Println(m)
		if _, err := m.WriteTo(c); err != nil {
			log.Fatal("Write failed:", err)
		}
	}
}
