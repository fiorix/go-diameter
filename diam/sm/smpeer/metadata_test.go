// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpeer

import (
	"context"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
)

func TestFromCER(t *testing.T) {
	cer := &smparser.CER{
		OriginHost:  datatype.DiameterIdentity("foobar"),
		OriginRealm: datatype.DiameterIdentity("test"),
	}
	meta := FromCER(cer)
	if meta.OriginHost != cer.OriginHost {
		t.Fatalf("Unexpected OriginHost. Want %q, have %q",
			cer.OriginHost, meta.OriginHost)
	}
	if meta.OriginRealm != cer.OriginRealm {
		t.Fatalf("Unexpected OriginRealm. Want %q, have %q",
			cer.OriginRealm, meta.OriginRealm)
	}
	if meta.CER != cer {
		t.Fatalf("Unexpected CER. Want %p, have %p", cer, meta.CER)
	}
	if meta.CEA != nil {
		t.Fatalf("Unexpected CEA. Want nil, have %#v", meta.CEA)
	}
	ctx := NewContext(context.Background(), meta)
	data, ok := FromContext(ctx)
	if !ok {
		t.Fatal("Metadata not present in this context")
	}
	if data != meta {
		t.Fatalf("Unexpected Metadata. Want %#v, have %#v", meta, data)
	}
}

func TestFromCEA(t *testing.T) {
	cer := &smparser.CEA{
		OriginHost:  datatype.DiameterIdentity("foobar"),
		OriginRealm: datatype.DiameterIdentity("test"),
	}
	meta := FromCEA(cer)
	if meta.OriginHost != cer.OriginHost {
		t.Fatalf("Unexpected OriginHost. Want %q, have %q",
			cer.OriginHost, meta.OriginHost)
	}
	if meta.OriginRealm != cer.OriginRealm {
		t.Fatalf("Unexpected OriginRealm. Want %q, have %q",
			cer.OriginRealm, meta.OriginRealm)
	}
	if meta.CEA != cer {
		t.Fatalf("Unexpected CEA. Want %p, have %p", cer, meta.CEA)
	}
	if meta.CER != nil {
		t.Fatalf("Unexpected CER. Want nil, have %#v", meta.CER)
	}
	ctx := NewContext(context.Background(), meta)
	data, ok := FromContext(ctx)
	if !ok {
		t.Fatal("Metadata not present in this context")
	}
	if data != meta {
		t.Fatalf("Unexpected Metadata. Want %#v, have %#v", meta, data)
	}
}
