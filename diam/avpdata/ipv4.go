// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"fmt"
	"net"
)

// IPv4 Type for Framed-IP-Address and alike.
type IPv4 struct {
	// Parsed IP address
	IP net.IP
}

// Data implements the Data interface.
func (addr *IPv4) Data() Generic {
	return addr.IP
}

// Put implements the Coded interface. It updates internal Family and IP.
func (addr *IPv4) Put(b []byte) {
	if len(b) == 4 {
		addr.IP = net.IPv4(b[0], b[1], b[2], b[3])
	}
}

// Bytes implement the Codec interface.
func (addr *IPv4) Bytes() []byte {
	if ip := addr.IP.To4(); ip != nil {
		return ip
	}
	return []byte{}
}

// Length implements the Codec interface. Returns length without padding.
func (addr *IPv4) Length() uint32 {
	if addr.IP.To4() != nil {
		return 4 // TODO: Fix this
	}
	return 0
}

// String returns a human readable version of the AVP.
func (addr *IPv4) String() string {
	return fmt.Sprintf("IPv4{Value:'%s'}", addr.IP.String())
}
