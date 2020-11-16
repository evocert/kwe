package chnls

import (
	"io"
	"strings"

	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/resources"
)

//Action - struct
type Action struct {
	rqst    *Request
	rsngpth *resources.ResourcingPath
	sttngs  map[string]interface{}
}

func newAction(rqst *Request, rsngpth *resources.ResourcingPath) (actn *Action) {
	actn = &Action{rqst: rqst, rsngpth: rsngpth}
	return
}

func executeAction(actn *Action, rqstTmpltLkp func(tmpltpath string, a ...interface{}) (rdr io.Reader)) (err error) {
	var rspath = actn.rsngpth.Path
	var isTextRequest = false
	if curactnhndlr := actn.ActionHandler(); curactnhndlr == nil {
		if rspth := actn.rsngpth.Path; rspth != "" {
			if _, ok := actn.rqst.rsngpthsref[rspth]; ok {
				actn.rqst.rsngpthsref[rspth] = nil
				delete(actn.rqst.rsngpthsref, rspth)
			}
		}
		if actn.rqst.isFirstRequest {
			actn.rqst.isFirstRequest = false
			if actn.rqst.mimetype == "" {
				actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
			}
			if rspath != "" {
				if strings.LastIndex(rspath, ".") == -1 {
					if !strings.HasSuffix(rspath, "/") {
						rspath = rspath + "/"
					}
					rspath = rspath + "index.html"
					actn.rsngpth.Path = rspath
					actn.rsngpth.LookupPath = actn.rsngpth.Path
					actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
					if curactnhndlr = actn.ActionHandler(); curactnhndlr == nil {
						actn.rqst.mimetype = "text/plain"
						isTextRequest = false
					} else {
						actn.rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
						if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
							isTextRequest = false
						}
						if isTextRequest {
							isTextRequest = false
							if actn.rqst.atv == nil {
								actn.rqst.atv = active.NewActive()
							}
							if actn.rqst.atv.ObjectMapRef == nil {
								actn.rqst.atv.ObjectMapRef = func() map[string]interface{} {
									return actn.rqst.objmap
								}
							}
							if actn.rqst.atv.LookupTemplate == nil {
								actn.rqst.atv.LookupTemplate = rqstTmpltLkp
							}
							actn.rqst.copy(curactnhndlr, nil, true)
						} else {
							actn.rqst.copy(curactnhndlr, nil, false)
						}
						curactnhndlr.Close()
						curactnhndlr = nil
					}
				} else {
					actn.Close()
				}
			} else {
				actn.Close()
			}
		} else {
			actn.Close()
		}
		actn = nil
	} else if curactnhndlr != nil {
		if actn.rqst.isFirstRequest {
			if actn.rqst.mimetype == "" {
				actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
			} else {
				_, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
			}
			actn.rqst.isFirstRequest = false
		}
		actn.rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
		if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
			isTextRequest = false
		}
		if isTextRequest {
			isTextRequest = false
			if actn.rqst.atv == nil {
				actn.rqst.atv = active.NewActive()
			}
			if actn.rqst.atv.ObjectMapRef == nil {
				actn.rqst.atv.ObjectMapRef = func() map[string]interface{} {
					return actn.rqst.objmap
				}
			}
			if actn.rqst.atv.LookupTemplate == nil {
				actn.rqst.atv.LookupTemplate = rqstTmpltLkp
			}
			actn.rqst.copy(curactnhndlr, nil, true)
		} else {
			actn.rqst.copy(curactnhndlr, nil, false)
		}
		if curactnhndlr != nil {
			curactnhndlr.Close()
			curactnhndlr = nil
		}
		actn.Close()
		actn = nil
	}
	return
}

//ActionHandler - handle individual action io
func (actn *Action) ActionHandler() (actnhndl *ActionHandler) {
	actnhndl = NewActionHandler(actn)
	return
}

//Close - action
func (actn *Action) Close() (err error) {
	if actn != nil {
		if actn.rqst != nil {
			actn.rqst = nil
		}
		if actn.rsngpth != nil {
			actn.rsngpth.Close()
			actn.rsngpth = nil
		}
	}
	return
}
