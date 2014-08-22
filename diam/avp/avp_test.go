// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diamtype"
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

func TestNew(t *testing.T) {
	a := New(
		OriginHost, // Code
		Mbit,       // Flags
		0,          // Vendor
		diamtype.DiameterIdentity("foobar"), // Data
	)
	if a.Length != 14 { // Length in the AVP header
		t.Fatalf("Unexpected length. Want 14, have %d", a.Length)
	}
	if a.Len() != 16 { // With padding
		t.Fatalf("Unexpected length (with padding). Want 16, have %d\n", a.Len())
	}
	t.Log(a)
}

func TestDecode(t *testing.T) {
	avp, err := Decode(testAVP[0], 1, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case avp.Code != OriginHost:
		t.Fatalf("Unexpected Code. Want %d, have %d", OriginHost, avp.Code)
	case avp.Flags != Mbit:
		t.Fatalf("Unexpected Code. Want %#x, have %#x", Mbit, avp.Flags)
	case avp.Length != 14:
		t.Fatalf("Unexpected Length. Want 14, have %d", avp.Length)
	case avp.Data.Padding() != 2:
		t.Fatalf("Unexpected Padding. Want 2, have %d", avp.Data.Padding())
	}
	t.Log(avp)
}

func TestEncode(t *testing.T) {
	avp := &AVP{
		Code:  OriginHost,
		Flags: Mbit,
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

func TestEncodeWithoutData(t *testing.T) {
	avp := &AVP{
		Code:  OriginHost,
		Flags: Mbit,
	}
	_, err := avp.Serialize()
	if err != nil {
		t.Log("Expected:", err)
	} else {
		t.Fatal("Unexpected serialization succeeded")
	}
}

func BenchmarkDecode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Decode(testAVP[0], 1, dict.Default)
	}
}

func BenchmarkEncode(b *testing.B) {
	a := New(OriginHost, Mbit, 0, diamtype.DiameterIdentity("client"))
	for n := 0; n < b.N; n++ {
		a.Serialize()
	}
}
