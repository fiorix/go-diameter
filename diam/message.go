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
	"sync"
	"time"

	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/avp/format"
	"github.com/fiorix/go-diameter/diam/dict"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Diameter message.
type Message struct {
	Header *Header
	AVP    []*AVP // AVPs in this message.

	// dictionary parser object used to encode and decode AVPs.
	dictionary *dict.Parser
}

// ReadMessage returns a Message. It uses the dictionary to parse the
// binary stream from the reader.
func ReadMessage(reader io.Reader, dictionary *dict.Parser) (*Message, error) {
	buf := newMessageBuffer()
	defer putMessageBuffer(buf)
	m := &Message{dictionary: dictionary}
	cmd, err := readAndParseHeader(m, buf, reader)
	if err != nil {
		return nil, err
	}
	if err = readAndParseBody(m, cmd, buf, reader); err != nil {
		return nil, err
	}
	return m, nil
}

var (
	messageBufferPool   sync.Pool
	MessageBufferLength = 1 << 10 // Default 1K per message.
)

func newMessageBuffer() *bytes.Buffer {
	if v := messageBufferPool.Get(); v != nil {
		return v.(*bytes.Buffer)
	}
	return bytes.NewBuffer(make([]byte, MessageBufferLength))
}

func putMessageBuffer(b *bytes.Buffer) {
	b.Reset()
	if cap(b.Bytes()) >= MessageBufferLength {
		messageBufferPool.Put(b)
	}
}

func messageBytes(b *bytes.Buffer, l int) []byte {
	p := b.Bytes()
	if l <= MessageBufferLength && cap(p) >= MessageBufferLength {
		return p[0:l]
	}
	return make([]byte, l)
}

func readAndParseHeader(m *Message, buf *bytes.Buffer, r io.Reader) (cmd *dict.Command, err error) {
	b := buf.Bytes()[:HeaderLength]
	if _, err = io.ReadFull(r, b); err != nil {
		return nil, err
	}
	m.Header, err = DecodeHeader(b)
	if err != nil {
		return nil, err
	}
	cmd, err = m.dictionary.FindCommand(
		m.Header.ApplicationId,
		m.Header.CommandCode,
	)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func readAndParseBody(m *Message, cmd *dict.Command, buf *bytes.Buffer, r io.Reader) error {
	b := messageBytes(buf, int(m.Header.MessageLength-HeaderLength))
	_, err := io.ReadFull(r, b)
	if err != nil {
		return err
	}
	if n := maxAVPs(m, cmd); n == 0 {
		// TODO: fail to load the dictionary instead.
		return fmt.Errorf(
			"Command %s (%d) has no AVPs defined in the dictionary.",
			cmd.Name, cmd.Code)
	} else {
		// Pre-allocate max # of AVPs for this message.
		m.AVP = make([]*AVP, 0, n)
	}
	if err = decodeAVPs(m, b); err != nil {
		return err
	}
	return nil
}

func maxAVPs(m *Message, cmd *dict.Command) int {
	if m.Header.CommandFlags&RequestFlag > 0 {
		return len(cmd.Request.Rule)
	} else {
		return len(cmd.Answer.Rule)
	}
}

func decodeAVPs(m *Message, pbytes []byte) error {
	var a *AVP
	var err error
	for n := 0; n < len(pbytes); {
		a, err = DecodeAVP(
			pbytes[n:],
			m.Header.ApplicationId,
			m.dictionary,
		)
		if err != nil {
			return fmt.Errorf("Failed to decode AVP: %s", err)
		}
		m.AVP = append(m.AVP, a)
		n += a.Len()
	}
	return nil
}

// NewMessage creates and initializes a Message.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, dictionary *dict.Parser) *Message {
	if dictionary == nil {
		dictionary = dict.Default // for safety.
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
		dictionary: dictionary,
	}
}

// NewRequest creates a new Message with the Request bit set.
func NewRequest(cmd uint32, appid uint32, dictionary *dict.Parser) *Message {
	return NewMessage(cmd, RequestFlag, appid, 0, 0, dictionary)
}

