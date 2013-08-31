// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/fiorix/go-diameter/diam/avpdata"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/util"
)

type rfcHdr2 struct {
	Code     uint32
	Flags    uint8
	Length   [3]uint8
	VendorId uint32
}

func newAVP(msg *Message, code interface{}, flags uint8, vendor uint32, data avpdata.Generic) (*AVP, error) {
	davp, err := msg.Dict.FindAVP(msg.Header.ApplicationId, code)
	if err != nil {
		return nil, err
	}
	avp := &AVP{
		Code:     davp.Code,
		Flags:    flags,
		VendorId: vendor,
		msg:      msg,
	}
	if flags&0x80 > 0 {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr2{}))
	} else {
		avp.Length = uint32(unsafe.Sizeof(rfcHdr1{}))
	}
	// Set the body.
	if err = avp.set(davp, data); err != nil {
		return nil, err
	}
	return avp, nil
}

// NewAVP allocates and returns a new AVP.
// @code can be either the AVP code (int, uint32) or name (string).
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data avpdata.Generic) (*AVP, error) {
	avp, err := newAVP(m, code, flags, vendor, data)
	if err != nil {
		return nil, err
	}
	m.AVP = append(m.AVP, avp)
	return avp, nil
}

var ErrNoParentMessage = errors.New("This AVP has no parent message")

// Set updates the internal AVP data (the body) with a new value.
func (avp *AVP) Set(data interface{}) error {
	if avp.msg == nil {
		return ErrNoParentMessage
	}
	davp, err := avp.msg.Dict.FindAVP(avp.msg.Header.ApplicationId, avp.Code)
	if err != nil {
		return err
	}
	return avp.set(davp, data)
}

