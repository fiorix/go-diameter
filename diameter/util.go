// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diameter

// uint24to32 converts b from BigEndian to LittleEndian and return a packed uint32.
func uint24to32(b [3]uint8) uint32 {
	return uint32(b[2]) | (uint32(b[1]) << 8) | (uint32(b[0]) << 16)
}

// pad4 returns n padded to 4 bytes
func pad4(n uint32) uint32 {
	return n + ((4 - n) & 3)
}
