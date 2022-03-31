package foundation

import (
	"embed"
	_ "embed"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/resources"
)

//go:embed js/*.js
var foundationjsfs embed.FS

//go:embed js/plugins/*.js
var foundationpluginsjsfs embed.FS

//go:embed css/*.css
var foundationcssfs embed.FS

func init() {
	gblrsngfs := resources.GLOBALRSNG().FS()
	foundationheadhtml := iorw.NewBuffer()
	gblrsngfs.MKDIR("/foundation/js", foundationjsfs)
	if dirs, dirserr := foundationjsfs.ReadDir("js"); dirserr == nil {
		for _, dr := range dirs {
			if !dr.IsDir() {
				foundationheadhtml.Println(`<script type="application/javascript" src="/foundation/js/` + dr.Name() + `"></script>`)
			}
		}
	}
	gblrsngfs.MKDIR("/foundation/css", foundationcssfs)
	if dirs, dirserr := foundationcssfs.ReadDir("css"); dirserr == nil {
		for _, dr := range dirs {
			if !dr.IsDir() {
				if !dr.IsDir() {
					foundationheadhtml.Println(`<link rel="stylesheet" type="text/CSS" href="/foundation/css/` + dr.Name() + `">`)
				}
			}
		}
	}
	gblrsngfs.MKDIR("/foundation/js/plugins", foundationpluginsjsfs)
	if dirs, dirserr := foundationpluginsjsfs.ReadDir("js/plugins"); dirserr == nil {
		for _, dr := range dirs {
			if !dr.IsDir() {
				foundationheadhtml.Println(`<script type="application/javascript" src="/foundation/js/plugins/` + dr.Name() + `"></script>`)
			}
		}
	}
	gblrsngfs.MKDIR("/foundation", "")
	gblrsngfs.SET("/foundation/head.html", foundationheadhtml.String())
	foundationheadhtml.Close()
}
