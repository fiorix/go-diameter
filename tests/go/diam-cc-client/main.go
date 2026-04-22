// diam-cc-client: Credit-Control client covering RFC 4006 (Ro/Gy,
// App-Id 4) and 3GPP Gx (App-Id 16777238). Sends an INITIAL_REQUEST,
// one UPDATE_REQUEST, and a TERMINATION_REQUEST, expecting 2001 on each.
package main

import (
	"bytes"
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

const tgppVendorID = 10415

const (
	ccrInitial    = 1
	ccrUpdate     = 2
	ccrTermination = 3
)

var (
	appID      = flag.Uint("app-id", diam.CHARGING_CONTROL_APP_ID, "4 (Ro/Gy) or 16777238 (Gx)")
	extraDicts = flag.String("extra-dicts", "", "comma-separated extra dictionary XML paths")
)

func main() {
	f := common.RegisterFlags("127.0.0.1:3868")
	flag.Parse()

	if err := loadExtraDicts(*extraDicts); err != nil {
		log.Fatalf("load dicts: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	settings := common.Settings(f, "go-diameter-cc-client")
	mux := sm.New(settings)

	done := make(chan uint32, 8)
	mux.HandleIdx(
		diam.CommandIndex{AppID: uint32(*appID), Code: diam.CreditControl, Request: false},
		diam.HandlerFunc(func(c diam.Conn, m *diam.Message) {
			rc := resultCode(m)
			log.Printf("CCA Result-Code=%d from %s", rc, c.RemoteAddr())
			done <- rc
		}))

	go common.LogErrors(mux.ErrorReports())

	cli := common.Client(mux, f, nil, nil)
	if *appID == diam.CHARGING_CONTROL_APP_ID {
		cli.AuthApplicationID = []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID)),
		}
	} else {
		cli.VendorSpecificApplicationID = []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(tgppVendorID)),
				},
			}),
		}
		cli.SupportedVendorID = []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(tgppVendorID)),
		}
	}

	conn, err := cli.DialNetwork(f.Network, f.Addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	log.Println("handshake OK")

	meta, _ := smpeer.FromContext(conn.Context())
	sid := "session;" + strconv.Itoa(int(rand.Uint32()))

	for i, rtype := range []uint32{ccrInitial, ccrUpdate, ccrTermination} {
		if err := sendCCR(conn, settings, meta, sid, rtype, uint32(i)); err != nil {
			log.Fatalf("send CCR type=%d: %v", rtype, err)
		}
		if err := waitFor(done, f.Timeout); err != nil {
			log.Fatal(err)
		}
	}

	_ = common.SendDPR(conn, settings)
	time.Sleep(200 * time.Millisecond)
	conn.Close()
	log.Println("OK")
}

func waitFor(done <-chan uint32, timeout time.Duration) error {
	select {
	case rc := <-done:
		if rc != diam.Success {
			return errors.New("non-success CCA Result-Code=" + strconv.Itoa(int(rc)))
		}
		return nil
	case <-time.After(timeout):
		os.Exit(1)
		return errors.New("CCA timeout")
	}
}

func sendCCR(c diam.Conn, s *sm.Settings, meta *smpeer.Metadata, sid string, rtype, rnum uint32) error {
	m := diam.NewRequest(diam.CreditControl, uint32(*appID), dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(int32(rtype)))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(rnum))
	if *appID == diam.CHARGING_CONTROL_APP_ID {
		m.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("32251@3gpp.org"))
	}
	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(1)), // MSISDN
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String("15551234567")),
		},
	})
	log.Printf("CCR type=%d num=%d", rtype, rnum)
	_, err := m.WriteTo(c)
	return err
}

func resultCode(m *diam.Message) uint32 {
	rcAVP, _ := m.FindAVP(avp.ResultCode, 0)
	if rcAVP == nil {
		return 0
	}
	u, _ := rcAVP.Data.(datatype.Unsigned32)
	return uint32(u)
}

func loadExtraDicts(csv string) error {
	if csv == "" {
		return nil
	}
	for _, p := range splitCSV(csv) {
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if err := dict.Default.Load(bytes.NewReader(data)); err != nil {
			return err
		}
	}
	return nil
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}
