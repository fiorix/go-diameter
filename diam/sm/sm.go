// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"fmt"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/sm/smpeer"
)

// Settings used to configure the state machine.
type Settings struct {
	OriginHost       datatype.DiameterIdentity
	OriginRealm      datatype.DiameterIdentity
	VendorID         datatype.Unsigned32
	ProductName      datatype.UTF8String
	FirmwareRevision datatype.Unsigned32
}

// StateMachine is a specialized type of diam.ServeMux that handles
// the CER/CEA handshake and DWR/DWA messages for clients or servers.
//
// Other handlers registered in the state machine are only executed
// after the peer has passed the initial CER/CEA handshake.
type StateMachine struct {
	cfg       *Settings
	mux       *diam.ServeMux
	hsNotifyc chan diam.Conn // handshake notifier
}

// New creates and initializes a new StateMachine for clients or servers.
func New(settings *Settings) *StateMachine {
	sm := &StateMachine{
		cfg:       settings,
		mux:       diam.NewServeMux(),
		hsNotifyc: make(chan diam.Conn),
	}
	sm.mux.Handle("CER", handleCER(sm))
	sm.mux.Handle("DWR", handshakeOK(handleDWR(sm)))
	return sm
}

// Settings return the Settings object used by this StateMachine.
func (sm *StateMachine) Settings() *Settings {
	return sm.cfg
}

// ServeDIAM implements the diam.Handler interface.
func (sm *StateMachine) ServeDIAM(c diam.Conn, m *diam.Message) {
	sm.mux.ServeDIAM(c, m)
}

// Handle implements the diam.Handler interface.
func (sm *StateMachine) Handle(cmd string, handler diam.Handler) {
	sm.HandleFunc(cmd, handler.ServeDIAM)
}

// HandleFunc implements the diam.Handler interface.
func (sm *StateMachine) HandleFunc(cmd string, handler diam.HandlerFunc) {
	switch cmd {
	case "CER", "CEA", "DWR", "DWA":
		sm.Error(&diam.ErrorReport{
			Error: fmt.Errorf("cannot overwrite %s command in the state machine", cmd),
		})
	default:
		sm.mux.Handle(cmd, handshakeOK(handler))
	}
}

// Error implements the diam.ErrorReporter interface.
func (sm *StateMachine) Error(err *diam.ErrorReport) {
	sm.mux.Error(err)
}

// ErrorReports implement the diam.ErrorReporter interface.
func (sm *StateMachine) ErrorReports() <-chan *diam.ErrorReport {
	return sm.mux.ErrorReports()
}

// HandshakeNotify implements the HandshakeNotifier interface.
func (sm *StateMachine) HandshakeNotify() <-chan diam.Conn {
	return sm.hsNotifyc
}

// The HandshakeNotifier interface is implemented by Handlers
// that allow detecting peers that have passed the CER/CEA
// handshake.
type HandshakeNotifier interface {
	// HandshakeNotify returns a channel that receives
	// a peer's diam.Conn after it passes the handshake.
	HandshakeNotify() <-chan diam.Conn
}

// handshakeOK is a wrapper for state machine handlers that only
// calls the designated handler function if the peer has passed the
// CER/CEA handshake.
type handshakeOK diam.HandlerFunc

// ServeDIAM implements the diam.Handler interface.
func (f handshakeOK) ServeDIAM(c diam.Conn, m *diam.Message) {
	if _, ok := smpeer.FromContext(c.Context()); ok {
		f(c, m)
	}
}
