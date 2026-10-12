// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/diamtest"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// testTLSConfigs builds a matching server/client pair from a certificate
// generated for this run, so the tests do not carry an expiring fixture.
func testTLSConfigs(t *testing.T) (server, client *tls.Config) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		DNSNames:              []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return &tls.Config{
			Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}},
		}, &tls.Config{
			RootCAs:    pool,
			ServerName: "localhost",
		}
}

func inbandSettings(host string, tlsCfg *tls.Config) *Settings {
	return &Settings{
		OriginHost:       datatype.DiameterIdentity(host),
		OriginRealm:      "test",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		TLSConfig:        tlsCfg,
	}
}

func inbandClient(sm *StateMachine, requested uint32) *Client {
	return &Client{
		Handler:            sm,
		MaxRetransmits:     3,
		RetransmitInterval: 500 * time.Millisecond,
		InbandSecurityID:   requested,
		AcctApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001)),
		},
	}
}

// The upgrade must work end to end over plain TCP, and the connection must
// still carry traffic afterwards.
func TestInbandTLSUpgradeEndToEnd(t *testing.T) {
	srvTLS, cliTLS := testTLSConfigs(t)
	srv := diamtest.NewServer(New(inbandSettings("srv", srvTLS)), dict.Default)
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	if conn.TLS() == nil {
		t.Fatal("conn.TLS() is nil: the connection was not upgraded")
	}
	assertDWRRoundTrip(t, cli, conn)
}

// With a CloseNotifier registered, the connection's reads come through a pipe
// fed by a copy goroutine. The upgrade has to read through that chain instead
// of the socket, or the handshake and the copy goroutine compete for the same
// bytes and the upgrade stalls.
func TestInbandTLSUpgradeWithCloseNotifier(t *testing.T) {
	srvTLS, cliTLS := testTLSConfigs(t)
	srv := diamtest.NewUnstartedServer(New(inbandSettings("srv", srvTLS)), dict.Default)
	notified := make(chan struct{}, 1)
	srv.Config.OnNewConnection = func(c diam.Conn) {
		cn, ok := c.(diam.CloseNotifier)
		if !ok {
			return
		}
		select {
		case notified <- struct{}{}:
		default:
		}
		go func() { <-cn.CloseNotify() }()
	}
	srv.Start()
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial with a CloseNotifier registered: %v", err)
	}
	defer conn.Close()
	select {
	case <-notified:
	case <-time.After(2 * time.Second):
		t.Fatal("OnNewConnection never ran")
	}
	if conn.TLS() == nil {
		t.Fatal("conn.TLS() is nil: the connection was not upgraded")
	}
	assertDWRRoundTrip(t, cli, conn)
}

func assertDWRRoundTrip(t *testing.T, cli *Client, conn diam.Conn) {
	t.Helper()
	done := make(chan *diam.Message, 1)
	cli.Handler.HandleFunc("DWA", func(c diam.Conn, m *diam.Message) {
		select {
		case done <- m:
		default:
		}
	})
	m := diam.NewRequest(diam.DeviceWatchdog, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cli.Handler.cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cli.Handler.cfg.OriginRealm)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, cli.Handler.cfg.OriginStateID)
	if _, err := m.WriteTo(conn); err != nil {
		t.Fatalf("send DWR over the upgraded connection: %v", err)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("no DWA over the upgraded connection")
	}
}

// The client must not start a handshake unless the peer selected TLS. A CEA
// that omits Inband-Security-Id, or answers 0, has declined; throwing a
// ClientHello at such a peer breaks the connection. The peer here is a raw
// socket, so the test sees exactly what the client puts on the wire.
func TestClientDoesNotUpgradeWithoutAgreement(t *testing.T) {
	for _, tc := range []struct {
		desc string
		echo bool
		id   uint32
	}{
		{"CEA omits Inband-Security-Id", false, 0},
		{"CEA answers NO_INBAND_SECURITY", true, 0},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()
			sawClientHello := make(chan bool, 1)
			go func() {
				c, err := ln.Accept()
				if err != nil {
					sawClientHello <- false
					return
				}
				defer c.Close()
				m, err := diam.ReadMessage(c, dict.Default)
				if err != nil {
					sawClientHello <- false
					return
				}
				a := m.Answer(diam.Success)
				a.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
				a.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
				a.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
				a.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
				a.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
				a.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
				if tc.echo {
					a.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(tc.id))
				}
				a.WriteTo(c)
				// 0x16 is a TLS handshake record: the client upgraded anyway.
				c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
				b := make([]byte, 1)
				n, _ := c.Read(b)
				sawClientHello <- n == 1 && b[0] == 0x16
			}()

			_, cliTLS := testTLSConfigs(t)
			cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
			cli.MaxRetransmits = 0
			cli.RetransmitInterval = 2 * time.Second
			conn, err := cli.Dial(ln.Addr().String())
			if conn != nil {
				defer conn.Close()
			}
			if <-sawClientHello {
				t.Fatal("client sent a TLS ClientHello although the peer did not select TLS")
			}
			// Continuing in cleartext would silently drop the protection
			// the client asked for, so the handshake must fail.
			if err == nil {
				t.Fatal("Dial = nil error, want a declined-security failure")
			}
		})
	}
}

