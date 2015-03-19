// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package parser

import (
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
)

// CEA is a Capabilities-Exchange-Answer message.
// See RFC 6733 section 5.3.2 for details.
type CEA struct {
	ResultCode    uint32                    `avp:"Result-Code"`
	OriginHost    datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm   datatype.DiameterIdentity `avp:"Origin-Realm"`
	OriginStateID uint32                    `avp:"Origin-State-Id"`
	AcctAppID     []*diam.AVP               `avp:"Acct-Application-Id"`
	AuthAppID     []*diam.AVP               `avp:"Auth-Application-Id"`
}

// Parse parses and validates the given message, and returns nil when
// all AVPs are ok, and all accounting or authentication applications
// in the CEA match the applications in our dictionary. If one or more
// mandatory AVPs are missing, it returns a nil failedAVP and a proper
// error. If all mandatory AVPs are present but no common application
// is found, then it returns the failedAVP (with the application that
// we don't support in our dictionary) and an error.
func (cea *CEA) Parse(m *diam.Message) (failedAVP *diam.AVP, err error) {
	if err = m.Unmarshal(cea); err != nil {
		return nil, err
	}
	if err = cea.sanityCheck(); err != nil {
		return nil, err
	}
	return nil, nil
}

// sanityCheck ensures mandatory AVPs are present.
func (cea *CEA) sanityCheck() error {
	if cea.ResultCode == 0 {
		return ErrMissingResultCode
	}
	if len(cea.OriginHost) == 0 {
		return ErrMissingOriginHost
	}
	if len(cea.OriginRealm) == 0 {
		return ErrMissingOriginRealm
	}
	if cea.OriginStateID == 0 {
		return ErrMissingOriginStateID
	}
	return nil
}
