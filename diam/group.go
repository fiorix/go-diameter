// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"fmt"

	"github.com/fiorix/go-diameter/diam/avp/format"
	"github.com/fiorix/go-diameter/diam/dict"
)

const GroupedAVPFormat = 50 // Must not conflict with other format.DataFormatId.

// GroupedAVP that is different from the dummy format.Grouped.
type GroupedAVP struct {
	AVP []*AVP
}

func DecodeGrouped(data format.Grouped, application uint32, dictionary *dict.Parser) (*GroupedAVP, error) {
	g := new(GroupedAVP)
	b := []byte(data)
	for n := 0; n < len(b); {
		avp, err := DecodeAVP(b[n:], application, dictionary)
		if err != nil {
			return nil, err
		}
		g.AVP = append(g.AVP, avp)
		n += avp.Len()
	}
	// TODO: handle nested groups?
	return g, nil
}

func (g *GroupedAVP) Serialize() []byte {
	b := make([]byte, g.Len())
	var n int
	for _, a := range g.AVP {
		a.SerializeTo(b[n:])
		n += a.Len()
	}
	return b
}

func (g *GroupedAVP) Len() int {
	var l int
	for _, a := range g.AVP {
		l += a.Len()
	}
	return l
}

func (g *GroupedAVP) Padding() int {
	return 0
}

func (g *GroupedAVP) Format() format.FormatId {
	return GroupedAVPFormat
}

func (g *GroupedAVP) String() string {
	var b bytes.Buffer
	for n, a := range g.AVP {
		if n > 0 {
			fmt.Fprint(&b, ",")
		}
		fmt.Fprint(&b, a)
	}
	return b.String()
}
