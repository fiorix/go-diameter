// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import "fmt"

// DiameterURI Diameter Type.
type DiameterURI struct {
	Value string
}

// Data implements the Generic interface.
func (du *DiameterURI) Data() Generic {
	return du.Value
}

// Put implements the Codec interface.
func (du *DiameterURI) Put(b []byte) {
	du.Value = string(b)
}

// Bytes implement the Codec interface.
func (du *DiameterURI) Bytes() []byte {
	return []byte(du.Value)
}

// Length implements the Codec interface.
func (du *DiameterURI) Length() uint32 {
	return uint32(len(du.Value))
}

// String returns a human readable version of the AVP.
func (du *DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{Value:'%s'}", du.Value)
}
