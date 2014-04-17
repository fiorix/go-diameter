// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Float32 Diameter Type
type Float32 float32

func DecodeFloat32(b []byte) (DataType, error) {
	var n float32
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &n)
	return Float32(n), err
}

func (n Float32) Serialize() []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, n)
	return b.Bytes()
}

func (n Float32) Len() int {
	return 4
}

func (n Float32) Padding() int {
	return 0
}

func (n Float32) Type() DataTypeId {
	return Float32Type
}

func (n Float32) String() string {
	return fmt.Sprintf("Float32{%0.4f}", n)
}
