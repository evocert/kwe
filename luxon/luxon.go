package luxon

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()
	gblrsfs.MKDIR("/raw:luxon/js", "")
	gblrsfs.SET("/luxon/js/luxon.js", luxonjs)
	gblrsfs.SET("/luxon/js/luxon.min.js", luxonjs)
	gblrsfs.MKDIR("/raw:luxon", "")
	gblrsfs.SET("/luxon/head.html", `<script type="appliaction/javascript" src="/luxon/js/luxon.js"></script>`)
}

//go:embed js/lunix.min.js
var luxonjs string
