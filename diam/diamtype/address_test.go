// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import (
	"bytes"
	"net"
	"testing"
)

func TestAddressIPv4(t *testing.T) {
	address := Address(net.ParseIP("10.0.0.1"))
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d", address.Padding())
	}
	t.Log(address)
}

func TestDecodeAddressIPv4(t *testing.T) {
	b := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	address, err := DecodeAddress(b)
	if err != nil {
		t.Fatal(err)
	}
	if ip := net.IP(address.(Address)).String(); ip != "10.0.0.1" {
		t.Fatalf("Unexpected value. Want 10.0.0.1, have %s", ip)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d", address.Padding())
	}
	t.Log(address)
}

func TestAddressIPv6(t *testing.T) {
	address := Address(net.ParseIP("2001:0db8::ff00:0042:8329"))
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	if v := address.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d", address.Padding())
	}
	t.Log(address)
}

func TestDecodeAddressIPv6(t *testing.T) {
	b := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	address, err := DecodeAddress(b)
	if err != nil {
		t.Fatal(err)
	}
	want := "2001:db8::ff00:42:8329"
	if ip := net.IP(address.(Address)).String(); ip != want {
		t.Fatalf("Unexpected value. Want %s, have %s", want, ip)
	}
	if address.Padding() != 2 {
		t.Fatalf("Unexpected padding. Want 2, have %d", address.Padding())
	}
	t.Log(address)
}

func BenchmarkAddressIPv4(b *testing.B) {
	address := Address(net.ParseIP("10.0.0.1"))
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkDecodeAddressIPv4(b *testing.B) {
	v := []byte{0x00, 0x01, 0x0a, 0x00, 0x00, 0x01}
	for n := 0; n < b.N; n++ {
		DecodeAddress(v)
	}
}

func BenchmarkAddressIPv6(b *testing.B) {
	address := Address(net.ParseIP("2001:db8::ff00:42:8329"))
	for n := 0; n < b.N; n++ {
		address.Serialize()
	}
}

func BenchmarkDecodeAddressIPv6(b *testing.B) {
	v := []byte{0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29,
	}
	for n := 0; n < b.N; n++ {
		DecodeAddress(v)
	}
}
