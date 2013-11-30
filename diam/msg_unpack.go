// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fiorix/go-diameter/diam/dict"
)

// ReadMessage reads a diameter message from the connection.
//
// Dictionary d is used to parse message's AVPs. If set to nil, the static
// dict.Default is used.
//
// Using the wrong dictionary might result in errors.
func ReadMessage(r io.Reader, d *dict.Parser) (*Message, error) {
	var (
		err     error
		avp     *AVP
		padding uint32
	)
	if d == nil {
		d = dict.Default
	}
	m := &Message{Dict: d}
	if m.Header, err = readHeader(r); err != nil {
		return nil, err
	}
	// n is how many bytes are left in this message.
	n := int32(m.Header.MessageLength() - hdrSize)
	// Read all AVPs in this message.
	for n != 0 {
		if padding, avp, err = ReadAVP(m, r); err != nil {
			return nil, err
		} else {
			n = n - int32(avp.Length+padding)
			if n < 0 {
				return nil, fmt.Errorf(
					"Malformed AVP: %s", avp.String())
			}
		}
		m.AVP = append(m.AVP, avp)
	}
	return m, nil
}

// readHeader reads one diameter header from the connection and return it.
func readHeader(r io.Reader) (*Header, error) {
	var hdr Header
	if err := binary.Read(r, binary.BigEndian, &hdr); err != nil {
		return nil, err
	}
	// Only supports diameter version 1.
	if hdr.Version != byte(1) {
		return nil, fmt.Errorf(
			"Unsupported diameter version %d", hdr.Version)
	}
	return &hdr, nil
}
