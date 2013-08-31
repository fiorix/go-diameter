// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Unsigned32 Diameter Type
type Unsigned32 struct {
	Value uint32
}

// Data implements the Generic interface.
func (n *Unsigned32) Data() Generic {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Uint32.
func (n *Unsigned32) Put(b []byte) {
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Uint32 and stored on Buffer.
func (n *Unsigned32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	return b.Bytes()
}

// Length implements the Codec interface.
func (n *Unsigned32) Length() uint32 {
	return 4
}

// String returns a human readable version of the AVP.
func (n *Unsigned32) String() string {
	return fmt.Sprintf("Unsigned32{Value:%d}", n.Value)
}
