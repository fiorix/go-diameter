// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"testing"

	"github.com/fiorix/go-diameter/diam/dict"
)

func TestGrouped(t *testing.T) {
	m := NewMessage(267, 0, 0, 0, 0, dict.Default)
	if _, err := m.NewAVP("Origin-Host", 0x40, 0, "z"); err != nil {
		t.Error(err)
	}
	gr := NewGroup(m)
	if _, err := gr.NewAVP("Vendor-Id", 0x40, 0, 1); err != nil {
		t.Error(err)
	}
	if _, err := gr.NewAVP("Vendor-Id", 0x40, 0, 2); err != nil {
		t.Error(err)
	}
	if _, err := m.NewAVP("Vendor-Specific-Application-Id", 0x40, 0, gr); err != nil {
		t.Error(err)
	}
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
