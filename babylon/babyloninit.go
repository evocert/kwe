package babylon

import (
	"github.com/evocert/kwe/resources"
)

func init() {
	rsmngrfs := resources.GLOBALRSNG().FS()
	rsmngrfs.MKDIR("/raw:babylon/js", "")
	rsmngrfs.SET("/babylon/js/babylon.js", babylonjs)
	rsmngrfs.MKDIR("/raw:babylon/html", "")
	rsmngrfs.SET("/babylon/html/head.html",
		`<script type="application/javascript" src="/babylon/js/babylon.js"></script>`)
}