// Advertising Inband-Security-Id without configuring TLSConfig must behave
// exactly as before: the AVP is sent, no upgrade is attempted, Dial succeeds.
func TestClientAdvertisesWithoutTLSConfig(t *testing.T) {
	var received uint32
	mux := diam.NewServeMux()
	mux.HandleFunc("CER", func(c diam.Conn, m *diam.Message) {
		if a, _ := m.FindAVP(avp.InbandSecurityID, 0); a != nil {
			received = uint32(a.Data.(datatype.Unsigned32))
		}
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
		a.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
		a.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
		a.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
		a.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
		a.WriteTo(c)
	})
	srv := diamtest.NewServer(mux, dict.Default)
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", nil)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	if received != 1 {
		t.Errorf("server saw Inband-Security-Id=%d, want 1", received)
	}
	if conn.TLS() != nil {
		t.Error("connection was upgraded without TLSConfig")
	}
}

// A server that cannot upgrade must refuse the mechanism before promising it,
// answering DIAMETER_NO_COMMON_SECURITY rather than sending a success CEA it
// cannot honour.
func TestServerRefusesInbandTLSWithoutTLSConfig(t *testing.T) {
	srv := diamtest.NewServer(New(inbandSettings("srv", nil)), dict.Default)
	defer srv.Close()

	mc := make(chan *diam.Message, 1)
	mux := diam.NewServeMux()
	mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
		select {
		case mc <- m:
		default:
		}
	})
	conn, err := diam.Dial(srv.Addr, mux, dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("cli"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	m.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err := m.WriteTo(conn); err != nil {
		t.Fatal(err)
	}
	select {
	case cea := <-mc:
		rc, err := cea.FindAVP(avp.ResultCode, 0)
		if err != nil {
			t.Fatalf("CEA has no Result-Code: %v", err)
		}
		if got := uint32(rc.Data.(datatype.Unsigned32)); got != diam.NoCommonSecurity {
			t.Errorf("Result-Code = %d, want %d (DIAMETER_NO_COMMON_SECURITY)",
				got, diam.NoCommonSecurity)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no CEA received")
	}
}

// The server must echo the selected mechanism only when it will upgrade, so
// the peer knows to start its handshake.
func TestServerEchoesSelectedMechanism(t *testing.T) {
	srvTLS, _ := testTLSConfigs(t)
	for _, tc := range []struct {
		desc     string
		tlsCfg   *tls.Config
		request  uint32
		wantEcho bool
	}{
		{"TLS requested and available", srvTLS, 1, true},
		{"no security requested", srvTLS, 0, false},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			sm := New(inbandSettings("srv", tc.tlsCfg))
			srv := diamtest.NewServer(sm, dict.Default)
			defer srv.Close()

			mc := make(chan *diam.Message, 1)
			mux := diam.NewServeMux()
			mux.HandleFunc("CEA", func(c diam.Conn, m *diam.Message) {
				select {
				case mc <- m:
				default:
				}
			})
			conn, err := diam.Dial(srv.Addr, mux, dict.Default)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			m := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
			m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("cli"))
			m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
			m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
			m.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(tc.request))
			m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
			m.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
			m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
			if _, err := m.WriteTo(conn); err != nil {
				t.Fatal(err)
			}
			select {
			case cea := <-mc:
				a, _ := cea.FindAVP(avp.InbandSecurityID, 0)
				if tc.wantEcho && a == nil {
					t.Error("CEA has no Inband-Security-Id, want it echoed")
				}
				if !tc.wantEcho && a != nil {
					t.Errorf("CEA echoed Inband-Security-Id=%v, want none", a.Data)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("no CEA received")
			}
		})
	}
}

