// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// AVP Data conversions between Diameter and Go.  Part of go-diameter.
// Based on database/sql types.

package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// OctetString Diameter Type.
type OctetString struct {
	Value   string
	Padding uint32 // Extra bytes to make the String a multiple of 4 octets
}

// Data implements the Data interface.
func (os *OctetString) Data() Data {
	return os.Value
}

// Put implements the Codec interface. It updates internal String and Padding.
func (os *OctetString) Put(d Data) {
	b := d.([]byte)
	l := uint32(len(b))
	os.Padding = pad4(l) - l
	os.Value = string(b)
}

// Bytes implement the Codec interface. Padding is always recalculated from
// the internal String.
func (os *OctetString) Bytes() []byte {
	os.updatePadding() // Do this every time? Geez.
	l := uint32(len(os.Value))
	b := make([]byte, l+os.Padding)
	copy(b, os.Value)
	return b
}

// Length implements the Codec interface. Returns length without padding.
func (os *OctetString) Length() uint32 {
	return uint32(len(os.Value)) - os.Padding
}

// update internal padding value.
func (os *OctetString) updatePadding() {
	if os.Padding == 0 {
		l := uint32(len(os.Value))
		os.Padding = pad4(l) - l
	}
}

// String returns a human readable version of the AVP.
func (os *OctetString) String() string {
	os.updatePadding() // Update padding
	return fmt.Sprintf("OctetString{Value:'%s',Padding:%d}",
		os.Value, os.Padding)
}

// Time Diameter Type.
type Time struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *Time) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf("Time{Value:'%s',Padding:%d}", p.Value, p.Padding)
}

// UTF8String Diameter Type.
type UTF8String struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *UTF8String) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf("UTF8String{Value:'%s',Padding:%d}",
		p.Value, p.Padding)
}

// DiameterIdentity Diameter Type.
type DiameterIdentity struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *DiameterIdentity) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf("DiameterIdentity{Value:'%s',Padding:%d}",
		p.Value, p.Padding)
}

// IPFilterRule Diameter Type.
type IPFilterRule struct {
	OctetString
}

// String returns a human readable version of the AVP.
func (p *IPFilterRule) String() string {
	p.updatePadding() // Update padding
	return fmt.Sprintf("IPFilterRule{Value:'%s',Padding:%d}",
		p.Value, p.Padding)
}

// Address Diameter Type.
type Address struct {
	// Address Family (e.g. AF_INET=1)
	// http://www.iana.org/assignments/address-family-numbers/address-family-numbers.xhtml
	Family []byte

	// Parsed IP address
	IP net.IP

	// Padding to 4 octets
	Padding int
}

// Data implements the Data interface.
func (addr *Address) Data() Data {
	return addr.IP
}

// Put implements the Coded interface. It updates internal Family and IP.
func (addr *Address) Put(d Data) {
	b := d.([]byte)
	// TODO: Support IPv6
	if len(b) >= 6 && b[1] == 1 { // AF_INET=1 IPv4 only.
		addr.Family = []byte{b[0], b[1]}
		addr.IP = net.IPv4(b[2], b[3], b[4], b[5])
		addr.Padding = 2
	}
}

// Bytes implement the Codec interface.
func (addr *Address) Bytes() []byte {
	if ip := addr.IP.To4(); ip != nil {
		// IPv4 always need 2 byte padding. (derived from OctetString)
		addr.Padding = 2 // TODO: Fix this
		b := []byte{
			addr.Family[0],
			addr.Family[1],
			ip[0],
			ip[1],
			ip[2],
			ip[3],
			0, // Padding
			0, // Padding
		}
		return b
	}
	return []byte{}
}

// Length implements the Codec interface. Returns length without padding.
func (addr *Address) Length() uint32 {
	if addr.IP.To4() != nil {
		return 6 // TODO: Fix this
	}
	return 0
}

// String returns a human readable version of the AVP.
func (addr *Address) String() string {
	addr.Bytes() // Update family and padding
	return fmt.Sprintf("Address{Value:'%s',Padding:%d}",
		addr.IP.String(), addr.Padding)
}

// DiameterURI Diameter Type.
type DiameterURI struct {
	Value string
}

// Data implements the Data interface.
func (du *DiameterURI) Data() Data {
	return du.Value
}

// Put implements the Codec interface.
func (du *DiameterURI) Put(d Data) {
	du.Value = string(d.([]byte))
}

// Bytes implement the Codec interface.
func (du *DiameterURI) Bytes() []byte {
	return []byte(du.Value)
}

// Length implements the Codec interface.
func (du *DiameterURI) Length() uint32 {
	return uint32(len(du.Value))
}

