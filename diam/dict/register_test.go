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

// registerTestXML declares app 0 and app 4, so that the inheritance chain in
// parentAppIds (4 -> 0) has both ends loaded, plus one AVP whose code is
// reused by a second vendor.
const registerTestXML = `<diameter>
 <application id="0">
  <avp name="Register-Base" code="1000" mandatory="must">
   <data type="Unsigned32"/>
  </avp>
  <avp name="Register-Vendor-Owned" code="1001" vendor-id="10415" mandatory="must">
   <data type="UTF8String"/>
  </avp>
 </application>
 <application id="4">
  <avp name="Register-App4" code="1002" mandatory="must">
   <data type="Unsigned32"/>
  </avp>
 </application>
 <vendor id="10415" name="TGPP"/>
</diameter>`

const registerTestXML2 = `<diameter>
 <application id="7">
  <avp name="Register-Second-Load" code="1003" mandatory="must">
   <data type="Unsigned32"/>
  </avp>
 </application>
</diameter>`

func newRegisterTestParser(t *testing.T) *Parser {
	t.Helper()
	p := new(Parser)
	p.Strict = true
	if err := p.Load(bytes.NewReader([]byte(registerTestXML))); err != nil {
		t.Fatalf("load test dictionary: %v", err)
	}
	return p
}

// A Parser that never loaded XML has nil maps, since Load is what allocates
// them. Registering must still work rather than panicking.
func TestRegisterGroupedAVPOnFreshParser(t *testing.T) {
	for _, tc := range []struct {
		name string
		p    *Parser
	}{
		{"zero value", new(Parser)},
		{"NewParser with no files", func() *Parser {
			p, err := NewParser()
			if err != nil {
				t.Fatalf("NewParser: %v", err)
			}
			return p
		}()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.p.RegisterGroupedAVP(0, 2000, 10415, "Fresh-Grouped"); err != nil {
				t.Fatalf("RegisterGroupedAVP: %v", err)
			}
			avp, err := tc.p.FindAVPByCode(0, 2000, 10415)
			if err != nil {
				t.Fatalf("FindAVPByCode: %v", err)
			}
			if avp.Name != "Fresh-Grouped" {
				t.Errorf("Name = %q, want Fresh-Grouped", avp.Name)
			}
			if avp.Data.Type != datatype.GroupedType {
				t.Errorf("Data.Type = %v, want GroupedType", avp.Data.Type)
			}
		})
	}
}

