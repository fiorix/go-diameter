// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

// Diameter codes for the Result-Code AVP.
const (
	MultiRoundAuth                = 1001
	Success                       = 2001
	LimitedSuccess                = 2002
	CommandUnsupported            = 3001
	UnableToDeliver               = 3002
	RealmNotServed                = 3003
	TooBusy                       = 3004
	LoopDetected                  = 3005
	RedirectIndication            = 3006
	ApplicationUnsupported        = 3007
	InvalidHDRBits                = 3008
	InvalidAVPBits                = 3009
	UnknownPeer                   = 3010
	AuthenticationRejected        = 4001
	OutOfSpace                    = 4002
	ElectionLost                  = 4003
	AuthenticationDataUnavailable = 4181
	AVPUnsupported                = 5001
	UnknownUser                   = 5001
	UnknownSessionID              = 5002
	AuthorizationRejected         = 5003
	RoamingNotAllowed             = 5004
	InvalidAVPValue               = 5018
	MissingAVP                    = 5005
	IdentityAlreadyRegistered     = 5005
	ResourcesExceeded             = 5006
	ContradictingAVPs             = 5007
	ErrorInAssignmentType	      = 5007
	AVPNotAllowed                 = 5008
	AVPOccursTooManyTimes         = 5009
	NoCommonApplication           = 5010
	UnsupportedVersion            = 5011
	UnableToComply                = 5012
	InvalidBitInHeader            = 5013
	InvalidAVPLenght              = 5014
	InvalidMessageLength          = 5015
	InvalidAVPBitCombo            = 5016
	NoCommonSecurity              = 5017
	UnknownEpsSubscription        = 5420
	RatNotAllowed                 = 5421
	UnknownEquipment              = 5422
    UnknownServingNode            = 5423
	AbsentUser                    = 5550
	UserBusyForMtSms              = 5551
	FacilityNotSupported          = 5552
	IllegalUser                   = 5553
	IllegalEquipment              = 5554
	SmDeliveryFailure             = 5555
	ServiceNotSubscribed          = 5556
	ServiceBarred                 = 5557
	MwdListFull                   = 5558
	UserDataNotRecognized         = 5100
	OperationNotAllowed           = 5101
	UserDataCannotBeRead          = 5102
	UserDataCannotBeModified      = 5103
	UserDataCannotBeNotified      = 5104
	TooMuchData                   = 5008
	TransparentDataOutOfSync      = 5105
	FeatureUnsupported            = 5011
	SubsDataAbsent                = 5106
	NoSubscriptionToData          = 5107
	DsaiNotAvailable              = 5108
	IdentitiesDontMatch           = 5002
)
