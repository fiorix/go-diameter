// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"crypto/tls"
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

func TestHandleCER_InbandSecurity_WithTLSConfig(t *testing.T) {
	// When the server has TLSConfig set, it should accept
	// Inband-Security-Id=1 and attempt TLS upgrade after CEA.
	// Since we can't easily do a real TLS upgrade in a unit test
	// (the client would also need to upgrade), we verify the server
	// accepts the CER and sends a success CEA instead of rejecting.
	tlsSettings := *serverSettings
	tlsSettings.TLSConfig = &tls.Config{} // non-nil signals TLS capability
	sm := New(&tlsSettings)
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
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	_, err = m.WriteTo(cli)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case resp := <-mc:
		// Should get a success CEA (not NO_COMMON_SECURITY)
		if testResultCode(resp, diam.NoCommonSecurity) {
			t.Fatal("Server rejected Inband-Security-Id=1 despite having TLSConfig")
		}
		if !testResultCode(resp, diam.Success) {
			t.Fatalf("Expected success CEA.\n%s", resp)
		}
	case err := <-mux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("No message received")
	}
}

func TestInbandTLSUpgrade_EndToEnd(t *testing.T) {
	// End-to-end test: plain TCP → CER/CEA with Inband-Security-Id=1 → TLS upgrade.
	// Both server and client perform the TLS handshake after CEA.

	// Generate a self-signed test certificate at runtime.
	cert, err := tls.X509KeyPair(testCert, testKey)
	if err != nil {
		t.Fatal(err)
	}

	serverTLSCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	clientTLSCfg := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Server settings with TLSConfig to enable inband upgrade.
	srvSettings := &Settings{
		OriginHost:       "srv",
		OriginRealm:      "test",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		TLSConfig:        serverTLSCfg,
	}
	sm := New(srvSettings)

	// Start a plain TCP server (not pre-TLS).
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()

	// Client settings with TLSConfig to enable inband upgrade.
	cliSettings := &Settings{
		OriginHost:       "cli",
		OriginRealm:      "test",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		TLSConfig:        clientTLSCfg,
	}

	cli := &Client{
		Handler:            New(cliSettings),
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		InbandSecurityID:   1, // Request TLS upgrade
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001)),
		},
	}

	// Dial plain TCP — the CER/CEA exchange + TLS upgrade happens inside.
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Client.Dial failed: %v", err)
	}
	defer conn.Close()

	// After successful dial, TLS should be active.
	if conn.TLS() == nil {
		t.Fatal("Expected TLS to be active after inband upgrade, but conn.TLS() is nil")
	}

	// Verify we can still exchange messages over the upgraded connection.
	// Send a DWR and expect a DWA back (the state machine handles DWR automatically).
	m := diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cliSettings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cliSettings.OriginRealm)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, cliSettings.OriginStateID)

	done := make(chan *diam.Message, 1)
	cli.Handler.HandleFunc("DWA", func(c diam.Conn, msg *diam.Message) {
		done <- msg
	})

	if _, err := m.WriteTo(conn); err != nil {
		t.Fatalf("Failed to send DWR over TLS connection: %v", err)
	}

	select {
	case dwa := <-done:
		if dwa == nil {
			t.Fatal("Received nil DWA")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for DWA over TLS connection")
	}
}

