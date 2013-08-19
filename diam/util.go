// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
)

// pad4 returns n padded to 4 bytes
func pad4(n uint32) uint32 {
	return n + ((4 - n) & 3)
}

// uint24to32 converts b from [3]uint8 to uint32.
func uint24to32(b [3]uint8) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// uint32to24 converts b from uint32 to [3]uint8.
func uint32to24(b uint32) [3]uint8 {
	var r [3]uint8
	r[0] = uint8(b >> 16)
	r[1] = uint8(b >> 8)
	r[2] = uint8(b)
	return r
}

// bytes2uint32 converts byte array to uint32.
func bytes2uint32(b []byte) uint32 {
	if len(b) == 4 {
		return uint32(b[0])<<24 |
			uint32(b[1])<<16 |
			uint32(b[2])<<8 |
			uint32(b[3])
	}
	return 0
}

// uint32tobytes converts uint32 to byte array.
func uint32tobytes(n uint32) []byte {
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, n) // ign error?
	return b.Bytes()
}

// bytes2uint64 converts 8 bytes uint64.
func bytes2uint64(b []byte) uint64 {
	if len(b) == 4 {
		return uint64(b[0])<<56 |
			uint64(b[1])<<48 |
			uint64(b[2])<<40 |
			uint64(b[3])<<32 |
			uint64(b[4])<<24 |
			uint64(b[5])<<16 |
			uint64(b[6])<<8 |
			uint64(b[7])
	}
	return 0
}

// uint64tobytes converts uint64 to byte array.
func uint64tobytes(n uint64) []byte {
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, n) // ign error?
	return b.Bytes()
}
