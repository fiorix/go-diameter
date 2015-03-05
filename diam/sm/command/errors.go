// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package command

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingOriginHost is returned by Parse when
	// the message does not contain an Origin-Host AVP.
	ErrMissingOriginHost = errors.New("missing Origin-Host")

	// ErrMissingOriginRealm is returned by Parse when
	// the message does not contain an Origin-Realm AVP.
	ErrMissingOriginRealm = errors.New("missing Origin-Realm")

	// ErrMissingOriginStateID is returned by Parse when
	// the message does not contain an Origin-State-Id AVP.
	ErrMissingOriginStateID = errors.New("missing Origin-State-Id")

	// ErrMissingApplication is returned by Parse when
	// the CER does not contain any Acct-Application-Id or
	// Auth-Application-Id, or their embedded versions in
	// the Vendor-Specific-Application-Id AVP.
	ErrMissingApplication = errors.New("missing application")

	// ErrNoCommonSecurity is returned by Parse when
	// the CER contains the Inband-Security-Id.
	// We currently don't support that.
	ErrNoCommonSecurity = errors.New("no common security")
)

// ErrNoCommonApplication is returned by Parse when the
// application IDs in the CER don't match the applications
// defined in our dictionary.
type ErrNoCommonApplication struct {
	ID   uint32
	Type string
}

// Error implements the error interface.
func (e *ErrNoCommonApplication) Error() string {
	return fmt.Sprintf("%s application %d is not supported", e.Type, e.ID)
}
