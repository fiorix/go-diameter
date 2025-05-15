// UAR_Handle.go
package main

import (
	"fmt"
	"io"
	"log"

	"github.com/rakeshgmtke/go-diameter-hss/v4/diam"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/avp"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/datatype"
	_ "github.com/rakeshgmtke/go-diameter-hss/v4/diam/dict"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/sm"
)

func handleUAR(settings sm.Settings, stats *DiameterStats, enableLogging bool) diam.HandlerFunc {

	type UAR struct {
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
		UserName            datatype.UTF8String       `avp:"User-Name"`
		OCSupportedFeatures struct {
			OCFeatureVector datatype.Unsigned64 `avp:"OC-Feature-Vector"`
		} `avp:"OC-Supported-Features"`
		SupportedFeatures struct {
			VendorId      datatype.Unsigned32 `avp:"Vendor-Id"`
			FeatureListID datatype.Unsigned32 `avp:"Feature-List-ID"`
			FeatureList   datatype.Unsigned32 `avp:"Feature-List"`
		} `avp:"Supported-Features"`
		PublicIdentity           datatype.UTF8String  `avp:"Public-Identity"`
		VisitedNetworkIdentifier datatype.OctetString `avp:"Visited-Network-Identifier"`
		UserAuthorizationType    datatype.Enumerated  `avp:"User-Authorization-Type,omitempty"`
		UARFlags                 datatype.Unsigned32  `avp:"UAR-Flags"`
	}

	return func(c diam.Conn, m *diam.Message) {
		var err error
		var req UAR
		var impu string
		var impi string
		var scscf_name string

		if err := m.Unmarshal(&req); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}

		stats.IncrementReceived("UAR", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType))))

		if enableLogging {
			log.Printf("Received UAR from %s\n%s", c.RemoteAddr(), m)
		}

		Is_PublicIdentity_AVP, _ := m.FindAVP(avp.PublicIdentity, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_PublicIdentity_AVP != nil {
			impu = string(req.PublicIdentity)
			scscf_name = readSCSCFNameData(impu)
			//log.Printf("for received IMPU: %s stored SCSCF_NAME: %s", impu, scscf_name)
		}

		Is_UserName_AVP, _ := m.FindAVP(avp.UserName, 0) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_UserName_AVP != nil {
			impi = string(req.UserName)
			success := addOrModifyIMPI(impu, impi)
			if success {
				//log.Printf("addOrModifyIMPI func success for IMPU: %s is IMPI: %s is Stored is Successful", impu, impi)
			} else {
				//log.Printf("addOrModifyIMPI func failed IMPU: %s IMPI: %s", impu, impi)
			}
		}

		//stats.IncrementReceived("UAR", "10.1.1.1", "uar_type_1")

		////log.Printf("for IMPU: %s stored SCSCF_NAME2: %s", impu, scscf_name)
		//scscf_name = readSCSCFNameData(impu)
		////log.Printf("for IMPU: %s stored SCSCF_NAME3: %s", impu, scscf_name)

		//Creating Response
		a := m.Answer(0)
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

		Is_UserAuthorizationType_AVP, _ := m.FindAVP(avp.UserAuthorizationType, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_UserAuthorizationType_AVP received: %s : %d", Is_UserAuthorizationType_AVP, req.UserAuthorizationType)
		if Is_UserAuthorizationType_AVP != nil {
			if req.UserAuthorizationType == 0 && scscf_name == "" {
				stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-2001")

				// for Debug added
				/*success := addOrModifySCSCFName(impu, "scscf.com")
				if success {
					//log.Printf("addOrModifySCSCFName func success for IMPU: %s is SCSCF_NAME DUMMY is Stored is Successful", impu)
				} else {
					//log.Printf("addOrModifySCSCFName func failed IMPU: %s SCSCF_NAME DUMMY ", impu)
				}
				*/
				_, err = sendUAA_UAT_INIT_REG(settings, c, a, enableLogging)
				if err != nil {
					//log.Printf("Failed to send sendUAA_UAT_INIT_REG: %s", err.Error())
				}
			} else if req.UserAuthorizationType == 0 && scscf_name != "" {
				a.NewAVP(avp.ServerName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String(scscf_name))
				stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-2002")
				_, err = sendUAA_UAT_SUBQ_REG(settings, c, a, enableLogging)
				if err != nil {
					//log.Printf("Failed to send sendUAA_UAT_SUBQ_REG: %s", err.Error())
				}
			} else if req.UserAuthorizationType == 1 && scscf_name == "" {
				stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-2004")
				_, err = sendUAA_UAT_1_NOTREGTRD(settings, c, a, enableLogging)
				if err != nil {
					//log.Printf("Failed to send sendUAA_UAT_1_NOTREGTRD: %s", err.Error())
				}
			} else if req.UserAuthorizationType == 1 && scscf_name != "" {
				a.NewAVP(avp.ServerName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String(scscf_name))
				stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-2001")
				_, err = sendUAA_UAT_1_REGTRD(settings, c, a, enableLogging)
				if err != nil {
					//log.Printf("Failed to send sendUAA_UAT_1_REGTRD: %s", err.Error())
				}
			} else if req.UserAuthorizationType == 2 {
				stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-2001")
				_, err = sendUAA_UAT_2(settings, c, a, enableLogging)
				if err != nil {
					//log.Printf("Failed to send sendUAA_UAT_2: %s", err.Error())
				}
			}
		} else {
			stats.IncrementReceived("UAA", string(req.OriginHost), "uar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.UserAuthorizationType)))+"-RESP-CODE-5012")
			_, err = sendUAA_default(settings, c, a, enableLogging)
			if err != nil {
				//log.Printf("Failed to send sendUAA_default when UAT AVP is missing: %s", err.Error())
			}
		}

	}

}

func sendUAA_UAT_INIT_REG(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
	})

	if enableLogging {
		log.Printf("inside sendUAA_UAT_INIT_REG func and SCSCF NAME EMPTY")
		log.Printf("Sending UAA to \n%s", m)
	}

	return m.WriteTo(w)
}

