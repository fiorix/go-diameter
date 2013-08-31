// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser.  Part of go-diameter.

package dict

import (
	"encoding/xml"
	"io/ioutil"
	"sync"
)

// Parser is the root element for dictionaries and supports multiple XML
// dictionary files loaded together. Diameter applications use dictionaries
// to parse messages received from peers as well as to encode generated
// messages before sending them over the wire.
//
// Parser can load multiple XML dictionary files, which in turn support
// multiple applications that are composed by multiple AVPs.
//
// The Parser element has an index to make pre-loaded AVPs searcheable per App.
type Parser struct {
	file    []*File          // Dict supports multiple XML dictionaries
	avpname map[nameIdx]*AVP // AVP index by name
	avpcode map[codeIdx]*AVP // AVP index by code
	cmd     map[codeIdx]*Cmd // Cmd index
	mu      sync.RWMutex
}

type codeIdx struct {
	AppId uint32
	Code  uint32
}

type nameIdx struct {
	AppId uint32
	Name  string
}

// New allocates a new Parser optionally loading dictionary XML files.
func New(filename ...string) (*Parser, error) {
	p := new(Parser)
	p.avpname = make(map[nameIdx]*AVP)
	p.avpcode = make(map[codeIdx]*AVP)
	p.cmd = make(map[codeIdx]*Cmd)
	var err error
	for _, f := range filename {
		if err = p.LoadFile(f); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// LoadFile loads a dictionary XML file. May be used multiple times.
func (p *Parser) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return p.Load(buf)
}

// Load loads a dictionary from byte array. May be used multiple times.
func (p *Parser) Load(buf []byte) error {
	f := new(File)
	if err := xml.Unmarshal(buf, f); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.file = append(p.file, f)
	for _, app := range f.App {
		for _, cmd := range app.Cmd {
			p.cmd[codeIdx{app.Id, cmd.Code}] = cmd
		}
		for _, avp := range app.AVP {
			// Link AVP to its Application
			avp.App = app
			p.avpname[nameIdx{app.Id, avp.Name}] = avp
			p.avpcode[codeIdx{app.Id, avp.Code}] = avp
		}
	}
	return nil
}
