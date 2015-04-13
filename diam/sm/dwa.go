// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/sm/parser"
)

// handleDWA handles Device-Watchdog-Answer messages.
func handleDWA(sm *StateMachine, osidc chan uint32) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		dwa := new(parser.DWA)
		if err := dwa.Parse(m); err != nil {
			sm.Error(&diam.ErrorReport{
				Conn:    c,
				Message: m,
				Error:   err,
			})
			return
		}
		if dwa.ResultCode != diam.Success {
			return
		}
		select {
		case osidc <- dwa.OriginStateID:
		default:
		}
	}
}
