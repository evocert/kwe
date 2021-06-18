package chnls

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/iorw"
)

//ActionHandler - struct
type ActionHandler struct {
	actn        *Action
	actnrdr     io.Reader
	actnrnrdr   io.RuneReader
	hndlMaxSize int64
	raw         bool
}

//NewActionHandler - for Action io
func NewActionHandler(actn *Action) (actnhndl *ActionHandler) {
	path := actn.rspath
	israw := false
	if path != "" && strings.Contains(path, "raw:") {
		israw = true
		path = strings.Replace(path, "raw:", "", 1)
		actn.rspath = path
	}
	pathext := filepath.Ext(path)
	path = path[:len(path)-len(pathext)]
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	hndlMaxSize := int64(-1)

	var lookuprs = func(lkppath string) bool {
		if rqstrs := actn.rqst.Resource(lkppath); rqstrs != nil {
			if eofclsr, eofclsrok := rqstrs.(*iorw.EOFCloseSeekReader); eofclsrok && eofclsr != nil {
				actnhndl = &ActionHandler{actn: actn, raw: israw, actnrdr: eofclsr, actnrnrdr: eofclsr, hndlMaxSize: hndlMaxSize}
			} else if bf, bfok := rqstrs.(*iorw.Buffer); bfok && bf != nil && bf.Size() > 0 {
				hndlMaxSize = bf.Size()
				rdr := bf.Reader()
				actnhndl = &ActionHandler{actn: actn, raw: israw, actnrdr: rdr, actnrnrdr: rdr, hndlMaxSize: hndlMaxSize}
			} else if fncr, fncrok := rqstrs.(func() io.Reader); fncrok && fncr != nil {
				eofrdr := iorw.NewEOFCloseSeekReader(fncr())
				actnhndl = &ActionHandler{actn: actn, raw: israw, actnrdr: eofrdr, actnrnrdr: eofrdr}
			} else if rd, rdok := rqstrs.(io.Reader); rdok {
				eofrdr := iorw.NewEOFCloseSeekReader(rd)
				actnhndl = &ActionHandler{actn: actn, raw: israw, actnrdr: eofrdr, actnrnrdr: eofrdr, hndlMaxSize: hndlMaxSize}
			}
		}
		return actnhndl != nil
	}

	if pathext != "" && path != "" && !strings.HasSuffix(path, "/") {
		lookuprs(path + pathext)
	} else {
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		var extrange []string = nil
		if pathext == "" {
			extrange = []string{".html", ".js", ".json", ".xml", ".svg"}
		} else {
			extrange = []string{pathext}
		}
		for _, pthext := range extrange {
			if lookuprs(path + "index" + pthext) {
				actn.rspath = path + "index" + pthext
				break
			}
		}
	}
	return
}

func (actnhndlr *ActionHandler) ReadRune() (r rune, size int, err error) {
	if actnhndlr.actnrnrdr != nil {
		r, size, err = actnhndlr.actnrnrdr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

func (actnhndlr *ActionHandler) Read(p []byte) (n int, err error) {
	if actnhndlr != nil {
		n, err = actnhndlr.actnrdr.Read(p)
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	if err == io.EOF {
		if actnhndlr != nil && actnhndlr.actnrdr != nil {
			if clsr, clsrok := actnhndlr.actnrdr.(io.Closer); clsrok {
				clsr.Close()
			}
			actnhndlr.actnrdr = nil
		}
	}
	return
}

//Close - refer to io.Closer
func (actnhndlr *ActionHandler) Close() (err error) {
	if actnhndlr != nil {
		if actnhndlr.actn != nil {
			actnhndlr.actn = nil
		}
		if actnhndlr.actnrdr != nil {
			if clsr, clsrok := actnhndlr.actnrdr.(io.Closer); clsrok {
				clsr.Close()
			}
			actnhndlr.actnrdr = nil
		}
		actnhndlr = nil
	}
	return
}
