package ml5

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/ml5.min.js
var ml5js string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/ml5/js", "")
	gblrsngfs.MKDIR("/ml5", "")
	gblrsngfs.SET("/ml5/js/ml5.min.js", ml5js)
	gblrsngfs.SET("/ml5/head.html", `<script type="application/javascript" src="/ml5/js/ml5.min.js"></script>`)
}
