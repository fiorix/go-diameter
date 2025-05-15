// LIR_Handle.go
package main

import (
	//"fmt"
	"io"
	"log"

	/*
		"github.com/fiorix/go-diameter/v4/diam"
		"github.com/fiorix/go-diameter/v4/diam/avp"
		"github.com/fiorix/go-diameter/v4/diam/datatype"
		_ "github.com/fiorix/go-diameter/v4/diam/dict"
		"github.com/fiorix/go-diameter/v4/diam/sm"

	*/
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/avp"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/datatype"
	_ "github.com/rakeshgmtke/go-diameter-hss/v4/diam/dict"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/sm"
)

func handleLIR(settings sm.Settings, stats *DiameterStats, enableLogging bool) diam.HandlerFunc {

	type LIR struct {
		SessionID                   datatype.UTF8String `avp:"Session-Id"`
		DRMP                        datatype.Enumerated `avp:"DRMP"`
		VendorSpecificApplicationID struct {
			VendorId          datatype.Unsigned32 `avp:"Vendor-Id"`
			AuthApplicationId datatype.Unsigned32 `avp:"Auth-Application-Id"`
		} `avp:"Vendor-Specific-Application-Id"`
		AuthSessionState datatype.Enumerated       `avp:"Auth-Session-State"`
		OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
		OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
		DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
		DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
		ProxyInfo        struct {
			ProxyHost  datatype.DiameterIdentity `avp:"Proxy-Host"`
			ProxyState datatype.OctetString      `avp:"Proxy-State"`
		} `avp:"Proxy-Info"`
		RouteRecord         datatype.DiameterIdentity `avp:"Route-Record,omitempty"`
		OCSupportedFeatures struct {
			OCFeatureVector datatype.Unsigned64 `avp:"OC-Feature-Vector"`
		} `avp:"OC-Supported-Features"`
		SupportedFeatures struct {
			VendorId      datatype.Unsigned32 `avp:"Vendor-Id"`
			FeatureListID datatype.Unsigned32 `avp:"Feature-List-ID"`
			FeatureList   datatype.Unsigned32 `avp:"Feature-List"`
		} `avp:"Supported-Features"`

		PublicIdentity        datatype.UTF8String `avp:"Public-Identity"`
		ServerAssignmentType  datatype.Enumerated `avp:"Server-Assignment-Type"`
		UserAuthorizationType datatype.Enumerated `avp:"User-Authorization-Type"`
		OriginatingRequest    datatype.Enumerated `avp:"Originating-Request"`
		SessionPriority       datatype.Enumerated `avp:"Session-Priority"`
	}

	return func(c diam.Conn, m *diam.Message) {
		var err error
		var req LIR
		var impu string
		//var impi string
		var scscf_name string
		//var msisdn string

		if err := m.Unmarshal(&req); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}

		if enableLogging {
			log.Printf("Received LIR from %s\n%s", c.RemoteAddr(), m)
		}
		Is_PublicIdentity_AVP, _ := m.FindAVP(avp.PublicIdentity, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_PublicIdentity_AVP != nil {
			impu = string(req.PublicIdentity)
		}

		scscf_name = readSCSCFNameData(impu)

		Is_UserAuthorizationType_AVP, _ := m.FindAVP(avp.UserAuthorizationType, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_UserAuthorizationType_AVP received: %s : %d", Is_UserAuthorizationType_AVP, req.UserAuthorizationType)
		if Is_UserAuthorizationType_AVP != nil && req.UserAuthorizationType == 2 {
			stats.IncrementReceived("LIR", string(req.OriginHost), "CAPABILITIES-QUERY")
		} else {
			//Is_UserAuthorizationType_AVP == nil && Is_PublicIdentity_AVP != nil  {
			stats.IncrementReceived("LIR", string(req.OriginHost), "SCSCF-NAME-QUERY")
		}

		if enableLogging {
			log.Printf("from LIR received IMPU: %s STORED SCSCF_NAME : %s ", impu, scscf_name)
		}

		//Creating Response
		a := m.Response()
		a.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)

		//Check is Vendor-Specific-Application-Id present in Request, if yes.. include Vendor-Specific-Application-Id in response also.
		Is_VendorSpecificApplicationID_AVP, _ := m.FindAVP(avp.VendorSpecificApplicationID, 0) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_VendorSpecificApplicationID_AVP != nil {
			a.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(req.VendorSpecificApplicationID.AuthApplicationId)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(req.VendorSpecificApplicationID.VendorId)),
				},
			})
		}

		//Check for AuthSessionState AVP present, if yes add AuthSessionState AVP in response
		Is_authSessionState_AVP, _ := m.FindAVP(avp.AuthSessionState, 0) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_authSessionState_AVP != nil {
			a.NewAVP(avp.AuthSessionState, avp.Mbit, 0, req.AuthSessionState)
		}

		Is_SupportedFeatures_AVP, _ := m.FindAVP(avp.SupportedFeatures, 0) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_SupportedFeatures_AVP != nil {
			a.NewAVP(avp.SupportedFeatures, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(req.VendorSpecificApplicationID.VendorId)),
					//	diam.NewAVP(avp.VendorId, avp.Mbit, 0, datatype.Unsigned32(req.SupportedFeatures.VendorId)),
					diam.NewAVP(avp.FeatureListID, avp.Mbit, 0, datatype.Unsigned32(req.SupportedFeatures.FeatureListID)),
					diam.NewAVP(avp.FeatureList, avp.Mbit, 0, datatype.Unsigned32(req.SupportedFeatures.FeatureList)),
				},
			})
		}

		//Is_UserAuthorizationType_AVP, _ := m.FindAVP(avp.UserAuthorizationType, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_UserAuthorizationType_AVP received: %s : %d", Is_UserAuthorizationType_AVP, req.UserAuthorizationType)
		//if Is_UserAuthorizationType_AVP != nil ||  {
		if req.UserAuthorizationType == 2 {
			stats.IncrementReceived("LIA", string(req.OriginHost), "CAPABILITIES-RESP-CODE-2001")
			m.NewAVP(avp.ServerCapabilities, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.MandatoryCapability, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(MandatoryCapability)),
					//diam.NewAVP(avp.OptionalCapability, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(1)),
					//diam.NewAVP(avp.ServerName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("sip:scscf.maavenir.com")),
				},
			})
			m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))

		} else if scscf_name != "" {
			//SCSCF_NAME is stored. returning with Stored SCSCF_NAME with Success
			if enableLogging {
				log.Printf("LIA sending SCSCF_NAME is stored", scscf_name)
			}
			stats.IncrementReceived("LIA", string(req.OriginHost), "SCSCF-NAME-RESP-CODE-2001")
			a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))
			a.NewAVP(avp.ServerName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String(scscf_name))
		} else {
			//SCSCF_NAME is not stored. returning with DIAMETER_UNREGISTERED_SERVICE.
			stats.IncrementReceived("LIA", string(req.OriginHost), "SCSCF-NAME-RESP-CODE-5003")
			if enableLogging {
				log.Printf("LIA sending SCSCF_NAME is Not stored sending DIAMETER_UNREGISTERED_SERVICE")
			}
			a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(0)),
					diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5003)),
				},
			})
		}
		//}

		_, err = sendLIA_default(settings, c, a, enableLogging)
		if err != nil {
			log.Printf("LIA sending ERROR %s", err.Error())

		}
	}
}

func sendLIA_default(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	if enableLogging {
		log.Printf("inside sendLIA_default func")
		log.Printf("Sending LIA to \n%s", m)
	}

	return m.WriteTo(w)
}
