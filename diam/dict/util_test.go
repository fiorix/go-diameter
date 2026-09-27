// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

func TestApps(t *testing.T) {
	apps := Default.Apps()
	if len(apps) != 10 {
		t.Fatalf("Unexpected # of apps. Want 10, have %d", len(apps))
	}
	// Base protocol.
	if apps[0].ID != 0 {
		t.Fatalf("Unexpected app.ID. Want 0, have %d", apps[0].ID)
	}
	// Base accounting
	if apps[1].ID != 3 {
		t.Fatalf("Unexpected app.ID. Want 3, have %d", apps[1].ID)
	}
	// Credit-Control applications.
	if apps[2].ID != 4 {
		t.Fatalf("Unexpected app.ID. Want 4, have %d", apps[2].ID)
	}
	// 3GPP Gx Charging Control applications
	if apps[3].ID != 16777238 {
		t.Fatalf("Unexpected app.ID. Want 16777238, have %d", apps[3].ID)
	}
	// NASREQ applications
	if apps[4].ID != 1 {
		t.Fatalf("Unexpected app.ID. Want 1, have %d", apps[4].ID)
	}
	// 3GPP Rx applications
	if apps[6].ID != 16777236 {
		t.Fatalf("Unexpected app.ID. Want 16777236, have %d", apps[6].ID)
	}
	// 3GPP S6a applications
	if apps[7].ID != 16777251 {
		t.Fatalf("Unexpected app.ID. Want 16777251, have %d", apps[7].ID)
	}
	// 3GPP S13 application
	if apps[8].ID != 16777252 {
		t.Fatalf("Unexpected app.ID. Want 16777252, have %d", apps[8].ID)
	}
	if apps[9].ID != 16777265 {
		t.Fatalf("Unexpected app.ID. Want 16777265, have %d", apps[9].ID)
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

func findAVPCodeTest(t *testing.T, app uint32, codeStr string, vendor, expectedCode uint32) {
	if avp, err := Default.FindAVPWithVendor(app, codeStr, vendor); err != nil {
		t.Fatalf("FindAVP error: %v for app %d & %s AVP", err, app, codeStr)
	} else if avp.Code != expectedCode {
		t.Fatalf(
			"Unexpected code %d for %s AVP and %d vendor. Expected: %d",
			avp.Code, codeStr, vendor, expectedCode)
	}
}

func TestFindAVPWithVendor(t *testing.T) {
	var nokiaXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="43">
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
	findAVPCodeTest(t, 4, "Session-Id", UndefinedVendorID, 263)
	findAVPCodeTest(t, 43, "Session-Start-Indicator", 94, 5105)
	findAVPCodeTest(t, 43, "Session-Start-Indicator", UndefinedVendorID, 5105)

	if _, err := Default.FindAVPWithVendor(4, "Session-Start-Indicator", 0); err == nil {
		t.Error("Should get not found")
	}
	findAVPCodeTest(t, 16777251, "Supported-Features", UndefinedVendorID, 628)

	// Test 'parent' AVP find - S6a app ID, tgpp_ro_rf dictionary
	findAVPCodeTest(t, 16777251, "GMLC-Address", UndefinedVendorID, 2405)

	if _, err := Default.FindAVPWithVendor(43, "User-Password", UndefinedVendorID); err == nil {
		t.Error("User-Password Should not be found for app 43")
	}
	findAVPCodeTest(t, 1, "User-Password", UndefinedVendorID, 2)
	findAVPCodeTest(t, 4, "User-Password", UndefinedVendorID, 2)
	findAVPCodeTest(t, 16777251, "User-Password", UndefinedVendorID, 2)
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

	if cmd, err := Default.FindCommand(16777251, 316); err != nil {
		t.Error(err)
	} else if cmd.Short != "UL" {
		t.Fatalf("Unexpected command: %#v", cmd)
	}

	if cmd, err := Default.FindCommand(16777251, 318); err != nil {
		t.Error(err)
	} else if cmd.Short != "AI" {
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

func TestFindAVPWithVendorNoInfiniteRecursion(t *testing.T) {
	// Looking up a non-existent AVP with a specific vendor ID should return
	// an Unknown AVP instead of causing infinite recursion between
	// FindAVPWithVendor and FindAVP.
	avp, _ := Default.FindAVPWithVendor(4, uint32(99999), 12345)
	if avp == nil {
		t.Fatal("Expected Unknown AVP, got nil")
	}
	if avp.Name != "Unknown-99999-12345" {
		t.Fatalf("Expected Unknown-99999-12345 AVP, got %s", avp.Name)
	}
}

func TestFindAVPByCode(t *testing.T) {
	// Exact (appid, code, vendorID) match.
	if avp, err := Default.FindAVPByCode(4, 461, UndefinedVendorID); err != nil {
		t.Fatalf("FindAVPByCode error for Service-Context-Id: %v", err)
	} else if avp.Name != "Service-Context-Id" {
		t.Fatalf("Unexpected AVP %q, expected Service-Context-Id", avp.Name)
	}

	// Inherited base AVP (app 4 → base) resolves via the pre-merged index.
	if avp, err := Default.FindAVPByCode(4, 263, 0); err != nil {
		t.Fatalf("FindAVPByCode error for inherited Session-Id: %v", err)
	} else if avp.Name != "Session-Id" {
		t.Fatalf("Unexpected AVP %q, expected Session-Id", avp.Name)
	}

	// Inheritance through the full parent chain: Gx (16777238) → 4 → base.
	if avp, err := Default.FindAVPByCode(16777238, 264, 0); err != nil {
		t.Fatalf("FindAVPByCode error for inherited Origin-Host: %v", err)
	} else if avp.Name != "Origin-Host" {
		t.Fatalf("Unexpected AVP %q, expected Origin-Host", avp.Name)
	}

	// Unknown vendor AVP resolves to Unknown, not cross-vendor to base NAS-Port (code 5, vendor 0).
	avp, err := Default.FindAVPByCode(4, 5, 10415)
	if err == nil {
		t.Fatal("Expected error for unknown vendor AVP code 5 / vendor 10415")
	}
	if avp == nil {
		t.Fatal("Expected Unknown AVP, got nil")
	}
	if avp.Name != "Unknown-5-10415" {
		t.Fatalf("Expected Unknown-5-10415, got %q (cross-vendor mismatch)", avp.Name)
	}
	if avp.Data.Type != datatype.UnknownType {
		t.Fatalf("Expected Unknown data type, got %v", avp.Data.Type)
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

// TestCreditControlMSCCUnbounded guards the Credit-Control (Gy) command
// rules for Multiple-Services-Credit-Control. RFC 4006 lists the AVP as
// *[ Multiple-Services-Credit-Control ] in both CCR and CCA, so the rule
// must not carry an upper bound. The dictionary used to declare max="1",
// which contradicted both the RFC and the Ro rule in tgpp_ro_rf.xml.
func TestCreditControlMSCCUnbounded(t *testing.T) {
	cmd, err := Default.FindCommand(4, 272) // Credit-Control
	if err != nil {
		t.Fatal(err)
	}
	find := func(rules []*Rule) *Rule {
		for _, r := range rules {
			if r.AVP == "Multiple-Services-Credit-Control" {
				return r
			}
		}
		return nil
	}
	for _, tc := range []struct {
		name  string
		rules []*Rule
	}{
		{"CCR", cmd.Request.Rule},
		{"CCA", cmd.Answer.Rule},
	} {
		r := find(tc.rules)
		if r == nil {
			t.Errorf("%s: Multiple-Services-Credit-Control rule not found", tc.name)
			continue
		}
		if r.Max != 0 {
			t.Errorf("%s: Multiple-Services-Credit-Control max=%d, want 0 (unbounded)",
				tc.name, r.Max)
		}
	}
}

// RFC 8506 §8 AVP table: codes 653-669, M bit optional, V bit forbidden.
var rfc8506AVPs = []struct {
	code uint32
	name string
	typ  datatype.TypeID
}{
	{653, "User-Equipment-Info-Extension", datatype.GroupedType},
	{654, "User-Equipment-Info-IMEISV", datatype.OctetStringType},
	{655, "User-Equipment-Info-MAC", datatype.OctetStringType},
	{656, "User-Equipment-Info-EUI64", datatype.OctetStringType},
	{657, "User-Equipment-Info-ModifiedEUI64", datatype.OctetStringType},
	{658, "User-Equipment-Info-IMEI", datatype.OctetStringType},
	{659, "Subscription-Id-Extension", datatype.GroupedType},
	{660, "Subscription-Id-E164", datatype.UTF8StringType},
	{661, "Subscription-Id-IMSI", datatype.UTF8StringType},
	{662, "Subscription-Id-SIP-URI", datatype.UTF8StringType},
	{663, "Subscription-Id-NAI", datatype.UTF8StringType},
	{664, "Subscription-Id-Private", datatype.UTF8StringType},
	{665, "Redirect-Server-Extension", datatype.GroupedType},
	{666, "Redirect-Address-IPAddress", datatype.AddressType},
	{667, "Redirect-Address-URL", datatype.UTF8StringType},
	{668, "Redirect-Address-SIP-URI", datatype.UTF8StringType},
	{669, "QoS-Final-Unit-Indication", datatype.GroupedType},
}

func TestCreditControlRFC8506AVPs(t *testing.T) {
	for _, appID := range []uint32{4, 16777238} { // Credit-Control, Gx (inherits 4)
		for _, want := range rfc8506AVPs {
			avp, err := Default.FindAVPByCode(appID, want.code, 0)
			if err != nil {
				t.Errorf("app %d: AVP %d: %v", appID, want.code, err)
				continue
			}
			if avp.Name != want.name || avp.Data.Type != want.typ {
				t.Errorf("app %d: AVP %d = %s/type %d, want %s/type %d", appID, want.code,
					avp.Name, avp.Data.Type, want.name, want.typ)
			}
			if strings.Contains(avp.Must, "M") || !strings.Contains(avp.May, "M") || avp.MustNot != "V" {
				t.Errorf("app %d: %s flags must=%q may=%q must-not=%q, want M optional and V forbidden",
					appID, want.name, avp.Must, avp.May, avp.MustNot)
			}
		}
	}
}

func TestCreditControlRFC8506Rules(t *testing.T) {
	cmd, err := Default.FindCommand(4, 272) // Credit-Control
	if err != nil {
		t.Fatal(err)
	}
	find := func(rules []*Rule, name string) *Rule {
		for _, r := range rules {
			if r.AVP == name {
				return r
			}
		}
		return nil
	}
	grouped := func(name string) []*Rule {
		avp, err := Default.FindAVP(4, name)
		if err != nil {
			t.Fatal(err)
		}
		return avp.Data.Rule
	}
	for _, tc := range []struct {
		where string
		rules []*Rule
		avp   string
		max   int
	}{
		{"CCR", cmd.Request.Rule, "Subscription-Id-Extension", 0},
		{"CCR", cmd.Request.Rule, "User-Equipment-Info-Extension", 1},
		{"CCA", cmd.Answer.Rule, "QoS-Final-Unit-Indication", 1},
		{"Multiple-Services-Credit-Control", grouped("Multiple-Services-Credit-Control"), "QoS-Final-Unit-Indication", 1},
		{"QoS-Final-Unit-Indication", grouped("QoS-Final-Unit-Indication"), "Redirect-Server-Extension", 1},
		{"Subscription-Id-Extension", grouped("Subscription-Id-Extension"), "Subscription-Id-E164", 1},
		{"User-Equipment-Info-Extension", grouped("User-Equipment-Info-Extension"), "User-Equipment-Info-IMEI", 1},
		{"Redirect-Server-Extension", grouped("Redirect-Server-Extension"), "Redirect-Address-SIP-URI", 1},
	} {
		r := find(tc.rules, tc.avp)
		if r == nil {
			t.Errorf("%s: no rule for %s", tc.where, tc.avp)
			continue
		}
		if r.Required || r.Max != tc.max {
			t.Errorf("%s: %s required=%v max=%d, want optional max=%d",
				tc.where, tc.avp, r.Required, r.Max, tc.max)
		}
	}
}

// AVPs the credit-control family (4 and its children) defines as copies of
// NASREQ (1) definitions, keyed by the application that holds the copy.
var nasreqCopies = map[uint32][]string{
	4:        {"Called-Station-Id", "Filter-Id", "Accounting-Input-Octets", "Accounting-Output-Octets"},
	16777238: {"Framed-IP-Address", "Framed-IPv6-Prefix"},
}

func TestCreditControlWithoutNASREQ(t *testing.T) {
	p, err := NewParser(
		"./testdata/base.xml",
		"./testdata/credit_control.xml",
		"./testdata/tgpp_ro_rf.xml",
		"./testdata/gx_credit_control.xml",
	)
	if err != nil {
		t.Fatal(err)
	}
	for appID, names := range nasreqCopies {
		for _, name := range names {
			if _, err := p.FindAVP(appID, name); err != nil {
				t.Errorf("app %d: %v", appID, err)
			}
		}
	}
}

func TestCreditControlRulesDoNotResolveThroughNASREQ(t *testing.T) {
	for _, app := range Default.Apps() {
		if _, scoped := parentAppIds[app.ID]; !scoped || app.ID == 1 {
			continue
		}
		var rules []*Rule
		for _, cmd := range app.Command {
			rules = append(rules, cmd.Request.Rule...)
			rules = append(rules, cmd.Answer.Rule...)
		}
		for _, avp := range app.AVP {
			rules = append(rules, avp.Data.Rule...)
		}
		for _, r := range rules {
			avp, err := Default.FindAVP(app.ID, r.AVP)
			if err == nil && avp.App.ID == 1 {
				t.Errorf("app %d (%s): rule %s resolves only through NASREQ", app.ID, app.Name, r.AVP)
			}
		}
	}
}

func TestNASREQCopiesMatchNASREQ(t *testing.T) {
	for appID, names := range nasreqCopies {
		for _, name := range names {
			cp, err := Default.FindAVP(appID, name)
			if err != nil {
				t.Fatal(err)
			}
			orig, err := Default.FindAVP(1, name)
			if err != nil {
				t.Fatal(err)
			}
			if cp.App.ID == 1 {
				t.Errorf("app %d: %s is not defined by the application", appID, name)
			}
			if cp.Code != orig.Code || cp.Must != orig.Must || cp.May != orig.May ||
				cp.MustNot != orig.MustNot || cp.MayEncrypt != orig.MayEncrypt ||
				cp.VendorID != orig.VendorID || cp.Data.Type != orig.Data.Type {
				t.Errorf("app %d: %s = %+v, NASREQ has %+v", appID, name, *cp, *orig)
			}
		}
	}
}
