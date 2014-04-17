// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser.  Part of go-diameter.

package dict

import (
	"encoding/xml"
	"io"
	"os"
	"sync"

	"github.com/fiorix/go-diameter/diam/datatypes"
)

// Parser is the root element for dictionaries and supports multiple XML
// dictionary files loaded together. Diameter applications use dictionaries
// to parse messages received from peers as well as to encode crafted
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
	command map[codeIdx]*CMD // CMD index
	mu      sync.RWMutex     // Protects all maps
	once    sync.Once
}

type decoderFunc func([]byte) (datatypes.DataType, error)

type codeIdx struct {
	appId uint32
	code  uint32
}

type nameIdx struct {
	appId uint32
	name  string
}

// New allocates a new Parser optionally loading dictionary XML files.
func NewParser(filename ...string) (*Parser, error) {
	p := new(Parser)
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
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return p.Load(fd)
}

// Load loads a dictionary from byte array. May be used multiple times.
func (p *Parser) Load(r io.Reader) error {
	p.once.Do(func() {
		p.avpname = make(map[nameIdx]*AVP)
		p.avpcode = make(map[codeIdx]*AVP)
		p.command = make(map[codeIdx]*CMD)
	})
	f := new(File)
	d := xml.NewDecoder(r)
	if err := d.Decode(f); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.file = append(p.file, f)
	for _, app := range f.App {
		// Cache commands.
		for _, cmd := range app.CMD {
			p.command[codeIdx{app.Id, cmd.Code}] = cmd
		}
		// Cache AVPs.
		for _, avp := range app.AVP {
			// Link AVP to its Application
			avp.App = app
			p.avpname[nameIdx{app.Id, avp.Name}] = avp
			p.avpcode[codeIdx{app.Id, avp.Code}] = avp
		}
	}
	return nil
}
