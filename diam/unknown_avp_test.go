package diam

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// helper: start a server with ALL handler, return listener + handled channel
func startTestServer(t *testing.T) (net.Listener, chan *Message) {
	t.Helper()
	return startTestServerWithIdentity(t, "", "")
}

// helper: like startTestServer but with an Origin-Host/Realm configured,
// so the server can emit RFC-compliant 3001 answers for unknown commands.
func startTestServerWithIdentity(t *testing.T, host, realm datatype.DiameterIdentity) (net.Listener, chan *Message) {
	t.Helper()
	handled := make(chan *Message, 10)
	mux := NewServeMux()
	mux.HandleFunc("ALL", func(c Conn, m *Message) {
		handled <- m
	})
	go func() {
		for err := range mux.ErrorReports() {
			t.Logf("SERVER ERROR: %v", err)
		}
	}()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &Server{Handler: mux, OriginHost: host, OriginRealm: realm}
	go srv.Serve(ln)
	return ln, handled
}

// helper: build a minimal valid CER
func buildCER(extra ...*AVP) *Message {
	m := NewRequest(257, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test.host"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test.realm"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	for _, a := range extra {
		m.AddAVP(a)
	}
	return m
}

func sendMsg(t *testing.T, conn net.Conn, m *Message) {
	t.Helper()
	b, err := m.Serialize()
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if _, err := conn.Write(b); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func expectHandled(t *testing.T, handled chan *Message, label string) *Message {
	t.Helper()
	select {
	case m := <-handled:
		t.Logf("%s: handled OK (AVPs: %d)", label, len(m.AVP))
		return m
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: NOT handled — connection likely dead", label)
		return nil
	}
}

// Test 1: Unknown AVP (M=0) at top level — should NOT kill connection
func TestUnknownAVP_Mbit0_TopLevel(t *testing.T) {
	ln, handled := startTestServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// CER with unknown AVP code 99999, M=0
	sendMsg(t, conn, buildCER(NewAVP(99999, 0, 0, datatype.OctetString("hello"))))
	expectHandled(t, handled, "unknown AVP M=0")

	// Prove connection still alive
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after M=0")
}

// Test 2: Unknown AVP (M=1) at top level — should NOT kill connection
func TestUnknownAVP_Mbit1_TopLevel(t *testing.T) {
	ln, handled := startTestServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// CER with unknown AVP code 99998, M=1
	sendMsg(t, conn, buildCER(NewAVP(99998, avp.Mbit, 0, datatype.OctetString("mandatory"))))

	// Give server time to process
	time.Sleep(500 * time.Millisecond)

	// Connection must still be alive — send another message
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after M=1 unknown")
}

// Test 3: Unknown vendor-specific AVP — should NOT kill connection
func TestUnknownAVP_VendorSpecific(t *testing.T) {
	ln, handled := startTestServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Unknown vendor AVP (vendor 12345, code 55555, M=1, V=1)
	sendMsg(t, conn, buildCER(NewAVP(55555, avp.Mbit|avp.Vbit, 12345, datatype.OctetString("vendor-data"))))

	time.Sleep(500 * time.Millisecond)

	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after vendor unknown")
}

// Test 4: Unknown AVP inside a Grouped AVP — should NOT kill connection
func TestUnknownAVP_InsideGrouped(t *testing.T) {
	ln, handled := startTestServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Build a Grouped AVP (Vendor-Specific-Application-Id, code 260) containing
	// a valid Auth-Application-Id plus an unknown sub-AVP
	innerKnown := NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	innerUnknown := NewAVP(88888, avp.Mbit, 0, datatype.OctetString("nested-unknown"))

	// Manually build the grouped payload
	b1, _ := innerKnown.Serialize()
	b2, _ := innerUnknown.Serialize()
	grouped := append(b1, b2...)

	groupedAVP := NewAVP(260, avp.Mbit, 0, datatype.Grouped(grouped))

	sendMsg(t, conn, buildCER(groupedAVP))

	time.Sleep(500 * time.Millisecond)

	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after grouped unknown")
}

// Test 5: Unknown command code — should get 3001 answer with E-bit, NOT kill connection
func TestUnknownCommand(t *testing.T) {
	ln, handled := startTestServerWithIdentity(t,
		datatype.DiameterIdentity("srv.example.net"),
		datatype.DiameterIdentity("example.net"))
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// First send a valid CER so we know the connection works
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "initial CER")

	// Now send a message with unknown command code 9999
	m := NewRequest(9999, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test.host"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test.realm"))
	b, _ := m.Serialize()

	if _, err := conn.Write(b); err != nil {
		t.Fatalf("write unknown cmd: %v", err)
	}

	// Read the 3001 answer
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	ans, err := ReadMessage(conn, dict.Default)
	if err != nil && err != ErrCommandUnsupported {
		t.Fatalf("reading answer: %v", err)
	}
	if ans == nil {
		t.Fatal("no answer received for unknown command")
	}
	if ans.Header.CommandCode != 9999 {
		t.Errorf("answer command code: got %d, want 9999", ans.Header.CommandCode)
	}
	if ans.Header.CommandFlags&RequestFlag != 0 {
		t.Error("answer has R-bit set — should be cleared")
	}
	if ans.Header.CommandFlags&ErrorFlag == 0 {
		t.Error("answer missing E-bit — RFC 6733 requires it for 3001")
	}
	// Check Result-Code AVP = 3001
	rcAVP, err := findFromAVP(ans.AVP, avp.ResultCode, false)
	if err != nil || len(rcAVP) == 0 {
		t.Fatal("answer missing Result-Code AVP")
	}
	rc, ok := rcAVP[0].Data.(datatype.Unsigned32)
	if !ok {
		t.Fatalf("Result-Code data type: %T", rcAVP[0].Data)
	}
	if uint32(rc) != CommandUnsupported {
		t.Errorf("Result-Code: got %d, want %d", rc, CommandUnsupported)
	}
	// RFC 6733 §7.2 <answer-message> requires Origin-Host and Origin-Realm.
	ohAVP, err := findFromAVP(ans.AVP, avp.OriginHost, false)
	if err != nil || len(ohAVP) == 0 {
		t.Error("answer missing Origin-Host AVP (RFC 6733 §7.2)")
	} else if got, want := string(ohAVP[0].Data.(datatype.DiameterIdentity)), "srv.example.net"; got != want {
		t.Errorf("Origin-Host: got %q, want %q", got, want)
	}
	orAVP, err := findFromAVP(ans.AVP, avp.OriginRealm, false)
	if err != nil || len(orAVP) == 0 {
		t.Error("answer missing Origin-Realm AVP (RFC 6733 §7.2)")
	} else if got, want := string(orAVP[0].Data.(datatype.DiameterIdentity)), "example.net"; got != want {
		t.Errorf("Origin-Realm: got %q, want %q", got, want)
	}
	t.Logf("Got RFC-compliant 3001 answer with E-bit for unknown command 9999")

	// Connection should still be alive — send another valid message
	conn.SetReadDeadline(time.Time{}) // clear deadline
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after unknown command")
}

