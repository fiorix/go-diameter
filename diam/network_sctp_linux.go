// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build linux && !386

package diam

import (
	"syscall"

	"github.com/ishidawataru/sctp"
)

// setSCTPNoDelay disables Nagle-like buffering so small Diameter messages are
// sent immediately.
//
// It sets the option inside the socket's syscall.RawConn.Control with a typed
// setter. sctp.SCTPConn.Setsockopt takes the value as a uintptr, which does
// not keep a stack variable in place: if the goroutine stack grows before the
// syscall runs, the kernel reads the old, freed stack.
func setSCTPNoDelay(c *sctp.SCTPConn) {
	rc, err := c.SyscallConn()
	if err != nil {
		return
	}
	_ = rc.Control(func(fd uintptr) {
		_ = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY, 1)
	})
}
