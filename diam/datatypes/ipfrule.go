// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import "fmt"

// IPFilterRule Diameter Type.
type IPFilterRule OctetString

func DecodeIPFilterRule(b []byte) (DataType, error) {
	return IPFilterRule(OctetString(b)), nil
}

func (s IPFilterRule) Serialize() []byte {
	return OctetString(s).Serialize()
}

func (s IPFilterRule) Len() int {
	return len(s)
}

func (s IPFilterRule) Padding() int {
	l := len(s)
	return pad4(l) - l
}

func (s IPFilterRule) Type() DataTypeId {
	return IPFilterRuleType
}

func (s IPFilterRule) String() string {
	return fmt.Sprintf("IPFilterRule{%s},Padding:%d", string(s), s.Padding())
}
