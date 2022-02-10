package ace

import (
	"embed"
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/*.js
var acels embed.FS

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:ace/js", acels)
	gblrsngfs.MKDIR("/raw:ace", "")
	gblrsngfs.SET("/ace/head.html", `<script type="application/javascript" src="/ace/js/ace.js"></script>`)
}
