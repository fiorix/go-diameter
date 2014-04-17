// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Float64 Diameter Type
type Float64 float64

func DecodeFloat64(b []byte) (DataType, error) {
	var n float64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &n)
	return Float64(n), err
}

func (n Float64) Serialize() []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, n)
	return b.Bytes()
}

func (n Float64) Len() int {
	return 8
}

func (n Float64) Padding() int {
	return 0
}

func (n Float64) String() string {
	return fmt.Sprintf("Float64{%0.4f}", n)
}
