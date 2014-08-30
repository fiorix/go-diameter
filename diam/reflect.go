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
// a format.DiameterIdentity, and stores the value of it in the struct.
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
// applicable. See the format sub-package for details. Usually, you want
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
	return scanStruct(m, v, m.AVP)
}

// newIndex returns a map of AVPs indexed by their code.
// TODO: make this part of the Message.
func newIndex(avps []*AVP) map[uint32][]*AVP {
	idx := make(map[uint32][]*AVP, len(avps))
	for _, a := range avps {
		idx[a.Code] = append(idx[a.Code], a)
	}
	return idx
}

func scanStruct(m *Message, field reflect.Value, avps []*AVP) error {
	base := reflect.Indirect(field)
	if base.Kind() != reflect.Struct {
		return errors.New("dst is not a pointer to struct")
	}
	idx := newIndex(avps)
	for n := 0; n < base.NumField(); n++ {
		f := base.Field(n)
		bt := base.Type().Field(n)
		avpname := bt.Tag.Get("avp")
		if len(avpname) == 0 || avpname == "-" {
			continue
		}
		// Lookup the AVP name (tag) in the dictionary.
		// The dictionary AVP has the code.
		d, err := m.dictionary.FindAVP(m.Header.ApplicationId, avpname)
		if err != nil {
			return err
		}
		// See if this AVP exist in the message.
		avps, exists := idx[d.Code]
		if !exists {
			continue
		}
		//log.Println("Handling", f, bt)
		unmarshal(m, f, avps)
	}
	return nil
}

func unmarshal(m *Message, f reflect.Value, avps []*AVP) {
	fieldType := f.Type()

	switch f.Kind() {
	case reflect.Slice:
		// Copy byte arrays.
		dv := reflect.ValueOf(avps[0].Data)
		if dv.Type().ConvertibleTo(fieldType) {
			f.Set(dv.Convert(fieldType))
			break
		}

		// Allocate new slice and copy all items.
		f.Set(reflect.MakeSlice(fieldType, len(avps), len(avps)))
		// TODO: optimize?
		for n := 0; n < len(avps); n++ {
			unmarshal(m, f.Index(n), avps[n:])
		}

	case reflect.Interface, reflect.Ptr:
		if f.IsNil() {
			f.Set(reflect.New(fieldType.Elem()))
		}
		unmarshal(m, f.Elem(), avps)

	case reflect.Struct:
		// Test for *AVP
		at := reflect.TypeOf(avps[0])
		if fieldType.AssignableTo(at) {
			f.Set(reflect.ValueOf(avps[0]))
			break
		}

		// Test for AVP
		at = reflect.TypeOf(*avps[0])
		if fieldType.ConvertibleTo(at) {
			f.Set(reflect.ValueOf(*avps[0]))
			break
		}

		// Handle grouped AVPs.
		if group, ok := avps[0].Data.(*GroupedAVP); ok {
			scanStruct(m, f, group.AVP)
		}

	default:
		// Test for AVP.Data (e.g. format.UTF8String, string)
		dv := reflect.ValueOf(avps[0].Data)
		if dv.Type().ConvertibleTo(fieldType) {
			f.Set(dv.Convert(fieldType))
		}
	}
}
