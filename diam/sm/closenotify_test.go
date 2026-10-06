// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// closeNotifySettings returns settings for one end of the CloseNotify tests.
func closeNotifySettings(host string) *Settings {
	return &Settings{
		OriginHost:       datatype.DiameterIdentity(host),
		OriginRealm:      "test",
		VendorID:         10415,
		ProductName:      "closenotify-test",
		FirmwareRevision: 1,
		HostIPAddresses:  []datatype.Address{datatype.Address(net.ParseIP("127.0.0.1"))},
	}
}

// startClosingServer accepts associations and closes each one 300ms after
// the CER/CEA handshake, so a client that calls CloseNotify straight after
// Dial returns is genuinely early, and one that waits longer is late.
func startClosingServer(t *testing.T) string {
	t.Helper()
	mux := New(closeNotifySettings("srv.test"))
	hs := mux.HandshakeNotify()
	go func() {
		for c := range hs {
			go func(c diam.Conn) { time.Sleep(300 * time.Millisecond); c.Close() }(c)
		}
	}()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go diam.Serve(ln, mux)
	return ln.Addr().String()
}

// dialNoWatchdog dials with the watchdog off: the watchdog also calls
// CloseNotify, which would create the channel early and hide the bug.
func dialNoWatchdog(t *testing.T, addr string) diam.Conn {
	t.Helper()
	cli := &Client{
		Dict:              dict.Default,
		Handler:           New(closeNotifySettings("cli.test")),
		EnableWatchdog:    false,
		SupportedVendorID: []*diam.AVP{diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(10415))},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(16777251)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
				},
			}),
		},
	}
	c, err := cli.DialNetwork("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func closeFires(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	case <-time.After(3 * time.Second):
		return false
	}
}

// TestCloseNotifyBeforeClose is the control: CloseNotify called before the
// peer closes must fire when it does.
func TestCloseNotifyBeforeClose(t *testing.T) {
	c := dialNoWatchdog(t, startClosingServer(t))
	if !closeFires(c.(diam.CloseNotifier).CloseNotify()) {
		t.Fatal("CloseNotify called before the close did not fire; the test server did not close the association")
	}
}

// TestCloseNotifyAfterClose: CloseNotify called AFTER the connection has
// already gone away must still report it. Before the fix the channel was
// created after the close and never closed, so a caller blocked forever.
func TestCloseNotifyAfterClose(t *testing.T) {
	c := dialNoWatchdog(t, startClosingServer(t))
	time.Sleep(1500 * time.Millisecond) // the server closes at +300ms
	if !closeFires(c.(diam.CloseNotifier).CloseNotify()) {
		t.Fatal("CloseNotify called after the connection closed never fired")
	}
}
