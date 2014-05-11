// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/diam/diamdict"
	"github.com/fiorix/go-diameter/diam/diamtype"
)

var testGroupedAVP = []byte{ // Vendor-Specific-Application-Id
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
	a, err := decodeAVP(testGroupedAVP, 0, diamdict.Default)
	if err != nil {
		t.Fatal(err)
	}
	if a.Data.Type() != GroupedType {
		t.Fatal("AVP is not grouped")
	}
	b, err := a.Serialize()
	if !bytes.Equal(b, testGroupedAVP) {
		t.Fatalf("Unexpected value.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testGroupedAVP), hex.Dump(b))
	}
	t.Log(a)
}

func TestDecodeMessageWithGroupedAVP(t *testing.T) {
	m := NewRequest(257, 0, diamdict.Default)
	m.NewAVP(264, 0x40, 0, diamtype.DiameterIdentity("client"))
	a, _ := decodeAVP(testGroupedAVP, 0, diamdict.Default)
	m.AddAVP(a)
	t.Log(m)
}

func TestMakeGroupedAVP(t *testing.T) {
	// Create empty Grouped AVP
	g := new(Grouped)
	// Add Auth-Application-Id
	g.AVP = append(g.AVP, NewAVP(258, 0x40, 0, diamtype.Unsigned32(4)))
	// Add Vendor-Id
	g.AVP = append(g.AVP, NewAVP(266, 0x40, 0, diamtype.Unsigned32(10415)))
	// Create Vendor-Specific-Application-Id
	a := NewAVP(260, 0x40, 0, g)
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
