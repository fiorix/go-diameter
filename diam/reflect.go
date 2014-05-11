// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"errors"
	"reflect"
)

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
	if dv.Type() == reflect.TypeOf(&Grouped{}) {
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
	if f.Kind() == reflect.Slice && !ft.AssignableTo(byteArrayType) {
		unmarshalSlice(f, avps)
		return
	}
}

var byteArrayType = reflect.TypeOf([]byte{})

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
