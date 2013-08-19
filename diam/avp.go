// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
	"unsafe"
)

type AVP struct {
	Code     uint32
	Flags    uint8
	Length   uint32
	VendorId uint32
	Data     interface{}
	RawData  []byte
	Padding  int
}

type Grouped []*AVP

type rfcHdr1 struct {
	Code   uint32
	Flags  uint8
	Length [3]uint8
}

type rfcHdr2 struct {
	Code     uint32
	Flags    uint8
	Length   [3]uint8
	VendorId uint32
}

// ReadAVP reads an AVP and returns the number of extra bytes read and parsed
// AVP, or an error.
// Extra bytes are read when the content of the AVP is OctetString and
// needs padding. Total bytes read is avp.Length + extra.
func ReadAVP(r io.Reader, dict *Dict) (uint32, *AVP, error) {
	var (
		err error
		raw rfcHdr1
	)
	if err = binary.Read(r, binary.BigEndian, &raw); err != nil {
		return 0, nil, err
	}
	avp := &AVP{
		Code:   raw.Code,
		Flags:  raw.Flags,
		Length: uint24To32(raw.Length),
	}
	dlen := avp.Length - uint32(unsafe.Sizeof(raw))
	// Read VendorId when necessary.
	if raw.Flags&0x20 > 0 {
		if err = binary.Read(r, binary.BigEndian, &avp.VendorId); err != nil {
			return 0, nil, err
		}
		dlen -= uint32(unsafe.Sizeof(avp.VendorId))
	}
	// Find this AVP in a pre-loaded dict so we know how to parse it,
	// pad it, or load grouped AVPs from it's data.
	var davp *DictAVP
	if davp, err = dict.AVP(avp.Code); err != nil {
		return 0, nil, fmt.Errorf("Unknown AVP code %d: missing dict?", avp.Code)
	}
	// Read grouped (embedded) AVPs.
	// Grouped AVPs are the only reason why this function returns
	// "extra" bytes, otherwise callers would have to walk through
	// these grouped AVPs and sum their padding + length to figure
	// out the total of bytes read.
	// TODO: Handle dynamically grouped AVPs in case of 260, 279 or 284.
	if davp.Data.Type == "Grouped" {
		if eb, gavp, err := ReadAVP(r, dict); err != nil {
			return 0, nil, err
		} else {
			if avp.Data == nil {
				avp.Data = Grouped{}
			}
			avp.Data = append(avp.Data.(Grouped), gavp)
			return eb, avp, nil
		}
	}
	var needPadding bool
	// Read binary AVP data.
	avp.Data = make([]byte, dlen)
	if err = binary.Read(r, binary.BigEndian, avp.Data); err != nil {
		return 0, nil, err
	}
	switch davp.Data.Type {
	case "Address":
		needPadding = true
		// TODO: Double check this.
		if len(avp.Data.([]byte)) == 6 {
			b := avp.Data.([]byte)
			if b[1] == 1 { // Family address 1 == IPv4
				avp.Data = net.IPv4(b[2], b[3], b[4], b[5])
			}
			break
		}
		// Fallback to string instead of IPv4.
		avp.Data = string(avp.Data.([]byte))
	case "Time":
		fallthrough
	case "DiameterIdentity":
		fallthrough
	case "DiameterURI":
		fallthrough
	case "OctetString":
		fallthrough
	case "UTF8String":
		needPadding = true
		avp.Data = string(avp.Data.([]byte))
	case "Enumerated":
		fallthrough
	case "Unsigned32":
		if dlen != 4 {
			return 0, nil, fmt.Errorf("Expecting Unsigned32, got %d instead", dlen*8)
		}
		var n uint32
		netBytesToN(avp.Data.([]byte), &n)
		avp.Data = n
	case "Unsigned64":
		if dlen != 8 {
			return 0, nil, fmt.Errorf("Expecting Unsigned64, got %d instead", dlen*8)
		}
		var n uint64
		netBytesToN(avp.Data.([]byte), &n)
		avp.Data = n
	}
	// Check if there's extra data to read due to padding of OctetString.
	//
	// http://tools.ietf.org/html/rfc3588#section-4
	//
	// Each AVP of type OctetString MUST be padded to align on a 32-bit
	// boundary, while other AVP types align naturally.  A number of zero-
	// valued bytes are added to the end of the AVP Data field till a word
	// boundary is reached.  The length of the padding is not reflected in
	// the AVP Length field.
	//
	// This also applies to subtypes of OctetString such as Address.
	var extrabytes uint32
	if needPadding {
		// Read and discard pad bytes.
		if avp.Padding = int(pad4(dlen) - dlen); avp.Padding > 0 {
			extrabytes += uint32(avp.Padding)
			b := make([]byte, avp.Padding)
			if _, err = io.ReadFull(r, b); err != nil {
				return 0, nil, err
			}
		}
	}
	return extrabytes, avp, nil
}

