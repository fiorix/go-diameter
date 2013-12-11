// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package dict

import (
	"fmt"
	"testing"
)

var Base *Parser

func init() {
	var err error
	Base, err = New("./diam_base.xml")
	if err != nil {
		panic(err)
	}
	if false {
		Base.PrettyPrint()
	}
}

func TestFindAVP(t *testing.T) {
	if _, err := Base.FindAVP(0, 263); err != nil {
		t.Error(err)
	}
}

func TestScanAVP(t *testing.T) {
	if avp, err := Base.ScanAVP("Session-Id"); err != nil {
		t.Error(err)
	} else if avp.Code != 263 {
		t.Error(fmt.Errorf(
			"Unexpected code %d for Session-Id AVP", avp.Code))
	}
}

func TestFindCmd(t *testing.T) {
	if cmd, err := Base.FindCmd(0, 257); err != nil {
		t.Error(err)
	} else if cmd.Short != "CE" {
		t.Error(fmt.Errorf("Unexpected command: %#v", cmd))
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
		t.Error(fmt.Errorf("Unexpected rule %#v", rule))
	}
}
