package jqueryui

import (
	"bytes"
	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed index.html
var indexhtml string

//go:embed jquery-ui.min.css
var jqueryuicss string

//go:embed jquery-ui.min.js
var jqueryuijs string

//go:embed images/ui-icons_444444_256x240.png
var uiicons_444444_256x240pngb64 []byte

//go:embed images/ui-icons_555555_256x240.png
var uiicons_555555_256x240pngb64 []byte

//go:embed images/ui-icons_777620_256x240.png
var uiicons_777620_256x240pngb64 []byte

//go:embed images/ui-icons_777777_256x240.png
var uiicons_777777_256x240pngb64 []byte

//go:embed images/ui-icons_cc0000_256x240.png
var uiicons_cc0000_256x240pngb64 []byte

//go:embed images/ui-icons_ffffff_256x240.png
var uiicons_ffffff_256x240pngb64 []byte

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/raw:jquery-ui/html", "")
	gblrsngfs.SET("/jquery-ui/html/head.html", `<link href="/jquery-ui/jquery-ui.min.css" rel="stylesheet">
<script src="/jquery/jquery.min.js"></script>
<script src="/jquery-ui/jquery-ui.min.js"></script>`)
	gblrsngfs.MKDIR("/raw:jquery-ui", "")
	gblrsngfs.SET("/jquery-ui/index.html", indexhtml)
	gblrsngfs.SET("/jquery-ui/jquery-ui.min.css", jqueryuicss)
	gblrsngfs.SET("/jquery-ui/jquery-ui.min.js", jqueryuijs)
	gblrsngfs.SET("/jquery-ui/jquery-ui.css", jqueryuicss)
	gblrsngfs.SET("/jquery-ui/jquery-ui.js", jqueryuijs)
	gblrsngfs.MKDIR("/jquery-ui/images", "")
	gblrsngfs.SET("/jquery-ui/images/ui-icons_444444_256x240.png", bytes.NewReader(uiicons_444444_256x240pngb64))
	gblrsngfs.SET("/jquery-ui/images/ui-icons_555555_256x240.png", bytes.NewReader(uiicons_555555_256x240pngb64))
	gblrsngfs.SET("/jquery-ui/images/ui-icons_777620_256x240.png", bytes.NewReader(uiicons_777620_256x240pngb64))
	gblrsngfs.SET("/jquery-ui/images/ui-icons_777777_256x240.png", bytes.NewReader(uiicons_777777_256x240pngb64))
	gblrsngfs.SET("/jquery-ui/images/ui-icons_cc0000_256x240.png", bytes.NewReader(uiicons_cc0000_256x240pngb64))
	gblrsngfs.SET("/jquery-ui/images/ui-icons_ffffff_256x240.png", bytes.NewReader(uiicons_ffffff_256x240pngb64))
}
