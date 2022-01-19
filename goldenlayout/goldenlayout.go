package goldenlayout

import (
	_ "embed"

	"github.com/evocert/kwe/resources"
)

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/raw:goldenlayout/css", "")
	gblrs.FS().MKDIR("/raw:goldenlayout/js", "")
	gblrs.FS().MKDIR("/raw:goldenlayout", "")
	gblrs.FS().SET("/goldenlayout/head.html", `<link rel="stylesheet" href="/goldenlayout/css/goldenlayout-base.min.css">
	<link rel="stylesheet" href="/goldenlayout/css/goldenlayout-base.min.css">
	<link rel="stylesheet" href="/goldenlayout/css/goldenlayout-dark-theme.min.css">
	<script type="application/javascript" src="/goldenlayout/js/goldenlayout.min.js"></script>`)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-base.css", goldenlayoutbasecss)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-base.min.css", goldenlayoutbasecss)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-light-theme.css", goldenlayoutlightthemecss)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-light-theme.min.css", goldenlayoutlightthemecss)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-dark-theme.css", goldenlayoutdarkthemecss)
	gblrs.FS().SET("/goldenlayout/css/goldenlayout-dark-theme.min.css", goldenlayoutdarkthemecss)
	gblrs.FS().SET("/goldenlayout/js/goldenlayout.min.js", goldenlayoutjs)
	gblrs.FS().SET("/goldenlayout/js/goldenlayout.js", goldenlayoutjs)
}

//go:embed css/goldenlayout-light-theme.min.css
var goldenlayoutlightthemecss string

//go:embed css/goldenlayout-dark-theme.min.css
var goldenlayoutdarkthemecss string

//go:embed css/goldenlayout-base.min.css
var goldenlayoutbasecss string

//go:embed js/goldenlayout.min.js
var goldenlayoutjs string
