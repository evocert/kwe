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
	gblrs.FS().MKDIR("/materialdesign", "")
	gblrs.FS().SET("/materialdesign/preview.html", previewhtml)
	gblrs.FS().MKDIR("/materialdesign/css", "")
	gblrs.FS().MKDIR("/materialdesign/fonts", "")
	gblrs.FS().SET("/materialdesign/css/materialdesignicons.css", materialcss)
	gblrs.FS().SET("/materialdesign/css/materialdesignicons.min.css", materialcss)
	gblrs.FS().SET("/materialdesign/fonts/materialdesignicons-webfont.eot", readFile("fonts/materialdesignicons-webfont.eot"))     // MaterialdesigniconsWebfontEOT())
	gblrs.FS().SET("/materialdesign/fonts/materialdesignicons-webfont.ttf", readFile("fonts/materialdesignicons-webfont.ttf"))     // MaterialdesigniconsWebfontTTF())
	gblrs.FS().SET("/materialdesign/fonts/materialdesignicons-webfont.woff", readFile("fonts/materialdesignicons-webfont.woff"))   // MaterialdesigniconsWebfontWOFF())
	gblrs.FS().SET("/materialdesign/fonts/materialdesignicons-webfont.woff2", readFile("fonts/materialdesignicons-webfont.woff2")) // MaterialdesigniconsWebfontWOFF2())

	gblrs.FS().MKDIR("/materialdesign/html", "")
	gblrs.FS().SET("/materialdesign/html/head.html", `<link rel="stylesheet" type="text/css" href="/materialdesign/css/materialdesignicons.min.css">`)
}
