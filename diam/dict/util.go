// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser, helper functions.  Part of go-diameter.

package dict

import (
	"fmt"

	"github.com/fiorix/go-diameter/diam/datatypes"
)

// Apps return a list of all applications loaded in the Parser instance.
// Apps must never be called concurrently with LoadFile or Load.
func (p Parser) Apps() []*App {
	//p.mu.Lock()
	//defer p.mu.Unlock()
	var apps []*App
	for _, f := range p.file {
		for _, app := range f.App {
			apps = append(apps, app)
		}
	}
	return apps
}

// FindAVP is a helper function that returns a pre-loaded AVP from the Dict.
// If the AVP code is not found in the given appid it tries with appid=0
// before returning an error.
// @code can be either the AVP code (int, uint32) or name (string).
//
// FindAVP must never be called concurrently with LoadFile or Load.
func (p Parser) FindAVP(appid uint32, code interface{}) (*AVP, error) {
	//p.mu.Lock()
	//defer p.mu.Unlock()
	var (
		avp *AVP
		ok  bool
		err error
	)
retry:
	switch code.(type) {
	case string:
		avp, ok = p.avpname[nameIdx{appid, code.(string)}]
		if !ok && appid == 0 {
			err = fmt.Errorf("Could not find AVP %s", code.(string))
		}
	case uint32:
		avp, ok = p.avpcode[codeIdx{appid, code.(uint32)}]
		if !ok && appid == 0 {
			err = fmt.Errorf("Could not find AVP %d", code.(uint32))
		}
	case int:
		avp, ok = p.avpcode[codeIdx{appid, uint32(code.(int))}]
		if !ok && appid == 0 {
			err = fmt.Errorf("Could not find AVP %d", code.(int))
		}
	default:
		return nil, fmt.Errorf("Unsupported AVP code type %#v", code)
	}
	if ok {
		return avp, nil
	} else if appid != 0 {
		// Try searching the back dictionary.
		appid = 0
		goto retry
	}
	return nil, err
}

// ScanAVP is a helper function that returns a pre-loaded AVP from the Dict.
// It's similar to FindAPI except that it scans the list of available AVPs
// instead of looking into one specific appid.
//
// ScanAVP is 20x or more slower than FindAVP. Use with care.
// @code can be either the AVP code (uint32) or name (string).
//
// ScanAVP must never be called concurrently with LoadFile or Load.
func (p Parser) ScanAVP(code interface{}) (*AVP, error) {
	//p.mu.Lock()
	//defer p.mu.Unlock()
	switch code.(type) {
	case string:
		for idx, avp := range p.avpname {
			if idx.name == code.(string) {
				return avp, nil
			}
		}
		return nil, fmt.Errorf("Could not find AVP %s", code.(string))
	case uint32:
		for idx, avp := range p.avpcode {
			if idx.code == code.(uint32) {
				return avp, nil
			}
		}
		return nil, fmt.Errorf("Could not find AVP code %d", code.(uint32))
	case int:
		for idx, avp := range p.avpcode {
			if idx.code == uint32(code.(int)) {
				return avp, nil
			}
		}
		return nil, fmt.Errorf("Could not find AVP code %d", code.(int))
	}
	return nil, fmt.Errorf("Unsupported AVP code type %#v", code)
}

// FindCMD is a helper function that returns a pre-loaded CMD from the Dict.
//
// FindCMD must never be called concurrently with LoadFile or Load.
func (p Parser) FindCMD(appid, code uint32) (*CMD, error) {
	//p.mu.Lock()
	//defer p.mu.Unlock()
	if cmd, ok := p.command[codeIdx{appid, code}]; ok {
		return cmd, nil
	} else if cmd, ok = p.command[codeIdx{0, code}]; ok {
		// Always fall back to base dict.
		return cmd, nil
	}
	return nil, fmt.Errorf("Could not find preloaded Cmd with code %d", code)
}

// Enum is a helper function that returns a pre-loaded Enum item for the
// given AVP appid, code and n. (n is the enum code in the dictionary)
//
// Enum must never be called concurrently with LoadFile or Load.
func (p Parser) Enum(appid, code uint32, n uint8) (*Enum, error) {
	avp, err := p.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	if avp.Data.Type != datatypes.EnumeratedType {
		return nil, fmt.Errorf(
			"Data of AVP %s (%d) data is not Enumerated.",
			avp.Name, avp.Code)
	}
	for _, item := range avp.Data.Enum {
		if item.Code == n {
			return item, nil
		}
	}
	return nil, fmt.Errorf(
		"Could not find preload Enum %d for AVP %s (%d)",
		n, avp.Name, avp.Code)
}

// Rule is a helper function that returns a pre-loaded Rule item for the
// given AVP code and name.
//
// Rule must never be called concurrently with LoadFile or Load.
func (p Parser) Rule(appid, code uint32, n string) (*Rule, error) {
	avp, err := p.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	if avp.Data.Type != datatypes.GroupedType {
		return nil, fmt.Errorf(
			"Data of AVP %s (%d) data is not Grouped.",
			avp.Name, avp.Code)
	}
	for _, item := range avp.Data.Rule {
		if item.AVP == n {
			return item, nil
		}
	}
	return nil, fmt.Errorf(
		"Could not find preload Rule for %s for AVP %s (%d)",
		n, avp.Name, avp.Code)
}
