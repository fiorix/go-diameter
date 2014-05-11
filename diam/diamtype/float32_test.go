// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import (
	"bytes"
	"testing"
)

func TestFloat32(t *testing.T) {
	n := Float32(3.1415)
	b := []byte{0x40, 0x49, 0x0e, 0x56}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeFloat32(t *testing.T) {
	b := []byte{0x40, 0x49, 0x0e, 0x56}
	n, err := DecodeFloat32(b)
	if err != nil {
		t.Fatal(err)
	}
	if v := n.(Float32); v != 3.1415 {
		t.Fatalf("Unexpected value. Want 3.1414, have %0.4f", v)
	}
	t.Log(n)
}

func BenchmarkFloat32(b *testing.B) {
	v := Float32(3.1415)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeFloat32(b *testing.B) {
	v := []byte{0x40, 0x49, 0x0e, 0x56}
	for n := 0; n < b.N; n++ {
		DecodeFloat32(v)
	}
}
