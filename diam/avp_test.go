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

var testAVP = [][]byte{ // Body of a CER message
	[]byte{ // Origin-Host
		0x00, 0x00, 0x01, 0x08,
		0x40, 0x00, 0x00, 0x0e,
		0x63, 0x6c, 0x69, 0x65,
		0x6e, 0x74, 0x00, 0x00,
	},
	[]byte{ // Origin-Realm
		0x00, 0x00, 0x01, 0x28,
		0x40, 0x00, 0x00, 0x11,
		0x6c, 0x6f, 0x63, 0x61,
		0x6c, 0x68, 0x6f, 0x73,
		0x74, 0x00, 0x00, 0x00,
	},
	[]byte{ // Host-IP-Address
		0x00, 0x00, 0x01, 0x01,
		0x40, 0x00, 0x00, 0x0e,
		0x00, 0x01, 0xc0, 0xa8,
		0xf2, 0x7a, 0x00, 0x00,
	},
	[]byte{ // Vendor-Id
		0x00, 0x00, 0x01, 0x0a,
		0x40, 0x00, 0x00, 0x0c,
		0x00, 0x00, 0x00, 0x0d,
	},
	[]byte{ // Product-Name
		0x00, 0x00, 0x01, 0x0d,
		0x40, 0x00, 0x00, 0x13,
		0x67, 0x6f, 0x2d, 0x64,
		0x69, 0x61, 0x6d, 0x65,
		0x74, 0x65, 0x72, 0x00,
	},
	[]byte{ // Origin-State-Id
		0x00, 0x00, 0x01, 0x16,
		0x40, 0x00, 0x00, 0x0c,
		0xe8, 0x3e, 0x3b, 0x84,
	},
}

func TestNewAVP(t *testing.T) {
	a := NewAVP(264, 0x40, 0, diamtype.DiameterIdentity("client"))
	if a.Length != 14 { // Length in the AVP header
		t.Fatalf("Unexpected length. Want 14, have %d", a.Length)
	}
	if a.Len() != 16 { // With padding
		t.Fatalf("Unexpected length (with padding). Want 16, have %d\n", a.Len())
	}
	t.Log(a)
}

func TestDecodeAVP(t *testing.T) {
	avp, err := decodeAVP(testAVP[0], 1, diamdict.Default)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case avp.Code != 264:
		t.Fatalf("Unexpected Code. Want 264, have %d", avp.Code)
	case avp.Flags != 0x40:
		t.Fatalf("Unexpected Code. Want 0x40, have 0x%x", avp.Flags)
	case avp.Length != 14:
		t.Fatalf("Unexpected Length. Want 14, have %d", avp.Length)
	case avp.Data.Padding() != 2:
		t.Fatalf("Unexpected Padding. Want 2, have %d", avp.Data.Padding())
	}
	t.Log(avp)
}

func TestEncodeAVP(t *testing.T) {
	avp := &AVP{
		Code:  264,
		Flags: 0x40,
		Data:  diamtype.DiameterIdentity("client"),
	}
	b, err := avp.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, testAVP[0]) {
		t.Fatalf("AVPs do not match.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testAVP[0]), hex.Dump(b))
	}
	t.Log(hex.Dump(b))
}

func TestEncodeEmptyAVP(t *testing.T) {
	avp := &AVP{
		Code:  264,
		Flags: 0x40,
	}
	_, err := avp.Serialize()
	if err != nil {
		t.Log("Expected:", err)
	} else {
		t.Fatal("Unexpected serialization succeeded")
	}
}

func BenchmarkDecodeAVP(b *testing.B) {
	for n := 0; n < b.N; n++ {
		decodeAVP(testAVP[0], 1, diamdict.Default)
	}
}

func BenchmarkEncodeAVP(b *testing.B) {
	a := NewAVP(264, 0x40, 0, diamtype.DiameterIdentity("client"))
	for n := 0; n < b.N; n++ {
		a.Serialize()
	}
}
