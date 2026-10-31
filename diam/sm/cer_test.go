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
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

// These tests use dictionary, settings and functions from sm_test.go.

func TestHandleCER_HandshakeMetadataTCP(t *testing.T) {
	testHandleCER_HandshakeMetadata(t, "tcp")
}

func testHandleCER_HandshakeMetadata(t *testing.T, network string) {
	sm := New(serverSettings)
	srv := diamtest.NewServerNetwork(network, sm, dict.Default)
	defer srv.Close()

	hsc := make(chan diam.Conn, 1)
	cli, err := diam.DialNetwork(network, srv.Addr, nil, dict.Default)
	if err != nil {
		t.Fatal(err)
	}

	ready := make(chan struct{})
	go func() {
		c := <-sm.HandshakeNotify()
		hsc <- c
		close(ready)
	}()

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
	<-ready

	c := <-hsc
	ctx := c.Context()
	meta, ok := smpeer.FromContext(ctx)
	if !ok {
		t.Fatal("Handshake ok but no context/metadata found")
	}
	if meta.OriginHost != clientSettings.OriginHost {
		t.Fatalf("Unexpected OriginHost. Want %q, have %q",
			clientSettings.OriginHost, meta.OriginHost)
	}
	if meta.OriginRealm != clientSettings.OriginRealm {
		t.Fatalf("Unexpected OriginRealm. Want %q, have %q",
			clientSettings.OriginRealm, meta.OriginRealm)
	}
}

func TestHandleCER_HandshakeMetadata_CustomIP(t *testing.T) {
	sm := New(serverSettings2)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()

	hsc := make(chan diam.Conn, 1)
	cli, err := diam.Dial(srv.Addr, nil, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	ready := make(chan struct{})
	go func() {
		close(ready)
		c := <-sm.HandshakeNotify()
		hsc <- c
	}()
	<-ready

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

	c := <-hsc
	ctx := c.Context()
	meta, ok := smpeer.FromContext(ctx)
	if !ok {
		t.Fatal("Handshake ok but no context/metadata found")
	}
	if meta.OriginHost != clientSettings.OriginHost {
		t.Fatalf("Unexpected OriginHost. Want %q, have %q",
			clientSettings.OriginHost, meta.OriginHost)
	}
	if meta.OriginRealm != clientSettings.OriginRealm {
		t.Fatalf("Unexpected OriginRealm. Want %q, have %q",
			clientSettings.OriginRealm, meta.OriginRealm)
	}
}

func TestHandleCER_Acct(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
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
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_Acct_Fail(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_Acct_Fail_CustomIP(t *testing.T) {
	sm := New(serverSettings2)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_VS_Acct(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001)),
		},
	})
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_VS_Acct_Fail(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000)),
		},
	})
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_Auth(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 1002, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(1002))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_Auth_Fail(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_VS_AuthTCP(t *testing.T) {
	testHandleCER_VS_Auth(t, "tcp")
}

func testHandleCER_VS_Auth(t *testing.T, network string) {
	sm := New(serverSettings)
	srv := diamtest.NewServerNetwork(network, sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.DialNetwork(network, srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 1002, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(1002)),
		},
	})
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_VS_Auth_FailTCP(t *testing.T) {
	testHandleCER_VS_Auth_Fail(t, "tcp")
}

func testHandleCER_VS_Auth_Fail(t *testing.T, network string) {
	sm := New(serverSettings)
	srv := diamtest.NewServerNetwork(network, sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.DialNetwork(network, srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000)),
		},
	})
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestHandleCER_InbandSecurity(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonSecurity) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

// TestHandleCER_ErrorCEAFlags checks that an error CEA (5010 here) is sent
// without the E bit: permanent failures use the command's own grammar
// (RFC 6733 §7.1.5), and the E bit marks a message that does not.
func TestHandleCER_ErrorCEAFlags(t *testing.T) {
	sm := New(serverSettings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()
	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		mc <- m
	})
	cli, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1000))
	if _, err = m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		if !testResultCode(resp, diam.NoCommonApplication) {
			t.Fatalf("Unexpected result code.\n%s", resp)
		}
		if resp.Header.CommandFlags&diam.ErrorFlag != 0 {
			t.Fatalf("Error CEA has the E bit set.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No CEA received")
	}
}

