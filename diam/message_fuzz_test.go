// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// crafted malformed message that previously crashed the dictionary decoder
// with a nil-pointer dereference in (*AVP).Len() via decodeAVPs. A peer that
// can reach the server could send this to crash it (remote DoS), so the
// decoder must return cleanly instead of panicking. See decodeAVPs in
// message.go.
const malformedNilDataAVP = "0\x00\x01\x000\x00\x01\x02" +
	"00000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000"

// TestReadMessageMalformed feeds a table of malformed buffers through
// ReadMessage and asserts it never panics, regardless of whether the buffer
// is rejected with an error or decoded into a (possibly partial) message.
func TestReadMessageMalformed(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
	}{
		{"empty", nil},
		{"short-header", []byte{0x01, 0x00, 0x00}},
		{"nil-data-avp", []byte(malformedNilDataAVP)},
		// Valid-looking 20-byte header (CER, length 28) followed by an AVP
		// whose payload fails to decode, exercising the DecodeError path.
		{"avp-bad-payload", []byte{
			0x01, 0x00, 0x00, 0x1c, // version 1, length 28
			0x80, 0x00, 0x01, 0x01, // R flag, command 257 (CER)
			0x00, 0x00, 0x00, 0x00, // application 0
			0x00, 0x00, 0x00, 0x01, // hop-by-hop
			0x00, 0x00, 0x00, 0x01, // end-to-end
			0x00, 0x00, 0x01, 0x0a, // AVP code 266 (Vendor-Id, Unsigned32)
			0x00, 0x00, 0x00, 0x09, // V=0, length 9 (1 payload byte: too short)
		}},
		// AVP with a zero length field must not spin the decode loop.
		{"avp-zero-length", []byte{
			0x01, 0x00, 0x00, 0x1c,
			0x80, 0x00, 0x01, 0x01,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x01, 0x0a,
			0x00, 0x00, 0x00, 0x00, // length 0
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// A panic here fails the test via the testing runtime; we only
			// need to ensure the call returns.
			ReadMessage(bytes.NewReader(tc.in), dict.Default)
		})
	}
}

// fuzzSeeds returns the golden Gy fixtures used as fuzz corpus seeds, plus
// the known crash input. Missing fixtures are skipped rather than failing,
// so the test stays robust to fixture changes.
func fuzzSeeds(t testing.TB) [][]byte {
	seeds := [][]byte{[]byte(malformedNilDataAVP)}
	matches, _ := filepath.Glob(filepath.Join("..", "bench", "fixtures", "*.bin"))
	for _, m := range matches {
		b, err := os.ReadFile(m)
		if err != nil {
			t.Logf("skipping seed %s: %v", m, err)
			continue
		}
		seeds = append(seeds, b)
	}
	return seeds
}

// FuzzReadMessage asserts that ReadMessage never panics on arbitrary input.
// It is seeded with the real Gy fixtures and the known crash input. Run with:
//
//	go test ./diam/ -run=^$ -fuzz=FuzzReadMessage
func FuzzReadMessage(f *testing.F) {
	for _, seed := range fuzzSeeds(f) {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		// The contract under test: ReadMessage may return an error or a
		// message, but it must not panic on any input.
		ReadMessage(bytes.NewReader(in), dict.Default)
	})
}
