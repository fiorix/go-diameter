// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import "fmt"

// DiameterURI Diameter Type.
type DiameterURI OctetString

func DecodeDiameterURI(b []byte) (DiameterURI, error) {
	return DiameterURI(OctetString(b)), nil
}

func (s DiameterURI) Serialize() []byte {
	return OctetString(s).Serialize()
}

func (s DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{%s}", string(s))
}
