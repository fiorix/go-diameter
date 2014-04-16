// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"net"
	"testing"
)

func TestIPv4(t *testing.T) {
	ip4 := IPv4(net.ParseIP("10.0.0.1"))
	b := []byte{0x0a, 0x00, 0x00, 0x01}
	if v := ip4.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	t.Log(ip4)
}

func TestDecodeIPv4(t *testing.T) {
	b := []byte{0x0a, 0x00, 0x00, 0x01}
	ip4, err := DecodeIPv4(b)
	if err != nil {
		t.Fatal(err)
	}
	if ip := net.IP(ip4.(IPv4)).String(); ip != "10.0.0.1" {
		t.Fatalf("Unexpected value. Want 10.0.0.1, have %s", ip)
	}
	t.Log(ip4)
}
