// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"math"
	"testing"
)

func TestUnsigned64(t *testing.T) {
	n := Unsigned64(math.MaxUint64)
	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeUnsigned64(t *testing.T) {
	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	n, err := DecodeUnsigned64(b)
	if err != nil {
		t.Fatal(err)
	}
	z := uint64(math.MaxUint64)
	if uint64(n.(Unsigned64)) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
	t.Log(n)
}

func BenchmarkUnsigned64(b *testing.B) {
	v := Unsigned64(math.MaxUint64)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeUnsigned64(b *testing.B) {
	v := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for n := 0; n < b.N; n++ {
		DecodeUnsigned64(v)
	}
}
