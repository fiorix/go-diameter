// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"fmt"
	"testing"
)

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
