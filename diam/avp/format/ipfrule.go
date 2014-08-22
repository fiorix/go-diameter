// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package format

import "fmt"

// IPFilterRule Diameter Format.
type IPFilterRule OctetString

func DecodeIPFilterRule(b []byte) (Format, error) {
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

func (s IPFilterRule) Format() FormatId {
	return IPFilterRuleFormat
}

func (s IPFilterRule) String() string {
	return fmt.Sprintf("IPFilterRule{%s},Padding:%d", string(s), s.Padding())
}
