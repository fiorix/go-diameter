// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"fmt"
	"testing"
)

// Tests

func TestPad4(t *testing.T) {
	switch {
	case pad4(3) != 4:
		t.Error(fmt.Errorf("pad(3) != 4"))
	case pad4(5) != 8:
		t.Error(fmt.Errorf("pad(5) != 8"))
	case pad4(10) != 12:
		t.Error(fmt.Errorf("pad(10) != 12"))
	}
}

func TestUintConversion(t *testing.T) {
	b := uint32To24(4096)
	n := uint24To32(b)
	if n != 4096 {
		t.Error(fmt.Errorf("uint24 0x%x != 0x%x uint32", b, n))
	}
}

const (
	testInt32   = int32((1 << 31) - 1)
	testInt64   = int64((1 << 63) - 1)
	testUint32  = uint32((1 << 32) - 1)
	testUint64  = uint64((1 << 64) - 1)
	testFloat32 = float32((1 << 31) - 1)
	testFloat64 = float64((1 << 63) - 1)
)

func TestEndianInt32(t *testing.T) {
	b, err := nToNetBytes(testInt32)
	if err != nil {
		t.Error(err)
		return
	}
	if len(b) != 4 {
		t.Error(fmt.Errorf("Length of 0x%x should be 4!", b))
		return
	}
	var v int32
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
		return
	}
	if v != testInt32 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testInt32))
		return
	}
}

func TestEndianInt64(t *testing.T) {
	b, err := nToNetBytes(testInt64)
	if err != nil {
		t.Error(err)
	}
	if len(b) != 8 {
		t.Error(fmt.Errorf("Length of 0x%x should be 8!", b))
		return
	}
	var v int64
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
	}
	if v != testInt64 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testInt64))
		return
	}
}

func TestEndianUint32(t *testing.T) {
	b, err := nToNetBytes(testUint32)
	if err != nil {
		t.Error(err)
		return
	}
	if len(b) != 4 {
		t.Error(fmt.Errorf("Length of 0x%x should be 4!", b))
		return
	}
	var v uint32
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
		return
	}
	if v != testUint32 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testUint32))
		return
	}
}

func TestEndianUint64(t *testing.T) {
	b, err := nToNetBytes(testUint64)
	if err != nil {
		t.Error(err)
	}
	if len(b) != 8 {
		t.Error(fmt.Errorf("Length of 0x%x should be 8!", b))
		return
	}
	var v uint64
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
	}
	if v != testUint64 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testUint64))
		return
	}
}

func TestEndianFloat32(t *testing.T) {
	b, err := nToNetBytes(testFloat32)
	if err != nil {
		t.Error(err)
		return
	}
	if len(b) != 4 {
		t.Error(fmt.Errorf("Length of 0x%x should be 4!", b))
		return
	}
	var v float32
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
		return
	}
	if v != testFloat32 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testFloat32))
		return
	}
}

func TestEndianFloat64(t *testing.T) {
	b, err := nToNetBytes(testFloat64)
	if err != nil {
		t.Error(err)
	}
	if len(b) != 8 {
		t.Error(fmt.Errorf("Length of 0x%x should be 8!", b))
		return
	}
	var v float64
	if err = netBytesToN(b, &v); err != nil {
		t.Error(err)
	}
	if v != testFloat64 {
		t.Error(fmt.Errorf("v 0x%x != 0x%x", v, testFloat64))
		return
	}
}
