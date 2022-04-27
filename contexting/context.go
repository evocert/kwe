package contexting

import "github.com/evocert/kwe/iorw"

type Context struct {
	rdr    iorw.Reader
	wtr    iorw.Printer
	ctxrdr *ContextReader
	ctxwtr *ContextWriter
}

type ContextReader struct {
	ctx *Context
}

func (ctxrdr *ContextReader) Seek(offset int64, whence int) (n int64, err error) {

	return
}

func (ctxrdr *ContextReader) SetMaxRead(maxlen int64) (err error) {

	return
}

func (ctxrdr *ContextReader) Read(p []byte) (n int, err error) {

	return
}

func (ctxrdr *ContextReader) ReadRune() (r rune, size int, err error) {

	return
}

func (ctxrdr *ContextReader) Readln() (line string, err error) {

	return
}

func (ctxrdr *ContextReader) Readlines() (lines []string, err error) {

	return
}

func (ctxrdr *ContextReader) ReadAll() (string, error)

type ContextWriter struct {
	ctx *Context
}
