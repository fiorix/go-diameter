// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package dict

import (
	"fmt"
	"testing"
)

func init() {
	if false {
		Base.PrettyPrint()
	}
}

func TestVendors(t *testing.T) {
	s := &Vendor{Id: 10415}
	v := Base.Vendors()
	if len(v) != 1 || v[0].Id != s.Id {
		t.Error(fmt.Errorf("Unexpected vendor: %#v", v[0]))
	}
}

func TestFindAVP(t *testing.T) {
	if _, err := Base.FindAVP(0, 263); err != nil {
		t.Error(err)
	}
}

func TestScanAVP(t *testing.T) {
	if _, err := Base.ScanAVP(263); err != nil {
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

func TestEnum(t *testing.T) {
	if item, err := Base.Enum(0, 274, 1); err != nil {
		t.Error(err)
		return
	} else if item.Name != "AUTHENTICATE_ONLY" {
		t.Error(fmt.Errorf(
			"Unexpected value %s, expected AUTHENTICATE_ONLY",
			item.Name))
	}
}

func TestRule(t *testing.T) {
	if rule, err := Base.Rule(0, 284, "Proxy-Host"); err != nil {
		t.Error(err)
		return
	} else if !rule.Required {
		t.Error(fmt.Errorf("Unexpected rule %s", rule))
	}
}
