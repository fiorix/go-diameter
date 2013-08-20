// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package base

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestOctetString(t *testing.T) {
	s := "hello, world!"
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
		// Padding 0x00, 0x00, 0x00,
	}
	os := &OctetString{Value: s}
	//os := &DiameterIdentity{OctetString{Value: s}}
	osb := os.Bytes()
	if os.Padding > 0 {
		osb = osb[:uint32(len(osb))-os.Padding]
	}
	if !bytes.Equal(osb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", osb, b))
		return
	}
	os.Put(b)
	if d := os.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
	if os.Padding != 3 {
		t.Error(fmt.Errorf("Padding length of '%s' is %d, expected 3.",
			s, os.Padding))
		return
	}
}

func TestAddress(t *testing.T) {
	s := net.ParseIP("192.168.4.20")
	b := []byte{
		0,
		1,                      // AF_INET=1
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

func TestDiameterURI(t *testing.T) {
	s := "aaa://diameter:3868;transport=tcp"
	du := DiameterURI{Value: s}
	du.Put(du.Bytes())
	if d := du.Data(); d != s {
		t.Error(fmt.Errorf("Data is '%s', expected '%s'", d, s))
	}
}

func TestInteger32(t *testing.T) {
	s := int32((1 << 31) - 1)
	b := []byte{0x7f, 0xff, 0xff, 0xff}
	n := Integer32{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func TestInteger64(t *testing.T) {
	s := int64((1 << 63) - 1)
	b := []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	n := Integer64{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func TestUnsigned32(t *testing.T) {
	s := uint32(0xffc0ffee)
	b := []byte{0xff, 0xc0, 0xff, 0xee}
	n := Unsigned32{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func TestUnsigned64(t *testing.T) {
	s := uint64(0xffffffffffc0ffee)
	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xc0, 0xff, 0xee}
	n := Unsigned64{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func TestFloat32(t *testing.T) {
	s := float32(4.20)
	b := []byte{0x40, 0x86, 0x66, 0x66}
	n := Float32{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func TestFloat64(t *testing.T) {
	s := float64(4.20)
	b := []byte{0x40, 0x10, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcd}
	n := Float64{Value: s}
	nb := n.Bytes()
	if !bytes.Equal(nb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", nb, b))
		return
	}
	n.Put(b)
	if d := n.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}
