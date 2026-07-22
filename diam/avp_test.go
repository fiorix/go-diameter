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

func TestDecodeAVPUnknownVendor(t *testing.T) {
	// Unknown vendor AVP decodes as type Unknown with the full payload,
	// not bound to same-code base NAS-Port.
	payload := []byte("08-21090003e80000") // 17 bytes, not a 4-byte Unsigned32
	a := NewAVP(5, avp.Vbit, 10415, datatype.OctetString(payload))
	b, err := a.Serialize()
	if err != nil {
		t.Fatal("Failed to serialize AVP:", err)
	}
	got, err := DecodeAVP(b, 4, dict.Default)
	if err != nil {
		t.Fatal("Failed to decode unknown vendor AVP:", err)
	}
	if got == nil || got.Data == nil {
		t.Fatal("Expected Unknown AVP with non-nil Data")
	}
	if got.Len() != len(b) {
		t.Fatalf("Unknown AVP Len() %d does not match wire length %d (desync)",
			got.Len(), len(b))
	}
	if u, ok := got.Data.(datatype.Unknown); !ok {
		t.Fatalf("Expected datatype.Unknown, got %T", got.Data)
	} else if !bytes.Equal([]byte(u), payload) {
		t.Fatalf("Unknown payload not preserved: have %x, want %x", []byte(u), payload)
	}
}

func TestReadMessageUnknownVendorAVP(t *testing.T) {
	// Regression: an unknown vendor AVP colliding with base code 5 (NAS-Port) must
	// parse without desyncing the stream or panicking in (*AVP).Len().
	m := NewRequest(CreditControl, 4, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("test;1;0;1"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("client"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(5, avp.Vbit, 10415, datatype.OctetString([]byte("08-21090003e80000")))

	var buf bytes.Buffer
	if _, err := m.WriteTo(&buf); err != nil {
		t.Fatal("Failed to serialize message:", err)
	}
	wire := buf.Bytes()

	got, err := ReadMessage(bytes.NewReader(wire), dict.Default)
	if err != nil {
		t.Fatal("Failed to read message:", err)
	}
	if len(got.AVP) != 6 {
		t.Fatalf("Expected 6 AVPs, got %d (stream desync)", len(got.AVP))
	}
	last := got.AVP[len(got.AVP)-1]
	if last.Code != 5 || last.VendorID != 10415 {
		t.Fatalf("Last AVP misframed: code=%d vendor=%d", last.Code, last.VendorID)
	}
	if last.Data == nil {
		t.Fatal("Unknown AVP decoded with nil Data")
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
