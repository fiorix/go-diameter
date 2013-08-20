// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP parser.  Part of go-diameter.

package base

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type rfcHdr2 struct {
	Code     uint32
	Flags    uint8
	Length   [3]uint8
	VendorId uint32
}

// NewAVP allocates and returns a new AVP. Used for building messages.
//
// The parent Message is required because it contains a link to the
// dictionary associated with the message, used for parsing AVP data.
// If nil, it will be automatically associated with the message when
// Message.Add(avp) is called.
func NewAVP(code uint32, flags uint8, vendor uint32, data Codec, m *Message) *AVP {
	avp := &AVP{
		Code:     code,
		Flags:    flags,
		VendorId: vendor,
		Data:     data,
		Message:  m,
	}
	if flags&0x20 > 0 {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr2{}))
	} else {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr1{}))
	}
	avp.Length += data.Length()
	return avp
}

// Bytes returns an AVP in binary form so it can be attached to a Message
// before sent to a connection.
func (avp *AVP) Bytes() []byte {
	b := bytes.NewBuffer(nil)
	if avp.Flags&0x20 > 0 {
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
	binary.Write(b, binary.BigEndian, avp.Data.Bytes())
	return b.Bytes()
}
