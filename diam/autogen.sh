#!/bin/sh -e
# Copyright 2013-2015 go-diameter authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
#
# Generate Diameter constants from our dictionaries.
#
# Run `sh autogen.sh` to re-generate these files after changing
# dictionary XML files.

dict=dict/testdata/*.xml


## Generate commands.go
src=commands.go

cat << EOF > $src
// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package diam

// Diameter command codes.
const (
EOF

cat $dict | sed \
	-e 's/-//g' \
	-ne 's/.*command code="\(.*\)" .* name="\(.*\)".*/\2 = \1/p' \
	| sort -u >> $src

echo ')' >> $src

go fmt $src


## Generate avp/codes.go
src=avp/codes.go

cat << EOF > $src
// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package avp

// Diameter AVP types.
const (
EOF

cat $dict | sed \
	-e 's/-Id/-ID/g' \
	-e 's/-//g' \
	-ne 's/.*avp name="\(.*\)" code="\([0-9]*\)".*/\1 = \2/p' \
	| sort -u >> $src

echo ')' >> $src

go fmt $src


## Generate dict/default.go
src=dict/default.go

cat << EOF > $src
// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package dict

import "bytes"

// Default is a Parser object with pre-loaded
// Base Protocol and Credit Control dictionaries.
var Default *Parser

func init() {
	Default, _ = NewParser()
	Default.Load(bytes.NewReader([]byte(base_xml)))
	Default.Load(bytes.NewReader([]byte(credit_control_xml)))
}

EOF

for f in $dict
do

var=`basename $f | sed 's/\.xml/_xml/g'`
cat << EOF >> $src
var $var=\``cat $f`\`

EOF

done

go fmt $src
