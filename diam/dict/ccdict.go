// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Default dictionary Parser and Base Protocol XML.

package dict

import "bytes"

// Default is a static Parser with a pre-loaded Base Protocol.
var CCDict *Parser

func init() {
	CCDict, _ = NewParser()
	CCDict.Load(bytes.NewReader(CCDictXML))
}

// DefaultXML is an embedded version of the Diameter Base Protocol.
//
// Copy of diam_base.xml
var CCDictXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<diameter>

  <application id="4"> <!-- Diameter Credit-Control Application -->
  </application>
</diameter>`)
