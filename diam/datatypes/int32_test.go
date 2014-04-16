// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"testing"
)

func TestInteger32(t *testing.T) {
	n := Integer32((1 << 31) - 1)
	b := []byte{0x7f, 0xff, 0xff, 0xff}
	if x := n.Serialize(); !bytes.Equal(b, x) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, x)
	}
	t.Log(n)
}

func TestDecodeInteger32(t *testing.T) {
	b := []byte{0x7f, 0xff, 0xff, 0xff}
	n, err := DecodeInteger32(b)
	if err != nil {
		t.Fatal(err)
	}
	z := int32((1 << 31) - 1)
	if int32(n) != z {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", z, n)
	}
	t.Log(n)
}

func BenchmarkInteger32(b *testing.B) {
	v := Integer32((1 << 31) - 1)
	for n := 0; n < b.N; n++ {
		v.Serialize()
	}
}

func BenchmarkDecodeInteger32(b *testing.B) {
	v := []byte{0x7f, 0xff, 0xff, 0xff}
	for n := 0; n < b.N; n++ {
		DecodeInteger32(v)
	}
}
