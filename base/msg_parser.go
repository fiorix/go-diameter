// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter message, header with multiple AVPs.  Part of go-diameter.

package base

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"

	"github.com/fiorix/go-diameter/dict"
)

// ReadHeader reads one diameter header from the connection and return it.
func ReadHeader(r io.Reader) (*Header, error) {
	hdr := new(Header)
	if err := binary.Read(r, binary.BigEndian, hdr); err != nil {
		return nil, err
	}
	// Only supports diameter version 1.
	if hdr.Version != byte(1) {
		return nil, fmt.Errorf(
			"Unsupported diameter version %d", hdr.Version)
	}
	return hdr, nil
}

// ReadMessage reads an entire diameter message from the connection with
// Header and AVPs and return it.
func ReadMessage(r io.Reader, d *dict.Parser) (*Message, error) {
	var (
		err   error
		avp   *AVP
		extra uint32
	)
	m := &Message{Dict: d}
	if m.Header, err = ReadHeader(r); err != nil {
		return nil, err
	}
	// b is how many bytes are left in this message.
	b := m.Header.MessageLength() - uint32(unsafe.Sizeof(*m.Header))
	// Read all AVPs in this message.
	for b != 0 {
		// Read may timeout some time.
		if extra, avp, err = ReadAVP(r, m); err != nil {
			return nil, err
		} else {
			b -= (avp.Length + extra)
			if b < 0 {
				return nil, fmt.Errorf("Malformed AVP %s", avp)
			}
		}
		m.AVP = append(m.AVP, avp)
	}
	return m, nil
}
