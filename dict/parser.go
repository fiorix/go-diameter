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
	file []*File        // Dict supports multiple XML dictionaries
	avp  map[index]*AVP // AVP index
	cmd  map[index]*Cmd // Cmd index
	mu   sync.RWMutex
}

// index AVPs and Cmds by their application id and code.
type index struct {
	AppId uint32
	Code  uint32
}

// New allocates a new Parser optionally loading dictionary XML files.
// Base Protocol dictionary is always present, and AVPs can be overloaded.
func New(filename ...string) (*Parser, error) {
	p := new(Parser)
	p.avp = make(map[index]*AVP)
	p.cmd = make(map[index]*Cmd)
	p.Load(BaseProtocolXML)
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
			p.cmd[index{app.Id, cmd.Code}] = cmd
		}
		for _, avp := range app.AVP {
			// Link AVP to its Application
			avp.App = app
			p.avp[index{app.Id, avp.Code}] = avp
		}
	}
	return nil
}
