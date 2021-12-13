package bootstrap

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/bootstrap/css", "")
	gblrs.FS().MKDIR("/bootstrap/js", "")
	gblrs.FS().SET("/bootstrap/css/bootstrap.css", bootstrapcss)
	gblrs.FS().SET("/bootstrap/css/bootstrap.min.css", bootstrapjs)
	gblrs.FS().SET("/bootstrap/js/bootstrap.js", bootstrapjs)
	gblrs.FS().SET("/bootstrap/js/bootstrap.min.js", bootstrapjs)
	gblrs.FS().SET("/bootstrap/js/bootstrap.bundle.js", bootstrapjs)
	gblrs.FS().SET("/bootstrap/js/bootstrap.bundle.min.js", bootstrapjs)

	gblrs.FS().MKDIR("/bootstrap/html", "")
	gblrs.FS().SET("/bootstrap/html/head.html",
		`<link rel="stylsheet" type="text/css" href="/bootstrap/css/bootstrap.min.css">
<script type="application/javascript" src="/bootstrap/js/bootstrap.bundle.min.js"></script>`)

}

//go:embed js/bootstrap.bundle.min.js
var bootstrapjs string

//go:embed css/bootstrap.min.css
var bootstrapcss string
