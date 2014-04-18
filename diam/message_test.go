// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net"
	"sync"
	"testing"

	"github.com/fiorix/go-diameter/diam/datatypes"
	"github.com/fiorix/go-diameter/diam/dict"
)

var testOnce sync.Once

func TestReadMessage(t *testing.T) {
	msg, err := ReadMessage(testMessage(), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg)
}

func TestNewMessage(t *testing.T) {
	want, _ := ReadMessage(testMessage(), dict.Default)
	m := NewMessage(257, 0x80, 1, 0x2c0b6149, 0xdbbfd385, dict.Default)
	m.NewAVP(264, 0x40, 0, datatypes.DiameterIdentity("client"))
	m.NewAVP(296, 0x40, 0, datatypes.DiameterIdentity("localhost"))
	m.NewAVP(257, 0x40, 0, datatypes.Address(net.ParseIP("192.168.242.122")))
	m.NewAVP(266, 0x40, 0, datatypes.Unsigned32(13))
	m.NewAVP(269, 0x40, 0, datatypes.UTF8String("go-diameter"))
	m.NewAVP(278, 0x40, 0, datatypes.Unsigned32(3896392580))
	if m.Len() != want.Len() {
		t.Fatalf("Unexpected message length.\nWant: %d\n%s\nHave: %d\n%s",
			want.Len(), want, m.Len(), m)
	}
	a, b := m.Serialize(), want.Serialize()
	if !bytes.Equal(a, b) {
		t.Fatalf("Unexpected message.\nWant:\n%s\n%s\nHave:\n%s\n%s",
			want, hex.Dump(b), m, hex.Dump(a))
	}
	t.Log(len(a), "bytes\n", m)
	t.Log(hex.Dump(a))
}

func TestMessageFindAVP(t *testing.T) {
	m, _ := ReadMessage(testMessage(), dict.Default)
	a, err := m.FindAVP(278)
	if err != nil {
		t.Fatal(err)
	}
	a, err = m.FindAVP("Origin-State-Id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)
}

func BenchmarkReadMessage(b *testing.B) {
	reader := testMessage()
	for n := 0; n < b.N; n++ {
		ReadMessage(reader, dict.Default)
		reader.Seek(0, 0)
	}
}

func BenchmarkWriteMessage(b *testing.B) {
	m, _ := ReadMessage(testMessage(), dict.Default)
	for n := 0; n < b.N; n++ {
		m.WriteTo(ioutil.Discard)
	}
}

func testMessage() *bytes.Reader {
	// build message with fragments from other tests
	buf := bytes.NewBuffer(testHeader)
	for _, b := range testAVP {
		buf.Write(b)
	}
	return bytes.NewReader(buf.Bytes())
}
