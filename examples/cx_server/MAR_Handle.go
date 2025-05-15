// MAR_Handle.go
package main

import (
	//"fmt"
	"io"
	"log"

	"github.com/rakeshgmtke/go-diameter-hss/v4/diam"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/avp"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/datatype"
	_ "github.com/rakeshgmtke/go-diameter-hss/v4/diam/dict"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/sm"
)

func handleMAR(settings sm.Settings, stats *DiameterStats, enableLogging bool) diam.HandlerFunc {

	type MAR struct {
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
		ServerName          datatype.UTF8String       `avp:"Server-Name"`
		OCSupportedFeatures struct {
			OCFeatureVector datatype.Unsigned64 `avp:"OC-Feature-Vector"`
		} `avp:"OC-Supported-Features"`
		SupportedFeatures struct {
			VendorId      datatype.Unsigned32 `avp:"Vendor-Id"`
			FeatureListID datatype.Unsigned32 `avp:"Feature-List-ID"`
			FeatureList   datatype.Unsigned32 `avp:"Feature-List"`
		} `avp:"Supported-Features"`
		PublicIdentity     datatype.UTF8String `avp:"Public-Identity"`
		SIPNumberAuthItems datatype.Unsigned32 `avp:"SIP-Number-Auth-Items"`
		SIPAuthDataItem    struct {
			SIPItemNumber           datatype.Unsigned32  `avp:"SIP-Item-Number"`
			SIPAuthenticationScheme datatype.UTF8String  `avp:"SIP-Authentication-Scheme"`
			SIPAuthorization        datatype.OctetString `avp:"SIP-Authorization"`
		} `avp:"SIP-Auth-Data-Item"`
	}

	return func(c diam.Conn, m *diam.Message) {
		var err error
		var req MAR
		var impu string
		var impi string
		var scscf_stored_name string

		if err := m.Unmarshal(&req); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}

		stats.IncrementReceived("MAR", string(req.OriginHost), string(req.SIPAuthDataItem.SIPAuthenticationScheme))

		if enableLogging {
			log.Printf("Received MAR from %s\n%s", c.RemoteAddr(), m)
		}

		Is_PublicIdentity_AVP, _ := m.FindAVP(avp.PublicIdentity, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_PublicIdentity_AVP != nil {
			impu = string(req.PublicIdentity)
			scscf_stored_name = readSCSCFNameData(impu)
			if enableLogging {
				log.Printf("for received IMPU: %s stored SCSCF_NAME: %s", impu, scscf_stored_name)
			}

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

		Is_ServerName_AVP, _ := m.FindAVP(avp.ServerName, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_ServerName_AVP != nil {
			success := addOrModifySCSCFName(impu, string(req.ServerName))
			if success {
				if enableLogging {
					log.Printf("addOrModifySCSCFName func success for IMPU: %s is SCSCF_NAME: %s is Stored is Successful", impu, req.ServerName)
				}
			} else {
				if enableLogging {
					log.Printf("addOrModifySCSCFName func failed IMPU: %s SCSCF_NAME : %s ", impu, req.ServerName)
				}
			}
		}
		//Creating Response
		a := m.Response()
		a.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		a.NewAVP(avp.UserName, avp.Mbit, 0, req.UserName)

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

		Is_SIPNumberAuthItems_AVP, _ := m.FindAVP(avp.SIPNumberAuthItems, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_SIPNumberAuthItems_AVP != nil {
			a.NewAVP(avp.SIPNumberAuthItems, avp.Mbit, VENDOR_3GPP, req.SIPNumberAuthItems)
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

		///SIPAuthDataItem
		Is_SIPAuthDataItem_AVP, _ := m.FindAVP(avp.SIPAuthDataItem, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_SIPAuthDataItem_AVP received: %s : %d", Is_SIPAuthDataItem_AVP, req.SIPAuthDataItem)
		if Is_SIPAuthDataItem_AVP != nil {
			if req.SIPAuthDataItem.SIPAuthenticationScheme == "Digest-AKAv1-MD5" {
				//Digest-AKAv1-MD5
				a.NewAVP(avp.SIPAuthDataItem, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SIPItemNumber, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(1)),
						diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("Digest-AKAv1-MD5")),
						diam.NewAVP(avp.SIPAuthenticate, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xed\x9a\xff\x45\xe7\xde\x34\xbb\xe5\xd1\x38\x99\x96\xc5\x0b\x2b\xe1\xba\x78\xc5\xe5\x35\xb9\xb9\xb2\xf1\x84\x64\x4a\x0b\xcc\x44")),
						diam.NewAVP(avp.SIPAuthorization, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xc3\x16\x26\xaa\xef\x74\x0c\xce")),
						diam.NewAVP(avp.ConfidentialityKey, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\x2e\x50\x5f\xc5\x92\xad\x9b\x79\xd0\xf0\xdb\x6c\xc3\xf2\x1f\xf4")),
						diam.NewAVP(avp.IntegrityKey, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xd3\x64\x45\x95\x32\xd4\x7f\x3c\x02\x60\x8b\xe6\x32\xe9\x02\xac")),
					},
				})
			} else if req.SIPAuthDataItem.SIPAuthenticationScheme == "Digest" {
				//Digest
				a.NewAVP(avp.SIPAuthDataItem, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SIPItemNumber, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(1)),
						diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("Digest")),
						diam.NewAVP(avp.SIPAuthenticationContext, avp.Mbit, VENDOR_3GPP, datatype.OctetString("ims*1234")),
					},
				})
			} else if req.SIPAuthDataItem.SIPAuthenticationScheme == "SIP Digest" {
				//SIP Digest
				a.NewAVP(avp.SIPAuthDataItem, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("SIP Digest")),
						diam.NewAVP(avp.SIPDigestAuthenticate, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								diam.NewAVP(avp.DigestRealm, avp.Mbit, 0, datatype.UTF8String("sha-it.mavenir.com")),
								diam.NewAVP(avp.DigestAlgorithm, avp.Mbit, 0, datatype.UTF8String("MD5")),
								diam.NewAVP(avp.DigestQop, avp.Mbit, 0, datatype.UTF8String("auth")),
								diam.NewAVP(avp.DigestHA1, avp.Mbit, 0, datatype.UTF8String("82110794a5fa5a00bf5a3af1eb5d3c14")),
							},
						}),
					},
				})
			} else {
				//Digest-AKAv1-MD5
				a.NewAVP(avp.SIPAuthDataItem, avp.Mbit, VENDOR_3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SIPItemNumber, avp.Mbit, VENDOR_3GPP, datatype.Unsigned32(1)),
						diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("Digest-AKAv1-MD5")),
						diam.NewAVP(avp.SIPAuthenticate, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xed\x9a\xff\x45\xe7\xde\x34\xbb\xe5\xd1\x38\x99\x96\xc5\x0b\x2b\xe1\xba\x78\xc5\xe5\x35\xb9\xb9\xb2\xf1\x84\x64\x4a\x0b\xcc\x44")),
						diam.NewAVP(avp.SIPAuthorization, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xc3\x16\x26\xaa\xef\x74\x0c\xce")),
						diam.NewAVP(avp.ConfidentialityKey, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\x2e\x50\x5f\xc5\x92\xad\x9b\x79\xd0\xf0\xdb\x6c\xc3\xf2\x1f\xf4")),
						diam.NewAVP(avp.IntegrityKey, avp.Mbit, VENDOR_3GPP, datatype.OctetString("\xd3\x64\x45\x95\x32\xd4\x7f\x3c\x02\x60\x8b\xe6\x32\xe9\x02\xac")),
					},
				})
			}
		}

		stats.IncrementReceived("MAA", string(req.OriginHost), string(req.SIPAuthDataItem.SIPAuthenticationScheme)+"-RESP-CODE-2001")

		_, err = sendMAA_default(settings, c, a, enableLogging)
		if err != nil {
			//log.Printf("MAA sending ERROR %s", err.Error())

		}
	}
}

func sendMAA_default(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))
	if enableLogging {
		log.Printf("inside sendMAA_default func")
		log.Printf("Sending MAA to \n%s", m)
	}

	return m.WriteTo(w)
}
