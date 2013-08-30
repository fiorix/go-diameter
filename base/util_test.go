// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"fmt"
	"testing"
)

func TestUintConversion(t *testing.T) {
	b := uint32To24(4096)
	n := uint24To32(b)
	if n != 4096 {
		t.Error(fmt.Errorf("uint24 0x%x != 0x%x uint32", b, n))
	}
}
