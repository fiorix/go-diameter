// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package format

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Float32 Diameter Format.
type Float32 float32

func DecodeFloat32(b []byte) (Format, error) {
	return Float32(math.Float32frombits(binary.BigEndian.Uint32(b))), nil
}

func (n Float32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(n)))
	return b
}

func (n Float32) Len() int {
	return 4
}

func (n Float32) Padding() int {
	return 0
}

func (n Float32) Format() FormatId {
	return Float32Format
}

func (n Float32) String() string {
	return fmt.Sprintf("Float32{%0.4f}", n)
}
