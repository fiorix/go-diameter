// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"fmt"
	"testing"
)

// Tests

// TestPad4 tests padding to 32 bits.
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

// TestIntConversion tests uint 24 to 32 and vice-versa.
func TestIntConversion(t *testing.T) {
	b := uint32to24(4096)
	n := uint24to32(b)
	if n != 4096 {
		t.Error(fmt.Errorf("uint24 0x%x != uint32 0x%x", b, n))
	}
}
