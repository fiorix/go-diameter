// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Skeleton of the dictionary file.  Part of go-diameter.

package dict

import (
	"encoding/xml"

	"github.com/fiorix/go-diameter/diam/datatypes"
)

// File is the dictionary root element of a XML file.  See diam_base.xml.
type File struct {
	XMLName xml.Name `xml:"diameter"`
	App     []*App   `xml:"application"` // Support for multiple applications
}

// App defines a diameter application in XML and its multiple AVPs.
type App struct {
	Id     uint32    `xml:"id,attr"`   // Application Id
	Type   string    `xml:"type,attr"` // Application type
	Name   string    `xml:"name,attr"` // Application name
	Vendor []*Vendor `xml:"vendor"`    // Support for multiple vendors
	CMD    []*CMD    `xml:"command"`   // Diameter commands
	AVP    []*AVP    `xml:"avp"`       // Each application support multiple AVPs
}

// Vendor defines diameter vendors in XML, that can be used to translate
// the VendorId AVP of incoming messages.
type Vendor struct {
	Id   uint32 `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// CMD defines a diameter command (CE, CC, etc)
type CMD struct {
	Code    uint32  `xml:"code,attr"`
	Name    string  `xml:"name,attr"`
	Short   string  `xml:"short,attr"`
	Request CMDRule `xml:"request"`
	Answer  CMDRule `xml:"answer"`
}

type CMDRule struct {
	Rule []*Rule `xml:"rule"`
}

// AVP represents a dictionary AVP that is loaded from XML.
type AVP struct {
	Name       string `xml:"name,attr"`
	Code       uint32 `xml:"code,attr"`
	Must       string `xml:"must,attr"`
	May        string `xml:"may,attr"`
	MustNot    string `xml:"must-not,attr"`
	MayEncrypt string `xml:"may-encrypt,attr"`
	Data       Data   `xml:"data"`
	App        *App   `xml:"none"` // Link back to diameter application
}

// Data of an AVP can be EnumItem or a Parser of multiple AVPs.
type Data struct {
	Type     datatypes.DataTypeId `xml:"-"`
	TypeName string               `xml:"type,attr"`
	Enum     []*Enum              `xml:"item"` // In case of Enumerated AVP data
	Rule     []*Rule              `xml:"rule"` // In case of Grouped AVPs
}

type Enum struct {
	Code uint8  `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Rule struct {
	AVP      string `xml:"avp,attr"` // AVP Name
	Required bool   `xml:"required,attr"`
	Min      int    `xml:"min,attr"`
	Max      int    `xml:"max,attr"`
}
