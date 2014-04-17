// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Address Diameter Type.
type Address net.IP

// DecodeAddress decodes the byte representation of a Diameter Address.
// Example:
// 	b := Address(net.ParseIP("10.0.0.1"))
func DecodeAddress(b []byte) (DataType, error) {
	if len(b) < 6 {
		return nil, errors.New("Not enough data to make an Address")
	}
	switch binary.BigEndian.Uint16(b[:2]) {
	case 0x01:
		if len(b[2:]) != 4 {
			return nil, errors.New("Invalid length for IPv4")
		}
	case 0x02:
		if len(b[2:]) != 16 {
			return nil, errors.New("Invalid length for IPv6")
		}
	default:
		return nil, fmt.Errorf("Unsupported address family: 0x%x", b[:2])
	}
	return Address(b[2:]), nil
}

// Serialize returns the byte representation of the Diameter Address.
// Example:
// 	ip := net.IP(addr.Serialize())
func (addr Address) Serialize() []byte {
	var b []byte
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		b = make([]byte, 6)
		b[1] = 0x01
		copy(b[2:], ip4)
	} else {
		b = make([]byte, 18)
		b[1] = 0x02
		copy(b[2:], addr)

	}
	return b
}

func (addr Address) Len() int {
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		return len(ip4) + 2 // two from address family
	} else {
		return len(addr) + 2
	}
}

func (addr Address) Padding() int {
	l := len(addr) + 2 // two bytes from the address family
	return pad4(l) - l
}

func (addr Address) Type() DataTypeId {
	return AddressType
}

func (addr Address) String() string {
	return fmt.Sprintf("Address{%s},Padding:%d", net.IP(addr), addr.Padding())
}
