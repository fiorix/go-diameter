// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Integer64 Diameter Type
type Integer64 int64

func DecodeInteger64(b []byte) (Integer64, error) {
	var n int64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &n)
	return Integer64(n), err
}

func (n Integer64) Serialize() []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, n)
	return b.Bytes()
}

func (n Integer64) String() string {
	return fmt.Sprintf("Integer64{%d}", n)
}
