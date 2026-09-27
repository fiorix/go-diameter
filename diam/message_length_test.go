// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diam

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// A well formed 20 byte header whose Message Length field is 0.
//
// Version 1, Request flag set, Command Code 257 and Application ID 0 — that is
// CER, the first message of any Diameter connection. The command and
// application both exist in the base dictionary, which is the only thing
// readHeader checks before handing the declared length to readBody, so this
// reaches the allocation before capabilities are exchanged and before any
// handler runs.
var messageLengthUnderflow = []byte{
	0x01, 0x00, 0x00, 0x00, // version 1, Message Length 0
	0x80, 0x00, 0x01, 0x01, // Request, Command Code 257 (CER)
	0x00, 0x00, 0x00, 0x00, // Application ID 0
	0x00, 0x00, 0x00, 0x01, // Hop-by-Hop
	0x00, 0x00, 0x00, 0x01, // End-to-End
}

// TestDecodeHeaderRejectsLengthBelowHeader is the unit half of the fix.
//
// RFC 6733 section 3 defines Message Length as the length of the message
// including the header, so anything below HeaderLength describes no message
// that can exist.
func TestDecodeHeaderRejectsLengthBelowHeader(t *testing.T) {
	for _, declared := range []uint32{0, 1, 19} {
		data := append([]byte(nil), messageLengthUnderflow...)
		data[1] = byte(declared >> 16)
		data[2] = byte(declared >> 8)
		data[3] = byte(declared)

		if _, err := DecodeHeader(data); err == nil {
			t.Errorf("a Message Length of %d was accepted; it is below the %d byte header",
				declared, HeaderLength)
		}
	}

	// The boundary is inclusive: a bodyless message is exactly HeaderLength.
	data := append([]byte(nil), messageLengthUnderflow...)
	data[3] = HeaderLength
	if _, err := DecodeHeader(data); err != nil {
		t.Errorf("a Message Length of exactly %d was rejected: %v", HeaderLength, err)
	}
}

// TestReadMessageDoesNotAllocateOnLengthUnderflow is the one that matters.
//
// readBody computes int(Header.MessageLength - HeaderLength) in unsigned
// arithmetic. Without the check in DecodeFromBytes, a declared length of 0
// makes that uint32(0)-20, or 4294967276, and readerBufferSlice make()s it: 20
// bytes on the wire become a 4 GiB allocation, measured at 4294970504 bytes and
// 12.5 seconds on an unpatched build of this commit's parent. Under a container
// memory limit it is a fatal out-of-memory rather than a slow decode — and an
// out-of-memory is a fatal runtime error, so no recover() above it can help.
//
// Asserted on allocation rather than only on the error, because an error is
// returned either way: the unpatched path fails with "readBody Error: EOF"
// *after* allocating. The allocation is the defect.
func TestReadMessageDoesNotAllocateOnLengthUnderflow(t *testing.T) {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	m, err := ReadMessage(bytes.NewReader(messageLengthUnderflow), dict.Default)

	runtime.ReadMemStats(&after)
	allocated := after.TotalAlloc - before.TotalAlloc

	if err == nil {
		t.Fatalf("a Message Length of 0 decoded into %v", m)
	}
	// Three orders of magnitude below the defect and three above what the
	// patched path actually uses (about 2 KiB), so this fails on the bug
	// returning and not on an unrelated allocation moving.
	if allocated > 1<<20 {
		t.Fatalf("decoding %d bytes allocated %d, which is the length underflow returning",
			len(messageLengthUnderflow), allocated)
	}
}
