// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import "fmt"

// Enumerated Diameter Type
type Enumerated struct {
	Integer32
}

// String returns a human readable version of the AVP.
func (p *Enumerated) String() string {
	return fmt.Sprintf("Enumerated{Value:%d}", p.Value)
}
