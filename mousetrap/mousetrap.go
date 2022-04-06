package mousetrap

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/mousetrap.min.js
var mousetrapjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/mousestrap/js", "")
	gblrsngfs.SET("//mousestrap/js/mousetrap.min.js", mousetrapjs)
	gblrsngfs.MKDIR("/mousestrap", "")
	gblrsngfs.SET("/mousetrap/head.html", `<script type="application/javascript" src="/mousetrap/js/mousetrap.min.js"></script>`)
}
