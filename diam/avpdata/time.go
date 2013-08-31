// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// Time Diameter Type.
type Time struct {
	Value time.Time
}

// Data implements the Data interface.
func (t *Time) Data() Generic {
	return t.Value
}

// Put implements the Codec interface. It updates internal Value.
func (t *Time) Put(b []byte) {
	if len(b) == 4 {
		t.Value = time.Unix(int64(binary.BigEndian.Uint32(b)), 0)
	}
}

// Bytes implement the Codec interface.
func (t *Time) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(t.Value.Unix()))
	return buf.Bytes()
}

// Length implements the Codec interface. Returns length without padding.
func (t *Time) Length() uint32 {
	return 4
}

// String returns a human readable version of the AVP.
func (t *Time) String() string {
	return fmt.Sprintf("Time{Value:'%s'}", t.Value.String())
}
