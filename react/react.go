package react

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/react.min.js
var reactjs string

//go:embed js/react-dom.min.js
var reactdomjs string

//go:embed js/react-jsonschema-form.js
var reactjsonschemaformjs string

func init() {
	gblrsfs := resources.GLOBALRSNG().FS()
	gblrsfs.MKDIR("/react/js", "")
	gblrsfs.MKDIR("/react", "")
	gblrsfs.SET("/react/js/react.min.js", reactjs)
	gblrsfs.SET("/react/js/react.js", reactjs)

	gblrsfs.SET("/react/js/react-dom.min.js", reactdomjs)
	gblrsfs.SET("/react/js/react-dom.js", reactdomjs)
	gblrsfs.SET("/react/js/react-jsonschema-form.js", reactjsonschemaformjs)

	gblrsfs.SET("/react/head.html",
		`<script crossorigin src="/react/js/react.min.js"></script>
<script crossorigin src="/react/js/react-dom.min.js"></script>`)
}
