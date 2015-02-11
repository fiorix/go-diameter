// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Time data type.
type Time time.Time

// DecodeTime decodes a Time data type from byte array.
func DecodeTime(b []byte) (Type, error) {
	return Time(time.Unix(int64(binary.BigEndian.Uint32(b)), 0)), nil
}

// Serialize implements the Type interface.
func (t Time) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(time.Time(t).Unix()))
	return b
}

// Len implements the Type interface.
func (t Time) Len() int {
	return 4
}

// Padding implements the Type interface.
func (t Time) Padding() int {
	return 0
}

// Type implements the Type interface.
func (t Time) Type() TypeID {
	return TimeType
}

// String implements the Type interface.
func (t Time) String() string {
	return fmt.Sprintf("Time{%s}", time.Time(t))
}
