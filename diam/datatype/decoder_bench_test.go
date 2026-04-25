// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "testing"

// decodeCase is a realistic (TypeID, payload) pair for benchmarking the
// shipped Decode function. Payloads are chosen to reflect what actually
// shows up on the wire in Diameter traffic: small numeric AVPs, variable
// length string / octet AVPs, Address AVPs with family prefix, and a
// Grouped payload large enough that the allocation inside DecodeGrouped
// is not dwarfed by a few bytes of stack work.
type decodeCase struct {
	name    string
	typeID  TypeID
	payload []byte
}

var decodeCases = []decodeCase{
	{"Unsigned32", Unsigned32Type, []byte{0x00, 0x00, 0x00, 0x2a}},
	{"Integer32", Integer32Type, []byte{0xff, 0xff, 0xff, 0xce}},
	{"Enumerated", EnumeratedType, []byte{0x00, 0x00, 0x00, 0x01}},
	{"Float32", Float32Type, []byte{0x40, 0x49, 0x0f, 0xdb}},
	{"AddressIPv4", AddressType, []byte{0x00, 0x01, 0x7f, 0x00, 0x00, 0x01}},
	{"AddressIPv6", AddressType, []byte{
		0x00, 0x02,
		0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}},
	{"OctetString", OctetStringType, []byte("example.com/octet")},
	{"UTF8String", UTF8StringType, []byte("subscriber@example.com")},
	{"DiameterIdentity", DiameterIdentityType, []byte("hss.example.com")},
	// Grouped payload sized like a small VendorSpecificApplicationID AVP
	// (~20 bytes of child AVP header+data). DecodeGrouped copies this,
	// so the benchmark captures the copy cost too.
	{"Grouped20B", GroupedType, make([]byte, 20)},
	{"Grouped256B", GroupedType, make([]byte, 256)},
}

// BenchmarkDecode exercises the shipped Decode function across the AVP
// data types that dominate real traffic. This is the regression gate: if
// Decode's dispatch order or fast-path changes, results here move.
//
// Use `-benchtime=3s` and benchstat across runs for stable comparisons.
func BenchmarkDecode(b *testing.B) {
	for _, tc := range decodeCases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := Decode(tc.typeID, tc.payload); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkDecodeMixed cycles through every case per iteration so the
// branch predictor cannot specialize on a single TypeID. Closer to what
// a real message decode loop does when AVP types alternate.
func BenchmarkDecodeMixed(b *testing.B) {
	n := len(decodeCases)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := decodeCases[i%n]
		if _, err := Decode(tc.typeID, tc.payload); err != nil {
			b.Fatal(err)
		}
	}
}

// The Dispatch* benchmarks isolate the cost of the dispatch mechanism
// itself — array index vs. map hash — with no decoder call, no allocation.
// They are not regression gates for Decode's shipped behavior; they exist
// to document the per-lookup delta between the two table types, which is
// the only thing that differs between them. Inlining the lookup here is
// intentional: calling Decode would bring in the actual decoder work and
// drown out the signal.
//
// BenchmarkDispatchMap uses a locally-constructed map rather than the
// exported Decoder map: Decoder is now extension-only and starts empty,
// so looking up built-in TypeIDs in it would measure a miss path instead
// of a real lookup. The local map is populated with every built-in
// TypeID so it has the same cardinality as decoderArray, keeping the
// comparison honest.

// dispatchTypes is the rotating set of TypeIDs both Dispatch* benchmarks
// look up, so only the table mechanism differs between them.
var dispatchTypes = []TypeID{Unsigned32Type, Integer32Type, EnumeratedType, Float32Type}

// dispatchMap mirrors decoderArray's contents in map form. Built at init
// time so the benchmark timer excludes the population cost.
var dispatchMap = func() map[TypeID]DecoderFunc {
	m := make(map[TypeID]DecoderFunc, maxTypeID)
	for t, f := range decoderArray {
		if f != nil {
			m[TypeID(t)] = f
		}
	}
	return m
}()

func BenchmarkDispatchArray(b *testing.B) {
	var sink DecoderFunc
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := dispatchTypes[i&3]
		if int(t) >= 0 && int(t) < len(decoderArray) {
			sink = decoderArray[t]
		}
	}
	_ = sink
}

func BenchmarkDispatchMap(b *testing.B) {
	var sink DecoderFunc
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := dispatchTypes[i&3]
		sink = dispatchMap[t]
	}
	_ = sink
}
