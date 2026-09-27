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

// A well formed 20 byte header whose Message Length field is 0. It declares
// a CER on application 0, both in the base dictionary, so readHeader hands
// the declared length to readBody before capabilities are exchanged and
// before any handler runs.
var messageLengthUnderflow = []byte{
	0x01, 0x00, 0x00, 0x00, // version 1, Message Length 0
	0x80, 0x00, 0x01, 0x01, // Request, Command Code 257 (CER)
	0x00, 0x00, 0x00, 0x00, // Application ID 0
	0x00, 0x00, 0x00, 0x01, // Hop-by-Hop
	0x00, 0x00, 0x00, 0x01, // End-to-End
}

// RFC 6733 section 3 defines Message Length as the length of the message
// including the header, so anything below HeaderLength is invalid.
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

// readBody computes int(Header.MessageLength - HeaderLength) in unsigned
// arithmetic, so without the check in DecodeFromBytes a declared length of 0
// turns 20 bytes on the wire into a near-4 GiB allocation. Asserted on
// allocation rather than only on the error, because the unpatched path also
// returns an error, after allocating.
func TestReadMessageDoesNotAllocateOnLengthUnderflow(t *testing.T) {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	m, err := ReadMessage(bytes.NewReader(messageLengthUnderflow), dict.Default)

	runtime.ReadMemStats(&after)
	allocated := after.TotalAlloc - before.TotalAlloc

	if err == nil {
		t.Fatalf("a Message Length of 0 decoded into %v", m)
	}
	// Far below the near-4 GiB defect and far above the roughly 2 KiB the
	// fixed path uses, so this fails only on the underflow returning.
	if allocated > 1<<20 {
		t.Fatalf("decoding %d bytes allocated %d, which is the length underflow returning",
			len(messageLengthUnderflow), allocated)
	}
}
