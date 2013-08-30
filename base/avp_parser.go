// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP parser.  Part of go-diameter.

package diam

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

type rfcHdr1 struct {
	Code   uint32
	Flags  uint8
	Length [3]uint8
}

var ErrInvalidAvpHdr = errors.New("Invalid AVP header size: probably a bad dict")

// ReadAVP reads an AVP and returns the number of extra bytes read and parsed
// AVP, or an error.
//
// Extra bytes are read when the content of the AVP is OctetString and
// has padding. Total bytes read is each avp.Length + extra.
//
// A pointer to the parent Message is required.
func ReadAVP(m *Message, r io.Reader) (uint32, *AVP, error) {
	if m == nil {
		panic("Can't read AVP without parent Message")
	}
	var raw rfcHdr1
	if err := binary.Read(r, binary.BigEndian, &raw); err != nil {
		return 0, nil, err
	}
	avp := &AVP{
		Code:   raw.Code,
		Flags:  raw.Flags,
		Length: uint24To32(raw.Length),
		dict:   m.Dict,
	}
	dlen := avp.Length - uint32(unsafe.Sizeof(raw))
	if dlen >= m.Header.MessageLength() {
		return 0, nil, ErrInvalidAvpHdr
	}
	// Read VendorId when necessary.
	if raw.Flags&0x80 > 0 {
		if err := binary.Read(r, binary.BigEndian, &avp.VendorId); err != nil {
			return 0, nil, err
		}
		dlen -= uint32(unsafe.Sizeof(avp.VendorId))
	}
	// Find this AVP in a pre-loaded dict so we know how to parse it,
	// pad it, or even recursively load grouped AVPs from it's data.
	//
	// If a previously parsed AVP has the incorrect size due to a bad
	// dictionary, this call might fail because the current header
	// will be malformed.
	davp, err := m.Dict.FindAVP(m.Header.ApplicationId, avp.Code)
	if err != nil {
		return 0, nil, fmt.Errorf(
			"Unknown AVP code %d for appid %d: missing dict",
			avp.Code, m.Header.ApplicationId)
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
		for dlen != 0 {
			ex, gravp, err := ReadAVP(m, r)
			if err != nil {
				return 0, nil, err
			} else if avp.body == nil {
				avp.body = &Grouped{Message: m}
			}
			avp.body.Put(gravp)
			dlen = dlen - (gravp.Length + ex)
		}
		// Lesson learned: there's never extra bytes in Grouped AVPs.
		// Only OctetString.
		return 0, avp, nil
	}
	// Read binary data of regular (non-grouped) AVPs.
	b := make([]byte, dlen)
	if err = binary.Read(r, binary.BigEndian, b); err != nil {
		return 0, nil, err
	}
	var pad bool // Indicates the data might have been padded.
	switch davp.Data.Type {
	case "OctetString":
		pad = true
		avp.body = new(OctetString)
	case "Integer32":
		avp.body = new(Integer32)
	case "Integer64":
		avp.body = new(Integer64)
	case "Unsigned32":
		avp.body = new(Unsigned32)
	case "Unsigned64":
		avp.body = new(Unsigned64)
	case "Float32":
		avp.body = new(Float32)
	case "Float64":
		avp.body = new(Float64)
	case "Address":
		pad = true
		avp.body = new(Address)
	case "IPv4": // To support Framed-IP-Address and alike.
		avp.body = new(IPv4)
	case "Time":
		pad = true
		avp.body = new(Time)
	case "UTF8String":
		pad = true
		avp.body = new(UTF8String)
	case "DiameterIdentity":
		pad = true
		avp.body = new(DiameterIdentity)
	case "DiameterURI":
		avp.body = new(DiameterURI)
	case "Enumerated":
		avp.body = new(Enumerated)
	case "IPFilterRule":
		pad = true
		avp.body = new(IPFilterRule)
	default:
		return 0, nil, fmt.Errorf(
			"Unsupported avp.body type: %s", davp.Data.Type)
	}
	// Put binary data in this AVP.
	avp.body.Put(b)
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
	var n uint32 // extra bytes to read
	if pad {
		// Read and discard pad bytes.
		if n = pad4(dlen) - dlen; n > 0 {
			b := make([]byte, n)
			if _, err = io.ReadFull(r, b); err != nil {
				return 0, nil, err
			}
		}
	}
	return n, avp, nil
}

// String returns the AVP in human readable format.
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
