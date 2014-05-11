// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/fiorix/go-diameter/diam/diamdict"
	"github.com/fiorix/go-diameter/diam/diamtype"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Diameter message.
type Message struct {
	Header *Header
	AVP    []*AVP

	Dictionary *diamdict.Parser // Used to encode and decode AVPs.
}

// ReadMessage returns a Message. It uses the dictionary to parse the
// binary stream from the reader.
func ReadMessage(reader io.Reader, dictionary *diamdict.Parser) (*Message, error) {
	var err error
	hbytes := make([]byte, HeaderLength)
	if _, err = io.ReadFull(reader, hbytes); err != nil {
		return nil, err
	}
	m := &Message{}
	m.Header, err = decodeHeader(hbytes)
	if err != nil {
		return nil, err
	}
	pbytes := make([]byte, m.Header.MessageLength-HeaderLength)
	if _, err = io.ReadFull(reader, pbytes); err != nil {
		return nil, err
	}
	m.Dictionary = dictionary
	return decodeAVPs(m, pbytes)
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
		n += avp.Len()
	}
	return m, nil
}

// NewMessage creates and initializes Message.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, dictionary *diamdict.Parser) *Message {
	if dictionary == nil {
		dictionary = diamdict.Default
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

// NewRequest is an alias to NewMessage.
func NewRequest(cmd uint32, appid uint32, dictionary *diamdict.Parser) *Message {
	return NewMessage(cmd, 0x80, appid, 0, 0, dictionary)
}

// NewAVP creates and initializes a new AVP and adds it to the Message.
// @code can be int, uint32 or string (e.g. 268 or Result-Code)
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data diamtype.DataType) (*AVP, error) {
	var a *AVP
	switch code.(type) {
	case int:
		a = NewAVP(uint32(code.(int)), flags, vendor, data)
	case uint32:
		a = NewAVP(code.(uint32), flags, vendor, data)
	case string:
		dictAVP, err := m.Dictionary.FindAVP(
			m.Header.ApplicationId,
			code.(string),
		)
		if err != nil {
			return nil, err
		}
		a = NewAVP(dictAVP.Code, flags, vendor, data)
	}
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
	return a, nil
}

// AddAVP adds the AVP to the Message.
func (m *Message) AddAVP(a *AVP) {
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
}

func (m *Message) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write(m.Serialize())
	return int64(n), err
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
	return l
}

// FindAVP searches the Message for a specific AVP.
// @code can be either the AVP code (int, uint32) or name (string).
//
// Example:
//
//	avp, err := m.FindAVP(264)
//	avp, err := m.FindAVP("Origin-Host")
func (m *Message) FindAVP(code interface{}) (*AVP, error) {
	dictAVP, err := m.Dictionary.FindAVP(m.Header.ApplicationId, code)
	if err != nil {
		return nil, err
	}
	for _, a := range m.AVP {
		if a.Code == dictAVP.Code {
			return a, nil
		}
	}
	return nil, errors.New("Not found")
}

// Answer creates an answer for the current Message with an embedded
// Result-Code AVP.
func (m *Message) Answer(resultCode uint32) *Message {
	nm := NewMessage(
		m.Header.CommandCode,
		m.Header.CommandFlags&^0x80,
		m.Header.ApplicationId,
		m.Header.HopByHopId,
		m.Header.EndToEndId,
		m.Dictionary,
	)
	nm.NewAVP(268, 0x40, 0x00, diamtype.Unsigned32(resultCode))
	return nm
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
			fmt.Fprintf(&b, "\tUnknown %s (%s)\n", avp, err)
		} else if avp.Data.Type() == GroupedType {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, printGrouped("\t", m, avp))
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, avp)
		}
	}
	return b.String()
}

func printGrouped(prefix string, m *Message, avp *AVP) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:Grouped{\n",
		avp.Code,
		avp.Flags,
		avp.Len(),
		avp.VendorId,
	)
	for _, a := range avp.Data.(*Grouped).AVP {
		if dictAVP, err := m.Dictionary.FindAVP(
			m.Header.ApplicationId,
			a.Code,
		); err != nil {
			fmt.Fprintf(&b, "%s\tUnknown %s (%s),\n", prefix, avp, err)
		} else {
			fmt.Fprintf(&b, "%s\t%s %s,\n", prefix, dictAVP.Name, a)
		}
	}
	fmt.Fprintf(&b, "%s}}", prefix)
	return b.String()
}
