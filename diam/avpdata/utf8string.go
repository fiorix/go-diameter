// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import "fmt"

// UTF8String Diameter Type.
type UTF8String struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *UTF8String) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf(
		"UTF8String{Value:'%s',Padding:%d}", p.Value, p.Padding)
}
