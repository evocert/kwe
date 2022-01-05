package build

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed build.js
var buildjs string

func init() {
	glbrsngfs := resources.GLOBALRSNG().FS()
	glbrsngfs.MKDIR("/kwe/build/js", "")
	glbrsngfs.SET("/kwe/build/js/build.js", buildjs)
}
