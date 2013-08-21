// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary file structure.  Part of go-diameter.

package dict

import (
	"encoding/xml"
)

// File is the dictionary root element of a XML file.  See diam_base.xml.
type File struct {
	XMLName xml.Name  `xml:"diameter"`
	Vendor  []*Vendor `xml:"vendor"`      // Support for multiple vendors
	App     []*App    `xml:"application"` // Support for multiple applications
}

// Vendor defines diameter vendors in XML, that can be used to translate
// the VendorId AVP of incoming messages.
type Vendor struct {
	Id   uint32 `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// App defines a diameter application in XML and its multiple AVPs.
type App struct {
	Id  uint32 `xml:"id,attr"` // Application Id
	Cmd []*Cmd `xml:"command"` // Diameter commands
	AVP []*AVP `xml:"avp"`     // Each application support multiple AVPs
}

// Cmd defines a diameter command (CE, CC, etc)
type Cmd struct {
	Code  uint32 `xml:"code,attr"`
	Name  string `xml:"name,attr"`
	Short string `xml:"short,attr"`
}

// AVP represents a dictionary AVP that is loaded from XML.
type AVP struct {
	Name    string `xml:"name,attr"`
	Code    uint32 `xml:"code,attr"`
	Must    string `xml:"must,attr"`
	May     string `xml:"may,attr"`
	MustNot string `xml:"must-not,attr"`
	Encr    string `xml:"encr,attr"`
	Data    Data   `xml:"data"`
	App     *App   `xml:"none"` // Link back to diameter application
}

// Data of an AVP can be EnumItem or a Parser of multiple AVPs.
type Data struct {
	Type     string      `xml:"type,attr"`
	EnumItem []*EnumItem `xml:"item"` // In case of Enumerated AVP data
	AVP      []*AVP      `xml:"avp"`  // In case of Parsered AVPs
}

type EnumItem struct {
	Code uint8  `xml:"code,attr"`
	Name string `xml:"name,attr"`
}