// Test certificate with 2048-bit RSA key for TLS upgrade tests.
// Generated with: openssl req -x509 -newkey rsa:2048 -nodes -subj '/CN=localhost' -addext 'subjectAltName=IP:127.0.0.1,IP:::1'
var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDLDCCAhSgAwIBAgIUQTnWs9bR1QQqy0SDFQfh//Qv/VEwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI2MDUwMTE5MTYwMloXDTM2MDQy
ODE5MTYwMlowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEA016HNlvsQ7gxfHhIXX/7+CZHOVStKsJPNx/0hW341nFu
69Z00U4Sl/cff8nJE07og5yDMkj1YEVi1fUvDFvyhiV2IM+672yEBi2DgVaoCm4d
u/31Gb1itPmXp/oxLDHuPaentM/qOCxFaoDLN3KtRKQKOYgYGGh112+pZW1iqLnE
FSQHZbe+1Byd76qqNTcYn0V2ieuJ6Mi2BrD1TXGClvdDOXFLcJ6ZP2yAfIQGjOTV
DVowVn/Tw5aoNO3QrYdWK737de9rtYHq33WSKANmPhfTsk9nCvoA1S/BbX12LY19
4rpmr0xRIVzeGBGnizaswEm+DCsNXdmXyXjRSkq6PwIDAQABo3YwdDAdBgNVHQ4E
FgQU5k+HVDdmq2WVnHoL4fy3qygu8xUwHwYDVR0jBBgwFoAU5k+HVDdmq2WVnHoL
4fy3qygu8xUwDwYDVR0TAQH/BAUwAwEB/zAhBgNVHREEGjAYhwR/AAABhxAAAAAA
AAAAAAAAAAAAAAABMA0GCSqGSIb3DQEBCwUAA4IBAQBcUip2F8D87R29m+g+6LbA
u1FaTIE4uPEMhqHDifcxkl/q1J7exSq8TMtwAQsNqg09scN//wpOCzYh1gHKmST5
+euP4X7VR/WPwkRqXZhUzVCQwJR5b46tJEnD7g+eNdGfCJJ1MhXDy/ghmC6BiMWg
5LyCAaPxaTTz0F6xccCPD0DTrBJNkydHJ9gbabKHxHUw0OkSXNHkd2jhiUfzgjQ+
BVhBRUY+f6sJ0uYmjrr5BWF2VvW/Fjng51ZScQVASv42R3VSER/grT466viHsRmT
+tPRn/CtBVO2P2ERJ/eWdJJFdW14yNbIrER56Xb59Z79WXI6vmIt9NJe79u+Ufq9
-----END CERTIFICATE-----`)

var testKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDTXoc2W+xDuDF8
eEhdf/v4Jkc5VK0qwk83H/SFbfjWcW7r1nTRThKX9x9/yckTTuiDnIMySPVgRWLV
9S8MW/KGJXYgz7rvbIQGLYOBVqgKbh27/fUZvWK0+Zen+jEsMe49p6e0z+o4LEVq
gMs3cq1EpAo5iBgYaHXXb6llbWKoucQVJAdlt77UHJ3vqqo1NxifRXaJ64noyLYG
sPVNcYKW90M5cUtwnpk/bIB8hAaM5NUNWjBWf9PDlqg07dCth1Yrvft172u1gerf
dZIoA2Y+F9OyT2cK+gDVL8FtfXYtjX3iumavTFEhXN4YEaeLNqzASb4MKw1d2ZfJ
eNFKSro/AgMBAAECggEAF2j4S8Z5j/SGEpWV2jkzFIRUzh45Qaucr2vMHr0T2thc
YyVw8b+WYptdsz8LlKZgLTd39mlLN/rnW/AYYmOKpF3gy/iF6T+ZDcAbuQb6fJE+
nNQfQdcOaCHesJ2OtajgDJcVhXqjo84PcCDMoRsD4r7SXRXcKVPkfVRiLBgl3a7l
wuXNOmkxcZOslvtDWAtoWPd3u8usc4XhxGIYKNwuG08N5pcJCcbq7W+Lf8fUD/pF
VukzqX5NLp/RoLyn6i+iKVEocuCMfmW5dczujbsYR/bqsCjFXJTvb710AJNEPwK9
QJIO0YUhni6P7M5+aS/vYwwp1y0fpBv+yEHezocsiQKBgQD6ydelW4tcgGpCOdi8
ab82J20KUpnebytmdYoq3PFnYxJRkyw9D0f5nlnGGm+0ZbAYGsl1v+RDoJ+opMOA
MguDrtakvqKjRlURWWsVHrVyWRfslULKKw0KDcuXiu8R79WvFmOL+Xo/0q+g52KO
JPFN5UCRXEhRoxaqFgfS4KtiuQKBgQDXwvtChPcUg5Z9Z34nTsQRyTMpq/38RPa5
oRAO83o8Mn7ZtuI6WYKdSfozZa7SgP02p4dzVor2eOFMiSeF1RCysdTHInuE19iL
UE/9uifYGbhEHXlUJ72P0o7Ll6Ko73Bg+i2TbckRs65Gk4u8ADU7ez5UMajB7bzr
QitsZwxotwKBgQDuAOg68eoMW4J8X1GlXeYtirUc+s80HeTeU+ZQT2Z6a7dS2408
VWhFKVahfy1L0sWP2rwel4IV/DYJYnR3EQeEbUUfDBxlP7YzxNyvKnmgj5T43Z6J
Jto1FGqG4z+HkkkE5QaMLLMsJtKurWkG5WBsQIlKan3nnBNCT64VH0sHYQKBgG4E
5K5UstDpEHG9thxBE8Wl/MrBAvACEnUxZcjZ6niLnxdRJCZwwiOGN2jB7tU0JOob
nvv3I0Du/qNSRK7/qFYWS9OHB8kDb04Kk99jbzHIW6eQB/Abm5Oc4Gd8WNsfzQQG
TfshPigioTknv1cMHBjKjUvNTqokmfK0eQP7v94dAoGACahWbGxB6zZbVwb2MQno
mgp27KEpMF7ObnCLDRBcw6dNjnbk+Sr6epxsV+QtzTcgxhF+O306fLiYeqL1o+uE
8oq9CIB6a4MN+huVhFnCLJ9bmpW/OSLPzgepXNZxWFP+UVEGONhn718fH0D4ytSc
LZo5zivmnoir47JtzE9yTp4=
-----END PRIVATE KEY-----`)
