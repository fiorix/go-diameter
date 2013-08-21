// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter message, header with multiple AVPs.  Part of go-diameter.

package base

import (
	"fmt"

	"github.com/fiorix/go-diameter/dict"
)

// Message is a diameter message, composed of a header and multiple AVPs.
type Message struct {
	Header *Header
	AVP    []*AVP
	Dict   *dict.Parser // Dictionary associated with this Message
}

// FindAVP is a helper that returns an AVP by looking up its code, or nil if
// the AVP is not in the message.
func (m *Message) FindAVP(code uint32) *AVP {
	for _, a := range m.AVP {
		if code == a.Code {
			return a
		}
	}
	return nil
}

// PrettyPrint prints messages in a human readable format.
func (m *Message) PrettyPrint() {
	fmt.Println(m.String())
	for _, avp := range m.AVP {
		fmt.Printf("  %s\n", avp)
	}
	fmt.Println()
}

// String returns a human readable version of the Message header.
func (m *Message) String() string {
	rflag := m.Header.CommandFlags&0x80 > 0
	pflag := m.Header.CommandFlags&0x40 > 0
	eflag := m.Header.CommandFlags&0x20 > 0
	tflag := m.Header.CommandFlags&0x10 > 0
	cmdName, cmdShort := findCmd(m.Dict, m.Header)
	return fmt.Sprintf(
		"%s (%s) Header{Code=%d,Version=%d,"+
			"MessageLength=%d,CommandFlags={r=%v,p=%v,e=%v,t=%v},"+
			"ApplicationId=%d,HopByHopId=%#v,EndToEndId=%#v}",
		cmdName,
		cmdShort,
		m.Header.CommandCode(),
		m.Header.Version,
		m.Header.MessageLength(),
		rflag, pflag, eflag, tflag,
		m.Header.ApplicationId,
		m.Header.HopByHopId,
		m.Header.EndToEndId,
	)
}

func findCmd(d *dict.Parser, h *Header) (string, string) {
	var cmdName, cmdShort string
	if d != nil {
		cmd, err := d.FindCmd(h.ApplicationId, h.CommandCode())
		if err == nil {
			cmdName = cmd.Name
			cmdShort = cmd.Short
		}
	}
	if cmdName == "" {
		cmdName, cmdShort = "Unknown", ""
	}
	if h.CommandFlags&0x80 > 0 {
		cmdName += "-Request"
		cmdShort += "R"
	} else {
		cmdName += "-Answer"
		cmdShort += "A"
	}
	return cmdName, cmdShort
}
