// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
// The commands in here are not fully implemented. For that you have
// to read the RFCs (base and credit control) and follow the spec.

// Generate SSL certificates:
// go run $GOROOT/src/pkg/crypto/tls/generate_cert.go --host localhost

package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

const (
	identity    = datatype.DiameterIdentity("server")
	realm       = datatype.DiameterIdentity("localhost")
	vendorID    = datatype.Unsigned32(13)
	productName = datatype.UTF8String("go-diameter")
)

var quiet bool

func main() {
	addr := flag.String("l", ":3868", "listen address and port")
	cert := flag.String("cert", "", "SSL cert file (e.g. cert.pem)")
	key := flag.String("key", "", "SSL key file (e.g. key.pem)")
	q := flag.Bool("q", false, "quiet, do not print messages")
	t := flag.Int("t", 0, "threads (0 means one per core)")
	flag.Parse()
	quiet = *q
	if *t == 0 {
		*t = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*t)
	// Message handlers:
	diam.HandleFunc("CER", OnCER)
	diam.HandleFunc("CCR", OnCCR)
	diam.HandleFunc("ALL", OnMSG) // Catch-all
	// Handle server errors.
	go func() {
		for {
			report := <-diam.ErrorReports()
			log.Println("Error:", report.Error)
		}
	}()
	// Start server.
	if *cert != "" && *key != "" {
		log.Println("Starting secure server on", *addr)
		log.Fatal(diam.ListenAndServeTLS(*addr, *cert, *key, nil, nil))
	} else {
		log.Println("Starting server on", *addr)
		log.Fatal(diam.ListenAndServe(*addr, nil, nil))
	}
}

// OnCER handles Capabilities-Exchange-Request messages.
func OnCER(c diam.Conn, m *diam.Message) {
	// Reject client if there's no Origin-Host.
	host, err := m.FindAVP(avp.OriginHost)
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Origin-Realm.
	crealm, err := m.FindAVP(avp.OriginRealm)
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Host-IP-Address.
	ipaddr, err := m.FindAVP(avp.HostIPAddress)
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Origin-State-Id.
	stateID, err := m.FindAVP(avp.OriginStateID)
	if err != nil {
		c.Close()
		return
	}
	if !quiet {
		//log.Println("Receiving message from %s", c.RemoteAddr().String())
		log.Printf("Receiving message from %s.%s (%s)", host, crealm, ipaddr)
		log.Println(m)
	}
	// Craft CEA with result code 2001 (OK).
	a := m.Answer(diam.Success)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	a.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	a.NewAVP(avp.ProductName, avp.Mbit, 0, productName)
	// Copy origin Origin-State-Id.
	a.AddAVP(stateID)
	if !quiet {
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		log.Println(a)
	}
	// Send message to the connection
	if _, err := a.WriteTo(c); err != nil {
		log.Println("Write failed:", err)
		c.Close()
	}
	go func() {
		<-c.(diam.CloseNotifier).CloseNotify()
		if !quiet {
			log.Printf("Client %s disconnected",
				c.RemoteAddr().String())
		}
	}()
}

// OnCCR handles Credit-Control-Request messages.
func OnCCR(c diam.Conn, m *diam.Message) {
	log.Println(m)
	// Craft a CCA with result code 3001.
	a := m.Answer(diam.CommandUnsupported)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	a.WriteTo(c)
}

// OnMSG handles all other messages and replies to them
// with a generic 2001 (OK) answer.
func OnMSG(c diam.Conn, m *diam.Message) {
	// Ignore message if there's no Origin-State-Id.
	stateID, err := m.FindAVP(avp.OriginStateID)
	if err != nil {
		log.Println("Invalid message: missing Origin-State-Id\n", m)
	}
	if !quiet {
		log.Printf("Receiving message from %s", c.RemoteAddr().String())
		log.Println(m)
	}
	// Craft answer with result code 2001 (OK).
	a := m.Answer(diam.Success)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	a.AddAVP(stateID)
	if !quiet {
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		log.Println(a)
	}
	// Send message to the connection.
	if _, err := a.WriteTo(c); err != nil {
		log.Println("Write failed:", err)
		c.Close()
	}
}
