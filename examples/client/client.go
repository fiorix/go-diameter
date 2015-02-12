// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client example. This is by no means a complete client.
// The commands in here are not fully implemented. For that you have
// to read the RFCs (base and credit control) and follow the spec.

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
	"github.com/fiorix/go-diameter/diam/datatype"
)

const (
	identity    = datatype.DiameterIdentity("client")
	realm       = datatype.DiameterIdentity("localhost")
	vendorID    = datatype.Unsigned32(13)
	productName = datatype.UTF8String("go-diameter")
)

func main() {
	ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Use: client [-ssl] host:port")
		return
	}
	// ALL incoming messages are handled here.
	sessionID := "fake_session_id"
	msisdn := "85589481811"
	diam.Handle("CEA", OnCEA(sessionID, msisdn))
	diam.HandleFunc("CCA", OnCCA)
	diam.HandleFunc("ALL", OnMSG) // Catch-all.
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
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	m.NewAVP(avp.ProductName, 0, 0, productName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(rand.Uint32()))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
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
		m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
		m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(rand.Uint32()))
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		log.Println(m)
		if _, err := m.WriteTo(c); err != nil {
			log.Fatal("Write failed:", err)
		}
	}
}

// OnCEA handles Capabilities-Exchange-Answer messages.
func OnCEA(sessionID string, msisdn string) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		rc, err := m.FindAVP(avp.ResultCode)
		if err != nil {
			log.Fatal(err)
		}
		if v, _ := rc.Data.(datatype.Unsigned32); v != diam.Success {
			log.Fatal("Unexpected response:", rc)
		}
		// Craft a CCR message.
		r := diam.NewRequest(diam.CreditControl, 4, nil)
		r.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionID))
		r.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
		r.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
		peerRealm, _ := m.FindAVP(avp.OriginRealm) // You should handle errors.
		r.NewAVP(avp.DestinationRealm, avp.Mbit, 0, peerRealm.Data)
		r.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
		r.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0x00)),
				diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(msisdn)),
			},
		})
		// Add Service-Context-Id and all other AVPs...
		r.WriteTo(c)
	}
}

// OnCCA handles Credit-Control-Answer messages.
func OnCCA(c diam.Conn, m *diam.Message) {
	log.Println(m)
}

// OnMSG handles all other messages and just print them.
func OnMSG(c diam.Conn, m *diam.Message) {
	log.Printf("Receiving message from %s", c.RemoteAddr().String())
	log.Println(m)
}
