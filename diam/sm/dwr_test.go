// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/diamtest"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
)

// These tests use dictionary, settings and functions from sm_test.go.

func TestHandleDWR(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	mux.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	// Send CER first.
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code for CEA.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No CEA received")
	}
	// Send DWR.
	m = diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code for DWA.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No DWA received")
	}
}

func TestHandleDWR_Fail(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	mux.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	// Send CER first.
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code for CEA.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No CEA received")
	}
	// Send broken DWR (missing Origin-Host, etc).
	m = diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-sm.ErrorReports():
		if err.Error != smparser.ErrMissingOriginHost {
			t.Fatalf("Unexpected error. Want ErrMissingOriginHost, have %#v", err.Error)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No DWA received")
	}
}

// TestHandleDWR_OriginStateID asserts that the DWA carries the configured
// Origin-State-Id. Regression test: handleDWR used to append the AVP to the
// received request instead of the answer, so it never reached the wire.
func TestHandleDWR_OriginStateID(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	mux.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	// Send CER first to complete the handshake.
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err = m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code for CEA.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No CEA received")
	}
	// Send DWR and inspect the answer.
	m = diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	if _, err = m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		osid, err := resp.FindAVP(avp.OriginStateID, 0)
		if err != nil {
			t.Fatalf("DWA is missing Origin-State-Id (server has it configured).\n%s", resp)
		}
		got, ok := osid.Data.(datatype.Unsigned32)
		if !ok {
			t.Fatalf("Origin-State-Id has unexpected type %T", osid.Data)
		}
		if got != serverSettings.OriginStateID {
			t.Errorf("Origin-State-Id: got %d, want %d", got, serverSettings.OriginStateID)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No DWA received")
	}
}
