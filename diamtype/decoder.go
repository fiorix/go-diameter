// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

var Decoder = map[DataTypeId]DecoderFunc{
	AddressType:          DecodeAddress,
	DiameterIdentityType: DecodeDiameterIdentity,
	DiameterURIType:      DecodeDiameterURI,
	EnumeratedType:       DecodeEnumerated,
	Float32Type:          DecodeFloat32,
	Float64Type:          DecodeFloat64,
	GroupedType:          DecodeGrouped,
	IPFilterRuleType:     DecodeIPFilterRule,
	IPv4Type:             DecodeIPv4,
	Integer32Type:        DecodeInteger32,
	Integer64Type:        DecodeInteger64,
	OctetStringType:      DecodeOctetString,
	TimeType:             DecodeTime,
	UTF8StringType:       DecodeUTF8String,
	Unsigned32Type:       DecodeUnsigned32,
	Unsigned64Type:       DecodeUnsigned64,
}

type DecoderFunc func([]byte) (DataType, error)

func Decode(datatype DataTypeId, b []byte) (DataType, error) {
	if f, exists := Decoder[datatype]; !exists {
		return nil, fmt.Errorf("Unknown data type: %d", datatype)
	} else {
		return f(b)
	}
}
