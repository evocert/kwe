package mousetrap

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/mousetrap.min.js
var mousetrapjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/mousetrap/js", "")
	gblrsngfs.SET("/mousetrap/js/mousetrap.min.js", mousetrapjs)
	gblrsngfs.MKDIR("/mousetrap", "")
	gblrsngfs.SET("/mousetrap/head.html", `<script type="application/javascript" src="/mousetrap/js/mousetrap.min.js"></script>`)
}
