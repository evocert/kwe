package d3

import (
	"strings"

	"github.com/evocert/kwe/resources"

	_ "embed"
)

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:d3/js", "")
	gblrsngfs.SET("/d3/d3.js", strings.Replace(d3js, "|'|", "`", -1))
	gblrsngfs.SET("/d3/d3.min.js", strings.Replace(d3js, "|'|", "`", -1))
	gblrsngfs.MKDIR("/raw:d3", "")
	gblrsngfs.SET("/d3/head.html", `<script type="application/javascript" src="/d3/js/d3.js"></script>`)
}

//go:embed js/d3.min.js
var d3js string
