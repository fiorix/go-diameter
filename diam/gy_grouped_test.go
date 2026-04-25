// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// buildGyCCA constructs a Gy Credit-Control-Answer with multiple
// Multiple-Services-Credit-Control AVPs containing deeply nested grouped
// structures (Granted-Service-Unit, Used-Service-Unit, CC-Money, Unit-Value).
func buildGyCCA() *Message {
	m := NewMessage(
		CreditControl,
		0, // Answer (no RequestFlag)
		CHARGING_CONTROL_APP_ID,
		0x1df4b75e,
		0xb1b50a4b,
		dict.Default,
	)

	// Mandatory top-level AVPs
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("session;1234567890"))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("ocs.example.com"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("example.com"))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(2)) // UPDATE_REQUEST
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))

	// MSCC #1: Data service with volume grants (3 levels deep)
	m.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(3600)),
					NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(1048576)),
					NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(524288)),
					NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(524288)),
				},
			}),
			NewAVP(avp.UsedServiceUnit, avp.Mbit, 0, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(1800)),
					NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(512000)),
					NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(256000)),
					NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(256000)),
				},
			}),
			NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(100)),
			NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(1000)),
			NewAVP(avp.ValidityTime, avp.Mbit, 0, datatype.Unsigned32(7200)),
			NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)),
		},
	})

	// MSCC #2: Monetary grant with CC-Money → Unit-Value (4 levels deep)
	m.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(7200)),
					NewAVP(avp.CCMoney, avp.Mbit, 0, &GroupedAVP{
						AVP: []*AVP{
							NewAVP(avp.UnitValue, avp.Mbit, 0, &GroupedAVP{
								AVP: []*AVP{
									NewAVP(avp.ValueDigits, avp.Mbit, 0, datatype.Integer64(50000)),
									NewAVP(avp.Exponent, avp.Mbit, 0, datatype.Integer32(-2)),
								},
							}),
							NewAVP(avp.CurrencyCode, avp.Mbit, 0, datatype.Unsigned32(978)), // EUR
						},
					}),
				},
			}),
			NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(200)),
			NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(2000)),
			NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)),
		},
	})

	// MSCC #3: Usage report only (no grant)
	m.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.UsedServiceUnit, avp.Mbit, 0, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(600)),
					NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(1024000)),
				},
			}),
			NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(300)),
			NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(3000)),
		},
	})

	return m
}

// TestGyNestedGroupedRoundTrip builds a Gy CCA with multiple MSCC AVPs
// containing up to 4 levels of nesting, serializes it, parses it back,
// and verifies every nested value survives the round trip.
func TestGyNestedGroupedRoundTrip(t *testing.T) {
	original := buildGyCCA()

	// Serialize to wire format.
	wireBytes, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Parse back from wire format.
	parsed, err := ReadMessage(bytes.NewReader(wireBytes), dict.Default)
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	// Verify header.
	if parsed.Header.CommandCode != CreditControl {
		t.Errorf("CommandCode: got %d, want %d", parsed.Header.CommandCode, CreditControl)
	}
	if parsed.Header.ApplicationID != CHARGING_CONTROL_APP_ID {
		t.Errorf("ApplicationID: got %d, want %d", parsed.Header.ApplicationID, CHARGING_CONTROL_APP_ID)
	}

	// Verify top-level AVP count matches.
	if len(parsed.AVP) != len(original.AVP) {
		t.Fatalf("AVP count: got %d, want %d", len(parsed.AVP), len(original.AVP))
	}

	// Re-serialize the parsed message and compare wire bytes.
	reserializedBytes, err := parsed.Serialize()
	if err != nil {
		t.Fatalf("Re-serialize failed: %v", err)
	}
	if !bytes.Equal(wireBytes, reserializedBytes) {
		t.Fatal("Round-trip wire bytes mismatch")
	}

	// --- Verify MSCC #1: Data service volume grant ---
	mscc1 := findGroupedAVP(t, parsed.AVP, avp.MultipleServicesCreditControl, 0)
	gsu1 := findGroupedAVP(t, mscc1, avp.GrantedServiceUnit, 0)
	assertUnsigned32(t, gsu1, avp.CCTime, 3600)
	assertUnsigned64(t, gsu1, avp.CCTotalOctets, 1048576)
	assertUnsigned64(t, gsu1, avp.CCInputOctets, 524288)
	assertUnsigned64(t, gsu1, avp.CCOutputOctets, 524288)

	usu1 := findGroupedAVP(t, mscc1, avp.UsedServiceUnit, 0)
	assertUnsigned32(t, usu1, avp.CCTime, 1800)
	assertUnsigned64(t, usu1, avp.CCTotalOctets, 512000)
	assertUnsigned64(t, usu1, avp.CCInputOctets, 256000)
	assertUnsigned64(t, usu1, avp.CCOutputOctets, 256000)

	assertUnsigned32(t, mscc1, avp.ServiceIdentifier, 100)
	assertUnsigned32(t, mscc1, avp.RatingGroup, 1000)
	assertUnsigned32(t, mscc1, avp.ValidityTime, 7200)
	assertUnsigned32(t, mscc1, avp.ResultCode, 2001)

	// --- Verify MSCC #2: Monetary grant (4 levels deep) ---
	mscc2 := findGroupedAVP(t, parsed.AVP, avp.MultipleServicesCreditControl, 1)
	gsu2 := findGroupedAVP(t, mscc2, avp.GrantedServiceUnit, 0)
	assertUnsigned32(t, gsu2, avp.CCTime, 7200)

	ccMoney := findGroupedAVP(t, gsu2, avp.CCMoney, 0)
	unitValue := findGroupedAVP(t, ccMoney, avp.UnitValue, 0)
	assertInteger64(t, unitValue, avp.ValueDigits, 50000)
	assertInteger32(t, unitValue, avp.Exponent, -2)
	assertUnsigned32(t, ccMoney, avp.CurrencyCode, 978)

	assertUnsigned32(t, mscc2, avp.ServiceIdentifier, 200)
	assertUnsigned32(t, mscc2, avp.RatingGroup, 2000)
	assertUnsigned32(t, mscc2, avp.ResultCode, 2001)

	// --- Verify MSCC #3: Usage report only ---
	mscc3 := findGroupedAVP(t, parsed.AVP, avp.MultipleServicesCreditControl, 2)
	usu3 := findGroupedAVP(t, mscc3, avp.UsedServiceUnit, 0)
	assertUnsigned32(t, usu3, avp.CCTime, 600)
	assertUnsigned64(t, usu3, avp.CCTotalOctets, 1024000)

	assertUnsigned32(t, mscc3, avp.ServiceIdentifier, 300)
	assertUnsigned32(t, mscc3, avp.RatingGroup, 3000)
}

