// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"net"

	"os"

	"github.com/fiorix/go-diameter/base"
	diam "github.com/fiorix/go-diameter/stack"
)

func init() {
	diam.HandleFunc("CEA", OnCEA)
}

func SendCER(c diam.Conn) {
	// Build CER message.
	m := base.NewMessage(257, 0x80, 0, 0, 0, base.Dict)
	// AVPs.
	m.NewAVP("Origin-Host", 0x40, 0x00,
		&base.DiameterIdentity{base.OctetString{Value: "go"}})
	m.NewAVP("Origin-Realm", 0x40, 0x00,
		&base.DiameterIdentity{base.OctetString{Value: "diameter"}})
	localIP, _, _ := net.SplitHostPort(c.LocalAddr().String())
	m.NewAVP("Host-IP-Address", 0x40, 0x0,
		&base.Address{Family: []byte("01"), IP: net.ParseIP(localIP)})
	m.NewAVP("Vendor-Id", 0x40, 0x0, &base.Unsigned32{Value: 131313})
	m.NewAVP("Product-Name", 0x40, 0x0,
		&base.OctetString{Value: "go-diameter client"})
	// Add random Origin-State-Id
	m.NewAVP("Origin-State-Id", 0x40, 0x0,
		&base.Unsigned32{Value: rand.Uint32()})
	// Write request
	//fmt.Println("Request:")
	//m.PrettyPrint()
	c.Write(m)
}

func OnCEA(c diam.Conn, m *base.Message) {
	log.Printf("Response from %s to %s\n",
		c.RemoteAddr().String(), c.LocalAddr().String())
	m.PrettyPrint()
	// Bye bye
	os.Exit(0)
}
