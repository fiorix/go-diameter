// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

// UTF8String Diameter Type.
type UTF8String OctetString

func DecodeUTF8String(b []byte) (DataType, error) {
	return UTF8String(OctetString(b)), nil
}

func (s UTF8String) Serialize() []byte {
	return OctetString(s).Serialize()
}

func (s UTF8String) Len() int {
	return len(s)
}

func (s UTF8String) Padding() int {
	l := len(s)
	return pad4(l) - l
}

func (s UTF8String) Type() DataTypeId {
	return UTF8StringType
}

func (s UTF8String) String() string {
	return fmt.Sprintf("UTF8String{%s},Padding:%d", string(s), s.Padding())
}
