package crypto

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()
	gblrsfs.MKDIR("/raw:crypto/js", "")
	gblrsfs.SET("/crypto/js/crypto-js.js", cryptojs)
	gblrsfs.MKDIR("/raw:crypto", "")
	gblrsfs.SET("/crypto/head.html", `<script type="application/javascript" src="/crypto/js/crypto-js.js"></script>`)
}

//go:embed crypto.js
var cryptojs string
