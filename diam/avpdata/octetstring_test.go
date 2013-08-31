// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOctetString(t *testing.T) {
	s := "hello, world!"
	b := []byte{
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
		// Padding 0x00, 0x00, 0x00,
	}
	os := &OctetString{Value: s}
	//os := &DiameterIdentity{OctetString{Value: s}}
	osb := os.Bytes()
	if os.Padding > 0 {
		osb = osb[:uint32(len(osb))-os.Padding]
	}
	if !bytes.Equal(osb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", osb, b))
		return
	}
	os.Put(b)
	if d := os.Data(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
	if os.Padding != 3 {
		t.Error(fmt.Errorf("Padding length of '%s' is %d, expected 3.",
			s, os.Padding))
		return
	}
}

func BenchmarkOctetStringParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(OctetString)
		p.Put([]byte{
			0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20,
			0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21,
		})
	}
}

func BenchmarkOctetStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &OctetString{Value: "Hello, world"}
		p.Bytes()
	}
}
