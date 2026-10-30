package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// testAIAWithOCOLR is an S6a Authentication-Information-Answer that carries
// both an RFC 7683 OC-OLR and a UTRAN authentication vector. The two share
// AVP codes: OC-Validity-Duration and Confidentiality-Key are both code 625,
// and OC-Report-Type and Integrity-Key are both code 626, told apart only by
// Vendor-Id (0 vs 10415).
var testAIAWithOCOLR, _ = hex.DecodeString("" +
	"010000a00000013e010000230102030405060708" + // header: AIA, app 16777251
	"0000010c4000000c000007d1" + // Result-Code 2001
	"0000026f00000030" + // OC-OLR (623)
	"0000027000000010" + "0000000000000001" + // OC-Sequence-Number 1
	"000002720000000c" + "00000000" + // OC-Report-Type (626) HOST_REPORT
	"000002710000000c" + "0000001e" + // OC-Validity-Duration (625) 30
	"00000585c0000050000028af" + // Authentication-Info (1413, 10415)
	"00000587c0000044000028af" + // UTRAN-Vector (1415, 10415)
	"00000271c000001c000028af" + "cccccccccccccccccccccccccccccccc" + // Confidentiality-Key (625, 10415)
	"00000272c000001c000028af" + "11111111111111111111111111111111") // Integrity-Key (626, 10415)

// TestMessageFindAVPVendorID checks that FindAVP, FindAVPs and
// FindAVPsWithPath return only AVPs with the Vendor-Id of the AVP they
// were asked for, not any AVP that shares its code.
func TestMessageFindAVPVendorID(t *testing.T) {
	m, err := ReadMessage(bytes.NewReader(testAIAWithOCOLR), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	ck := datatype.OctetString(bytes.Repeat([]byte{0xcc}, 16))
	ik := datatype.OctetString(bytes.Repeat([]byte{0x11}, 16))

	for _, tc := range []struct {
		code     interface{}
		vendorID uint32
		want     datatype.Type
	}{
		{"Confidentiality-Key", 10415, ck},
		{"Integrity-Key", 10415, ik},
		{"OC-Validity-Duration", 0, datatype.Unsigned32(30)},
		{"OC-Report-Type", 0, datatype.Enumerated(0)},
	} {
		a, err := m.FindAVP(tc.code, tc.vendorID)
		if err != nil {
			t.Errorf("FindAVP(%v, %d): %v", tc.code, tc.vendorID, err)
			continue
		}
		if a.VendorID != tc.vendorID || a.Data != tc.want {
			t.Errorf("FindAVP(%v, %d) = %v, want Vendor-Id %d and %v",
				tc.code, tc.vendorID, a, tc.vendorID, tc.want)
		}
		avps, err := m.FindAVPs(tc.code, tc.vendorID)
		if err != nil || len(avps) != 1 || avps[0] != a {
			t.Errorf("FindAVPs(%v, %d) = %v, %v; want only %v",
				tc.code, tc.vendorID, avps, err, a)
		}
	}

	path := []interface{}{avp.AuthenticationInfo, avp.UTRANVector, avp.ConfidentialityKey}
	avps, err := m.FindAVPsWithPath(path, 10415)
	if err != nil || len(avps) != 1 || avps[0].Data != ck {
		t.Errorf("FindAVPsWithPath(%v, 10415) = %v, %v; want the Confidentiality-Key", path, avps, err)
	}
	// OC-OLR holds an OC-Validity-Duration (625), not a Confidentiality-Key.
	// A uint32 code not defined for the Vendor-Id resolves to the AVP of
	// another vendor (here OC-OLR, Vendor-Id 0), which lets a path cross
	// vendors.
	path = []interface{}{uint32(avp.OCOLR), uint32(avp.ConfidentialityKey)}
	if avps, err := m.FindAVPsWithPath(path, 10415); err != nil || len(avps) != 0 {
		t.Errorf("FindAVPsWithPath(%v, 10415) = %v, %v; want no AVPs", path, avps, err)
	}
}