// BenchmarkGyNestedGroupedRead benchmarks parsing a Gy CCA with multiple
// deeply nested MSCC grouped AVPs.
func BenchmarkGyNestedGroupedRead(b *testing.B) {
	original := buildGyCCA()
	wireBytes, err := original.Serialize()
	if err != nil {
		b.Fatal(err)
	}
	reader := bytes.NewReader(wireBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadMessage(reader, dict.Default)
		reader.Seek(0, 0)
	}
}

// BenchmarkGyNestedGroupedWrite benchmarks serializing a Gy CCA with
// multiple deeply nested MSCC grouped AVPs.
func BenchmarkGyNestedGroupedWrite(b *testing.B) {
	m := buildGyCCA()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Serialize()
	}
}

// --- Test helpers ---

// findGroupedAVP returns the child AVPs of the nth occurrence of a grouped
// AVP with the given code within avps. Fatals on missing or wrong type.
func findGroupedAVP(t *testing.T, avps []*AVP, code uint32, nth int) []*AVP {
	t.Helper()
	count := 0
	for _, a := range avps {
		if a.Code == code {
			if count == nth {
				g, ok := a.Data.(*GroupedAVP)
				if !ok {
					t.Fatalf("AVP code %d occurrence %d is not Grouped", code, nth)
				}
				return g.AVP
			}
			count++
		}
	}
	t.Fatalf("AVP code %d occurrence %d not found (found %d)", code, nth, count)
	return nil
}

// findAVPByCode returns the first AVP with the given code in avps.
func findAVPByCode(t *testing.T, avps []*AVP, code uint32) *AVP {
	t.Helper()
	for _, a := range avps {
		if a.Code == code {
			return a
		}
	}
	t.Fatalf("AVP code %d not found", code)
	return nil
}

func assertUnsigned32(t *testing.T, avps []*AVP, code uint32, want uint32) {
	t.Helper()
	a := findAVPByCode(t, avps, code)
	got, ok := a.Data.(datatype.Unsigned32)
	if !ok {
		t.Fatalf("AVP %d: expected Unsigned32, got %T", code, a.Data)
	}
	if uint32(got) != want {
		t.Errorf("AVP %d: got %d, want %d", code, got, want)
	}
}

func assertUnsigned64(t *testing.T, avps []*AVP, code uint32, want uint64) {
	t.Helper()
	a := findAVPByCode(t, avps, code)
	got, ok := a.Data.(datatype.Unsigned64)
	if !ok {
		t.Fatalf("AVP %d: expected Unsigned64, got %T", code, a.Data)
	}
	if uint64(got) != want {
		t.Errorf("AVP %d: got %d, want %d", code, got, want)
	}
}

func assertInteger32(t *testing.T, avps []*AVP, code uint32, want int32) {
	t.Helper()
	a := findAVPByCode(t, avps, code)
	got, ok := a.Data.(datatype.Integer32)
	if !ok {
		// Exponent is Integer32 in the dictionary, but check Enumerated too
		if e, eok := a.Data.(datatype.Enumerated); eok {
			if int32(e) != want {
				t.Errorf("AVP %d: got %d, want %d", code, e, want)
			}
			return
		}
		t.Fatalf("AVP %d: expected Integer32, got %T", code, a.Data)
	}
	if int32(got) != want {
		t.Errorf("AVP %d: got %d, want %d", code, got, want)
	}
}

func assertInteger64(t *testing.T, avps []*AVP, code uint32, want int64) {
	t.Helper()
	a := findAVPByCode(t, avps, code)
	got, ok := a.Data.(datatype.Integer64)
	if !ok {
		t.Fatalf("AVP %d: expected Integer64, got %T", code, a.Data)
	}
	if int64(got) != want {
		t.Errorf("AVP %d: got %d, want %d", code, got, want)
	}
}
