// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Grouped AVP type.  Part of go-diameter.

package base

import (
	"bytes"
	"fmt"
	"strings"
)

// Grouped Diameter Type
type Grouped struct {
	Message *Message // Parent Message of this group.
	AVP     []*AVP   // Group AVPs
	Buffer  []byte   // len(Buffer) might be bigger than length below due to padding.
	length  uint32   // Aggregate length of all AVPs without padding.
}

// Data implements the Data interface.
func (gr *Grouped) Data() Data {
	return gr.AVP
}

// Put implements the Codec interface. It updates internal Buffer and Length.
// Takes an AVP as input.
func (gr *Grouped) Put(d Data) {
	avp := d.(*AVP)
	gr.length += avp.Length
	if avp.dict != gr.Message.Dict {
		avp.dict = gr.Message.Dict
	}
	gr.AVP = append(gr.AVP, avp)
	gr.Buffer = bytes.Join([][]byte{gr.Buffer, avp.Bytes()}, []byte{})
}

// Bytes implement the Codec interface. Bytes are always returned from
// internal Buffer cache.
func (gr *Grouped) Bytes() []byte {
	return gr.Buffer
}

// Length implements the Codec interface.
func (gr *Grouped) Length() uint32 {
	return gr.length
}

// String returns a human readable version of the AVP.
func (gr *Grouped) String() string {
	s := make([]string, len(gr.AVP))
	for n, avp := range gr.AVP {
		s[n] = avp.String()
	}
	return fmt.Sprintf("Grouped{%s}", strings.Join(s, ","))
}

// NewAVP allocates an AVP and appends to the group.
// @code can be either the AVP code (int, uint32) or name (string).
func (gr *Grouped) NewAVP(code interface{}, flags uint8, vendor uint32, body Codec) (*AVP, error) {
	if gr.Message == nil {
		panic("Grouped.NewAVP requires a parent Message, not nil")
	}
	avp, err := newAVP(
		gr.Message.Header.ApplicationId,
		gr.Message.Dict,
		code,
		flags,
		vendor,
		body,
	)
	if err != nil {
		return nil, err
	}
	gr.Put(avp)
	return avp, nil
}

// Append is an alias to Put.
func (gr *Grouped) Append(avp *AVP) {
	gr.Put(avp)
}

// NewGroup allocates a new Grouped AVP. Same as &Grouped{Message: m}
func NewGroup(m *Message) *Grouped {
	return &Grouped{Message: m}
}
