// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import "crypto/tls"

// TLSConfigClone returns a deep copy of cfg, or nil if cfg is nil.
func TLSConfigClone(cfg *tls.Config) *tls.Config {
	if cfg != nil {
		return cfg.Clone()
	}
	return nil
}
