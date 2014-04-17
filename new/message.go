// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fiorix/go-diameter/diam/dict"
)

type Message struct {
	Header *Header
	AVP    []*AVP

	Dictionary *dict.Parser // Used to encode and decode AVPs.
}

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
