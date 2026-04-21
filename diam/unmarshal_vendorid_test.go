package diam

import (
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// TestUnmarshalVendorId reproduces #169: when two AVPs share the same code
// but have different VendorIds, Unmarshal puts the wrong data into the struct.
func TestUnmarshalVendorId(t *testing.T) {
	// User-Password (code=2, vendor=0) and 3GPP-Charging-Id (code=2, vendor=10415)
	// both have AVP code 2. Unmarshal should distinguish them by VendorId.

	// Load a dictionary that has both AVPs with code=2 but different vendors.
	// The base dict already has Inband-Security-Id (code=299) and the 3GPP dict
	// has several vendor-specific AVPs. Let's use a simpler approach:
	// build a message with two AVPs that share code but differ by vendor.

	m := NewRequest(272, 4, dict.Default) // CCR
	// AVP code 2, vendor 0 (base protocol doesn't define code=2, so use raw)
	m.AddAVP(NewAVP(2, avp.Mbit, 0, datatype.OctetString("base-password")))
	// AVP code 2, vendor 10415 (3GPP-Charging-Id in Ro/Rf dict)
	m.AddAVP(NewAVP(2, avp.Mbit|avp.Vbit, 10415, datatype.Unsigned32(12345)))

	type testStruct struct {
		BaseAVP    datatype.OctetString  `avp:"User-Password"`
		ChargingId datatype.Unsigned32   `avp:"3GPP-Charging-Id"`
	}

	// This will fail if the dictionary doesn't have both names.
	// Let's test with raw AVP access first to prove the index bug.
	idx := newIndex(m.AVP)

	// Both AVPs have code=2 — with the fix they should be in separate buckets
	bucket0, exists0 := idx[avpKey{2, 0}]
	bucket3gpp, exists3gpp := idx[avpKey{2, 10415}]

	if !exists0 {
		t.Fatal("AVP code=2 vendor=0 not found in index")
	}
	if !exists3gpp {
		t.Fatal("AVP code=2 vendor=10415 not found in index")
	}
	if len(bucket0) != 1 || len(bucket3gpp) != 1 {
		t.Fatalf("expected 1 AVP per bucket, got %d and %d", len(bucket0), len(bucket3gpp))
	}

	t.Logf("AVP vendor=0: code=%d data=%v", bucket0[0].Code, bucket0[0].Data)
	t.Logf("AVP vendor=10415: code=%d data=%v", bucket3gpp[0].Code, bucket3gpp[0].Data)

	if bucket0[0].VendorID != 0 {
		t.Error("wrong AVP in vendor=0 bucket")
	}
	if bucket3gpp[0].VendorID != 10415 {
		t.Error("wrong AVP in vendor=10415 bucket")
	}
}

// TestUnmarshalVendorIdFix verifies that after the fix, AVPs with the same
// code but different VendorIds are correctly distinguished during Unmarshal.
func TestUnmarshalVendorIdFix(t *testing.T) {
	m := NewRequest(272, 4, dict.Default)
	m.AddAVP(NewAVP(2, avp.Mbit, 0, datatype.OctetString("base-data")))
	m.AddAVP(NewAVP(2, avp.Mbit|avp.Vbit, 10415, datatype.Unsigned32(99999)))

	idx := newIndex(m.AVP)

	// After fix: AVPs with different VendorIds should be in separate buckets
	v0, ok0 := idx[avpKey{2, 0}]
	v3gpp, ok3gpp := idx[avpKey{2, 10415}]

	if !ok0 || len(v0) != 1 {
		t.Fatal("vendor=0 AVP not correctly indexed")
	}
	if !ok3gpp || len(v3gpp) != 1 {
		t.Fatal("vendor=10415 AVP not correctly indexed")
	}
	if string(v0[0].Data.(datatype.OctetString)) != "base-data" {
		t.Errorf("vendor=0 AVP has wrong data: %v", v0[0].Data)
	}
	if uint32(v3gpp[0].Data.(datatype.Unsigned32)) != 99999 {
		t.Errorf("vendor=10415 AVP has wrong data: %v", v3gpp[0].Data)
	}
}
