// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"net"
	"testing"

	"github.com/fiorix/go-diameter/diam/diamdict"
	"github.com/fiorix/go-diameter/diam/diamtype"
)

func TestUnmarshalAVP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type Data struct {
		OriginHost1 AVP  `avp:"Origin-Host"`
		OriginHost2 *AVP `avp:"Origin-Host"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v := d.OriginHost1.Data.(diamtype.DiameterIdentity); v != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", v)
	}
	if v := d.OriginHost2.Data.(diamtype.DiameterIdentity); v != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", v)
	}
}

func TestUnmarshalString(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type Data struct {
		OriginHost string `avp:"Origin-Host"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if d.OriginHost != "test" {
		t.Fatalf("Unexpected value. Want test, have %s", d.OriginHost)
	}
}

func TestUnmarshalNetIP(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type Data struct {
		HostIP1 *AVP   `avp:"Host-IP-Address"`
		HostIP2 net.IP `avp:"Host-IP-Address"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v := d.HostIP1.Data.(diamtype.Address); net.IP(v).String() != "10.1.0.1" {
		t.Fatalf("Unexpected value. Want 10.1.0.1, have %s", v)
	}
	if v := d.HostIP2.String(); v != "10.1.0.1" {
		t.Fatalf("Unexpected value. Want 10.1.0.1, have %s", v)
	}
}

func TestUnmarshalInt(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type Data struct {
		VendorId1 *AVP `avp:"Vendor-Id"`
		VendorId2 int  `avp:"Vendor-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v := d.VendorId1.Data.(diamtype.Unsigned32); v != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", v)
	}
	if d.VendorId2 != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", d.VendorId2)
	}
}

