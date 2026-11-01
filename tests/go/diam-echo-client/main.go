// diam-echo-client: CER, then N Test-Requests to exercise the FD
// test_app.fdx echo server. Exits 0 if all answers return Result-Code
// 2001, non-zero otherwise.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
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
	testVendorID = 999999
	testAppID    = 123456
	testCmdCode  = 234567
	testAVPCode  = 345678
)

func main() {
	f := common.RegisterFlags("127.0.0.1:3868")
	dictFile := flag.String("dict", "/etc/go-diameter/test_app.xml", "test_app dictionary XML path")
	flag.Parse()

	if err := loadDict(*dictFile); err != nil {
		log.Fatalf("load dict: %v", err)
	}

	settings := common.Settings(f, "go-diameter-echo-client")
	mux := sm.New(settings)

	var ok, nok int64
	done := make(chan struct{}, f.N)
	mux.HandleIdx(
		diam.CommandIndex{AppID: testAppID, Code: testCmdCode, Request: false},
		diam.HandlerFunc(func(c diam.Conn, m *diam.Message) {
			rc, err := m.FindAVP(avp.ResultCode, 0)
			if err != nil || rc == nil {
				atomic.AddInt64(&nok, 1)
				log.Printf("answer missing Result-Code")
			} else if u, _ := rc.Data.(datatype.Unsigned32); u == diam.Success {
				atomic.AddInt64(&ok, 1)
			} else {
				atomic.AddInt64(&nok, 1)
				log.Printf("answer Result-Code=%v", rc.Data)
			}
			done <- struct{}{}
		}))

	go common.LogErrors(mux.ErrorReports())

	// FD test_app advertises its app id standalone in CEA (Auth-Application-Id
	// without Vendor-Specific-Application-Id), so match that on the CER side.
	cli := common.Client(mux, f, []uint32{testAppID}, nil)

	conn, err := cli.DialNetwork(f.Network, f.Addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	log.Println("handshake OK, sending", f.N, "Test-Requests")

	meta, _ := smpeer.FromContext(conn.Context())
	for i := 0; i < f.N; i++ {
		if err := sendTest(conn, settings, meta, i); err != nil {
			log.Fatalf("send Test #%d: %v", i, err)
		}
	}

	deadline := time.After(f.Timeout)
	for got := 0; got < f.N; got++ {
		select {
		case <-done:
		case <-deadline:
			log.Fatalf("timeout after %d/%d answers (ok=%d nok=%d)",
				got, f.N, atomic.LoadInt64(&ok), atomic.LoadInt64(&nok))
		}
	}
	log.Printf("done: ok=%d nok=%d", ok, nok)

	_ = common.SendDPR(conn, settings)
	time.Sleep(200 * time.Millisecond)
	conn.Close()

	if nok > 0 || ok != int64(f.N) {
		os.Exit(1)
	}
}

func sendTest(c diam.Conn, s *sm.Settings, meta *smpeer.Metadata, i int) error {
	m := diam.NewRequest(testCmdCode, testAppID, dict.Default)
	sid := fmt.Sprintf("session;%d;%d", time.Now().Unix(), i)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, s.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
	if meta != nil {
		m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
		m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	}
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("test-"+strconv.Itoa(i)))
	// test_app.fdx registers its own Test-AVP as INTEGER32; match that
	// to avoid 5014 INVALID_AVP_LENGTH on the echo path.
	m.NewAVP(testAVPCode, avp.Mbit|avp.Vbit, testVendorID, datatype.Unsigned32(i))
	_, err := m.WriteTo(c)
	return err
}

func loadDict(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return dict.Default.Load(bytes.NewReader(data))
}
