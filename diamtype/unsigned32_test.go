// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import (
	"bytes"
	"math"
	"testing"
)

type zz struct{ uint32 }

func TestUnsigned32(t *testing.T) {
	n := Unsigned32(math.MaxUint32)
	b := []byte{0xff, 0xff, 0xff, 0xff}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeUnsigned32(t *testing.T) {
	b := []byte{0xff, 0xff, 0xff, 0xff}
	n, err := DecodeUnsigned32(b)
	if err != nil {
		t.Fatal(err)
	}
	z := uint32(math.MaxUint32)
	if uint32(n.(Unsigned32)) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
	t.Log(n)
}

func BenchmarkUnsigned32(b *testing.B) {
	v := Unsigned32(math.MaxUint32)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeUnsigned32(b *testing.B) {
	v := []byte{0xff, 0xc0, 0xff, 0xee}
	for n := 0; n < b.N; n++ {
		DecodeUnsigned32(v)
	}
}
