// +build linux,!386

// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import "testing"

func TestHandleCER_VS_AuthSCTP(t *testing.T) {
	testHandleCER_VS_Auth(t, "sctp")
}
