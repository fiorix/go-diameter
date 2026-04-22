// diam-cc-server: minimal Credit-Control responder for App-Id 4
// (Ro/Gy) and App-Id 16777238 (Gx). Mirrors CC-Request-Type /
// CC-Request-Number from the request and returns Result-Code=2001.
package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

var (
	appID      = flag.Uint("app-id", diam.CHARGING_CONTROL_APP_ID, "4 (Ro/Gy) or 16777238 (Gx)")
	extraDicts = flag.String("extra-dicts", "", "comma-separated extra dictionary XML paths")
)

func main() {
	f := common.RegisterFlags(":3868")
	flag.Parse()

	if err := loadExtraDicts(*extraDicts); err != nil {
		log.Fatalf("load dicts: %v", err)
	}

	settings := common.Settings(f, "go-diameter-cc-server")
	mux := sm.New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: uint32(*appID), Code: diam.CreditControl, Request: true},
		diam.HandlerFunc(handleCCR(settings)))

	mux.HandleFunc("DPR", func(c diam.Conn, m *diam.Message) {
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		_, _ = a.WriteTo(c)
	})

	go common.LogErrors(mux.ErrorReports())

	log.Printf("cc server listening on %s (%s) app-id=%d", f.Addr, f.Network, *appID)
	if err := diam.ListenAndServeNetwork(f.Network, f.Addr, mux, nil); err != nil {
		log.Fatal(err)
	}
}

func handleCCR(s *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		sid, _ := m.FindAVP(avp.SessionID, 0)
		reqType, _ := m.FindAVP(avp.CCRequestType, 0)
		reqNum, _ := m.FindAVP(avp.CCRequestNumber, 0)
		a := m.Answer(diam.Success)
		if sid != nil {
			a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, sid.Data))
		}
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
		a.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(*appID))
		if reqType != nil {
			a.NewAVP(avp.CCRequestType, avp.Mbit, 0, reqType.Data)
		}
		if reqNum != nil {
			a.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, reqNum.Data)
		}
		if _, err := a.WriteTo(c); err != nil {
			log.Printf("CCA write: %v", err)
		}
	}
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
