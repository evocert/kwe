package chnls

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	active "github.com/evocert/kwe/iorw/active"
	mimes "github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/resources"
)

//Request -
type Request struct {
	atv            *active.Active
	rsngpaths      []*resources.ResourcingPath
	rsngpthsref    map[string]*resources.ResourcingPath
	currshndlr     *resources.ResourceHandler
	chnl           *Channel
	settings       map[string]interface{}
	args           []interface{}
	startedWriting bool
	mimetype       string
	httpw          http.ResponseWriter
	flshr          http.Flusher
	wbytes         []byte
	wbytesi        int
	rqstw          io.Writer
	httpr          *http.Request
	zpw            *gzip.Writer
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
		if rqst.atv != nil {
			rqst.atv.Close()
			rqst.atv = nil
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

func (rqst *Request) Write(p []byte) (n int, err error) {
	if rqst != nil {
		if !rqst.startedWriting {
			rqst.startWriting()
		}
		if rqst.httpw != nil {
			if rqst.zpw != nil {
				n, err = rqst.zpw.Write(p)
			} else {
				n, err = rqst.httpw.Write(p)
				if rqst.flshr != nil && n > 0 && err == nil {
					rqst.flshr.Flush()
				}
			}
		}
	}
	return
}

func (rqst *Request) processPaths() {
	for len(rqst.rsngpaths) > 0 && !rqst.Interrupted {
		var rsngpth = rqst.rsngpaths[0]
		if rqst.mimetype == "" {
			rqst.mimetype = mimes.FindMimeType(rsngpth.Path, "application/*")
		}
		rqst.rsngpaths = rqst.rsngpaths[1:]
		if rqst.currshndlr = rsngpth.ResourceHandler(); rqst.currshndlr == nil {
			if _, ok := rqst.rsngpthsref[rsngpth.Path]; ok {
				rqst.rsngpthsref[rsngpth.Path] = nil
				delete(rqst.rsngpthsref, rsngpth.Path)
			}
			rsngpth.Close()
			rsngpth = nil
			continue
		} else if rqst.currshndlr != nil {
			io.Copy(rqst, rqst.currshndlr)
		}
	}
	if !rqst.startedWriting {
		rqst.startWriting()
	}
	rqst.wrapup()
}

func (rqst *Request) wrapup() (err error) {
	if rqst != nil && rqst.httpw != nil {
		if rqst.zpw != nil {
			err = rqst.zpw.Close()
			rqst.zpw = nil
		}
		if err == nil {
			if wflsh, wflshok := rqst.httpw.(http.Flusher); wflshok {
				wflsh.Flush()
			}
		}
	}
	return
}

func (rqst *Request) startWriting() {
	if httpw := rqst.httpw; rqst.httpr != nil && httpw != nil {
		if rqst.startedWriting {
			return
		}
		rqst.startedWriting = true

		httpw.Header().Set("Content-Type", rqst.mimetype)
		httpw.Header().Del("Content-Length")
		//httpw.Header().Set("Transfer-Encoding", "chunked")
		//rqst.zpw = gzip.NewWriter(httpw)
		httpw.WriteHeader(200)

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
	var httpflshr http.Flusher = nil
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
				if flshr, flshrok := httpw.(http.Flusher); flshrok {
					httpflshr = flshr
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
	rqst = &Request{mimetype: "text/javascript", zpw: nil, atv: active.NewActive(), Interrupted: false, currshndlr: nil, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, httpw: httpw, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, rsngpaths: []*resources.ResourcingPath{}, args: make([]interface{}, len(a))}
	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}
	return
}
