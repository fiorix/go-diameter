package diam

import (
	"net"
	"testing"
	"time"

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
					NewAVP(avp.StartTime, avp.Mbit, 10415, datatype.Time(time.Unix(1377093974, 0))),
				},
			}),
		}})

	// Existing String() print
	t.Logf("Message:\n%s", msg)

	// New PrettyDump() print
	t.Logf("Message:\n%s", msg.PrettyDump())

	// TODO Maybe make PrettyDump() testable and assert the output
}
