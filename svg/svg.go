package svg

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/svg.min.js
var svgjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:svg/js", "")
	gblrsngfs.SET("/svg/js/svg.min.js", svgjs)
	gblrsngfs.MKDIR("/raw:svg", "")
	gblrsngfs.SET("/svg/head.html", `<script src="/svg/js/svg.min.js"></script>`)
}
