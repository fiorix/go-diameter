// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import "fmt"

// DiameterIdentity Diameter Type.
type DiameterIdentity struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *DiameterIdentity) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf(
		"DiameterIdentity{Value:'%s',Padding:%d}", p.Value, p.Padding)
}
