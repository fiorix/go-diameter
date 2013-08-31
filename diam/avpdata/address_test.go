// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestAddress(t *testing.T) {
	s := net.ParseIP("192.168.4.20")
	b := []byte{
		0, 1, // AF_INET=1
		0xc0, 0xa8, 0x04, 0x14, // IP
	}
	addr := new(Address)
	addr.Put(b)
	if d := addr.Data(); !bytes.Equal(d.(net.IP).To4(), s.To4()) {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
	ab := addr.Bytes()
	if addr.Padding > 0 {
		ab = ab[:len(ab)-addr.Padding]
	}
	if !bytes.Equal(ab, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", ab, b))
		return
	}
}

func BenchmarkAddressParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(Address)
		p.Put([]byte{0, 1, 0xc0, 0xa8, 0x04, 0x14})
	}
}

func BenchmarkAddressBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &Address{IP: net.ParseIP("192.168.4.20")}
		p.Bytes()
	}
}
