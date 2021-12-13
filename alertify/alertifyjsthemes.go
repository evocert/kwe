package alertify

import (
	"io"
	"strings"

	_ "embed"

	"github.com/evocert/kwe/resources"
)

//go:embed css/themes/default.min.css
var alertifythemesdefaultcss string

//go:embed css/themes/default.rtl.min.css
var alertifythemesdefaultrtlcss string

func AlertifyThemesDefaultCSS() io.Reader {
	return strings.NewReader(alertifythemesdefaultcss)
}

func AlertifyThemesDefaultRtlCSS() io.Reader {
	return strings.NewReader(alertifythemesdefaultrtlcss)
}

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/alertify/css/themes", "")
	gblrs.FS().SET("/alertify/css/themes/default.css", alertifythemesdefaultcss)
	gblrs.FS().SET("/alertify/css/themes/default.min.css", alertifythemesdefaultcss)
	gblrs.FS().SET("/alertify/css/themes/default.rtl.css", alertifythemesdefaultrtlcss)
	gblrs.FS().SET("/alertify/css/themes/default.rtl.min.css", alertifythemesdefaultrtlcss)
}
