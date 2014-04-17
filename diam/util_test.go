// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"testing"
)

func TestUint24to32(t *testing.T) {
	if v := uint24to32([]byte{218, 190, 239}); v != 0xdabeef {
		t.Fatalf("Unexpected result. Want 0xdeadbeef, have 0x%x", v)
	}
}

func TestUint32to24(t *testing.T) {
	if v := uint32to24(0xdabeef); !bytes.Equal(v, []byte{218, 190, 239}) {
		t.Fatalf("Unexpected result. Want [218 190 239], have %v", v)
	}
}
