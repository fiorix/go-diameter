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

func newAVP(appid uint32, d *dict.Parser, code interface{}, flags uint8, vendor uint32, body Codec) (*AVP, error) {
	davp, err := d.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	avp := &AVP{
		Code:     davp.Code,
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
	return avp, nil
}

// NewAVP allocates and returns a new AVP.
// @code can be either the AVP code (int, uint32) or name (string).
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, body Codec) (*AVP, error) {
	avp, err := newAVP(
		m.Header.ApplicationId,
		m.Dict,
		code,
		flags,
		vendor,
		body,
	)
	if err != nil {
		return nil, err
	}
	m.AVP = append(m.AVP, avp)
	return avp, nil
}

// Set updates the internal AVP data (the body) with a new value.
func (avp *AVP) Set(body Codec) {
	if avp.body != nil {
		avp.Length -= avp.body.Length()
	}
	avp.body = body
	avp.Length += body.Length()
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
