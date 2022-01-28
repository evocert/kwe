package babel

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/babel.min.js
var babeljs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:babel/js", "")
	gblrsngfs.MKDIR("/raw:babel", "")
	gblrsngfs.SET("/babel/js/babel.min.js", babeljs)
	gblrsngfs.SET("/babel/head.html", `<script crossorigin type="application/javascript" src="/babel/js/babel.min.js"></script>`)

}
