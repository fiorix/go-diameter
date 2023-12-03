package diam

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func TestPrettyDump(t *testing.T) {

	msg := NewMessage(CreditControl, RequestFlag, CHARGING_CONTROL_APP_ID, 0xa8cc407d, 0xa8c1b2b4, dict.Default)
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))
	msg.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("10.1.0.1")))
	msg.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13))
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("sess;123456789"))
	msg.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(1397760650))
	msg.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1))
	msg.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(1000))
	msg.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(7786)),
			NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(7786)),
			NewAVP(avp.TGPPRATType, avp.Mbit, 10415, datatype.OctetString("1234")),
		},
	})
	msg.NewAVP(avp.ServiceInformation, avp.Mbit, 10415, &GroupedAVP{
		AVP: []*AVP{
			NewAVP(avp.PSInformation, avp.Mbit, 10415, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.CalledStationID, avp.Mbit, 0, datatype.UTF8String("10999")),
					NewAVP(avp.StartTime, 0, 10415, datatype.Time(time.Unix(1377093974, 0))),
				},
			}),
		}})

	// Existing String() print
	t.Logf("Message:\n%s", msg)

	// New PrettyDump() print
	t.Logf("Message:\n%s", msg.PrettyDump())
}

func TestPrettyDumpAVP(t *testing.T) {
	tests := []struct {
		name     string
		avp      *AVP
		expected string
	}{
		{
			name:     "Unsigned32",
			avp:      NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(13)),
			expected: "  Vendor-Id                                       0   266  ✗ ✓ ✗  Unsigned32          13",
		},
		{
			name:     "UTF8String",
			avp:      NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("abc-1234567")),
			expected: "  Session-Id                                      0   263  ✗ ✓ ✗  UTF8String          abc-1234567",
		},
		{
			name:     "Address",
			avp:      NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("10.1.0.1"))),
			expected: "  Host-IP-Address                                 0   257  ✗ ✓ ✗  Address             10.1.0.1",
		},
		{
			name:     "AddressIPv6",
			avp:      NewAVP(avp.GGSNAddress, avp.Mbit, 10415, datatype.Address(net.ParseIP("2001:0db8::ff00:0042:8329"))),
			expected: "  GGSN-Address                                10415   847  ✓ ✓ ✗  Address             2001:db8::ff00:42:8329",
		},
		{
			name:     "Enumerated",
			avp:      NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1)),
			expected: "  CC-Request-Type                                 0   416  ✗ ✓ ✗  Enumerated          1",
		},
		{
			name: "GroupedAVP",
			avp: NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(7786)),
					NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(7786)),
					NewAVP(avp.TGPPRATType, avp.Mbit, 10415, datatype.OctetString("1234")),
				},
			}),
			expected: strings.Join([]string{
				"  Multiple-Services-Credit-Control                0   456  ✗ ✓ ✗  Grouped",
				"    Service-Identifier                            0   439  ✗ ✓ ✗  Unsigned32          7786",
				"    Rating-Group                                  0   432  ✗ ✓ ✗  Unsigned32          7786",
				"    TGPP-RAT-Type                             10415    21  ✓ ✓ ✗  OctetString         1234",
			}, "\n"),
		},
		{
			name: "NestedGroupedAVP",
			avp: NewAVP(avp.ServiceInformation, avp.Mbit, 10415, &GroupedAVP{
				AVP: []*AVP{
					NewAVP(avp.PSInformation, avp.Mbit, 10415, &GroupedAVP{
						AVP: []*AVP{
							NewAVP(avp.CalledStationID, avp.Mbit, 0, datatype.UTF8String("10999")),
							NewAVP(avp.StartTime, 0, 10415, datatype.Time(time.Date(2023, 8, 21, 22, 06, 14, 0, time.UTC))),
						},
					}),
				}}),
			expected: strings.Join([]string{
				"  Service-Information                         10415   873  ✓ ✓ ✗  Grouped",
				"    PS-Information                            10415   874  ✓ ✓ ✗  Grouped",
				"      Called-Station-Id                           0    30  ✗ ✓ ✗  UTF8String          10999",
				"      Start-Time                              10415  2041  ✓ ✗ ✗  Time                2023-08-21 22:06:14 +0000 UTC",
			}, "\n"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := NewMessage(CreditControl, RequestFlag, CHARGING_CONTROL_APP_ID, 0xa8cc407d, 0xa8c1b2b4, dict.Default)

			var b bytes.Buffer
			prettyDumpAVP(&b, msg, tc.avp, 0)

			lines := strings.Split(b.String(), "\n")
			for i, line := range lines {
				lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
			}
			actual := strings.Join(lines, "\n")
			actual = strings.TrimSuffix(actual, "\n")

			if actual != tc.expected {
				t.Errorf("\nActual:\n%v\nExpected:\n%v\n", actual, tc.expected)
			}
		})
	}
}
