// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// testGroupedAVP is a Vendor-Specific-Application-Id Grouped AVP.
var testGroupedAVP = []byte{
	0x00, 0x00, 0x01, 0x04,
	0x40, 0x00, 0x00, 0x20,
	0x00, 0x00, 0x01, 0x02, // Auth-Application-Id
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x00, 0x04,
	0x00, 0x00, 0x01, 0x0a, // Vendor-Id
	0x40, 0x00, 0x00, 0x0c,
	0x00, 0x00, 0x28, 0xaf,
}

func TestGroupedAVP(t *testing.T) {
	a, err := DecodeAVP(testGroupedAVP, 0, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	if a.Data.Type() != GroupedAVPType {
		t.Fatal("AVP is not grouped")
	}
	b, err := a.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testGroupedAVP) {
		t.Fatalf("Unexpected value.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testGroupedAVP), hex.Dump(b))
	}
	t.Log(a)
}

func TestDecodeMessageWithGroupedAVP(t *testing.T) {
	m := NewRequest(257, 0, dict.Default)
	m.NewAVP(264, 0x40, 0, datatype.DiameterIdentity("client"))
	a, _ := DecodeAVP(testGroupedAVP, 0, dict.Default)
	m.AddAVP(a)
	t.Logf("Message:\n%s", m)
}

func TestDecodeGroupedFromBytesHeaderTooShort(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03} // 3 bytes — cannot form an AVP header
	g, err := DecodeGroupedFromBytes(b, 0, dict.Default)
	if err == nil {
		t.Fatal("Expected error for truncated header, got nil")
	}
	if len(g.AVP) != 0 {
		t.Fatalf("Expected 0 sub-AVPs, got %d", len(g.AVP))
	}
}

func TestDecodeGroupedFromBytesDataTooShort(t *testing.T) {
	b := []byte{
		0x01, 0x02, 0x03, 0x04, // Code
		0x40,             // Flags: M-bit
		0x00, 0x00, 0xff, // Length: 255 — far exceeds the 8 bytes available
	}
	g, err := DecodeGroupedFromBytes(b, 0, dict.Default)
	if err == nil {
		t.Fatal("Expected error for data-too-short sub-AVP, got nil")
	}
	if len(g.AVP) != 0 {
		t.Fatalf("Expected 0 sub-AVPs (break before append), got %d", len(g.AVP))
	}
}

func TestDecodeGroupedFromBytesValidThenTruncated(t *testing.T) {
	b := []byte{
		// Auth-Application-Id (code 258, Unsigned32=4) — 12 bytes, valid
		0x00, 0x00, 0x01, 0x02, 0x40, 0x00, 0x00, 0x0c,
		0x00, 0x00, 0x00, 0x04,
		// 3 trailing bytes — too few to form a sub-AVP header (need 8)
		0x01, 0x02, 0x03,
	}
	g, err := DecodeGroupedFromBytes(b, 0, dict.Default)
	if err == nil {
		t.Fatal("Expected error for truncated trailing bytes, got nil")
	}
	if len(g.AVP) != 1 {
		t.Fatalf("Expected 1 decoded sub-AVP before the truncated remainder, got %d", len(g.AVP))
	}
	if g.AVP[0].Code != avp.AuthApplicationID {
		t.Fatalf("Expected Auth-Application-Id (code %d), got %d", avp.AuthApplicationID, g.AVP[0].Code)
	}
	// g.Len() reflects only the successfully decoded sub-AVP.
	if got := g.Len(); got != 12 {
		t.Fatalf("Expected g.Len()=12, got %d", got)
	}
}

func TestMakeGroupedAVP(t *testing.T) {
	g := &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
			NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
		},
	}
	a := NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, g)
	b, err := a.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testGroupedAVP) {
		t.Fatalf("Unexpected value.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testGroupedAVP), hex.Dump(b))
	}
	t.Logf("Message:\n%s", a)
}

// TestDecodeGroupedSurvivesAnUndecodableAVP covers the same defect as
// TestDecodeAVPsSurvivesAnUndecodableAVP reached through a grouped AVP, which
// is how the interesting messages arrive: Media-Component-Description on Rx
// nests to arbitrary depth. Minimised by a fuzzer to 60 bytes.
func TestDecodeGroupedSurvivesAnUndecodableAVP(t *testing.T) {
	message := []byte("0\x00\x0070\x00\x01\x12\x01\x00\x00100000000\x00\x00\x01)0" +
		"\x00\x00 00000000000000000000000000000000000000000000")

	m, err := ReadMessage(bytes.NewReader(message), dict.Default)
	if err == nil {
		t.Fatal("Expected a decode error for a grouped AVP with an undecodable member, got nil")
	}
	if !strings.Contains(err.Error(), "Failed to decode one or more AVPs") {
		t.Fatalf("Unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("Expected a non-nil *Message alongside the decode error")
	}
}

// nestedMSCC returns a Credit-Control request whose body is depth
// Multiple-Services-Credit-Control AVPs, each containing the next.
func nestedMSCC(depth int) []byte {
	total := HeaderLength + 8*depth
	b := make([]byte, total)
	b[0] = 1
	putUint24(b[1:4], uint32(total))
	b[4] = RequestFlag
	putUint24(b[5:8], CreditControl)
	binary.BigEndian.PutUint32(b[8:12], CHARGING_CONTROL_APP_ID)
	for i := 0; i < depth; i++ {
		o := HeaderLength + 8*i
		binary.BigEndian.PutUint32(b[o:o+4], avp.MultipleServicesCreditControl)
		b[o+4] = avp.Mbit
		putUint24(b[o+5:o+8], uint32(total-o))
	}
	return b
}

func TestDecodeGroupedNestingLimit(t *testing.T) {
	for _, tc := range []struct {
		name    string
		depth   int
		wantErr bool
	}{
		{"at limit", dict.DefaultMaxGroupedDepth, false},
		{"one past limit", dict.DefaultMaxGroupedDepth + 1, true},
		{"hostile", 20000, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadMessage(bytes.NewReader(nestedMSCC(tc.depth)), dict.Default)
			if (err != nil) != tc.wantErr {
				t.Fatalf("depth %d: err = %v, wantErr %v", tc.depth, err, tc.wantErr)
			}
		})
	}
}

func TestDecodeGroupedNestingLimitFromParser(t *testing.T) {
	newParser := func(limit int) *dict.Parser {
		p, err := dict.NewParser("./dict/testdata/base.xml", "./dict/testdata/credit_control.xml")
		if err != nil {
			t.Fatal(err)
		}
		p.MaxGroupedDepth = limit
		return p
	}
	for _, tc := range []struct {
		name    string
		limit   int
		depth   int
		wantErr bool
	}{
		{"raised, at limit", 64, 64, false},
		{"raised, one past limit", 64, 65, true},
		{"lowered, at limit", 4, 4, false},
		{"lowered, one past limit", 4, 5, true},
		{"zero means default", 0, dict.DefaultMaxGroupedDepth + 1, true},
		{"negative means default", -1, dict.DefaultMaxGroupedDepth, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadMessage(bytes.NewReader(nestedMSCC(tc.depth)), newParser(tc.limit))
			if (err != nil) != tc.wantErr {
				t.Fatalf("limit %d, depth %d: err = %v, wantErr %v", tc.limit, tc.depth, err, tc.wantErr)
			}
		})
	}
}
