// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"fmt"
	"testing"
)

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
