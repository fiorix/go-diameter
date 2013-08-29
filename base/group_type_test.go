// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package base

import "testing"

func TestGrouped(t *testing.T) {
	m := NewMessage(267, 0, 0, 0, 0, Dict)
	m.NewAVP("Origin-Host", 0x40, 0,
		&DiameterIdentity{OctetString{Value: "z"}})
	gr := NewGroup(m)
	gr.NewAVP("Vendor-Id", 0x40, 0, &Unsigned32{Value: 1})
	gr.NewAVP("Vendor-Id", 0x40, 0, &Unsigned32{Value: 2})
	m.NewAVP("Vendor-Specific-Application-Id", 0x40, 0, gr)
	if len(m.AVP) != 2 {
		t.Error("Missing newly created AVPs")
		return
	}
	// No drama, pls.
	if m.AVP[1].Data().([]*AVP)[0].Data() != uint32(1) {
		t.Error("Missing grouped AVP")
		return
	}
}
