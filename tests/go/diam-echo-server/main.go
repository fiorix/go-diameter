// diam-echo-server: echoes the test_app Test command with Result-Code
// 2001 and the same Test-AVP payload. Interoperates with FD
// test_app.fdx mode=client.
package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

const (
	testVendorID = 999999
	testAppID    = 123456
	testCmdCode  = 234567
	testAVPCode  = 345678
)

func main() {
	f := common.RegisterFlags(":3868")
	dictFile := flag.String("dict", "/etc/go-diameter/test_app.xml", "test_app dictionary XML path")
	flag.Parse()

	data, err := os.ReadFile(*dictFile)
	if err != nil {
		log.Fatalf("read dict: %v", err)
	}
	if err := dict.Default.Load(bytes.NewReader(data)); err != nil {
		log.Fatalf("load dict: %v", err)
	}

	settings := common.Settings(f, "go-diameter-echo-server")
	mux := sm.New(settings)

	mux.HandleIdx(
		diam.CommandIndex{AppID: testAppID, Code: testCmdCode, Request: true},
		diam.HandlerFunc(handleTest(settings)))

	mux.HandleFunc("DPR", func(c diam.Conn, m *diam.Message) {
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		_, _ = a.WriteTo(c)
	})

	go common.LogErrors(mux.ErrorReports())

	log.Printf("echo server listening on %s (%s)", f.Addr, f.Network)
	if err := diam.ListenAndServeNetwork(f.Network, f.Addr, mux, nil); err != nil {
		log.Fatal(err)
	}
}

func handleTest(s *sm.Settings) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
		// Echo Test-AVP if present.
		if testAVP, err := m.FindAVP(testAVPCode, testVendorID); err == nil && testAVP != nil {
			a.AddAVP(diam.NewAVP(testAVPCode, avp.Mbit|avp.Vbit, testVendorID, testAVP.Data))
		}
		if _, err := a.WriteTo(c); err != nil {
			log.Printf("write answer: %v", err)
		}
	}
}
