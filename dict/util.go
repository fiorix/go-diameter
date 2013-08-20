// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser, helper functions.  Part of go-diameter.

package dict

import "fmt"

// FindAVP is a helper function that returns a pre-loaded AVP from the Dict.
func (p Parser) FindAVP(appid, code uint32) (*AVP, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if avp, ok := p.avp[index{appid, code}]; ok {
		return avp, nil
	}
	return nil, fmt.Errorf("Could not find preload AVP with code %d", code)
}

// CodeFor is a helper function that returns the code for the given AVP name.
func (p *Parser) CodeFor(name string) uint32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// TODO: Cache this and invalidate when new dict is loaded.
	for avp, v := range p.avp {
		if name == v.Name {
			return avp.Code
		}
	}
	return 0
}

// AppFor is a helper function that returns the DictApp for the given AVP name.
func (p *Parser) AppFor(name string) *App {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// TODO: Cache this and invalidate when new dict is loaded.
	for _, v := range p.avp {
		if name == v.Name {
			return v.App
		}
	}
	return nil // TODO: Return error as well?
}

// Enum is a helper function that returns a pre-loaded DictEnumItem for the
// given AVP code and n. (n is the enum code in the dictionary definition)
func (p *Parser) Enum(appid, code uint32, n uint8) (*EnumItem, error) {
	avp, err := p.FindAVP(appid, code)
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
func (p *Parser) PrettyPrint() {
	fmt.Printf("Vendors:\n")
	for _, f := range p.file {
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

func printAVP(avp *AVP, Parsered bool) {
	var space string
	if Parsered {
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
	// Parsered AVPs
	if len(avp.Data.AVP) > 0 {
		fmt.Printf("  Parsered AVPs:\n")
		for _, ParseredAVP := range avp.Data.AVP {
			printAVP(ParseredAVP, true)
		}
	}
}
