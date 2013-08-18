// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diameter

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Header struct {
	Version          uint8
	RawMessageLength [3]uint8
	CommandFlags     uint8
	RawCommandCode   [3]uint8
	ApplicationId    uint32
	HopByHopId       uint32
	EndToEndId       uint32
}

// MessageLength returns the RawMessageLength as int.
func (hdr *Header) MessageLength() uint32 {
	return uint24to32(hdr.RawMessageLength)
}

// CommandCode returns the RawCommandCode as int.
func (hdr *Header) CommandCode() uint32 {
	return uint24to32(hdr.RawCommandCode)
}

type Command struct {
	Name   string
	Abbrev string
}

var commandCodes = map[uint32]Command{
	274: {"Abbort-Session", "AS"},
	271: {"Accounting", "AC"},
	257: {"Capabilities-Exchange", "CE"},
	280: {"Device-Watchdog", "DW"},
	282: {"Disconnect-Peer", "DP"},
	258: {"Re-Auth", "RA"},
	275: {"Session-Termination", "ST"},
}

func CommandName(hdr *Header) (*Command, error) {
	var nameSuffix, abbrevSuffix string
	if hdr.CommandFlags&0x80 > 0 {
		nameSuffix = "-Request"
		abbrevSuffix = "R"
	} else {
		nameSuffix = "-Answer"
		abbrevSuffix = "A"
	}
	code := hdr.CommandCode()
	if cmd, ok := commandCodes[code]; ok {
		return &Command{
			cmd.Name + nameSuffix,
			cmd.Abbrev + abbrevSuffix,
		}, nil
	}
	return nil, fmt.Errorf("Unknown diameter command code: %d\n", code)
}

// ReadHeader reads one diameter header and return it.
func ReadHeader(r io.Reader) (*Header, error) {
	hdr := new(Header)
	if err := binary.Read(r, binary.BigEndian, hdr); err != nil {
		return nil, err
	}
	if hdr.Version != byte(1) {
		return nil,
			fmt.Errorf("Invalid diameter version %d", hdr.Version)
	}
	return hdr, nil
}

// Returns the diameter header in human readable format.
func (hdr *Header) String() string {
	rflag := hdr.CommandFlags&0x80 > 0
	pflag := hdr.CommandFlags&0x40 > 0
	eflag := hdr.CommandFlags&0x20 > 0
	tflag := hdr.CommandFlags&0x10 > 0
	cmd, err := CommandName(hdr)
	if err != nil {
		cmd = &Command{err.Error(), ""}
	}
	return fmt.Sprintf("DiameterHeader{"+
		"Version:%d, "+
		"MessageLength:%d, "+
		"CommandFlags{r=%v,p=%v,e=%v,t=%v}, "+
		"Command:%s (%s), "+
		"ApplicationId:%d, "+
		"HopByHopId:%#v, "+
		"EndToEndId:%#v}",
		hdr.Version,
		hdr.MessageLength(),
		rflag,
		pflag,
		eflag,
		tflag,
		cmd.Name,
		cmd.Abbrev,
		hdr.ApplicationId,
		hdr.HopByHopId,
		hdr.EndToEndId)
}
