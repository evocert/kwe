package chnls

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/requirejs"
	"github.com/evocert/kwe/web"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/resources"
)

//Request -
type Request struct {
	atv               *active.Active
	actns             []*Action
	lstexctngactng    *Action
	rsngpthsref       map[string]*resources.ResourcingPath
	embeddedResources map[string]interface{}
	//curactnhndlr      *ActionHandler
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
	intrnbuffs       map[*iorw.Buffer]*iorw.Buffer
	isFirstRequest   bool
	//dbms
	activecns map[string]*database.Connection
}

//Resource - return mapped resource interface{} by path
func (rqst *Request) Resource(path string) (rs interface{}) {
	if path != "" {
		rs, _ = rqst.embeddedResources[path]
		if rs == nil && (strings.HasSuffix(path, "require.js") || strings.HasSuffix(path, "require.min.js")) {
			if strings.HasSuffix(path, "require.js") {
				rqst.MapResource(path, requirejs.RequireJS())
			} else {
				rqst.MapResource(path, requirejs.RequireMinJS())
			}
			rs, _ = rqst.embeddedResources[path]
		}
	}
	return
}

//RemoveResource - remove inline resource - true if found and removed and false if not exists
func (rqst *Request) RemoveResource(path string) (rmvd bool) {
	if path != "" {
		if rs, rsok := rqst.embeddedResources[path]; rsok {
			rmvd = rsok
			rqst.embeddedResources[path] = nil
			delete(rqst.embeddedResources, path)
			if rs != nil {
				if bf, bfok := rs.(*iorw.Buffer); bfok && bf != nil {
					bf.Close()
					bf = nil
				}
			}
		}
	}
	return
}

//Resources list of embedded resource paths
func (rqst *Request) Resources() (rsrs []string) {
	if lrsrs := len(rqst.embeddedResources); lrsrs > 0 {
		rsrs = make([]string, lrsrs)
		rsrsi := 0
		for rsrsk := range rqst.embeddedResources {
			rsrs[rsrsi] = rsrsk
			rsrsi++
		}
	}
	return
}

