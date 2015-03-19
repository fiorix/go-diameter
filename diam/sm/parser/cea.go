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
}

// Parse parses and validates the given message.
func (cea *CEA) Parse(m *diam.Message) (err error) {
	if err = m.Unmarshal(cea); err != nil {
		return err
	}
	if err = cea.sanityCheck(); err != nil {
		return err
	}
	return nil
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
