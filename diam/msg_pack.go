// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"math/rand"

	"github.com/fiorix/go-diameter/diam/avpdata"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/util"
)

// NewMessage allocates a new Message object. Used for building messages
// that will be sent to a connection later.
//
// Arguments hopbyhop and endtoend are optional. If set to 0, random values
// are used.
//
// Dictionary d is used for encoding the AVPs added to the message. If set to
// nil, static dict.Default is used.
//
// TODO: Support command short names like CER, CEA.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, d *dict.Parser) *Message {
	if d == nil {
		d = dict.Default
	}
	if hopbyhop == 0 {
		hopbyhop = rand.Uint32()
	}
	if endtoend == 0 {
		endtoend = rand.Uint32()
	}
	return &Message{
		Header: &Header{
			Version:        1, // Supports diameter version 1.
			CommandFlags:   flags,
			RawCommandCode: util.Uint32To24(cmd),
			ApplicationId:  appid,
			HopByHopId:     hopbyhop,
			EndToEndId:     endtoend,
		},
		Dict: d,
	}
}

// Answer creates an answer for the current Message with an embedded
// Result-Code AVP.
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
	nm.NewAVP("Result-Code", 0x40, 0x00,
		&avpdata.Unsigned32{Value: resultCode})
	return nm
}

// Add adds an AVP to the Message and set its parent Message to the current.
func (m *Message) Add(avp *AVP) {
	// This is required when copying AVPs from one Message to another.
	if avp.msg != m {
		avp.msg = m
	}
	m.AVP = append(m.AVP, avp)
}

// Pack returns the Message in binary form to be sent to a connection.
func (m *Message) Pack() []byte {
	var buf [][]byte
	var length uint32
	for _, avp := range m.AVP {
		b := avp.Pack()
		length += uint32(len(b))
		buf = append(buf, b)
	}
	m.Header.SetMessageLength(length)
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, m.Header)
	binary.Write(b, binary.BigEndian, bytes.Join(buf, []byte{}))
	return b.Bytes()
}
