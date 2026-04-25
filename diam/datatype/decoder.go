// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"errors"
	"fmt"
)

// DecoderFunc is an adapter to decode a byte array to an AVP data type.
type DecoderFunc func([]byte) (Type, error)

// Decoder is a map of AVP data types indexed by TypeID.
//
// Writing to this map directly only affects new TypeIDs (those outside the
// built-in range covered by decoderArray). To override a built-in decoder
// on the fast path, use RegisterDecoder instead; direct writes to this map
// for built-in TypeIDs are ignored by Decode.
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
	IPv6Type:             DecodeIPv6,
	Integer32Type:        DecodeInteger32,
	Integer64Type:        DecodeInteger64,
	OctetStringType:      DecodeOctetString,
	QoSFilterRuleType:    DecodeQoSFilterRule,
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
	decoderArray[IPv6Type] = DecodeIPv6
	decoderArray[Integer32Type] = DecodeInteger32
	decoderArray[Integer64Type] = DecodeInteger64
	decoderArray[OctetStringType] = DecodeOctetString
	decoderArray[QoSFilterRuleType] = DecodeQoSFilterRule
	decoderArray[TimeType] = DecodeTime
	decoderArray[UTF8StringType] = DecodeUTF8String
	decoderArray[Unsigned32Type] = DecodeUnsigned32
	decoderArray[Unsigned64Type] = DecodeUnsigned64
}

// Decode decodes a specific AVP data type from byte array to a DataType.
// The fast path dispatches through decoderArray; TypeIDs not covered there
// (including custom types registered in the exported Decoder map) fall back
// to the map lookup.
func Decode(t TypeID, b []byte) (Type, error) {
	if int(t) >= 0 && int(t) < len(decoderArray) {
		if f := decoderArray[t]; f != nil {
			return f(b)
		}
	}
	if f, ok := Decoder[t]; ok {
		return f(b)
	}
	return nil, fmt.Errorf("Unknown data type: %d", t)
}

// ErrCannotUnregisterBuiltin is returned by RegisterDecoder when the caller
// attempts to unregister a built-in TypeID by passing a nil DecoderFunc.
// The original built-in decoder is not retained anywhere reachable by the
// caller, so honoring the request would silently disable Decode for that
// type with no path to recover.
var ErrCannotUnregisterBuiltin = errors.New(
	"datatype: cannot unregister a built-in TypeID; pass a real DecoderFunc " +
		"or register a custom TypeID instead")

// RegisterDecoder installs f as the decoder for TypeID t.
//
// Dispatch targets:
//   - If t falls within the built-in range (0 <= t < maxTypeID), f is
//     written to decoderArray so Decode picks it up on the zero-allocation
//     fast path. This is how you override a built-in decoder without
//     giving up the array dispatch performance.
//   - Otherwise (custom TypeIDs beyond the built-in range), f is written
//     to the Decoder map, which Decode consults as a fallback.
//
// In both cases the Decoder map is also kept in sync, so iterating
// Decoder gives an accurate view of what is actually registered.
//
// Typical usage:
//
//	// Override a built-in for a specialized decode path:
//	if err := datatype.RegisterDecoder(datatype.Unsigned32Type, myFastUnsigned32); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Register a brand-new type:
//	const MyCustomType datatype.TypeID = 100
//	_ = datatype.RegisterDecoder(MyCustomType, myCustomDecoder)
//
// Passing f == nil unregisters a custom (non-built-in) TypeID by removing
// the entry from the Decoder map and returns nil. Passing nil for a
// built-in TypeID leaves the registration unchanged and returns
// ErrCannotUnregisterBuiltin: the caller has no reference to the original
// built-in decoder and would otherwise silently break the Decode fast path
// with no way to recover. Re-register with a real function, or fork the
// library if truly disabling a built-in type is required.
//
// Concurrency: RegisterDecoder mutates package-global state without
// locking and must not be called concurrently with Decode. Register
// all custom or overriding decoders from init() or during startup,
// before any goroutine begins decoding messages.
func RegisterDecoder(t TypeID, f DecoderFunc) error {
	inBuiltin := int(t) >= 0 && int(t) < len(decoderArray)
	if f == nil && inBuiltin {
		return fmt.Errorf("%w: TypeID %d", ErrCannotUnregisterBuiltin, t)
	}
	if inBuiltin {
		decoderArray[t] = f
		Decoder[t] = f
		return nil
	}
	if f == nil {
		delete(Decoder, t)
		return nil
	}
	Decoder[t] = f
	return nil
}