// String returns a human readable version of the AVP.
func (du *DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{Value:'%s'}", du.Value)
}

// Integer32 Diameter Type
type Integer32 struct {
	Value  int32
	Buffer []byte
}

// Data *implements the Data interface.
func (n Integer32) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Int32.
func (n *Integer32) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Int32 and stored on Buffer.
func (n *Integer32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Integer32) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Integer32) String() string {
	return fmt.Sprintf("Integer32{Value:%d}", n.Value)
}

// Integer64 Diameter Type
type Integer64 struct {
	Value  int64
	Buffer []byte
}

// Data implements the Data interface.
func (n *Integer64) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Int64.
func (n *Integer64) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Int64 and stored on Buffer.
func (n *Integer64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Integer64) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Integer64) String() string {
	return fmt.Sprintf("Integer64{Value:%d}", n.Value)
}

// Unsigned32 Diameter Type
type Unsigned32 struct {
	Value  uint32
	Buffer []byte
}

// Data implements the Data interface.
func (n *Unsigned32) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Uint32.
func (n *Unsigned32) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Uint32 and stored on Buffer.
func (n *Unsigned32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Unsigned32) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Unsigned32) String() string {
	return fmt.Sprintf("Unsigned32{Value:%d}", n.Value)
}

// Unsigned64 Diameter Type
type Unsigned64 struct {
	Value  uint64
	Buffer []byte
}

// Data implements the Data interface.
func (n *Unsigned64) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Uint64.
func (n *Unsigned64) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Uint64 and stored on Buffer.
func (n *Unsigned64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Unsigned64) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Unsigned64) String() string {
	return fmt.Sprintf("Unsigned64{Value:%d}", n.Value)
}

// Float32 Diameter Type
type Float32 struct {
	Value  float32
	Buffer []byte
}

// Data implements the Data interface.
func (n *Float32) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Float32.
func (n *Float32) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float32 and stored on Buffer.
func (n *Float32) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Float32) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Float32) String() string {
	return fmt.Sprintf("Float32{Value:%d}", n.Value)
}

// Float64 Diameter Type
type Float64 struct {
	Value  float64
	Buffer []byte
}

// Data implements the Data interface.
func (n *Float64) Data() Data {
	return n.Value
}

// Put implements the Codec interface. It updates internal Buffer and Float64.
func (n *Float64) Put(d Data) {
	n.Buffer = d.([]byte)
	binary.Read(bytes.NewBuffer(n.Buffer), binary.BigEndian, &n.Value)
}

// Bytes implement the Codec interface. Bytes are always rewritten from
// the internal Float64 and stored on Buffer.
func (n *Float64) Bytes() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, n.Value)
	n.Buffer = b.Bytes()
	return n.Buffer
}

// Length implements the Codec interface.
func (n *Float64) Length() uint32 {
	if n.Buffer == nil {
		n.Bytes()
	}
	return uint32(len(n.Buffer))
}

// String returns a human readable version of the AVP.
func (n *Float64) String() string {
	return fmt.Sprintf("Float64{Value:%d}", n.Value)
}

// Enumerated Diameter Type
type Enumerated struct {
	Integer32
}

// String returns a human readable version of the AVP.
func (p *Enumerated) String() string {
	return fmt.Sprintf("Enumerated{Value:%d}", p.Value)
}

// Grouped Diameter Type
type Grouped struct {
	AVP    []*AVP
	Buffer []byte // len(Buffer) might be bigger than length below due to padding.
	length uint32 // Aggregate length of all AVPs without padding.
}

// Data implements the Data interface.
func (gr *Grouped) Data() Data {
	return gr.AVP
}

// Put implements the Codec interface. It updates internal Buffer and Length.
// Takes an AVP as input.
func (gr *Grouped) Put(d Data) {
	avp := d.(*AVP)
	gr.length += avp.Length
	gr.AVP = append(gr.AVP, avp)
	gr.Buffer = bytes.Join([][]byte{gr.Buffer, avp.Bytes()}, []byte{})
}

// Bytes implement the Codec interface. Bytes are always returned from
// internal Buffer cache.
func (gr *Grouped) Bytes() []byte {
	return gr.Buffer
}

// Length implements the Codec interface.
func (gr *Grouped) Length() uint32 {
	return gr.length
}

// String returns a human readable version of the AVP.
func (gr *Grouped) String() string {
	s := make([]string, len(gr.AVP))
	for n, avp := range gr.AVP {
		s[n] = avp.String()
	}
	return fmt.Sprintf("Grouped{%s}", strings.Join(s, ","))
}
