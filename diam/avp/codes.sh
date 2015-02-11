#!/bin/sh -e
# Copyright 2013-2015 go-diameter authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
#
# Generate DIAMETER AVP constants from our dictionaries.

src=codes.go

cat << EOF > $src
// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

// AVP types. Auto-generated from our dictionaries.
const (
EOF

dict=../dict/testdata/*.xml
cat $dict | sed \
	-e 's/-Id/-ID/g' \
	-e 's/-//g' \
	-ne 's/.*avp name="\(.*\)" code="\([0-9]*\)".*/\1 = \2/p' | \
	sort -u >> $src

echo ')' >> $src

go fmt $src
