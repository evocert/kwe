package datepicker

import (
	"github.com/evocert/kwe/resources"

	_ "embed"
)

// embedded https://cdn.jsdelivr.net/npm/js-datepicker@5.18.0/dist/datepicker.min.css
//go:embed css/datepicker.css
var datepickercss string

// embedded https://cdn.jsdelivr.net/npm/js-datepicker@5.18.0/dist/datepicker.min.js
//go:embed js/datepicker.js
var datepickerjs string

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/raw:datepicker/css", "")
	gblrs.FS().MKDIR("/raw:datepicker/js", "")
	gblrs.FS().MKDIR("/raw:datepicker", "")
	gblrs.FS().SET("/datepicker/head.html",
		`<link rel="stylesheet" type="ext/css" href="/datepicker/css/datepicker.css">
<script type="application/javascript" src="/datepicker/js/datepicker.js"></script>`)
	gblrs.FS().SET("/datepicker/css/datepicker.css", datepickercss)
	gblrs.FS().SET("/datepicker/js/datepicker.js", datepickerjs)
}
