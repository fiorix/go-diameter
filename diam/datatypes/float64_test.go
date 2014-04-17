// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"testing"
)

func TestFloat64(t *testing.T) {
	n := Float64(3.1415926535)
	b := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x41, 0x17, 0x44}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeFloat64(t *testing.T) {
	b := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x41, 0x17, 0x44}
	n, err := DecodeFloat64(b)
	if err != nil {
		t.Fatal(err)
	}
	if v := n.(Float64); v != 3.1415926535 {
		t.Fatalf("Unexpected value. Want 3.1415926535, have %0.4f", v)
	}
	t.Log(n)
}

func BenchmarkFloat64(b *testing.B) {
	v := Float64(3.1415926535)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeFloat64(b *testing.B) {
	v := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x41, 0x17, 0x44}
	for n := 0; n < b.N; n++ {
		DecodeFloat64(v)
	}
}
