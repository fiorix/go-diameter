// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import "testing"

func TestFindAVP(t *testing.T) {
	if _, err := Default.FindAVP(999, 263); err != nil {
		t.Fatal(err)
	}
}

func TestScanAVP(t *testing.T) {
	if avp, err := Default.ScanAVP("Session-Id"); err != nil {
		t.Error(err)
	} else if avp.Code != 263 {
		t.Fatalf("Unexpected code %d for Session-Id AVP", avp.Code)
	}
}

func TestFindCmd(t *testing.T) {
	if cmd, err := Default.FindCMD(999, 257); err != nil {
		t.Error(err)
	} else if cmd.Short != "CE" {
		t.Fatalf("Unexpected command: %#v", cmd)
	}
}

func TestEnum(t *testing.T) {
	if item, err := Default.Enum(0, 274, 1); err != nil {
		t.Fatal(err)
	} else if item.Name != "AUTHENTICATE_ONLY" {
		t.Errorf(
			"Unexpected value %s, expected AUTHENTICATE_ONLY",
			item.Name,
		)
	}
}

func TestRule(t *testing.T) {
	if rule, err := Default.Rule(0, 284, "Proxy-Host"); err != nil {
		t.Fatal(err)
	} else if !rule.Required {
		t.Errorf("Unexpected rule %#v", rule)
	}
}

func BenchmarkFindAVPName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.FindAVP(0, "Session-Id")
	}
}

func BenchmarkFindAVPCode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.FindAVP(0, 263)
	}
}

func BenchmarkScanAVPName(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.ScanAVP("Session-Id")
	}
}

func BenchmarkScanAVPCode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Default.ScanAVP(263)
	}
}
