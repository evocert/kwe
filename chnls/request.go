package chnls

import (
	"bufio"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/resources"
)

//Request -
type Request struct {
	atv              *active.Active
	actns            []*Action
	rsngpthsref      map[string]*resources.ResourcingPath
	curactnhndlr     *ActionHandler
	chnl             *Channel
	settings         map[string]interface{}
	args             []interface{}
	startedWriting   bool
	mimetype         string
	httpw            http.ResponseWriter
	flshr            http.Flusher
	prms             *parameters.Parameters
	wbytes           []byte
	wbytesi          int
	rqstw            io.Writer
	httpr            *http.Request
	cchdrqstcntnt    *iorw.Buffer
	cchdrqstcntntrdr *iorw.BuffReader
	prtclmethod      string
	prtcl            string
	zpw              *gzip.Writer
	rqstr            io.Reader
	Interrupted      bool
	wgtxt            *sync.WaitGroup
	objmap           map[string]interface{}
	isFirstRequest   bool
	//dbms
	activecns map[string]*database.Connection
}

//ProtoMethod - http e.g request METHOD
func (rqst *Request) ProtoMethod() string {
	return rqst.prtclmethod
}

//Proto - protocol of request e.g HTTP/1.1
func (rqst *Request) Proto() string {
	return rqst.prtcl
}

//Interrupt - Request execution
func (rqst *Request) Interrupt() {
	if rqst.atv != nil {
		rqst.atv.Interrupt()
	}
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
							rqst.actns = append(rqst.actns, newAction(rqst, rsngpth))
						} else if rsngpth := resources.NewResourcingPath(pth, nil); rsngpth != nil {
							rqst.actns = append(rqst.actns, newAction(rqst, rsngpth))
							rqst.rsngpthsref[pth] = rsngpth
						}
					}
				}
			}
		}
	}
}

//ResponseHeaders wrap arround current ResponseWriter.Header
func (rqst *Request) ResponseHeaders() (hdrs []string) {
	hdrs = []string{}
	for k := range rqst.httpw.Header() {
		hdrs = append(hdrs, k)
	}
	return
}

//ResponseHeader wrap arround current ResponseWriter.Header
func (rqst *Request) ResponseHeader() http.Header {
	return rqst.httpw.Header()
}

//RequestHeaders wrap arround current Request.Header
func (rqst *Request) RequestHeaders() (hdrs []string) {
	hdrs = []string{}
	for k := range rqst.httpr.Header {
		hdrs = append(hdrs, k)
	}
	return
}

//RequestHeader wrap arround current Request.Header
func (rqst *Request) RequestHeader() http.Header {
	return rqst.httpr.Header
}

//Parameters - Request web Parameters
func (rqst *Request) Parameters() *parameters.Parameters {
	return rqst.prms
}

//RequestBodyS - wrap request.RequestBody() as string
func (rqst *Request) RequestBodyS(cached ...bool) (s string) {
	if bf := rqst.RequestBody(cached...); bf != nil {
		var rns = make([]rune, 1024)
		var rnsi = 0
		for {
			r, size, rerr := bf.ReadRune()
			if size > 0 {
				rns[rnsi] = r
				rnsi++
				if rnsi == len(rns) {
					rnsi = 0
					s += string(rns)
				}
			}
			if rerr != nil {
				break
			}
		}
		if rnsi > 0 {
			s += string(rns[:rnsi])
		}
	}
	return s
}

