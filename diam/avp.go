// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter AVP.  http://tools.ietf.org/html/rfc6733#section-4

package diam

import (
	"fmt"

	"github.com/fiorix/go-diameter/diam/avpdata"
	"github.com/fiorix/go-diameter/diam/dict"
)

// AVP represents an AVP header and data.
type AVP struct {
	Code     uint32
	Flags    uint8
	Length   uint32
	VendorId uint32
	body     Codec        // AVP body
	dict     *dict.Parser // Dictionary associated with this AVP
}

// Data returns internal AVP body data.  It's a short for AVP.Body().Data().
func (avp *AVP) Data() avpdata.Generic {
	return avp.body.Data()
}

// Body returns the internal AVP body.
func (avp *AVP) Body() Codec {
	return avp.body
}

// String returns the AVP in human readable format.
//
// The AVP name is "guessed" by scanning the list of available AVPs in the
// dictionary that was used to build this AVP. It might return the wrong
// AVP name if the same code is used by different dictionaries in different
// applications, with a different name - yet, very unlikely.
func (avp *AVP) String() string {
	// TODO: Lookup the vendor id from AVP in the dictionary.
	var name string
	if avp.dict != nil {
		if davp, err := avp.dict.ScanAVP(avp.Code); davp != nil && err == nil {
			name = davp.Name
		}
	}
	if name == "" {
		name = "Unknown"
	}
	return fmt.Sprintf("%s AVP{Code=%d,Flags=%#x,Length=%d,VendorId=%#x,%s}",
		name,
		avp.Code,
		avp.Flags,
		avp.Length,
		avp.VendorId,
		avp.body,
	)
}
