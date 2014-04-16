// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser.  Part of go-diameter.

package dict

import (
	"encoding/xml"
	"fmt"
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
	file    []*File                // Dict supports multiple XML dictionaries
	avpname map[nameIdx]*AVP       // AVP index by name
	avpcode map[codeIdx]*AVP       // AVP index by code
	command map[codeIdx]*Cmd       // Cmd index
	decoder map[string]decoderFunc // AVP data decoders (eg Int32)
	mu      sync.RWMutex           // Protects all maps
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
		p.command = make(map[codeIdx]*Cmd)
		p.decoder = make(map[string]decoderFunc)
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
		for _, cmd := range app.Cmd {
			p.command[codeIdx{app.Id, cmd.Code}] = cmd
		}
		// Cache AVPs.
		for _, avp := range app.AVP {
			// Link AVP to its Application
			avp.App = app
			p.avpname[nameIdx{app.Id, avp.Name}] = avp
			p.avpcode[codeIdx{app.Id, avp.Code}] = avp
			// Set AVP data decoders.
			if err := p.setDecoder(avp.Data.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

// setDecoder updates the decoder map for the given data type.
// This allows applications to easily decode their AVPs by using the
// pre-loaded decoder in the dictionary.
//
// TODO: Allow custom types?
func (p *Parser) setDecoder(datatype string) error {
	if _, exists := p.decoder[datatype]; exists {
		return nil
	}
	var f decoderFunc
	switch datatype {
	case "Address":
		f = datatypes.DecodeAddress
	case "DiameterIdentity":
		f = datatypes.DecodeDiameterIdentity
	case "DiameterURI":
		f = datatypes.DecodeDiameterURI
	case "Enumerated":
		f = datatypes.DecodeEnumerated
	case "Float32":
		f = datatypes.DecodeFloat32
	case "Float64":
		f = datatypes.DecodeFloat64
	case "Grouped":
		f = datatypes.DecodeGrouped
	case "Int32":
		f = datatypes.DecodeInteger32
	case "Int64":
		f = datatypes.DecodeInteger64
	case "IPFilterRule":
		f = datatypes.DecodeIPFilterRule
	case "IPv4":
		f = datatypes.DecodeIPv4
	case "OctetString":
		f = datatypes.DecodeOctetString
	case "Time":
		f = datatypes.DecodeTime
	case "Unsigned32":
		f = datatypes.DecodeUnsigned32
	case "Unsigned64":
		f = datatypes.DecodeUnsigned64
	case "UTF8String":
		f = datatypes.DecodeUTF8String
	default:
		return fmt.Errorf("Unsupported data type: %s", datatype)
	}
	p.decoder[datatype] = f
	return nil
}

func (p *Parser) Decode(datatype string, b []byte) (datatypes.DataType, error) {
	p.mu.RLock()
	f, exists := p.decoder[datatype]
	p.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("Unknown data type: %s", datatype)
	}
	return f(b)
}
