// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

func TestCEA_MissingResultCode(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	cea := new(CEA)
	_, err := cea.Parse(m)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingResultCode {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingOriginHost(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	cea := new(CEA)
	_, err := cea.Parse(m)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingOriginHost {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingOriginRealm(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	cea := new(CEA)
	_, err := cea.Parse(m)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingOriginRealm {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA_MissingOriginStateID(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	cea := new(CEA)
	_, err := cea.Parse(m)
	if err == nil {
		t.Fatal("Broken CEA was parsed with no errors")
	}
	if err != nil && err != ErrMissingOriginStateID {
		t.Fatal("Unexpected error:", err)
	}
}

func TestCEA(t *testing.T) {
	m := diam.NewMessage(diam.CapabilitiesExchange, 0, 0, 0, 0, nil)
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("foobar"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	cea := new(CEA)
	if _, err := cea.Parse(m); err != nil {
		t.Fatal(err)
	}
	if cea.ResultCode != diam.Success {
		t.Fatalf("Unexpected Result-Code. Want %d, have %d",
			diam.Success, cea.ResultCode)
	}
	if cea.OriginStateID != 1 {
		t.Fatalf("Unexpected Result-Code. Want 1, have %d", cea.OriginStateID)
	}
}
