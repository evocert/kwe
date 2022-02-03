package x2js

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/x2js.min.js
var x2jsjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/x2js/js", "")
	gblrsngfs.MKDIR("/x2js", "")
	gblrsngfs.SET("/x2js/head.html",
		`<script type="application/javascript" src="/x2js/js/x2js.min.js"></script>`)
	gblrsngfs.SET("x2js/js/x2js.min.js", x2jsjs)
	gblrsngfs.SET("x2js/js/x2js.js", x2jsjs)
}
