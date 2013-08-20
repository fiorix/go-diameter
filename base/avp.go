// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter AVP.  Part of go-diameter.
// http://tools.ietf.org/html/rfc6733#section-4

package base

// AVP represents an AVP header and data.
type AVP struct {
	Code     uint32
	Flags    uint8
	Length   uint32
	VendorId uint32
	Data     Codec
	Message  *Message // Link to parent Message
}

// Data is an interface for AVP Data types.
// These are the supported types and their equivalent native types:
//
// Diameter             Go
// --------             --
// OctetString          string
// DiameterIdentity     string
// Integer32            int32
// Integer64            int64
// Unsigned32           uint32
// Unsigned64           uint64
// Enumerated           []EnumItem
// Grouped              []*AVP
type Data interface{}

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
