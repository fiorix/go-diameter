// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

// Enumerated Diameter Type
type Enumerated Integer32

func DecodeEnumerated(b []byte) (DataType, error) {
	v, err := DecodeInteger32(b)
	return Enumerated(v.(Integer32)), err
}

func (n Enumerated) Serialize() []byte {
	return Integer32(n).Serialize()
}

func (n Enumerated) Len() int {
	return 4
}

func (n Enumerated) Padding() int {
	return 0
}

func (n Enumerated) Type() DataTypeId {
	return EnumeratedType
}

func (n Enumerated) String() string {
	return fmt.Sprintf("Enumerated{%d}", n)
}
