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

// handleCEA handles the Capabilities-Exchange-Answer to the client's CER.
// errc reports a handshake failure to Client.handshake, or is closed once the
// CEA has been accepted.
func handleCEA(sm *StateMachine, errc chan error) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		cea := new(smparser.CEA)
		if err := cea.Parse(m, smparser.Client); err != nil {
			errc <- err
			return
		}
		// RFC 6733 §13: a TLS handshake in the open state follows only a
		// CER/CEA exchange that both carry Inband-Security-Id=TLS. A CEA
		// that omits the AVP, or answers NO_INBAND_SECURITY (the default
		// when omitted, §6.10), has declined.
		requested := c.TLS() == nil &&
			requestedInbandSecurity(c.Context()) == smparser.InbandSecurityTLS
		if requested && sm.cfg.TLSConfig == nil &&
			cea.InbandSecurityID == smparser.InbandSecurityTLS {
			// Without TLSConfig the AVP is only advertised. A peer that
			// selected TLS will now start the handshake, so the connection
			// cannot continue in cleartext.
			errc <- fmt.Errorf("peer selected in-band TLS but Settings.TLSConfig is nil")
			return
		}
		if requested && sm.cfg.TLSConfig != nil {
			if cea.InbandSecurityID != smparser.InbandSecurityTLS {
				// Continuing in cleartext would silently give up the
				// protection this client asked for.
				errc <- fmt.Errorf(
					"peer declined the requested in-band TLS: Inband-Security-Id=%d in CEA",
					cea.InbandSecurityID)
				return
			}
			u, ok := c.(diam.TLSUpgrader)
			if !ok || !u.CanStartTLS() {
				errc <- fmt.Errorf("peer selected in-band TLS but this connection cannot be upgraded")
				return
			}
			if err := u.StartTLSClient(sm.cfg.TLSConfig); err != nil {
				errc <- fmt.Errorf("in-band TLS upgrade failed: %w", err)
				return
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
