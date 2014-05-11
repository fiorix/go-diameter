// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

// OctetString Diameter Type.
type OctetString string

func DecodeOctetString(b []byte) (DataType, error) {
	return OctetString(b), nil
}

func (s OctetString) Serialize() []byte {
	return []byte(s)
}

func (s OctetString) Len() int {
	return len(s)
}

func (s OctetString) Padding() int {
	l := len(s)
	return pad4(l) - l
}

func (s OctetString) Type() DataTypeId {
	return OctetStringType
}

func (s OctetString) String() string {
	return fmt.Sprintf("OctetString{%s},Padding:%d", string(s), s.Padding())
}
