// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Float32 data type.
type Float32 float32

// DecodeFloat32 decodes a Float32 data type from a byte array.
func DecodeFloat32(b []byte) (DataType, error) {
	return Float32(math.Float32frombits(binary.BigEndian.Uint32(b))), nil
}

// Serialize implements the DataType interface.
func (n Float32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(n)))
	return b
}

// Len implements the DataType interface.
func (n Float32) Len() int {
	return 4
}

// Padding implements the DataType interface.
func (n Float32) Padding() int {
	return 0
}

// Type implements the DataType interface.
func (n Float32) Type() TypeID {
	return Float32Type
}

// String implements the DataType interface.
func (n Float32) String() string {
	return fmt.Sprintf("Float32{%0.4f}", n)
}
