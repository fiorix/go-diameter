// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Default dictionary Parser and Base Protocol XML.

package dict

import "bytes"

// CreditControl is a static Parser with a pre-loaded Base Protocol.
var CreditControl *Parser

func init() {
	CreditControl, _ = NewParser()
	CreditControl.Load(bytes.NewReader(CreditControlXML))
}

// CreditControlXML is an embedded version of the Diameter Credit-Control Application
//
// Copy of diam_credit_control.xml
var CreditControlXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="4"> <!-- Diameter Credit-Control Application -->
  </application>
</diameter>`)
