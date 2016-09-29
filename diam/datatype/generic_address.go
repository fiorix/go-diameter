// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
	"net"
)

// GenericAddress data type
type GenericAddress struct {
	AddrType  uint16
	AddrValue []byte
}

// DecodeGenericAddress is not used anywhere, it is here so that GenericAddress
// satisfies Type interface
func DecodeGenericAddress(b []byte) (Type, error) {
	return GenericAddress{}, nil
}

// Serialize implements the Type interface.
func (addr GenericAddress) Serialize() []byte {
	var b []byte
	switch addr.AddrType {
	case 0x01:
		b = make([]byte, 6)
		binary.BigEndian.PutUint16(b[:2], uint16(addr.AddrType))
		if ip4 := net.IP(addr.AddrValue).To4(); ip4 != nil {
			copy(b[2:], ip4)
		} else {
			fmt.Errorf("Error during serialization: not an IPv4 address")
		}
	case 0x02:
		b = make([]byte, 18)
		binary.BigEndian.PutUint16(b[:2], uint16(addr.AddrType))
		if ip16 := net.IP(addr.AddrValue).To16(); ip16 != nil {
			copy(b[2:], ip16)
		} else {
			fmt.Errorf("Error during serialization: not an IPv6 address")
		}
	default:
		b = make([]byte, (len(addr.AddrValue) + 2))
		binary.BigEndian.PutUint16(b[:2], uint16(addr.AddrType))
		copy(b[2:], addr.AddrValue)
	}
	return b
}

// Len implements the Type interface.
func (addr GenericAddress) Len() int {
	return len(addr.AddrValue) + 2
}

// Padding implements the Type interface.
func (addr GenericAddress) Padding() int {
	l := len(addr.AddrValue) + 2 // two bytes from the address family
	return pad4(l) - l
}

// Type implements the Type interface.
func (addr GenericAddress) Type() TypeID {
	return AddressType
}

// String implements the Type interface.
func (addr GenericAddress) String() string {
	return fmt.Sprintf("Address{%x},Type{%d},Padding:%d", addr.AddrValue, addr.AddrType, addr.Padding())
}
