// Copyright 2013-2015 go-diameter authors. All rights reserved.
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
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// MessageBufferLength is the default buffer length for Diameter messages.
var MessageBufferLength = 1 << 10

// Message represents a Diameter message.
type Message struct {
	Header *Header
	AVP    []*AVP // AVPs in this message.

	// dictionary parser object used to encode and decode AVPs.
	dictionary *dict.Parser
}

// ReadMessage returns a Message. It uses the dictionary to parse the
// binary stream from the reader.
func ReadMessage(reader io.Reader, dictionary *dict.Parser) (*Message, error) {
	buf := newReaderBuffer()
	defer putReaderBuffer(buf)
	m := &Message{dictionary: dictionary}
	cmd, err := readAndParseHeader(reader, buf, m)
	if err != nil {
		return nil, err
	}
	if err = readAndParseBody(reader, buf, cmd, m); err != nil {
		return nil, err
	}
	return m, nil
}

var readerBufferPool sync.Pool

func newReaderBuffer() *bytes.Buffer {
	if v := readerBufferPool.Get(); v != nil {
		return v.(*bytes.Buffer)
	}
	return bytes.NewBuffer(make([]byte, MessageBufferLength))
}

func putReaderBuffer(b *bytes.Buffer) {
	if cap(b.Bytes()) == MessageBufferLength {
		b.Reset()
		readerBufferPool.Put(b)
	}
}

func readerBufferSlice(b *bytes.Buffer, l int) []byte {
	p := b.Bytes()
	if l <= MessageBufferLength && cap(p) >= MessageBufferLength {
		return p[0:l]
	}
	return make([]byte, l)
}

func readAndParseHeader(r io.Reader, buf *bytes.Buffer, m *Message) (cmd *dict.Command, err error) {
	b := buf.Bytes()[:HeaderLength]
	if _, err = io.ReadFull(r, b); err != nil {
		return nil, err
	}
	m.Header, err = DecodeHeader(b)
	if err != nil {
		return nil, err
	}
	cmd, err = m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode,
	)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func readAndParseBody(r io.Reader, buf *bytes.Buffer, cmd *dict.Command, m *Message) error {
	b := readerBufferSlice(buf, int(m.Header.MessageLength-HeaderLength))
	_, err := io.ReadFull(r, b)
	if err != nil {
		return err
	}
	n := maxAVPs(m, cmd)
	if n == 0 {
		// TODO: fail to load the dictionary instead.
		return fmt.Errorf(
			"Command %s (%d) has no AVPs defined in the dictionary.",
			cmd.Name, cmd.Code)
	}
	// Pre-allocate max # of AVPs for this message.
	m.AVP = make([]*AVP, 0, n)
	if err = decodeAVPs(m, b); err != nil {
		return err
	}
	return nil
}

func maxAVPs(m *Message, cmd *dict.Command) int {
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		return len(cmd.Request.Rule)
	}
	return len(cmd.Answer.Rule)
}

func decodeAVPs(m *Message, pbytes []byte) error {
	var a *AVP
	var err error
	for n := 0; n < len(pbytes); {
		a, err = DecodeAVP(pbytes[n:],
			m.Header.ApplicationID, m.Dictionary())
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
			ApplicationID: appid,
			HopByHopID:    hopbyhop,
			EndToEndID:    endtoend,
		},
		dictionary: dictionary,
	}
}

// NewRequest creates a new Message with the Request bit set.
func NewRequest(cmd uint32, appid uint32, dictionary *dict.Parser) *Message {
	return NewMessage(cmd, RequestFlag, appid, 0, 0, dictionary)
}

// Dictionary returns the dictionary parser object associated with this
// message. This dictionary is used to encode and decode the message.
// If no dictionary is associated then it returns the default dictionary.
func (m *Message) Dictionary() *dict.Parser {
	if m.dictionary == nil {
		return dict.Default
	}
	return m.dictionary
}

