// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

import (
	"bytes"
	"fmt"

	"github.com/fiorix/go-diameter/diam/avp/format"
	"github.com/fiorix/go-diameter/diam/dict"
)

const GroupedFormat = 50 // Must not conflict with other format.DataFormatId.

// Grouped AVP.  This is different from the dummy format.Grouped.
type Grouped struct {
	AVP []*AVP
}

func DecodeGrouped(data format.Grouped, application uint32, dictionary *dict.Parser) (*Grouped, error) {
	g := new(Grouped)
	b := []byte(data)
	for n := 0; n < len(b); {
		avp, err := Decode(b[n:], application, dictionary)
		if err != nil {
			return nil, err
		}
		g.AVP = append(g.AVP, avp)
		n += avp.Len()
	}
	// TODO: handle nested groups?
	return g, nil
}

func (g *Grouped) Serialize() []byte {
	b := make([]byte, g.Len())
	var n int
	for _, a := range g.AVP {
		a.SerializeTo(b[n:])
		n += a.Len()
	}
	return b
}

func (g *Grouped) Len() int {
	var l int
	for _, a := range g.AVP {
		l += a.Len()
	}
	return l
}

func (g *Grouped) Padding() int {
	return 0
}

func (g *Grouped) Format() format.FormatId {
	return GroupedFormat
}

func (g *Grouped) String() string {
	var b bytes.Buffer
	for n, a := range g.AVP {
		if n > 0 {
			fmt.Fprint(&b, ",")
		}
		fmt.Fprint(&b, a)
	}
	return b.String()
}
