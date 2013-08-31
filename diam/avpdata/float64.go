// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Float64 Diameter Type
type Float64 struct {
	Value float64
}

// Data implements the Generic interface.
func (n *Float64) Data() Generic {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Float64.
func (n *Float64) Put(b []byte) {
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float64 and stored on Buffer.
func (n *Float64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	return b.Bytes()
}

// Length implements the Codec interface.
func (n *Float64) Length() uint32 {
	return 8
}

// String returns a human readable version of the AVP.
func (n *Float64) String() string {
	return fmt.Sprintf("Float64{Value:%d}", n.Value)
}
