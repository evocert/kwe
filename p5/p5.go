package p5

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/p5.min.js
var p5js string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/p5/js", "")
	gblrsngfs.MKDIR("/p5", "")
	gblrsngfs.SET("/p5/js/p5.min.js", p5js)
	gblrsngfs.SET("/p5/head.html", `<script type="application/javascript" src="/p5/js/p5.min.js"></script>`)
}
