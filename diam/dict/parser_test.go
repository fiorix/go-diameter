// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dict

import (
	"os"
	"testing"
)

const testDict = "./testdata/base.xml"

func TestNewParser(t *testing.T) {
	p, err := NewParser(testDict)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}

func TestLoadFile(t *testing.T) {
	p, _ := NewParser()
	if err := p.LoadFile(testDict); err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	f, err := os.Open(testDict)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p, _ := NewParser()
	if err = p.Load(f); err != nil {
		t.Fatal(err)
	}
}
