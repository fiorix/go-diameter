// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Dictionary parser.  Part of go-diameter.

package diamdict

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/fiorix/go-diameter/diam/diamtype"
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
	mu      sync.Mutex       // Protects all maps
	once    sync.Once
}

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
	p.mu.Lock()
	defer p.mu.Unlock()
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
			// Check the AVP type.
			if err := updateType(avp); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateType(a *AVP) error {
	id, exists := diamtype.Available[a.Data.TypeName]
	if !exists {
		return fmt.Errorf("Unsupported data type: %s", a.Data.TypeName)
	}
	a.Data.Type = id
	return nil
}

// String returns the Parser represented in a human readable form.
func (p Parser) String() string {
	var b bytes.Buffer
	for _, f := range p.file {
		for _, app := range f.App {
			fmt.Fprintf(&b, "Application Id: %d\n", app.Id)
			fmt.Fprintf(&b, "\tVendors:\n")
			for _, vendor := range app.Vendor {
				fmt.Fprintf(&b, "\t\tId=%d Name=%s\n", vendor.Id, vendor.Name)
			}
			fmt.Fprintf(&b, "\tCommands:\n")
			for _, cmd := range app.CMD {
				printCMD(&b, cmd)
			}
			fmt.Fprintf(&b, "\tAVPs:\n")
			for _, avp := range app.AVP {
				printAVP(&b, avp)
			}
		}
	}
	return b.String()
}

func printCMD(w io.Writer, cmd *CMD) {
	fmt.Fprintf(w, "\t\t%-4d %s-Request (%sR)\n", cmd.Code, cmd.Name, cmd.Short)
	for _, rule := range cmd.Request.Rule {
		if rule.Required && rule.Min == 0 {
			rule.Min = 1
		}
		fmt.Fprintf(w, "\t\t\t% -40s required=%-5t min=%d max=%d\n",
			rule.AVP, rule.Required, rule.Min, rule.Max)
	}
	fmt.Fprintf(w, "\t\t%-4d %s-Answer (%sA)\n", cmd.Code, cmd.Name, cmd.Short)
	for _, rule := range cmd.Answer.Rule {
		if rule.Required && rule.Min == 0 {
			rule.Min = 1
		}
		fmt.Fprintf(w, "\t\t\t% -40s required=%-5t min=%d max=%d\n",
			rule.AVP, rule.Required, rule.Min, rule.Max)
	}
}

func printAVP(w io.Writer, avp *AVP) {
	fmt.Fprintf(w, "\t%-4d %s: %s\n",
		avp.Code, avp.Name, avp.Data.TypeName)
	// Enumerated
	if len(avp.Data.Enum) > 0 {
		fmt.Fprintf(w, "\t\tItems:\n")
		for _, item := range avp.Data.Enum {
			fmt.Fprintf(w, "\t\t\t% -2d %s\n", item.Code, item.Name)
		}
	}
	// Grouped AVPs
	if len(avp.Data.Rule) > 0 {
		fmt.Fprintf(w, "\t\tRules:\n")
		for _, rule := range avp.Data.Rule {
			fmt.Fprintf(w, "\t\t\t% -40s required=%-5t min=%d max=%d\n",
				rule.AVP, rule.Required, rule.Min, rule.Max)
		}
	}
}
