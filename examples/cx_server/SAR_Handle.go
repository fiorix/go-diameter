// SAR_Handle.go
package main

import (
	_ "encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"io"

	"github.com/rakeshgmtke/go-diameter-hss/v4/diam"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/avp"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/datatype"
	_ "github.com/rakeshgmtke/go-diameter-hss/v4/diam/dict"
	"github.com/rakeshgmtke/go-diameter-hss/v4/diam/sm"
)

func handleSAR(settings sm.Settings, stats *DiameterStats, enableLogging bool, ifcXmlFile string) diam.HandlerFunc {

	type SAR struct {
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

		UserName                       datatype.UTF8String  `avp:"User-Name"`
		ServerName                     datatype.UTF8String  `avp:"Server-Name"`
		PublicIdentity                 datatype.UTF8String  `avp:"Public-Identity"`
		WildcardedPublicIdentity       datatype.UTF8String  `avp:"Wildcarded-Public-Identity"`
		ServerAssignmentType           datatype.Enumerated  `avp:"Server-Assignment-Type"`
		UserDataAlreadyAvailable       datatype.Enumerated  `avp:"User-Data-Already-Available"`
		UserData                       datatype.OctetString `avp:"User-Data"`
		SessionPriority                datatype.Enumerated  `avp:"Session-Priority"`
		MultipleRegistrationIndication datatype.Enumerated  `avp:"Multiple-Registration-Indication"`
		SCSCFRestorationInfo           *diam.AVP            `avp:"SCSCF-Restoration-Info"`

		/*SCSCFRestorationInfo           struct {
			UserName                datatype.UTF8String `avp:"User-Name"`
			SIPAuthenticationScheme datatype.UTF8String `avp:"SIP-Authentication-Scheme"`
			RestorationInfo         struct {
				Path             datatype.OctetString `avp:"Path"`
				Contact          datatype.OctetString `avp:"Contact"`
				SubscriptionInfo struct {
					CallIDSIPHeader datatype.OctetString `avp:"Call-ID-SIP-Header"`
					FromSIPHeader   datatype.OctetString `avp:"From-SIP-Header"`
					ToSIPHeader     datatype.OctetString `avp:"To-SIP-Header"`
					RecordRoute     datatype.OctetString `avp:"Record-Route"`
					Contact         datatype.OctetString `avp:"Contact"`
				} `avp:"Subscription-Info"`
			} `avp:"Restoration-Info"`
		} `avp:"SCSCF-Restoration-Info"`
		*/
	}

	return func(c diam.Conn, m *diam.Message) {
		var err error
		var req SAR
		//var code uint32
		var impu string
		var impi string
		var scscf_name string
		var scscf_stored_name string
		var msisdn string
		var tel string
		//var xmlUserData string

		//avpSCSCFRestorationInfo := make(map[string]*diam.AVP, 5000000)
		avpSCSCFRestorationInfo := make(map[string]*diam.AVP)

		//log.Printf("xmlUserData :%s", xmlUserData)

		if err := m.Unmarshal(&req); err != nil {
			log.Printf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}

		//SAR Count
		stats.IncrementReceived("SAR", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType))))

		if enableLogging {
			log.Printf("Received SAR from %s\n%s", c.RemoteAddr(), m)
		}

		Is_PublicIdentity_AVP, _ := m.FindAVP(avp.PublicIdentity, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_PublicIdentity_AVP != nil {
			impu = string(req.PublicIdentity)
		}

		/*		Is_UserName_AVP, _ := m.FindAVP(avp.UserName, 0) // Provide both AVP code and Vendor ID (0 for standard)
				if Is_UserName_AVP != nil {
					impi = string(req.UserName)
				}
		*/
		Is_ServerName_AVP, _ := m.FindAVP(avp.ServerName, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_ServerName_AVP != nil {
			scscf_name = string(req.ServerName)
		}

		//Creating Response
		a := m.Response()
		a.NewAVP(avp.SessionID, avp.Mbit, 0, req.SessionID)
		a.NewAVP(avp.OriginHost, avp.Mbit, 0, settings.OriginHost)
		a.NewAVP(avp.OriginRealm, avp.Mbit, 0, settings.OriginRealm)
		//a.NewAVP(avp.UserName, avp.Mbit, 0, req.UserName)
		//a.AddAVP(avpSCSCFRestorationInfo[impu])

		Is_UserName_AVP, _ := m.FindAVP(avp.UserName, 0) // Provide both AVP code and Vendor ID (0 for standard)
		if Is_UserName_AVP != nil {
			impi = string(req.UserName)
			a.NewAVP(avp.UserName, avp.Mbit, 0, req.UserName)
		} else if req.ServerAssignmentType == 3 {
			//impi = strings.TrimPrefix(impu, "sip:")
			if strings.HasPrefix(impu, "sip:+") {
				impi = strings.TrimPrefix(impu, "sip:+")
			} else {
				impi = strings.TrimPrefix(impu, "sip:")
			}
			a.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(impi))
		}

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

		//Is_UserDataAlreadyAvailable_AVP, _ := m.FindAVP(avp.UserDataAlreadyAvailable, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)

		//ServerAssignmentType check send UserData if UserDataAlreadyAvailable is request
		Is_ServerAssignmentType_AVP, _ := m.FindAVP(avp.ServerAssignmentType, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_ServerAssignmentType_AVP received: %s : %d", Is_ServerAssignmentType_AVP, req.ServerAssignmentType)
		if Is_ServerAssignmentType_AVP != nil {
			if req.ServerAssignmentType == 0 || req.ServerAssignmentType == 1 || req.ServerAssignmentType == 2 || req.ServerAssignmentType == 3 {
				//ifc xml reading and preparing for sending
				if enableLogging {
					log.Printf("XML file file name: %s", ifcXmlFile)
				}
				xmlUserData1, _ := os.ReadFile(ifcXmlFile)
				xmlUserData := string(xmlUserData1)
				scscf_stored_name = readSCSCFNameData(impu)
				msisdn = onlyNumbers(impu)
				tel = getTelNumber(impu)
				xmlUserData = strings.ReplaceAll(xmlUserData, "IMPU", impu)
				xmlUserData = strings.ReplaceAll(xmlUserData, "IMPI", impi)
				xmlUserData = strings.ReplaceAll(xmlUserData, "TEL", tel)

				if enableLogging {
					log.Printf("from SAR received IMPU: %s IMPI: %s SCSCF_NAME : %s MSISDN no.: %s", impu, impi, scscf_name, msisdn)
					log.Printf("for received IMPU: %s stored SCSCF_NAME: %s and received SCSCF NAME: %s : SAR_TYPE: %d", impu, scscf_stored_name, scscf_name, datatype.Enumerated(req.ServerAssignmentType))
					//log.Printf("xmlUserData :%s", xmlUserData)
				}
				//		log.Printf("USER-DATA Added:")
				a.NewAVP(avp.UserData, avp.Mbit, VENDOR_3GPP, datatype.OctetString(xmlUserData))
			}
		}

		///ServerAssignmentType to check server name
		//Is_ServerAssignmentType_AVP, _ := m.FindAVP(avp.ServerAssignmentType, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
		//log.Printf("Is_ServerAssignmentType_AVP received: %s : %d", Is_ServerAssignmentType_AVP, req.ServerAssignmentType)
		if Is_ServerAssignmentType_AVP != nil {
			if req.ServerAssignmentType == 0 {
				//log.Printf("ServerAssignmentType received:", req.ServerAssignmentType)
				//#NO_ASSIGNMENT  "S-CSCF restoration information for a registered Public User Identity"
				//https://realtimecommunication.wordpress.com/2016/05/25/rainy-day-scenarios-s-cscf-restoration/
				//if restoration info present then Add it in reponse

				/*If it indicates NO_ASSIGNMENT, the HSS checks whether the Public Identity is assigned for the S-CSCF requesting the data.
				If the requesting S-CSCF is not the same as the assigned S-CSCF and the S-CSCF reassignment pending flag is not set, the Result-Code shall be set to DIAMETER_UNABLE_TO_COMPLY,
				otherwise the HSS shall download the relevant user information and the Result-Code shall be set to DIAMETER_SUCCESS.
				If relevant S-CSCF Restoration Information is stored in the HSS and IMS Restoration Procedures are supported,
				it shall be added to the answer message. If there is S-CSCF Restoration Information associated with several Private User Identities,
				the HSS shall include all the S-CSCF Restoration Information groups in the response. If there are multiple Private User Identities, which belong to
				the served IMS subscription the Associated-Identities AVP should be added to the answer message and it shall contain all Private User Identities associated to the IMS subscription.
				*/

				if scscf_stored_name != "" {
					if scscf_stored_name == scscf_name {
						a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)) //DIAMETER_SUCCESS

						//if restoration info present for impu then Add it in response
						if avpSCSCFRestorationInfo[impu] != nil {
							if enableLogging {
								log.Printf("adding avpSCSCFRestorationInfo for IMPU SAR_TYPE 3:", impu)
							}
							a.AddAVP(avpSCSCFRestorationInfo[impu])
							stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-0"+"restore-info"+"-2001")
						} else {
							stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-0"+"-2001")
						}

						// need to add RESTORATION INFO if present.
					} else {
						a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
								diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5012)), //DIAMETER_UNABLE_TO_COMPLY
							},
						})
						stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-0"+"-5012")
					}
				} else {
					stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-0"+"-failed")
				}
			} else if req.ServerAssignmentType == 1 || req.ServerAssignmentType == 2 {
				//log.Printf("ServerAssignmentType received:", req.ServerAssignmentType)
				//#REGISTER or RE_REGISTER
				/*Store SCSCFNAME AND restoration : and Check scscf name changes if different send error reponse and exit
				Check the Server Assignment Type value received in the request:
				 - If it indicates REGISTRATION or RE_REGISTRATION, the HSS shall check whether the Public Identity is assigned for the S-CSCF requesting the data. If there is already an S-CSCF assigned to the user and the requesting S-CSCF is not the same as the previously assigned S-CSCF, and IMS restoration procedures are not supported or the S-CSCF reassignment pending flag is not set, the HSS shall include the name of the previously assigned S-CSCF in the response message and the Experimental-Result-Code shall be set to DIAMETER_ERROR_IDENTITY_ALREADY_REGISTERED.

				If it is REGISTRATION and the HSS implements IMS Restoration procedures, if multiple registration indication is included in the request and the Public User Identity is stored as registered in the HSS, and there is restoration information related to the Private User Identity, the HSS shall not overwrite the stored S-CSCF Restoration Information, instead, it shall send the stored S-CSCF restoration information together with the user profile in the SAA. The Experimental-Result-Code shall be set to DIAMETER_ERROR_IN_ASSIGNMENT_TYPE (5007). Otherwise, the HSS shall store the received S-CSCF restoration information. The  Result-Code shall be set to DIAMETER_SUCCESS.
				*/

				//Check is SCSCFRestorationInfo info present and store in avpSCSCFRestorationInfo map
				if req.SCSCFRestorationInfo != nil {
					avpSCSCFRestorationInfo[impu] = req.SCSCFRestorationInfo
					if enableLogging {
						log.Printf("imsRestorationInfo IMPU store: %s imsRestorationInfo : %s ", impu, avpSCSCFRestorationInfo[impu])
					}
				}
				if scscf_stored_name != "" {
					if scscf_stored_name == scscf_name {
						a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)) //DIAMETER_SUCCESS
						stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType)))+"-2001")
						// need to store RESTORATION INFO if present.
					} else {
						a.NewAVP(avp.ServerName, avp.Mbit, 0, datatype.UTF8String(scscf_stored_name)) //Stored SCSCF NAME added
						a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
								diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5005)), //DIAMETER_ERROR_IDENTITY_ALREADY_REGISTERED
							},
						})
						stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType)))+"-5005")
					}
				} else {
					stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType)))+"-failed")
				}
			} else if req.ServerAssignmentType == 3 {
				//log.Printf("ServerAssignmentType received:", req.ServerAssignmentType)
				//#UNREGISTERED_USER

				/*
					If it indicates UNREGISTERED_USER, the HSS shall check whether the Public Identity is assigned for the S-CSCF requesting the data.
					If there is already an S-CSCF assigned to the user and the requesting S-CSCF is not the same as the previously assigned S-CSCF, and IMS restoration procedures are not supported or the S-CSCF reassignment pending flag is not set, the HSS shall include the name of the previously assigned S-CSCF in the response message and the Experimental-Result-Code shall be set to	DIAMETER_ERROR_IDENTITY_ALREADY_REGISTERED.

					If there is already an S-CSCF assigned to the user and the requesting S-CSCF is not the same as the previously assigned S-CSCF and IMS restoration procedures are supported, and the S-CSCF reassignment pending flag is set, the HSS shall overwrite the S-CSCF name and shall reset the S-CSCF reassignment pending flag.

					If there is no S-CSCF assigned to the user or the requesting S-CSCF is the same as the previously assigned SCSCF stored in the HSS, the HSS shall store the S-CSCF name.

					If the registration state of the Public Identity is not registered, the HSS shall set the registration state of the Public Identity as unregistered, i.e. registered as a consequence of an originating or terminating request and download the relevant user information. The Result-Code shall be set to DIAMETER_SUCCESS. If there are multiple Private User Identities associated to the Public User Identity in the HSS, the HSS shall arbitrarily select one of the Private User Identities and put it into the response message.
				*/

				if scscf_stored_name != "" {
					if scscf_stored_name == scscf_name {
						a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)) //DIAMETER_SUCCESS
						//if restoration info present for impu then Add it in response
						if avpSCSCFRestorationInfo[impu] != nil {
							if enableLogging {
								log.Printf("adding avpSCSCFRestorationInfo for IMPU SAR_TYPE 3:", impu)
							}
							a.AddAVP(avpSCSCFRestorationInfo[impu])
							stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-3"+"restore-info"+"-2001")
						}
						// need to add RESTORATION INFO if present.
					} else {
						a.NewAVP(avp.ServerName, avp.Mbit, 0, datatype.UTF8String(scscf_stored_name)) //Stored SCSCF NAME added
						a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
								diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5005)), //DIAMETER_ERROR_IDENTITY_ALREADY_REGISTERED
							},
						})
						stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-3"+"-5005")
					}
				} else {

					Is_ServerName_AVP, _ := m.FindAVP(avp.ServerName, VENDOR_3GPP) // Provide both AVP code and Vendor ID (0 for standard)
					if Is_ServerName_AVP != nil {
						success := addOrModifySCSCFName(impu, string(req.ServerName))
						if success {
							if enableLogging {
								log.Printf("SAR TYPE 3 addOrModifySCSCFName func success for IMPU: %s is SCSCF_NAME: %s is Stored is Successful", impu, req.ServerName)
							}
						} else {
							if enableLogging {
								log.Printf("SAR TYPE 3 addOrModifySCSCFName func failed IMPU: %s SCSCF_NAME : %s ", impu, req.ServerName)
							}
						}
					}
					stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-3"+"-server-name-empty-2001")
					a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)) //DIAMETER_SUCCESS
				}

			} else if req.ServerAssignmentType == 4 || req.ServerAssignmentType == 5 || req.ServerAssignmentType == 8 ||
				req.ServerAssignmentType == 9 || req.ServerAssignmentType == 10 || req.ServerAssignmentType == 11 {
				//log.Printf("ServerAssignmentType received:", req.ServerAssignmentType)

				//#TIMEOUT_DEREGISTRATION or USER_DEREGISTRATION or ADMINISTRATIVE_DEREGISTRATION
				//#AUTHENTICATION_FAILURE  or AUTHENTICATION_TIMEOUT or DEREGISTRATION_TOO_MUCH_DATA

				// delete scscf name and delete IMS restoration info ++ DIAMETER_RESULT_SUCCESS

				deleteIMPUData(impu)
				delete(avpSCSCFRestorationInfo, impu)

				// tmp
				if enableLogging {
					log.Printf("delete deleteIMPUData for IMPU SAR_TYPE IMPU :%s :SCSCF NAME: %s", impu, readSCSCFNameData(impu))
					log.Printf("delete avpSCSCFRestorationInfo for IMPU SAR_TYPE IMPU:%s :restoration info: %s", impu, avpSCSCFRestorationInfo[impu])
				}

				a.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001)) //DIAMETER_SUCCESS
				stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType)))+"-2001")

				// it will delete SCSCFName/IMPI/IMSRestorationInfo for that IMPU

			} else if req.ServerAssignmentType == 6 || req.ServerAssignmentType == 7 {
				//log.Printf("ServerAssignmentType received:", req.ServerAssignmentType)

				//#TIMEOUT_DEREGISTRATION_STORE_SERVER_NAME  or USER_DEREGISTRATION_STORE_SERVER_NAME

				/*- If it indicates TIMEOUT_DEREGISTRATION_STORE_SERVER_NAME or USER_DEREGISTRATION_STORE_SERVER_NAME the HSS decides
				whether to keep the S-CSCF name associated to the Private User Identity stored or not for all the Public User Identities that the S-CSCF indicated in the request.
				If no Public User Identity is present in the request, the Private User Identity shall be present.
				- If the HSS decides to keep the S-CSCF name stored the HSS shall keep the S-CSCF name stored for all
				the Public User Identities associated to the Private User Identity. The Result-Code shall be set to DIAMETER_SUCCESS.

				If the HSS decides not to keep the S-CSCF name the Experimental-Result-Code shall be set to DIAMETER_SUCCESS_SERVER_NAME_NOT_STORED(2004).

				saa.append(AVP_Grouped(ProtocolConstants.DI_EXPERIMENTAL_RESULT,\
					[AVP_Unsigned32(ProtocolConstants.DI_VENDOR_ID,10415).setM(),\
					AVP_Unsigned32(ProtocolConstants.DI_EXPERIMENTAL_RESULT_CODE,\
					 #2006).setM()])) #is correct response as per https://www.rfc-editor.org/rfc/rfc4740.html#section-9.4
					 2004).setM()])) #from 3gpp https://www.arib.or.jp/english/html/overview/doc/STD-T63v9_20/5_Appendix/Rel6/29/29229-6b0.pdf

				*/

				stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-"+string(fmt.Sprintf("%d", datatype.Enumerated(req.ServerAssignmentType)))+"-2004")
				a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
						diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(2004)), //DIAMETER_SUCCESS_SERVER_NAME_NOT_STORED
					},
				})

			} else {
				//log.Printf("ServerAssignmentType Not received:")

				//ServerAssignmentType header is missing.. send error response
				a.NewAVP(avp.ExperimentalResult, avp.Mbit, 0, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(VENDOR_3GPP)),
						diam.NewAVP(avp.ExperimentalResultCode, avp.Mbit, 0, datatype.Unsigned32(5012)), //DIAMETER_UNABLE_TO_COMPLY
					},
				})
				stats.IncrementReceived("SAA", string(req.OriginHost), "sar-type-nil"+"-5012")
			}
		}

		//Adding ChargingInformation AVP
		//		a.NewAVP(avp.ChargingInformation, avp.Mbit, 0, &diam.GroupedAVP{
		//			AVP: []*diam.AVP{
		//				diam.NewAVP(avp.PrimaryChargingCollectionFunctionName, avp.Mbit, VENDOR_3GPP, datatype.UTF8String("pri_ccf_addr")),
		//			},
		//		})

		_, err = sendSAA_default(settings, c, a, enableLogging)
		if err != nil {
			//log.Printf("SAA sending ERROR %s", err.Error())

		}
	}
}

func sendSAA_default(settings sm.Settings, w io.Writer, m *diam.Message, enableLogging bool) (n int64, err error) {

	//m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))

	if enableLogging {
		log.Printf("inside sendSAA_default func")
		log.Printf("Sending SAA to \n%s", m)
	}

	return m.WriteTo(w)
}
