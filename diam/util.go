// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// pad4 returns n padded to 4 bytes
func pad4(n uint32) uint32 {
	return n + ((4 - n) & 3)
}

// uint24To32 converts b from [3]uint8 to uint32 in network byte order.
func uint24To32(b [3]uint8) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// uint32To24 converts b from uint32 to [3]uint8 in network byte order.
func uint32To24(b uint32) [3]uint8 {
	var r [3]uint8
	r[0] = uint8(b >> 16)
	r[1] = uint8(b >> 8)
	r[2] = uint8(b)
	return r
}

var invalid = errors.New("Invalid type for conversion")

// nToNetBytes converts any numeric type to byte array in network byte order.
func nToNetBytes(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	switch v.(type) {
	case int32:
		binary.Write(buf, binary.BigEndian, v.(int32))
	case int64:
		binary.Write(buf, binary.BigEndian, v.(int64))
	case uint32:
		binary.Write(buf, binary.BigEndian, v.(uint32))
	case uint64:
		binary.Write(buf, binary.BigEndian, v.(uint64))
	case float32:
		binary.Write(buf, binary.BigEndian, v.(float32))
	case float64:
		binary.Write(buf, binary.BigEndian, v.(float64))
	default:
		return nil, invalid
	}
	return buf.Bytes(), nil
}

// netBytesToN converts byte array in network byte order to any numeric type.
func netBytesToN(b []byte, v interface{}) error {
	switch v.(type) {
	case *int32:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*int32))
	case *int64:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*int64))
	case *uint32:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*uint32))
	case *uint64:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*uint64))
	case *float32:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*float32))
	case *float64:
		binary.Read(bytes.NewBuffer(b), binary.BigEndian, v.(*float64))
	default:
		return invalid
	}
	return nil
}
