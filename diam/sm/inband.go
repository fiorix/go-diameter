// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import "context"

// inbandSecurityKey keys the security mechanism this side requested in its
// CER. It is recorded per connection, because the answer has to be compared
// against what was actually asked for: a peer that omits Inband-Security-Id,
// or answers NO_INBAND_SECURITY, has not agreed to an upgrade.
type inbandSecurityKey struct{}

// withRequestedInbandSecurity records the mechanism sent in the CER.
func withRequestedInbandSecurity(ctx context.Context, id uint32) context.Context {
	return context.WithValue(ctx, inbandSecurityKey{}, id)
}

// requestedInbandSecurity returns the mechanism recorded by
// withRequestedInbandSecurity, or 0 (NO_INBAND_SECURITY) when none was.
func requestedInbandSecurity(ctx context.Context) uint32 {
	if v, ok := ctx.Value(inbandSecurityKey{}).(uint32); ok {
		return v
	}
	return 0
}