// A diam.Conn may be inspected from any goroutine -- the handshake timer
// closes it, and applications log RemoteAddr or poll TLS from watchers they
// start in OnNewConnection. The upgrade swaps the connection's transport, so
// those accessors must stay safe while it does.
func TestInbandTLSUpgradeConcurrentAccessors(t *testing.T) {
	srvTLS, cliTLS := testTLSConfigs(t)
	srv := diamtest.NewUnstartedServer(New(inbandSettings("srv", srvTLS)), dict.Default)
	upgraded := make(chan struct{})
	srv.Config.OnNewConnection = func(c diam.Conn) {
		go func() {
			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) {
				_ = c.RemoteAddr()
				_ = c.LocalAddr()
				if c.TLS() != nil {
					close(upgraded)
					return
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}
	srv.Start()
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	select {
	case <-upgraded:
	case <-time.After(3 * time.Second):
		t.Fatal("server side never observed the upgraded connection")
	}
}

// If the peer accepts the upgrade but never starts the handshake, the
// handshake timer must still be able to close the connection rather than
// leaving the server blocked in the TLS handshake.
func TestInbandTLSStalledHandshakeIsClosed(t *testing.T) {
	srvTLS, _ := testTLSConfigs(t)
	settings := inbandSettings("srv", srvTLS)
	settings.HandshakeTimeout = 200 * time.Millisecond
	srv := diamtest.NewServer(New(settings), dict.Default)
	defer srv.Close()

	raw, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	cer := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	cer.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("cli"))
	cer.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	cer.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	cer.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	cer.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	cer.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
	cer.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err := cer.WriteTo(raw); err != nil {
		t.Fatal(err)
	}
	raw.SetReadDeadline(time.Now().Add(3 * time.Second))
	cea, err := diam.ReadMessage(raw, dict.Default)
	if err != nil {
		t.Fatalf("read CEA: %v", err)
	}
	if a, _ := cea.FindAVP(avp.InbandSecurityID, 0); a == nil {
		t.Fatal("CEA does not select Inband-Security-Id")
	}
	// Send nothing. The server is now waiting for a ClientHello.
	buf := make([]byte, 1)
	if _, err := raw.Read(buf); err == nil {
		t.Fatal("expected the server to close the stalled connection")
	} else if ne, ok := err.(net.Error); ok && ne.Timeout() {
		t.Fatal("server did not close a connection stalled in the TLS handshake")
	}
}

// noUpgradeConn stands in for a transport that cannot be upgraded, such as an
// SCTP association. It records whether an upgrade was attempted anyway.
type noUpgradeConn struct {
	diam.Conn
	attempted *bool
}

func (c noUpgradeConn) CanStartTLS() bool { return false }
func (c noUpgradeConn) StartTLS(*tls.Config) error {
	*c.attempted = true
	return errors.New("not upgradable")
}
func (c noUpgradeConn) StartTLSClient(*tls.Config) error {
	*c.attempted = true
	return errors.New("not upgradable")
}

// A server with TLSConfig set must still refuse in-band TLS on a connection
// it cannot upgrade, answering DIAMETER_NO_COMMON_SECURITY instead of sending
// a success CEA and then failing (or panicking) in the handshake.
func TestServerRefusesInbandTLSOnNonUpgradableConn(t *testing.T) {
	srvTLS, _ := testTLSConfigs(t)
	smach := New(inbandSettings("srv", srvTLS))
	var attempted bool
	wrapped := diam.HandlerFunc(func(c diam.Conn, m *diam.Message) {
		smach.ServeDIAM(noUpgradeConn{Conn: c, attempted: &attempted}, m)
	})
	srv := diamtest.NewServer(wrapped, dict.Default)
	defer srv.Close()

	raw, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	cer := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	cer.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("cli"))
	cer.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	cer.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	cer.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	cer.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	cer.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
	cer.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err := cer.WriteTo(raw); err != nil {
		t.Fatal(err)
	}
	raw.SetReadDeadline(time.Now().Add(3 * time.Second))
	cea, err := diam.ReadMessage(raw, dict.Default)
	if err != nil {
		t.Fatalf("read CEA: %v", err)
	}
	rc, err := cea.FindAVP(avp.ResultCode, 0)
	if err != nil {
		t.Fatalf("CEA has no Result-Code: %v", err)
	}
	if got := uint32(rc.Data.(datatype.Unsigned32)); got != diam.NoCommonSecurity {
		t.Errorf("Result-Code = %d, want %d (DIAMETER_NO_COMMON_SECURITY)", got, diam.NoCommonSecurity)
	}
	if attempted {
		t.Error("StartTLS was attempted on a connection that reported it cannot upgrade")
	}
}

