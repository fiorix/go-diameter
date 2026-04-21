package diam

import (
	"net"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func benchmarkDiamTransport(b *testing.B, network string) {
	done := make(chan struct{}, 1)
	mux := NewServeMux()
	mux.HandleFunc("ALL", func(c Conn, m *Message) {
		done <- struct{}{}
	})

	ln, err := MultistreamListen(network, "127.0.0.1:0")
	if err != nil {
		b.Skipf("cannot listen on %s: %v", network, err)
	}
	defer ln.Close()

	srv := &Server{Handler: mux}
	go srv.Serve(ln)

	// Use the same dialer the library uses internally
	dialer := getMultistreamDialer(network, 0, nil)
	rwc, err := dialer.Dial(network, ln.Addr().String())
	if err != nil {
		b.Skipf("cannot dial %s: %v", network, err)
	}
	defer rwc.Close()

	msg := NewRequest(257, 0, dict.Default)
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("bench.host"))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("bench.realm"))
	msg.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	msg.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	msg.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("bench"))
	payload, _ := msg.Serialize()

	b.SetBytes(int64(len(payload)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rwc.Write(payload)
		<-done
	}
}

func BenchmarkDiamTransport_TCP(b *testing.B) {
	benchmarkDiamTransport(b, "tcp")
}

func BenchmarkDiamTransport_SCTP(b *testing.B) {
	benchmarkDiamTransport(b, "sctp")
}
