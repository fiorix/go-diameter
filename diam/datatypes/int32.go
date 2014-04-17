// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Integer32 Diameter Type
type Integer32 int32

func DecodeInteger32(b []byte) (DataType, error) {
	var n int32
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &n)
	return Integer32(n), err
}

func (n Integer32) Serialize() []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, n)
	return b.Bytes()
}

func (n Integer32) Len() int {
	return 4
}

func (n Integer32) Padding() int {
	return 0
}

func (n Integer32) String() string {
	return fmt.Sprintf("Integer32{%d}", n)
}