func sendUAA_UAT_SUBQ_REG(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2002)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
	})

	if enableLogging {
		log.Printf("inside sendUAA_UAT_SUBQ_REG func and SCSCF NAME PREESENT")
		log.Printf("Sending UAA to \n%s", m)
	}

	return m.WriteTo(w)
}

func sendUAA_UAT_1_NOTREGTRD(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2004)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
	})

	if enableLogging {
		log.Printf("inside sendUAA_UAT_1_NOTREGTRD func and SCSCF NAME EMPTY")
		log.Printf("Sending UAA to \n%s", m)
	}
	return m.WriteTo(w)
}

func sendUAA_UAT_1_REGTRD(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(diam.Success))

	/*	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)),
				diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
			},
		})
	*/
	if enableLogging {
		log.Printf("inside sendUAA_UAT_1_REGTRD func and SCSCF NAME PREESENT")
		log.Printf("Sending UAA to \n%s", m)
	}

	return m.WriteTo(w)
}

func sendUAA_UAT_2(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ServerCapabilities, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MandatoryCapability, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(MandatoryCapability)),
			//diam.NewAVP(avp.OptionalCapability, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(1)),
			//diam.NewAVP(avp.ServerName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("sip:scscf.maavenir.com")),
		},
	})
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))

	if enableLogging {
		log.Printf("inside sendUAA_UAT_2 func")
		log.Printf("Sending UAA to \n%s", m)
	}
	return m.WriteTo(w)
}

func sendUAA_default(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5012)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
		},
	})
	if enableLogging {
		log.Printf("User-Authorization-Type: AVP is missing")
		log.Printf("inside sendUAA_default func")
		log.Printf("Sending UAA to \n%s", m)
	}

	return m.WriteTo(w)
}
