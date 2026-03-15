// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// DecoderFunc is an adapter to decode a byte array to an AVP data type.
type DecoderFunc func([]byte) (Type, error)

// Decoder is a map of AVP data types indexed by TypeID.
// Kept for backward compatibility.
var Decoder = map[TypeID]DecoderFunc{
	UnknownType:          DecodeUnknown,
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

// maxTypeID is one past the highest TypeID, used to size the decoder array.
const maxTypeID = IPv6Type + 1

// decoderArray is an array-indexed decoder table for fast O(1) dispatch
// without map hashing overhead.
var decoderArray [maxTypeID]DecoderFunc

func init() {
	decoderArray[UnknownType] = DecodeUnknown
	decoderArray[AddressType] = DecodeAddress
	decoderArray[DiameterIdentityType] = DecodeDiameterIdentity
	decoderArray[DiameterURIType] = DecodeDiameterURI
	decoderArray[EnumeratedType] = DecodeEnumerated
	decoderArray[Float32Type] = DecodeFloat32
	decoderArray[Float64Type] = DecodeFloat64
	decoderArray[GroupedType] = DecodeGrouped
	decoderArray[IPFilterRuleType] = DecodeIPFilterRule
	decoderArray[IPv4Type] = DecodeIPv4
	decoderArray[Integer32Type] = DecodeInteger32
	decoderArray[Integer64Type] = DecodeInteger64
	decoderArray[OctetStringType] = DecodeOctetString
	decoderArray[QoSFilterRuleType] = DecodeQoSFilterRule
	decoderArray[TimeType] = DecodeTime
	decoderArray[UTF8StringType] = DecodeUTF8String
	decoderArray[Unsigned32Type] = DecodeUnsigned32
	decoderArray[Unsigned64Type] = DecodeUnsigned64
	decoderArray[IPv6Type] = DecodeIPv6
}

// Decode decodes a specific AVP data type from byte array to a DataType.
func Decode(t TypeID, b []byte) (Type, error) {
	if int(t) >= 0 && int(t) < len(decoderArray) {
		if f := decoderArray[t]; f != nil {
			return f(b)
		}
	}
	return nil, fmt.Errorf("Unknown data type: %d", t)
}
