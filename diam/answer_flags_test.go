package diam

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// testRetransmittedAIR is an S6a Authentication-Information-Request with
// the R, P and T flags set (0xd0), as an MME sends it again after a
// transport failover (RFC 6733 section 5.5.4).
var testRetransmittedAIR, _ = hex.DecodeString("" +
	"0100007cd000013e010000230000002a0000002b" + // header: AIR, flags R|P|T
	"000001074000000f6d6d653b313b3100" + // Session-Id "mme;1;1"
	"000001154000000c00000001" + // Auth-Session-State NO_STATE_MAINTAINED
	"000001084000000b6d6d6500" + // Origin-Host "mme"
	"000001284000000b65706300" + // Origin-Realm "epc"
	"0000011b4000000b65706300" + // Destination-Realm "epc"
	"0000000140000017303031303130313233343536373839" + "00" + // User-Name
	"0000057fc000000f000028af00f11000") // Visited-PLMN-Id 001/01

// TestMessageAnswerFlags checks that Answer keeps only the P flag of the
// request. RFC 6733 section 3: the T flag "MUST NOT be set in answer
// messages", the E flag "MUST NOT be set in request messages" and so is
// only set in an answer for a protocol error, and the reserved bits "MUST
// be set to zero". Section 6.2: the P bit is set to the value it has in
// the request.
func TestMessageAnswerFlags(t *testing.T) {
	req, err := ReadMessage(bytes.NewReader(testRetransmittedAIR), dict.Default)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		request, answer uint8
	}{
		{RequestFlag | ProxiableFlag | RetransmittedFlag, ProxiableFlag},
		{RequestFlag | RetransmittedFlag, 0},
		{RequestFlag | ProxiableFlag, ProxiableFlag},
		{RequestFlag | ErrorFlag, 0},
		{RequestFlag | ProxiableFlag | 0x0f, ProxiableFlag}, // reserved bits
	} {
		req.Header.CommandFlags = tc.request
		ans := req.Answer(2001)
		if got := ans.Header.CommandFlags; got != tc.answer {
			t.Errorf("request flags %#02x: answer flags = %#02x, want %#02x",
				tc.request, got, tc.answer)
		}
	}
}
