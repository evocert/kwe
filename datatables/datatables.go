package datatables

import (
	"bytes"

	"github.com/evocert/kwe/resources"

	_ "embed"
)

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()

	gblrsngfs.MKDIR("/raw:datatables", "")
	gblrsngfs.SET("/datatables/head.html", `<link rel="stylesheet" type="text/css" href="/datatables/css/jquery.Datatables.min.css">
	<script type="application/javascript" src="/datatables/js/jquery.Datatables.min.js"></script>`)

	gblrsngfs.MKDIR("/raw:datatables/css", "")
	gblrsngfs.SET("/datatables/css/jquery.Datatables.css", datatablescss)
	gblrsngfs.SET("/datatables/css/jquery.Datatables.min.css", datatablescss)

	gblrsngfs.MKDIR("/raw:datatables/js", "")
	gblrsngfs.SET("/datatables/js/jquery.Datatables.js", datatablesjs)
	gblrsngfs.SET("/datatables/js/jquery.Datatables.min.js", datatablesjs)

	gblrsngfs.MKDIR("/raw:datatables/images", "")

	gblrsngfs.SET("/datatables/images/details_open.png", bytes.NewReader(details_openpngbytes))
	gblrsngfs.SET("/datatables/images/details_close.png", bytes.NewReader(details_closepngbytes))

	gblrsngfs.SET("/datatables/images/sort_asc.png", bytes.NewReader(sort_ascpngbytes))
	gblrsngfs.SET("/datatables/images/sort_asc_disabled.png", bytes.NewReader(sort_asc_disabledpngbytes))
	gblrsngfs.SET("/datatables/images/sort_both.png", bytes.NewReader(sort_bothpngbytes))
	gblrsngfs.SET("/datatables/images/sort_desc.png", bytes.NewReader(sort_descpngbytes))
	gblrsngfs.SET("/datatables/images/sort_desc_disabled.png", bytes.NewReader(sort_desc_disabledpngbytes))
}

//go:embed images/sort_asc.png
var sort_ascpngbytes []byte

//go:embed images/sort_asc_disabled.png
var sort_asc_disabledpngbytes []byte

//go:embed images/sort_both.png
var sort_bothpngbytes []byte

//go:embed images/sort_desc.png
var sort_descpngbytes []byte

//go:embed images/sort_desc_disabled.png
var sort_desc_disabledpngbytes []byte

//go:embed images/details_open.png
var details_openpngbytes []byte

//go:embed images/details_close.png
var details_closepngbytes []byte

//go:embed css/jquery.Datatables.min.css
var datatablescss string

//go:embed js/jquery.Datatables.min.js
var datatablesjs string
