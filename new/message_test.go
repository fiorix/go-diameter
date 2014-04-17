// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"sync"
	"testing"

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

func BenchmarkReadMessage(b *testing.B) {
	reader := testMessage()
	for n := 0; n < b.N; n++ {
		ReadMessage(reader, dict.Default)
		reader.Seek(0, 0)
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
