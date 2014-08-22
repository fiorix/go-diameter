// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package format

// Format is an interface for supporting multiple AVP data formats.
type Format interface {
	Serialize() []byte
	Len() int
	Padding() int
	Format() FormatId
	String() string
}

type FormatId int

const (
	UnknownFormat FormatId = iota
	AddressFormat
	DiameterIdentityFormat
	DiameterURIFormat
	EnumeratedFormat
	Float32Format
	Float64Format
	GroupedFormat
	IPFilterRuleFormat
	IPv4Format
	Integer32Format
	Integer64Format
	OctetStringFormat
	TimeFormat
	UTF8StringFormat
	Unsigned32Format
	Unsigned64Format
)

var Available = map[string]FormatId{
	"Address":          AddressFormat,
	"DiameterIdentity": DiameterIdentityFormat,
	"DiameterURI":      DiameterURIFormat,
	"Enumerated":       EnumeratedFormat,
	"Float32":          Float32Format,
	"Float64":          Float64Format,
	"Grouped":          GroupedFormat,
	"IPFilterRule":     IPFilterRuleFormat,
	"IPv4":             IPv4Format,
	"Integer32":        Integer32Format,
	"Integer64":        Integer64Format,
	"OctetString":      OctetStringFormat,
	"Time":             TimeFormat,
	"UTF8String":       UTF8StringFormat,
	"Unsigned32":       Unsigned32Format,
	"Unsigned64":       Unsigned64Format,
}
