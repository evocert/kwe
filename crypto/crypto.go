package crypto

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()
	gblrsfs.MKDIR("/crypto/js", "")
	gblrsfs.SET("/crypto/js/crypto-js.js", cryptojs)
	gblrsfs.MKDIR("/crypto/html", "")
	gblrsfs.SET("/crypto/html/head.html", `<script type="application/javascript" src="/crypto/js/crypto-js.js"></script>`)
}

//go:embed crypto.js
var cryptojs string
