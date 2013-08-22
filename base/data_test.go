// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests and Benchmarks

package base

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

// Tests
// go test -v

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

func TestTime(t *testing.T) {
	s := int64(1377093974)
	b := []byte{0x52, 0x14, 0xc9, 0x56}
	tm := Time{Value: time.Unix(s, 0)}
	tmb := tm.Bytes()
	if !bytes.Equal(tmb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", tmb, b))
		return
	}
	tm.Put(b)
	if d := tm.Data().(time.Time).Unix(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
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
	b := []byte{
		0x61, 0x61, 0x61, 0x3a, 0x2f, 0x2f, 0x64, 0x69, 0x61, 0x6d,
		0x65, 0x74, 0x65, 0x72, 0x3a, 0x33, 0x38, 0x36, 0x38, 0x3b,
		0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x3d,
		0x74, 0x63, 0x70,
	}
	du := DiameterURI{Value: s}
	dub := du.Bytes()
	if !bytes.Equal(dub, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", dub, b))
		return
	}
	du.Put(b)
	if d := du.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
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

// Benchmarks
// go test -bench .

func BenchmarkOctetStringParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(OctetString)
		p.Put([]byte{
			0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20,
			0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
		})
	}
}

func BenchmarkOctetStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &OctetString{Value: "Hello, world"}
		p.Bytes()
	}
}

func BenchmarkTimeParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(Time)
		p.Put([]byte{0x52, 0x14, 0xc9, 0x56})
	}
}

func BenchmarkTimeBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &Time{Value: time.Unix(1377093974, 0)}
		p.Bytes()
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

func BenchmarkIPv4Parser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(IPv4)
		p.Put([]byte{0xc0, 0xa8, 0x04, 0x14})
	}
}

func BenchmarkIPv4Builder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &IPv4{IP: net.ParseIP("192.168.4.20")}
		p.Bytes()
	}
}

func BenchmarkDiameterURIParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(DiameterURI)
		p.Put([]byte{
			0x61, 0x61, 0x61, 0x3a, 0x2f, 0x2f, 0x64, 0x69, 0x61,
			0x6d, 0x65, 0x74, 0x65, 0x72, 0x3a, 0x33, 0x38, 0x36,
			0x38, 0x3b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
			0x72, 0x74, 0x3d, 0x74, 0x63, 0x70,
		})
	}
}

func BenchmarkDiameterURIBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &DiameterURI{Value: "aaa://diameter:3868;transport=tcp"}
		p.Bytes()
	}
}
