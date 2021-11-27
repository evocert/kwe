package html

import (
	"github.com/evocert/kwe/requirejs"
	"github.com/evocert/kwe/resources"
)

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()
	gblrsfs.MKDIR("/require/js", "")
	gblrsfs.SET("/require/js/require.js", requirejs.RequireJS())
	gblrsfs.MKDIR("/require/html", "")
	gblrsfs.SET("/require/html/head.html", `<script type="application/javascript" src="/require/js/require.js"></script>`)
}
