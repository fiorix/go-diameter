// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

// Grouped Diameter Type
type Grouped []byte

func DecodeGrouped(b []byte) (DataType, error) {
	return Grouped(b), nil
}

func (g Grouped) Serialize() []byte {
	return g
}

func (g Grouped) Len() int {
	return len(g)
}

func (g Grouped) Padding() int {
	return 0
}

func (g Grouped) Type() DataTypeId {
	return GroupedType
}

func (g Grouped) String() string {
	return fmt.Sprint("Grouped{...}")
}
