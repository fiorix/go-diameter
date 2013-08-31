// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import "github.com/fiorix/go-diameter/diam/avpdata"

// Codec provides an interface for converting data from network bytes to
// Go native and vice-versa.
//
// This interface is used to encode and decode AVP data such as OctetString,
// and is implemented by all data types in the avpdata sub package.
type Codec interface {
	// Decode binary data from the network and store as internal AVP data.
	Put([]byte)

	// Encode internal AVP data to binary and return it.
	Bytes() []byte

	// Returns internal AVP data.
	Data() avpdata.Generic

	// Length of the internal data without padding.
	// Might be diffent from len(Bytes()).
	Length() uint32

	// String returns the internal AVP data in human readable format.
	String() string
}
