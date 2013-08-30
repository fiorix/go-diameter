// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP parser.  Part of go-diameter.

package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/fiorix/go-diameter/dict"
)

type rfcHdr2 struct {
	Code     uint32
	Flags    uint8
	Length   [3]uint8
	VendorId uint32
}

func newAVP(appid uint32, d *dict.Parser, code interface{}, flags uint8, vendor uint32, data Data) (*AVP, error) {
	davp, err := d.FindAVP(appid, code)
	if err != nil {
		return nil, err
	}
	var body Codec
	switch davp.Data.Type {
	case "OctetString":
		switch data.(type) {
		case *OctetString:
			body = data.(*OctetString)
		case []byte:
			body = &OctetString{Value: string(data.([]byte))}
		case string:
			body = &OctetString{Value: data.(string)}
		}
	case "Integer32":
		switch data.(type) {
		case *Integer32:
			body = data.(*Integer32)
		case int:
			body = &Integer32{Value: int32(data.(int))}
		case int32:
			body = &Integer32{Value: data.(int32)}
		case uint32:
			body = &Integer32{Value: int32(data.(uint32))}
		}
	case "Integer64":
		switch data.(type) {
		case *Integer64:
			body = data.(*Integer64)
		case int:
			body = &Integer64{Value: int64(data.(int))}
		case int64:
			body = &Integer64{Value: data.(int64)}
		case uint64:
			body = &Integer64{Value: int64(data.(uint64))}
		}
	case "Unsigned32":
		switch data.(type) {
		case *Unsigned32:
			body = data.(*Unsigned32)
		case int:
			body = &Unsigned32{Value: uint32(data.(int))}
		case int32:
			body = &Unsigned32{Value: uint32(data.(int32))}
		case uint32:
			body = &Unsigned32{Value: data.(uint32)}
		}
	case "Unsigned64":
		switch data.(type) {
		case *Unsigned64:
			body = data.(*Unsigned64)
		case int:
			body = &Unsigned64{Value: uint64(data.(int))}
		case int64:
			body = &Unsigned64{Value: uint64(data.(int64))}
		case uint64:
			body = &Unsigned64{Value: data.(uint64)}
		}
	case "Float32":
		switch data.(type) {
		case *Float32:
			body = data.(*Float32)
		case float32:
			body = &Float32{Value: data.(float32)}
		}
	case "Float64":
		switch data.(type) {
		case *Float64:
			body = data.(*Float64)
		case float64:
			body = &Float64{Value: data.(float64)}
		}
	case "Address":
		switch data.(type) {
		case *Address:
			body = data.(*Address)
		case net.IP:
			body = &Address{IP: data.(net.IP)}
		case string:
			body = &Address{IP: net.ParseIP(data.(string))}
		}
	case "IPv4": // To support Framed-IP-Address and alike.
		switch data.(type) {
		case *IPv4:
			body = data.(*IPv4)
		case net.IP:
			body = &IPv4{IP: data.(net.IP)}
		case string:
			body = &IPv4{IP: net.ParseIP(data.(string))}
		}
	case "Time":
		switch data.(type) {
		case *Time:
			body = data.(*Time)
		case int:
			body = &Time{Value: time.Unix(int64(data.(int)), 0)}
		case int32:
			body = &Time{Value: time.Unix(int64(data.(int32)), 0)}
		case int64:
			body = &Time{Value: time.Unix(data.(int64), 0)}
		case time.Time:
			body = &Time{Value: data.(time.Time)}
		}
	case "UTF8String":
		switch data.(type) {
		case *UTF8String:
			body = data.(*UTF8String)
		case []byte:
			body = &UTF8String{
				OctetString{Value: string(data.([]byte))}}
		case string:
			body = &UTF8String{
				OctetString{Value: data.(string)}}
		}
	case "DiameterIdentity":
		switch data.(type) {
		case *DiameterIdentity:
			body = data.(*DiameterIdentity)
		case []byte:
			body = &DiameterIdentity{
				OctetString{Value: string(data.([]byte))}}
		case string:
			body = &DiameterIdentity{
				OctetString{Value: data.(string)}}
		}
	case "DiameterURI":
		switch data.(type) {
		case *DiameterURI:
			body = data.(*DiameterURI)
		case []byte:
			body = &DiameterURI{Value: string(data.([]byte))}
		case string:
			body = &DiameterURI{Value: data.(string)}
		}
	case "Enumerated":
		switch data.(type) {
		case *Enumerated:
			body = data.(*Enumerated)
		case int:
			body = &Enumerated{Integer32{Value: int32(data.(int))}}
		}
	case "IPFilterRule":
		switch data.(type) {
		case *IPFilterRule:
			body = data.(*IPFilterRule)
		case []byte:
			body = &IPFilterRule{
				OctetString{Value: string(data.([]byte))}}
		case string:
			body = &IPFilterRule{
				OctetString{Value: data.(string)}}
		}
	case "Grouped":
		if gr, ok := data.(*Grouped); ok {
			body = gr
		}
	}
	if body == nil {
		return nil, fmt.Errorf("Unsupported data type: %s", data)
	}
	avp := &AVP{
		Code:     davp.Code,
		Flags:    flags,
		VendorId: vendor,
		body:     body,
		dict:     d,
	}
	if flags&0x80 > 0 {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr2{}))
	} else {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr1{}))
	}
	avp.Length += body.Length()
	return avp, nil
}

// NewAVP allocates and returns a new AVP.
// @code can be either the AVP code (int, uint32) or name (string).
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data Data) (*AVP, error) {
	avp, err := newAVP(
		m.Header.ApplicationId,
		m.Dict,
		code,
		flags,
		vendor,
		data,
	)
	if err != nil {
		return nil, err
	}
	m.AVP = append(m.AVP, avp)
	return avp, nil
}

// Set updates the internal AVP data (the body) with a new value.
func (avp *AVP) Set(body Codec) {
	if avp.body != nil {
		avp.Length -= avp.body.Length()
	}
	avp.body = body
	avp.Length += body.Length()
}

// Bytes returns an AVP in binary form so it can be attached to a Message
// before sent to a connection.
func (avp *AVP) Bytes() []byte {
	b := bytes.NewBuffer(nil)
	if avp.Flags&0x80 > 0 {
		hdr := rfcHdr2{
			Code:     avp.Code,
			Flags:    avp.Flags,
			Length:   uint32To24(avp.Length),
			VendorId: avp.VendorId,
		}
		binary.Write(b, binary.BigEndian, hdr)
	} else {
		hdr := rfcHdr1{
			Code:   avp.Code,
			Flags:  avp.Flags,
			Length: uint32To24(avp.Length),
		}
		binary.Write(b, binary.BigEndian, hdr)
	}
	binary.Write(b, binary.BigEndian, avp.body.Bytes())
	return b.Bytes()
}
