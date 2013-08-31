// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package util

// Pad4 returns n padded to 4 bytes.
func Pad4(n uint32) uint32 {
	return n + ((4 - n) & 3)
}
