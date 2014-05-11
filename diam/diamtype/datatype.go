// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

type DataType interface {
	Serialize() []byte
	Len() int
	Padding() int
	Type() DataTypeId
	String() string
}

type DataTypeId int

const (
	UnknownType DataTypeId = iota
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

var Available = map[string]DataTypeId{
	"Address":          AddressType,
	"DiameterIdentity": DiameterIdentityType,
	"DiameterURI":      DiameterURIType,
	"Enumerated":       EnumeratedType,
	"Float32":          Float32Type,
	"Float64":          Float64Type,
	"Grouped":          GroupedType,
	"IPFilterRule":     IPFilterRuleType,
	"IPv4":             IPv4Type,
	"Integer32":        Integer32Type,
	"Integer64":        Integer64Type,
	"OctetString":      OctetStringType,
	"Time":             TimeType,
	"UTF8String":       UTF8StringType,
	"Unsigned32":       Unsigned32Type,
	"Unsigned64":       Unsigned64Type,
}
