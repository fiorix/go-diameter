// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"fmt"

	"github.com/fiorix/go-diameter/diam/util"
)

// OctetString Diameter Type.
type OctetString struct {
	Value   string
	Padding uint32 // Extra bytes to make the Value a multiple of 4 octets
}

// Data implements the Data interface.
func (os *OctetString) Data() Generic {
	return os.Value
}

// Put implements the Codec interface. It updates internal Value and Padding.
func (os *OctetString) Put(b []byte) {
	l := uint32(len(b))
	os.Padding = util.Pad4(l) - l
	os.Value = string(b)
}

// Bytes implement the Codec interface. Padding is always recalculated from
// the internal Value.
func (os *OctetString) Bytes() []byte {
	os.updatePadding() // Do this every time? Geez.
	l := uint32(len(os.Value))
	b := make([]byte, l+os.Padding)
	copy(b, os.Value)
	return b
}

// Length implements the Codec interface. Returns length without padding.
func (os *OctetString) Length() uint32 {
	return uint32(len(os.Value))
}

// update internal padding value.
func (os *OctetString) updatePadding() {
	if os.Padding == 0 {
		l := uint32(len(os.Value))
		os.Padding = util.Pad4(l) - l
	}
}

// String returns a human readable version of the AVP.
func (os *OctetString) String() string {
	os.updatePadding() // Update padding
	return fmt.Sprintf("OctetString{Value:'%s',Padding:%d}",
		os.Value, os.Padding)
}
