package ammo

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/ammo.js
var ammojs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()

	gblrsngfs.MKDIR("/ammo", "")
	gblrsngfs.SET("/ammo/head.html", `<script type="application/text" src="/ammo/js/ammo.js"></scrip>`)
	gblrsngfs.MKDIR("/ammo/js", "")
	gblrsngfs.SET("/ammo/js/ammo.js", ammojs)
}
