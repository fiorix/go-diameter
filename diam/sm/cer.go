// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
)

// handleCER handles Capabilities-Exchange-Request messages.
//
// If mandatory AVPs such as Origin-Host or Origin-Realm
// are missing, we close the connection.
//
// See RFC 6733 section 5.3 for details.
func handleCER(sm *StateMachine) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		ctx := c.Context()
		if _, ok := smpeer.FromContext(ctx); ok {
			// Ignore retransmission.
			return
		}
		cer := new(smparser.CER)
		// Decide whether an in-band upgrade is possible before answering,
		// so that a request this side cannot honour is refused with
		// DIAMETER_NO_COMMON_SECURITY rather than accepted and then failed.
		// An already-encrypted connection satisfies the request without an
		// upgrade.
		tlsActive := c.TLS() != nil
		upgrader, canUpgrade := c.(diam.TLSUpgrader)
		willUpgrade := !tlsActive && canUpgrade &&
			sm.cfg.TLSConfig != nil && upgrader.CanStartTLS()
		_, err := cer.ParseWithSecurity(m, smparser.Server, tlsActive || willUpgrade)
		if err != nil {
			err = errorCEA(sm, c, m, cer, err)
			if err != nil {
				sm.Error(&diam.ErrorReport{
					Conn:    c,
					Message: m,
					Error:   err,
				})
			}
			c.Close()
			return
		}
		// Only upgrade when the peer actually asked for TLS.
		upgrade := willUpgrade && cer.RequestedSecurity() == smparser.InbandSecurityTLS
		err = successCEA(sm, c, m, cer, upgrade)

		if err != nil {
			sm.Error(&diam.ErrorReport{
				Conn:    c,
				Message: m,
				Error:   err,
			})
			return
		}
		if upgrade {
			// RFC 6733 §13: the handshake begins once both ends reach the
			// open state, and if it fails both ends MUST move to the closed
			// state.
			if err := upgrader.StartTLS(sm.cfg.TLSConfig); err != nil {
				sm.Error(&diam.ErrorReport{
					Conn:    c,
					Message: m,
					Error:   fmt.Errorf("in-band TLS upgrade failed: %w", err),
				})
				c.Close()
				return
			}
		}
		meta := smpeer.FromCER(cer)
		c.SetContext(smpeer.NewContext(ctx, meta))
		// Notify about peer passing the handshake.
		select {
		case sm.hsNotifyc <- c:
		default:
		}
	}
}

// errorCEA sends an error answer indicating that the CER failed due to
// an unsupported (acct/auth) application, and includes the AVP that
// caused the failure in the message.
func errorCEA(sm *StateMachine, c diam.Conn, m *diam.Message, cer *smparser.CER, errMessage error) error {
	var (
		hostAddresses []datatype.Address
		err           error
	)
	if len(sm.cfg.HostIPAddresses) > 0 {
		hostAddresses = sm.cfg.HostIPAddresses
	} else {
		hostAddresses, err = getLocalAddresses(c)
		if err != nil {
			return fmt.Errorf("Error CEA '%s' create failure: %v", errMessage, err)
		}
	}

	var a *diam.Message
	switch errMessage {
	case smparser.ErrNoCommonSecurity:
		a = m.Answer(diam.NoCommonSecurity)
	case smparser.ErrNoCommonApplication:
		a = m.Answer(diam.NoCommonApplication)
	default:
		a = m.Answer(diam.UnableToComply)
	}
	// Fix for Same H2H and E2E Identifier in success response
	a.Header.HopByHopID = m.Header.HopByHopID
	a.Header.EndToEndID = m.Header.EndToEndID
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, sm.cfg.OriginHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, sm.cfg.OriginRealm)
	for _, hostAddress := range hostAddresses {
		a.NewAVP(avp.HostIPAddress, avp.Mbit, 0, hostAddress)
	}
	a.NewAVP(avp.VendorID, avp.Mbit, 0, sm.cfg.VendorID)
	a.NewAVP(avp.ProductName, 0, 0, sm.cfg.ProductName)
	if cer.OriginStateID != nil {
		a.AddAVP(cer.OriginStateID)
	}
	if sm.cfg.FirmwareRevision != 0 {
		a.NewAVP(avp.FirmwareRevision, 0, 0, sm.cfg.FirmwareRevision)
	}
	if sm.cfg.OnCEA != nil {
		sm.cfg.OnCEA(c, a)
	}
	_, err = a.WriteTo(c)
	if err != nil {
		err = fmt.Errorf("Error CEA '%s' send failure: %v", errMessage, err)
	}
	return err
}

// successCEA sends a success answer indicating that the CER was successfully
// parsed and accepted by the server. upgrade reports whether an in-band TLS
// handshake will follow, in which case the selected mechanism is echoed back
// so the peer knows to start it: RFC 6733 §13 requires the CER/CEA exchange
// to include Inband-Security-Id=TLS for TLS to be established in the open
// state.
func successCEA(sm *StateMachine, c diam.Conn, m *diam.Message, cer *smparser.CER, upgrade bool) error {
	var (
		hostAddresses []datatype.Address
		err           error
	)
	if len(sm.cfg.HostIPAddresses) > 0 {
		hostAddresses = sm.cfg.HostIPAddresses
	} else {
		hostAddresses, err = getLocalAddresses(c)
		if err != nil {
			return err
		}
	}

	a := m.Answer(diam.Success)
	// Fix for Same H2H and E2E Identifier in success response
	a.Header.HopByHopID = m.Header.HopByHopID
	a.Header.EndToEndID = m.Header.EndToEndID
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, sm.cfg.OriginHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, sm.cfg.OriginRealm)
	for _, hostAddress := range hostAddresses {
		a.NewAVP(avp.HostIPAddress, avp.Mbit, 0, hostAddress)
	}
	a.NewAVP(avp.VendorID, avp.Mbit, 0, sm.cfg.VendorID)
	a.NewAVP(avp.ProductName, 0, 0, sm.cfg.ProductName)
	if cer.OriginStateID != nil {
		a.AddAVP(cer.OriginStateID)
	}
	for _, app := range sm.supportedApps {
		var typ uint32
		switch app.AppType {
		case "auth":
			typ = avp.AuthApplicationID
		case "acct":
			typ = avp.AcctApplicationID
		}
		if app.Vendor != 0 {
			a.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(app.Vendor))
			a.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(app.Vendor)),
					diam.NewAVP(typ, avp.Mbit, 0, datatype.Unsigned32(app.ID)),
				},
			})
		} else {
			a.NewAVP(typ, avp.Mbit, 0, datatype.Unsigned32(app.ID))
		}
	}
	if sm.cfg.FirmwareRevision != 0 {
		a.NewAVP(avp.FirmwareRevision, 0, 0, sm.cfg.FirmwareRevision)
	}
	if upgrade {
		a.NewAVP(avp.InbandSecurityID, avp.Mbit, 0,
			datatype.Unsigned32(smparser.InbandSecurityTLS))
	}
	if sm.cfg.OnCEA != nil {
		sm.cfg.OnCEA(c, a)
	}
	_, err = a.WriteTo(c)
	return err
}
