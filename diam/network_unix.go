// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package diam

import (
	"log"
	"strings"
	"syscall"
)

// setReuseTcpAddr is a "best effort" attempt to set TCP socket SO_REUSEADDR option
// it always returns nil & logs an error if syscall fails
func setReuseTcpAddr(network, address string, c syscall.RawConn) error {
	if c != nil && len(address) > 0 && strings.HasPrefix(network, "tcp") {
		c.Control(func(fd uintptr) {
			err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if err != nil {
				log.Printf("Setting SO_REUSEADDR for %s:%s error: %v", network, address, err)
			}
		})
	}
	return nil
}
