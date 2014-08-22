// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package format

import "fmt"

var Decoder = map[FormatId]DecoderFunc{
	AddressFormat:          DecodeAddress,
	DiameterIdentityFormat: DecodeDiameterIdentity,
	DiameterURIFormat:      DecodeDiameterURI,
	EnumeratedFormat:       DecodeEnumerated,
	Float32Format:          DecodeFloat32,
	Float64Format:          DecodeFloat64,
	GroupedFormat:          DecodeGrouped,
	IPFilterRuleFormat:     DecodeIPFilterRule,
	IPv4Format:             DecodeIPv4,
	Integer32Format:        DecodeInteger32,
	Integer64Format:        DecodeInteger64,
	OctetStringFormat:      DecodeOctetString,
	TimeFormat:             DecodeTime,
	UTF8StringFormat:       DecodeUTF8String,
	Unsigned32Format:       DecodeUnsigned32,
	Unsigned64Format:       DecodeUnsigned64,
}

type DecoderFunc func([]byte) (Format, error)

func Decode(Format FormatId, b []byte) (Format, error) {
	if f, exists := Decoder[Format]; !exists {
		return nil, fmt.Errorf("Unknown data type: %d", Format)
	} else {
		return f(b)
	}
}
