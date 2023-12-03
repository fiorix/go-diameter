package diam

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func (m *Message) PrettyDump() string {
	var b bytes.Buffer
	prettyDumpMessage(&b, m, 0)
	return b.String()
}

func prettyDumpMessage(w io.Writer, m *Message, depth int) {
	requestFlag, errorFlag, proxyFlag, retransmittedFlag := flagsToString(m.Header)

	// Print Header
	fmt.Fprintf(w, "%s(%d) %s(%d) %s%s%s%s %d, %d\n",
		cmdToString(m.Dictionary(), m.Header),
		m.Header.CommandCode,
		appIdToString(int(m.Header.ApplicationID)),
		m.Header.ApplicationID,
		requestFlag,
		errorFlag,
		proxyFlag,
		retransmittedFlag,
		m.Header.HopByHopID,
		m.Header.EndToEndID)

	// Print Titles
	fmt.Fprintf(w, "  %-40s %8s %5s  %s %s %s  %-18s  %s\n",
		"AVP", "Vendor", "Code", "V", "M", "P", "Type", "Value")

	// Print AVPs
	for _, a := range m.AVP {
		prettyDumpAVP(w, m, a, depth)
	}
}

func prettyDumpGroupedAVP(w io.Writer, m *Message, a *AVP, depth int) {
	for _, ga := range a.Data.(*GroupedAVP).AVP {
		prettyDumpAVP(w, m, ga, depth)
	}
}

func prettyDumpAVP(w io.Writer, m *Message, a *AVP, depth int) {
	indent := strings.Repeat("  ", max(0, depth))

	avpName, avpType, avpData, isGrouped := avpToString(m, a)

	fmt.Fprintf(w, "  %-40s %8d %5d  %s %s %s  %-18s  %s\n",
		indent+avpName,
		a.VendorID,
		a.Code,
		boolToSymbol(a.Flags&avp.Vbit == avp.Vbit),
		boolToSymbol(a.Flags&avp.Mbit == avp.Mbit),
		boolToSymbol(a.Flags&avp.Pbit == avp.Pbit),
		avpType,
		avpData)

	if isGrouped {
		prettyDumpGroupedAVP(w, m, a, depth+1)
	}
}

func cmdToString(dictionary *dict.Parser, header *Header) string {
	if dictCMD, err := dictionary.FindCommand(
		header.ApplicationID,
		header.CommandCode,
	); err != nil {
		return "Unknown"
	} else {
		return dictCMD.Name
	}
}

func appIdToString(appId int) string {
	switch appId {
	case BASE_APP_ID:
		return "Common"
	case NETWORK_ACCESS_APP_ID:
		return "Network-Access"
	case BASE_ACCOUNTING_APP_ID:
		return "Accounting"
	case CHARGING_CONTROL_APP_ID:
		return "Charging-Control"
	//case TGPP_APP_ID:
	//	return "TGPP_APP_ID"
	case GX_CHARGING_CONTROL_APP_ID:
		return "Gx"
	case TGPP_S6A_APP_ID:
		return "S6A"
	case TGPP_SWX_APP_ID:
		return "SWX"
	case DIAMETER_SY_APP_ID:
		return "Sy"
	default:
		return "Unknown"
	}
}

func flagsToString(header *Header) (string, string, string, string) {
	var requestFlag string
	if header.CommandFlags&RequestFlag == RequestFlag {
		requestFlag = "request"
	} else {
		requestFlag = "answer"
	}

	var errorFlag string
	if header.CommandFlags&ErrorFlag == ErrorFlag {
		errorFlag = "error"
	} else {
		errorFlag = ""
	}

	var proxyFlag string
	if header.CommandFlags&ProxiableFlag == ProxiableFlag {
		proxyFlag = "proxiable"
	} else {
		proxyFlag = ""
	}

	var retransmittedFlag string
	if header.CommandFlags&RetransmittedFlag == RetransmittedFlag {
		retransmittedFlag = "retransmitted"
	} else {
		retransmittedFlag = ""
	}

	return requestFlag, errorFlag, proxyFlag, retransmittedFlag
}

func avpToString(m *Message, a *AVP) (string, string, string, bool) {

	var avpName string
	var avpType string
	var avpData string
	var isGrouped bool

	if dictAVP, err := m.Dictionary().FindAVPWithVendor(
		m.Header.ApplicationID,
		a.Code,
		a.VendorID,
	); err != nil {
		avpName = "Unknown"
		avpType = "Unknown"
		avpData = a.Data.String()
		isGrouped = false
	} else if a.Data.Type() == GroupedAVPType {
		avpName = dictAVP.Name
		avpType = "Grouped"
		avpData = ""
		isGrouped = true
	} else {
		for k, v := range datatype.Available {
			if v == a.Data.Type() {
				avpType = k
				break
			}
		}
		avpName = dictAVP.Name
		avpData = dataValueToString(a.Data)
		isGrouped = false
	}

	return avpName, avpType, avpData, isGrouped
}

func dataValueToString(data datatype.Type) string {

	switch data.Type() {
	case datatype.Integer32Type,
		datatype.Integer64Type,
		datatype.Unsigned32Type,
		datatype.Unsigned64Type,
		datatype.EnumeratedType:
		return fmt.Sprintf("%d", data)

	case datatype.Float32Type,
		datatype.Float64Type:
		return fmt.Sprintf("%0.4f", data)

	case datatype.OctetStringType:
		return string(data.(datatype.OctetString))

	case datatype.UTF8StringType:
		return string(data.(datatype.UTF8String))

	case datatype.DiameterIdentityType:
		return string(data.(datatype.DiameterIdentity))

	case datatype.DiameterURIType:
		return string(data.(datatype.DiameterURI))

	case datatype.IPFilterRuleType:
		return string(data.(datatype.IPFilterRule))

	case datatype.QoSFilterRuleType:
		return string(data.(datatype.QoSFilterRule))

	case datatype.TimeType:
		return fmt.Sprintf("%s", time.Time(data.(datatype.Time)))

	case datatype.AddressType:
		addr := string(data.(datatype.Address))
		if ip4 := net.IP(addr).To4(); ip4 != nil {
			return fmt.Sprintf("%s", net.IP(addr))
		}
		if ip6 := net.IP(addr).To16(); ip6 != nil {
			return fmt.Sprintf("%s", net.IP(addr))
		}
		return fmt.Sprintf("%#v, %#v", addr[2:], addr[:2])

	case datatype.IPv4Type:
		addr := string(data.(datatype.IPv4))
		return fmt.Sprintf("%s", net.IP(addr))

	case datatype.IPv6Type:
		addr := string(data.(datatype.IPv6))
		return fmt.Sprintf("%s", net.IP(addr))
	}

	return data.String()
}

func boolToSymbol(flag bool) string {
	if flag {
		return "\u2713" // ✓
	}
	return "\u2717" // ✗
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
