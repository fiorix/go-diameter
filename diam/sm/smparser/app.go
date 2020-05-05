// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// Role stores information whether SM is initialized as a Client or a Server
type Role uint8

// ServerRole and ClientRole enums are passed to smparser for proper CER/CEA verification
const (
	Server Role = iota + 1
	Client
)

// Application validates accounting, auth, and vendor specific application IDs.
type Application struct {
	AcctApplicationID           []*diam.AVP
	AuthApplicationID           []*diam.AVP
	VendorSpecificApplicationID []*diam.AVP
	id                          []uint32 // List of supported application IDs.
}

// Parse ensures at least one common acct or auth applications in the CE
// exist in this server's dictionary.
func (app *Application) Parse(d *dict.Parser, localRole Role) (failedAVP *diam.AVP, err error) {
	failedAVP, err = app.validateAll(d, avp.AcctApplicationID, app.AcctApplicationID, localRole)
	oneFound := err == nil
	a, e := app.validateAll(d, avp.AuthApplicationID, app.AuthApplicationID, localRole)
	failedAVP, err, oneFound = chooseErr(failedAVP, err, oneFound, a, e)
	for _, vs := range app.VendorSpecificApplicationID {
		a, e = app.handleGroup(d, vs)
		failedAVP, err, oneFound = chooseErr(failedAVP, err, oneFound, a, e)
	}
	if !oneFound {
		return failedAVP, err
	}
	if app.ID() == nil {
		if localRole == Client {
			return failedAVP, ErrMissingApplication
		}
		return failedAVP, ErrNoCommonApplication

	}
	return nil, nil
}

func chooseErr(curAVP *diam.AVP, curErr error, oneFound bool, newAvp *diam.AVP, newErr error) (*diam.AVP, error, bool) {
	if newErr == nil {
		return curAVP, curErr, true
	}
	if (curErr == nil) || (curAVP == nil && newAvp != nil) {
		return newAvp, newErr, oneFound
	}
	return curAVP, curErr, oneFound
}

// handleGroup handles the VendorSpecificApplicationID grouped AVP and
// validates accounting or auth applications.
func (app *Application) handleGroup(d *dict.Parser, gavp *diam.AVP) (failedAVP *diam.AVP, err error) {
	group, ok := gavp.Data.(*diam.GroupedAVP)
	if !ok {
		return gavp, &ErrUnexpectedAVP{gavp}
	}
	var success bool
	for _, a := range group.AVP {
		switch a.Code {
		case avp.AcctApplicationID:
			failedAVP, err = app.validate(d, a.Code, a)
		case avp.AuthApplicationID:
			failedAVP, err = app.validate(d, a.Code, a)
		}
		success = success || (err == nil)
	}
	if success {
		return nil, nil
	}
	return failedAVP, err
}

// validateAll is a convenience method to test a slice of application IDs.
// according to https://tools.ietf.org/html/rfc6733#page-60:
//   A receiver of a Capabilities-Exchange-Request (CER) message that does
//   not have any applications in common with the sender MUST return a
//   Capabilities-Exchange-Answer (CEA) with the Result-Code AVP set to
//   DIAMETER_NO_COMMON_APPLICATION and SHOULD disconnect the transport
//   layer connection.
// so, we need to find at least one App ID in common
func (app *Application) validateAll(
	d *dict.Parser, appType uint32, appAVPs []*diam.AVP, localRole Role) (failedAVP *diam.AVP, err error) {

	if len(appAVPs) > 0 {
		var oneFound bool
		for _, a := range appAVPs {
			a, e := app.validate(d, appType, a)
			failedAVP, err, oneFound = chooseErr(failedAVP, err, oneFound, a, e)
		}
		if oneFound {
			return nil, nil
		}
		return
	}
	if localRole == Client {
		return nil, ErrMissingApplication
	}
	return nil, ErrNoCommonApplication
}

// validate ensures the given acct or auth application ID exists in
// the given dictionary.
func (app *Application) validate(d *dict.Parser, appType uint32, appAVP *diam.AVP) (failedAVP *diam.AVP, err error) {
	if appAVP == nil {
		return nil, nil
	}
	var typ string
	switch appType {
	case avp.AcctApplicationID:
		typ = "acct"
	case avp.AuthApplicationID:
		typ = "auth"
	}
	if appAVP.Code != appType {
		return appAVP, &ErrUnexpectedAVP{appAVP}
	}
	appID, ok := appAVP.Data.(datatype.Unsigned32)
	if !ok {
		return appAVP, &ErrUnexpectedAVP{appAVP}
	}
	id := uint32(appID)
	if id == 0xffffffff { // relay application id
		app.id = append(app.id, id)
		return nil, nil
	}
	_, err = d.App(id, typ)
	if err != nil {
		return appAVP, ErrNoCommonApplication
	}
	app.id = append(app.id, id)
	return nil, nil
}

// ID returns a list of supported application IDs.
// Must be called after Parse, otherwise it returns an empty array.
func (app *Application) ID() []uint32 {
	return app.id
}