//MapResource - inline resource -  can be either func() io.Reader, *iorw.Buffer
func (rqst *Request) MapResource(path string, resource interface{}) {
	if path != "" && resource != nil {
		var validResource = false
		var strng = ""
		var isReader = false
		var r io.Reader = nil
		var isBuffer = false
		var buff *iorw.Buffer = nil

		if strng, validResource = resource.(string); !validResource {
			if _, validResource = resource.(func() io.Reader); !validResource {
				if buff, validResource = resource.(*iorw.Buffer); !validResource {
					if r, validResource = resource.(io.Reader); validResource {
						validResource = (r != nil)
					}
					isReader = validResource
				} else {
					isBuffer = true
				}
			}
		} else {
			if strng != "" {
				r = strings.NewReader(strng)
				isReader = true
			} else {
				validResource = false
			}
		}
		if validResource {
			if isReader {
				buff := iorw.NewBuffer()
				io.Copy(buff, r)
				resource = buff
			}
			if _, resourceok := rqst.embeddedResources[path]; resourceok && rqst.embeddedResources[path] != resource {
				if rqst.embeddedResources[path] != nil {
					if buff, isBuffer = rqst.embeddedResources[path].(*iorw.Buffer); isBuffer {
						buff.Close()
						buff = nil
					}
					rqst.embeddedResources[path] = resource
				} else {
					rqst.embeddedResources[path] = resource
				}
			} else {
				rqst.embeddedResources[path] = resource
			}
		}
	}
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
						if strings.Index(pth, ":raw/") > -1 {
							pth = strings.Replace(pth, ":raw/", "/", -1)
							if !(strings.Index(pth, "@") > -1 && strings.Index(pth, "@") < strings.LastIndex(pth, "@")) {
								pth += "@@"
							}
						}
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
			if rqst.cchdrqstcntnt != nil && rqst.cchdrqstcntntrdr != nil {
				rqst.cchdrqstcntntrdr.Seek(0, io.SeekStart)
				bf = bufio.NewReader(rqst.cchdrqstcntntrdr)
			} else if bdy := rqst.httpr.Body; bdy != nil {
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
		if rqst.intrnbuffs != nil {
			if il := len(rqst.intrnbuffs); il > 0 {
				bfs := make([]*iorw.Buffer, il)
				bfsi := 0
				for bf := range rqst.intrnbuffs {
					bfs[bfsi] = bf
					bfsi++
				}
				for len(bfs) > 0 {
					bf := bfs[0]
					bf.Close()
					bf = nil
					bfs = bfs[1:]
				}
			}
			rqst.intrnbuffs = nil
		}
		if rqst.embeddedResources != nil {
			if emdbrsrs := rqst.Resources(); len(emdbrsrs) > 0 {
				for _, embdk := range emdbrsrs {
					rqst.RemoveResource(embdk)
				}
				emdbrsrs = nil
			}
			rqst.embeddedResources = nil
		}
		if rqst.lstexctngactng != nil {
			for rqst.lstexctngactng != nil {
				rqst.lstexctngactng.Close()
			}
			rqst.lstexctngactng = nil
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
			isCancelled := false
			ctx, cancel := context.WithCancel(rqst.httpr.Context())
			defer func() {
				if r := recover(); r != nil {

				}
				if !isCancelled {
					isCancelled = true
					cancel()
				}
			}()

			go func() {
				defer func() {
					if r := recover(); r != nil {

					}
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
	} else if rqst.rqstw != nil && rqst.rqstr != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {

				}
			}()
			isCancelled := false
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				defer func() {
					if r := recover(); r != nil {

					}
					isCancelled = true
					cancel()
				}()
				if rwerr := rqst.executeRW(interrupt); rwerr != nil {
					fmt.Println(rwerr)
				}
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
		} else if rqst.httpw != nil {
			n, err = rqst.httpw.Write(p)
			if rqst.flshr != nil && n > 0 && err == nil {
				rqst.flshr.Flush()
			}
		} else if rqst.rqstw != nil {
			n, err = rqst.rqstw.Write(p)
			if rqst.flshr != nil && n > 0 && err == nil {
				rqst.flshr.Flush()
			}
		}
	} else if rqst.rqstw != nil {
		n, err = rqst.rqstw.Write(p)
		if rqst.flshr != nil && n > 0 && err == nil {
			rqst.flshr.Flush()
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

func (rqst *Request) templateLookup(actn *Action, tmpltpath string, a ...interface{}) (rdr io.Reader, rdrerr error) {
	tmpltpath = strings.Replace(tmpltpath, "\\", "/", -1)
	if strings.HasPrefix(tmpltpath, "http://") || strings.HasPrefix(tmpltpath, "https://") {
		rdr, rdrerr = web.DefaultClient.Send(tmpltpath, nil)
	} else if actn != nil {
		if strings.Index(tmpltpath, ":raw/") > -1 {
			tmpltpath = strings.Replace(tmpltpath, ":raw/", "/", -1)
			if !(strings.Index(tmpltpath, "@") > -1 && strings.Index(tmpltpath, "@") < strings.LastIndex(tmpltpath, "@")) {
				tmpltpath += "@@"
			}
		}
		var tmpltpathroot = ""
		var tmpltext = filepath.Ext(tmpltpath)
		if tmpltext == "" {
			tmpltext = filepath.Ext(actn.rsngpth.LookupPath)
		}

		if strings.HasPrefix(tmpltpath, "/") {
			tmpltpath = tmpltpath[1:]
			tmpltpathroot = "/"
		}
		if !strings.HasPrefix(tmpltpath, "/") {
			if tmpltpathroot == "" {
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
			}
			if tmpltpath = tmpltpathroot + tmpltpath; /*+ tmpltext*/ tmpltpath != "" {
				rdr = actn.rsngpth.ResourceHandler(tmpltpath)
				tmpltpath = ""
			}
		}
	}
	return
}

func (rqst *Request) processPaths(wrapup bool) {
	var actn *Action = nil
	for len(rqst.actns) > 0 && !rqst.Interrupted {
		actn = rqst.actns[0]
		rqst.actns = rqst.actns[1:]
		func() {
			defer func() {
				/*if r := recover(); r != nil {
					if !rqst.Interrupted {
						rqst.Interrupt()
					}
					actn.Close()
				}*/
			}()
			executeAction(actn)
		}()
	}
	if rqst.wbytesi > 0 {
		_, _ = rqst.internWrite(rqst.wbytes[:rqst.wbytesi])
	}
	if !rqst.startedWriting {
		rqst.startWriting()
	}
	if wrapup {
		rqst.wrapup()
	}
}

//Print helper Print(...interface) over *Request
func (rqst *Request) Print(a ...interface{}) {
	iorw.Fprint(rqst, a...)
}

//Println helper Println(...interface) over *Request
func (rqst *Request) Println(a ...interface{}) {
	iorw.Fprintln(rqst, a...)
}

func (rqst *Request) copy(r io.Reader, altw io.Writer, istext bool) {
	if rqst != nil {
		if istext {
			rqst.invokeAtv()
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
	if rqst.httpr != nil && rqst.httpw != nil {
		if rqst.startedWriting {
			return
		}
		rqst.startedWriting = true
		if rqst.httpw.Header().Get("Content-Type") == "" {
			rqst.httpw.Header().Set("Content-Type", rqst.mimetype)
		}
		rqst.httpw.Header().Del("Content-Length")
		rqst.httpw.Header().Set("Cache-Control", "no-cache")
		rqst.httpw.Header().Set("Expires", time.Now().Format(http.TimeFormat))
		rqst.httpw.Header().Set("Connection", "close")
		//httpw.Header().Set("Transfer-Encoding", "chunked")
		//rqst.zpw = gzip.NewWriter(httpw)
		rqst.httpw.WriteHeader(200)

	}
}

func (rqst *Request) executeHTTP(interrupt func()) {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		parameters.LoadParametersFromHTTPRequest(rqst.prms, rqst.httpr)
		rqst.AddPath(rqst.httpr.URL.Path)
		rqst.processPaths(true)
	}
}

func (rqst *Request) executeRW(interrupt func()) (err error) {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		if rqststdio := newrequeststdio(rqst); rqststdio != nil {
			func() {
				defer rqststdio.dispose()
				err = rqststdio.executeStdIO()
			}()
			rqst.wrapup()
		}
	}
	return
}

func newRequest(chnl *Channel, rdr io.Reader, wtr io.Writer, a ...interface{}) (rqst *Request, interrupt func()) {
	var rqstsettings map[string]interface{} = nil
	var ai = 0
	var httpw http.ResponseWriter = nil
	var httpflshr http.Flusher = nil
	var httpr *http.Request = nil

	if wtr != nil {
		if dhttpw, dhttpwok := wtr.(http.ResponseWriter); dhttpwok {
			if httpw == nil {
				httpw = dhttpw
				if wtr == nil {
					wtr = httpw
				}
				if flshr, flshrok := httpw.(http.Flusher); flshrok {
					httpflshr = flshr
				}
			}
		}
	}

	for ai < len(a) {
		if da, daok := a[ai].([]interface{}); daok {
			if al := len(da); al > 0 {
				a = append(da, a[1:]...)
				ai = 0
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
	rqst = &Request{isFirstRequest: true, mimetype: "", zpw: nil, Interrupted: false, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, rqstw: wtr, httpw: httpw, rqstr: rdr, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, actns: []*Action{}, args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}, embeddedResources: map[string]interface{}{}, activecns: map[string]*database.Connection{}}
	rqst.invokeAtv()
	nmspce := ""
	if rqst.atv != nil {
		nmspce = rqst.atv.Namespace
		if nmspce != "" {
			nmspce = nmspce + "."
		}
	}
	rqst.objmap[nmspce+"request"] = rqst
	rqst.objmap[nmspce+"channel"] = chnl
	rqst.objmap[nmspce+"dbms"] = database.GLOBALDBMS()
	rqst.objmap[nmspce+"resourcing"] = resources.GLOBALRSNG()
	rqst.objmap[nmspce+"newrqstbuffer"] = func() (buff *iorw.Buffer) {
		buff = iorw.NewBuffer()
		buff.OnClose = rqst.removeBuffer
		rqst.intrnbuffs[buff] = buff
		return
	}
	rqst.objmap[nmspce+"action"] = func() *Action {
		return rqst.lstexctngactng
	}
	rqst.objmap[nmspce+"webing"] = web.DefaultClient
	for cobjk, cobj := range chnl.objmap {
		rqst.objmap[cobjk] = cobj
	}

	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}

	interrupt = func() {
		rqst.Interrupt()
	}
	return
}

func (rqst *Request) detachAction(actn *Action) {
	if actn.prvactn != nil {
		rqst.lstexctngactng = actn.prvactn
		actn.prvactn = nil
	} else {
		rqst.lstexctngactng = nil
	}
}

func (rqst *Request) invokeAtv() {
	if rqst.atv == nil {
		rqst.atv = active.NewActive()
	}

	if rqst.atv.ObjectMapRef == nil {
		rqst.atv.ObjectMapRef = func() map[string]interface{} {
			return rqst.objmap
		}
	}
	if rqst.atv.LookupTemplate == nil {
		rqst.atv.LookupTemplate = func(tmpltpath string, a ...interface{}) (rdr io.Reader, rdrerr error) {
			return rqst.templateLookup(rqst.lstexctngactng, tmpltpath, a...)
		}
	}
}

func (rqst *Request) removeBuffer(buff *iorw.Buffer) {
	if len(rqst.intrnbuffs) > 0 {
		if bf, bfok := rqst.intrnbuffs[buff]; bfok && bf == buff {
			rqst.intrnbuffs[buff] = nil
			delete(rqst.intrnbuffs, buff)
		}
	}
}

//Response - struct
type Response struct {
	r              *http.Request
	w              io.Writer
	statusCode     int
	header         http.Header
	canWriteHeader bool
}

//NewResponse - Instance of Response http.ResponseWriter helper
func NewResponse(w io.Writer, r *http.Request) (resp *Response) {
	resp = &Response{w: w, header: http.Header{}, r: r, canWriteHeader: true}
	return resp
}

//Header refer to http.Header
func (resp *Response) Header() http.Header {
	return resp.header
}

//Writer refer to io.Writer
func (resp *Response) Write(p []byte) (n int, err error) {
	if resp.w != nil {
		n, err = resp.w.Write(p)
	}
	return 0, nil
}

//WriteHeader - refer to http.ResponseWriter -> WriteHeader
func (resp *Response) WriteHeader(statusCode int) {
	resp.statusCode = statusCode

	if resp.w != nil {
		if resp.canWriteHeader {
			var statusLine = resp.r.Proto + " " + fmt.Sprintf("%d", statusCode) + " " + http.StatusText(statusCode)
			fmt.Fprintln(resp.w, statusLine)
			if resp.header != nil {
				resp.header.Write(resp.w)
			}
			fmt.Fprintln(resp.w)
		}
	}
}

//Flush refer to http.Flusher
func (resp *Response) Flush() {
	if resp.w != nil {
		if flshr, flshrok := resp.w.(http.Flusher); flshrok {
			flshr.Flush()
		}
	}
}
