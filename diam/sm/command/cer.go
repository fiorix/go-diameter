package command

import (
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
)

// CER is a Capabilities-Exchange-Request message.
// For details see section 5.3.1 of RFC 6733.
type CER struct {
	OriginHost          datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm         datatype.DiameterIdentity `avp:"Origin-Realm"`
	OriginStateID       *diam.AVP                 `avp:"Origin-State-Id"`
	InbandSecurityID    *diam.AVP                 `avp:"Inband-Security-Id"`
	AcctAppID           []*diam.AVP               `avp:"Acct-Application-Id"`
	AuthAppID           []*diam.AVP               `avp:"Auth-Application-Id"`
	VendorSpecificAppID []struct {
		AcctAppID *diam.AVP `avp:"Acct-Application-Id"`
		AuthAppID *diam.AVP `avp:"Auth-Application-Id"`
	} `avp:"Vendor-Specific-Application-Id"`

	appID []uint32 // List of supported application IDs.
}

// Parse parses and validates the given message, and returns nil when
// all AVPs are ok, and all accounting or authentication applications
// in the CER match the applications in our dictionary. If one or more
// mandatory AVPs are missing, it returns a nil failedAVP and a proper
// error. If all mandatory AVPs are present but no common application
// is found, then it returns the failedAVP (with the application that
// we don't support in our dictionary) and an error. Another cause
// for error is the presence of Inband Security, we don't support that.
func (cer *CER) Parse(m *diam.Message) (failedAVP *diam.AVP, err error) {
	if err = m.Unmarshal(cer); err != nil {
		return nil, err
	}
	if cer.InbandSecurityID != nil {
		return cer.InbandSecurityID, ErrNoCommonSecurity
	}
	if err = cer.sanityCheck(); err != nil {
		return nil, err
	}
	if failedAVP, err = cer.appsOK(m.Dictionary()); err != nil {
		return failedAVP, err
	}
	if len(cer.appID) == 0 {
		return nil, ErrMissingApplication
	}
	return nil, nil
}

// sanityCheck ensures all mandatory AVPs are present.
func (cer *CER) sanityCheck() error {
	if len(cer.OriginHost) == 0 {
		return ErrMissingOriginHost
	}
	if len(cer.OriginRealm) == 0 {
		return ErrMissingOriginRealm
	}
	if cer.OriginStateID == nil {
		return ErrMissingOriginStateID
	}
	return nil
}

// appsOK ensures all acct or auth applications in the CER
// exist in this server's dictionary.
func (cer *CER) appsOK(d *dict.Parser) (failedAVP *diam.AVP, err error) {
	failedAVP, err = cer.validateApps(d, "acct", cer.AcctAppID)
	if err != nil {
		return failedAVP, err
	}
	failedAVP, err = cer.validateApps(d, "auth", cer.AuthAppID)
	if err != nil {
		return failedAVP, err
	}
	if cer.VendorSpecificAppID != nil {
		for _, vs := range cer.VendorSpecificAppID {
			failedAVP, err = cer.validateApp(d, "acct", vs.AcctAppID)
			if err != nil {
				return failedAVP, err
			}
			failedAVP, err = cer.validateApp(d, "auth", vs.AuthAppID)
			if err != nil {
				return failedAVP, err
			}
		}
	}
	return nil, nil
}

// validateApps is a convenience method to test an array of application IDs.
func (cer *CER) validateApps(d *dict.Parser, appType string, appAVPs []*diam.AVP) (failedAVP *diam.AVP, err error) {
	if appAVPs != nil {
		for _, app := range appAVPs {
			failedAVP, err = cer.validateApp(d, appType, app)
			if err != nil {
				return failedAVP, err
			}
		}
	}
	return nil, nil
}

// validateApp ensures the given acct or auth application ID exists in
// this server's dictionary.
func (cer *CER) validateApp(d *dict.Parser, appType string, appAVP *diam.AVP) (failedAVP *diam.AVP, err error) {
	if appAVP != nil {
		id := uint32(appAVP.Data.(datatype.Unsigned32))
		app, err := d.App(id)
		if err != nil {
			return appAVP, &ErrNoCommonApplication{id, appType}
		}
		if len(app.Type) > 0 && app.Type != appType {
			return appAVP, &ErrNoCommonApplication{id, appType}
		}
		cer.appID = append(cer.appID, id)
	}
	return nil, nil
}

// Appliations return a list of supported application IDs from this CER.
// Must be called after Parse, otherwise it returns an empty array.
func (cer *CER) Applications() []uint32 {
	return cer.appID
}

var (
	// ErrMissingOriginHost is returned by CER.OK when the CER
	// does not contain an Origin-Host AVP.
	ErrMissingOriginHost = errors.New("missing Origin-Host")

	// ErrMissingOriginRealm is returned by CER.OK when the CER
	// does not contain an Origin-Realm AVP.
	ErrMissingOriginRealm = errors.New("missing Origin-Realm")

	// ErrMissingOriginStateID is returned by CER.OK when the CER
	// does not contain an Origin-State-Id AVP.
	ErrMissingOriginStateID = errors.New("missing Origin-State-Id")

	// ErrMissingApplication is returned by CER.OK when the CER
	// does not contain any Acct-Application-Id or Auth-Application-Id,
	// or their embedded versions in Vendor-Specific-Application-Id.
	ErrMissingApplication = errors.New("missing application")

	// ErrNoCommonSecurity is returned by CER.OK when the CER
	// contains the Inband-Security-Id. We don't support that.
	ErrNoCommonSecurity = errors.New("no common security")
)

// ErrNoCommonApplication is returned by CER.OK when the
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
