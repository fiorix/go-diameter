// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatypes

import (
	"fmt"
	"net"
)

// IPv4 Diameter Type for Framed-IP-Address AVP.
type IPv4 net.IP

func DecodeIPv4(b []byte) (DataType, error) {
	return IPv4(b), nil
}

func (ip IPv4) Serialize() []byte {
	if ip4 := net.IP(ip).To4(); ip4 != nil {
		return ip4
	}
	return ip
}

func (ip IPv4) Len() int {
	return 4
}

func (ip IPv4) Padding() int {
	return 0
}

func (ip IPv4) Type() DataTypeId {
	return IPv4Type
}

func (ip IPv4) String() string {
	return fmt.Sprintf("IPv4{%s}", net.IP(ip))
}