// Registering before any Load, then loading, must keep both: the index is
// allocated once and neither entry point resets it.
func TestRegisterGroupedAVPSurvivesLaterLoad(t *testing.T) {
	p := new(Parser)
	if err := p.RegisterGroupedAVP(0, 2000, 10415, "Survives-Load"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	if err := p.Load(bytes.NewReader([]byte(registerTestXML))); err != nil {
		t.Fatalf("Load after registering: %v", err)
	}
	if _, err := p.FindAVPByCode(0, 2000, 10415); err != nil {
		t.Errorf("registration lost after a later Load: %v", err)
	}
	if _, err := p.FindAVPByCode(0, 1000, 0); err != nil {
		t.Errorf("loaded dictionary missing: %v", err)
	}
	// And again for a second dictionary, which Load explicitly supports.
	if err := p.Load(bytes.NewReader([]byte(registerTestXML2))); err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if _, err := p.FindAVPByCode(0, 2000, 10415); err != nil {
		t.Errorf("registration lost after a second Load: %v", err)
	}
	if _, err := p.FindAVPByCode(7, 1003, 0); err != nil {
		t.Errorf("second dictionary not loaded: %v", err)
	}
}

// The registered entry is keyed by {Code, VendorID} and must not be reachable
// through another vendor's lookup, nor displace an AVP that reuses its code.
func TestRegisterGroupedAVPIsVendorScoped(t *testing.T) {
	p := newRegisterTestParser(t)
	// Code 1000 already belongs to Register-Base in the IETF space.
	if err := p.RegisterGroupedAVP(0, 1000, 10415, "Register-Shadow"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	base, err := p.FindAVPByCode(0, 1000, 0)
	if err != nil {
		t.Fatalf("vendor 0 lookup: %v", err)
	}
	if base.Name != "Register-Base" {
		t.Errorf("vendor 0 code 1000 = %q, want Register-Base: the registration displaced it", base.Name)
	}
	if base.Data.Type != datatype.Unsigned32Type {
		t.Errorf("vendor 0 code 1000 type = %v, want Unsigned32", base.Data.Type)
	}
	// The vendor-agnostic slot must be untouched, or an unknown vendor would
	// resolve cross-vendor.
	if generic, ok := p.avpcode[codeIdx{0, 1000, UndefinedVendorID}]; ok {
		if generic.Name != "Register-Base" {
			t.Errorf("generic slot for code 1000 = %q, want Register-Base", generic.Name)
		}
	}
	// A third vendor must not reach the registration.
	if avp, err := p.FindAVPByCode(0, 1000, 99999); err == nil {
		t.Errorf("vendor 99999 resolved to %q, want not found", avp.Name)
	}
	// The registration itself is reachable under its own vendor.
	got, err := p.FindAVPByCode(0, 1000, 10415)
	if err != nil {
		t.Fatalf("vendor 10415 lookup: %v", err)
	}
	if got.Name != "Register-Shadow" {
		t.Errorf("vendor 10415 code 1000 = %q, want Register-Shadow", got.Name)
	}
}

// FindAVPByCode resolves one index entry and does not walk the parent chain,
// so a base-app registration has to be propagated into descendant apps.
func TestRegisterGroupedAVPInheritedByChildApp(t *testing.T) {
	p := newRegisterTestParser(t)
	if err := p.RegisterGroupedAVP(0, 2000, 10415, "Inherited-Grouped"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	// App 4 descends from app 0 via parentAppIds.
	avp, err := p.FindAVPByCode(4, 2000, 10415)
	if err != nil {
		t.Fatalf("child app 4 cannot see the base registration: %v", err)
	}
	if avp.Name != "Inherited-Grouped" {
		t.Errorf("app 4 code 2000 = %q, want Inherited-Grouped", avp.Name)
	}
	// An app-specific registration must not leak into an unrelated app.
	if err := p.RegisterGroupedAVP(4, 2001, 10415, "App4-Only"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	if avp, err := p.FindAVPByCode(0, 2001, 10415); err == nil {
		t.Errorf("app 0 resolved app-4 registration to %q, want not found", avp.Name)
	}
}

// A registration must be reachable by name, with its vendor.
func TestRegisterGroupedAVPNameLookup(t *testing.T) {
	p := newRegisterTestParser(t)
	if err := p.RegisterGroupedAVP(0, 2000, 10415, "Named-Grouped"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	avp, err := p.FindAVPWithVendor(uint32(0), "Named-Grouped", 10415)
	if err != nil {
		t.Fatalf("FindAVPWithVendor by name: %v", err)
	}
	if avp.Code != 2000 {
		t.Errorf("Code = %d, want 2000", avp.Code)
	}
}

func TestRegisterGroupedAVPIdempotent(t *testing.T) {
	p := newRegisterTestParser(t)
	if err := p.RegisterGroupedAVP(0, 2000, 10415, "Idem-Grouped"); err != nil {
		t.Fatalf("first: %v", err)
	}
	if err := p.RegisterGroupedAVP(0, 2000, 10415, "Idem-Grouped"); err != nil {
		t.Errorf("re-registering the same definition: %v, want nil", err)
	}
}

func TestRegisterGroupedAVPRejectsInvalid(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		appID    uint32
		code     uint32
		vendorID uint32
		name     string
		wantErr  string
	}{
		{"empty name", 0, 2000, 10415, "", "name is empty"},
		{"zero code", 0, 0, 10415, "Zero-Code", "code 0 is reserved"},
		{"wildcard vendor", 0, 2000, UndefinedVendorID, "Wildcard", "not a vendor"},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			p := newRegisterTestParser(t)
			err := p.RegisterGroupedAVP(tc.appID, tc.code, tc.vendorID, tc.name)
			if err == nil {
				t.Fatalf("RegisterGroupedAVP = nil, want an error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestRegisterGroupedAVPRejectsConflict(t *testing.T) {
	t.Run("code taken by a loaded AVP", func(t *testing.T) {
		p := newRegisterTestParser(t)
		err := p.RegisterGroupedAVP(0, 1001, 10415, "Different-Name")
		if err == nil {
			t.Fatal("RegisterGroupedAVP = nil, want a conflict error")
		}
		if !strings.Contains(err.Error(), "Register-Vendor-Owned") {
			t.Errorf("error = %q, want it to name the existing AVP", err)
		}
		// The loaded AVP must be unchanged.
		avp, err := p.FindAVPByCode(0, 1001, 10415)
		if err != nil {
			t.Fatalf("FindAVPByCode: %v", err)
		}
		if avp.Data.Type != datatype.UTF8StringType {
			t.Errorf("existing AVP type = %v, want UTF8String", avp.Data.Type)
		}
	})
	t.Run("name taken by another code", func(t *testing.T) {
		p := newRegisterTestParser(t)
		if err := p.RegisterGroupedAVP(0, 2000, 10415, "Dup-Name"); err != nil {
			t.Fatalf("first: %v", err)
		}
		if err := p.RegisterGroupedAVP(0, 2001, 10415, "Dup-Name"); err == nil {
			t.Error("RegisterGroupedAVP = nil, want a name conflict error")
		}
		// The first registration must survive the rejected second one.
		if _, err := p.FindAVPByCode(0, 2000, 10415); err != nil {
			t.Errorf("first registration lost: %v", err)
		}
	})
}

// The AVP must carry the fields the rest of the package reads: a resolved
// type, its printable name, and a link to its application.
func TestRegisterGroupedAVPFields(t *testing.T) {
	p := newRegisterTestParser(t)
	if err := p.RegisterGroupedAVP(4, 2000, 10415, "Fields-Grouped"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	avp, err := p.FindAVPByCode(4, 2000, 10415)
	if err != nil {
		t.Fatalf("FindAVPByCode: %v", err)
	}
	if avp.Data.TypeName != "Grouped" {
		t.Errorf("Data.TypeName = %q, want Grouped", avp.Data.TypeName)
	}
	if avp.Data.Type != datatype.Available[avp.Data.TypeName] {
		t.Errorf("Data.Type = %v, inconsistent with TypeName %q",
			avp.Data.Type, avp.Data.TypeName)
	}
	if avp.App == nil {
		t.Fatal("App is nil")
	}
	if avp.App.ID != 4 {
		t.Errorf("App.ID = %d, want 4", avp.App.ID)
	}
	if avp.VendorID != 10415 {
		t.Errorf("VendorID = %d, want 10415", avp.VendorID)
	}
}

// Registering must not make an application look supported during capabilities
// exchange, which reads the application index.
func TestRegisterGroupedAVPDoesNotDeclareApplication(t *testing.T) {
	p := newRegisterTestParser(t)
	const unloadedApp = 16777251
	if _, err := p.App(unloadedApp); err != ErrApplicationUnsupported {
		t.Skipf("app %d already loaded in the test dictionary", unloadedApp)
	}
	if err := p.RegisterGroupedAVP(unloadedApp, 2000, 10415, "No-App-Grouped"); err != nil {
		t.Fatalf("RegisterGroupedAVP: %v", err)
	}
	if _, err := p.App(unloadedApp); err != ErrApplicationUnsupported {
		t.Errorf("App(%d) = supported after registering an AVP, want %v",
			unloadedApp, ErrApplicationUnsupported)
	}
}
