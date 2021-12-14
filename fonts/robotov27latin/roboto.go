package robotov27latin

import (
	"bytes"
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed index.css
var indexcss string

//go:embed roboto-v27-latin-regular.woff
var roboto_v27_latin_regular_woff []byte

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/roboto/css", "")
	gblrs.FS().MKDIR("/roboto/fonts", "")
	gblrs.FS().SET("/roboto/css/index.css", indexcss)
	gblrs.FS().SET("/roboto/fonts/roboto-v27-latin-regular.woff", bytes.NewReader(roboto_v27_latin_regular_woff))
	gblrs.FS().MKDIR("/roboto/html", "")
	gblrs.FS().SET("/roboto/html/head.html", `<link rel="stylesheet" type="text/css" href="/roboto/css/index.css">`)
}
