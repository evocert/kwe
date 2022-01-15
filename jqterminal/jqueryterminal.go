package jqterminal

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/jquery.terminal.js
var jqueryterminaljs string

//go:embed css/jquery.terminal.css
var jqueryterminalcss string

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/jqterminal/html", "")
	gblrs.FS().MKDIR("/jqterminal/css", "")
	gblrs.FS().MKDIR("/jqterminal/js", "")
	gblrs.FS().SET("/jqterminal/css/jquery.terminal.css", jqueryterminalcss)
	gblrs.FS().SET("/jqterminal/js/jquery.terminal.js", jqueryterminaljs)
	gblrs.FS().SET("/jqterminal/html/head.html", `<link rel="stylsheet" type="text/css" href="/jqterminal/css/jquery.terminal.css">
	<script type="application/javascript" src="/jqterminal/js/jquery.terminal.js"></script>`)
}