//RequestBody - RequestBody as bufio.Reader
func (rqst *Request) RequestBody(cached ...bool) (bf *bufio.Reader) {
	if rqst.httpr != nil {
		if len(cached) == 1 && cached[0] {
			if rqst.cchdrqstcntnt == nil {
				if rqst.cchdrqstcntntrdr != nil {
					rqst.cchdrqstcntntrdr.Close()
					rqst.cchdrqstcntntrdr = nil
				}
				rqst.cchdrqstcntnt = iorw.NewBuffer()
				pi, po := io.Pipe()
				go func() {
					defer po.Close()
					if bdy := rqst.httpr.Body; bdy != nil {
						io.Copy(io.MultiWriter(po, rqst.cchdrqstcntnt), bdy)
					}
				}()
				bf = bufio.NewReader(pi)
			} else {
				if rqst.cchdrqstcntntrdr == nil {
					rqst.cchdrqstcntntrdr = rqst.cchdrqstcntnt.Reader()
				} else {
					rqst.cchdrqstcntntrdr.Seek(0, io.SeekStart)
				}
			}
			if rqst.cchdrqstcntntrdr != nil {
				bf = bufio.NewReader(rqst.cchdrqstcntntrdr)
			}
		} else {
			if bdy := rqst.httpr.Body; bdy != nil {
				bf = bufio.NewReader(bdy)
			}
		}
	}
	return
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
		if rqst.actns != nil {
			for len(rqst.actns) > 0 {
				rqst.actns[0].Close()
				rqst.actns[0] = nil
				rqst.actns = rqst.actns[1:]
			}
			rqst.actns = nil
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
		if rqst.cchdrqstcntntrdr != nil {
			rqst.cchdrqstcntntrdr.Close()
			rqst.cchdrqstcntntrdr = nil
		}
		if rqst.cchdrqstcntnt != nil {
			rqst.cchdrqstcntnt.Close()
			rqst.cchdrqstcntnt = nil
		}
		if rqst.httpr != nil {
			rqst.httpr = nil
		}
		if rqst.httpw != nil {
			rqst.httpw = nil
		}
		if rqst.objmap != nil {
			if l := len(rqst.objmap); l > 0 {
				var ks = make([]string, l)
				var ksi = 0
				for k := range rqst.objmap {
					ks[ksi] = k
					ksi++
				}
				for _, k := range ks {
					rqst.objmap[k] = nil
					delete(rqst.objmap, k)
				}
				ks = nil
			}
			rqst.objmap = nil
		}
		if rqst.activecns != nil {
			if l := len(rqst.activecns); l > 0 {
				var ks = make([]string, l)
				var ksi = 0
				for k := range rqst.activecns {
					ks[ksi] = k
					ksi++
				}
				for _, k := range ks {
					rqst.activecns[k] = nil
					delete(rqst.activecns, k)
				}
				ks = nil
			}
			rqst.activecns = nil
		}
		rqst = nil
	}
	return
}

func (rqst *Request) execute(interrupt func()) {
	if rqst.httpr != nil && rqst.httpw != nil {
		rqst.prtcl = rqst.httpr.Proto
		rqst.prtclmethod = rqst.httpr.Method
		func() {
			defer func() {
				if r := recover(); r != nil {

				}
			}()
			isCancelled := false
			ctx, cancel := context.WithCancel(rqst.httpr.Context())
			go func() {
				defer func() {
					isCancelled = true
					cancel()
				}()
				rqst.executeHTTP(interrupt)
			}()
			select {
			case <-ctx.Done():
				if ctxerr := ctx.Err(); ctxerr != nil {
					if !isCancelled {
						if interrupt != nil {
							interrupt()
						}
					}
				}
			}
		}()
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
		}
	}
	return
}

func processDbmsPath(rqst *Request, rsngpth *resources.ResourcingPath, path string) (didprocess bool, err error) {
	if strings.Index(path, "/dbms-") > -1 {
		var alias = path[strings.Index(path, "/dbms-")+1:]
		if strings.Index(alias, "/") > 0 {
			sqlpath := path[strings.Index(path, "/dbms-")+len(alias)+1:]
			alias = alias[len("dbms-"):strings.Index(alias, "/")]

			dbcn, exists := rqst.activecns[alias]
			if !exists {
				if exists, dbcn = database.GLOBALDBMS().AliasExists(alias); exists {
					rqst.activecns[alias] = dbcn
				}
			}

			if exists {
				if sqlpath != "" {
					path = sqlpath
					if !strings.HasPrefix(path, "/") {
						path = path + "/"
					}
				} else {

				}
			} else if alias == "all" {

			}
		}
	}
	return
}