// A client that only advertises Inband-Security-Id=1 (no TLSConfig) cannot
// carry on in cleartext once the peer has selected TLS, because the peer will
// now start the handshake (RFC 6733 §13). Dial must fail instead of hanging.
func TestClientWithoutTLSConfigFailsWhenPeerSelectsTLS(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		m, err := diam.ReadMessage(c, dict.Default)
		if err != nil {
			return
		}
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
		a.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
		a.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
		a.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
		a.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
		a.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
		a.WriteTo(c)
		time.Sleep(time.Second)
	}()

	cli := inbandClient(New(inbandSettings("cli", nil)), 1)
	cli.MaxRetransmits = 0
	cli.RetransmitInterval = 2 * time.Second
	conn, err := cli.Dial(ln.Addr().String())
	if err == nil {
		conn.Close()
		t.Fatal("Dial = nil error, want a failure: the peer selected TLS and this client cannot do it")
	}
}

// A Handler that wraps the StateMachine without forwarding InbandTLSEnabled
// gets concurrent CER dispatch, where the serve loop may be reading during an
// upgrade. The upgrade must then be refused up front with
// DIAMETER_NO_COMMON_SECURITY, not attempted.
func TestServerRefusesInbandTLSUnderConcurrentDispatchWithoutOptIn(t *testing.T) {
	srvTLS, _ := testTLSConfigs(t)
	smach := New(inbandSettings("srv", srvTLS))
	wrapped := diam.HandlerFunc(smach.ServeDIAM) // hides InbandTLSEnabled
	srv := diamtest.NewUnstartedServer(wrapped, dict.Default)
	srv.Config.MaxConcurrentHandlers = 8
	srv.Start()
	defer srv.Close()

	raw, err := net.Dial("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	cer := diam.NewRequest(diam.CapabilitiesExchange, 0, dict.Default)
	cer.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("cli"))
	cer.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	cer.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	cer.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	cer.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	cer.NewAVP(avp.InbandSecurityID, avp.Mbit, 0, datatype.Unsigned32(1))
	cer.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1001))
	if _, err := cer.WriteTo(raw); err != nil {
		t.Fatal(err)
	}
	raw.SetReadDeadline(time.Now().Add(3 * time.Second))
	cea, err := diam.ReadMessage(raw, dict.Default)
	if err != nil {
		t.Fatalf("read CEA: %v", err)
	}
	rc, err := cea.FindAVP(avp.ResultCode, 0)
	if err != nil {
		t.Fatalf("CEA has no Result-Code: %v", err)
	}
	if got := uint32(rc.Data.(datatype.Unsigned32)); got != diam.NoCommonSecurity {
		t.Errorf("Result-Code = %d, want %d (DIAMETER_NO_COMMON_SECURITY)", got, diam.NoCommonSecurity)
	}
}

// With the opt-in visible, in-band TLS works under concurrent dispatch too.
func TestInbandTLSUpgradeUnderConcurrentDispatch(t *testing.T) {
	srvTLS, cliTLS := testTLSConfigs(t)
	srv := diamtest.NewUnstartedServer(New(inbandSettings("srv", srvTLS)), dict.Default)
	srv.Config.MaxConcurrentHandlers = 8
	srv.Start()
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	if conn.TLS() == nil {
		t.Fatal("conn.TLS() is nil: the connection was not upgraded")
	}
	assertDWRRoundTrip(t, cli, conn)
}

// CloseNotify registered while the CER is being handled -- after the serve
// loop's last read -- leaves the notifier's pipe set up but not yet running.
// The upgrade must start it on the cleartext socket; otherwise the next read
// would start it with the TLS stream replaced by raw ciphertext.
func TestInbandTLSUpgradeWithCloseNotifyDuringCER(t *testing.T) {
	srvTLS, cliTLS := testTLSConfigs(t)
	settings := inbandSettings("srv", srvTLS)
	gone := make(chan struct{})
	settings.OnCER = func(c diam.Conn, m *diam.Message) {
		cn, ok := c.(diam.CloseNotifier)
		if !ok {
			return
		}
		ch := cn.CloseNotify()
		go func() { <-ch; close(gone) }()
	}
	srv := diamtest.NewServer(New(settings), dict.Default)
	defer srv.Close()

	cli := inbandClient(New(inbandSettings("cli", cliTLS)), 1)
	conn, err := cli.Dial(srv.Addr)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	if conn.TLS() == nil {
		t.Fatal("conn.TLS() is nil: the connection was not upgraded")
	}
	assertDWRRoundTrip(t, cli, conn)
	conn.Close()
	select {
	case <-gone:
	case <-time.After(3 * time.Second):
		t.Fatal("CloseNotify registered during the CER never fired after the upgrade")
	}
}
