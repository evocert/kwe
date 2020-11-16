package chnls

import (
	"io"

	"github.com/evocert/kwe/resources"
)

//ActionHandler - struct
type ActionHandler struct {
	actn    *Action
	rshndlr *resources.ResourceHandler
}

//NewActionHandler - for Action io
func NewActionHandler(actn *Action) (actnhndl *ActionHandler) {
	if rshndl := actn.rsngpth.ResourceHandler(); rshndl != nil {
		actnhndl = &ActionHandler{actn: actn, rshndlr: rshndl}
	}
	return
}

func (actnhndlr *ActionHandler) Read(p []byte) (n int, err error) {
	if actnhndlr != nil {
		if actnhndlr.rshndlr != nil {
			n, err = actnhndlr.rshndlr.Read(p)
		}
	}
	if n == 0 && err == nil {
		err = io.EOF
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
		actnhndlr = nil
	}
	return
}
