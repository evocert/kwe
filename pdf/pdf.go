package pdf

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed pdf.js
var pdfjs string

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:pdf/js", "")
	gblrsngfs.SET("/pdf/js/pdf.js", pdfjs)
	gblrsngfs.SET("/pdf/js/pdf.min.js", pdfjs)
	gblrsngfs.MKDIR("/raw:pdf", "")
	gblrsngfs.SET("/pdf/head.html", `<script type="application/javascript" src="/pdf/js/pdf.min.js"></script>`)
}
