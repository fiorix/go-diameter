// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter message, header with multiple AVPs.  Part of go-diameter.

package base

import (
	"fmt"

	"github.com/fiorix/go-diameter/dict"
)

// Message is a diameter message, composed of a header and multiple AVPs.
type Message struct {
	Header *Header
	AVP    []*AVP
	Dict   *dict.Parser // Dictionary associated with this Message
}

// Find is a helper that returns an AVP by looking up its code, or nil if
// the AVP is not in the message.
func (m *Message) FindAVP(code uint32) *AVP {
	for _, a := range m.AVP {
		if code == a.Code {
			return a
		}
	}
	return nil
}

// PrettyPrint prints messages in a human readable format.
func (m *Message) PrettyPrint() {
	fmt.Println(m.Header)
	for _, avp := range m.AVP {
		fmt.Printf("  %s\n", avp)
	}
	fmt.Println()
}
