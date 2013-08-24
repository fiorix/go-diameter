package main

import (
	"fmt"
	"net"

	"github.com/fiorix/go-diameter/base"
	diam "github.com/fiorix/go-diameter/stack"
)

func main() {
	diam.HandleFunc(257, OnCER)
	diam.ListenAndServe(":3868", nil, nil)
}

// On CER reply with CEA.
// http://tools.ietf.org/html/rfc6733#section-5.3.2
func OnCER(w diam.ResponseWriter, r *diam.Request) {
	fmt.Println("Request from", r.RemoteAddr, "to", r.LocalAddr, ":")
	r.Msg.PrettyPrint()
	m := base.NewMessage(
		257,  // Command Code
		0x00, // Flags
		0x00, // Application Id
		r.Msg.Header.HopByHopId,
		r.Msg.Header.EndToEndId,
		r.Msg.Dict,
	)
	// Build message by attaching multiple AVPs to it.
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Result-Code"), 0x40, 0x0,
		&base.Unsigned32{Value: 2001}))
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Origin-Host"), 0x40, 0x0,
		&base.DiameterIdentity{base.OctetString{Value: "go"}}))
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Origin-Realm"), 0x40, 0x0,
		&base.DiameterIdentity{base.OctetString{Value: "server"}}))
	localIP, _, _ := net.SplitHostPort(r.LocalAddr)
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Host-IP-Address"), 0x40, 0x0,
		&base.Address{Family: []byte("01"), IP: net.ParseIP(localIP)}))
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Vendor-Id"), 0x40, 0x0,
		&base.Unsigned32{Value: 131313}))
	m.Add(base.NewAVP(r.Dict, r.Dict.CodeFor("Product-Name"), 0x40, 0x0,
		&base.OctetString{Value: "go-diameter"}))
	// Reply with the same Origin-State-Id
	code := r.Dict.CodeFor("Origin-State-Id")
	OriginStateId, err := r.Msg.FindAVP(code)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}
	m.Add(base.NewAVP(r.Dict, code, 0x40, 0x0, OriginStateId.Body()))

	// Write response
	fmt.Println("Response:")
	m.PrettyPrint()
	w.Write(m)
}
