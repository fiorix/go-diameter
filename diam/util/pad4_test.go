// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package util

import (
	"fmt"
	"testing"
)

func TestPad4(t *testing.T) {
	switch {
	case Pad4(3) != 4:
		t.Error(fmt.Errorf("pad(3) != 4"))
	case Pad4(5) != 8:
		t.Error(fmt.Errorf("pad(5) != 8"))
	case Pad4(10) != 12:
		t.Error(fmt.Errorf("pad(10) != 12"))
	}
}
