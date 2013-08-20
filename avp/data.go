// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP Data conversions between Diameter and Go.
// Based on database/sql types.

package avp

import (
	"bytes"
	"encoding/binary"
)

// OctetString Diameter Type.
type OctetString struct {
	String  string
	Padding int // Extra bytes to make the String a multiple of 4 octets
}

// Data implements the Data interface.
func (os OctetString) Data() Data {
	return os.String
}

// Put implements the Codec interface. It updates internal String and Padding.
func (os *OctetString) Put(b []byte) {
	l := len(b)
	os.Padding = pad4(l) - l
	os.String = string(b)
}

// Bytes implement the Codec interface. Padding is always recalculated from
// the internal String.
func (os *OctetString) Bytes() []byte {
	if os.Padding == 0 { // Recalculate every time?
		l := len(os.String)
		os.Padding = pad4(l) - l
	}
	b := make([]byte, len(os.String)+os.Padding)
	copy(b, os.String)
	return b
}

// Address Diameter Type.
type Address OctetString

// Time Diameter Type.
type Time OctetString

// UTF8String Diameter Type.
type UTF8String OctetString

// DiameterIdentity Diameter Type.
type DiameterIdentity OctetString

// DiameterURI Diameter Type.
type DiameterURI struct {
	String string
}

// Data implements the Data interface.
func (du *DiameterURI) Data() Data {
	return du.String
}

// Put implements the Codec interface.
func (du *DiameterURI) Put(b []byte) {
	du.String = string(b)
}

// Bytes implement the Codec interface.
func (du *DiameterURI) Bytes() []byte {
	return []byte(du.String)
}

// Integer32 Diameter Type
type Integer32 struct {
	Int32  int32
	Buffer []byte
}

// Data implements the Data interface.
func (n Integer32) Data() Data {
	return n.Int32
}

// Put implements the Codec interface. It updates internal Buffer and Int32.
func (n *Integer32) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Int32)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Int32 and stored on Buffer.
func (n *Integer32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Int32)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Integer64 Diameter Type
type Integer64 struct {
	Int64  int64
	Buffer []byte
}

// Data implements the Data interface.
func (n Integer64) Data() Data {
	return n.Int64
}

// Put implements the Codec interface. It updates internal Buffer and Int64.
func (n *Integer64) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Int64)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Int64 and stored on Buffer.
func (n *Integer64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Int64)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Unsigned32 Diameter Type
type Unsigned32 struct {
	Uint32 uint32
	Buffer []byte
}

// Data implements the Data interface.
func (n Unsigned32) Data() Data {
	return n.Uint32
}

// Put implements the Codec interface. It updates internal Buffer and Uint32.
func (n *Unsigned32) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Uint32)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Uint32 and stored on Buffer.
func (n *Unsigned32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Uint32)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Unsigned64 Diameter Type
type Unsigned64 struct {
	Uint64 uint64
	Buffer []byte
}

// Data implements the Data interface.
func (n Unsigned64) Data() Data {
	return n.Uint64
}

// Put implements the Codec interface. It updates internal Buffer and Uint64.
func (n *Unsigned64) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Uint64)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Uint64 and stored on Buffer.
func (n *Unsigned64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Uint64)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Float32 Diameter Type
type Float32 struct {
	Float32 float32
	Buffer  []byte
}

// Data implements the Data interface.
func (n Float32) Data() Data {
	return n.Float32
}

// Put implements the Codec interface. It updates internal Buffer and Float32.
func (n *Float32) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Float32)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float32 and stored on Buffer.
func (n *Float32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Float32)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Float64 Diameter Type
type Float64 struct {
	Float64 float64
	Buffer  []byte
}

// Data implements the Data interface.
func (n Float64) Data() Data {
	return n.Float64
}

// Put implements the Codec interface. It updates internal Buffer and Float64.
func (n *Float64) Put(b []byte) {
	n.Buffer = b
	binary.Read(bytes.NewBuffer(b), binary.BigEndian, &n.Float64)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float64 and stored on Buffer.
func (n *Float64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Float64)
	n.Buffer = b.Bytes()
	return n.Buffer
}
