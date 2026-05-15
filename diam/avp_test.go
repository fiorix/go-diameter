// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

var testAVP = [][]byte{ // Body of a CER message
	{ // Origin-Host
		0x00, 0x00, 0x01, 0x08,
		0x40, 0x00, 0x00, 0x0e,
		0x63, 0x6c, 0x69, 0x65,
		0x6e, 0x74, 0x00, 0x00,
	},
	{ // Origin-Realm
		0x00, 0x00, 0x01, 0x28,
		0x40, 0x00, 0x00, 0x11,
		0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x6f, 0x73,
		0x74, 0x00, 0x00, 0x00,
	},
	{ // Host-IP-Address
		0x00, 0x00, 0x01, 0x01,
		0x40, 0x00, 0x00, 0x0e,
		0x00, 0x01, 0xc0, 0xa8,
		0xf2, 0x7a, 0x00, 0x00,
	},
	{ // Vendor-Id
		0x00, 0x00, 0x01, 0x0a,
		0x40, 0x00, 0x00, 0x0c,
		0x00, 0x00, 0x00, 0x0d,
	},
	{ // Product-Name
		0x00, 0x00, 0x01, 0x0d,
		0x40, 0x00, 0x00, 0x13,
		0x67, 0x6f, 0x2d, 0x64,
		0x69, 0x61, 0x6d, 0x65,
		0x74, 0x65, 0x72, 0x00,
	},
	{ // Origin-State-Id
		0x00, 0x00, 0x01, 0x16,
		0x40, 0x00, 0x00, 0x0c,
		0xe8, 0x3e, 0x3b, 0x84,
	},
}

func TestNewAVP(t *testing.T) {
	a := NewAVP(
		avp.OriginHost,                      // Code
		avp.Mbit,                            // Flags
		0,                                   // Vendor
		datatype.DiameterIdentity("foobar"), // Data
	)
	if a.Length != 14 { // Length in the AVP header
		t.Fatalf("Unexpected length. Want 14, have %d", a.Length)
	}
	if a.Len() != 16 { // With padding
		t.Fatalf("Unexpected length (with padding). Want 16, have %d\n", a.Len())
	}
	t.Log(a)
}

func TestDecodeAVP(t *testing.T) {
	a, err := DecodeAVP(testAVP[0], 1, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case a.Code != avp.OriginHost:
		t.Fatalf("Unexpected Code. Want %d, have %d", avp.OriginHost, a.Code)
	case a.Flags != avp.Mbit:
		t.Fatalf("Unexpected Code. Want %#x, have %#x", avp.Mbit, a.Flags)
	case a.Length != 14:
		t.Fatalf("Unexpected Length. Want 14, have %d", a.Length)
	case a.Data.Padding() != 2:
		t.Fatalf("Unexpected Padding. Want 2, have %d", a.Data.Padding())
	}
	t.Log(a)
}

func TestDecodeAVPMalformed(t *testing.T) {
	_, err := DecodeAVP(testAVP[0][:1], 1, dict.Default)
	if err == nil {
		t.Fatal("Malformed AVP decoded with no error")
	}
}