// Dictionary returns the dictionary parser object associated with this message.
// This dictionary is used to encode and decode the message.
func (m *Message) Dictionary() *dict.Parser { return m.dictionary }

// NewAVP creates and initializes a new AVP and adds it to the Message.
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data format.Format) (*AVP, error) {
	var a *AVP
	switch code.(type) {
	case int:
		a = NewAVP(uint32(code.(int)), flags, vendor, data)
	case uint32:
		a = NewAVP(code.(uint32), flags, vendor, data)
	case string:
		dictAVP, err := m.dictionary.FindAVP(
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

// WriteTo serializes the Message and writes into the writer.
func (m *Message) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write(m.Serialize())
	if err != nil {
		return 0, err
	}
	return int64(n), err
}

// Serialize returns the serialized bytes of the Message.
func (m *Message) Serialize() []byte {
	b := make([]byte, m.Len())
	m.SerializeTo(b)
	return b
}

// SerializeTo writes the serialized bytes of the Message into b.
func (m *Message) SerializeTo(b []byte) {
	m.Header.SerializeTo(b[0:HeaderLength])
	offset := HeaderLength
	for _, avp := range m.AVP {
		avp.SerializeTo(b[offset:])
		offset += avp.Len()
	}
}

// Len returns the length of the Message in bytes.
func (m *Message) Len() int {
	l := HeaderLength
	for _, avp := range m.AVP {
		l += avp.Len()
	}
	return l
}

// FindAVP searches the Message for a specific AVP.
// The code can be either the AVP code (int, uint32) or name (string).
//
// Example:
//
//	avp, err := m.FindAVP(264)
//	avp, err := m.FindAVP("Origin-Host")
func (m *Message) FindAVP(code interface{}) (*AVP, error) {
	dictAVP, err := m.dictionary.FindAVP(m.Header.ApplicationId, code)
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
		m.Header.CommandFlags&^RequestFlag, // Reset the Request bit.
		m.Header.ApplicationId,
		m.Header.HopByHopId,
		m.Header.EndToEndId,
		m.dictionary,
	)
	nm.NewAVP(avp.ResultCode, avp.Mbit, 0, format.Unsigned32(resultCode))
	return nm
}

func (m *Message) String() string {
	var b bytes.Buffer
	var typ string
	if m.Header.CommandFlags&RequestFlag > 0 {
		typ = "Request"
	} else {
		typ = "Answer"
	}
	if dictCMD, err := m.dictionary.FindCommand(
		m.Header.ApplicationId,
		m.Header.CommandCode,
	); err != nil {
		fmt.Fprintf(&b, "Unknown-%s\n%s\n", typ, m.Header)
	} else {
		fmt.Fprintf(&b, "%s-%s (%s%c)\n%s\n",
			dictCMD.Name,
			typ,
			dictCMD.Short,
			typ[0],
			m.Header,
		)
	}
	for _, a := range m.AVP {
		if dictAVP, err := m.dictionary.FindAVP(
			m.Header.ApplicationId,
			a.Code,
		); err != nil {
			fmt.Fprintf(&b, "\tUnknown %s (%s)\n", a, err)
		} else if a.Data.Format() == GroupedAVPFormat {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, printGrouped("\t", m, a))
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, a)
		}
	}
	return b.String()
}

func printGrouped(prefix string, m *Message, a *AVP) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:Grouped{\n",
		a.Code,
		a.Flags,
		a.Len(),
		a.VendorId,
	)
	for _, ga := range a.Data.(*GroupedAVP).AVP {
		if dictAVP, err := m.dictionary.FindAVP(
			m.Header.ApplicationId,
			ga.Code,
		); err != nil {
			fmt.Fprintf(&b, "%s\tUnknown %s (%s),\n", prefix, ga, err)
		} else {
			fmt.Fprintf(&b, "%s\t%s %s,\n", prefix, dictAVP.Name, ga)
		}
	}
	fmt.Fprintf(&b, "%s}}", prefix)
	return b.String()
}