// Test 6: ReadMessage directly — verify unknown AVP is decoded as Unknown type
func TestReadMessage_UnknownAVP_Decoded(t *testing.T) {
	m := buildCER(
		NewAVP(99999, 0, 0, datatype.OctetString("m0-data")),
		NewAVP(99998, avp.Mbit, 0, datatype.OctetString("m1-data")),
	)
	b, err := m.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := ReadMessage(bytes.NewReader(b), dict.Default)
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	// Should have all AVPs including the unknown ones
	foundM0, foundM1 := false, false
	for _, a := range parsed.AVP {
		if a.Code == 99999 {
			foundM0 = true
			if a.Data.Type() != datatype.UnknownType {
				t.Errorf("AVP 99999: expected UnknownType, got %d", a.Data.Type())
			}
		}
		if a.Code == 99998 {
			foundM1 = true
			if a.Flags&avp.Mbit != avp.Mbit {
				t.Error("AVP 99998: M-bit not preserved")
			}
		}
	}
	if !foundM0 {
		t.Error("Unknown AVP 99999 (M=0) not found in parsed message")
	}
	if !foundM1 {
		t.Error("Unknown AVP 99998 (M=1) not found in parsed message")
	}
}

// Test 7: Craft raw bytes with unknown AVP that has string-type code lookup
// (exercises the FindAVPWithVendor path where code is NOT uint32)
func TestUnknownAVP_RawBytes_NonStandardVendor(t *testing.T) {
	// Build a CER, then manually inject a raw unknown AVP with a weird vendor
	m := buildCER()
	b, _ := m.Serialize()

	// Append a raw AVP: code=77777, flags=0xC0 (V+M), length=20, vendor=99999, data="test"
	rawAVP := make([]byte, 20)
	binary.BigEndian.PutUint32(rawAVP[0:4], 77777)          // code
	rawAVP[4] = avp.Vbit | avp.Mbit                          // flags
	copy(rawAVP[5:8], uint32to24(20))                         // length
	binary.BigEndian.PutUint32(rawAVP[8:12], 99999)          // vendor
	copy(rawAVP[12:16], []byte("test"))                       // data (4 bytes)

	// Update message length in header
	newLen := len(b) + len(rawAVP)
	b[0] = 1 // version
	b[1] = byte(newLen >> 16)
	b[2] = byte(newLen >> 8)
	b[3] = byte(newLen)
	b = append(b, rawAVP...)

	parsed, err := ReadMessage(bytes.NewReader(b), dict.Default)
	if err != nil {
		t.Fatalf("ReadMessage with raw vendor AVP failed: %v — this is the crash!", err)
	}

	found := false
	for _, a := range parsed.AVP {
		if a.Code == 77777 && a.VendorID == 99999 {
			found = true
			t.Logf("Found unknown vendor AVP: code=%d vendor=%d flags=0x%x data=%v",
				a.Code, a.VendorID, a.Flags, a.Data)
		}
	}
	if !found {
		t.Error("Unknown vendor AVP 77777/99999 not found in parsed message")
	}
}

