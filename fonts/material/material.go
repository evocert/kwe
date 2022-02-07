package material

import (
	"context"
	"embed"
	"io"

	"github.com/evocert/kwe/resources"
)

//go:embed css/materialdesignicons.min.css
var materialcss string

//go:embed preview.html
var previewhtml string

//go:embed fonts/*
var assetFonts embed.FS

func init() {
	var readFile = func(filepath string) (rdc io.ReadCloser) {
		pi, pw := io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func() {
			defer pw.Close()
			ctxcancel()
			if f, ferr := assetFonts.Open(filepath); ferr == nil && f != nil {
				func() {
					defer f.Close()
					io.Copy(pw, f)
				}()
			}
		}()
		<-ctx.Done()
		rdc = pi
		return
	}

	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/raw:fonts/material", "")
	gblrs.FS().SET("/fonts/material/preview.html", previewhtml)
	gblrs.FS().MKDIR("/raw:/fonts/material/css", "")
	gblrs.FS().MKDIR("/raw:/fonts/material", "")
	gblrs.FS().SET("/fonts/material/css/materialdesignicons.css", materialcss)
	gblrs.FS().SET("/fonts/material/css/materialdesignicons.min.css", materialcss)
	gblrs.FS().SET("/fonts/material/fonts/materialdesignicons-webfont.eot", readFile("fonts/materialdesignicons-webfont.eot"))     // MaterialdesigniconsWebfontEOT())
	gblrs.FS().SET("/fonts/material/fonts/materialdesignicons-webfont.ttf", readFile("fonts/materialdesignicons-webfont.ttf"))     // MaterialdesigniconsWebfontTTF())
	gblrs.FS().SET("/fonts/material/fonts/materialdesignicons-webfont.woff", readFile("fonts/materialdesignicons-webfont.woff"))   // MaterialdesigniconsWebfontWOFF())
	gblrs.FS().SET("/fonts/material/fonts/materialdesignicons-webfont.woff2", readFile("fonts/materialdesignicons-webfont.woff2")) // MaterialdesigniconsWebfontWOFF2())

	gblrs.FS().SET("/fonts/material/head.html", `<link rel="stylesheet" type="text/css" href="/fonts/material/css/materialdesignicons.min.css">`)
}
