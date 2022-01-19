package babylon

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/babylon.js
var babylonjs string

func init() {
	rsmngrfs := resources.GLOBALRSNG().FS()
	rsmngrfs.MKDIR("/raw:babylon/js", "")
	rsmngrfs.SET("/babylon/js/babylon.js", babylonjs)
	rsmngrfs.MKDIR("/raw:babylon", "")
	rsmngrfs.SET("/babylon/head.html",
		`<script type="application/javascript" src="/babylon/js/babylon.js"></script>`)
}
