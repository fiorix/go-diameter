// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

// DataType is an interface to support Diameter AVP data types.
type DataType interface {
	Serialize() []byte
	Len() int
	Padding() int
	Type() TypeID
	String() string
}

// TypeID is the identifier of an AVP data type.
type TypeID int

// List of available AVP data types.
const (
	UnknownType TypeID = iota
	AddressType
	DiameterIdentityType
	DiameterURIType
	EnumeratedType
	Float32Type
	Float64Type
	GroupedType
	IPFilterRuleType
	IPv4Type
	Integer32Type
	Integer64Type
	OctetStringType
	TimeType
	UTF8StringType
	Unsigned32Type
	Unsigned64Type
)
