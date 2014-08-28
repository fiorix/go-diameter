// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

// AVP types. Auto-generated from our dictionaries.
const (
	// List of AVP codes auto-generated from our dictionaries:
	// cat ../dict/testdata/*.xml | sed -n 's/.*avp name="\(.*\)" code="\([0-9]*\)".*/\1 = \2/p' | tr -d '-' | sort -u

	AccountingRealtimeRequired    = 483
	AccountingRecordNumber        = 485
	AccountingRecordType          = 480
	AccountingSessionId           = 44
	AccountingSubSessionId        = 287
	AcctApplicationId             = 259
	AcctInterimInterval           = 85
	AcctMultiSessionId            = 50
	AuthApplicationId             = 258
	AuthGracePeriod               = 276
	AuthRequestType               = 274
	AuthSessionState              = 277
	AuthorizationLifetime         = 291
	CCCorrelationId               = 411
	CCInputOctets                 = 412
	CCMoney                       = 413
	CCOutputOctets                = 414
	CCRequestNumber               = 415
	CCRequestType                 = 416
	CCServiceSpecificUnits        = 417
	CCSessionFailover             = 418
	CCSubSessionId                = 419
	CCTime                        = 420
	CCTotalOctets                 = 421
	CCUnitType                    = 454
	CheckBalanceResult            = 422
	Class                         = 25
	CostInformation               = 423
	CostUnit                      = 424
	CreditControl                 = 426
	CreditControlFailureHandling  = 427
	CurrencyCode                  = 425
	DestinationHost               = 293
	DestinationRealm              = 283
	DirectDebitingFailureHandling = 428
	DisconnectCause               = 273
	ErrorMessage                  = 281
	ErrorReportingHost            = 294
	EventTimestamp                = 55
	ExperimentalResult            = 297
	ExperimentalResultCode        = 298
	Exponent                      = 429
	FailedAVP                     = 279
	FinalUnitAction               = 449
	FinalUnitIndication           = 430
	FirmwareRevision              = 267
	GSUPoolIdentifier             = 453
	GSUPoolReference              = 457
	GrantedServiceUnit            = 431
	HostIPAddress                 = 257
	InbandSecurityId              = 299
	MultiRoundTimeOut             = 272
	MultipleServicesCreditControl = 456
	MultipleServicesIndicator     = 455
	OriginHost                    = 264
	OriginRealm                   = 296
	OriginStateId                 = 278
	ProductName                   = 269
	ProxyHost                     = 280
	ProxyInfo                     = 284
	ProxyState                    = 33
	RatingGroup                   = 432
	ReAuthRequestType             = 285
	RedirectAddressType           = 433
	RedirectHost                  = 292
	RedirectHostUsage             = 261
	RedirectMaxCacheTime          = 262
	RedirectServer                = 434
	RedirectServerAddress         = 435
	RequestedAction               = 436
	RequestedServiceUnit          = 437
	RestrictionFilterRule         = 438
	ResultCode                    = 268
	RouteRecord                   = 282
	ServiceContextId              = 461
	ServiceIdentifier             = 439
	ServiceParameterInfo          = 440
	ServiceParameterType          = 441
	ServiceParameterValue         = 442
	SessionBinding                = 270
	SessionId                     = 263
	SessionServerFailover         = 271
	SessionTimeout                = 27
	SubscriptionId                = 443
	SubscriptionIdData            = 444
	SubscriptionIdType            = 450
	SupportedVendorId             = 265
	TariffChangeUsage             = 452
	TariffTimeChange              = 451
	TerminationCause              = 295
	UnitValue                     = 445
	UsedServiceUnit               = 446
	UserEquipmentInfo             = 458
	UserEquipmentInfoType         = 459
	UserEquipmentInfoValue        = 460
	UserName                      = 1
	ValidityTime                  = 448
	ValueDigits                   = 447
	VendorId                      = 266
	VendorSpecificApplicationId   = 260
)
