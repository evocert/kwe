package paper

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/paper-full.min.js
var paperjs string

func init() {
	glbrscngmngrfs := resources.GLOBALRSNG().FS()
	glbrscngmngrfs.MKDIR("/paper/js", "")
	glbrscngmngrfs.MKDIR("/paper", "")
	glbrscngmngrfs.SET("/paper/js/paper.min.js", paperjs)
	glbrscngmngrfs.SET("/paper/js/paper.js", paperjs)
	glbrscngmngrfs.SET("/paper/head.html", `<script type="application/javascript" src="/paper/js/paper.min.js"></script>`)
}
