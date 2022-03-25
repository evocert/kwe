package crossfilter

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/crossfilter2.min.js
var crossfilterjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:crossfilter/js", "")
	gblrsngfs.MKDIR("/raw:crossfilter", "")
	gblrsngfs.SET("/crossfilter/js/crossfilter.min.js", crossfilterjs)
	gblrsngfs.SET("/crossfilter/head.html", `<script type="application/javascript" src="/crossfilter/js/crossfilter.min.js"></script>`)
}
