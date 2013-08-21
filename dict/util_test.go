// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package dict

import (
	"fmt"
	"testing"
)

func TestFindAVP(t *testing.T) {
	if _, err := Base.FindAVP(0, 263); err != nil {
		t.Error(err)
	}
}

func TestFindCmd(t *testing.T) {
	if cmd, err := Base.FindCmd(0, 257); err != nil {
		t.Error(err)
	} else if cmd.Short != "CE" {
		t.Error(fmt.Errorf("Unexpected command: %#v", cmd))
	}
}

func TestCodeFor(t *testing.T) {
	if n := Base.CodeFor("Session-Id"); n != 263 {
		t.Error(fmt.Errorf("Unexpected code %d for Session-Id AVP", n))
	}
}

func TestAppFor(t *testing.T) {
	if app := Base.AppFor("Session-Id"); app == nil {
		t.Error("Could not find app for Session-Id AVP")
	}
}
