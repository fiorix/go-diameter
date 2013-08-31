// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avpdata

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	s := int64(1377093974)
	b := []byte{0x52, 0x14, 0xc9, 0x56}
	tm := Time{Value: time.Unix(s, 0)}
	tmb := tm.Bytes()
	if !bytes.Equal(tmb, b) {
		t.Error(fmt.Errorf("Bytes are 0x%x, expected 0x%x", tmb, b))
		return
	}
	tm.Put(b)
	if d := tm.Data().(time.Time).Unix(); d != s {
		t.Error(fmt.Errorf("Data is 0x%x, expected 0x%x", d, s))
		return
	}
}

func BenchmarkTimeParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := new(Time)
		p.Put([]byte{0x52, 0x14, 0xc9, 0x56})
	}
}

func BenchmarkTimeBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := &Time{Value: time.Unix(1377093974, 0)}
		p.Bytes()
	}
}
