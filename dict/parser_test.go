// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Tests

package dict

import "testing"

func TestNew(t *testing.T) {
	if _, err := New("./diam_base.xml"); err != nil {
		t.Error(err)
	}
}