func (avp *AVP) set(davp *dict.AVP, data interface{}) error {
	var body Codec
	switch davp.Data.Type {
	case "OctetString":
		switch data.(type) {
		case *avpdata.OctetString:
			body = data.(*avpdata.OctetString)
		case []byte:
			body = &avpdata.OctetString{
				Value: string(data.([]byte))}
		case string:
			body = &avpdata.OctetString{Value: data.(string)}
		}
	case "Integer32":
		switch data.(type) {
		case *avpdata.Integer32:
			body = data.(*avpdata.Integer32)
		case int:
			body = &avpdata.Integer32{Value: int32(data.(int))}
		case int32:
			body = &avpdata.Integer32{Value: data.(int32)}
		case uint32:
			body = &avpdata.Integer32{Value: int32(data.(uint32))}
		}
	case "Integer64":
		switch data.(type) {
		case *avpdata.Integer64:
			body = data.(*avpdata.Integer64)
		case int:
			body = &avpdata.Integer64{Value: int64(data.(int))}
		case int64:
			body = &avpdata.Integer64{Value: data.(int64)}
		case uint64:
			body = &avpdata.Integer64{Value: int64(data.(uint64))}
		}
	case "Unsigned32":
		switch data.(type) {
		case *avpdata.Unsigned32:
			body = data.(*avpdata.Unsigned32)
		case int:
			body = &avpdata.Unsigned32{Value: uint32(data.(int))}
		case int32:
			body = &avpdata.Unsigned32{Value: uint32(data.(int32))}
		case uint32:
			body = &avpdata.Unsigned32{Value: data.(uint32)}
		}
	case "Unsigned64":
		switch data.(type) {
		case *avpdata.Unsigned64:
			body = data.(*avpdata.Unsigned64)
		case int:
			body = &avpdata.Unsigned64{Value: uint64(data.(int))}
		case int64:
			body = &avpdata.Unsigned64{Value: uint64(data.(int64))}
		case uint64:
			body = &avpdata.Unsigned64{Value: data.(uint64)}
		}
	case "Float32":
		switch data.(type) {
		case *avpdata.Float32:
			body = data.(*avpdata.Float32)
		case float32:
			body = &avpdata.Float32{Value: data.(float32)}
		}
	case "Float64":
		switch data.(type) {
		case *avpdata.Float64:
			body = data.(*avpdata.Float64)
		case float64:
			body = &avpdata.Float64{Value: data.(float64)}
		}
	case "Address":
		switch data.(type) {
		case *avpdata.Address:
			body = data.(*avpdata.Address)
		case net.IP:
			body = &avpdata.Address{IP: data.(net.IP)}
		case string:
			body = &avpdata.Address{IP: net.ParseIP(data.(string))}
		}
	case "IPv4": // To support Framed-IP-Address and alike.
		switch data.(type) {
		case *avpdata.IPv4:
			body = data.(*avpdata.IPv4)
		case net.IP:
			body = &avpdata.IPv4{IP: data.(net.IP)}
		case string:
			body = &avpdata.IPv4{IP: net.ParseIP(data.(string))}
		}
	case "Time":
		switch data.(type) {
		case *avpdata.Time:
			body = data.(*avpdata.Time)
		case int:
			body = &avpdata.Time{
				Value: time.Unix(int64(data.(int)), 0)}
		case int32:
			body = &avpdata.Time{
				Value: time.Unix(int64(data.(int32)), 0)}
		case int64:
			body = &avpdata.Time{
				Value: time.Unix(data.(int64), 0)}
		case time.Time:
			body = &avpdata.Time{
				Value: data.(time.Time)}
		}
	case "UTF8String":
		switch data.(type) {
		case *avpdata.UTF8String:
			body = data.(*avpdata.UTF8String)
		case []byte:
			body = &avpdata.UTF8String{avpdata.OctetString{
				Value: string(data.([]byte))}}
		case string:
			body = &avpdata.UTF8String{
				avpdata.OctetString{Value: data.(string)}}
		}
	case "DiameterIdentity":
		switch data.(type) {
		case *avpdata.DiameterIdentity:
			body = data.(*avpdata.DiameterIdentity)
		case []byte:
			body = &avpdata.DiameterIdentity{avpdata.OctetString{
				Value: string(data.([]byte))}}
		case string:
			body = &avpdata.DiameterIdentity{
				avpdata.OctetString{Value: data.(string)}}
		}
	case "DiameterURI":
		switch data.(type) {
		case *avpdata.DiameterURI:
			body = data.(*avpdata.DiameterURI)
		case []byte:
			body = &avpdata.DiameterURI{
				Value: string(data.([]byte))}
		case string:
			body = &avpdata.DiameterURI{Value: data.(string)}
		}
	case "Enumerated":
		switch data.(type) {
		case *avpdata.Enumerated:
			body = data.(*avpdata.Enumerated)
		case int:
			body = &avpdata.Enumerated{
				avpdata.Integer32{Value: int32(data.(int))}}
		}
	case "IPFilterRule":
		switch data.(type) {
		case *avpdata.IPFilterRule:
			body = data.(*avpdata.IPFilterRule)
		case []byte:
			body = &avpdata.IPFilterRule{avpdata.OctetString{
				Value: string(data.([]byte))}}
		case string:
			body = &avpdata.IPFilterRule{
				avpdata.OctetString{Value: data.(string)}}
		}
	case "Grouped":
		if gr, ok := data.(*Grouped); ok {
			body = gr
		}
	}
	if body == nil {
		fmt.Errorf("Unsupported data type: %s", data)
	}
	if avp.body != nil {
		avp.Length -= avp.body.Length()
	}
	avp.body = body
	avp.Length += body.Length()
	return nil
}

// Bytes returns an AVP in binary form so it can be attached to a Message
// before sent to a connection.
func (avp *AVP) Bytes() []byte {
	b := bytes.NewBuffer(nil)
	if avp.Flags&0x80 > 0 {
		hdr := rfcHdr2{
			Code:     avp.Code,
			Flags:    avp.Flags,
			Length:   util.Uint32To24(avp.Length),
			VendorId: avp.VendorId,
		}
		binary.Write(b, binary.BigEndian, hdr)
	} else {
		hdr := rfcHdr1{
			Code:   avp.Code,
			Flags:  avp.Flags,
			Length: util.Uint32To24(avp.Length),
		}
		binary.Write(b, binary.BigEndian, hdr)
	}
	binary.Write(b, binary.BigEndian, avp.body.Bytes())
	return b.Bytes()
}
