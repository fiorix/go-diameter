// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package dict

import "bytes"

// Default is a Parser object with pre-loaded
// Base Protocol and Credit Control dictionaries.
var Default *Parser

func init() {
	Default, _ = NewParser()
	Default.Load(bytes.NewReader([]byte(baseXML)))
	Default.Load(bytes.NewReader([]byte(creditcontrolXML)))
}

var baseXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="0"> <!-- Diameter Common Messages -->

    <vendor id="10415" name="3GPP"/>

    <command code="257" short="CE" name="Capabilities-Exchange">
      <request>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Host-IP-Address" required="true" max="1"/>
        <rule avp="Vendor-Id" required="true" max="1"/>
        <rule avp="Product-Name" required="true" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Supported-Vendor-Id" required="False"/>
        <rule avp="Auth-Application-Id" required="False"/>
        <rule avp="Inband-Security-Id" required="False"/>
        <rule avp="Acct-Application-Id" required="False"/>
        <rule avp="Vendor-Specific-Application-Id" required="False"/>
        <rule avp="Firmware-Revision" required="False" max="1"/>
      </request>
      <answer>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Host-IP-Address" required="true" max="1"/>
        <rule avp="Vendor-Id" required="true" max="1"/>
        <rule avp="Product-Name" required="true" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Supported-Vendor-Id" required="False"/>
        <rule avp="Auth-Application-Id" required="False"/>
        <rule avp="Inband-Security-Id" required="False"/>
        <rule avp="Acct-Application-Id" required="False"/>
        <rule avp="Vendor-Specific-Application-Id" required="False"/>
        <rule avp="Firmware-Revision" required="False" max="1"/>
      </answer>
    </command>

    <command code="258" short="RA" name="Re-Auth">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Destination-Host" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Re-Auth-Request-Type" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
        <rule avp="Route-Record" required="false"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Redirect-Host" required="false"/>
        <rule avp="Redirect-Host-Usage" required="false" max="1"/>
        <rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
      </answer>
    </command>

    <command code="271" short="AC" name="Accounting">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Accounting-Record-Type" required="true" max="1"/>
        <rule avp="Accounting-Record-Number" required="true" max="1"/>
        <rule avp="Acct-Application-Id" required="false" max="1"/>
        <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Destination-Host" required="false" max="1"/>
        <rule avp="Accounting-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Interim-Interval" required="false" max="1"/>
        <rule avp="Accounting-Realtime-Required" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
        <rule avp="Route-Record" required="false"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Accounting-Record-Type" required="true" max="1"/>
        <rule avp="Accounting-Record-Number" required="true" max="1"/>
        <rule avp="Acct-Application-Id" required="false" max="1"/>
        <rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Accounting-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Acct-Interim-Interval" required="false" max="1"/>
        <rule avp="Accounting-Realtime-Required" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
      </answer>
    </command>

    <command code="274" short="AS" name="Abort-Session">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Destination-Host" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
        <rule avp="Route-Record" required="false"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Redirect-Host" required="false"/>
        <rule avp="Redirect-Host-Usage" required="false" max="1"/>
        <rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
      </answer>
    </command>

    <command code="275" short="ST" name="Session-Termination">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Termination-Cause" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Destination-Host" required="false" max="1"/>
        <rule avp="Class" required="false"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
        <rule avp="Route-Record" required="false"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="Class" required="false"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Error-Reporting-Host" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Redirect-Host" required="false"/>
        <rule avp="Redirect-Host-Usage" required="false" max="1"/>
        <rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false"/>
      </answer>
    </command>

    <command code="280" short="DW" name="Device-Watchdog">
      <request>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
      </request>
      <answer>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
      </answer>
    </command>

    <command code="282" short="DP" name="Disconnect-Peer">
      <request>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Disconnect-Cause" required="false" max="1"/>
      </request>
      <answer>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Error-Message" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
      </answer>
    </command>

    <avp name="Acct-Interim-Interval" code="85" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Accounting-Realtime-Required" code="483" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="1" name="DELIVER_AND_GRANT"/>
        <item code="2" name="GRANT_AND_STORE"/>
        <item code="3" name="GRANT_AND_LOSE"/>
      </data>
    </avp>

    <avp name="Acct-Multi-Session-Id" code="50" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Accounting-Record-Number" code="485" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Accounting-Record-Type" code="480" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="1" name="EVENT_RECORD"/>
        <item code="2" name="START_RECORD"/>
        <item code="3" name="INTERIM_RECORD"/>
        <item code="4" name="STOP_RECORD"/>
      </data>
    </avp>

    <avp name="Accounting-Session-Id" code="44" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Accounting-Sub-Session-Id" code="287" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="Acct-Application-Id" code="259" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Application-Id" code="258" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Request-Type" code="274" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="1" name="AUTHENTICATE_ONLY"/>
        <item code="2" name="AUTHORIZE_ONLY"/>
        <item code="3" name="AUTHORIZE_AUTHENTICATE"/>
      </data>
    </avp>

    <avp name="Authorization-Lifetime" code="291" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Grace-Period" code="276" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Session-State" code="277" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="0" name="STATE_MAINTAINED"/>
        <item code="1" name="NO_STATE_MAINTAINED"/>
      </data>
    </avp>

    <avp name="Re-Auth-Request-Type" code="285" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="0" name="AUTHORIZE_ONLY"/>
        <item code="1" name="AUTHORIZE_AUTHENTICATE"/>
      </data>
    </avp>

    <avp name="Class" code="25" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Destination-Host" code="293" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Destination-Realm" code="283" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Disconnect-Cause" code="273" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="0" name="REBOOTING"/>
        <item code="1" name="BUSY"/>
        <item code="2" name="DO_NOT_WANT_TO_TALK_TO_YOU"/>
      </data>
    </avp>

    <avp name="Error-Message" code="281" must="-" may="P" must-not="V,M" may-encrypt="-">
      <data type="UTF8String"/>
    </avp>

    <avp name="Error-Reporting-Host" code="294" must="-" may="P" must-not="V,M" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Event-Timestamp" code="55" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Time"/>
    </avp>

    <avp name="Experimental-Result" code="297" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Grouped">
        <rule avp="Vendor-Id" required="true" max="1"/>
        <rule avp="Experimental-Result-Code" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Experimental-Result-Code" code="298" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Failed-AVP" code="279" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Grouped"/>
    </avp>

    <avp name="Firmware-Revision" code="267" must="-" may="-" must-not="P,V,M" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Host-IP-Address" code="257" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Address"/>
    </avp>

    <avp name="Inband-Security-Id" code="299" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Multi-Round-Time-Out" code="272" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Origin-Host" code="264" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Origin-Realm" code="296" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Origin-State-Id" code="278" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Product-Name" code="269" must="-" may="-" must-not="P,V,M" may-encrypt="-">
      <data type="UTF8String"/>
    </avp>

    <avp name="Proxy-Host" code="280" must="M" may="-" must-not="P,V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Proxy-Info" code="284" must="M" may="-" must-not="P,V" may-encrypt="-">
      <data type="Grouped">
        <rule avp="Proxy-Host" required="true" max="1"/>
        <rule avp="Proxy-State" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Proxy-State" code="33" must="M" may="-" must-not="P,V" may-encrypt="-">
      <data type="OctetString"/>
    </avp>

    <avp name="Redirect-Host" code="292" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="DiameterURI"/>
    </avp>

    <avp name="Redirect-Host-Usage" code="261" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="0" name="DONT_CACHE"/>
        <item code="1" name="ALL_SESSION"/>
        <item code="2" name="ALL_REALM"/>
        <item code="3" name="REALM_AND_APPLICATION"/>
        <item code="4" name="ALL_APPLICATION"/>
        <item code="5" name="ALL_HOST"/>
        <item code="6" name="ALL_USER"/>
      </data>
    </avp>

    <avp name="Redirect-Max-Cache-Time" code="262" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Result-Code" code="268" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Route-Record" code="282" must="M" may="-" must-not="P,V" may-encrypt="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Session-Id" code="263" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Session-Timeout" code="27" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Session-Binding" code="270" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Session-Server-Failover" code="271" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="REFUSE_SERVICE"/>
        <item code="1" name="TRY_AGAIN"/>
        <item code="2" name="ALLOW_SERVICE"/>
        <item code="3" name="TRY_AGAIN_ALLOW_SERVICE"/>
      </data>
    </avp>

    <avp name="Supported-Vendor-Id" code="265" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Termination-Cause" code="295" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Enumerated">
        <item code="1" name="DIAMETER_LOGOUT"/>
        <item code="2" name="DIAMETER_SERVICE_NOT_PROVIDED"/>
        <item code="3" name="DIAMETER_BAD_ANSWER"/>
        <item code="4" name="DIAMETER_ADMINISTRATIVE"/>
        <item code="5" name="DIAMETER_LINK_BROKEN"/>
        <item code="6" name="DIAMETER_AUTH_EXPIRED"/>
        <item code="7" name="DIAMETER_USER_MOVED"/>
        <item code="8" name="DIAMETER_SESSION_TIMEOUT"/>
      </data>
    </avp>

    <avp name="User-Name" code="1" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Vendor-Id" code="266" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Vendor-Specific-Application-Id" code="260" must="M" may="P" must-not="V" may-encrypt="-">
      <data type="Grouped">
        <rule avp="Vendor-Id" required="false" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Acct-Application-Id" required="true" max="1"/>
      </data>
    </avp>

  </application>
