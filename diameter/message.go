// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diameter

import (
	"io"
	"unsafe"
)

type Message struct {
	Header *Header
	AVP    []*AVP
}

// ReadMessage reads an entire diameter message with Header and AVPs.
func ReadMessage(r io.Reader, dict *Dict) (*Message, error) {
	var (
		err   error
		avp   *AVP
		extra uint32
	)
	msg := new(Message)
	if msg.Header, err = ReadHeader(r); err != nil {
		return nil, err
	}
	// b is how many bytes are left in this message.
	b := msg.Header.MessageLength() - uint32(unsafe.Sizeof(*msg.Header))
	// Read all AVPs in this message.
	for b > 1 {
		if extra, avp, err = ReadAVP(r, dict); err != nil {
			return nil, err
		} else {
			b -= (avp.Length + extra)
		}
		msg.AVP = append(msg.AVP, avp)
	}
	return msg, nil
}
