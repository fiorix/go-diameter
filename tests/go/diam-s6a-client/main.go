// diam-s6a-client: sends AIR + ULR to an S6a peer (FD or go-diameter)
// and exits 0 iff both get Result-Code=2001 within -timeout.
package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

const (
	s6aVendorID = 10415
	s6aAppID    = diam.TGPP_S6A_APP_ID
	ulrFlags    = 1<<1 | 1<<5
)

var (
	imsi    = flag.String("imsi", "001010000000001", "User IMSI")
	plmnID  = flag.String("plmnid", "\x00\xF1\x10", "Visited-PLMN-Id")
	vectors = flag.Uint("vectors", 1, "Number-Of-Requested-Vectors")
)

func main() {
	f := common.RegisterFlags("127.0.0.1:3868")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	settings := common.Settings(f, "go-diameter-s6a-client")
	mux := sm.New(settings)

	done := make(chan uint32, 4)
	mux.HandleIdx(
		diam.CommandIndex{AppID: s6aAppID, Code: diam.AuthenticationInformation, Request: false},
		answerHandler("AIA", done))
	mux.HandleIdx(
		diam.CommandIndex{AppID: s6aAppID, Code: diam.UpdateLocation, Request: false},
		answerHandler("ULA", done))

	go common.LogErrors(mux.ErrorReports())

	cli := common.Client(mux, f, nil, []uint32{s6aVendorID})
	cli.VendorSpecificApplicationID = []*diam.AVP{
		diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(s6aAppID)),
				diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(s6aVendorID)),
			},
		}),
	}

	conn, err := cli.DialNetwork(f.Network, f.Addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	log.Println("handshake OK")

	meta, ok := smpeer.FromContext(conn.Context())
	if !ok {
		log.Fatal("no peer metadata")
	}

	if err := sendAIR(conn, settings, meta); err != nil {
		log.Fatalf("send AIR: %v", err)
	}
	if err := waitFor("AIA", done, f.Timeout); err != nil {
		log.Fatal(err)
	}
	if err := sendULR(conn, settings, meta); err != nil {
		log.Fatalf("send ULR: %v", err)
	}
	if err := waitFor("ULA", done, f.Timeout); err != nil {
		log.Fatal(err)
	}

	_ = common.SendDPR(conn, settings)
	time.Sleep(200 * time.Millisecond)
	conn.Close()
	log.Println("OK")
}

func answerHandler(tag string, done chan<- uint32) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		rcAVP, err := m.FindAVP(avp.ResultCode, 0)
		var rc uint32
		if err == nil && rcAVP != nil {
			if u, ok := rcAVP.Data.(datatype.Unsigned32); ok {
				rc = uint32(u)
			}
		}
		log.Printf("%s Result-Code=%d from %s", tag, rc, c.RemoteAddr())
		done <- rc
	}
}

func waitFor(tag string, done <-chan uint32, timeout time.Duration) error {
	select {
	case rc := <-done:
		if rc != diam.Success {
			return errors.New(tag + " non-success Result-Code=" + strconv.Itoa(int(rc)))
		}
		return nil
	case <-time.After(timeout):
		os.Exit(1)
		return errors.New(tag + " timeout")
	}
}

func sendAIR(c diam.Conn, s *sm.Settings, meta *smpeer.Metadata) error {
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.AuthenticationInformation, s6aAppID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(*imsi))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, s6aVendorID, datatype.OctetString(*plmnID))
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, s6aVendorID, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.NumberOfRequestedVectors, avp.Vbit|avp.Mbit, s6aVendorID, datatype.Unsigned32(*vectors)),
			diam.NewAVP(avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, s6aVendorID, datatype.Unsigned32(0)),
		},
	})
	log.Printf("AIR to %s", c.RemoteAddr())
	_, err := m.WriteTo(c)
	return err
}

func sendULR(c diam.Conn, s *sm.Settings, meta *smpeer.Metadata) error {
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))
	m := diam.NewRequest(diam.UpdateLocation, s6aAppID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(*imsi))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(0))
	m.NewAVP(avp.RATType, avp.Mbit, s6aVendorID, datatype.Enumerated(1004))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, s6aVendorID, datatype.Unsigned32(ulrFlags))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, s6aVendorID, datatype.OctetString(*plmnID))
	log.Printf("ULR to %s", c.RemoteAddr())
	_, err := m.WriteTo(c)
	return err
}