// String returns the AVP in human readable format.
func (avp *AVP) String() string {
	// TODO: Lookup the vendor id from AVP in the dictionary.
	var name string
	if davp, err := BaseDict.AVP(avp.Code); err != nil {
		name = "Unknown"
	} else {
		name = davp.Name
	}
	v := fmt.Sprintf("%s AVP{Code=%d,Flags=%#x,Length=%d,VendorId=%#x,Padding=%d,",
		name, avp.Code, avp.Flags, avp.Length, avp.VendorId, avp.Padding)
	switch avp.Data.(type) {
	case string:
		v += fmt.Sprintf("string('%s')", avp.Data.(string))
	case DictEnumItem:
		v += fmt.Sprintf("enum(%s)", avp.Data.(DictEnumItem))
	case time.Time:
		v += fmt.Sprintf("time(%s)", avp.Data.(time.Time))
	case uint32:
		v += fmt.Sprintf("uint32(%d)", avp.Data.(uint32))
	case uint64:
		v += fmt.Sprintf("uint64(%d)", avp.Data.(uint64))
	case net.IP:
		v += fmt.Sprintf("net.IP(%s)", avp.Data.(net.IP))
	case Grouped:
		v += fmt.Sprintf("Grouped(%s)", avp.Data.(Grouped))
	default:
		v += fmt.Sprintf("Unknown(%s)", avp.Data)
	}
	return v + "}"
}

// NewAVP allocates and returns a new AVP. Used for building messages.
func NewAVP(code uint32, flags uint8, vendor uint32, data interface{}) *AVP {
	avp := &AVP{
		Code:     code,
		Flags:    flags,
		VendorId: vendor,
		Data:     data,
	}
	if flags&0x20 > 0 {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr2{}))
	} else {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr1{}))
	}
	switch data.(type) {
	case string:
		avpString(avp, []byte(data.(string)))
	case net.IP:
		ip := data.(net.IP).To4()
		// TODO: Fix this static family 1 + IPv4 address.
		avpString(avp, []byte{0, 1, ip[0], ip[1], ip[2], ip[3]})
	case uint32:
		avp.Length += 4
		avp.RawData, _ = nToNetBytes(data.(uint32))
	case uint64:
		avp.Length += 8
		avp.RawData, _ = nToNetBytes(data.(uint64))
	}
	return avp
}

// avpString encodes the AVP data with padding if needed, and updates the
// Length header accordingly. It also sets the Padding attribute.
func avpString(avp *AVP, s []byte) {
	length := uint32(len(s))
	avp.Length += length
	if extra := pad4(length) - length; extra > 0 {
		avp.Padding = int(extra)
		avp.RawData = make([]byte, length+extra)
		copy(avp.RawData, s)
	} else {
		avp.RawData = s
	}
}

// Marshal returns an AVP in binary form so it can be attached to a Message
// before sent to a connection.
func (avp *AVP) Marshal() []byte {
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
	binary.Write(b, binary.BigEndian, avp.RawData)
	return b.Bytes()
}
