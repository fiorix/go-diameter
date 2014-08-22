// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import (
	"encoding/binary"
	"fmt"
)

// Integer32 Diameter Type
type Integer32 int32

func DecodeInteger32(b []byte) (DataType, error) {
	return Integer32(binary.BigEndian.Uint32(b)), nil
}

func (n Integer32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func (n Integer32) Len() int {
	return 4
}

func (n Integer32) Padding() int {
	return 0
}

func (n Integer32) Type() DataTypeId {
	return Integer32Type
}

func (n Integer32) String() string {
	return fmt.Sprintf("Integer32{%d}", n)
}