func TestDecodeAVPWithVendorID(t *testing.T) {
	var userNameVendorXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="1">
    <avp name="Session-Start-Indicator" code="1" vendor-id="999">
      <data type="UTF8String" />
    </avp>
  </application>
</diameter>`
	dict.Default.Load(bytes.NewReader([]byte(userNameVendorXML)))
	a := NewAVP(avp.UserName, avp.Mbit|avp.Vbit, 999, datatype.UTF8String("foobar"))
	b, err := a.Serialize()
	if err != nil {
		t.Fatal("Failed to serialize AVP:", err)
	}
	a, err = DecodeAVP(b, 1, dict.Default)
	if err != nil {
		t.Fatal("Failed to decode AVP:", err)
	}
	if a.VendorID != 999 {
		t.Fatalf("Unexpected VendorID. Want 999, have %d", a.VendorID)
	}
}

func TestEncodeAVP(t *testing.T) {
	a := &AVP{
		Code:  avp.OriginHost,
		Flags: avp.Mbit,
		Data:  datatype.DiameterIdentity("client"),
	}
	b, err := a.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testAVP[0]) {
		t.Fatalf("AVPs do not match.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testAVP[0]), hex.Dump(b))
	}
	t.Log(hex.Dump(b))
}

func TestEncodeAVPWithoutData(t *testing.T) {
	a := &AVP{
		Code:  avp.OriginHost,
		Flags: avp.Mbit,
	}
	_, err := a.Serialize()
	if err != nil {
		t.Log("Expected:", err)
	} else {
		t.Fatal("Unexpected serialization succeeded")
	}
}

func TestDecodeAVPVbitShortLength(t *testing.T) {
	// AVP with Vbit set but length < 12 should return error, not panic.
	data := []byte{
		0x00, 0x00, 0x00, 0x01, // Code
		0x80, 0x00, 0x00, 0x08, // Flags (Vbit set), Length=8
	}
	_, err := DecodeAVP(data, 1, dict.Default)
	if err == nil {
		t.Fatal("Expected error for AVP with Vbit and length < 12")
	}
}

func TestDecodeAVPEmptyPayloadNoVendor(t *testing.T) {
	// AVP with no vendor and length=8 (empty payload) should not crash.
	data := []byte{
		0x00, 0x00, 0x01, 0x08, // Code: Origin-Host (264)
		0x40, 0x00, 0x00, 0x08, // Flags (Mbit), Length=8
	}
	_, err := DecodeAVP(data, 1, dict.Default)
	// May return a decode error for empty data, but must not panic.
	_ = err
}

func TestDecodeAVPSizeMismatch(t *testing.T) {
	// VendorId (code 266) is Unsigned32 in the dictionary (expects 4 bytes).
	// This AVP has a 12-byte payload — a deliberate size mismatch.
	raw := []byte{
		0x00, 0x00, 0x01, 0x0a, // Code: 266 (VendorId)
		0x40, 0x00, 0x00, 0x14, // Flags: M-bit, Length: 20 (8 header + 12 payload)
		0x00, 0x00, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, // 12-byte payload
	}
	a, err := DecodeAVP(raw, 1, dict.Default)
	if err == nil {
		t.Fatal("Expected DecodeError for type size mismatch, got nil")
	}
	if _, ok := a.Data.(datatype.Unknown); !ok {
		t.Fatalf("Expected Unknown data type for mismatched AVP, got %T", a.Data)
	}
	if got := a.Len(); got != 20 {
		t.Fatalf("Expected Len()=20 (correct wire size), got %d", got)
	}
}

func TestDecodeAVPGroupedWithNonGroupedBytes(t *testing.T) {
	// Vendor-Specific-Application-Id (code 260) is Grouped in the dictionary.
	// Payload bytes have a Length field (0xff=255) that exceeds available data,
	// which is a fatal sub-AVP decode error.
	raw := []byte{
		0x00, 0x00, 0x01, 0x04, // Code: 260 (Vendor-Specific-Application-Id)
		0x40, 0x00, 0x00, 0x10, // Flags: M-bit, Length: 16 (8 header + 8 payload)
		0x01, 0x02, 0x03, 0x04, // garbage: sub-AVP code
		0x40, 0x00, 0x00, 0xff, // garbage: sub-AVP Length=255 > 4 remaining bytes
	}
	a, err := DecodeAVP(raw, 0, dict.Default)
	if err == nil {
		t.Fatal("Expected DecodeError for grouped mismatch, got nil")
	}
	if _, ok := a.Data.(datatype.Unknown); !ok {
		t.Fatalf("Expected Unknown data type for grouped mismatch AVP, got %T", a.Data)
	}
	if got := a.Len(); got != 16 {
		t.Fatalf("Expected Len()=16 (correct wire size), got %d", got)
	}
}

// testMismatchMessage is a raw CER message where the middle AVP (VendorId, code
// 266) carries a 12-byte payload instead of the expected 4 bytes (Unsigned32).
var testMismatchMessage = []byte{
	// Diameter header (20 bytes)
	0x01,             // Version
	0x00, 0x00, 0x48, // Message Length: 72
	0x80,             // Command Flags: Request
	0x00, 0x01, 0x01, // Command Code: 257 (CER)
	0x00, 0x00, 0x00, 0x00, // Application-ID: 0
	0x00, 0x00, 0x00, 0x01, // Hop-by-Hop ID
	0x00, 0x00, 0x00, 0x01, // End-to-End ID
	// AVP 1: Origin-Host (code 264), "client" — 16 bytes
	0x00, 0x00, 0x01, 0x08, // Code: 264
	0x40, 0x00, 0x00, 0x0e, // Flags: M, Length: 14 (8+6)
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x00, 0x00, // "client" + 2 padding
	// AVP 2: VendorId (code 266, Unsigned32 in dict) with 12-byte payload — 20 bytes
	0x00, 0x00, 0x01, 0x0a, // Code: 266
	0x40, 0x00, 0x00, 0x14, // Flags: M, Length: 20 (8+12)
	0x00, 0x00, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, // 12-byte payload (mismatch)
	// AVP 3: Origin-Realm (code 296), "realm" — 16 bytes
	0x00, 0x00, 0x01, 0x28, // Code: 296
	0x40, 0x00, 0x00, 0x0d, // Flags: M, Length: 13 (8+5)
	0x72, 0x65, 0x61, 0x6c, 0x6d, 0x00, 0x00, 0x00, // "realm" + 3 padding
}

func TestDecodeMessageWithSizeMismatchedAVP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMismatchMessage), dict.Default)
	if m == nil {
		t.Fatal("Expected non-nil message even with decode error")
	}
	if len(m.AVP) != 3 {
		t.Fatalf("Expected 3 AVPs, got %d", len(m.AVP))
	}
	if _, ok := m.AVP[1].Data.(datatype.Unknown); !ok {
		t.Fatalf("Expected Unknown type for mismatched VendorId AVP, got %T", m.AVP[1].Data)
	}
	if got := m.AVP[1].Len(); got != 20 {
		t.Fatalf("Expected mismatched AVP Len()=20, got %d", got)
	}
	if m.AVP[2].Code != avp.OriginRealm {
		t.Fatalf("Expected Origin-Realm after mismatched AVP, got code %d", m.AVP[2].Code)
	}
	if got := string(m.AVP[2].Data.(datatype.DiameterIdentity)); got != "realm" {
		t.Fatalf("Expected Origin-Realm=realm, got %q", got)
	}
}

func TestDecodeGroupedWithTruncatedSubAVP(t *testing.T) {
	// Auth-Application-Id sub-AVP (12 bytes) + 3 garbage bytes = 15 payload bytes.
	// Outer AVP: code 260 (Vendor-Specific-Application-Id, Grouped), Length=23.
	// One padding byte follows to align to 4 bytes.
	raw := []byte{
		0x00, 0x00, 0x01, 0x04, // Code: 260 (Vendor-Specific-Application-Id)
		0x40, 0x00, 0x00, 0x17, // Flags: M, Length: 23 (8 header + 15 payload)
		// Sub-AVP: Auth-Application-Id (code 258, Unsigned32=4) — 12 bytes
		0x00, 0x00, 0x01, 0x02, // Code: 258
		0x40, 0x00, 0x00, 0x0c, // Flags: M, Length: 12
		0x00, 0x00, 0x00, 0x04, // Value: 4
		// 3 trailing bytes — too few to form a sub-AVP header (need 8)
		0x01, 0x02, 0x03,
		0x00, // padding byte (not part of the AVP payload)
	}
	a, err := DecodeAVP(raw, 0, dict.Default)
	if err == nil {
		t.Fatal("Expected DecodeError for grouped with truncated sub-AVP, got nil")
	}
	if _, ok := a.Data.(datatype.Unknown); !ok {
		t.Fatalf("Expected Unknown data type, got %T", a.Data)
	}
	// Len() must equal 24: 8 header + 15 payload + 1 padding (pad4(23)=24)
	if got := a.Len(); got != 24 {
		t.Fatalf("Expected Len()=24, got %d", got)
	}
}

func BenchmarkDecodeAVP(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DecodeAVP(testAVP[0], 1, dict.Default)
	}
}

func BenchmarkEncodeAVP(b *testing.B) {
	a := NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("client"))
	for n := 0; n < b.N; n++ {
		a.Serialize()
	}
}
