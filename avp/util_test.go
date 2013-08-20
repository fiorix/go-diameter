// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

import (
	"fmt"
	"testing"
)

// Tests

func TestPad4(t *testing.T) {
	switch {
	case pad4(3) != 4:
		t.Error(fmt.Errorf("pad(3) != 4"))
	case pad4(5) != 8:
		t.Error(fmt.Errorf("pad(5) != 8"))
	case pad4(10) != 12:
		t.Error(fmt.Errorf("pad(10) != 12"))
	}
}

func TestUintConversion(t *testing.T) {
	b := uint32To24(4096)
	n := uint24To32(b)
	if n != 4096 {
		t.Error(fmt.Errorf("uint24 0x%x != 0x%x uint32", b, n))
	}
}
