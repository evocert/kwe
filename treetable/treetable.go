package treetable

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/jquery.treetable.min.js
var jquerytreetablejs string

//go:embed css/jquery.treetable.min.css
var jquerytreetablecss string

//go:embed css/jquery.treetable.theme.default.min.css
var jquerytreetablethemedeaultcss string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/treetable/js", "")
	gblrsngfs.SET("/treetable/js/jquery.treetable.js", jquerytreetablejs)
	gblrsngfs.MKDIR("/treetable/css", "")
	gblrsngfs.SET("/treetable/css/jquery.treetable.css", jquerytreetablecss)
	gblrsngfs.SET("/treetable/css/jquery.treetable.theme.default.css", jquerytreetablethemedeaultcss)
	gblrsngfs.MKDIR("/treetable", "")
	gblrsngfs.SET("/treetable/head.html",
		`<link rel="stylesheet" type="text/css" href="/treetable/css/jquery.treetable.css">
	<link rel="stylesheet" type="text/css" href="/treetable/css/jquery.treetable.theme.default.css">
	<script type="application/text" src="/treetable/js/jquery.treetable.js"></script>`)
}
