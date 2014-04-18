// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var testHeader = []byte{
	0x01, 0x00, 0x00, 0x74,
	0x80, 0x00, 0x01, 0x01,
	0x00, 0x00, 0x00, 0x01,
	0x2c, 0x0b, 0x61, 0x49,
	0xdb, 0xbf, 0xd3, 0x85,
}

func TestDecodeHeader(t *testing.T) {
	hdr, err := decodeHeader(testHeader)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hdr)
	switch {
	case hdr.Version != 1:
		t.Fatalf("Unexpected Version. Want 1, have %d", hdr.Version)
	case hdr.MessageLength != 116:
		t.Fatalf("Unexpected MessageLength. Want 116, have %d", hdr.MessageLength)
	case hdr.CommandFlags != 0x80:
		t.Fatalf("Unexpected CommandFlags. Want 0x80, have 0x%x", hdr.CommandFlags)
	case hdr.CommandCode != 257:
		t.Fatalf("Unexpected CommandCode. Want 257, have %d", hdr.CommandCode)
	case hdr.ApplicationId != 1:
		t.Fatalf("Unexpected ApplicationId. Want 1, have %d", hdr.ApplicationId)
	case hdr.HopByHopId != 0x2c0b6149:
		t.Fatalf("Unexpected HopByHopId. Want 0x2c0b6149, have 0x%x", hdr.HopByHopId)
	case hdr.EndToEndId != 0xdbbfd385:
		t.Fatalf("Unexpected EndToEndId. Want 0xdbbf0385, have 0x%x", hdr.EndToEndId)
	}
}

func TestEncodeHeader(t *testing.T) {
	hdr := &Header{
		Version:       1,
		MessageLength: 116,
		CommandFlags:  0x80,
		CommandCode:   257,
		ApplicationId: 1,
		HopByHopId:    0x2c0b6149,
		EndToEndId:    0xdbbfd385,
	}
	b := hdr.Serialize()
	if !bytes.Equal(testHeader, b) {
		t.Fatalf("Unexpected packet.\nWant:\n%s\nHave:\n%s",
			hex.Dump(testHeader), hex.Dump(b))
	}
}

func BenchmarkDecodeHeader(b *testing.B) {
	for n := 0; n < b.N; n++ {
		decodeHeader(testHeader)
	}
}

func BenchmarkEncodeHeader(b *testing.B) {
	hdr := &Header{
		Version:       1,
		MessageLength: 116,
		CommandFlags:  0x80,
		CommandCode:   257,
		ApplicationId: 1,
		HopByHopId:    0x2c0b6149,
		EndToEndId:    0xdbbfd385,
	}
	for n := 0; n < b.N; n++ {
		hdr.Serialize()
	}
}
