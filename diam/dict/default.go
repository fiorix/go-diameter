// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package dict

import (
	"github.com/fiorix/go-diameter/v4/diam/dict/library"
)

// Default is a Parser object with pre-loaded dictionaries.
var Default *Parser

func init() {
	var err error
	Default, err = NewParserFromLibrary(
		library.Base,
		library.CreditControl,
		library.GxCreditControl,
		library.NetworkAccessServer,
		library.TgppRoRf,
		library.TgppS6a,
		library.TgppSwx,
	)
	if err != nil {
		panic(err)
	}
}
