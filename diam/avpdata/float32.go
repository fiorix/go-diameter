// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Float32 Diameter Type
type Float32 struct {
	Value float32
}

// Data implements the Data interface.
func (n *Float32) Data() Generic {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Float32.
func (n *Float32) Put(b []byte) {
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float32 and stored on Buffer.
func (n *Float32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	return b.Bytes()
}

// Length implements the Codec interface.
func (n *Float32) Length() uint32 {
	return 4
}

// String returns a human readable version of the AVP.
func (n *Float32) String() string {
	return fmt.Sprintf("Float32{Value:%d}", n.Value)
}
