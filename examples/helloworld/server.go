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
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Result-Code"), 0x40, 0x0,
		&base.Unsigned32{Value: 2001}, nil,
	))
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Origin-Host"), 0x40, 0x0,
		&base.DiameterIdentity{base.OctetString{Value: "go"}}, nil,
	))
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Origin-Realm"), 0x40, 0x0,
		&base.DiameterIdentity{base.OctetString{Value: "server"}}, nil,
	))
	localIP, _, _ := net.SplitHostPort(r.LocalAddr)
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Host-IP-Address"), 0x40, 0x0,
		&base.Address{Family: []byte("01"), IP: net.ParseIP(localIP)}, nil,
	))
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Vendor-Id"), 0x40, 0x0,
		&base.Unsigned32{Value: 131313}, nil,
	))
	m.Add(base.NewAVP(
		r.Dict.CodeFor("Product-Name"), 0x40, 0x0,
		&base.OctetString{Value: "go-diameter"}, nil,
	))
	// Reply with the same Origin-State-Id
	code := r.Dict.CodeFor("Origin-State-Id")
	OriginStateId := r.Msg.FindAVP(code)
	m.Add(base.NewAVP(code, 0x40, 0x0, OriginStateId.Data, nil))

	// Write response
	fmt.Println("Response:")
	m.PrettyPrint()
	w.Write(m)
}
