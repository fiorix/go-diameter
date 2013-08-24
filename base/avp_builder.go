// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP parser.  Part of go-diameter.

package base

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/fiorix/go-diameter/dict"
)

type rfcHdr2 struct {
	Code     uint32
	Flags    uint8
	Length   [3]uint8
	VendorId uint32
}

// NewAVP allocates and returns a new AVP.
func NewAVP(d *dict.Parser, code uint32, flags uint8, vendor uint32, body Codec) *AVP {
	if d == nil {
		panic("AVP requires a valid dictionary, not nil")
	}
	avp := &AVP{
		Code:     code,
		Flags:    flags,
		VendorId: vendor,
		body:     body,
		dict:     d,
	}
	if flags&0x80 > 0 {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr2{}))
	} else {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr1{}))
	}
	avp.Length += body.Length()
	return avp
}

// Bytes returns an AVP in binary form so it can be attached to a Message
// before sent to a connection.
func (avp *AVP) Bytes() []byte {
	b := bytes.NewBuffer(nil)
	if avp.Flags&0x80 > 0 {
		hdr := rfcHdr2{
			Code:     avp.Code,
			Flags:    avp.Flags,
			Length:   uint32To24(avp.Length),
			VendorId: avp.VendorId,
		}
		binary.Write(b, binary.BigEndian, hdr)
	} else {
		hdr := rfcHdr1{
			Code:   avp.Code,
			Flags:  avp.Flags,
			Length: uint32To24(avp.Length),
		}
		binary.Write(b, binary.BigEndian, hdr)
	}
	binary.Write(b, binary.BigEndian, avp.body.Bytes())
	return b.Bytes()
}
