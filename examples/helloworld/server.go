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
// http://tools.ietf.org/html/rfc6733#section-5.3.2
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
	// Add Grouped AVPs as received
	m.Add(diam.NewAVP(r.Dict.CodeFor("Vendor-Specific-Application-Id"), 0x00, 0x00, []*diam.AVP{
		diam.NewAVP(r.Dict.CodeFor("Acct-Application-Id"), 0x00, 0x00, uint32(16777299)),
	}))

	//Vendor-Specific-Application-Id AVP{Code=260,Flags=0x40,Length=44,VendorId=0x0,Padding=0,Grouped([Acct-Application-Id AVP{Code=259,Flags=0x40,Length=12,VendorId=0x0,Padding=0,uint32(16777299)}])}

	// Write response
	fmt.Println("Response:")
	m.PrettyPrint()
	w.Write(m)
}
