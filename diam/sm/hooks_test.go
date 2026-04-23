// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/diamtest"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// TestOnCERHook verifies that Settings.OnCER fires on a received CER before
// the state machine's default handler, and that the default handshake still
// completes. Regression test for #150.
func TestOnCERHook(t *testing.T) {
	var onCERCalls int32
	settings := *serverSettings
	settings.OnCER = func(c diam.Conn, m *diam.Message) {
		if m.Header.CommandCode != diam.CapabilitiesExchange {
			t.Errorf("OnCER invoked for non-CER command %d", m.Header.CommandCode)
		}
		atomic.AddInt32(&onCERCalls, 1)
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
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	m.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	if _, err := m.WriteTo(cli); err != nil {
		t.Fatal(err)
	}

	select {
	case <-sm.HandshakeNotify():
	case <-time.After(2 * time.Second):
		t.Fatal("handshake did not complete after OnCER hook")
	}

	if got := atomic.LoadInt32(&onCERCalls); got != 1 {
		t.Fatalf("OnCER calls = %d, want 1", got)
	}
}

// TestOnDWRHook verifies that Settings.OnDWR fires on a received DWR after
// the handshake and before the default DWR handler responds with DWA.
func TestOnDWRHook(t *testing.T) {
	onDWR := make(chan struct{}, 1)
	settings := *serverSettings
	settings.OnDWR = func(c diam.Conn, m *diam.Message) {
		if m.Header.CommandCode != diam.DeviceWatchdog {
			t.Errorf("OnDWR invoked for non-DWR command %d", m.Header.CommandCode)
		}
		onDWR <- struct{}{}
	}

	sm := New(&settings)
	srv := diamtest.NewServer(sm, dict.Default)
	defer srv.Close()

	cli, err := diam.Dial(srv.Addr, nil, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	cer := diam.NewRequest(diam.CapabilitiesExchange, 1001, dict.Default)
	cer.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	cer.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	cer.NewAVP(avp.HostIPAddress, avp.Mbit, 0, localhostAddress)
	cer.NewAVP(avp.VendorID, avp.Mbit, 0, clientSettings.VendorID)
	cer.NewAVP(avp.ProductName, 0, 0, clientSettings.ProductName)
	cer.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1))
	cer.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	cer.NewAVP(avp.FirmwareRevision, 0, 0, clientSettings.FirmwareRevision)
	if _, err := cer.WriteTo(cli); err != nil {
		t.Fatal(err)
	}

	select {
	case <-sm.HandshakeNotify():
	case <-time.After(2 * time.Second):
		t.Fatal("handshake did not complete")
	}

	dwr := diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	dwr.NewAVP(avp.OriginHost, avp.Mbit, 0, clientSettings.OriginHost)
	dwr.NewAVP(avp.OriginRealm, avp.Mbit, 0, clientSettings.OriginRealm)
	if _, err := dwr.WriteTo(cli); err != nil {
		t.Fatal(err)
	}

	select {
	case <-onDWR:
	case <-time.After(2 * time.Second):
		t.Fatal("OnDWR hook did not fire")
	}
}
