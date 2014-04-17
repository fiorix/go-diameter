// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/diam/dict"
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
	a, err := decodeAVP(testGroupedAVP, 0, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	if a.Data.Type() != GroupedType {
		t.Fatal("AVP is not grouped")
	}
	t.Log(a)
	s, err := a.Serialize()
	if !bytes.Equal(s, testGroupedAVP) {
		t.Fatalf("Unexpected value.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testGroupedAVP), hex.Dump(s))
	}
}
