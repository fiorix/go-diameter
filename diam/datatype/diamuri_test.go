// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"bytes"
	"testing"
)

func TestDiameterURI(t *testing.T) {
	s := DiameterURI("hello")
	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}
	if v := s.Serialize(); !bytes.Equal(v, b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, v)
	}
	if s.Len() != 5 {
		t.Fatalf("Unexpected len. Want 5, have %d", s.Len())
	}
	if s.Padding() != 3 {
		t.Fatalf("Unexpected padding. Want 3, have %d", s.Padding())
	}
	if s.Type() != DiameterURIType {
		t.Fatalf("Unexpected type. Want %d, have %d",
			DiameterURIType, s.Type())
	}
	if len(s.String()) == 0 {
		t.Fatalf("Unexpected empty string")
	}
}

func TestDecodeDiameterURI(t *testing.T) {
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
	}
	s, err := DecodeDiameterURI(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(s.(DiameterURI)), b) {
		t.Fatalf("Unexpected value. Want 0x%x, have 0x%x", b, s)
	}
	if s.Len() != 12 {
		t.Fatalf("Unexpected len. Want 12, have %d", s.Len())
	}
	if s.Padding() != 0 {
		t.Fatalf("Unexpected padding. Want 0, have %d", s.Padding())
	}
	if v := string(s.(DiameterURI)); v != "hello, world" {
		t.Fatalf("Unexpected string. Want 'hello, world', have %q", v)
	}
}

func TestDiameterURI_Parse(t *testing.T) {
	tests := []struct {
		input     string
		secure    bool
		fqdn      string
		port      uint16
		transport string
		protocol  string
		wantErr   bool
	}{
		{"aaa://host.example.com", false, "host.example.com", 3868, "tcp", "diameter", false},
		{"aaas://host.example.com", true, "host.example.com", 5868, "tcp", "diameter", false},
		{"aaa://host.example.com:6666", false, "host.example.com", 6666, "tcp", "diameter", false},
		{"aaas://host.example.com:5868;transport=tcp", true, "host.example.com", 5868, "tcp", "diameter", false},
		{"aaa://host.example.com;transport=sctp", false, "host.example.com", 3868, "sctp", "diameter", false},
		{"aaa://host.example.com:3868;transport=tcp;protocol=diameter", false, "host.example.com", 3868, "tcp", "diameter", false},
		{"aaa://host.example.com;protocol=radius", false, "host.example.com", 3868, "tcp", "radius", false},
		// Case-insensitive scheme (RFC 6733 §4.3)
		{"AAA://host.example.com", false, "host.example.com", 3868, "tcp", "diameter", false},
		{"AAAS://host.example.com", true, "host.example.com", 5868, "tcp", "diameter", false},
		{"Aaa://host.example.com", false, "host.example.com", 3868, "tcp", "diameter", false},
		// Error cases
		{"http://host.example.com", false, "", 0, "", "", true},
		{"", false, "", 0, "", "", true},
		{"aaa://", false, "", 0, "", "", true},
		{"aaa://host.example.com:", false, "", 0, "", "", true},      // empty port
		{"aaa://host.example.com:99999", false, "", 0, "", "", true}, // port overflow
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			uri := DiameterURI(tt.input)
			p, err := uri.Parse()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if p.Secure != tt.secure {
				t.Errorf("Secure: got %v, want %v", p.Secure, tt.secure)
			}
			if p.FQDN != tt.fqdn {
				t.Errorf("FQDN: got %q, want %q", p.FQDN, tt.fqdn)
			}
			if p.Port != tt.port {
				t.Errorf("Port: got %d, want %d", p.Port, tt.port)
			}
			if p.Transport != tt.transport {
				t.Errorf("Transport: got %q, want %q", p.Transport, tt.transport)
			}
			if p.Protocol != tt.protocol {
				t.Errorf("Protocol: got %q, want %q", p.Protocol, tt.protocol)
			}
		})
	}
}
