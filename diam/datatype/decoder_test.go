// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"testing"
)

// TestDecodeCustomDecoderFallback verifies that a decoder registered in the
// exported Decoder map for a TypeID outside the built-in range is resolved
// via the map fallback in Decode.
func TestDecodeCustomDecoderFallback(t *testing.T) {
	const customType TypeID = maxTypeID + 1
	if _, ok := Decoder[customType]; ok {
		t.Fatalf("precondition failed: TypeID %d already registered", customType)
	}

	want := []byte{0xde, 0xad, 0xbe, 0xef}
	Decoder[customType] = func(b []byte) (Type, error) {
		return OctetString(b), nil
	}
	defer delete(Decoder, customType)

	got, err := Decode(customType, want)
	if err != nil {
		t.Fatalf("Decode returned error for custom TypeID: %v", err)
	}
	if !bytes.Equal(got.Serialize(), want) {
		t.Fatalf("custom decoder not invoked. want %x, have %x", want, got.Serialize())
	}
}

// TestDecodeBuiltinFastPath verifies that a built-in TypeID still dispatches
// through the array fast path and returns the correct data type.
func TestDecodeBuiltinFastPath(t *testing.T) {
	b := []byte{0x00, 0x00, 0x00, 0x2a}
	got, err := Decode(Unsigned32Type, b)
	if err != nil {
		t.Fatal(err)
	}
	if got.Type() != Unsigned32Type {
		t.Fatalf("unexpected type. want %d, have %d", Unsigned32Type, got.Type())
	}
	if n := uint32(got.(Unsigned32)); n != 42 {
		t.Fatalf("unexpected value. want 42, have %d", n)
	}
}

// TestDecodeUnknownTypeID verifies that an unregistered TypeID returns an
// error rather than silently dispatching.
func TestDecodeUnknownTypeID(t *testing.T) {
	const unknownType TypeID = maxTypeID + 99
	if _, ok := Decoder[unknownType]; ok {
		t.Fatalf("precondition failed: TypeID %d is registered", unknownType)
	}
	if _, err := Decode(unknownType, nil); err == nil {
		t.Fatal("expected error for unregistered TypeID, got nil")
	}
}