// NewAVP creates and initializes a new AVP and adds it to the Message.
// It is not safe for concurrent calls.
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data datatype.Type) (*AVP, error) {
	var a *AVP
	switch code.(type) {
	case int:
		a = NewAVP(uint32(code.(int)), flags, vendor, data)
	case uint32:
		a = NewAVP(code.(uint32), flags, vendor, data)
	case string:
		dictAVP, err := m.Dictionary().FindAVP(
			m.Header.ApplicationID,
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

// AddAVP adds the AVP to the Message. It is not safe for concurrent calls.
func (m *Message) AddAVP(a *AVP) {
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
}

// WriteTo serializes the Message and writes into the writer.
func (m *Message) WriteTo(writer io.Writer) (int64, error) {
	l := m.Len()
	buf := newWriterBuffer(l)
	defer putWriterBuffer(buf)
	b := buf.Bytes()[0:l]
	if err := m.SerializeTo(b); err != nil {
		return 0, err
	}
	n, err := writer.Write(b)
	if err != nil {
		return 0, err
	}
	return int64(n), err
}

var writerBufferPool sync.Pool

func newWriterBuffer(min int) *bytes.Buffer {
	if min > MessageBufferLength {
		return bytes.NewBuffer(make([]byte, min))
	}
	if v := writerBufferPool.Get(); v != nil {
		return v.(*bytes.Buffer)
	}
	return bytes.NewBuffer(make([]byte, MessageBufferLength))
}

func putWriterBuffer(b *bytes.Buffer) {
	b.Reset()
	if cap(b.Bytes()) == MessageBufferLength {
		writerBufferPool.Put(b)
	}
}

// Serialize returns the serialized bytes of the Message.
func (m *Message) Serialize() ([]byte, error) {
	b := make([]byte, m.Len())
	if err := m.SerializeTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// SerializeTo writes the serialized bytes of the Message into b.
func (m *Message) SerializeTo(b []byte) (err error) {
	m.Header.SerializeTo(b[0:HeaderLength])
	offset := HeaderLength
	for _, avp := range m.AVP {
		if err = avp.SerializeTo(b[offset:]); err != nil {
			return err
		}
		offset += avp.Len()
	}
	return nil
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
//	avp, err := m.FindAVP(avp.OriginHost)
//	avp, err := m.FindAVP("Origin-Host")
//
func (m *Message) FindAVP(code interface{}) (*AVP, error) {
	dictAVP, err := m.Dictionary().FindAVP(m.Header.ApplicationID, code)
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
		m.Header.ApplicationID,
		m.Header.HopByHopID,
		m.Header.EndToEndID,
		m.Dictionary(),
	)
	nm.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(resultCode))
	return nm
}

func (m *Message) String() string {
	var b bytes.Buffer
	var typ string
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		typ = "Request"
	} else {
		typ = "Answer"
	}
	if dictCMD, err := m.Dictionary().FindCommand(
		m.Header.ApplicationID,
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
		if dictAVP, err := m.Dictionary().FindAVP(
			m.Header.ApplicationID,
			a.Code,
		); err != nil {
			fmt.Fprintf(&b, "\tUnknown %s (%s)\n", a, err)
		} else if a.Data.Type() == GroupedAVPType {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, printGrouped("\t", m, a, 1))
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, a)
		}
	}
	return b.String()
}

func printGrouped(prefix string, m *Message, a *AVP, indent int) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:Grouped{\n",
		a.Code,
		a.Flags,
		a.Len(),
		a.VendorID,
	)
	for _, ga := range a.Data.(*GroupedAVP).AVP {
		if dictAVP, err := m.Dictionary().FindAVP(
			m.Header.ApplicationID,
			ga.Code,
		); err != nil {
			fmt.Fprintf(&b, "%s\tUnknown %s (%s),\n", prefix, ga, err)
		} else {
			if ga.Data.Type() == GroupedAVPType {
				indent++
				tabs := indentTabs(indent)
				fmt.Fprintf(&b, "%s%s %s\n", tabs, dictAVP.Name, printGrouped(tabs, m, ga, indent))
			} else {
				fmt.Fprintf(&b, "%s\t%s %s,\n", prefix, dictAVP.Name, ga)
			}
		}
	}
	fmt.Fprintf(&b, "%s}}", prefix)
	return b.String()
}

func indentTabs(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += "\t"
	}
	return s
}
