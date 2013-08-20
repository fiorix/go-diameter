// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP parser.  Part of go-diameter.

package base

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"

	"github.com/fiorix/go-diameter/dict"
)

type rfcHdr1 struct {
	Code   uint32
	Flags  uint8
	Length [3]uint8
}

// ReadAVP reads an AVP and returns the number of extra bytes read and parsed
// AVP, or an error.
//
// Extra bytes are read when the content of the AVP is OctetString and
// has padding. Total bytes read is each avp.Length + extra.
//
// A pointer to the parent Message is required.
func ReadAVP(r io.Reader, m *Message) (uint32, *AVP, error) {
	if m == nil {
		panic("Can't read AVP without parent Message")
	}
	var (
		err error
		raw rfcHdr1
	)
	if err = binary.Read(r, binary.BigEndian, &raw); err != nil {
		return 0, nil, err
	}
	avp := &AVP{
		Code:    raw.Code,
		Flags:   raw.Flags,
		Length:  uint24To32(raw.Length),
		Message: m,
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
	// pad it, or even recursively load grouped AVPs from it's data.
	var davp *dict.AVP
	if davp, err = m.Dict.FindAVP(m.Header.ApplicationId, avp.Code); err != nil {
		return 0, nil, fmt.Errorf("Unknown AVP code %d: missing dict?", avp.Code)
	}
	// Read grouped (embedded) AVPs.
	//
	// Grouped AVPs are the only reason why this function returns
	// "extra" bytes, otherwise callers would have to walk through
	// these grouped AVPs and sum their padding + length to figure
	// out the total of bytes read. The value of "extra" represents
	// padded bytes read but not stored anywhere, but count when
	// check summing the entire message length.
	//
	// TODO: Double check the handling of dynamically grouped AVPs
	//       in case of 260, 279 or 284. Should work as is.
	if davp.Data.Type == "Grouped" {
		var (
			extra uint32
			gravp *AVP
			err   error
		)
		for dlen > 0 {
			if extra, gravp, err = ReadAVP(r, m); err != nil {
				return 0, nil, err
			}
			if avp.Data == nil {
				avp.Data = new(Grouped)
			}
			avp.Data.Put(gravp)
			dlen -= gravp.Length + extra
		}
		return extra, avp, nil
	}
	// Read binary data of regualr (non-grouped) AVPs.
	b := make([]byte, dlen)
	if err = binary.Read(r, binary.BigEndian, b); err != nil {
		return 0, nil, err
	}
	var pad bool // Indicates the data might have been padded.
	switch davp.Data.Type {
	case "OctetString":
		pad = true
		avp.Data = new(OctetString)
	case "Integer32":
		avp.Data = new(Integer32)
	case "Integer64":
		avp.Data = new(Integer64)
	case "Unsigned32":
		avp.Data = new(Unsigned32)
	case "Unsigned64":
		avp.Data = new(Unsigned64)
	case "Float32":
		avp.Data = new(Float32)
	case "Float64":
		avp.Data = new(Float64)
	case "Address":
		pad = true
		avp.Data = new(Address)
	case "Time":
		pad = true
		avp.Data = new(Time)
	case "UTF8String":
		pad = true
		avp.Data = new(OctetString)
	case "DiameterIdentity":
		pad = true
		avp.Data = new(DiameterIdentity)
	case "DiameterURI":
		avp.Data = new(DiameterURI)
	case "Enumerated":
		avp.Data = new(Enumerated)
	case "IPFilterRule":
		pad = true
		avp.Data = new(IPFilterRule)
	}
	avp.Data.Put(b)
	// Check if there's extra data to read due to padding of OctetString.
	//
	// http://tools.ietf.org/html/rfc6733#section-4
	//
	// Each AVP of type OctetString MUST be padded to align on a 32-bit
	// boundary, while other AVP types align naturally.  A number of zero-
	// valued bytes are added to the end of the AVP Data field till a word
	// boundary is reached.  The length of the padding is not reflected in
	// the AVP Length field.
	//
	// This also applies to subtypes of OctetString such as Address.
	var extrabytes uint32
	if pad {
		// Read and discard pad bytes.
		if n := pad4(dlen) - dlen; n > 0 {
			extrabytes += n
			b := make([]byte, n)
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
	if davp, err := avp.Message.Dict.FindAVP(
		avp.Message.Header.ApplicationId,
		avp.Code,
	); err != nil {
		name = "Unknown"
	} else {
		name = davp.Name
	}
	return fmt.Sprintf("%s AVP{Code=%d,Flags=%#x,Length=%d,VendorId=%#x,%s}",
		name,
		avp.Code,
		avp.Flags,
		avp.Length,
		avp.VendorId,
		avp.Data,
	)
}
