// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/diam/datatypes"
	"github.com/fiorix/go-diameter/diam/dict"
)

const GroupedType = 50 // Must not conflict with other datatypes.DataTypeId.

// Grouped AVP.  This is different from the dummy datatypes.Grouped.
type Grouped struct {
	AVP []*AVP
}

func DecodeGrouped(a *AVP, application uint32, dictionary *dict.Parser) (*Grouped, error) {
	if a.Data.Type() != datatypes.GroupedType {
		return nil, errors.New("Invalid AVP is not Grouped")
	}
	g := &Grouped{}
	b := []byte(a.Data.(datatypes.Grouped))
	for n := 0; n < cap(b); {
		avp, err := decodeAVP(b[n:], application, dictionary)
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

func (g *Grouped) Type() datatypes.DataTypeId {
	return GroupedType
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
