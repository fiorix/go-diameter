// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

// handleCEA handles Capabilities-Exchange-Answer messages.
func handleCEA(sm *StateMachine, errc chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		cea := new(smparser.CEA)
		if err := cea.Parse(m, smparser.Client); err != nil {
			errc <- err
			return
		}
		// RFC 6733 §6.2: If we requested Inband-Security-Id=1 and
		// received a success CEA, upgrade to TLS on the client side.
		if sm.cfg.TLSConfig != nil && c.TLS() == nil {
			if u, ok := c.(diam.TLSUpgrader); ok {
				if err := u.StartTLSClient(sm.cfg.TLSConfig); err != nil {
					errc <- fmt.Errorf("post-CEA TLS upgrade failed: %w", err)
					return
				}
			}
		}
		meta := smpeer.FromCEA(cea)
		c.SetContext(smpeer.NewContext(c.Context(), meta))
		// Notify about peer passing the handshake.
		select {
		case sm.hsNotifyc <- c:
		default:
		}
		// Done receiving and validating this CEA.
		close(errc)
	}
}
