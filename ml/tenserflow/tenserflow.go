package tenserflow

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/tf.min.js
var tfjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:tflw/js", "")
	gblrsngfs.MKDIR("/raw:tflw", "")
	gblrsngfs.SET("/tflw/js/tf.min.js", "var global={};\r\n", tfjs)
	gblrsngfs.SET("/tflw/head.html", `<script type="application/javascript" src="/tflw/js/tf.min.js"></script>`)
}
