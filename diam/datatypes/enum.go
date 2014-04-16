// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import "fmt"

// Enumerated Diameter Type
type Enumerated Integer32

func DecodeEnumerated(b []byte) (Enumerated, error) {
	v, err := DecodeInteger32(b)
	return Enumerated(v), err
}

func (n Enumerated) Serialize() []byte {
	return Integer32(n).Serialize()
}

func (n Enumerated) String() string {
	return fmt.Sprintf("Enumerated{%d}", n)
}
