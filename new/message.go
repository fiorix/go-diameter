// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/fiorix/go-diameter/diam/datatypes"
	"github.com/fiorix/go-diameter/diam/dict"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Diameter message.
type Message struct {
	Header *Header
	AVP    []*AVP

	Dictionary *dict.Parser // Used to encode and decode AVPs.
}

// ReadMessage returns a Message. It uses the dictionary to parse the
// binary stream from the reader.
func ReadMessage(reader io.Reader, dictionary *dict.Parser) (*Message, error) {
	var err error
	hbytes := make([]byte, HeaderLength)
	if _, err = io.ReadFull(reader, hbytes); err != nil {
		return nil, fmt.Errorf("Failed to read Header: %s", err)
	}
	m := &Message{}
	m.Header, err = decodeHeader(hbytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode Header: %s", err)
	}
	pbytes := make([]byte, m.Header.MessageLength-HeaderLength)
	if _, err = io.ReadFull(reader, pbytes); err != nil {
		return nil, fmt.Errorf("Failed to read Payload: %s", err)
	}
	m.Dictionary = dictionary
	return decodeAVPs(m, pbytes)
}

// NewMessage creates and initializes Message.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, dictionary *dict.Parser) *Message {
	if dictionary == nil {
		dictionary = dict.Default
	}
	if hopbyhop == 0 {
		hopbyhop = rand.Uint32()
	}
	if endtoend == 0 {
		endtoend = rand.Uint32()
	}
	return &Message{
		Header: &Header{
			Version:       1,
			MessageLength: HeaderLength,
			CommandFlags:  flags,
			CommandCode:   cmd,
			ApplicationId: appid,
			HopByHopId:    hopbyhop,
			EndToEndId:    endtoend,
		},
		Dictionary: dictionary,
	}
}

// NewAVP creates and initializes a new AVP and adds it to the Message.
func (m *Message) NewAVP(code uint32, flags uint8, vendor uint32, data datatypes.DataType) {
	a := NewAVP(code, flags, vendor, data)
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
}

// AddAVP adds the AVP to the Message.
func (m *Message) AddAVP(a *AVP) {
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
}

func (m *Message) Write(writer io.Writer) (int, error) {
	return writer.Write(m.Serialize())
}

func (m *Message) Serialize() []byte {
	b := make([]byte, m.Len())
	m.SerializeTo(b)
	return b
}

func (m *Message) SerializeTo(b []byte) {
	m.Header.SerializeTo(b[0:HeaderLength])
	offset := HeaderLength
	for _, avp := range m.AVP {
		avp.SerializeTo(b[offset:])
		offset += avp.Len()
	}
}

func (m *Message) Len() int {
	l := HeaderLength
	for _, avp := range m.AVP {
		l += avp.Len()
	}
	return HeaderLength + l
}

func decodeAVPs(m *Message, pbytes []byte) (*Message, error) {
	var avp *AVP
	var err error
	for n := 0; n < cap(pbytes); {
		avp, err = decodeAVP(
			pbytes[n:],
			m.Header.ApplicationId,
			m.Dictionary,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode AVP: %s", err)
		}
		m.AVP = append(m.AVP, avp)
		n += avp.Length + avp.Data.Padding()
		// TODO: Handle grouped AVPs.
	}
	return m, nil
}

func (m *Message) String() string {
	var b bytes.Buffer
	var typ string
	if m.Header.CommandFlags&0x80 > 0 {
		typ = "Request"
	} else {
		typ = "Answer"
	}
	if dictCMD, err := m.Dictionary.FindCMD(
		m.Header.ApplicationId,
		m.Header.CommandCode,
	); err != nil {
		fmt.Fprintf(&b, "Unknown-%s\n%s\n", typ, m.Header)
	} else {
		fmt.Fprintf(&b,
			"%s-%s (%s%c)\n%s\n",
			dictCMD.Name,
			typ,
			dictCMD.Short,
			typ[0],
			m.Header)
	}
	for _, avp := range m.AVP {
		if dictAVP, err := m.Dictionary.FindAVP(
			m.Header.ApplicationId,
			avp.Code,
		); err != nil {
			fmt.Fprintf(&b, "\tUnknown %s\n", err)
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, avp)
		}
	}
	return b.String()
}