// Test: Unknown command ANSWER (R-bit=0) must not kill the connection.
// Per RFC 6733 §6.2 spirit: no response to send, just discard.
func TestUnknownCommand_Answer(t *testing.T) {
	ln, handled := startTestServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Prove the connection works.
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "initial CER")

	// Build an "answer" with an unknown command code (R-bit cleared).
	m := NewMessage(9999, 0, 0, 0, 0, dict.Default) // flags=0 → answer, not request
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test.host"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test.realm"))
	m.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(2001))
	b, _ := m.Serialize()
	if _, err := conn.Write(b); err != nil {
		t.Fatalf("write unknown answer: %v", err)
	}

	// Server must not send a 3001 back for an answer, and must not close.
	// Wait briefly then prove the connection is still alive with a valid CER.
	time.Sleep(200 * time.Millisecond)
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "follow-up after unknown answer")
}

// Test: without an Origin-Host/Realm configured, the server cannot build
// a RFC-compliant 3001 answer. It must fall back to closing the
// connection rather than emitting a non-compliant half-answer.
func TestUnknownCommand_NoIdentityClosesConn(t *testing.T) {
	ln, handled := startTestServer(t) // no identity
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Prove the connection works.
	sendMsg(t, conn, buildCER())
	expectHandled(t, handled, "initial CER")

	// Send a request with an unknown command code.
	m := NewRequest(9999, 0, dict.Default)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test.host"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test.realm"))
	b, _ := m.Serialize()
	if _, err := conn.Write(b); err != nil {
		t.Fatalf("write unknown cmd: %v", err)
	}

	// Server must close the connection; our next read should yield EOF.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1)
	if n, err := conn.Read(buf); err == nil {
		t.Fatalf("expected EOF/closed connection, read %d bytes: 0x%x", n, buf[:n])
	}
}
