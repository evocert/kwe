package jspanel

import (
	"github.com/evocert/kwe/resources"

	_ "embed"
)

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()

	gblrsfs.MKDIR("/raw:jspanel/css", "")
	gblrsfs.MKDIR("/raw:jspanel/js", "")
	gblrsfs.SET("/jspanel/js/jspanel.min.js", jspaneljs)
	gblrsfs.SET("/jspanel/js/jspanel.js", jspaneljs)
	gblrsfs.SET("/jspanel/css/jspanel.min.css", jspanelcss)
	gblrsfs.SET("/jspanel/css/jspanel.css", jspanelcss)

	gblrsfs.MKDIR("/raw:jspanel/js/extensions/modal", "")
	gblrsfs.SET("/jspanel/js/extensions/modal/jspanel.modal.min.js", jspanelmodaljs)
	gblrsfs.SET("/jspanel/js/extensions/modal/jspanel.modal.js", jspanelmodaljs)

	gblrsfs.MKDIR("/raw:jspanel", "")
	gblrsfs.SET("/jspanel/head.html",
		`<link rel="stylesheet" type="text/css" href="/jspanel/css/jspanel.min.css">
<script type="application/javascript" src="/jspanel/js/jspanel.min.js"></script>
<script type="application/javascript" src="/jspanel/js/extensions/modal/jspanel.modal.min.js"></script>`)
}

//go:embed css/jspanel.min.css
var jspanelcss string

//go:embed js/jspanel.min.js
var jspaneljs string
