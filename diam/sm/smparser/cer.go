// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

// CER is a Capabilities-Exchange-Request message.
// See RFC 6733 section 5.3.1 for details.
type CER struct {
	OriginHost                  datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm                 datatype.DiameterIdentity `avp:"Origin-Realm"`
	OriginStateID               *diam.AVP                 `avp:"Origin-State-Id"`
	InbandSecurityID            *diam.AVP                 `avp:"Inband-Security-Id"`
	AcctApplicationID           []*diam.AVP               `avp:"Acct-Application-Id"`
	AuthApplicationID           []*diam.AVP               `avp:"Auth-Application-Id"`
	VendorSpecificApplicationID []*diam.AVP               `avp:"Vendor-Specific-Application-Id"`
	appID                       []uint32                  // List of supported application IDs.
	requestedSecurity           uint32                    // Negotiated Inband-Security-Id value.
}

// RequestedSecurity returns the Inband-Security-Id value from the CER
// (0=NO_INBAND_SECURITY, 1=TLS). Valid after a successful Parse.
func (cer *CER) RequestedSecurity() uint32 {
	return cer.requestedSecurity
}

// Parse parses and validates the given message, and returns nil when
// all AVPs are ok, and all accounting or authentication applications
// in the CER match the applications in our dictionary. If one or more
// mandatory AVPs are missing, it returns a nil failedAVP and a proper
// error. If all mandatory AVPs are present but no common application
// is found, then it returns the failedAVP (with the application that
// we don't support in our dictionary) and an error. Another cause
// for error is the presence of Inband Security when TLS is not active.
func (cer *CER) Parse(m *diam.Message, localRole Role) (failedAVP *diam.AVP, err error) {
	return cer.ParseWithSecurity(m, localRole, false)
}

// ParseWithSecurity is like Parse but accepts a tlsActive flag. When
// tlsActive is true, Inband-Security-Id=1 (TLS) is accepted because
// the transport is already secured — per RFC 6733 §5.3.1 the peer is
// simply declaring its TLS capability which is already satisfied.
//
// When tlsActive is false and tlsConfig is available on the server,
// Inband-Security-Id=1 signals that a TLS upgrade should occur after
// CEA is sent (RFC 6733 §6.2). The caller should check
// RequestedSecurity() after a successful parse.
func (cer *CER) ParseWithSecurity(m *diam.Message, localRole Role, tlsActive bool) (failedAVP *diam.AVP, err error) {
	if err = m.Unmarshal(cer); err != nil {
		return nil, err
	}
	if err = cer.sanityCheck(); err != nil {
		return nil, err
	}
	if cer.InbandSecurityID != nil {
		v := uint32(cer.InbandSecurityID.Data.(datatype.Unsigned32))
		cer.requestedSecurity = v
		if v != 0 && !tlsActive {
			return nil, ErrNoCommonSecurity
		}
	}
	app := &Application{
		AcctApplicationID:           cer.AcctApplicationID,
		AuthApplicationID:           cer.AuthApplicationID,
		VendorSpecificApplicationID: cer.VendorSpecificApplicationID,
	}
	if failedAVP, err = app.Parse(m.Dictionary(), localRole); err != nil {
		return failedAVP, err
	}
	cer.appID = app.ID()
	return nil, nil
}

// sanityCheck ensures mandatory AVPs are present.
func (cer *CER) sanityCheck() error {
	if len(cer.OriginHost) == 0 {
		return ErrMissingOriginHost
	}
	if len(cer.OriginRealm) == 0 {
		return ErrMissingOriginRealm
	}
	return nil
}

// Applications return a list of supported Application IDs.
func (cer *CER) Applications() []uint32 {
	return cer.appID
}
