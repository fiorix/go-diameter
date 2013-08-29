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
// TODO: Support command short names like CER, CEA.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, d *dict.Parser) *Message {
	if d == nil {
		panic("NewMessage requires a valid dictionary, not nil")
	}
	if hopbyhop == 0 {
		hopbyhop = RandUint32()
	}
	if endtoend == 0 {
		endtoend = RandUint32()
	}
	return &Message{
		Header: &Header{
			Version:        1, // Supports diameter version 1.
			CommandFlags:   flags,
			RawCommandCode: uint32To24(cmd),
			ApplicationId:  appid,
			HopByHopId:     hopbyhop,
			EndToEndId:     endtoend,
		},
		Dict: d,
	}
}

// Answer creates an answer for the current Message with the Result-Code AVP.
func (m *Message) Answer(resultCode uint32) *Message {
	nm := &Message{
		Header: &Header{
			Version:        m.Header.Version,
			CommandFlags:   (m.Header.CommandFlags &^ 0x80),
			RawCommandCode: m.Header.RawCommandCode,
			ApplicationId:  m.Header.ApplicationId,
			HopByHopId:     m.Header.HopByHopId,
			EndToEndId:     m.Header.EndToEndId,
		},
		Dict: m.Dict,
	}
	nm.NewAVP("Result-Code", 0x40, 0x00, &Unsigned32{Value: resultCode})
	return nm
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
