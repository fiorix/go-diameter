// diam-s6a-server: HSS-style S6a responder. Answers AIR and ULR with
// canonical 3GPP AVPs. Interoperates with FD dict_s6a.fdx or another
// go-diameter peer.
package main

import (
	"flag"
	"log"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

const (
	s6aVendorID = 10415
	s6aAppID    = diam.TGPP_S6A_APP_ID
)

func main() {
	f := common.RegisterFlags(":3868")
	flag.Parse()

	settings := common.Settings(f, "go-diameter-s6a-server")
	mux := sm.New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: s6aAppID, Code: diam.AuthenticationInformation, Request: true},
		diam.HandlerFunc(handleAIR(settings)))
	mux.HandleIdx(
		diam.CommandIndex{AppID: s6aAppID, Code: diam.UpdateLocation, Request: true},
		diam.HandlerFunc(handleULR(settings)))

	mux.HandleFunc("DPR", func(c diam.Conn, m *diam.Message) {
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		_, _ = a.WriteTo(c)
	})

	go common.LogErrors(mux.ErrorReports())

	log.Printf("s6a server listening on %s (%s)", f.Addr, f.Network)
	if err := diam.ListenAndServeNetwork(f.Network, f.Addr, mux, nil); err != nil {
		log.Fatal(err)
	}
}

func originAVPs(a *diam.Message, s *sm.Settings) {
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
}

func handleAIR(s *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		sid, _ := m.FindAVP(avp.SessionID, 0)
		a := m.Answer(diam.Success)
		if sid != nil {
			a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, sid.Data))
		}
		originAVPs(a, s)
		a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
		a.NewAVP(avp.AuthenticationInfo, avp.Mbit, s6aVendorID, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.EUTRANVector, avp.Mbit, s6aVendorID, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, s6aVendorID, datatype.OctetString("\x94\xbf/T\xc3v\xf3\x0e\x87\x83\x06k'\x18Z\x19")),
						diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, s6aVendorID, datatype.OctetString("F\xf0\"\xb9%#\xf58")),
						diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, s6aVendorID, datatype.OctetString("\xc7G!;\xad~\x80\x00)\x08o%\x11\x0cP_")),
						diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, s6aVendorID, datatype.OctetString("\xbf\x00\xf9\x80h3\"\x0e\xa1\x1c\xfa\x93\x03@\xd6\xf8\x02\xd51Y\xeb\xc4\x9d=\t\x14{\xeb!\xec\xcb:")),
					},
				}),
			},
		})
		if _, err := a.WriteTo(c); err != nil {
			log.Printf("AIA write: %v", err)
		}
	}
}

func handleULR(s *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		sid, _ := m.FindAVP(avp.SessionID, 0)
		a := m.Answer(diam.Success)
		if sid != nil {
			a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, sid.Data))
		}
		originAVPs(a, s)
		a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
		a.NewAVP(avp.ULAFlags, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(1))
		a.NewAVP(avp.SubscriptionData, avp.Mbit, s6aVendorID, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.MSISDN, avp.Mbit|avp.Vbit, s6aVendorID, datatype.OctetString("12345")),
				diam.NewAVP(avp.AccessRestrictionData, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(47)),
				diam.NewAVP(avp.SubscriberStatus, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(0)),
				diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(2)),
				diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, s6aVendorID, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(500)),
						diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, s6aVendorID, datatype.Unsigned32(500)),
					},
				}),
			},
		})
		if _, err := a.WriteTo(c); err != nil {
			log.Printf("ULA write: %v", err)
		}
	}
}
