// diam-bad-client: deliberately malformed message sender for error-path
// scenarios. Supports:
//   -mode omit-origin-host : valid CER handshake, then CCR without
//                            Origin-Host (expect 5005 DIAMETER_MISSING_AVP)
//   -mode bad-avp-length   : valid CER handshake, then raw write of a
//                            CCR header with a truncated AVP length
//                            (expect 5014 DIAMETER_INVALID_AVP_LENGTH)
package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"log"
	"math/rand"
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

var mode = flag.String("mode", "omit-origin-host", "omit-origin-host | bad-avp-length")

func main() {
	f := common.RegisterFlags("127.0.0.1:3868")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	settings := common.Settings(f, "go-diameter-bad-client")
	mux := sm.New(settings)

	got := make(chan uint32, 4)
	mux.HandleIdx(
		diam.CommandIndex{AppID: 4, Code: diam.CreditControl, Request: false},
		diam.HandlerFunc(func(_ diam.Conn, m *diam.Message) {
			rc, _ := m.FindAVP(avp.ResultCode, 0)
			if rc != nil {
				if u, ok := rc.Data.(datatype.Unsigned32); ok {
					got <- uint32(u)
					return
				}
			}
			got <- 0
		}))

	go common.LogErrors(mux.ErrorReports())

	cli := common.Client(mux, f, []uint32{4}, nil)
	cli.AuthApplicationID = []*diam.AVP{
		diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
	}

	conn, err := cli.DialNetwork(f.Network, f.Addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	meta, _ := smpeer.FromContext(conn.Context())

	switch *mode {
	case "omit-origin-host":
		if err := sendCCRNoOrigin(conn, settings, meta); err != nil {
			log.Fatalf("send: %v", err)
		}
	case "bad-avp-length":
		if err := sendBadAVP(conn); err != nil {
			log.Fatalf("send: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}

	select {
	case rc := <-got:
		log.Printf("got Result-Code=%d", rc)
		if err := expected(rc); err != nil {
			log.Fatal(err)
		}
	case <-time.After(f.Timeout):
		// Connection drop is also an acceptable error response.
		log.Println("no answer (connection closed by peer)")
	}

	conn.Close()
}

func expected(rc uint32) error {
	switch *mode {
	case "omit-origin-host":
		if rc == 5005 {
			return nil
		}
	case "bad-avp-length":
		if rc == 5014 {
			return nil
		}
	}
	return errors.New("unexpected Result-Code " + strconv.Itoa(int(rc)))
}

func sendCCRNoOrigin(c diam.Conn, s *sm.Settings, meta *smpeer.Metadata) error {
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("session;bad;1"))
	// Deliberately omit Origin-Host.
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, s.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("32251@3gpp.org"))
	_, err := m.WriteTo(c)
	return err
}

// sendBadAVP writes a minimal CCR frame with one AVP whose declared
// length is longer than the frame permits. Bypasses the marshaller.
func sendBadAVP(c diam.Conn) error {
	// Diameter header: Version(1)=1, MsgLen(3)=36, Flags(1)=0x80, CmdCode(3)=272,
	// App-Id(4)=4, HBH(4)=rand, E2E(4)=rand, + one bad AVP (12 bytes).
	buf := make([]byte, 36)
	buf[0] = 1
	// MsgLen = 36
	binary.BigEndian.PutUint32(buf[0:4], (1<<24)|36)
	buf[4] = 0x80
	binary.BigEndian.PutUint32(buf[4:8], (uint32(0x80)<<24)|272)
	binary.BigEndian.PutUint32(buf[8:12], 4)
	binary.BigEndian.PutUint32(buf[12:16], rand.Uint32())
	binary.BigEndian.PutUint32(buf[16:20], rand.Uint32())
	// AVP: code=264 (Origin-Host), flags=M, length=999 (> remaining bytes)
	binary.BigEndian.PutUint32(buf[20:24], 264)
	// flags (1) + length (3); len=999 is impossible here
	buf[24] = 0x40 // M bit
	buf[25] = 0x00
	buf[26] = 0x03
	buf[27] = 0xe7 // 999
	// payload bytes (4) — irrelevant, the length will trip validation
	copy(buf[28:32], []byte("xxxx"))
	binary.BigEndian.PutUint32(buf[32:36], 0)

	_, err := c.Write(buf)
	return err
}
