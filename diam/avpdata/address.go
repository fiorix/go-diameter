// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"fmt"
	"net"
)

// Address Diameter Type.
type Address struct {
	// Address Family (e.g. AF_INET=1)
	// http://www.iana.org/assignments/address-family-numbers/address-family-numbers.xhtml
	Family []byte

	// Parsed IP address
	IP net.IP

	// Padding to 4 octets
	Padding int
}

// AF_INET represents IPv4 address family.
var AF_INET = []byte{0, 1}

// Data implements the Data interface.
func (addr *Address) Data() Generic {
	return addr.IP
}

// Put implements the Coded interface. It updates internal Family and IP.
func (addr *Address) Put(b []byte) {
	// TODO: Support IPv6
	if len(b) >= 6 && b[1] == 1 { // AF_INET=1 IPv4 only.
		addr.Family = []byte{b[0], b[1]}
		addr.IP = net.IPv4(b[2], b[3], b[4], b[5])
		addr.Padding = 2
	}
}

// Bytes implement the Codec interface.
func (addr *Address) Bytes() []byte {
	if ip := addr.IP.To4(); ip != nil {
		if addr.Family == nil {
			addr.Family = AF_INET
		}
		// IPv4 always need 2 byte padding. (derived from OctetString)
		addr.Padding = 2 // TODO: Fix this
		b := []byte{
			addr.Family[0],
			addr.Family[1],
			ip[0],
			ip[1],
			ip[2],
			ip[3],
			0, // Padding
			0, // Padding
		}
		return b
	}
	return []byte{}
}

// Length implements the Codec interface. Returns length without padding.
func (addr *Address) Length() uint32 {
	if addr.IP.To4() != nil {
		return 6 // TODO: Fix this
	}
	return 0
}

// String returns a human readable version of the AVP.
func (addr *Address) String() string {
	addr.Bytes() // Update family and padding
	return fmt.Sprintf(
		"Address{IP:'%s',Padding:%d}", addr.IP.String(), addr.Padding)
}
