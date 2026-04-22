// Package common holds shared helpers for go-diameter interop test
// binaries under tests/go/.
package common

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
)

type Flags struct {
	Addr        string
	OriginHost  string
	OriginRealm string
	PeerHost    string
	PeerRealm   string
	HostIP      string
	Watchdog    time.Duration
	Timeout     time.Duration
	N           int
	VendorID    uint
	Network     string
}

func RegisterFlags(defaultAddr string) *Flags {
	f := &Flags{}
	flag.StringVar(&f.Addr, "addr", defaultAddr, "address ip:port")
	flag.StringVar(&f.OriginHost, "origin-host", "peer.test.local", "Origin-Host FQDN")
	flag.StringVar(&f.OriginRealm, "origin-realm", "test.local", "Origin-Realm")
	flag.StringVar(&f.PeerHost, "peer-host", "", "peer Origin-Host (for Destination-Host; empty = omit)")
	flag.StringVar(&f.PeerRealm, "peer-realm", "test.local", "peer Origin-Realm (for Destination-Realm)")
	flag.StringVar(&f.HostIP, "host-ip", "127.0.0.1", "Host-IP-Address AVP value")
	flag.DurationVar(&f.Watchdog, "watchdog", 0, "watchdog interval (0 disables)")
	flag.DurationVar(&f.Timeout, "timeout", 10*time.Second, "overall client timeout")
	flag.IntVar(&f.N, "n", 1, "message count (for loops)")
	flag.UintVar(&f.VendorID, "vendor", 13, "Vendor-Id")
	flag.StringVar(&f.Network, "network", "tcp", "tcp or sctp")
	return f
}

func Settings(f *Flags, productName string) *sm.Settings {
	ip := net.ParseIP(f.HostIP)
	if ip == nil {
		log.Fatalf("invalid host-ip: %q", f.HostIP)
	}
	return &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(f.OriginHost),
		OriginRealm:      datatype.DiameterIdentity(f.OriginRealm),
		VendorID:         datatype.Unsigned32(f.VendorID),
		ProductName:      datatype.UTF8String(productName),
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses:  []datatype.Address{datatype.Address(ip)},
	}
}

func Client(mux *sm.StateMachine, f *Flags, authAppIDs, vendorIDs []uint32) *sm.Client {
	c := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     f.Watchdog > 0,
		WatchdogInterval:   f.Watchdog,
	}
	for _, v := range vendorIDs {
		c.SupportedVendorID = append(c.SupportedVendorID,
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(v)))
	}
	for _, a := range authAppIDs {
		c.AuthApplicationID = append(c.AuthApplicationID,
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(a)))
	}
	return c
}

// SendDPR sends a Disconnect-Peer-Request and waits briefly for DPA
// before returning. The caller should Close the connection after.
func SendDPR(c diam.Conn, settings *sm.Settings) error {
	m := diam.NewRequest(diam.DisconnectPeer, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
	m.NewAVP(avp.DisconnectCause, avp.Mbit, 0, datatype.Enumerated(0)) // REBOOTING
	_, err := m.WriteTo(c)
	return err
}

func LogErrors(ec <-chan *diam.ErrorReport) {
	for e := range ec {
		log.Printf("error: %v", e.Error)
	}
}
