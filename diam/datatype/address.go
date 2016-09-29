// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Address data type.
type Address []byte

//type Address net.IP

// DecodeAddress decodes an Address data type from byte array.
func DecodeAddress(b []byte) (Type, error) {
	if len(b) < 3 {
		return nil, errors.New("Not enough data to make an Address")
	}
	if binary.BigEndian.Uint16(b[:2]) <= 0 || binary.BigEndian.Uint16(b[:2]) >= 65535 {
		return nil, errors.New("Invalid address type received")
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
		// Do nothing, we are not validating every possible address type
	}
	return Address(b[2:]), nil
}

// Serialize implements the Type interface.
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

// Len implements the Type interface.
func (addr Address) Len() int {
	ip4 := net.IP(addr).To4()
	if ip4 != nil {
		return len(ip4) + 2 // Two from address family.
	}
	return len(addr) + 2
}

// Padding implements the Type interface.
func (addr Address) Padding() int {
	l := len(addr) + 2 // two bytes from the address family
	return pad4(l) - l
}

// Type implements the Type interface.
func (addr Address) Type() TypeID {
	return AddressType
}

// String implements the Type interface.
func (addr Address) String() string {
	return fmt.Sprintf("Address{%s}, Padding:%d", string(addr), addr.Padding())
}
