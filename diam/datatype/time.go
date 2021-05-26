// Copyright 2013-2015 go-diameter authors. All rights reserved.
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

const rfc868offset = 2208988800 // Diff. between 1970 and 1900 in seconds.
//UTC time is reckoned from 6h 28m 16s UTC on 7 February 2036 because overload happens
const rfc2030offset = 2085978496 // 2085978496 comes from FFFFFFFF – 2208988800

// DecodeTime decodes a Time data type from byte array.
func DecodeTime(b []byte) (Type, error) {
	if len(b) != 4 {
		return &Time{}, nil
	}
	if (b[0] >> 7) == 0 {
		return Time(time.Unix(int64(binary.BigEndian.Uint32(b))+rfc2030offset, 0)), nil
	} else {
		return Time(time.Unix(int64(binary.BigEndian.Uint32(b))-rfc868offset, 0)), nil
	}

}

// Serialize implements the Type interface.
func (t Time) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(time.Time(t).Unix())+rfc868offset)
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
