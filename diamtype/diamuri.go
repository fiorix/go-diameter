// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtype

import "fmt"

// DiameterURI Diameter Type.
type DiameterURI OctetString

func DecodeDiameterURI(b []byte) (DataType, error) {
	return DiameterURI(OctetString(b)), nil
}

func (s DiameterURI) Serialize() []byte {
	return OctetString(s).Serialize()
}

func (s DiameterURI) Len() int {
	return len(s)
}

func (s DiameterURI) Padding() int {
	l := len(s)
	return pad4(l) - l
}

func (s DiameterURI) Type() DataTypeId {
	return DiameterURIType
}

func (s DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{%s},Padding:%d", string(s), s.Padding())
}
