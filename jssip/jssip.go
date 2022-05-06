package jssip

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/jssip.min.js
var jssipjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()

	gblrsngfs.MKDIR("/raw:jssip/js")
	gblrsngfs.MKDIR("/raw:jssip")
	gblrsngfs.SET("/jssip/js/jssip.min.js", jssipjs)
	gblrsngfs.SET("/jssip/head.html", `<script src="/jssip/js/jssip.min.js"></script>`)
}
