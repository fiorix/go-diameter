// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diameter

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
}

func (avp *AVP) String() string {
	var name string
	if davp, err := BaseDict.AVP(avp.Code); err != nil {
		name = "Unknown"
	} else {
		name = davp.Name
	}
	v := fmt.Sprintf("%s AVP{Code=%d,Flags=%#x,Length=%d,VendorId=%#x,",
		name, avp.Code, avp.Flags, avp.Length, avp.VendorId)
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
		v += fmt.Sprintf("ip(%s)", avp.Data.(net.IP))
	case GroupedAVP:
		v += fmt.Sprintf("grouped(%s)", avp.Data.(GroupedAVP))
	default:
		v += fmt.Sprintf("unknown(%s)", avp.Data)
	}
	return v + "}"
}

type GroupedAVP []*AVP

type rawAVP struct {
	Code   uint32
	Flags  uint8
	Length [3]uint8
}

// ReadAVP reads an AVP and returns the number of extra bytes read and parsed
// AVP, or an error.
// Extra bytes are read when the content of the AVP is OctetString and
// needs padding. Total bytes read is avp.Length + extra.
func ReadAVP(r io.Reader, dict *Dict) (uint32, *AVP, error) {
	var (
		err error
		raw rawAVP
	)
	if err = binary.Read(r, binary.BigEndian, &raw); err != nil {
		return 0, nil, err
	}
	avp := &AVP{
		Code:   raw.Code,
		Flags:  raw.Flags,
		Length: uint24to32(raw.Length),
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
	extrabytes := uint32(0)
	// Read grouped (embedded) AVPs.
	// TODO: Handle dynamically grouped AVPs in case of 260, 279 or 284.
	if davp.Data.Type == "Grouped" {
		if eb, gavp, err := ReadAVP(r, dict); err != nil {
			return 0, nil, err
		} else {
			if avp.Data == nil {
				avp.Data = GroupedAVP{}
			}
			avp.Data = append(avp.Data.(GroupedAVP), gavp)
			extrabytes += eb
		}
	} else {
		// Read binary AVP data.
		avp.Data = make([]byte, dlen)
		if err = binary.Read(r, binary.BigEndian, avp.Data); err != nil {
			return 0, nil, err
		}
	}
	// If grouped AVPs were read, there's no need to proceed.
	if _, ok := avp.Data.(GroupedAVP); ok {
		return extrabytes, avp, nil
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
	var (
		needPadding bool
		n           int64
	)
	switch davp.Data.Type {
	case "Address":
		needPadding = true
		// TODO: Double check this.
		if len(avp.Data.([]byte)) == 6 {
			b := avp.Data.([]byte)
			avp.Data = net.IPv4(b[2], b[3], b[4], b[5])
			break
		}
		// Fallback to string instead of IPv4.
		fallthrough
	case "DiameterIdentity":
		fallthrough
	case "DiameterURI":
		fallthrough
	case "OctetString":
		needPadding = true
		avp.Data = string(avp.Data.([]byte))
	case "UTF8String":
		needPadding = true
		avp.Data = string(avp.Data.([]byte))
	case "Enumerated":
		n, err = binary.ReadVarint(bytes.NewBuffer(avp.Data.([]byte)))
		if err != nil {
			return 0, nil, err
		}
		if avp.Data, err = dict.Enum(avp.Code, uint8(n)); err != nil {
			return 0, nil, err
		}
	case "Time":
		n, err = binary.ReadVarint(bytes.NewBuffer(avp.Data.([]byte)))
		if err != nil {
			return 0, nil, err
		}
		avp.Data = time.Unix(n, 0)
	case "Unsigned32":
		n, err = binary.ReadVarint(bytes.NewBuffer(avp.Data.([]byte)))
		if err != nil {
			return 0, nil, err
		}
		avp.Data = uint32(n)
	case "Unsigned64":
		n, err = binary.ReadVarint(bytes.NewBuffer(avp.Data.([]byte)))
		if err != nil {
			return 0, nil, err
		}
		avp.Data = uint64(n)
	}
	if needPadding {
		// Read and discard pad bytes.
		eb := pad4(dlen) - dlen
		if eb > 0 {
			extrabytes += eb
			b := make([]byte, eb)
			if _, err = io.ReadFull(r, b); err != nil {
				return 0, nil, err
			}
		}
	}
	return extrabytes, avp, nil
}
