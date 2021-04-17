package chnls

import (
	"io"

	"github.com/evocert/kwe/iorw"
)

//ActionHandler - struct
type ActionHandler struct {
	actn        *Action
	actnrdr     io.Reader
	hndlMaxSize int64
}

//NewActionHandler - for Action io
func NewActionHandler(actn *Action) (actnhndl *ActionHandler) {
	path := actn.rspath
	if path != "" && path[0] == '/' {
		path = path[1:]
	}
	if path != "" {
		hndlMaxSize := int64(-1)
		if rqstrs := actn.rqst.Resource(path); rqstrs != nil {
			if eofclsr, eofclsrok := rqstrs.(*iorw.EOFCloseSeekReader); eofclsrok && eofclsr != nil {
				actnhndl = &ActionHandler{actn: actn, actnrdr: eofclsr, hndlMaxSize: hndlMaxSize}
			} else if bf, bfok := rqstrs.(*iorw.Buffer); bfok && bf != nil && bf.Size() > 0 {
				hndlMaxSize = bf.Size()
				actnhndl = &ActionHandler{actn: actn, actnrdr: bf.Reader(), hndlMaxSize: hndlMaxSize}
			} else if fncr, fncrok := rqstrs.(func() io.Reader); fncrok && fncr != nil {
				actnhndl = &ActionHandler{actn: actn, actnrdr: iorw.NewEOFCloseSeekReader(fncr())}
			} else if rd, rdok := rqstrs.(io.Reader); rdok {
				actnhndl = &ActionHandler{actn: actn, actnrdr: iorw.NewEOFCloseSeekReader(rd), hndlMaxSize: hndlMaxSize}
			}
		}
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
