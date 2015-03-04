package command

import (
	"bytes"

	"github.com/fiorix/go-diameter/diam/dict"
)

func init() {
	dict.Default.Load(bytes.NewReader([]byte(acctDictionary)))
	dict.Default.Load(bytes.NewReader([]byte(authDictionary)))
}

var acctDictionary = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
	<application id="1001" type="acct">
	</application>
</diameter>
`

var authDictionary = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
	<application id="1002" type="auth">
	</application>
</diameter>
`
