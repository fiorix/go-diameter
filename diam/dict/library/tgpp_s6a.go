package library

var TgppS6a = DictInfo{
	Name: "3GPP S6a",
	XML: `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
    <!--
        3GPP TS 29.272
        See: http://www.etsi.org/deliver/etsi_ts/129200_129299/129272/12.06.00_60/ts_129272v120600p.pdf
    -->
    <application id="16777251" type="auth" name="TGPP S6A">
        <vendor id="10415" name="TGPP"/>
        <command code="316" short="UL" name="Update-Location">
            <request>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="Destination-Host" required="false" max="1"/>
                <rule avp="Destination-Realm" required="true" max="1"/>
                <rule avp="User-Name" required="true" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="Terminal-Information" required="false" max="1"/>
                <rule avp="RAT-Type" required="true" max="1"/>
                <rule avp="ULR-Flags" required="true" max="1"/>
                <rule avp="UE-SRVCC-Capability" required="false" max="1"/>
                <rule avp="Visited-PLMN-Id" required="true" max="1"/>
                <rule avp="SGSN-Number" required="false" max="1"/>
                <rule avp="Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions" required="false" max="1"/>
                <rule avp="GMLC-Address" required="false" max="1"/>
                <rule avp="Active-APN" required="false"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </request>
            <answer>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Result-Code" required="false" max="1"/>
                <rule avp="Experimental-Result" required="false" max="1"/>
                <rule avp="Error-Diagnostic" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="ULA-Flags" required="false" max="1"/>
                <rule avp="Subscription-Data" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Failed-AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </answer>
        </command>

        <command code="317" short="CL" name="Cancel-Location">
            <request>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="Destination-Host" required="true" max="1"/>
                <rule avp="Destination-Realm" required="true" max="1"/>
                <rule avp="User-Name" required="true" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="Cancellation-Type" required="true" max="1"/>
                <rule avp="CLR-Flags" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </request>
            <answer>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="Result-Code" required="false" max="1"/>
                <rule avp="Experimental-Result" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Failed-AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </answer>
        </command>

        <command code="318" short="AI" name="Authentication-Information">
            <request>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="Destination-Host" required="false" max="1"/>
                <rule avp="Destination-Realm" required="true" max="1"/>
                <rule avp="User-Name" required="true" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="Requested-EUTRAN-Authentication-Info" required="false" max="1"/>
                <rule avp="Requested-UTRAN-GERAN-Authentication-Info" required="false" max="1"/>
                <rule avp="Visited-PLMN-Id" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </request>
            <answer>
                <rule avp="Session-Id" required="true" max="1"/>
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
                <rule avp="Result-Code" required="false" max="1"/>
                <rule avp="Experimental-Result" required="false" max="1"/>
                <rule avp="Error-Diagnostic" required="false" max="1"/>
                <rule avp="Auth-Session-State" required="true" max="1"/>
                <rule avp="Origin-Host" required="true" max="1"/>
                <rule avp="Origin-Realm" required="true" max="1"/>
                <rule avp="Supported-Features" required="false"/>
                <rule avp="Authentication-Info" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
                <rule avp="Failed-AVP" required="false"/>
                <rule avp="Proxy-Info" required="false"/>
                <rule avp="Route-Record" required="false"/>
            </answer>
        </command>

        <command code="321" short="PU" name="Purge-UE">
            <!--
                < Purge-UE-Request> ::=	< Diameter Header: 321, REQ, PXY, 16777251 >
                < Session-Id >
                [ DRMP ]
                [ Vendor-Specific-Application-Id ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                [ Destination-Host ]
                { Destination-Realm }
                { User-Name }
                [ OC-Supported-Features ]
                [ PUR-Flags ]
                *[ Supported-Features ]
                [ EPS-Location-Information ]
                *[ AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <request>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="DRMP" required="false" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="Destination-Host" required="false" max="1" />
                <rule avp="Destination-Realm" required="true" max="1" />
                <rule avp="User-Name" required="true" max="1" />
                <rule avp="OC-Supported-Features" required="false" max="1" />
                <rule avp="PUR-Flags" required="false" max="1" />
                <rule avp="Supported-Features" required="false" />
                <rule avp="EPS-Location-Information" required="false" max="1" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </request>

            <!--
                < Purge-UE-Answer> ::=	< Diameter Header: 321, PXY, 16777251 >
                < Session-Id >
                [ DRMP ]
                [ Vendor-Specific-Application-Id ]
                *[ Supported-Features ]
                [ Result-Code ]
                [ Experimental-Result ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                [ OC-Supported-Features ]
                [ OC-OLR ]
                *[ Load ]
                [ PUA-Flags ]
                *[ AVP ]
                [ Failed-AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <answer>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="DRMP" required="false" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="Supported-Features" required="false" />
                <rule avp="Result-Code" required="false" max="1" />
                <rule avp="Experimental-Result" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="OC-Supported-Features" required="false" max="1" />
                <rule avp="OC-OLR" required="false" max="1" />
                <!-- rule avp="Load" required="false" /-->
                <rule avp="PUA-Flags" required="false" max="1" />
                <rule avp="Failed-AVP" required="false" max="1" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </answer>
        </command>

        <command code="323" short="NO" name="Notify">
            <!--
                < Notify-Request> ::=	< Diameter Header: 323, REQ, PXY, 16777251 >
                < Session-Id >
                [ Vendor-Specific-Application-Id ]
                [ DRMP ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                [ Destination-Host ]
                { Destination-Realm }
                { User-Name }
                [ OC-Supported-Features ]
                * [ Supported-Features ]
                [ Terminal-Information ]
                [ MIP6-Agent-Info ]
                [ Visited-Network-Identifier ]
                [ Context-Identifier ]
                [ Service-Selection ]
                [ Alert-Reason ]
                [ UE-SRVCC-Capability ]
                [ NOR-Flags ]
                [ Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions ]
                [ Maximum-UE-Availability-Time ]
                *[ Monitoring-Event-Config-Status ]
                [ Emergency-Services ]
                *[ AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <request>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="DRMP" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="Destination-Host" required="false" max="1" />
                <rule avp="Destination-Realm" required="true" max="1" />
                <rule avp="User-Name" required="true" max="1" />
                <rule avp="OC-Supported-Features" required="false" max="1" />
                <rule avp="Supported-Features" required="false" />
                <rule avp="Terminal-Information" required="false" max="1" />
                <rule avp="MIP6-Agent-Info" required="false" max="1" />
                <rule avp="Visited-Network-Identifier" required="false" max="1" />
                <rule avp="Context-Identifier" required="false" max="1" />
                <rule avp="Service-Selection" required="false" max="1" />
                <rule avp="Alert-Reason" required="false" max="1" />
                <rule avp="UE-SRVCC-Capability" required="false" max="1" />
                <rule avp="NOR-Flags" required="false" max="1" />
                <rule avp="Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions" required="false" max="1" />
                <rule avp="Maximum-UE-Availability-Type" required="false" max="1" />
                <rule avp="Monitoring-Event-Config-Status" required="false" />
                <rule avp="Emergency-Services" required="false" max="1" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </request>
            <!--
                < Notify-Answer> ::=	< Diameter Header: 323, PXY, 16777251 >
                < Session-Id >
                [ DRMP ]
                [ Vendor-Specific-Application-Id ]
                [ Result-Code ]
                [ Experimental-Result ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                [ OC-Supported-Features ]
                [ OC-OLR ]
                *[ Load ]
                *[ Supported-Features ]
                *[ AVP ]
                [ Failed-AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <answer>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="DRMP" required="false" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="Result-Code" required="false" max="1" />
                <rule avp="Experimental-Result" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="OC-Supported-Features" required="false" max="1" />
                <rule avp="OC-OLR" required="false" max="1" />
                <!-- rule avp="Load" required="false" /-->
                <rule avp="Supported-Features" required="false" />
                <rule avp="Failed-AVP" required="false" max="1" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </answer>
        </command>

        <command code="322" short="RS" name="Reset">
            <!--
              < Reset-Request> ::= < Diameter Header: 322, REQ, PXY, 16777251 >
                < Session-Id >
                [ Vendor-Specific-Application-Id ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                { Destination-Host }
                { Destination-Realm }
                *[ Supported-Features ]
                *[ User-Id ]
                *[ AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <request>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="Destination-Host" required="false" max="1" />
                <rule avp="Destination-Realm" required="true" max="1" />
                <rule avp="Supported-Features" required="false" />
                <rule avp="User-Id" required="false" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </request>
            <!--
              < Reset-Answer> ::= < Diameter Header: 322, PXY, 16777251 >
                < Session-Id >
                [ Vendor-Specific-Application-Id ]
                *[ Supported-Features ]
                [ Result-Code ]
                [ Experimental-Result ]
                { Auth-Session-State }
                { Origin-Host }
                { Origin-Realm }
                *[ AVP ]
                *[ Failed-AVP ]
                *[ Proxy-Info ]
                *[ Route-Record ]
            -->
            <answer>
                <rule avp="Session-Id" required="true" max="1" />
                <rule avp="Vendor-Specific-Application-Id" required="false" max="1" />
                <rule avp="Supported-Features" required="false" />
                <rule avp="Result-Code" required="false" max="1" />
                <rule avp="Experimental-Result" required="false" max="1" />
                <rule avp="Auth-Session-State" required="true" max="1" />
                <rule avp="Origin-Host" required="true" max="1" />
                <rule avp="Origin-Realm" required="true" max="1" />
                <rule avp="Failed-AVP" required="false" max="1" />
                <rule avp="Proxy-Info" required="false" />
                <rule avp="Route-Record" required="false" />
            </answer>
        </command>

        <avp name="Subscription-Data" code="1400" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Subscriber-Status" required="false" max="1"/>
                <rule avp="MSISDN" required="false" max="1"/>
                <rule avp="STN-SR" required="false" max="1"/>
                <rule avp="ICS-Indicator" required="false" max="1"/>
                <rule avp="Network-Access-Mode" required="false" max="1"/>
                <rule avp="Operator-Determined-Barring" required="false" max="1"/>
                <rule avp="HPLMN-ODB" required="false" max="1"/>
                <rule avp="Regional-Subscription-Zone-Code" required="false" max="10"/>
                <rule avp="Access-Restriction-Data" required="false" max="1"/>
                <rule avp="APN-OI-Replacement" required="false" max="1"/>
                <rule avp="LCS-Info" required="false" max="1"/>
                <rule avp="Teleservice-List" required="false" max="1"/>
                <rule avp="Call-Barring-Info" required="false"/>
                <rule avp="TGPP-Charging-Characteristics" required="false" max="1"/>
                <rule avp="AMBR" required="false" max="1"/>
                <rule avp="APN-Configuration-Profile" required="false" max="1"/>
                <rule avp="RAT-Frequency-Selection-Priority-ID" required="false" max="1"/>
                <rule avp="Trace-Data" required="false" max="1"/>
                <rule avp="GPRS-Subscription-Data" required="false" max="1"/>
                <rule avp="CSG-Subscription-Data" required="false"/>
                <rule avp="Roaming-Restricted-Due-To-Unsupported-Feature" required="false" max="1"/>
                <rule avp="Subscribed-Periodic-RAU-TAU-Timer" required="false" max="1"/>
                <rule avp="MPS-Priority" required="false" max="1"/>
                <rule avp="VPLMN-LIPA-Allowed" required="false" max="1"/>
                <rule avp="Relay-Node-Indicator" required="false" max="1"/>
                <rule avp="MDT-User-Consent" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Subscriber-Status" code="1424" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="SERVICE_GRANTED"/>
                <item code="1" name="OPERATOR_DETERMINED_BARRING"/>
            </data>
        </avp>

        <avp name="STN-SR" code="1433" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="ICS-Indicator" code="1491" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="FALSE"/>
                <item code="1" name="TRUE"/>
            </data>
        </avp>

        <avp name="Network-Access-Mode" code="1417" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="PACKET_AND_CIRCUIT"/>
                <item code="1" name="RESERVED"/>
                <item code="2" name="ONLY_PACKET"/>
            </data>
        </avp>

        <avp name="Operator-Determined-Barring" code="1425" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="HPLMN-ODB" code="1418" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="Regional-Subscription-Zone-Code" code="1446" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Access-Restriction-Data" code="1426" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="APN-OI-Replacement" code="1427" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="UTF8String"/>
        </avp>

        <avp name="LCS-Info" code="1473" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="GMLC-Number" required="false"/>
                <rule avp="LCS-PrivacyException" required="false"/>
                <rule avp="MO-LR" required="false"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>avp>

        <avp name="GMLC-Number" code="1474" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="LCS-PrivacyException" code="1475" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="SS-Code" required="true" max="1"/>
                <rule avp="SS-Status" required="true" max="1"/>
                <rule avp="Notification-To-UE-User" required="false" max="1"/>
                <rule avp="External-Client" required="false"/>
                <rule avp="PLMN-Client" required="false"/>
                <rule avp="TGPP-Service-Type" required="false"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>avp>

        <avp name="SS-Code" code="1476" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="SS-Status" code="1477" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Notification-To-UE-User" code="1478" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="NOTIFY_LOCATION_ALLOWED"/>
                <item code="1" name="NOTIFYANDVERIFY_LOCATION_ALLOWED_IF_NO_RESPONSE"/>
                <item code="2" name="NOTIFYANDVERIFY_LOCATION_NOT_ALLOWED_IF_NO_RESPONSE"/>
                <item code="3" name="LOCATION_NOT_ALLOWED"/>
            </data>
        </avp>

        <avp name="External-Client" code="1479" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Client-Identity" required="true" max="1"/>
                <rule avp="GMLC-Restriction" required="false" max="1"/>
                <rule avp="Notification-To-UE-User" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>avp>

        <avp name="Client-Identity" code="1480" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="GMLC-Restriction" code="1481" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="GMLC_LIST"/>
                <item code="1" name="HOME_COUNTRY"/>
            </data>
        </avp>

        <avp name="PLMN-Client" code="1482" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="BROADCAST_SERVICE"/>
                <item code="1" name="O_AND_M_HPLMN"/>
                <item code="2" name="O_AND_M_VPLMN"/>
                <item code="3" name="ANONYMOUS_LOCATION"/>
                <item code="3" name="TARGET_UE_SUBSCRIBED_SERVICE"/>
            </data>
        </avp>

        <avp name="TGPP-Service-Type" code="1483" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="ServiceTypeIdentity" required="true" max="1"/>
                <rule avp="GMLC-Restriction" required="false" max="1"/>
                <rule avp="Notification-To-UE-User" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="ServiceTypeIdentity" code="1484" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="MO-LR" code="1485" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="SS-Code" required="true" max="1"/>
                <rule avp="SS-Status" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Teleservice-List" code="1486" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="TS-Code" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="TS-Code" code="1487" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Call-Barring-Info" code="1488" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="SS-Code" required="true" max="1"/>
                <rule avp="SS-Status" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="AMBR" code="1435" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Max-Requested-Bandwidth-UL" required="true" max="1"/>
                <rule avp="Max-Requested-Bandwidth-DL" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="APN-Configuration-Profile" code="1429" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Context-Identifier" required="true" max="1"/>
                <rule avp="All-APN-Configurations-Included-Indicator" required="true" max="1"/>
                <rule avp="APN-Configuration" required="true"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Context-Identifier" code="1423" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="All-APN-Configurations-Included-Indicator" code="1428" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="All_APN_CONFIGURATIONS_INCLUDED"/>
                <item code="1" name="MODIFIED|ADDED_APN_CONFIGURATIONS_INCLUDED"/>
            </data>
        </avp>

        <avp name="APN-Configuration" code="1430" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Context-Identifier" required="true" max="1"/>
                <rule avp="Served-Party-IP-Address" required="false" max="2"/>
                <rule avp="PDN-Type" required="true" max="1"/>
                <rule avp="Service-Selection" required="true" max="1"/>
                <rule avp="EPS-Subscribed-QoS-Profile" required="false" max="1"/>
                <rule avp="VPLMN-Dynamic-Address-Allowed" required="false" max="1"/>
                <rule avp="MIP6-Agent-Info" required="false" max="1"/>
                <rule avp="Visited-Network-Identifier" required="false" max="1"/>
                <rule avp="PDN-GW-Allocation-Type" required="false" max="1"/>
                <rule avp="TGPP-Charging-Characteristics" required="false" max="1"/>
                <rule avp="AMBR" required="false" max="1"/>
                <rule avp="Specific-APN-Info" required="false"/>
                <rule avp="APN-OI-Replacement" required="false" max="1"/>
                <rule avp="SIPTO-Permission" required="false" max="1"/>
                <rule avp="LIPA-Permission" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Served-Party-IP-Address" code="848" must="M,V" may="P" may-encrypt="N" vendor-id="10415">
            <data type="Address"/>
        </avp>

        <avp name="PDN-Type" code="1456" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="IPv4"/>
                <item code="1" name="IPv6"/>
                <item code="2" name="IPv4v6"/>
                <item code="3" name="IPv4_OR_IPv6"/>
            </data>
        </avp>

        <avp name="EPS-Subscribed-QoS-Profile" code="1431" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="QoS-Class-Identifier" required="true" max="1"/>
                <rule avp="Allocation-Retention-Priority" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="VPLMN-Dynamic-Address-Allowed" code="1432" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="NOTALLOWED"/>
                <item code="1" name="ALLOWED"/>
            </data>
        </avp>

        <avp name="PDN-GW-Allocation-Type" code="1438" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="STATIC"/>
                <item code="1" name="DYNAMIC"/>
            </data>
        </avp>

        <avp name="SIPTO-Permission" code="1613" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="SIPTO_ALLOWED"/>
                <item code="1" name="SIPTO_NOTALLOWED"/>
            </data>
        </avp>

        <avp name="LIPA-Permission" code="1618" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="LIPA-PROHIBITED"/>
                <item code="1" name="LIPA-ONLY"/>
                <item code="2" name="LIPA-CONDITIONAL"/>
            </data>
        </avp>

        <avp name="RAT-Frequency-Selection-Priority-ID" code="1440" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="Trace-Data" code="1458" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Trace-Reference" required="true" max="1"/>
                <rule avp="Trace-Depth" required="true" max="1"/>
                <rule avp="Trace-NE-Type-List" required="true" max="1"/>
                <rule avp="Trace-Interface-List" required="false" max="1"/>
                <rule avp="Trace-Event-List" required="true" max="1"/>
                <rule avp="OMC-Id" required="false" max="1"/>
                <rule avp="Trace-Collection-Entity" required="true" max="1"/>
                <rule avp="MDT-Configuration" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>avp>

        <avp name="Trace-Reference" code="1459" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Trace-Depth" code="1462" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="LIPA-PROHIBITED"/>
                <item code="1" name="LIPA-ONLY"/>
                <item code="2" name="LIPA-CONDITIONAL"/>
            </data>
        </avp>

        <avp name="Trace-NE-Type-List" code="1463" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Trace-Interface-List" code="1464" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Trace-Event-List" code="1465" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="OMC-Id" code="1466" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Trace-Event-List" code="1465" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Trace-Collection-Entity" code="1452" must="M,V" may="P" may-encrypt="N" vendor-id="10415">
            <data type="Address"/>
        </avp>

        <avp name="MDT-Configuration" code="1622" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="QoS-Class-Identifier" required="true" max="1"/>
                <rule avp="Allocation-Retention-Priority" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="GPRS-Subscription-Data" code="1467" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Complete-Data-List-Included-Indicator" required="true" max="1"/>
                <rule avp="PDP-Context" required="true" max="50"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Complete-Data-List-Included-Indicator" code="1468" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="All_PDP_CONTEXTS_INCLUDED"/>
                <item code="1" name="MODIFIED/ADDED_PDP CONTEXTS_INCLUDED"/>
            </data>
        </avp>

        <avp name="PDP-Context" code="1469" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Context-Identifier" required="true" max="1"/>
                <rule avp="PDP-Type" required="true" max="1"/>
                <rule avp="PDP-Address" required="false" max="1"/>
                <rule avp="QoS-Subscribed" required="true" max="1"/>
                <rule avp="VPLMN-Dynamic-Address-Allowed" required="false" max="1"/>
                <rule avp="Service-Selection" required="true" max="1"/>
                <rule avp="TGPP-Charging-Characteristics" required="false" max="1"/>
                <rule avp="Ext-PDP-Type" required="false" max="1"/>
                <rule avp="Ext-PDP-Address" required="false" max="1"/>
                <rule avp="AMBR" required="false" max="1"/>
                <rule avp="SIPTO-Permission" required="false" max="1"/>
                <rule avp="LIPA-Permission" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="PDP-Type" code="1470" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="QoS-Subscribed" code="1404" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="VPLMN-Dynamic-Address-Allowed" code="1432" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="NOTALLOWED"/>
                <item code="1" name="ALLOWED"/>
            </data>
        </avp>

        <avp name="Ext-PDP-Type" code="1620" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Ext-PDP-Address" code="1621" must="V,M" may="P" must-not="-" may-encrypt="Y" vendor-id="10415">
            <data type="Address"/>
        </avp>

        <avp name="SIPTO-Permission" code="1613" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="SIPTO_ALLOWED"/>
                <item code="1" name="SIPTO_NOTALLOWED"/>
            </data>
        </avp>

        <avp name="CSG-Subscription-Data" code="1436" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="CSG-Id" required="true" max="1"/>
                <rule avp="Expiration-Date" required="false" max="1"/>
                <rule avp="Service-Selection" required="false"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Expiration-Date" code="1439" must="V,M" may-encrypt="N" vendor-id="10415">
            <data type="Time"/>
        </avp>

        <avp name="Roaming-Restricted-Due-To-Unsupported-Feature" code="1457" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="Roaming-Restricted-Due-To-Unsupported-Feature"/>
            </data>
        </avp>

        <avp name="Subscribed-Periodic-RAU-TAU-Timer" code="1619" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="MPS-Priority" code="1616" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="VPLMN-LIPA-Allowed" code="1617" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="LIPA-NOTALLOWED"/>
                <item code="1" name="LIPA-ALLOWED"/>
            </data>
        </avp>

        <avp name="Relay-Node-Indicator" code="1633" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="NOT_RELAY_NODE"/>
                <item code="1" name="RELAY_NODE"/>
            </data>
        </avp>

        <avp name="MDT-User-Consent" code="1634" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="CONSENT_NOT_GIVEN"/>
                <item code="1" name="CONSENT_GIVEN"/>
            </data>
        </avp>

        <avp name="Requested-EUTRAN-Authentication-Info" code="1408" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Number-Of-Requested-Vectors" required="false" max="1"/>
                <rule avp="Immediate-Response-Preferred" required="false" max="1"/>
                <rule avp="Re-synchronization-Info" required="false" max="1"/>
            </data>
        </avp>

        <avp name="Number-Of-Requested-Vectors" code="1410" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="Re-synchronization-Info" code="1411" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Immediate-Response-Preferred" code="1412" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="Requested-UTRAN-GERAN-Authentication-Info" code="1409" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Number-Of-Requested-Vectors" required="false" max="1"/>
                <rule avp="Immediate-Response-Preferred" required="false" max="1"/>
                <rule avp="Re-synchronization-Info" required="false" max="1"/>
            </data>
        </avp>

        <avp name="Visited-PLMN-Id" code="1407" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Error-Diagnostic" code="1614" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="GPRS_DATA_SUBSCRIBED"/>
                <item code="1" name="NO_GPRS_DATA_SUBSCRIBED"/>
                <item code="2" name="ODB-ALL-APN"/>
                <item code="3" name="ODB-HPLMN-APN"/>
                <item code="4" name="ODB-VPLMN-APN"/>
            </data>
        </avp>

        <avp name="Authentication-Info" code="1413" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="E-UTRAN-Vector" required="false"/>
                <rule avp="UTRAN-Vector" required="false"/>
                <rule avp="GERAN-Vector" required="false"/>
            </data>
        </avp>

        <avp name="E-UTRAN-Vector" code="1414" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Item-Number" required="false" max="1"/>
                <rule avp="RAND" required="true" max="1"/>
                <rule avp="XRES" required="true" max="1"/>
                <rule avp="AUTN" required="true" max="1"/>
                <rule avp="KASME" required="true" max="1"/>
            </data>
        </avp>

        <avp name="UTRAN-Vector" code="1415" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Item-Number" required="false"/>
                <rule avp="RAND" required="true"/>
                <rule avp="XRES" required="true"/>
                <rule avp="AUTN" required="true"/>
                <rule avp="Confidentiality-Key" required="true"/>
                <rule avp="Integrity-Key" required="true"/>
            </data>
        </avp>

        <avp name="GERAN-Vector" code="1416" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Item-Number" required="false" max="1"/>
                <rule avp="RAND" required="true" max="1"/>
                <rule avp="SRES" required="true" max="1"/>
                <rule avp="Kc" required="true" max="1"/>
            </data>
        </avp>

        <avp name="Item-Number" code="1419" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="Cancellation-Type" code="1420" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="MME_UPDATE_PROCEDURE"/>
                <item code="1" name="SGSN_UDPATE_PROCEDURE"/>
                <item code="2" name="SUBSCRIPTION_WITHDRAWAL"/>
                <item code="3" name="UPDATE_PROCEDURE_IWF"/>
                <item code="4" name="INITIAL_ATTACH_PROCEDURE"/>
            </data>
        </avp>

        <avp name="RAND" code="1447" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="XRES" code="1448" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="AUTN" code="1449" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="KASME" code="1450" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Kc" code="1453" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="SRES" code="1454" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Confidentiality-Key" code="625" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="Integrity-Key" code="626" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="ULR-Flags" code="1405" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="ULA-Flags" code="1406" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="CLR-Flags" code="1638" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="UE-SRVCC-Capability" code="1615" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="UE-SRVCC-NOT-SUPPORTED"/>
                <item code="1" name="UE-SRVCC-SUPPORTED"/>
            </data>
        </avp>

        <avp name="Homogeneous-Support-of-IMS-Voice-Over-PS-Sessions" code="1493" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="NOT-SUPPORTED"/>
                <item code="1" name="SUPPORTED"/>
            </data>
        </avp>

        <avp name="Active-APN" code="1612" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Context-Identifier" required="true" max="1"/>
                <rule avp="Service-Selection" required="false" max="1"/>
                <rule avp="MIP6-Agent-Info" required="false" max="1"/>
                <rule avp="Visited-Network-Identifier" required="false" max="1"/>
                <rule avp="Specific-APN-Info" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Specific-APN-Info" code="1472" vendor-id="10415" must="M,V" may-encrypt="N">
            <data type="Grouped">
                <rule avp="Service-Selection" required="true" max="1"/>
                <rule avp="MIP6-Agent-Info" required="true" max="1"/>
                <rule avp="Visited-Network-Identifier" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Context-Identifier" code="1423" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32"/>
        </avp>

        <avp name="PUR-Flags" code="1635" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32" />
        </avp>

        <avp name="PUA-Flags" code="1442" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32" />
        </avp>

        <avp name="NOR-Flags" code="1443" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="Unsigned32" />
        </avp>

        <avp name="Subscribed-VSRVCC" code="1636" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="Enumerated">
                <item code="0" name="VSRVCC_SUBSCRIBED" />
            </data>
        </avp>

        <avp name="MIP-Home-Agent-Address" code="334" must="M" must-not="V" vendor-id="10415">
            <data type="Address"/>
        </avp>

        <!-- RFC 4004 -->
        <avp name="MIP-Home-Agent-Host" code="348" must="M" may="P" must-not="V" may-encrypt="Y" vendor-id="10415">
            <data type="Grouped">
                <rule avp="Destination-Realm" required="true" max="1"/>
                <rule avp="Destination-Host" required="true" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <!-- RFC 5447 Diameter Mobile IPv6: Support for Network Access Server to Diameter Server Interaction -->
        <avp name="MIP6-Agent-Info" code="486" must="M" may="P" must-not="V" may-encrypt="Y" vendor-id="10415">
            <data type="Grouped">
                <rule avp="MIP-Home-Agent-Address" required="false" max="2"/>
                <rule avp="MIP-Home-Agent-Host" required="false" max="1"/>
                <rule avp="MIP6-Home-Link-Prefix" required="false" max="1"/>
                <rule avp="AVP" required="false"/>
            </data>
        </avp>

        <avp name="Service-Selection" code="493" must="M" may="P" must-not="V" may-encrypt="Y" vendor-id="10415">
            <data type="UTF8String"/>
        </avp>

        <avp name="Visited-Network-Identifier" code="600" must="M,V" may-encrypt="N" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="MIP6-Home-Link-Prefix" code="125" vendor-id="10415">
            <data type="OctetString"/>
        </avp>

        <avp name="User-Id" code="1444" must="V" must-not="M" may-encrypt="N" vendor-id="10415">
            <data type="UTF8String"/>
        </avp>

    </application>
</diameter>`,
}