// TestHandleCER_PeerContextSetBeforeCEA checks that the connection carries the
// peer metadata before the CEA is written. With Server.MaxConcurrentHandlers
// set, the peer's first request after the CEA is dispatched on its own
// goroutine and handshakeOK drops it if the context is not yet set. OnCEA
// runs immediately before the CEA write, so it observes the order.
func TestHandleCER_PeerContextSetBeforeCEA(t *testing.T) {
	contextSet := make(chan bool, 1)
	settings := *serverSettings
	settings.OnCEA = func(c diam.Conn, m *diam.Message) {
		_, ok := smpeer.FromContext(c.Context())
		contextSet <- ok
	}
	sm := New(&settings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()

	cli, err := diam.Dial(srv.Addr, nil, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	m := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err := m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}
	select {
	case ok := <-contextSet:
		if !ok {
			t.Fatal("peer context not set when the CEA is written")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnCEA did not fire")
	}
}

// TestHandleCER_OriginStateID checks that the CEA carries the server's own
// Origin-State-Id, the value its DWAs carry, and never the one from the CER
// (RFC 6733 section 8.16). A peer that tracks Origin-State-Id resets the
// connection when the value changes between the CEA and a DWA.
func TestHandleCER_OriginStateID(t *testing.T) {
	const peerOriginStateID = datatype.Unsigned32(1)
	for _, tc := range []struct {
		name   string
		osid   datatype.Unsigned32 // the server's Settings.OriginStateID
		appID  uint32              // 1001 is supported, 1000 is not
		result uint32
	}{
		{"success", 1234, 1001, diam.Success},
		{"error", 1234, 1000, diam.NoCommonApplication},
		{"unset", 0, 1001, diam.Success},
	} {
		t.Run(tc.name, func(t *testing.T) {
			settings := *serverSettings
			settings.OriginStateID = tc.osid
			sm := New(&settings)
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
			receive := func(name string) *diam.Message {
				t.Helper()
				select {
				case m := <-mc:
					return m
				case err := <-mux.ErrorReports():
					t.Fatal(err)
				case <-time.After(time.Second):
					t.Fatalf("No %s received", name)
				}
				return nil
			}
			checkOriginStateID := func(name string, m *diam.Message) {
				t.Helper()
				a, err := m.FindAVP(avp.OriginStateID, 0)
				if tc.osid == 0 {
					if err == nil {
						t.Fatalf("%s has Origin-State-Id %v, want none.\n%s", name, a.Data, m)
					}
					return
				}
				if err != nil {
					t.Fatalf("%s is missing Origin-State-Id.\n%s", name, m)
				}
				if got := a.Data.(datatype.Unsigned32); got != tc.osid {
					t.Fatalf("%s Origin-State-Id: got %d, want %d", name, got, tc.osid)
				}
			}

			m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
			m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
			m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
			m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
			m.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
			m.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
			m.NewAVP(avp.OriginStateID, avp.Mbit, 0, peerOriginStateID)
			m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(tc.appID))
			if _, err := m.WriteTo(cli); err != nil {
				t.Fatal(err)
			}
			cea := receive("CEA")
			if !testResultCode(cea, tc.result) {
				t.Fatalf("Unexpected result code.\n%s", cea)
			}
			checkOriginStateID("CEA", cea)
			if tc.result != diam.Success {
				return
			}

			m = diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
			m.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
			m.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
			m.NewAVP(avp.OriginStateID, avp.Mbit, 0, peerOriginStateID)
			if _, err := m.WriteTo(cli); err != nil {
				t.Fatal(err)
			}
			checkOriginStateID("DWA", receive("DWA"))
		})
	}
}
