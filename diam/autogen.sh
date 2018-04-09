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

echo ')\n// Short Command Names\nconst (\n' >> $src

cat $dict | sed \
	-e 's/-//g' \
	-ne 's/.*command code="[0-9]*".*\s.*short="\([^"]*\).*/\1R = "\1R"\n\1A = "\1A"/p' \
	| sort -u >> $src

echo ')' >> $src
go fmt $src

## Generate applications.go
src=applications.go

cat << EOF > $src
// Copyright 2013-2018 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package diam

// Diameter application IDs.
const (
EOF

cat $dict | sed \
    -e :1 -e 's/\("[^"]*\)[[:space:]]\([^"]*"\)/\1_\2/g;t1' \
    -ne 's/\s*<application\s*id="\([0-9]*\)".*name="\(.*\)".*/\U\2_APP_ID = \1/p' \
    | sort -u | sort -nk 3 >> $src

echo ')\n' >> $src
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
	-e 's/-Id\([-"s]\)/-ID\1/g' \
	-e 's/-//g' \
	-ne 's/.*avp name="\(.*\)" code="\([0-9]*\)".*/\1 = \2/p' \
	| sort -u >> $src

echo ')\n' >> $src

go fmt $src


## Generate dict/default.go
src=dict/default.go

cat << EOF > $src
// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This file is auto-generated from our dictionaries.

package dict

import (
	"bytes"
	"fmt"
)

// Default is a Parser object with pre-loaded
// Base Protocol and Credit Control dictionaries.
var Default *Parser

func init() {
	var dictionaries = []struct{ name, xml string }{
		{"Base", baseXML},
		{"Credit Control", creditcontrolXML},
		{"Network Access Server", networkaccessserverXML},
		{"TGPP", tgpprorfXML},
		{"TGPP_S6a", tgpps6aXML},
	}
	var err error
	Default, err = NewParser()
	if err != nil {
		panic(err)
	}
	for _, dict := range dictionaries {
		err = Default.Load(bytes.NewReader([]byte(dict.xml)))
		if err != nil {
			panic(fmt.Sprintf("Cannot load %s dictionary: %s", dict.name, err))
		}
	}
}

EOF

for f in $dict
do

var=`basename $f | sed -e 's/\.xml/XML/g' -e 's/_//g'`
cat << EOF >> $src
var $var=\``cat $f`\`

EOF

done

go fmt $src
