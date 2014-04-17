// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"encoding/binary"
	"fmt"
)

// Unsigned32 Diameter Type
type Unsigned32 uint32

func DecodeUnsigned32(b []byte) (DataType, error) {
	return Unsigned32(binary.BigEndian.Uint32(b)), nil
}

func (n Unsigned32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func (n Unsigned32) Len() int {
	return 4
}

func (n Unsigned32) Padding() int {
	return 0
}

func (n Unsigned32) Type() DataTypeId {
	return Unsigned32Type
}

func (n Unsigned32) String() string {
	return fmt.Sprintf("Unsigned32{%d}", n)
}