func (rqst *Request) processPaths() {
	//var isFirstRequest = true
	//var isTextRequest = false
	var actn *Action = nil
	var rqstTmpltLkp = func(tmpltpath string, a ...interface{}) (rdr io.Reader) {
		if actn != nil {
			var tmpltpathroot = ""
			var tmpltext = filepath.Ext(tmpltpath)
			if tmpltext == "" {
				tmpltext = filepath.Ext(actn.rsngpth.LookupPath)
			}
			tmpltpath = strings.Replace(tmpltpath, "\\", "/", -1)
			if !strings.HasPrefix(tmpltpath, "/") {
				tmpltpathroot = actn.rsngpth.LookupPath
				if strings.LastIndex(tmpltpathroot, ".") > strings.LastIndex(tmpltpathroot, "/") {
					if strings.LastIndex(tmpltpathroot, "/") > -1 {
						tmpltpathroot = tmpltpathroot[:strings.LastIndex(tmpltpathroot, "/")+1]
						if tmpltpathroot != "/" && !strings.HasPrefix(tmpltpathroot, "/") {
							tmpltpathroot = "/" + tmpltpathroot
						}
					} else {
						tmpltpathroot = "/"
					}
				}
				if tmpltpath = tmpltpathroot + tmpltpath + tmpltext; tmpltpath != "" {
					rdr = actn.rsngpth.ResourceHandler(tmpltpath)
					tmpltpath = ""
				}
			}
		}
		return
	}
	for len(rqst.actns) > 0 && !rqst.Interrupted {
		actn = rqst.actns[0]
		rqst.actns = rqst.actns[1:]
		executeAction(actn, rqstTmpltLkp)
		/*
			var rspath = actn.rsngpth.Path
			isTextRequest = false
			if rqst.curactnhndlr = actn.ActionHandler(); rqst.curactnhndlr == nil {
				if rspth := actn.rsngpth.Path; rspth != "" {
					if _, ok := rqst.rsngpthsref[rspth]; ok {
						rqst.rsngpthsref[rspth] = nil
						delete(rqst.rsngpthsref, rspth)
					}
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
							actn.rsngpth.Path = rspath
							actn.rsngpth.LookupPath = actn.rsngpth.Path
							rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
							if rqst.curactnhndlr = actn.ActionHandler(); rqst.curactnhndlr == nil {
								rqst.mimetype = "text/plain"
								isTextRequest = false
							} else {
								rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
								if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
									isTextRequest = false
								}
								if isTextRequest {
									isTextRequest = false
									if rqst.atv == nil {
										rqst.atv = active.NewActive()
									}
									if rqst.atv.ObjectMapRef == nil {
										rqst.atv.ObjectMapRef = func() map[string]interface{} {
											return rqst.objmap
										}
									}
									if rqst.atv.LookupTemplate == nil {
										rqst.atv.LookupTemplate = rqstTmpltLkp
									}
									rqst.copy(rqst.curactnhndlr, nil, true)
								} else {
									rqst.copy(rqst.curactnhndlr, nil, false)
								}
								rqst.curactnhndlr.Close()
								rqst.curactnhndlr = nil
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
				continue
			} else if rqst.curactnhndlr != nil {
				if isFirstRequest {
					if rqst.mimetype == "" {
						rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
					} else {
						_, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
					}
					isFirstRequest = false
				}
				rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
				if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
					isTextRequest = false
				}
				if isTextRequest {
					isTextRequest = false
					if rqst.atv == nil {
						rqst.atv = active.NewActive()
					}
					if rqst.atv.ObjectMapRef == nil {
						rqst.atv.ObjectMapRef = func() map[string]interface{} {
							return rqst.objmap
						}
					}
					if rqst.atv.LookupTemplate == nil {
						rqst.atv.LookupTemplate = rqstTmpltLkp
					}
					rqst.copy(rqst.curactnhndlr, nil, true)
				} else {
					rqst.copy(rqst.curactnhndlr, nil, false)
				}
				if rqst.curactnhndlr != nil {
					rqst.curactnhndlr.Close()
					rqst.curactnhndlr = nil
				}
			}*/
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

func (rqst *Request) executeHTTP(interrupt func()) {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		parameters.LoadParametersFromHTTPRequest(rqst.prms, rqst.httpr)
		rqst.AddPath(rqst.httpr.URL.Path)
		rqst.processPaths()
	}
}

func newRequest(chnl *Channel, a ...interface{}) (rqst *Request, interrupt func()) {
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
	rqst = &Request{isFirstRequest: true, mimetype: "", zpw: nil, atv: active.NewActive(), Interrupted: false, curactnhndlr: nil, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, httpw: httpw, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, actns: []*Action{}, args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, activecns: map[string]*database.Connection{}}
	rqst.objmap["request"] = rqst
	rqst.objmap["channel"] = chnl
	rqst.objmap["dbms"] = database.GLOBALDBMS()

	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}

	interrupt = func() {
		rqst.Interrupt()
	}
	return
}
