// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter header parser.  Part of go-diameter.

package base

import (
	"encoding/binary"
	"fmt"
	"io"
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
