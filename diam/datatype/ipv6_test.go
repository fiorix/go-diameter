// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"net"
	"testing"
)

func TestIPv6(t *testing.T) {
	ip6 := IPv6(net.ParseIP("2001:0db8::ff00:0042:8329"))
	b := []byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29}
	if v := ip6.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if ip6.Len() != 16 {
		t.Fatalf("Unexpected leip6. Want 4, have %d", ip6.Len())
	}
	if ip6.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", ip6.Padding())
	}
	if ip6.Type() != IPv6Type {
		t.Fatalf("Unexpected type. Want %d, have %d",
			IPv6Type, ip6.Type())
	}
	if len(ip6.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestIPv6Malformed(t *testing.T) {
	ip6 := IPv6(net.ParseIP("10.0.0.1"))
	b := []byte{0x0a, 0x00, 0x00, 0x01}
	if v := ip6.Serialize(); bytes.Equal(v, b) {
		t.Fatalf("IPv6 match, that's unexpected")
	}
}

func TestDecodeIPv6(t *testing.T) {
	b := []byte{0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x42, 0x83, 0x29}
	ip6, err := DecodeIPv6(b)
	if err != nil {
		t.Fatal(err)
	}
	if ip := net.IP(ip6.(IPv6)).String(); ip != "2001:db8::ff00:42:8329" {
		t.Fatalf("Unexpected value. Want '2001:db8::ff00:42:8329', have %s", ip)
	}
}
