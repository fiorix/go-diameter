// diam-smoke-client: connects, handshakes CER/CEA, optionally drives
// DWR via the built-in watchdog, then sends DPR and exits 0.
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/sm"

	"github.com/fiorix/go-diameter/v4/tests/go/internal/common"
)

func main() {
	f := common.RegisterFlags("127.0.0.1:3868")
	sendDPR := flag.Bool("dpr", true, "send Disconnect-Peer-Request before exit")
	flag.Parse()

	settings := common.Settings(f, "go-diameter-smoke-client")
	mux := sm.New(settings)

	go common.LogErrors(mux.ErrorReports())

	cli := common.Client(mux, f, nil, nil)
	conn, err := cli.DialNetwork(f.Network, f.Addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	log.Println("CER/CEA handshake OK")

	// Hold the connection open for the requested duration so DWRs fire.
	if f.Watchdog > 0 && f.Timeout > 0 {
		time.Sleep(f.Timeout)
	}

	if *sendDPR {
		if err := common.SendDPR(conn, settings); err != nil {
			log.Printf("send DPR: %v", err)
			os.Exit(1)
		}
		log.Println("DPR sent")
		time.Sleep(500 * time.Millisecond)
	}
	conn.Close()
}
