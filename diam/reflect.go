// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"errors"
	"reflect"
)

// Unmarshal stores the result of a diameter message in the struct
// pointed to by dst.
//
// Unmarshal can not only decode AVPs into the struct, but also their
// Go equivalent data types, directly.
//
// For example:
//
//	type CER struct {
//		OriginHost  AVP    `avp:"Origin-Host"`
//		.. or
//		OriginHost  *AVP   `avp:"Origin-Host"`
//		.. or
//		OriginHost  string `avp:"Origin-Host"`
//	}
//	var d CER
//	err := diam.Unmarshal(&d)
//
// This decodes the Origin-Host AVP as three different types. The first, AVP,
// makes a copy of the AVP in the message and stores in the struct. The
// second, *AVP, stores a pointer to the original AVP in the message. If you
// change the values of it, you're actually changing the message.
// The third decodes the inner contents of AVP.Data, which in this case is
// a diamtype.DiameterIdentity, and stores the value of it in the struct.
//
// Unmarshal supports all the basic Go types, including slices, for multiple
// AVPs of the same type) and structs, for grouped AVPs.
//
// Slices:
//
//	type CER struct {
//		Vendors  []*AVP `avp:"Supported-Vendor-Id"`
//	}
//	var d CER
//	err := diam.Unmarshal(&d)
//
// Slices have the same principles of other types. If they're of type
// []*AVP it'll store references in the struct, while []AVP makes
// copies and []int (or []string, etc) decodes the AVP data for you.
//
// Grouped AVPs:
//
//	type VSA struct {
//		AuthAppId int `avp:"Auth-Application-Id"`
//		VendorId  int `avp:"Vendor-Id"`
//	}
//	type CER struct {
//		VSA VSA  `avp:"Vendor-Specific-Application-Id"`
//		.. or
//		VSA *VSA `avp:"Vendor-Specific-Application-Id"`
//		.. or
//		VSA struct {
//			AuthAppId int `avp:"Auth-Application-Id"`
//			VendorId  int `avp:"Vendor-Id"`
//		} `avp:"Vendor-Specific-Application-Id"`
//	}
//	var d CER
//	err := m.Unmarshal(&d)
//
// Other types are supported as well, such as net.IP and time.Time where
// applicable. See the diamtype sub-module for details. Usually, you want
// to decode values to their native Go type when the AVPs don't have to be
// re-used in an answer, such as Origin-Host and friends. The ones that are
// usually added to responses, such as Origin-State-Id are better decoded to
// just AVP or *AVP, making it easier to re-use them in the answer.
//
// Note that decoding values to *AVP is much faster and more efficient than
// decoding to AVP or the native Go types.
func (m *Message) Unmarshal(dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return errors.New("dst is not a pointer to struct")
	}
	return scanStruct(m, m.AVP, v)
}

func scanStruct(m *Message, avps []*AVP, field reflect.Value) error {
	base := reflect.Indirect(field)
	if base.Kind() != reflect.Struct {
		return errors.New("dst is not a pointer to struct")
	}
	idx := newIndex(avps)
	for n := 0; n < base.NumField(); n++ {
		f := base.Field(n)
		bf := base.Type().Field(n)
		avpname := bf.Tag.Get("avp")
		if avpname == "" || avpname == "-" {
			continue
		}
		// Lookup the AVP name (tag) in the dictionary.
		// The dictionary AVP has the code.
		d, err := m.Dictionary.FindAVP(m.Header.ApplicationId, avpname)
		if err != nil {
			return err
		}
		// See if this AVP exist in the message.
		avps, exists := idx[d.Code]
		if !exists {
			continue
		}
		//log.Println("Handling", f, bf)
		unmarshal(m, f, avps)
	}
	return nil
}

func newIndex(avps []*AVP) map[uint32][]*AVP {
	idx := make(map[uint32][]*AVP)
	for _, avp := range avps {
		idx[avp.Code] = append(idx[avp.Code], avp)
	}
	return idx
}

func unmarshal(m *Message, f reflect.Value, avps []*AVP) {
	ft := f.Type()
	// Case 1: f is *AVP
	at := reflect.TypeOf(avps[0])
	if ft.AssignableTo(at) {
		f.Set(reflect.ValueOf(avps[0]))
		return
	}
	// Case 2: f is AVP
	at = reflect.TypeOf(*avps[0])
	if ft.AssignableTo(at) {
		f.Set(reflect.ValueOf(*avps[0]))
		return
	}
	// Case 3: f is the type of AVP.Data
	dv := reflect.ValueOf(avps[0].Data)
	if dv.Type().ConvertibleTo(ft) {
		f.Set(dv.Convert(ft))
		return
	}
	if dv.Type() == _groupedType {
		if f.Kind() == reflect.Ptr {
			nf := reflect.New(f.Type().Elem())
			scanStruct(m, avps[0].Data.(*Grouped).AVP, nf)
			f.Set(nf)
		} else {
			scanStruct(m, avps[0].Data.(*Grouped).AVP, f)
		}
		return
	}
	// Look for slices, except []byte (e.g. net.IP, time.Time)
	if f.Kind() == reflect.Slice && !ft.AssignableTo(_byteSliceType) {
		unmarshalSlice(f, avps)
		return
	}
}

func unmarshalSlice(f reflect.Value, avps []*AVP) {
	ft := f.Type()
	// Case 1: f is []*AVP
	at := reflect.TypeOf(avps)
	if at.AssignableTo(ft) {
		f.Set(reflect.ValueOf(avps))
		return
	}
	// Case 2: f is []AVP
	eft := f.Type().Elem()
	at = reflect.TypeOf(*avps[0])
	if at.AssignableTo(eft) {
		tmp := make([]AVP, 0, len(avps))
		for _, avp := range avps {
			tmp = append(tmp, *avp)
		}
		f.Set(reflect.ValueOf(tmp))
		return
	}
	// Case 3: f[n] is the type of AVP.Data
	tmp := reflect.MakeSlice(ft, 0, len(avps))
	for _, avp := range avps {
		// sanity check each avp
		v := reflect.ValueOf(avp.Data)
		if v.Type().ConvertibleTo(eft) {
			tmp = reflect.Append(tmp, v.Convert(eft))
		}
	}
	if tmp.Len() > 0 {
		f.Set(tmp)
	}
}

var (
	_groupedType   = reflect.TypeOf(&Grouped{})
	_byteSliceType = reflect.TypeOf([]byte{})
)
