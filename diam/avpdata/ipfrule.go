// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import "fmt"

// IPFilterRule Diameter Type.
type IPFilterRule struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *IPFilterRule) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf("IPFilterRule{Value:'%s',Padding:%d}",
		p.Value, p.Padding)
}
