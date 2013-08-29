// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Static dictionary Parser and Base Protocol XML.

package base

import "github.com/fiorix/go-diameter/dict"

// Dict is a static Parser with a pre-loaded Base Protocol.
var Dict *dict.Parser

func init() {
	Dict, _ = dict.New()
	Dict.Load(DictXML)
}

// DictXML is an embedded version of the Diameter Base Protocol.
//
// Copy of ../dict/diam_base.xml
var DictXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="0"> <!-- Diameter Common Messages -->

    <vendor id="10415" name="3GPP"/>

    <command code="257" short="CE" name="Capabilities-Exchange"/>
    <command code="258" short="RA" name="Re-Auth"/>
    <command code="271" short="AC" name="Accounting"/>
    <command code="274" short="AS" name="Abort-Session"/>
    <command code="275" short="ST" name="Session-Termination"/>
    <command code="280" short="DW" name="Device-Watchdog"/>
    <command code="282" short="DP" name="Disconnect-Peer"/>

    <avp name="Acct-Interim-Interval" code="85" must="M" may="P" must-not="V" encr="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Accounting-Realtime-Required" code="483" must="M" may="P" must-not="V" encr="Y">
      <data type="Enumerated">
        <item code="1" name="DELIVER_AND_GRANT"/>
        <item code="2" name="GRANT_AND_STORE"/>
        <item code="3" name="GRANT_AND_LOSE"/>
      </data>
    </avp>

    <avp name="Acct-Multi-Session-Id" code="50" must="M" may="P" must-not="V" encr="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Accounting-Record-Number" code="485" must="M" may="P" must-not="V" encr="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Accounting-Record-Type" code="480" must="M" may="P" must-not="V" encr="Y">
      <data type="Enumerated">
        <item code="1" name="EVENT_RECORD"/>
        <item code="2" name="START_RECORD"/>
        <item code="3" name="INTERIM_RECORD"/>
        <item code="4" name="STOP_RECORD"/>
      </data>
    </avp>

    <avp name="Accounting-Session-Id" code="44" must="M" may="P" must-not="V" encr="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Accounting-Sub-Session-Id" code="287" must="M" may="P" must-not="V" encr="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="Acct-Application-Id" code="259" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Application-Id" code="258" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Request-Type" code="274" must="M" may="P" must-not="V" encr="-">
      <data type="Enumerated">
        <item code="1" name="AUTHENTICATE_ONLY"/>
        <item code="2" name="AUTHORIZE_ONLY"/>
        <item code="3" name="AUTHORIZE_AUTHENTICATE"/>
      </data>
    </avp>

    <avp name="Authorization-Lifetime" code="291" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Grace-Period" code="276" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Auth-Session-State" code="277" must="M" may="P" must-not="V" encr="-">
      <data type="Enumerated">
        <item code="0" name="STATE_MAINTAINED"/>
        <item code="1" name="NO_STATE_MAINTAINED"/>
      </data>
    </avp>

    <avp name="Re-Auth-Request-Type" code="285" must="M" may="P" must-not="V" encr="-">
      <data type="Enumerated">
        <item code="0" name="AUTHORIZE_ONLY"/>
        <item code="1" name="AUTHORIZE_AUTHENTICATE"/>
      </data>
    </avp>

    <avp name="Class" code="25" must="M" may="P" must-not="V" encr="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Destination-Host" code="293" must="M" may="P" must-not="V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Destination-Realm" code="283" must="M" may="P" must-not="V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Disconnect-Cause" code="273" must="M" may="P" must-not="V" encr="-">
      <data type="Enumerated">
        <item code="0" name="REBOOTING"/>
        <item code="1" name="BUSY"/>
        <item code="2" name="DO_NOT_WANT_TO_TALK_TO_YOU"/>
      </data>
    </avp>

    <avp name="Error-Message" code="281" must="-" may="P" must-not="V,M" encr="-">
      <data type="UTF8String"/>
    </avp>

    <avp name="Error-Reporting-Host" code="294" must="-" may="P" must-not="V,M" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Event-Timestamp" code="55" must="M" may="P" must-not="V" encr="-">
      <data type="Time"/>
    </avp>

    <avp name="Experimental-Result" code="297" must="M" may="P" must-not="V" encr="-">
      <data type="Grouped">
        <rule avp="Vendor-Id" required="true" max="1"/>
        <rule avp="Experimental-Result-Code" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Experimental-Result-Code" code="298" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Failed-AVP" code="279" must="M" may="P" must-not="V" encr="-">
      <data type="Grouped"/>
    </avp>

    <avp name="Firmware-Revision" code="267" must="-" may="-" must-not="P,V,M" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Host-IP-Address" code="257" must="M" may="P" must-not="V" encr="-">
      <data type="Address"/>
    </avp>

    <avp name="Inband-Security-Id" code="299" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Multi-Round-Time-Out" code="272" must="M" may="P" must-not="V" encr="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Origin-Host" code="264" must="M" may="P" must-not="V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Origin-Realm" code="296" must="M" may="P" must-not="V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Origin-State-Id" code="278" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Product-Name" code="269" must="-" may="-" must-not="P,V,M" encr="-">
      <data type="UTF8String"/>
    </avp>

    <avp name="Proxy-Host" code="280" must="M" may="-" must-not="P,V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Proxy-Info" code="284" must="M" may="-" must-not="P,V" encr="-">
      <data type="Grouped">
        <rule avp="Proxy-Host" required="true" max="1"/>
        <rule avp="Proxy-State" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Proxy-State" code="33" must="M" may="-" must-not="P,V" encr="-">
      <data type="OctetString"/>
    </avp>

    <avp name="Redirect-Host" code="292" must="M" may="P" must-not="V" encr="-">
      <data type="DiameterURI"/>
    </avp>

    <avp name="Redirect-Host-Usage" code="261" must="M" may="P" must-not="V" encr="-">
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

    <avp name="Redirect-Max-Cache-Time" code="262" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Result-Code" code="268" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Route-Record" code="282" must="M" may="-" must-not="P,V" encr="-">
      <data type="DiameterIdentity"/>
    </avp>

    <avp name="Session-Id" code="263" must="M" may="P" must-not="V" encr="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Session-Timeout" code="27" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Session-Binding" code="270" must="M" may="P" must-not="V" encr="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Session-Server-Failover" code="271" must="M" may="P" must-not="V" encr="Y">
      <data type="Enumerated">
        <item code="0" name="REFUSE_SERVICE"/>
        <item code="1" name="TRY_AGAIN"/>
        <item code="2" name="ALLOW_SERVICE"/>
        <item code="3" name="TRY_AGAIN_ALLOW_SERVICE"/>
      </data>
    </avp>

    <avp name="Supported-Vendor-Id" code="265" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Termination-Cause" code="295" must="M" may="P" must-not="V" encr="-">
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

    <avp name="User-Name" code="1" must="M" may="P" must-not="V" encr="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Vendor-Id" code="266" must="M" may="P" must-not="V" encr="-">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Vendor-Specific-Application-Id" code="260" must="M" may="P" must-not="V" encr="-">
      <data type="Grouped">
        <rule avp="Vendor-Id" required="false" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Acct-Application-Id" required="true" max="1"/>
      </data>
    </avp>

  </application>
</diameter>`)
