package middleware

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/deb-2025/go-diameter/diam"
	"github.com/deb-2025/go-diameter/diam/avp"
	"github.com/deb-2025/go-diameter/diam/datatype"
	"github.com/deb-2025/go-diameter/diam/diamtest"
	"github.com/opentracing/opentracing-go"
)

func TestTracer(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", NewTracer(newCERHandler(errc)))

	srv := diamtest.NewServer(smux, nil)
	defer srv.Close()

	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.HandleIdx(diam.CommandIndex{AppID: 0, Code: diam.CapabilitiesExchange, Request: false}, handleCEA(errc, wait))

	cli, err := diam.Dial(srv.Addr, cmux, nil)
	if err != nil {
		t.Fatal(err)
	}

	sendCER(cli)

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

func TestTracerFunc(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", TracerFunc(handleCER(errc, false)))

	srv := diamtest.NewServer(smux, nil)
	defer srv.Close()

	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.HandleIdx(diam.CommandIndex{AppID: 0, Code: diam.CapabilitiesExchange, Request: false}, handleCEA(errc, wait))

	cli, err := diam.Dial(srv.Addr, cmux, nil)
	if err != nil {
		t.Fatal(err)
	}

	sendCER(cli)

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case err := <-smux.ErrorReports():
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

// copy from diam/server_test.go
func sendCER(w io.Writer) (n int64, err error) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("cli"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1234))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(1))
	return m.WriteTo(w)
}

type CER struct {
	OriginHost        string    `avp:"Origin-Host"`
	OriginRealm       string    `avp:"Origin-Realm"`
	VendorID          int       `avp:"Vendor-Id"`
	ProductName       string    `avp:"Product-Name"`
	OriginStateID     *diam.AVP `avp:"Origin-State-Id"`
	AcctApplicationID *diam.AVP `avp:"Acct-Application-Id"`
}

type cerHandler struct {
	errc chan error
}

func newCERHandler(errc chan error) diam.Handler {
	return &cerHandler{errc: errc}
}

func (h *cerHandler) ServeDIAM(c diam.Conn, m *diam.Message) {
	span := opentracing.SpanFromContext(m.Context())
	if span == nil {
		h.errc <- fmt.Errorf("Span is nil")
	}
	var req CER
	err := m.Unmarshal(&req)
	if err != nil {
		h.errc <- err
		return
	}
	a := m.Answer(diam.Success)
	_, err = sendCEA(c, a, req.OriginStateID, req.AcctApplicationID)
	if err != nil {
		h.errc <- err
	}
	c.(diam.CloseNotifier).CloseNotify()
	go func() {
		<-c.(diam.CloseNotifier).CloseNotify()
	}()
}

func handleCER(errc chan error, useTLS bool) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		span := opentracing.SpanFromContext(m.Context())
		if span == nil {
			errc <- fmt.Errorf("Span is nil")
		}

		var req CER
		err := m.Unmarshal(&req)
		if err != nil {
			errc <- err
			return
		}
		a := m.Answer(diam.Success)
		_, err = sendCEA(c, a, req.OriginStateID, req.AcctApplicationID)
		if err != nil {
			errc <- err
		}
		c.(diam.CloseNotifier).CloseNotify()
		go func() {
			<-c.(diam.CloseNotifier).CloseNotify()
		}()
	}
}

func sendCEA(w io.Writer, m *diam.Message, OriginStateID, AcctApplicationID *diam.AVP) (n int64, err error) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.AddAVP(OriginStateID)
	m.AddAVP(AcctApplicationID)
	return m.WriteTo(w)
}

func handleCEA(errc chan error, wait chan struct{}) diam.HandlerFunc {
	type CEA struct {
		OriginHost        string `avp:"Origin-Host"`
		OriginRealm       string `avp:"Origin-Realm"`
		VendorID          int    `avp:"Vendor-Id"`
		ProductName       string `avp:"Product-Name"`
		OriginStateID     int    `avp:"Origin-State-Id"`
		AcctApplicationID int    `avp:"Acct-Application-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {

		var resp CEA
		err := m.Unmarshal(&resp)
		if err != nil {
			errc <- err
			return
		}
		// Initialize & start close notifier
		closeNotifyChan := c.(diam.CloseNotifier).CloseNotify()
		// Wait on close notify chan outside of main serve loop, closeNotifier routine is started by
		// liveSwitchReader.Read to avoid io.Pipe deadlock issue
		go func() {
			<-closeNotifyChan // wait on c.Close to complete
			select {          // close only if not already closed
			case <-wait:
			default:
				close(wait)
			}
		}()
		c.Close()
	}
}
