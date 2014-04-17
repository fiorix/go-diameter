// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/diam/datatypes"
	"github.com/fiorix/go-diameter/diam/dict"
)

// Diameter AVP.
type AVP struct {
	Code     uint32             // Code of this AVP
	Flags    uint8              // Flags of this AVP
	Length   int                // Length of this AVP's payload
	VendorId uint32             // VendorId of this AVP
	Data     datatypes.DataType // Data of this AVP (payload)
}

// NewAVP creates and initializes a new AVP.
func NewAVP(code uint32, flags uint8, vendor uint32, data datatypes.DataType) *AVP {
	a := &AVP{
		Code:     code,
		Flags:    flags,
		VendorId: vendor,
		Data:     data,
	}
	a.Length = a.headerLen()
	return a
}

func decodeAVP(data []byte, application uint32, dictionary *dict.Parser) (*AVP, error) {
	avp := new(AVP)
	if err := avp.DecodeFromBytes(data, application, dictionary); err != nil {
		return nil, err
	}
	return avp, nil
}

// DecodeFromBytes decodes the bytes of a Diameter AVP.
// It requires a parent Header to be able to decode the AVP data by
// consulting the ApplicationId and Dictionary of the Header.
func (a *AVP) DecodeFromBytes(data []byte, application uint32, dictionary *dict.Parser) error {
	dl := len(data)
	if len(data) < 8 {
		return fmt.Errorf("Not enough data to decode AVP header: %d bytes", dl)
	}
	a.Code = binary.BigEndian.Uint32(data[0:4])
	// Find this code in the dictionary.
	dictAVP, err := dictionary.FindAVP(application, a.Code)
	if err != nil {
		return err
	}
	a.Flags = data[4]
	a.Length = int(uint24to32(data[5:8]))
	if dl < int(a.Length) {
		return fmt.Errorf("Not enough data to decode AVP: %d != %d",
			dl, a.Length)
	}
	data = data[:a.Length] // this cuts padded bytes off
	var hdrLength int
	var payload []byte
	// Read VendorId when required.
	if a.Flags&0x80 > 0 {
		a.VendorId = binary.BigEndian.Uint32(data[8:12])
		payload = data[12:]
		hdrLength = 12
	} else {
		payload = data[8:]
		hdrLength = 8
	}
	bodyLen := a.Length - hdrLength
	if n := len(payload); n < bodyLen {
		return fmt.Errorf(
			"Not enough data to decode AVP: %d != %d",
			hdrLength, n,
		)
	}
	a.Data, err = datatypes.Decode(dictAVP.Data.Type, payload)
	if err != nil {
		return err
	}
	// Handle grouped AVPs.
	if a.Data.Type() == datatypes.GroupedType {
		a.Data, err = DecodeGrouped(a.Data.(datatypes.Grouped), application, dictionary)
		if err != nil {
			return err
		}
	}
	return nil
}

// Serialize returns the byte sequence that represents this AVP.
// It requires at least the Code, Flags and Data fields set.
func (a *AVP) Serialize() ([]byte, error) {
	if a.Data == nil {
		return nil, errors.New("Failed to serialize AVP: Data is nil")
	}
	payload := a.Data.Serialize()
	payloadLen := len(payload)
	var b []byte
	if a.VendorId > 0 {
		b = make([]byte, 12+payloadLen+a.Data.Padding())
		copy(b[5:8], uint32to24(uint32(12+payloadLen)))
		binary.BigEndian.PutUint32(b[8:12], a.VendorId)
		copy(b[12:], payload)
	} else {
		b = make([]byte, 8+payloadLen+a.Data.Padding())
		copy(b[5:8], uint32to24(uint32(8+payloadLen)))
		copy(b[8:], payload)
	}
	binary.BigEndian.PutUint32(b[0:4], a.Code)
	b[4] = a.Flags
	return b, nil
}

func (a *AVP) SerializeTo(b []byte) error {
	if a.Data == nil {
		return errors.New("Failed to serialize AVP: Data is nil")
	}
	payload := a.Data.Serialize()
	payloadLen := len(payload)
	if a.VendorId > 0 {
		copy(b[5:8], uint32to24(uint32(12+payloadLen)))
		binary.BigEndian.PutUint32(b[8:12], a.VendorId)
		copy(b[12:], payload)
	} else {
		copy(b[5:8], uint32to24(uint32(8+payloadLen)))
		copy(b[8:], payload)
	}
	binary.BigEndian.PutUint32(b[0:4], a.Code)
	b[4] = a.Flags
	return nil
}

func (a *AVP) Len() int {
	return a.headerLen() + a.Data.Padding()
}

func (a *AVP) headerLen() int {
	if a.Flags&0x80 > 0 {
		return 12 + a.Data.Len()
	} else {
		return 8 + a.Data.Len()
	}
}

func (a *AVP) String() string {
	return fmt.Sprintf("{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:%s}",
		a.Code,
		a.Flags,
		a.Len(),
		a.VendorId,
		a.Data,
	)
}
