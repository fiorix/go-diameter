// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"bytes"
	"testing"
)

func TestApps(t *testing.T) {
	apps := Default.Apps()
	if len(apps) != 5 {
		t.Fatalf("Unexpected # of apps. Want 5, have %d", len(apps))
	}
	// Base protocol.
	if apps[0].ID != 0 {
		t.Fatalf("Unexpected app.ID. Want 0, have %d", apps[0].ID)
	}
	// Base accounting.
	if apps[1].ID != 3 {
		t.Fatalf("Unexpected app.ID. Want 3, have %d", apps[1].ID)
	}
	// Credit-Control applications.
	if apps[2].ID != 4 {
		t.Fatalf("Unexpected app.ID. Want 4, have %d", apps[2].ID)
	}
	// Network Access Server
	if apps[3].ID != 1 {
		t.Fatalf("Unexpected app.ID. Want 1, have %d", apps[3].ID)
	}
}

func TestApp(t *testing.T) {
	// Base protocol.
	if _, err := Default.App(0); err != nil {
		t.Fatal(err)
	}
	// Credit-Control applications.
	if _, err := Default.App(4); err != nil {
		t.Fatal(err)
	}
}

func TestFindAVPWithVendor(t *testing.T) {
	var nokiaXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="4">
    <vendor id="94" name="Nokia" />
    <avp name="Session-Start-Indicator" code="5105" must="V" may="P,M" must-not="-" may-encrypt="N" vendor-id="94">
      <data type="UTF8String" />
    </avp>
  </application>
</diameter>`
	Default.Load(bytes.NewReader([]byte(nokiaXML)))
	if _, err := Default.FindAVPWithVendor(4, 999, UndefinedVendorID); err == nil {
		t.Error("Should get not found")
	}
	if avp, err := Default.FindAVPWithVendor(4, "Session-Id", UndefinedVendorID); err != nil {
		t.Fatal(err)
	} else if avp.Code != 263 {
		t.Fatalf("Unexpected code %d for Session-Id AVP", avp.Code)
	}
	if avp, err := Default.FindAVPWithVendor(4, "Session-Start-Indicator", 94); err != nil {
		t.Fatal(err)
	} else if avp.Code != 5105 {
		t.Fatalf("Unexpected code %d for Session-Id AVP", avp.Code)
	}
	if avp, err := Default.FindAVPWithVendor(4, "Session-Start-Indicator", UndefinedVendorID); err != nil {
		t.Fatal(err)
	} else if avp.Code != 5105 {
		t.Fatalf("Unexpected code %d for Session-Id AVP", avp.Code)
	}
	if _, err := Default.FindAVPWithVendor(4, "Session-Start-Indicator", 0); err == nil {
		t.Error("Should get not found")
	}
}

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

func TestFindCommand(t *testing.T) {
	if cmd, err := Default.FindCommand(999, 257); err != nil {
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
