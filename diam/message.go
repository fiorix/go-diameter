// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"

	"github.com/fiorix/go-diameter/avp"
	"github.com/fiorix/go-diameter/dict"
)

// Message is a diameter message, composed of a header and multiple AVPs.
type Message struct {
	Header *Header
	AVP    []*avp.AVP
}

// ReadMessage reads an entire diameter message from the connection with
// Header and AVPs and return it.
func ReadMessage(r io.Reader, d *dict.Dict) (*Message, error) {
	var (
		err   error
		a     *avp.AVP
		extra uint32
	)
	msg := new(Message)
	if msg.Header, err = ReadHeader(r); err != nil {
		return nil, err
	}
	// b is how many bytes are left in this message.
	b := msg.Header.MessageLength() - uint32(unsafe.Sizeof(*msg.Header))
	// Read all AVPs in this message.
	for b != 0 {
		// Read may timeout some time.
		if extra, a, err = avp.ReadAVP(r, d); err != nil {
			return nil, err
		} else {
			b -= (a.Length + extra)
			if b < 0 {
				return nil, fmt.Errorf("Malformed AVP %s", a)
			}
		}
		msg.AVP = append(msg.AVP, a)
	}
	return msg, nil
}

// Find is a helper function that returns an AVP by looking up its code, or nil.
func (m *Message) Find(code uint32) *avp.AVP {
	for _, a := range m.AVP {
		if code == a.Code {
			return a
		}
	}
	return nil
}

// NewMessage allocates a new Message. Used for building messages that will
// be sent to a connection.
func NewMessage(code uint32, flags uint8, appid, hopbyhop, endtoend uint32) *Message {
	msg := new(Message)
	msg.Header = new(Header)
	msg.Header.Version = 1
	msg.Header.CommandFlags = flags
	msg.Header.RawCommandCode = uint32To24(code)
	msg.Header.HopByHopId = hopbyhop
	msg.Header.EndToEndId = endtoend
	return msg
}

// Add adds an AVP to the given Message.
func (m *Message) Add(a *avp.AVP) {
	m.AVP = append(m.AVP, a)
}

// Bytes returns the Message in binary form to be sent to a connection.
func (m *Message) Bytes() []byte {
	var buf [][]byte
	var length uint32
	for _, avp := range m.AVP {
		b := avp.Marshal()
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

// PrettyPrint prints messages in a human readable format.
func (m *Message) PrettyPrint() {
	fmt.Println(m.Header)
	for _, avp := range m.AVP {
		fmt.Printf("  %s\n", avp)
	}
	fmt.Println()
}
