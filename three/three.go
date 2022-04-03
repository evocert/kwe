package three

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/three.js
var threejs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()

	gblrsngfs.MKDIR("/three", "")
	gblrsngfs.SET("/three/head.html", `<script type="application/text" src="/three/js/three.js"></scrip>`)
	gblrsngfs.MKDIR("/three/js", "")
	gblrsngfs.SET("/three/js/three.js", threejs)
}
