// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import "fmt"

// Grouped Diameter Type
type Grouped []byte

func DecodeGrouped(b []byte) (DataType, error) {
	return Grouped(b), nil
}

func (b Grouped) Serialize() []byte {
	return b
}

func (b Grouped) String() string {
	return fmt.Sprint("Grouped{...}")
}
