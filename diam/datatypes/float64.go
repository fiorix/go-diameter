// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Float64 Diameter Type
type Float64 float64

func DecodeFloat64(b []byte) (DataType, error) {
	return Float64(math.Float64frombits(binary.BigEndian.Uint64(b))), nil
}

func (n Float64) Serialize() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(float64(n)))
	return b
}

func (n Float64) Len() int {
	return 4
}

func (n Float64) Padding() int {
	return 0
}

func (n Float64) Type() DataTypeId {
	return Float64Type
}

func (n Float64) String() string {
	return fmt.Sprintf("Float64{%0.4f}", n)
}
