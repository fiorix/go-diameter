// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server.

// Generate SSL certificates:
// go run $GOROOT/src/pkg/crypto/tls/generate_cert.go --host localhost

package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/fiorix/go-diameter/diam"
)

const (
	Identity    = "server"
	Realm       = "localhost"
	VendorId    = 13
	ProductName = "go-diameter"
)

var Quiet bool

func main() {
	addr := flag.String("l", ":3868", "listen address and port")
	cert := flag.String("cert", "", "SSL cert file (e.g. cert.pem)")
	key := flag.String("key", "", "SSL key file (e.g. key.pem)")
	q := flag.Bool("q", false, "quiet, do not print messages")
	flag.Parse()
	Quiet = *q
	// Use all CPUs.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Message handlers:
	diam.HandleFunc("CER", OnCER)
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
	host, err := m.FindAVP("Origin-Host")
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Origin-Realm.
	realm, err := m.FindAVP("Origin-Realm")
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Host-IP-Address.
	ipaddr, err := m.FindAVP("Host-IP-Address")
	if err != nil {
		c.Close()
		return
	}
	// Reject client if there's no Origin-State-Id.
	stateId, err := m.FindAVP("Origin-State-Id")
	if err != nil {
		c.Close()
		return
	}
	if !Quiet {
		//log.Println("Receiving message from %s", c.RemoteAddr().String())
		log.Printf("Receiving message from %s.%s (%s)",
			host.Data().(string),
			realm.Data().(string),
			ipaddr.Data().(net.IP).String(),
		)
		m.PrettyPrint()
	}
	// Craft CEA with result code 2001 (OK).
	a := m.Answer(2001)
	a.NewAVP("Origin-Host", 0x40, 0x00, Identity)
	a.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
	IP, _, _ := net.SplitHostPort(c.LocalAddr().String())
	a.NewAVP("Host-IP-Address", 0x40, 0x0, IP)
	a.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
	a.NewAVP("Product-Name", 0x40, 0x0, ProductName)
	// Copy origin Origin-State-Id.
	a.Add(stateId)
	if !Quiet {
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		a.PrettyPrint()
	}
	// Send message to the connection
	if _, err := c.Write(a); err != nil {
		log.Println("Write failed:", err)
		c.Close()
	}
	go func() {
		<-c.(diam.CloseNotifier).CloseNotify()
		if !Quiet {
			log.Printf("Client %s disconnected",
				c.RemoteAddr().String())
		}
	}()
}

// OnMSG handles all other messages and replies to them
// with a generic 2001 (OK) answer.
func OnMSG(c diam.Conn, m *diam.Message) {
	// Ignored message if there's no Origin-State-Id.
	stateId, err := m.FindAVP("Origin-State-Id")
	if err != nil {
		return
	}
	if !Quiet {
		log.Printf(
			"Receiving message from %s", c.RemoteAddr().String())
		m.PrettyPrint()
	}
	// Craft DWA with result code 2001 (OK).
	a := m.Answer(2001)
	a.NewAVP("Origin-Host", 0x40, 0x00, Identity)
	a.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
	a.Add(stateId)
	if !Quiet {
		log.Printf("Sending message to %s", c.RemoteAddr().String())
		a.PrettyPrint()
	}
	// Send message to the connection
	if _, err := c.Write(a); err != nil {
		log.Println("Write failed:", err)
		c.Close()
	}
}
