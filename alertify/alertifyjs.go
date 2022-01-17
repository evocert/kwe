package alertify

import (
	"io"
	"strings"

	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed js/alertify.min.js
var alertifyjs string

//go:embed css/alertify.min.css
var alertifycss string

func AlertifyJS() io.Reader {
	return strings.NewReader(alertifyjs)
}

func AlertifyCSS() io.Reader {
	return strings.NewReader(alertifycss)
}

const alertifyheadhtml string = `<link rel="stylesheet" href="/alertify/css/alertify.css">
<script type="application/javascript" src="/alertify/js/alertify.js"></script>`

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/raw:alertify/html", "")
	gblrs.FS().MKDIR("/raw:alertify/css", "")
	gblrs.FS().MKDIR("/raw:alertify/js", "")
	gblrs.FS().SET("/alertify/css/alertify.css", alertifycss)
	gblrs.FS().SET("/alertify/js/alertify.js", alertifyjs)

	gblrs.FS().SET("/alertify/html/head.html", strings.NewReader(alertifyheadhtml))
}
