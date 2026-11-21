// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build !linux || 386

package diam

import "github.com/ishidawataru/sctp"

// setSCTPNoDelay is a no-op where github.com/ishidawataru/sctp has no SCTP
// support.
func setSCTPNoDelay(*sctp.SCTPConn) {}
