// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Time Diameter Type.
type Time time.Time

func DecodeTime(b []byte) (DataType, error) {
	return Time(time.Unix(int64(binary.BigEndian.Uint32(b)), 0)), nil
}

func (t Time) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(time.Time(t).Unix()))
	return b
}

func (t Time) Len() int {
	return 4
}

func (t Time) Padding() int {
	return 0
}

func (t Time) Type() DataTypeId {
	return TimeType
}

func (t Time) String() string {
	return fmt.Sprintf("Time{%s}", time.Time(t))
}