</diameter>`

var creditcontrolXML = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="4">
    <!-- Diameter Credit Control Application -->
    <!-- http://tools.ietf.org/html/rfc4006 -->

    <vendor id="10415" name="3GPP"/>

    <command code="272" short="CC" name="Credit-Control">
      <request>
        <!-- http://tools.ietf.org/html/rfc4006#section-3.1 -->
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Service-Context-Id" required="true" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="Destination-Host" required="false" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="CC-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Subscription-Id" required="false" max="1"/>
        <rule avp="Service-Identifier" required="false" max="1"/>
        <rule avp="Termination-Cause" required="false" max="1"/>
        <rule avp="Requested-Service-Unit" required="false" max="1"/>
        <rule avp="Requested-Action" required="false" max="1"/>
        <rule avp="Used-Service-Unit" required="false" max="1"/>
        <rule avp="Multiple-Services-Indicator" required="false" max="1"/>
        <rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
        <rule avp="Service-Parameter-Info" required="false" max="1"/>
        <rule avp="CC-Correlation-Id" required="false" max="1"/>
        <rule avp="User-Equipment-Info" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </request>
      <answer>
        <!-- http://tools.ietf.org/html/rfc4006#section-3.2 -->
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="CC-Session-Failover" required="false" max="1"/>
        <rule avp="CC-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Granted-Service-Unit" required="false" max="1"/>
        <rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
        <rule avp="Cost-Information" required="false" max="1"/>
        <rule avp="Final-Unit-Indication" required="false" max="1"/>
        <rule avp="Check-Balance-Result" required="false" max="1"/>
        <rule avp="Credit-Control-Failure-Handling" required="false" max="1"/>
        <rule avp="Direct-Debiting-Failure-Handling" required="false" max="1"/>
        <rule avp="Validity-Time" required="false" max="1"/>
        <rule avp="Redirect-Host" required="false" max="1"/>
        <rule avp="Redirect-Host-Usage" required="false" max="1"/>
        <rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
      </answer>
    </command>

    <avp name="CC-Correlation-Id" code="411" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.1 -->
      <data type="OctetString"/>
    </avp>

    <avp name="CC-Input-Octets" code="412" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.24 -->
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Money" code="413" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.22 -->
      <data type="Grouped">
        <rule avp="Unit-Value" required="true" max="1"/>
        <rule avp="Currency-Code" required="true" max="1"/>
      </data>
    </avp>

    <avp name="CC-Output-Octets" code="414" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.25 -->
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Request-Number" code="415" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.2 -->
      <data type="Unsigned32"/>
    </avp>

    <avp name="CC-Request-Type" code="416" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.3 -->
      <data type="Enumerated">
        <item code="1" name="INITIAL_REQUEST"/>
        <item code="2" name="UPDATE_REQUEST"/>
        <item code="3" name="TERMINATION_REQUEST"/>
      </data>
    </avp>

    <avp name="CC-Service-Specific-Units" code="417" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Session-Failover" code="418" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.4 -->
      <data type="Enumerated">
        <item code="0" name="FAILOVER_NOT_SUPPORTED"/>
        <item code="1" name="FAILOVER_SUPPORTED"/>
      </data>
    </avp>

    <avp name="CC-Sub-Session-Id" code="419" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.5 -->
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Time" code="420" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.21 -->
      <data type="Unsigned32"/>
    </avp>

    <avp name="CC-Total-Octets" code="421" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.23 -->
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Unit-Type" code="454" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.32 -->
      <data type="Enumerated">
        <item code="0" name="TIME"/>
        <item code="1" name="MONEY"/>
        <item code="2" name="TOTAL-OCTETS"/>
        <item code="3" name="INPUT-OCTETS"/>
        <item code="4" name="OUTPUT-OCTETS"/>
        <item code="5" name="SERVICE-SPECIFIC-UNITS"/>
      </data>
    </avp>

    <avp name="Check-Balance-Result" code="422" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.6 -->
      <data type="Enumerated">
        <item code="0" name="ENOUGH_CREDIT"/>
        <item code="1" name="NO_CREDIT"/>
      </data>
    </avp>

    <avp name="Cost-Information" code="423" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.7 -->
      <data type="Grouped">
        <rule avp="Unit-Value" required="true" max="1"/>
        <rule avp="Currency-Code" required="true" max="1"/>
        <rule avp="Cost-Unit" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Cost-Unit" code="424" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.12 -->
      <data type="UTF8String"/>
    </avp>

    <avp name="Credit-Control" code="426" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.13 -->
      <data type="Enumerated">
        <item code="0" name="CREDIT_AUTHORIZATION"/>
        <item code="1" name="RE_AUTHORIZATION"/>
      </data>
    </avp>

    <avp name="Credit-Control-Failure-Handling" code="427" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.14 -->
      <data type="Enumerated">
        <item code="0" name="TERMINATE"/>
        <item code="1" name="CONTINUE"/>
        <item code="2" name="RETRY_AND_TERMINATE"/>
      </data>
    </avp>

    <avp name="Currency-Code" code="425" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.11 -->
      <data type="Unsigned32"/>
    </avp>

    <avp name="Direct-Debiting-Failure-Handling" code="428" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.15 -->
      <data type="Enumerated">
        <item code="0" name="TERMINATE_OR_BUFFER"/>
        <item code="1" name="CONTINUE"/>
      </data>
    </avp>

    <avp name="Exponent" code="429" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.9 -->
      <data type="Integer32"/>
    </avp>

    <avp name="Final-Unit-Action" code="449" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.35 -->
      <data type="Enumerated">
        <item code="0" name="TERMINATE"/>
        <item code="1" name="REDIRECT"/>
        <item code="2" name="RESTRICT_ACCESS"/>
      </data>
    </avp>

    <avp name="Final-Unit-Indication" code="430" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.34 -->
      <data type="Grouped">
        <rule avp="Final-Unit-Action" required="true" max="1"/>
        <rule avp="Restriction-Filter-Rule" required="false" max="1"/>
        <rule avp="Filter-Id" required="false" max="1"/>
        <rule avp="Redirect-Server" required="false" max="1"/>
      </data>
    </avp>

    <avp name="Granted-Service-Unit" code="431" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.17 -->
      <data type="Grouped">
        <rule avp="Tariff-Time-Change" required="false" max="1"/>
        <rule avp="CC-Time" required="false" max="1"/>
        <rule avp="CC-Money" required="false" max="1"/>
        <rule avp="CC-Total-Octets" required="false" max="1"/>
        <rule avp="CC-Input-Octets" required="false" max="1"/>
        <rule avp="CC-Output-Octets" required="false" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="false" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="G-S-U-Pool-Identifier" code="453" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.31 -->
      <data type="Unsigned32"/>
    </avp>

    <avp name="G-S-U-Pool-Reference" code="457" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.30 -->
      <data type="Grouped">
        <rule avp="G-S-U-Pool-Identifier" required="true" max="1"/>
        <rule avp="CC-Unit-Type" required="true" max="1"/>
        <rule avp="Unit-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Multiple-Services-Credit-Control" code="456" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.16 -->
      <data type="Grouped">
        <rule avp="Granted-Service-Unit" required="false" max="1"/>
        <rule avp="Requested-Service-Unit" required="false" max="1"/>
        <rule avp="Used-Service-Unit" required="false" max="1"/>
        <rule avp="Tariff-Change-Usage" required="false" max="1"/>
        <rule avp="Service-Identifier" required="false" max="1"/>
        <rule avp="Rating-Group" required="false" max="1"/>
        <rule avp="G-S-U-Pool-Reference" required="false" max="1"/>
        <rule avp="Validity-Time" required="false" max="1"/>
        <rule avp="Result-Code" required="false" max="1"/>
        <rule avp="Final-Unit-Indication" required="false" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="Multiple-Services-Indicator" code="455" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.40 -->
      <data type="Enumerated">
        <item code="0" name="MULTIPLE_SERVICES_NOT_SUPPORTED"/>
        <item code="1" name="MULTIPLE_SERVICES_SUPPORTED"/>
      </data>
    </avp>

    <avp name="Rating-Group" code="432" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.29 -->
      <data type="Unsigned32"/>
    </avp>

    <avp name="Redirect-Address-Type " code="433" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.38 -->
      <data type="Enumerated">
        <item code="0" name="IPv4 Address"/>
        <item code="1" name="IPv6 Address"/>
        <item code="2" name="URL"/>
        <item code="3" name="SIP URI"/>
      </data>
    </avp>

    <avp name="Redirect-Server" code="434" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.37 -->
      <data type="Grouped">
        <rule avp="Redirect-Address-Type" required="true" max="1"/>
        <rule avp="Redirect-Server-Address" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Redirect-Server-Address" code="435" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.39 -->
      <data type="UTF8String"/>
    </avp>

    <avp name="Requested-Action" code="436" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.41 -->
      <data type="Enumerated">
        <item code="0" name="DIRECT_DEBITING"/>
        <item code="1" name="REFUND_ACCOUNT"/>
        <item code="2" name="CHECK_BALANCE"/>
        <item code="3" name="PRICE_ENQUIRY"/>
      </data>
    </avp>

    <avp name="Requested-Service-Unit" code="437" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.18-->
      <data type="Grouped">
        <rule avp="CC-Time" required="false" max="1"/>
        <rule avp="CC-Money" required="false" max="1"/>
        <rule avp="CC-Total-Octets" required="false" max="1"/>
        <rule avp="CC-Input-Octets" required="false" max="1"/>
        <rule avp="CC-Output-Octets" required="false" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="false" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="Restriction-Filter-Rule" code="438" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.36-->
      <data type="IPFilterRule"/>
    </avp>

    <avp name="Service-Context-Id" code="461" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.42-->
      <data type="UTF8String"/>
    </avp>

    <avp name="Service-Identifier" code="439" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.28-->
      <data type="Unsigned32"/>
    </avp>

    <avp name="Service-Parameter-Info" code="440" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.43-->
      <data type="Grouped">
        <rule avp="Service-Parameter-Type" required="true" max="1"/>
        <rule avp="Service-Parameter-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Service-Parameter-Type" code="441" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.44-->
      <data type="Unsigned32"/>
    </avp>

    <avp name="Service-Parameter-Value" code="442" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.45-->
      <data type="OctetString"/>
    </avp>

    <avp name="Subscription-Id" code="443" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.46-->
      <data type="Grouped">
        <rule avp="Subscription-Id-Type" required="true" max="1"/>
        <rule avp="Subscription-Id-Data" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Subscription-Id-Data" code="444" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.48-->
      <data type="UTF8String"/>
    </avp>

    <avp name="Subscription-Id-Type" code="450" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.47-->
      <data type="Enumerated">
        <item code="0" name="END_USER_E164"/>
        <item code="1" name="END_USER_IMSI"/>
        <item code="2" name="END_USER_SIP_URI"/>
        <item code="3" name="END_USER_NAI"/>
      </data>
    </avp>

    <avp name="Tariff-Change-Usage" code="452" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.27-->
      <data type="Enumerated">
        <item code="0" name="UNIT_BEFORE_TARIFF_CHANGE"/>
        <item code="1" name="UNIT_AFTER_TARIFF_CHANGE"/>
        <item code="2" name="UNIT_INDETERMINATE"/>
      </data>
    </avp>

    <avp name="Tariff-Time-Change" code="451" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.20-->
      <data type="Time"/>
    </avp>

    <avp name="Unit-Value" code="445" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.8-->
      <data type="Grouped">
        <rule avp="Value-Digits" required="true" max="1"/>
        <rule avp="Exponent" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Used-Service-Unit" code="446" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.19-->
      <data type="Grouped">
        <rule avp="Tariff-Change-Usage" required="false" max="1"/>
        <rule avp="CC-Time" required="false" max="1"/>
        <rule avp="CC-Money" required="false" max="1"/>
        <rule avp="CC-Total-Octets" required="false" max="1"/>
        <rule avp="CC-Input-Octets" required="false" max="1"/>
        <rule avp="CC-Output-Octets" required="false" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="false" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="User-Equipment-Info" code="458" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.49-->
      <data type="Grouped">
        <rule avp="User-Equipment-Info-Type" required="true" max="1"/>
        <rule avp="User-Equipment-Info-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="User-Equipment-Info-Type" code="459" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.50-->
      <data type="Enumerated">
        <item code="0" name="IMEISV"/>
        <item code="1" name="MAC"/>
        <item code="2" name="EUI64"/>
        <item code="3" name="MODIFIED_EUI64"/>
      </data>
    </avp>

    <avp name="User-Equipment-Info-Value" code="460" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.51-->
      <data type="OctetString"/>
    </avp>

    <avp name="Value-Digits" code="447" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.10-->
      <data type="Integer64"/>
    </avp>

    <avp name="Validity-Time" code="448" must="M" may="P" must-not="V" may-encrypt="Y">
      <!-- http://tools.ietf.org/html/rfc4006#section-8.33-->
      <data type="Unsigned32"/>
    </avp>
  </application>
</diameter>`
