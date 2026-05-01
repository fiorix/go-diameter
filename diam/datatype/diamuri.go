// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"fmt"
	"strconv"
	"strings"
)

// DiameterURI data type.
type DiameterURI OctetString

// ParsedDiameterURI holds the components of a DiameterURI as defined in RFC 6733 §4.3.
type ParsedDiameterURI struct {
	Secure    bool   // true for "aaas://" scheme, false for "aaa://"
	FQDN      string // hostname
	Port      uint16 // default 3868 (aaa) or 5868 (aaas)
	Transport string // "tcp" or "sctp"; default "tcp"
	Protocol  string // "diameter", "radius", or "tacacsplus"; default "diameter"
}

// Parse parses the DiameterURI per RFC 6733 §4.3 ABNF:
//
//	"aaa://" FQDN [ ":" port ] [ ";transport=" transport ] [ ";protocol=" protocol ]
//	"aaas://" FQDN [ ":" port ] [ ";transport=" transport ] [ ";protocol=" protocol ]
func (s DiameterURI) Parse() (*ParsedDiameterURI, error) {
	raw := string(s)
	p := &ParsedDiameterURI{Transport: "tcp", Protocol: "diameter"}
	switch {
	case len(raw) > 7 && strings.EqualFold(raw[:7], "aaas://"):
		p.Secure = true
		p.Port = 5868
		raw = raw[7:]
	case len(raw) > 6 && strings.EqualFold(raw[:6], "aaa://"):
		p.Port = 3868
		raw = raw[6:]
	default:
		return nil, fmt.Errorf("diamuri: invalid scheme in %q", string(s))
	}
	// Extract ";protocol=" (must come after ";transport=" per ABNF)
	if i := strings.Index(raw, ";protocol="); i >= 0 {
		p.Protocol = raw[i+10:]
		raw = raw[:i]
	}
	// Extract ";transport="
	if i := strings.Index(raw, ";transport="); i >= 0 {
		p.Transport = raw[i+11:]
		raw = raw[:i]
	}
	// Remaining is FQDN [ ":" port ]
	if len(raw) == 0 {
		return nil, fmt.Errorf("diamuri: missing FQDN in %q", string(s))
	}
	if raw[0] == '[' {
		// IPv6 literal: [addr]:port
		if end := strings.Index(raw, "]"); end >= 0 {
			p.FQDN = raw[1:end]
			raw = raw[end+1:]
		}
	} else if i := strings.LastIndexByte(raw, ':'); i >= 0 {
		p.FQDN = raw[:i]
		portStr := raw[i+1:]
		if portStr == "" {
			return nil, fmt.Errorf("diamuri: empty port in %q", string(s))
		}
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("diamuri: invalid port in %q", string(s))
		}
		p.Port = uint16(port)
		raw = ""
	}
	if p.FQDN == "" {
		p.FQDN = raw
	}
	if p.FQDN == "" {
		return nil, fmt.Errorf("diamuri: missing FQDN in %q", string(s))
	}
	return p, nil
}

// DecodeDiameterURI decodes a DiameterURI from byte array.
func DecodeDiameterURI(b []byte) (Type, error) {
	d := make([]byte, len(b))
	copy(d, b)
	return DiameterURI(OctetString(d)), nil
}

// Serialize implements the Type interface.
func (s DiameterURI) Serialize() []byte {
	return OctetString(s).Serialize()
}

// Len implements the Type interface.
func (s DiameterURI) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s DiameterURI) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s DiameterURI) Type() TypeID {
	return DiameterURIType
}

// String implements the Type interface.
func (s DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{%s},Padding:%d", string(s), s.Padding())
}
