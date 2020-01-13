// +build linux,!386

// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtest

import "testing"

func TestNewServerSCTP(t *testing.T) {
	srv := NewServerNetwork("sctp", nil, nil)
	srv.Close()
}
