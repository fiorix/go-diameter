// diam-smoke-server: accepts CER, answers DWR via the state machine,
// handles DPR cleanly, and logs inbound commands.
package main

import (
	"flag"
	"log"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/sm"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

func main() {
	f := common.RegisterFlags(":3868")
	flag.Parse()

	settings := common.Settings(f, "go-diameter-smoke-server")
	mux := sm.New(settings)

	mux.HandleFunc("DPR", func(c diam.Conn, m *diam.Message) {
		log.Println("DPR received; sending DPA")
		a := m.Answer(diam.Success)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		if _, err := a.WriteTo(c); err != nil {
			log.Printf("DPA write: %v", err)
		}
	})

	mux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		log.Printf("unexpected message from %s:\n%s", c.RemoteAddr(), m)
	})

	go common.LogErrors(mux.ErrorReports())

	log.Printf("listening on %s (%s)", f.Addr, f.Network)
	if err := diam.ListenAndServeNetwork(f.Network, f.Addr, mux, nil); err != nil {
		log.Fatal(err)
	}
}
