// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"net"
	"testing"
)

func TestGenericAddressIPv4(t *testing.T) {
	address := GenericAddress{AddrType: 0x01, AddrValue: net.ParseIP("10.0.0.1").To4()}
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 6 {
		t.Fatalf("Unexpected len. Want 6, have %d", address.Len())
	}
	if len(address.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestGenericAddressIPv6(t *testing.T) {
	address := GenericAddress{AddrType: 0x02, AddrValue: net.ParseIP("2001:0db8::ff00:0042:8329").To16()}
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 18 {
		t.Fatalf("Unexpected len. Want 18, have %d", address.Len())
	}
	//t.Log(address)
}

func TestGenericAddressE164(t *testing.T) {
	address := GenericAddress{AddrType: 0x08, AddrValue: []byte("48602007060")}
	b := []byte{0x00, 0x08, 0x34, 0x38, 0x36, 0x30, 0x32, 0x30, 0x30, 0x37, 0x30, 0x36, 0x30}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 3 {
		t.Fatalf("Unexpected padding. Want 3, have %d",
			address.Padding())
	}
	if address.Type() != AddressType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			AddressType, address.Type())
	}
	if address.Len() != 13 {
		t.Fatalf("Unexpected len. Want 13, have %d", address.Len())
	}
	if len(address.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func BenchmarkGenericAddressIPv4(b *testing.B) {
	address := GenericAddress{AddrType: 0x01, AddrValue: net.ParseIP("10.0.0.1").To4()}
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkGenericAddressIPv6(b *testing.B) {
	address := GenericAddress{AddrType: 0x02, AddrValue: net.ParseIP("2001:db8::ff00:42:8329").To16()}
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkGenericAddressE164(b *testing.B) {
	address := GenericAddress{AddrType: 0x08, AddrValue: []byte("48602007060")}
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}
