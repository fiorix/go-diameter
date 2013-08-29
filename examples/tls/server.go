// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fiorix/go-diameter/base"
	diam "github.com/fiorix/go-diameter/stack"
)

func init() {
	diam.HandleFunc("CER", OnCER)
}

// On CER reply with CEA.
// http://tools.ietf.org/html/rfc6733#section-5.3.2
func OnCER(c diam.Conn, m *base.Message) {
	log.Printf("Request from %s to %s\n",
		c.RemoteAddr().String(), c.LocalAddr().String())
	m.PrettyPrint()
	// Build answer with Result-Code 2001.
	a := m.Answer(2001)
	a.NewAVP("Origin-Host", 0x40, 0x00,
		&base.DiameterIdentity{base.OctetString{Value: "go"}})
	a.NewAVP("Origin-Realm", 0x40, 0x00,
		&base.DiameterIdentity{base.OctetString{Value: "diameter"}})
	localIP, _, _ := net.SplitHostPort(c.LocalAddr().String())
	a.NewAVP("Host-IP-Address", 0x40, 0x0,
		&base.Address{Family: []byte("01"), IP: net.ParseIP(localIP)})
	a.NewAVP("Vendor-Id", 0x40, 0x0, &base.Unsigned32{Value: 131313})
	a.NewAVP("Product-Name", 0x40, 0x0,
		&base.OctetString{Value: "go-diameter server"})
	// Reply with the same Origin-State-Id
	OriginStateId, err := m.FindAVP("Origin-State-Id")
	if err != nil {
		fmt.Println("Err:", err)
		return
	}
	a.NewAVP("Origin-State-Id", 0x40, 0x0, OriginStateId.Body())
	// Write response
	//fmt.Println("Response:")
	//m.PrettyPrint()
	c.Write(a)
}
