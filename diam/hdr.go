// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter header.  Part of go-diameter.
// http://tools.ietf.org/html/rfc6733#section-3

package diam

import (
	"unsafe"

	"github.com/fiorix/go-diameter/diam/util"
)

type Header struct {
	Version          uint8
	RawMessageLength [3]uint8
	CommandFlags     uint8
	RawCommandCode   [3]uint8
	ApplicationId    uint32
	HopByHopId       uint32
	EndToEndId       uint32
}

const hdrSize = uint32(unsafe.Sizeof(Header{}))

// MessageLength helper function returns RawMessageLength as int.
func (hdr *Header) MessageLength() uint32 {
	return util.Uint24To32(hdr.RawMessageLength)
}

// UpdateLength updates RawMessageLength from an int.
func (hdr *Header) SetMessageLength(length uint32) {
	hdr.RawMessageLength = util.Uint32To24(hdrSize + length)
}

// CommandCode returns RawCommandCode as int.
func (hdr *Header) CommandCode() uint32 {
	return util.Uint24To32(hdr.RawCommandCode)
}
