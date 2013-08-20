// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

import (
	"bytes"
	"fmt"
	"testing"
)

func init() {
	fmt.Printf("oi avp_data_test\n")
}

// Tests

func TestOctetString(t *testing.T) {
	s := "hello, world!"
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
		// Padding 0x00, 0x00, 0x00,
	}
	os := OctetString{String: s}
	osb := os.Bytes()
	if os.Padding > 0 {
		osb = osb[:len(osb)-os.Padding]
	}
	if !bytes.Equal(osb, b) {
		t.Error(fmt.Errorf("Bytes are '%s', expected '%s'", osb, b))
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

func TestDiameterURI(t *testing.T) {
	s := "aaa://diameter:3868;transport=tcp"
	du := DiameterURI{String: s}
	du.Put(du.Bytes())
	if d := du.Data(); d != s {
		t.Error(fmt.Errorf("Data is '%s', expected '%s'", d, s))
	}
}

func TestInteger32(t *testing.T) {
	s := int32((1 << 31) - 1)
	b := []byte{0x7f, 0xff, 0xff, 0xff}
	n := Integer32{Int32: s}
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
	n := Integer64{Int64: s}
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
	n := Unsigned32{Uint32: s}
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
	n := Unsigned64{Uint64: s}
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
	n := Float32{Float32: s}
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
	n := Float64{Float64: s}
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
