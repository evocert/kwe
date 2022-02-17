package tablesorter

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/tablesorter.min.js
var tablesorterjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:tablesorter/js", "")
	gblrsngfs.SET("/tablesorter/js/tablesorter.min.js", tablesorterjs)
	gblrsngfs.MKDIR("/raw:tablesorter", "")
	gblrsngfs.SET("/tablesorter/head.html", `<script type="application/javascript" src="/tablesorter/js/tablesorter.min.js"></script>`)

}
