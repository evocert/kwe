package chnls

import (
	"io"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/resources"
)

//ActionHandler - struct
type ActionHandler struct {
	actn    *Action
	rshndlr *resources.ResourceHandler
	altr    io.Reader
}

//NewActionHandler - for Action io
func NewActionHandler(actn *Action) (actnhndl *ActionHandler) {
	if rshndl := actn.rsngpth.ResourceHandler(); rshndl != nil {
		actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl}
	} else {
		path := actn.rsngpth.Path
		if path != "" && path[0] == '/' {
			path = path[1:]
		}
		if path != "" {
			if rqstrs := actn.rqst.Resource(path); rqstrs != nil {
				if eofclsr, eofclsrok := rqstrs.(*iorw.EOFCloseSeekReader); eofclsrok && eofclsr != nil {
					actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl, altr: eofclsr}
				} else if bf, bfok := rqstrs.(*iorw.Buffer); bfok && bf != nil && bf.Size() > 0 {
					actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl, altr: bf.Reader()}
				} else if fncr, fncrok := rqstrs.(func() io.Reader); fncrok && fncr != nil {
					actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl, altr: fncr()}
				} else if rd, rdok := rqstrs.(io.Reader); rdok {
					actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl, altr: iorw.NewEOFCloseSeekReader(rd)}
				}
			}
		}
	}
	return
}

func (actnhndlr *ActionHandler) Read(p []byte) (n int, err error) {
	if actnhndlr != nil {
		if actnhndlr.rshndlr != nil {
			n, err = actnhndlr.rshndlr.Read(p)
		} else if actnhndlr.altr != nil {
			n, err = actnhndlr.altr.Read(p)
		}
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	if err == io.EOF {
		if actnhndlr != nil && actnhndlr.altr != nil {
			if clsr, clsrok := actnhndlr.altr.(io.Closer); clsrok {
				clsr.Close()
			}
			actnhndlr.altr = nil
		}
	}
	return
}

//Close - refer to io.Closer
func (actnhndlr *ActionHandler) Close() (err error) {
	if actnhndlr != nil {
		if actnhndlr.rshndlr != nil {
			actnhndlr.rshndlr.Close()
			actnhndlr.rshndlr = nil
		}
		if actnhndlr.actn != nil {
			actnhndlr.actn = nil
		}
		if actnhndlr.altr != nil {
			if clsr, clsrok := actnhndlr.altr.(io.Closer); clsrok {
				clsr.Close()
			}
			actnhndlr.altr = nil
		}
		actnhndlr = nil
	}
	return
}
