package library

var NetworkAccessServer = DictInfo{
	Name: "Network Access Server",
	XML: `<?xml version="1.0" encoding="UTF-8"?>
<diameter>

	<application id="1" type="auth" name="Network Access">
		<!-- Diameter Network Access Server Application -->
		<!-- http://tools.ietf.org/html/rfc7155 -->

		<command code="265" short="AA" name="AA">
			<request>
				<!-- https://tools.ietf.org/html/rfc7155#section-3.1 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Auth-Request-Type" required="true" max="1"/>
				<rule avp="Destination-Host" required="false" max="1"/>
				<rule avp="NAS-Identifier" required="false" max="1"/>
				<rule avp="NAS-IP-Address" required="true" max="1"/>
				<rule avp="NAS-IPv6-Address" required="false" max="1"/>
				<rule avp="NAS-Port" required="false" max="1"/>
				<rule avp="NAS-Port-Id" required="false" max="1"/>
				<rule avp="NAS-Port-Type" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Port-Limit" required="false" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="User-Password" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Authorization-Lifetime" required="false" max="1"/>
				<rule avp="Auth-Grace-Period" required="false" max="1"/>
				<rule avp="Auth-Session-State" required="false" max="1"/>
				<rule avp="Callback-Number" required="false" max="1"/>
				<rule avp="Called-Station-Id" required="false" max="1"/>
				<rule avp="Calling-Station-Id" required="false" max="1"/>
				<rule avp="Originating-Line-Info" required="false" max="1"/>
				<rule avp="Connect-Info" required="false" max="1"/>
				<rule avp="CHAP-Auth" required="false" max="1"/>
				<rule avp="CHAP-Challenge" required="false" max="1"/>
				<rule avp="Framed-Compression" required="false"/>
				<rule avp="Framed-Interface-Id" required="false" max="1"/>
				<rule avp="Framed-IP-Address" required="false" max="1"/>
				<rule avp="Framed-IPv6-Prefix" required="false"/>
				<rule avp="Framed-IP-Netmask" required="false" max="1"/>
				<rule avp="Framed-MTU" required="false" max="1"/>
				<rule avp="Framed-Protocol" required="false" max="1"/>
				<rule avp="ARAP-Password" required="false" max="1"/>
				<rule avp="ARAP-Security" required="false" max="1"/>
				<rule avp="ARAP-Security-Data" required="false"/>
				<rule avp="Login-IP-Host" required="false"/>
				<rule avp="Login-IPv6-Host" required="false"/>
				<rule avp="Login-LAT-Group" required="false" max="1"/>
				<rule avp="Login-LAT-Node" required="false" max="1"/>
				<rule avp="Login-LAT-Port" required="false" max="1"/>
				<rule avp="Login-LAT-Service" required="false" max="1"/>
				<rule avp="Tunneling" required="false"/>
				<rule avp="Proxy-Info" required="false"/>
				<rule avp="Route-Record" required="false"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.2 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Auth-Request-Type" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Configuration-Token" required="false"/>
				<rule avp="Acct-Interim-Interval" required="false" max="1"/>
				<rule avp="Error-Message" required="false" max="1"/>
				<rule avp="Error-Reporting-Host" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false"/>
				<rule avp="Idle-Timeout" required="false" max="1"/>
				<rule avp="Authorization-Lifetime" required="false" max="1"/>
				<rule avp="Auth-Grace-Period" required="false" max="1"/>
				<rule avp="Auth-Session-State" required="false" max="1"/>
				<rule avp="Re-Auth-Request-Type" required="false" max="1"/>
				<rule avp="Multi-Round-Time-Out" required="false" max="1"/>
				<rule avp="Session-Timeout" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Reply-Message" required="false"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Filter-Id" required="false"/>
				<rule avp="Password-Retry" required="false" max="1"/>
				<rule avp="Port-Limit" required="false" max="1"/>
				<rule avp="Prompt" required="false" max="1"/>
				<rule avp="ARAP-Challenge-Response" required="false" max="1"/>
				<rule avp="ARAP-Features" required="false" max="1"/>
				<rule avp="ARAP-Security" required="false" max="1"/>
				<rule avp="ARAP-Security-Data" required="false"/>
				<rule avp="ARAP-Zone-Access" required="false" max="1"/>
				<rule avp="Callback-Id" required="false" max="1"/>
				<rule avp="Callback-Number" required="false" max="1"/>
				<rule avp="Framed-Appletalk-Link" required="false" max="1"/>
				<rule avp="Framed-Appletalk-Network" required="false"/>
				<rule avp="Framed-Appletalk-Zone" required="false" max="1"/>
				<rule avp="Framed-Compression" required="false"/>
				<rule avp="Framed-Interface-Id" required="false" max="1"/>
				<rule avp="Framed-IP-Address" required="false" max="1"/>
				<rule avp="Framed-IPv6-Prefix" required="false"/>
				<rule avp="Framed-IPv6-Pool" required="false" max="1"/>
				<rule avp="Framed-IPv6-Route" required="false"/>
				<rule avp="Framed-IP-Netmask" required="false" max="1"/>
				<rule avp="Framed-Route" required="false"/>
				<rule avp="Framed-Pool" required="false" max="1"/>
				<rule avp="Framed-IPX-Network" required="false" max="1"/>
				<rule avp="Framed-MTU" required="false" max="1"/>
				<rule avp="Framed-Protocol" required="false" max="1"/>
				<rule avp="Framed-Routing" required="false" max="1"/>
				<rule avp="Login-IP-Host" required="false"/>
				<rule avp="Login-IPv6-Host" required="false"/>
				<rule avp="Login-LAT-Group" required="false" max="1"/>
				<rule avp="Login-LAT-Node" required="false" max="1"/>
				<rule avp="Login-LAT-Port" required="false" max="1"/>
				<rule avp="Login-LAT-Service" required="false" max="1"/>
				<rule avp="Login-Service" required="false" max="1"/>
				<rule avp="Login-TCP-Port" required="false" max="1"/>
				<rule avp="NAS-Filter-Rule" required="false"/>
				<rule avp="QoS-Filter-Rule" required="false"/>
				<rule avp="Tunneling" required="false"/>
				<rule avp="Redirect-Host" required="false"/>
				<rule avp="Redirect-Host-Usage" required="false" max="1"/>
				<rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
			</answer>
		</command>

		<command code="258" short="RA" name="Re-Auth">
			<request>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.3 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Destination-Host" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Re-Auth-Request-Type" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="NAS-Identifier" required="false" max="1"/>
				<rule avp="NAS-IP-Address" required="true" max="1"/>
				<rule avp="NAS-IPv6-Address" required="false" max="1"/>
				<rule avp="NAS-Port" required="false" max="1"/>
				<rule avp="NAS-Port-Id" required="false" max="1"/>
				<rule avp="NAS-Port-Type" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Framed-IP-Address" required="false" max="1"/>
				<rule avp="Framed-IPv6-Prefix" required="false" max="1"/>
				<rule avp="Framed-Interface-Id" required="false" max="1"/>
				<rule avp="Called-Station-Id" required="false" max="1"/>
				<rule avp="Calling-Station-Id" required="false" max="1"/>
				<rule avp="Originating-Line-Info" required="false" max="1"/>
				<rule avp="Acct-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Reply-Message" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
				<rule avp="Route-Record" required="false"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.4 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Error-Message" required="false" max="1"/>
				<rule avp="Error-Reporting-Host" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false"/>
				<rule avp="Redirect-Host" required="false"/>
				<rule avp="Redirect-Host-Usage" required="false" max="1"/>
				<rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Configuration-Token" required="false"/>
				<rule avp="Idle-Timeout" required="false" max="1"/>
				<rule avp="Authorization-Lifetime" required="false" max="1"/>
				<rule avp="Auth-Grace-Period" required="false" max="1"/>
				<rule avp="Re-Auth-Request-Type" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Reply-Message" required="false"/>
				<rule avp="Prompt" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
			</answer>
		</command>

		<command code="275" short="ST" name="Session-Termination">
			<request>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.5 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Termination-Cause" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Destination-Host" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
				<rule avp="Route-Record" required="false"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.6 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Error-Message" required="false" max="1"/>
				<rule avp="Error-Reporting-Host" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Redirect-Host" required="false"/>
				<rule avp="Redirect-Host-Usage" required="false" max="1"/>
				<rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
			</answer>
		</command>

		<command code="274" short="AS" name="Abort-Session">
			<request>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.7 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Destination-Host" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="NAS-Identifier" required="false" max="1"/>
				<rule avp="NAS-IP-Address" required="true" max="1"/>
				<rule avp="NAS-IPv6-Address" required="false" max="1"/>
				<rule avp="NAS-Port" required="false" max="1"/>
				<rule avp="NAS-Port-Id" required="false" max="1"/>
				<rule avp="NAS-Port-Type" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Framed-IP-Address" required="false" max="1"/>
				<rule avp="Framed-IPv6-Prefix" required="false" max="1"/>
				<rule avp="Framed-Interface-Id" required="false" max="1"/>
				<rule avp="Called-Station-Id" required="false" max="1"/>
				<rule avp="Calling-Station-Id" required="false" max="1"/>
				<rule avp="Originating-Line-Info" required="false" max="1"/>
				<rule avp="Acct-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Reply-Message" required="false"/>
				<rule avp="Proxy-Info" required="false"/>
				<rule avp="Route-Record" required="false"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.8 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="State" required="false" max="1"/>
				<rule avp="Error-Message" required="false" max="1"/>
				<rule avp="Error-Reporting-Host" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false"/>
				<rule avp="Redirect-Host" required="false"/>
				<rule avp="Redirect-Host-Usage" required="false" max="1"/>
				<rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false"/>
			</answer>
		</command>

		<command code="271" short="AC" name="Accounting">
			<request>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.9 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Accounting-Record-Type" required="true" max="1"/>
				<rule avp="Accounting-Record-Number" required="true" max="1"/>
				<rule avp="Acct-Application-Id" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Accounting-Sub-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Destination-Host" required="false" max="1"/>
				<rule avp="Event-Timestamp" required="false" max="1"/>
				<rule avp="Acct-Delay-Time" required="false" max="1"/>
				<rule avp="NAS-Identifier" required="false" max="1"/>
				<rule avp="NAS-IP-Address" required="true" max="1"/>
				<rule avp="NAS-IPv6-Address" required="false" max="1"/>
				<rule avp="NAS-Port" required="false" max="1"/>
				<rule avp="NAS-Port-Id" required="false" max="1"/>
				<rule avp="NAS-Port-Type" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Termination-Cause" required="false" max="1"/>
				<rule avp="Accounting-Input-Octets" required="false" max="1"/>
				<rule avp="Accounting-Input-Packets" required="false" max="1"/>
				<rule avp="Accounting-Output-Octets" required="false" max="1"/>
				<rule avp="Accounting-Output-Packets" required="false" max="1"/>
				<rule avp="Acct-Authentic" required="false" max="1"/>
				<rule avp="Accounting-Auth-Method" required="false" max="1"/>
				<rule avp="Acct-Link-Count" required="false" max="1"/>
				<rule avp="Acct-Session-Time" required="false" max="1"/>
				<rule avp="Acct-Tunnel-Connection" required="false" max="1"/>
				<rule avp="Acct-Tunnel-Packets-Lost" required="false" max="1"/>
				<rule avp="Callback-Id" required="false" max="1"/>
				<rule avp="Callback-Number" required="false" max="1"/>
				<rule avp="Called-Station-Id" required="false" max="1"/>
				<rule avp="Calling-Station-Id" required="false" max="1"/>
				<rule avp="Connection-Info" required="false"/>
				<rule avp="Originating-Line-Info" required="false" max="1"/>
				<rule avp="Authorization-Lifetime" required="false" max="1"/>
				<rule avp="Session-Timeout" required="false" max="1"/>
				<rule avp="Idle-Timeout" required="false" max="1"/>
				<rule avp="Port-Limit" required="false" max="1"/>
				<rule avp="Accounting-Realtime-Required" required="false" max="1"/>
				<rule avp="Acct-Interim-Interval" required="false" max="1"/>
				<rule avp="Filter-Id" required="false"/>
				<rule avp="NAS-Filter-Rule" required="false"/>
				<rule avp="QoS-Filter-Rule" required="false"/>
				<rule avp="Framed-Appletalk-Link" required="false" max="1"/>
				<rule avp="Framed-Appletalk-Network" required="false" max="1"/>
				<rule avp="Framed-Appletalk-Zone" required="false" max="1"/>
				<rule avp="Framed-Compression" required="false" max="1"/>
				<rule avp="Framed-Interface-Id" required="false" max="1"/>
				<rule avp="Framed-IP-Address" required="false" max="1"/>
				<rule avp="Framed-IP-Netmask" required="false" max="1"/>
				<rule avp="Framed-IPv6-Prefix" required="false"/>
				<rule avp="Framed-IPv6-Pool" required="false" max="1"/>
				<rule avp="Framed-IPv6-Route" required="false"/>
				<rule avp="Framed-IPX-Network" required="false" max="1"/>
				<rule avp="Framed-MTU" required="false" max="1"/>
				<rule avp="Framed-Pool" required="false" max="1"/>
				<rule avp="Framed-Protocol" required="false" max="1"/>
				<rule avp="Framed-Route" required="false"/>
				<rule avp="Framed-Routing" required="false" max="1"/>
				<rule avp="Login-IP-Host" required="false"/>
				<rule avp="Login-IPv6-Host" required="false"/>
				<rule avp="Login-LAT-Group" required="false" max="1"/>
				<rule avp="Login-LAT-Node" required="false" max="1"/>
				<rule avp="Login-LAT-Port" required="false" max="1"/>
				<rule avp="Login-LAT-Service" required="false" max="1"/>
				<rule avp="Login-Service" required="false" max="1"/>
				<rule avp="Login-TCP-Port" required="false" max="1"/>
				<rule avp="Tunneling" required="false"/>
				<rule avp="Proxy-Info" required="false"/>
				<rule avp="Route-Record" required="false"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc7155#section-3.10 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Accounting-Record-Type" required="true" max="1"/>
				<rule avp="Accounting-Record-Number" required="true" max="1"/>
				<rule avp="Acct-Application-Id" required="true" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Accounting-Sub-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Session-Id" required="false" max="1"/>
				<rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
				<rule avp="Event-Timestamp" required="false" max="1"/>
				<rule avp="Error-Message" required="false" max="1"/>
				<rule avp="Error-Reporting-Host" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false"/>
				<rule avp="Origin-AAA-Protocol" required="false" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="NAS-Identifier" required="false" max="1"/>
				<rule avp="NAS-IP-Address" required="true" max="1"/>
				<rule avp="NAS-IPv6-Address" required="false" max="1"/>
				<rule avp="NAS-Port" required="false" max="1"/>
				<rule avp="NAS-Port-Id" required="false" max="1"/>
				<rule avp="NAS-Port-Type" required="false" max="1"/>
				<rule avp="Service-Type" required="false" max="1"/>
				<rule avp="Termination-Cause" required="false" max="1"/>
				<rule avp="Accounting-Realtime-Required" required="false" max="1"/>
				<rule avp="Acct-Interim-Interval" required="false" max="1"/>
				<rule avp="Class" required="false"/>
				<rule avp="Proxy-Info" required="false"/>
			</answer>
		</command>



		<avp name="NAS-Port" code="5" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.2 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="NAS-Port-Id" code="87" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.3 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="NAS-Port-Type" code="61" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.4 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-13 -->
				<item code="0" name="Async"/>
				<item code="1" name="Sync"/>
				<item code="2" name="ISDN Sync"/>
				<item code="3" name="ISDN Async V.120"/>
				<item code="4" name="ISDN Async V.110"/>
				<item code="5" name="Virtual"/>
				<item code="6" name="PIAFS"/>
				<item code="7" name="HDLC Clear Channel"/>
				<item code="8" name="X.25"/>
				<item code="9" name="X.75"/>
				<item code="10" name="G.3 Fax"/>
				<item code="11" name="SDSL - Symmetric DSL"/>
				<item code="12" name="ADSL-CAP - Asymmetric DSL, Carrierless Amplitude Phase Modulation"/>
				<item code="13" name="ADSL-DMT - Asymmetric DSL, Discrete Multi-Tone"/>
				<item code="14" name="IDSL - ISDN Digital Subscriber Line"/>
				<item code="15" name="Ethernet"/>
				<item code="16" name="xDSL - Digital Subscriber Line of unknown type"/>
				<item code="17" name="Cable"/>
				<item code="18" name="Wireless - Other"/>
				<item code="19" name="Wireless - IEEE 802.11"/>
				<item code="20" name="Token-Ring"/>
				<item code="21" name="FDDI"/>
				<item code="22" name="Wireless - CDMA2000"/>
				<item code="23" name="Wireless - UMTS"/>
				<item code="24" name="Wireless - 1X-EV"/>
				<item code="25" name="IAPP"/>
				<item code="26" name="FTTP - Fiber to the Premises"/>
				<item code="27" name="Wireless - IEEE 802.16"/>
				<item code="28" name="Wireless - IEEE 802.20"/>
				<item code="29" name="Wireless - IEEE 802.22"/>
				<item code="30" name="PPPoA - PPP over ATM"/>
				<item code="31" name="PPPoEoA - PPP over Ethernet over ATM"/>
				<item code="32" name="PPPoEoE - PPP over Ethernet over Ethernet"/>
				<item code="33" name="PPPoEoVLAN - PPP over Ethernet over VLAN"/>
				<item code="34" name="PPPoEoQinQ - PPP over Ethernet over IEEE 802.1QinQ"/>
				<item code="35" name="xPON - Passive Optical Network"/>
				<item code="36" name="Wireless - XGP"/>
				<item code="37" name="WiMAX Pre-Release 8 IWK Function"/>
				<item code="38" name="WIMAX-WIFI-IWK: WiMAX WIFI Interworking"/>
				<item code="39" name="WIMAX-SFF: Signaling Forwarding Function for LTE/3GPP2"/>
				<item code="40" name="WIMAX-HA-LMA: WiMAX HA and or LMA function"/>
				<item code="41" name="WIMAX-DHCP: WIMAX DCHP service"/>
				<item code="42" name="WIMAX-LBS: WiMAX location based service"/>
				<item code="43" name="WIMAX-WVS: WiMAX voice service"/>
			</data>
		</avp>

		<avp name="Called-Station-Id" code="30" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.5 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Calling-Station-Id" code="31" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.6 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Connect-Info" code="77" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.7 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Originating-Line-Info" code="94" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.8 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Reply-Message" code="18" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.2.9 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="User-Password" code="2" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.1 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Password-Retry" code="75" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.2 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Prompt" code="76" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.3 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-17 -->
				<item code="0" name="No Echo"/>
				<item code="1" name="Echo"/>
			</data>
		</avp>

		<avp name="CHAP-Auth" code="402" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.4 -->
			<data type="Grouped">
				<rule avp="CHAP-Algorithm" required="true" max="1"/>
				<rule avp="CHAP-Ident" required="true" max="1"/>
				<rule avp="CHAP-Response" required="true" max="1"/>
			</data>
		</avp>


		<avp name="CHAP-Algorithm" code="403" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.5 -->
			<data type="Enumerated">
				<item code="5" name="CHAP with MD5"/>
			</data>
		</avp>

		<avp name="CHAP-Ident" code="404" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.6 -->
			<data type="OctetString"/>
		</avp>

		<avp name="CHAP-Response" code="405" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.7 -->
			<data type="OctetString"/>
		</avp>

		<avp name="CHAP-Challenge" code="60" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.8 -->
			<data type="OctetString"/>
		</avp>

		<avp name="ARAP-Password" code="70" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.9 -->
			<data type="OctetString"/>
		</avp>

		<avp name="ARAP-Challenge-Response" code="84" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.10 -->
			<data type="OctetString"/>
		</avp>

		<avp name="ARAP-Security" code="73" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.11 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="ARAP-Security-Data" code="74" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.3.12 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Service-Type" code="6" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.1 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-4 -->
				<item code="1" name="Login"/>
				<item code="2" name="Framed"/>
				<item code="3" name="Callback Login"/>
				<item code="4" name="Callback Framed"/>
				<item code="5" name="Outbound"/>
				<item code="6" name="Administrative"/>
				<item code="7" name="NAS Prompt"/>
				<item code="8" name="Authenticate Only"/>
				<item code="9" name="Callback NAS Prompt"/>
				<item code="10" name="Call Check"/>
				<item code="11" name="Callback Administrative"/>
				<item code="12" name="Voice"/>
				<item code="13" name="Fax"/>
				<item code="14" name="Modem Relay"/>
				<item code="15" name="IAPP-Register"/>
				<item code="16" name="IAPP-AP-Check"/>
				<item code="17" name="Authorize Only"/>
				<item code="18" name="Framed-Management"/>
				<item code="19" name="Additional-Authorization"/>
			</data>
		</avp>

		<avp name="Callback-Number" code="19" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.2 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Callback-Id" code="20" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.3 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Idle-Timeout" code="28" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.4 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Port-Limit" code="62" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.5 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="NAS-Filter-Rule" code="400" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.6 -->
			<data type="IPFilterRule"/>
		</avp>

		<avp name="Filter-Id" code="11" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.7 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Configuration-Token" code="78" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.8 -->
			<data type="OctetString"/>
		</avp>

		<!--avp name="QoS-Filter-Rule" code="407" must="-" may="" must-not="-" may-encrypt="Y"-->
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.9 -->
			<!--data type="QoSFilterRule"/-->
		<!--/avp-->


		<avp name="Framed-Protocol" code="7" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.1 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-5 -->
				<item code="1" name="PPP"/>
				<item code="2" name="SLIP"/>
				<item code="3" name="AppleTalk Remote Access Protocol (ARAP)"/>
				<item code="4" name="Gandalf proprietary SingleLink/MultiLink protocol	"/>
				<item code="5" name="Xylogics proprietary IPX/SLIP"/>
				<item code="6" name="X.75 Synchronous"/>
				<item code="7" name="GPRS PDP Context"/>
			</data>
		</avp>

		<avp name="Framed-Routing" code="10" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.2 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-6 -->
				<item code="0" name="None"/>
				<item code="1" name="Send routing packets"/>
				<item code="2" name="Listen for routing packets"/>
				<item code="3" name="Send and Listen"/>
			</data>
		</avp>

		<avp name="Framed-MTU" code="12" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.3 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Framed-Compression" code="13" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.4 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-7 -->
				<item code="0" name="None"/>
				<item code="1" name="VJ TCP/IP header compression	"/>
				<item code="2" name="IPX header compression"/>
				<item code="3" name="Stac-LZS compression"/>
			</data>
		</avp>

		<avp name="Framed-IP-Address" code="8" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.1 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Framed-IP-Netmask" code="9" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.2 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Framed-Route" code="22" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.3 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Framed-Pool" code="88" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.4 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Framed-Interface-Id" code="96" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.5 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="Framed-IPv6-Prefix" code="97" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.6 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Framed-IPv6-Route" code="99" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.7 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Framed-IPv6-Pool" code="100" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.5.8 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Framed-IPX-Network" code="23" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.6.1-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Framed-Appletalk-Link" code="37" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.7.1-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Framed-Appletalk-Network" code="38" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.7.2-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Framed-Appletalk-Zone" code="39" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.7.3-->
			<data type="OctetString"/>
		</avp>

		<avp name="ARAP-Features" code="71" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.8.1-->
			<data type="OctetString"/>
		</avp>

		<avp name="ARAP-Zone-Access" code="72" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.10.8.2 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-16 -->
				<item code="1" name="Only allow access to default zone"/>
				<item code="2" name="Use zone filter inclusively"/>
				<item code="3" name="Not used"/>
				<item code="4" name="Use zone filter exclusively"/>
			</data>
		</avp>

		<avp name="Login-IP-Host" code="14" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.1-->
			<data type="OctetString"/>
		</avp>

		<avp name="Login-IPv6-Host" code="98" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.2-->
			<data type="OctetString"/>
		</avp>

		<avp name="Login-Service" code="15" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.3 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-8 -->
				<item code="0" name="Telnet"/>
				<item code="1" name="Rlogin"/>
				<item code="2" name="TCP Clear"/>
				<item code="3" name="PortMaster (proprietary)"/>
				<item code="4" name="LAT"/>
				<item code="5" name="X25-PAD"/>
				<item code="6" name="X25-T3POS"/>
				<item code="7" name="Unassigned"/>
				<item code="8" name="TCP Clear Quiet (suppresses any NAS-generated connect string)"/>
			</data>
		</avp>

		<avp name="Login-TCP-Port" code="16" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.4.1-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Login-LAT-Service" code="34" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.5.1-->
			<data type="OctetString"/>
		</avp>

		<avp name="Login-LAT-Node" code="35" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.5.2-->
			<data type="OctetString"/>
		</avp>

		<avp name="Login-LAT-Group" code="36" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.5.3-->
			<data type="OctetString"/>
		</avp>

		<avp name="Login-LAT-Port" code="63" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.11.5.4-->
			<data type="OctetString"/>
		</avp>

		<avp name="Tunneling" code="401" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.1-->
			<data type="Grouped">
				<rule avp="Tunnel-Type" required="true" max="1"/>
				<rule avp="Tunnel-Medium-Type" required="true" max="1"/>
				<rule avp="Tunnel-Client-Endpoint" required="true" max="1"/>
				<rule avp="Tunnel-Server-Endpoint" required="true" max="1"/>
				<rule avp="Tunnel-Preference" required="false" max="1"/>
				<rule avp="Tunnel-Client-Auth-Id" required="false" max="1"/>
				<rule avp="Tunnel-Server-Auth-Id" required="false" max="1"/>
				<rule avp="Tunnel-Assignment-Id" required="false" max="1"/>
				<rule avp="Tunnel-Password" required="false" max="1"/>
				<rule avp="Tunnel-Private-Group-Id" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Tunnel-Type" code="64" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.2 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-14 -->
				<item code="1" name="Point-to-Point Tunneling Protocol (PPTP)"/>
				<item code="2" name="Layer Two Forwarding (L2F)"/>
				<item code="3" name="Layer Two Tunneling Protocol (L2TP)"/>
				<item code="4" name="Ascend Tunnel Management Protocol (ATMP)"/>
				<item code="5" name="Virtual Tunneling Protocol (VTP)"/>
				<item code="6" name="IP Authentication Header in the Tunnel-mode (AH)"/>
				<item code="7" name="IP-in-IP Encapsulation (IP-IP)"/>
				<item code="8" name="Minimal IP-in-IP Encapsulation (MIN-IP-IP)"/>
				<item code="9" name="IP Encapsulating Security Payload in the Tunnel-mode (ESP)"/>
				<item code="10" name="Generic Route Encapsulation (GRE)"/>
				<item code="11" name="Bay Dial Virtual Services (DVS)"/>
				<item code="12" name="IP-in-IP Tunneling"/>
				<item code="13" name="Virtual LANs (VLAN)"/>
			</data>
		</avp>

		<avp name="Tunnel-Medium-Type" code="65" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.3 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-15 -->
				<item code="1" name="IPv4 (IP version 4)"/>
				<item code="2" name="IPv6 (IP version 6)"/>
				<item code="3" name="NSAP"/>
				<item code="4" name="HDLC (8-bit multidrop)"/>
				<item code="5" name="BBN 1822"/>
				<item code="6" name="802 (includes all 802 media plus Ethernet 'canonical format')"/>
				<item code="7" name="E.163 (POTS)"/>
				<item code="8" name="E.164 (SMDS, Frame Relay, ATM)"/>
				<item code="9" name="F.69 (Telex)"/>
				<item code="10" name="X.121 (X.25, Frame Relay)"/>
				<item code="11" name="IPX"/>
				<item code="12" name="Appletalk"/>
				<item code="13" name="Decnet IV"/>
				<item code="14" name="Banyan Vines"/>
				<item code="15" name="E.164 with NSAP format subaddress"/>
			</data>
		</avp>

		<avp name="Tunnel-Client-Endpoint" code="66" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.4 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Tunnel-Server-Endpoint" code="67" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.5 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Tunnel-Password" code="69" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.6 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Tunnel-Private-Group-Id" code="81" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.7 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Tunnel-Assignment-Id" code="82" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.8 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Tunnel-Preference" code="83" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.9 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Tunnel-Client-Auth-Id" code="90" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.10 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Tunnel-Server-Auth-Id" code="91" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.5.11 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Accounting-Input-Octets" code="363" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.1 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="Accounting-Output-Octets" code="364" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.2 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="Accounting-Input-Packets" code="365" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.3 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="Accounting-Output-Packets" code="366" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.4 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="Acct-Session-Time" code="46" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.5 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Acct-Authentic" code="45" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.6 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/radius-types/radius-types.xhtml#radius-types-11 -->
				<item code="1" name="RADIUS"/>
				<item code="2" name="Local"/>
				<item code="3" name="Remote"/>
				<item code="4" name="Diameter"/>
			</data>
		</avp>

		<avp name="Accounting-Auth-Method" code="406" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.7 -->
			<data type="Enumerated">
				<!-- http://www.iana.org/assignments/aaa-parameters/aaa-parameters.xhtml#aaa-parameters-26 -->
				<item code="1" name="PAP"/>
				<item code="2" name="CHAP"/>
				<item code="3" name="MS-CHAP-1"/>
				<item code="4" name="MS-CHAP-2"/>
				<item code="5" name="EAP"/>
				<item code="7" name="None"/>
			</data>
		</avp>

		<avp name="Acct-Delay-Time" code="41" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.8 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Acct-Link-Count" code="51" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.9 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Acct-Tunnel-Connection" code="68" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.10 -->
			<data type="OctetString"/>
		</avp>

		<avp name="Acct-Tunnel-Packets-Lost" code="86" must="M" may="-" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.6.11 -->
			<data type="Unsigned32"/>
		</avp>

	</application>
</diameter>`,
}
