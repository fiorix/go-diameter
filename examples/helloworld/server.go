package main

import (
	"fmt"
	"net"

	"github.com/fiorix/go-diameter/diam"
)

func main() {
	diam.HandleFunc(257, OnCER)
	diam.ListenAndServe(":3868", nil, nil)
}

// On CER reply with CEA.
// http://tools.ietf.org/html/rfc3588#section-5.3.2
func OnCER(w diam.ResponseWriter, r *diam.Request) {
	fmt.Println("Request from", r.RemoteAddr, "to", r.LocalAddr, ":")
	host, _, _ := net.SplitHostPort(r.LocalAddr)
	r.Msg.PrettyPrint()
	m := diam.NewMessage(
		257,  // Command Code
		0x00, // Flags
		0x00, // Application Id
		r.Msg.Header.HopByHopId,
		r.Msg.Header.EndToEndId,
	)
	m.Add(diam.NewAVP(r.Dict.CodeFor("Result-Code"), 0x40, 0x0, uint32(2001)))
	m.Add(diam.NewAVP(r.Dict.CodeFor("Origin-Host"), 0x40, 0x0, "go"))
	m.Add(diam.NewAVP(r.Dict.CodeFor("Origin-Realm"), 0x40, 0x0, "server"))
	//m.Add(diam.NewAVP(r.Dict.CodeFor("Host-IP-Address"), 0x40, 0x0, net.ParseIP("10.0.0.1")))
	m.Add(diam.NewAVP(r.Dict.CodeFor("Host-IP-Address"), 0x40, 0x0, net.ParseIP(host)))
	m.Add(diam.NewAVP(r.Dict.CodeFor("Vendor-Id"), 0x40, 0x0, uint32(131313)))
	m.Add(diam.NewAVP(r.Dict.CodeFor("Product-Name"), 0x40, 0x0, "go-diameter"))
	// Reply with the same Origin-State-Id
	code := r.Dict.CodeFor("Origin-State-Id")
	osAVP := r.Msg.Find(code)
	m.Add(diam.NewAVP(code, 0x40, 0x0, osAVP.Data))
	// Write response
	fmt.Println("Response:")
	m.PrettyPrint()
	w.Write(m)
}
