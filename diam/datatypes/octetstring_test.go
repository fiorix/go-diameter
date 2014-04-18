// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"testing"
)

func TestOctetString(t *testing.T) {
	s := OctetString("hello, world")
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	if v := s.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	t.Log(s)
	t.Log(string(s)) // ?
}

func TestDecodeOctetString(t *testing.T) {
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	s, err := DecodeOctetString(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(s.(OctetString)), b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, s)
	}
	t.Log(s)
}

func BenchmarkOctetString(b *testing.B) {
	v := OctetString("hello, world")
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeOctetString(b *testing.B) {
	v := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	for n := 0; n < b.N; n++ {
		DecodeOctetString(v)
	}
}
