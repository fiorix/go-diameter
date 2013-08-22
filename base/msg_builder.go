// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter message, header with multiple AVPs.  Part of go-diameter.

package base

import (
	"bytes"
	"encoding/binary"

	"github.com/fiorix/go-diameter/dict"
)

// NewMessage allocates a new Message. Used for building messages that will
// be sent to a connection.
func NewMessage(code uint32, flags uint8, appid, hopbyhop, endtoend uint32, d *dict.Parser) *Message {
	return &Message{
		Header: &Header{
			Version:        1, // Supports diameter version 1.
			CommandFlags:   flags,
			RawCommandCode: uint32To24(code),
			HopByHopId:     hopbyhop,
			EndToEndId:     endtoend,
		},
		Dict: d,
	}
}

// Add adds an AVP to the given Message.
func (m *Message) Add(avp *AVP) {
	// Set AVP's parent Message to this.
	// This is required when copying AVPs from one Message to another.
	if avp.Message != m {
		avp.Message = m
	}
	m.AVP = append(m.AVP, avp)
}

// Bytes returns the Message in binary form to be sent to a connection.
func (m *Message) Bytes() []byte {
	var buf [][]byte
	var length uint32
	for _, avp := range m.AVP {
		b := avp.Bytes()
		bl := uint32(len(b))
		length += bl
		buf = append(buf, b)
	}
	m.Header.SetMessageLength(length)
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, m.Header)
	binary.Write(b, binary.BigEndian, bytes.Join(buf, []byte{}))
	return b.Bytes()
}