func TestUnmarshalSlice(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type Data struct {
		Vendors1 []*AVP `avp:"Supported-Vendor-Id"`
		Vendors2 []int  `avp:"Supported-Vendor-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if len(d.Vendors1) != 2 {
		t.Fatalf("Unexpected value. Want 2, have %d", len(d.Vendors1))
	}
	if v := d.Vendors1[0].Data.(diamtype.Unsigned32); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := d.Vendors1[1].Data.(diamtype.Unsigned32); v != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", v)
	}
	if len(d.Vendors2) != 2 {
		t.Fatalf("Unexpected value. Want 2, have %d", len(d.Vendors2))
	}
	if d.Vendors2[0] != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.Vendors2[0])
	}
	if d.Vendors2[1] != 13 {
		t.Fatalf("Unexpected value. Want 13, have %d", d.Vendors2[1])
	}
}

func TestUnmarshalGrouped(t *testing.T) {
	m, _ := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	type VSA struct {
		AuthAppId1 AVP  `avp:"Auth-Application-Id"`
		AuthAppId2 *AVP `avp:"Auth-Application-Id"`
		AuthAppId3 int  `avp:"Auth-Application-Id"`
		VendorId1  AVP  `avp:"Vendor-Id"`
		VendorId2  *AVP `avp:"Vendor-Id"`
		VendorId3  int  `avp:"Vendor-Id"`
	}
	type Data struct {
		VSA1 AVP  `avp:"Vendor-Specific-Application-Id"`
		VSA2 *AVP `avp:"Vendor-Specific-Application-Id"`
		VSA3 VSA  `avp:"Vendor-Specific-Application-Id"`
		VSA4 *VSA `avp:"Vendor-Specific-Application-Id"`
		VSA5 struct {
			AuthAppId1 AVP  `avp:"Auth-Application-Id"`
			AuthAppId2 *AVP `avp:"Auth-Application-Id"`
			AuthAppId3 int  `avp:"Auth-Application-Id"`
			VendorId1  AVP  `avp:"Vendor-Id"`
			VendorId2  *AVP `avp:"Vendor-Id"`
			VendorId3  int  `avp:"Vendor-Id"`
		} `avp:"Vendor-Specific-Application-Id"`
	}
	var d Data
	if err := m.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	if v, ok := d.VSA1.Data.(*Grouped); !ok {
		t.Fatalf("Unexpected value. Want Grouped, have %v", d.VSA1)
	} else if len(v.AVP) != 2 { // There must be 2 AVPs in it.
		t.Fatalf("Unexpected value. Want 2, have %d", len(v.AVP))
	}
	if v, ok := d.VSA2.Data.(*Grouped); !ok {
		t.Fatalf("Unexpected value. Want Grouped, have %s", d.VSA2)
	} else if len(v.AVP) != 2 { // There must be 2 AVPs in it.
		t.Fatalf("Unexpected value. Want 2, have %d", len(v.AVP))
	}
	if v := int(d.VSA3.AuthAppId1.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA3.AuthAppId2.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA3.AuthAppId3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA3.AuthAppId3)
	}
	if v := int(d.VSA3.VendorId1.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA3.VendorId2.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA3.VendorId3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA3.VendorId3)
	}
	if v := int(d.VSA4.AuthAppId1.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA4.AuthAppId2.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA4.AuthAppId3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA4.AuthAppId3)
	}
	if v := int(d.VSA4.VendorId1.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA4.VendorId2.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA4.VendorId3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA4.VendorId3)
	}
	if v := int(d.VSA5.AuthAppId1.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if v := int(d.VSA5.AuthAppId2.Data.(diamtype.Unsigned32)); v != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", v)
	}
	if d.VSA5.AuthAppId3 != 4 {
		t.Fatalf("Unexpected value. Want 4, have %d", d.VSA5.AuthAppId3)
	}
	if v := int(d.VSA5.VendorId1.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if v := int(d.VSA5.VendorId2.Data.(diamtype.Unsigned32)); v != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", v)
	}
	if d.VSA5.VendorId3 != 10415 {
		t.Fatalf("Unexpected value. Want 10415, have %d", d.VSA5.VendorId3)
	}
}

func TestUnmarshalCER(t *testing.T) {
	msg, err := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	if err != nil {
		t.Fatal(err)
	}
	type CER struct {
		OriginHost  string `avp:"Origin-Host"`
		OriginRealm string `avp:"Origin-Realm"`
		HostIP      net.IP `avp:"Host-IP-Address"`
		VendorId    int    `avp:"Vendor-Id"`
		ProductName string `avp:"Product-Name"`
		StateId     int    `avp:"Origin-State-Id"`
		Vendors     []int  `avp:"Supported-Vendor-Id"`
		AuthAppId   int    `avp:"Auth-Application-Id"`
		InbandSecId int    `avp:"Inband-Security-Id"`
		VSA         struct {
			AuthAppId int `avp:"Auth-Application-Id"`
			VendorId  int `avp:"Vendor-Id"`
		} `avp:"Vendor-Specific-Application-Id"`
		Firmware int `avp:"Firmware-Revision"`
	}
	var d CER
	if err := msg.Unmarshal(&d); err != nil {
		t.Fatal(err)
	}
	switch {
	case d.OriginHost != "test":
		t.Fatalf("Unexpected Code. Want test, have %s", d.OriginHost)
	case d.OriginRealm != "localhost":
		t.Fatalf("Unexpected Code. Want localhost, have %s", d.OriginRealm)
	case d.HostIP.String() != "10.1.0.1":
		t.Fatalf("Unexpected Host-IP-Address. Want 10.1.0.1, have %s", d.HostIP)
	case d.VendorId != 13:
		t.Fatalf("Unexpected Host-Vendor-Id. Want 13, have %d", d.VendorId)
	case d.ProductName != "go-diameter":
		t.Fatalf("Unexpected Product-Name. Want go-diameter, have %s", d.ProductName)
	case d.StateId != 1397760650:
		t.Fatalf("Unexpected Origin-State-Id. Want 1397760650, have %d", d.StateId)
	case d.Vendors[0] != 10415:
		t.Fatalf("Unexpected Origin-State-Id. Want 10415, have %d", d.StateId)
	case d.Vendors[1] != 13:
		t.Fatalf("Unexpected Origin-State-Id. Want 13, have %d", d.StateId)
	case d.AuthAppId != 4:
		t.Fatalf("Unexpected Origin-State-Id. Want 4, have %d", d.AuthAppId)
	case d.InbandSecId != 0:
		t.Fatalf("Unexpected Origin-State-Id. Want 0, have %d", d.InbandSecId)
	case d.VSA.AuthAppId != 4:
		t.Fatalf("Unexpected Origin-State-Id. Want 4, have %d", d.VSA.AuthAppId)
	case d.VSA.VendorId != 10415:
		t.Fatalf("Unexpected Origin-State-Id. Want 10415, have %d", d.VSA.VendorId)
	case d.Firmware != 1:
		t.Fatalf("Unexpected Origin-State-Id. Want 1, have %d", d.Firmware)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	msg, err := ReadMessage(bytes.NewReader(testMessage), diamdict.Default)
	if err != nil {
		b.Fatal(err)
	}
	type CER struct {
		OriginHost  AVP    `avp:"Origin-Host"`
		OriginRealm *AVP   `avp:"Origin-Realm"`
		HostIP      net.IP `avp:"Host-IP-Address"`
		VendorId    int    `avp:"Vendor-Id"`
		ProductName string `avp:"Product-Name"`
		StateId     int    `avp:"Origin-State-Id"`
	}
	var cer CER
	for n := 0; n < b.N; n++ {
		msg.Unmarshal(&cer)
	}
}
