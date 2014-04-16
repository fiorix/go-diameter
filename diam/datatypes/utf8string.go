// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import "fmt"

// UTF8String Diameter Type.
type UTF8String OctetString

func DecodeUTF8String(b []byte) (UTF8String, error) {
	return UTF8String(OctetString(b)), nil
}

func (s UTF8String) Serialize() []byte {
	return OctetString(s).Serialize()
}

func (s UTF8String) String() string {
	return fmt.Sprintf("UTF8String{%s}", string(s))
}
