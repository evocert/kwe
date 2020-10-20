package chnls

import (
	"io"
	"net/http"
	"strings"

	active "github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/resources"
)

//Request -
type Request struct {
	*active.Active
	rsngpaths      []*resources.ResourcingPath
	rsngpthsref    map[string]*resources.ResourcingPath
	currshndlr     *resources.ResourceHandler
	chnl           *Channel
	settings       map[string]interface{}
	args           []interface{}
	startedWriting bool
	httpw          http.ResponseWriter
	rqstw          io.Writer
	httpr          *http.Request
	rqstr          io.Reader
	Interrupted    bool
}

//AddPath - next resource path(s) to process
func (rqst *Request) AddPath(path ...string) {
	if rqst != nil {
		if len(path) > 0 {
			for len(path) > 0 {
				var pth = path[0]
				path = path[1:]
				if pth != "" {
					if strings.Index(pth, "|") > -1 {
						path = append(strings.Split(pth, "|"), path...)
						continue
					} else {
						if rsngpth, rsngpthok := rqst.rsngpthsref[pth]; rsngpthok {
							rqst.rsngpaths = append(rqst.rsngpaths, rsngpth)
						} else if rsngpth := resources.NewResourcingPath(pth, nil); rsngpth != nil {
							rqst.rsngpaths = append(rqst.rsngpaths, rsngpth)
							rqst.rsngpthsref[pth] = rsngpth
						}
					}
				}
			}
		}
	}
}

//Close - refer io.Closer
func (rqst *Request) Close() (err error) {
	if rqst != nil {
		if rqst.Active != nil {
			rqst.Active.Close()
			rqst.Active = nil
		}
		if rqst.chnl != nil {
			rqst.chnl = nil
		}
		if rqst.settings != nil {
			if len(rqst.settings) > 0 {
				var stngsks = make([]string, len(rqst.settings))
				var ski = 0
				for sk := range rqst.settings {
					stngsks[ski] = sk
					ski++
				}
				for _, sk := range stngsks {
					rqst.settings[sk] = nil
					delete(rqst.settings, sk)
				}
				stngsks = nil
			}
			rqst.settings = nil
		}
		if rqst.args != nil {
			for len(rqst.args) > 0 {
				rqst.args = rqst.args[1:]
			}
			rqst.args = nil
		}
		if rqst.rqstw != nil {
			rqst.rqstw = nil
		}
		if rqst.rqstr != nil {
			rqst.rqstr = nil
		}
		if rqst.rsngpaths != nil {
			for len(rqst.rsngpaths) > 0 {
				rqst.rsngpaths[0].Close()
				rqst.rsngpaths[0] = nil
				rqst.rsngpaths = rqst.rsngpaths[1:]
			}
			rqst.rsngpaths = nil
		}
		if rqst.rsngpthsref != nil {
			if len(rqst.rsngpthsref) > 0 {
				var rsngpaths = make([]string, len(rqst.rsngpthsref))
				var rsngpathsi = 0
				for rsngpathk := range rqst.rsngpthsref {
					rsngpaths[rsngpathsi] = rsngpathk
					rsngpathsi++
				}
				for _, rsngpathk := range rsngpaths {
					rqst.rsngpthsref[rsngpathk].Close()
					rqst.rsngpthsref[rsngpathk] = nil
					delete(rqst.rsngpthsref, rsngpathk)
				}
				rsngpaths = nil
			}
			rqst.rsngpthsref = nil
		}
		rqst = nil
	}
	return
}

func (rqst *Request) execute() {
	if rqst.httpr != nil && rqst.httpw != nil {
		rqst.executeHTTP()
	}
}

func (rqst *Request) processPaths() {
	for len(rqst.rsngpaths) > 0 && !rqst.Interrupted {
		var rsngpth = rqst.rsngpaths[0]
		rqst.rsngpaths = rqst.rsngpaths[1:]
		if rqst.currshndlr = rsngpth.ResourceHandler(); rqst.currshndlr == nil {
			if _, ok := rqst.rsngpthsref[rsngpth.Path]; ok {
				rqst.rsngpthsref[rsngpth.Path] = nil
				delete(rqst.rsngpthsref, rsngpth.Path)
			}
			rsngpth.Close()
			rsngpth = nil
			continue
		}
	}
}

func (rqst *Request) startWriting() {
	if rqst.httpr != nil && rqst.httpw != nil {
		if rqst.startedWriting {
			return
		}
		rqst.startedWriting = true
	}
}

func (rqst *Request) executeHTTP() {
	if rqst != nil {
		rqst.AddPath(rqst.httpr.URL.Path)
		rqst.processPaths()
	}
}

func newRequest(chnl *Channel, a ...interface{}) (rqst *Request) {
	var rqstsettings map[string]interface{} = nil
	var ai = 0
	var httpw http.ResponseWriter = nil
	var httpr *http.Request = nil

	var rdr io.Reader = nil
	var wtr io.Writer = nil

	for ai < len(a) {
		if da, daok := a[ai].([]interface{}); daok {
			if len(da) > 0 {
				a = append(da, a[1:])
			} else {
				a = a[1:]
			}
			continue
		} else if rstngs, rstngsok := a[ai].(map[string]interface{}); rstngsok {
			if rqstsettings == nil {
				rqstsettings = rstngs
			}
			a = a[1:]
			continue
		} else if dhttpr, dhttprok := a[ai].(*http.Request); dhttprok {
			if httpr == nil {
				httpr = dhttpr
				if rdr == nil {
					rdr = httpr.Body
				}
			}
			a = a[1:]
			continue
		} else if dhttpw, dhttpwok := a[ai].(http.ResponseWriter); dhttpwok {
			if httpw == nil {
				httpw = dhttpw
				if wtr == nil {
					wtr = httpw
				}
			}
			a = a[1:]
			continue
		} else if dr, drok := a[ai].(io.Reader); drok {
			if rdr == nil {
				rdr = dr
			}
			a = a[1:]
			continue
		} else if dw, dwok := a[ai].(io.Writer); dwok {
			if wtr == nil {
				wtr = dw
			}
			a = a[1:]
			continue
		}
		ai++
	}
	if rqstsettings == nil {
		rqstsettings = map[string]interface{}{}
	}
	rqst = &Request{Active: active.NewActive(), Interrupted: false, currshndlr: nil, startedWriting: false, httpw: httpw, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, rsngpaths: []*resources.ResourcingPath{}, args: make([]interface{}, len(a))}
	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}
	return
}
