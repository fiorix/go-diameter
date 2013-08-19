// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser and helper functions.

package diam

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sync"
)

type Dict struct {
	file []*DictFile         // Dict supports multiple XML dictionaries
	avp  map[uint32]*DictAVP // AVP index
	mu   sync.RWMutex
}

type DictFile struct {
	XMLName xml.Name      `xml:"diameter"`
	Vendor  []*DictVendor `xml:"vendor"`      // Support for multiple vendors
	App     []*DictApp    `xml:"application"` // Support for multiple applications
}

type DictVendor struct {
	Id   uint32 `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type DictApp struct {
	Id  int        `xml:"id,attr"` // Application Id
	AVP []*DictAVP `xml:"avp"`     // Each application support multiple AVPs
}

type DictAVP struct {
	Name    string   `xml:"name,attr"`
	Code    uint32   `xml:"code,attr"`
	Must    string   `xml:"must,attr"`
	May     string   `xml:"may,attr"`
	MustNot string   `xml:"must-not,attr"`
	Encr    string   `xml:"encr,attr"`
	Data    DictData `xml:"data"`
	App     *DictApp `xml:"none"` // Link back to diameter application
}

type DictData struct {
	Type     string          `xml:"type,attr"`
	EnumItem []*DictEnumItem `xml:"item"` // In case of Enumerated AVP data
	AVP      []*DictAVP      `xml:"avp"`  // In case of Grouped AVPs
}

type DictEnumItem struct {
	Code uint8  `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

// NewDict allocates a new Dict optionally loading dictionary XML files.
// Base Protocol dictionary is always present, and AVPs can be overloaded.
func NewDict(filename ...string) (*Dict, error) {
	dict := new(Dict)
	dict.avp = make(map[uint32]*DictAVP)
	dict.Load(BaseProtocolXML)
	var err error
	for _, f := range filename {
		if err = dict.LoadFile(f); err != nil {
			return nil, err
		}
	}
	return dict, nil
}

// LoadFile loads a dictionary XML file. May be used multiple times.
func (dict *Dict) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return dict.Load(buf)
}

// Load loads a dictionary from byte array. May be used multiple times.
func (dict *Dict) Load(buf []byte) error {
	f := new(DictFile)
	if err := xml.Unmarshal(buf, f); err != nil {
		return err
	}
	dict.mu.Lock()
	defer dict.mu.Unlock()
	dict.file = append(dict.file, f)
	for _, app := range f.App {
		for _, avp := range app.AVP {
			// Link AVP to its Application
			avp.App = app
			// TODO: Perhaps index by ApplicationId + AVP code?
			dict.avp[avp.Code] = avp
		}
	}
	return nil
}

// AVP is a helper function that returns a pre-loaded AVP from the Dict.
func (dict *Dict) AVP(code uint32) (*DictAVP, error) {
	dict.mu.RLock()
	defer dict.mu.RUnlock()
	if avp, ok := dict.avp[code]; ok {
		return avp, nil
	}
	return nil, fmt.Errorf("Could not find preload AVP with code %d", code)
}

// CodeFor is a helper function that returns the code for the given AVP name.
func (dict *Dict) CodeFor(name string) uint32 {
	dict.mu.RLock()
	defer dict.mu.RUnlock()
	// TODO: Cache this and invalidate when new dict is loaded.
	for code, v := range dict.avp {
		if name == v.Name {
			return code
		}
	}
	return 0
}

// AppFor is a helper function that returns the DictApp for the given AVP name.
func (dict *Dict) AppFor(name string) *DictApp {
	dict.mu.RLock()
	defer dict.mu.RUnlock()
	// TODO: Cache this and invalidate when new dict is loaded.
	for _, v := range dict.avp {
		if name == v.Name {
			return v.App
		}
	}
	return nil // TODO: Return error as well?
}

// Enum is a helper function that returns a pre-loaded DictEnumItem for the
// given AVP code and n. (n is the enum code in the dictionary definition)
func (dict *Dict) Enum(code uint32, n uint8) (*DictEnumItem, error) {
	avp, err := dict.AVP(code)
	if err != nil {
		return nil, err
	}
	if avp.Data.Type != "Enumerated" {
		return nil, fmt.Errorf("AVP %s (%d) data is not Enumerated.", avp.Name, avp.Code)
	}
	for _, item := range avp.Data.EnumItem {
		if item.Code == n {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Could not find preload Enum %d for AVP %s (%d)", n, avp.Name, avp.Code)
}

// PrettyPrint prints the Dict in a human readable form.
func (dict *Dict) PrettyPrint() {
	fmt.Printf("Vendors:\n")
	for _, f := range dict.file {
		for _, vendor := range f.Vendor {
			fmt.Printf("Id=%d Name=%s\n", vendor.Id, vendor.Name)
		}
		fmt.Println()
		for _, app := range f.App {
			fmt.Printf("Application Id: %d\nAVPs:\n", app.Id)
			for _, avp := range app.AVP {
				printAVP(avp, false)
			}
		}
	}
}

func printAVP(avp *DictAVP, grouped bool) {
	var space string
	if grouped {
		space = "  "
	}
	fmt.Printf("%s%s AVP{Code=%d,Type=%s}\n",
		space, avp.Name, avp.Code, avp.Data.Type)
	// Enumerated
	if len(avp.Data.EnumItem) > 0 {
		fmt.Printf("  Items:\n")
		for _, item := range avp.Data.EnumItem {
			fmt.Printf("  %d %s\n", item.Code, item.Name)
		}
	}
	// Grouped AVPs
	if len(avp.Data.AVP) > 0 {
		fmt.Printf("  Grouped AVPs:\n")
		for _, groupedAVP := range avp.Data.AVP {
			printAVP(groupedAVP, true)
		}
	}
}
