package diam

import (
	"bytes"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// TestUnmarshalVendorId is a regression test for #169.
//
// newIndex used to key AVPs by Code only, so AVPs with the same code but
// different VendorIds ended up in the same bucket. scanStruct then picked
// the first match regardless of VendorId, corrupting the destination
// struct. After the fix, {Code, VendorID} is the index key and Unmarshal
// distinguishes the two AVPs correctly. This test drives the public
// Unmarshal API so it catches any regression in scanStruct's use of the
// index, not just in newIndex.
func TestUnmarshalVendorId(t *testing.T) {
	// Two AVPs sharing code 777 but differing by VendorID.
	const dictXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="4">
    <avp name="Example-Base" code="777" must="M" may="-" must-not="V" may-encrypt="Y">
      <data type="OctetString" />
    </avp>
    <avp name="Example-Vendor" code="777" must="M,V" may="-" must-not="-" may-encrypt="Y" vendor-id="10415">
      <data type="Unsigned32" />
    </avp>
  </application>
</diameter>`
	if err := dict.Default.Load(bytes.NewReader([]byte(dictXML))); err != nil {
		t.Fatalf("load dict: %v", err)
	}

	m := NewRequest(272, 4, dict.Default)
	m.AddAVP(NewAVP(777, avp.Mbit, 0, datatype.OctetString("base-data")))
	m.AddAVP(NewAVP(777, avp.Mbit|avp.Vbit, 10415, datatype.Unsigned32(12345)))

	var dst struct {
		Base   datatype.OctetString `avp:"Example-Base"`
		Vendor datatype.Unsigned32  `avp:"Example-Vendor"`
	}
	if err := m.Unmarshal(&dst); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got, want := string(dst.Base), "base-data"; got != want {
		t.Errorf("Base: got %q, want %q (wrong AVP routed to Base field)", got, want)
	}
	if got, want := uint32(dst.Vendor), uint32(12345); got != want {
		t.Errorf("Vendor: got %d, want %d (wrong AVP routed to Vendor field)", got, want)
	}
}
