// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter AVP.  Part of go-diameter.
// http://tools.ietf.org/html/rfc6733#section-4

package diam

import (
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/dict"
)

// AVP represents an AVP header and data.
type AVP struct {
	Code     uint32
	Flags    uint8
	Length   uint32
	VendorId uint32
	body     Codec        // AVP data
	dict     *dict.Parser // Dictionary (for PrettyPrint)
}

// Data returns internal AVP body data.  It's a short for AVP.Body().Data().
func (avp *AVP) Data() datatype.Generic {
	return avp.body.Data()
}

// Body returns the internal AVP body.
func (avp *AVP) Body() Codec {
	return avp.body
}

// Codec provides an interface for converting Data from network bytes to
// native and vice-versa.
type Codec interface {
	// Write binary data from the network to this AVP Data.
	Put(Data)

	// Encode this AVP Data into binary data.
	Bytes() []byte

	// Returns its internal Data.
	Data() Data

	// Length without padding. Might be diffent from len(Bytes()).
	Length() uint32

	// String represents the AVP data in human readable format.
	String() string
}
