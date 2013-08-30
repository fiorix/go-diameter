// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter message, header with multiple AVPs.  Part of go-diameter.

package diam

import (
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/dict"
)

var ErrNotFound = errors.New("Not found")

// Message is a diameter message, composed of a header and multiple AVPs.
type Message struct {
	Header *Header
	AVP    []*AVP
	Dict   *dict.Parser // Dictionary associated with this Message
}

// FindAVP is a helper that returns an AVP by looking up its code, or nil if
// the AVP is not in the message.
// @code can be either the AVP code (int, uint32) or name (string).
func (m *Message) FindAVP(code interface{}) (*AVP, error) {
	davp, err := m.Dict.FindAVP(m.Header.ApplicationId, code)
	if err != nil {
		return nil, err
	}
	for _, a := range m.AVP {
		if davp.Code == a.Code {
			return a, nil
		}
	}
	var (
		name string
		avp  *dict.AVP
	)
	if avp, err = m.Dict.ScanAVP(code); err != nil {
		name = "Unknown"
	} else {
		name = avp.Name
	}
	return nil, fmt.Errorf("AVP %d (%s) not found", code, name)
}

// PrettyPrint prints messages in a human readable format.
func (m *Message) PrettyPrint() {
	// Update header length and other fields.
	m.Bytes()
	fmt.Println(m.String())
	for _, avp := range m.AVP {
		fmt.Printf("  %s\n", avp)
	}
	fmt.Println()
}

// String returns a human readable version of the Message header.
func (m *Message) String() string {
	cmdName, cmdShort := findCmd(m.Dict, m.Header)
	return fmt.Sprintf(
		"%s (%s) Header{Code=%d,Version=%d,"+
			"MessageLength=%d,CommandFlags=%#v,"+
			"ApplicationId=%d,HopByHopId=%#v,EndToEndId=%#v}",
		cmdName,
		cmdShort,
		m.Header.CommandCode(),
		m.Header.Version,
		m.Header.MessageLength(),
		m.Header.CommandFlags,
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
