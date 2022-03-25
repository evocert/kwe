package dc

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed css/dc.min.css
var dccss string

//go:embed js/dc.min.js
var dcjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:dc/css", "")
	gblrsngfs.MKDIR("/raw:dc/js", "")
	gblrsngfs.MKDIR("/raw:dc", "")
	gblrsngfs.SET("/dc/js/dc.min.js", dcjs)
	gblrsngfs.SET("/dc/css/dc.min.css", dccss)
	gblrsngfs.SET("/dc/head.html", `<link rel="stylesheet" href="/dc/css/dc.min.css"><script type="application/javascript" src="/dc/js/dc.min.js"></script>`)
}
