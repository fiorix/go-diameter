//go:build linux && !386

// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"testing"

	"github.com/fiorix/go-diameter/v4/diam"
)

func requireSCTP(t *testing.T) {
	t.Helper()
	ln, err := diam.MultistreamListen("sctp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("SCTP not available: %v", err)
	}
	ln.Close()
}

func TestHandleCER_HandshakeMetadataSCTP(t *testing.T) {
	requireSCTP(t)
	testHandleCER_HandshakeMetadata(t, "sctp")
}

func testClient_Handshake_CustomIP_SCTP(t *testing.T) {
	requireSCTP(t)
	testClient_Handshake_CustomIP(t, "sctp")
}

// TestStateMachineSCTP establishes a connection with a test SCTP server and
// sends a Re-Auth-Request message to ensure the handshake was
// completed and that the RAR handler has context from the peer.
func TestStateMachineSCTP(t *testing.T) {
	requireSCTP(t)
	testStateMachine(t, "sctp")
}
