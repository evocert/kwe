package chnls

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/chnls/parameters"
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
	prms           *parameters.Parameters
	wbytes         []byte
	wbytesi        int
	rqstw          io.Writer
	httpr          *http.Request
	zpw            *gzip.Writer
	rqstr          io.Reader
	Interrupted    bool
	wgtxt          *sync.WaitGroup
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

//ResponseHeader wrap arround current ResponseWriter.Header
func (rqst *Request) ResponseHeader() http.Header {
	return rqst.httpw.Header()
}

//RequestHeader wrap arround current Request.Header
func (rqst *Request) RequestHeader() http.Header {
	return rqst.httpr.Header
}

//Parameters - Request web Parameters
func (rqst *Request) Parameters() *parameters.Parameters {
	return rqst.prms
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
		if rqst.wgtxt != nil {
			rqst.wgtxt = nil
		}
		if rqst.prms != nil {
			rqst.prms.CleanupParameters()
			rqst.prms = nil
		}
		if rqst.httpr != nil {
			rqst.httpr = nil
		}
		if rqst.httpw != nil {
			rqst.httpw = nil
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

func (rqst *Request) internWrite(p []byte) (n int, err error) {
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
	return
}

func (rqst *Request) Write(p []byte) (n int, err error) {
	if rqst != nil {
		if pl := len(p); pl > 0 {
			if !rqst.startedWriting {
				rqst.startWriting()
			}
			n, err = rqst.internWrite(p)
			/*
				for n < pl && err == nil {
					if cl := (len(rqst.wbytes) - rqst.wbytesi); (pl - n) >= cl {
						copy(rqst.wbytes[rqst.wbytesi:rqst.wbytesi+cl], p[n:n+cl])
						n += cl
						rqst.wbytesi += cl
					} else if cl := (pl - n); (pl - n) < (len(rqst.wbytes) - rqst.wbytesi) {
						copy(rqst.wbytes[rqst.wbytesi:rqst.wbytesi+cl], p[n:n+cl])
						n += cl
						rqst.wbytesi += cl
					}
					if rqst.wbytesi == len(rqst.wbytes) {
						_, err = rqst.internWrite(rqst.wbytes[:rqst.wbytesi])
					}
				}*/
		}
	}
	return
}

func (rqst *Request) processPaths() {
	var isFirstRequest = true
	var isTextRequest = false
	for len(rqst.rsngpaths) > 0 && !rqst.Interrupted {
		var rsngpth = rqst.rsngpaths[0]
		rqst.rsngpaths = rqst.rsngpaths[1:]
		var rspath = rsngpth.Path
		isTextRequest = false
		if rqst.currshndlr = rsngpth.ResourceHandler(); rqst.currshndlr == nil {
			if _, ok := rqst.rsngpthsref[rsngpth.Path]; ok {
				rqst.rsngpthsref[rsngpth.Path] = nil
				delete(rqst.rsngpthsref, rsngpth.Path)
			}
			if isFirstRequest {
				isFirstRequest = false
				if rqst.mimetype == "" {
					rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
				}
				if rspath != "" {
					if strings.LastIndex(rspath, ".") == -1 {
						if !strings.HasSuffix(rspath, "/") {
							rspath = rspath + "/"
						}
						rspath = rspath + "index.html"
						rsngpth.Path = rspath
						rsngpth.LookupPath = rsngpth.Path
						rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
						if rqst.currshndlr = rsngpth.ResourceHandler(); rqst.currshndlr == nil {
							rqst.mimetype = "text/plain"
							isTextRequest = false
						} else {
							rqst.rsngpthsref[rsngpth.Path] = rsngpth
							if isTextRequest && rsngpth.Path != rsngpth.LookupPath {
								isTextRequest = false
							}
							if isTextRequest {
								isTextRequest = false
								if rqst.atv == nil {
									rqst.atv = active.NewActive()
								}
								rqst.copy(rqst.currshndlr, nil, true)
							} else {
								rqst.copy(rqst.currshndlr, nil, false)
							}
						}
					} else {
						rsngpth.Close()
					}
				} else {
					rsngpth.Close()
				}

			} else {
				rsngpth.Close()
			}
			rsngpth = nil
			continue
		} else if rqst.currshndlr != nil {
			if isFirstRequest {
				if rqst.mimetype == "" {
					rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
				} else {
					_, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
				}
				isFirstRequest = false
			}
			rqst.rsngpthsref[rsngpth.Path] = rsngpth
			if isTextRequest && rsngpth.Path != rsngpth.LookupPath {
				isTextRequest = false
			}
			if isTextRequest {
				isTextRequest = false
				if rqst.atv == nil {
					rqst.atv = active.NewActive()
				}
				rqst.copy(rqst.currshndlr, nil, true)
			} else {
				rqst.copy(rqst.currshndlr, nil, false)
			}
		}
	}
	if rqst.wbytesi > 0 {
		_, _ = rqst.internWrite(rqst.wbytes[:rqst.wbytesi])
	}
	if !rqst.startedWriting {
		rqst.startWriting()
	}
	rqst.wrapup()
}

func (rqst *Request) copy(r io.Reader, altw io.Writer, istext bool) {
	if rqst != nil {
		if istext {
			if altw == nil {
				rqst.atv.Eval(rqst, r)
			} else {
				rqst.atv.Eval(altw, r)
			}
		} else {
			if altw == nil {
				io.Copy(rqst, r)
			} else {
				io.Copy(altw, r)
			}
		}
	}
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
		httpw.Header().Set("Cache-Control", "no-cache")
		httpw.Header().Set("Expires", time.Now().Format(http.TimeFormat))
		httpw.Header().Set("Connection", "close")
		//httpw.Header().Set("Transfer-Encoding", "chunked")
		//rqst.zpw = gzip.NewWriter(httpw)
		httpw.WriteHeader(200)

	}
}

func (rqst *Request) executeHTTP() {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		parameters.LoadParametersFromHTTPRequest(rqst.prms, rqst.httpr)
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
	rqst = &Request{mimetype: "", zpw: nil, atv: active.NewActive(), Interrupted: false, currshndlr: nil, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, httpw: httpw, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, rsngpaths: []*resources.ResourcingPath{}, args: make([]interface{}, len(a))}
	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}
	return
}
