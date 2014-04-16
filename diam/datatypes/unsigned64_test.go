// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"testing"
)

func TestUnsigned64(t *testing.T) {
	n := Unsigned64(0xffffffffffc0ffee)
	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xc0, 0xff, 0xee}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeUnsigned64(t *testing.T) {
	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xc0, 0xff, 0xee}
	n, err := DecodeUnsigned64(b)
	if err != nil {
		t.Fatal(err)
	}
	z := uint64(0xffffffffffc0ffee)
	if uint64(n.(Unsigned64)) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
	t.Log(n)
}

func BenchmarkUnsigned64(b *testing.B) {
	v := Unsigned64(0xffffffffffc0ffee)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeUnsigned64(b *testing.B) {
	v := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xc0, 0xff, 0xee}
	for n := 0; n < b.N; n++ {
		DecodeUnsigned64(v)
	}
}
