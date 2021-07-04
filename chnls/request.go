package chnls

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/requirejs"
	"github.com/evocert/kwe/scheduling"
	"github.com/evocert/kwe/web"
	"github.com/evocert/kwe/webactions"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/resources"
)

//Request -
type Request struct {
	atv              *active.Active
	actnslst         *enumeration.List
	lstexctngactng   *Action
	rqstrsngmngr     *resources.ResourcingManager
	chnl             *Channel
	rqstoffset       int64
	rqstendoffset    int64
	rqstoffsetmax    int64
	rqstmaxsize      int64
	mediarqst        bool
	initPath         string
	settings         map[string]interface{}
	args             []interface{}
	startedWriting   bool
	startedReading   bool
	mimetype         string
	httpw            http.ResponseWriter
	httpstatus       int
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
	rmtHost          string
	lclHost          string
	zpw              *gzip.Writer
	rqstr            iorw.Reader
	Interrupted      bool
	objmap           map[string]interface{}
	intrnbuffs       map[*iorw.Buffer]*iorw.Buffer
	isFirstRequest   bool
	//caching
	mphndlr *caching.MapHandler
	//dbms
	dbms      *rqstdbms
	activecns map[string]*database.Connection
	//commands
	cmnds map[int]*osprc.Command
	//scheduling
	schdl    *scheduling.Schedule
	prntrqst *Request
	//fsutils
	fsutils *fsutils.FSUtils
	//webing
	webclient *web.ClientHandle
	//mqtt
	//mqttmngr  *mqtt.MQTTManager
	//mqtttopic mqtt.Topic
	mqttmsg   mqtt.Message
	mqttevent mqtt.MqttEvent
}

//RemoteAddr return remote Address of any network request
func (rqst *Request) RemoteAddr() string {
	return rqst.rmtHost
}

//LocalAddr return local Address of any network request
func (rqst *Request) LocalAddr() string {
	return rqst.lclHost
}

//Resource - return mapped resource interface{} by path
func (rqst *Request) Resource(path string) (rs interface{}) {
	if path != "" {
		if rs = rqst.FS().CAT(path); rs == nil && (strings.HasSuffix(path, "require.js") || strings.HasSuffix(path, "require.min.js")) {
			if strings.HasSuffix(path, "require.js") {
				path = "require.js"
			} else if strings.HasSuffix(path, "require.min.js") {
				path = "require.min.js"
			}
			rqst.FS().MKDIR("require")
			if rs = rqst.FS().CAT("require/" + path); rs == nil {
				if path == "require.js" || path == "require.min.js" {
					rqst.FS().SET("require/"+path, requirejs.RequireJS())
				}
			}
			rs = rqst.FS().CAT("require/" + path)
		} else if rs == nil {
			if rs = resources.GLOBALRSNG().FS().CAT(path); rs == nil {
				if strings.HasSuffix(path, "jquery.js") {
					rs = webactions.JqueryJS()
				} else if strings.HasSuffix(path, "webactions.js") {
					rs = webactions.WebactionsJS()
				} else if strings.HasSuffix(path, "parsing.js") {
					rs = webactions.ParsingJS()
				} else if strings.HasSuffix(path, "webactions.bundle.js") {
					rs = webactions.WebactionsBundleJS()
				}
			}
		}
	}
	return
}

//ProtoMethod - http e.g request METHOD
func (rqst *Request) ProtoMethod() string {
	return rqst.prtclmethod
}

//Proto - protocol of request e.g HTTP/1.1
func (rqst *Request) Proto() string {
	return rqst.prtcl
}

//MimeType - mimetype of extension - defaultext
func (rqst *Request) MimeType(ext string, defaultext string) (mimetype string) {
	mimetype, _ = mimes.FindMimeType(ext, defaultext)
	return
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
						rqst.actnslst.Push(nil, nil, newAction(rqst, pth))
					}
				}
			}
		}
	}
}

//ResponseHeaders wrap arround current ResponseWriter.Header
func (rqst *Request) ResponseHeaders() (hdrs []string) {
	hdrs = []string{}
	if hdr := rqst.ResponseHeader(); hdr != nil {
		for k := range hdr {
			hdrs = append(hdrs, k)
		}
	}
	return
}

//ResponseHeader wrap arround current ResponseWriter.Header
func (rqst *Request) ResponseHeader() (hdr http.Header) {
	if httpw := rqst.httpw; httpw != nil {
		hdr = httpw.Header()
	}
	return
}

//SetResponseStatus  set Response Status
func (rqst *Request) SetResponseStatus(status int) {
	if httpw := rqst.httpw; httpw != nil {
		rqst.httpstatus = status
	}
}

//RequestHeaders wrap arround current Request.Header
func (rqst *Request) RequestHeaders() (hdrs []string) {
	hdrs = []string{}
	if hdr := rqst.RequestHeader(); hdr != nil {
		for k := range hdr {
			hdrs = append(hdrs, k)
		}
	}
	return
}

//RequestHeader wrap arround current Request.Header
func (rqst *Request) RequestHeader() (hdr http.Header) {
	if httpr := rqst.httpr; httpr != nil {
		hdr = httpr.Header
	}
	return hdr
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

//RequestBody - RequestBody as *iorw.EOFCloseSeekReader
func (rqst *Request) RequestBody(cached ...bool) (bf *iorw.EOFCloseSeekReader) {
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
					if httpr := rqst.httpr; httpr != nil {
						if bdy := httpr.Body; bdy != nil {
							io.Copy(io.MultiWriter(po, rqst.cchdrqstcntnt), bdy)
						}
					}
				}()
				bf = iorw.NewEOFCloseSeekReader(pi)
			} else {
				if rqst.cchdrqstcntntrdr == nil {
					rqst.cchdrqstcntntrdr = rqst.cchdrqstcntnt.Reader()
				} else {
					rqst.cchdrqstcntntrdr.Seek(0, io.SeekStart)
				}
			}
			if rqst.cchdrqstcntntrdr != nil {
				bf = iorw.NewEOFCloseSeekReader(rqst.cchdrqstcntntrdr)
			}
		} else {
			if rqst.cchdrqstcntnt != nil && rqst.cchdrqstcntntrdr != nil {
				rqst.cchdrqstcntntrdr.Seek(0, io.SeekStart)
				bf = iorw.NewEOFCloseSeekReader(rqst.cchdrqstcntntrdr)
			} else if httpr := rqst.httpr; httpr != nil {
				if bdy := httpr.Body; bdy != nil {
					bf = iorw.NewEOFCloseSeekReader(bdy)
				}
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
		if rqst.actnslst != nil {
			var actntodispose *Action
			rqst.actnslst.Dispose(
				nil,
				func(nde *enumeration.Node, val interface{}) {
					if actntodispose, _ = val.(*Action); actntodispose != nil {
						actntodispose.Close()
						actntodispose = nil
					}
				})
			actntodispose = nil
			rqst.actnslst = nil
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
		if rqst.dbms != nil {
			rqst.dbms.dbms = nil
			rqst.dbms.rqst = nil
			rqst.dbms = nil
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
		if rqst.cmnds != nil {
			if il := len(rqst.cmnds); il > 0 {
				cms := make([]int, il)
				cmsi := 0
				for cmi := range rqst.cmnds {
					cms[cmsi] = cmi
					cmsi++
				}
				for len(cms) > 0 {
					cmi := cms[0]
					rqst.cmnds[cmi].Close()
					cms = cms[1:]
				}
			}
			rqst.cmnds = nil
		}
		if rqst.rqstrsngmngr != nil {
			rqst.rqstrsngmngr.Close()
			rqst.rqstrsngmngr = nil
		}
		if rqst.lstexctngactng != nil {
			for rqst.lstexctngactng != nil {
				rqst.lstexctngactng.Close()
				if rqst.lstexctngactng != nil && rqst.lstexctngactng.rqst == nil {
					rqst.lstexctngactng = nil
				}
			}
			rqst.lstexctngactng = nil
		}
		if rqst.schdl != nil {
			rqst.schdl = nil
		}
		if rqst.prntrqst != nil {
			rqst.prntrqst = nil
		}
		if rqst.webclient != nil {
			rqst.webclient.Close()
			rqst.webclient = nil
		}
		if rqst.mphndlr != nil {
			rqst.mphndlr.Close()
			rqst.mphndlr = nil
		}
		if rqst.mqttmsg != nil {
			rqst.mqttmsg = nil
		}
		if rqst.mqttevent != nil {
			rqst.mqttevent = nil
		}
		rqst = nil
	}
	return
}

func (rqst *Request) execute(interrupt func()) {
	if httpr := rqst.httpr; httpr != nil {
		if httpw := rqst.httpw; httpw != nil {
			rqst.prtcl = httpr.Proto
			rqst.prtclmethod = httpr.Method
			func() {
				isCancelled := false
				ctx, cancel := context.WithCancel(httpr.Context())
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
		}
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
	if httpw := rqst.httpw; httpw != nil {
		if rqst.zpw != nil {
			n, err = rqst.zpw.Write(p)
		} else if httpw != nil {
			n, err = httpw.Write(p)
		} else if rqstw := rqst.rqstw; rqstw != nil {
			n, err = rqstw.Write(p)
		}
	} else if rqstw := rqst.rqstw; rqstw != nil {
		n, err = rqstw.Write(p)
	}
	return
}

func (rqst *Request) internRead(p []byte) (n int, err error) {
	if rqstr := rqst.rqstr; rqstr != nil {
		n, err = rqstr.Read(p)
	}
	return
}

func (rqst *Request) Write(p []byte) (n int, err error) {
	if rqst != nil {
		if pl := len(p); pl > 0 {
			if !rqst.startedWriting {
				if err = rqst.startWriting(); err != nil {
					return
				}
			}
			n, err = rqst.internWrite(p)
		}
	}
	return
}

func (rqst *Request) startReading() (err error) {
	defer func() {
		if rv := recover(); rv != nil {
			err = fmt.Errorf("%v", rv)
		}
		rqst.startedReading = true
	}()
	return
}

func (rqst *Request) Read(p []byte) (n int, err error) {
	if rqst != nil {
		if pl := len(p); pl > 0 {
			if !rqst.startedReading {
				if err = rqst.startReading(); err != nil {
					return
				}
			}
			n, err = rqst.internRead(p)
		}
	}
	return
}

func (rqst *Request) WebClient() (webclient *web.ClientHandle) {
	if rqst.webclient == nil {
		webclient = &web.ClientHandle{
			Client: web.NewClient(),
		}
		webclient.Send = func(rqstpath string, a ...interface{}) (rspr io.Reader, err error) {
			if len(a) == 0 {
				if a == nil {
					a = []interface{}{rqst.atv}
				}
			} else {
				a = append([]interface{}{rqst.atv}, a...)
			}
			return rqst.webclient.Client.Send(rqstpath, a...)
		}
		webclient.SendRespondString = func(rqstpath string, a ...interface{}) (rspstr string, err error) {
			if len(a) == 0 {
				if a == nil {
					a = []interface{}{rqst.atv}
				}
			} else {
				a = append([]interface{}{rqst.atv}, a...)
			}
			return rqst.webclient.Client.SendRespondString(rqstpath, a...)
		}
		webclient.SendReceive = func(rqstpath string, a ...interface{}) (rw web.ReaderWriter, err error) {
			if len(a) == 0 {
				if a == nil {
					a = []interface{}{rqst.atv}
				}
			} else {
				a = append([]interface{}{rqst.atv}, a...)
			}
			return rqst.webclient.Client.SendReceive(rqstpath, a...)
		}
		webclient.Close = func() {
			if webclient.Client != nil {
				webclient.Client.Close()
				webclient.Client = nil
			}
		}
		rqst.webclient = webclient
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
			tmpltext = filepath.Ext(actn.rspath)
		}

		if strings.HasPrefix(tmpltpath, "/") {
			tmpltpath = tmpltpath[1:]
			tmpltpathroot = "/"
		}
		if !strings.HasPrefix(tmpltpath, "/") {
			if tmpltpathroot == "" {
				tmpltpathroot = actn.rspath //.rsngpth.LookupPath
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
				if elipsePathi := strings.Index(tmpltpath, "../"); tmpltpathroot != "/" && strings.HasSuffix(tmpltpathroot, "/") && elipsePathi > -1 {
					removecnt := 1
					for {
						if strings.HasPrefix(tmpltpath[elipsePathi+len("../"):], "../") {
							elipsePathi += len("../")
							removecnt++
						} else {
							break
						}
					}
					for removecnt > 0 {
						if pthi := strings.LastIndex(tmpltpathroot[:len(tmpltpathroot)-1], "/"); pthi > -1 {
							tmpltpathroot = tmpltpathroot[:pthi+1]
							removecnt--
						} else {
							break
						}
					}
					if tmpltpathroot == "" {
						tmpltpathroot = "/"
					}
				}
			}
			if tmpltpath = tmpltpathroot + tmpltpath; /*+ tmpltext*/ tmpltpath != "" {
				if rdr = rqst.FS().CAT(tmpltpath); rdr == nil {
					rdr = resources.GLOBALRSNG().FS().CAT(tmpltpath)
				}
				tmpltpath = ""
			}
		}
	}
	return
}

func (rqst *Request) processPaths(wrapup bool) {
	if rqst.actnslst.Length() > 0 {
		var actn *Action = nil
		rqst.actnslst.Do( //RemovingNode
			func(nde *enumeration.Node, val interface{}) bool {
				if actn, _ = val.(*Action); actn != nil {
					executeAction(actn)
					actn.Close()
					nde.Set(nil)
				}
				return true
			},
			//RemovedNode
			func(nde *enumeration.Node, val interface{}) {
				if actn, _ = val.(*Action); actn != nil {
					actn.Close()
					nde.Set(nil)
				}
			},
			//DisposingNode
			func(nde *enumeration.Node, val interface{}) {
				if actn, _ = val.(*Action); actn != nil {
					actn.Close()
				}
			})
		actn = nil
	}
	//var actn *Action = nil
	/*for len(rqst.actns) > 0 && !rqst.Interrupted {
		actn = rqst.actns[0]
		rqst.actns = rqst.actns[1:]
		func() {
			defer func() {

			}()
			executeAction(actn)
		}()
	}*/
	if rqst.wbytesi > 0 {
		_, _ = rqst.internWrite(rqst.wbytes[:rqst.wbytesi])
	}
	if !rqst.startedWriting {
		if rqst.httpw != nil {
			rqst.httpw.Header().Set("Content-Length", "0")
		}
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

//Readln helper Readln() over *Request
func (rqst *Request) Readln() (ln string, err error) {
	ln, err = iorw.ReadLine(rqst)
	return
}

//ReadLines helper ReadLines() over *Request
func (rqst *Request) ReadLines() (lines []string, err error) {
	lines, err = iorw.ReadLines(rqst)
	return
}

//ReadAll helper ReadAll() over *Request
func (rqst *Request) ReadAll() (s string, err error) {
	if rqst.rqstr != nil {

	} else {
		s, err = iorw.ReaderToString(rqst)
	}
	return
}

func (rqst *Request) copy(r io.Reader, altw io.Writer, istext bool, isactive bool, initpath string) {
	if rqst != nil {
		if istext {
			rqst.invokeAtv()
			if altw == nil {
				if isactive {
					if err := rqst.atv.Eval(rqst, rqst, initpath, "<@", r, "@>"); err != nil {
						if err != io.EOF {
							fmt.Println(err)
						}
					}
				} else {
					if err := rqst.atv.Eval(rqst, rqst, initpath, r); err != nil {
						if err != io.EOF {
							fmt.Println(err)
						}
					}
				}
			} else {
				rqst.atv.Eval(altw, rqst, initpath, r)
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
			if rqst.flshr != nil {
				rqst.flshr.Flush()
			} else {
				if httpw := rqst.httpw; httpw != nil {
					if wflsh, wflshok := httpw.(http.Flusher); wflshok {
						wflsh.Flush()
					}
				}
			}
		}
	}
	return
}

func (rqst *Request) startWriting() (err error) {
	defer func() {
		if rv := recover(); rv != nil {
			err = fmt.Errorf("%v", rv)
		}
	}()
	//if rqst.httpr != nil && rqst.httpw != nil {
	if rqst.startedWriting {
		return
	}
	rqst.startedWriting = true
	if hdr := rqst.ResponseHeader(); hdr != nil {
		if hdr.Get("Content-Type") == "" {
			hdr.Set("Content-Type", rqst.mimetype)
		}
		if cntntl := hdr.Get("Content-Length"); cntntl != "" {
			if cntntl != "0" {
				hdr.Del("Content-Length")
			}
		}
		hdr.Set("Cache-Control", "no-cache")
		hdr.Set("Expires", time.Now().Format(http.TimeFormat))
		hdr.Set("Connection", "close")
	}
	//httpw.Header().Set("Transfer-Encoding", "chunked")
	//rqst.zpw = gzip.NewWriter(httpw)
	if httpw := rqst.httpw; httpw != nil {
		httpw.WriteHeader(rqst.httpstatus)
	}
	//}
	return
}

func (rqst *Request) executeHTTP(interrupt func()) {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		if httpr := rqst.httpr; httpr != nil {
			parameters.LoadParametersFromHTTPRequest(rqst.prms, httpr)
			rqst.executePath(httpr.URL.Path, interrupt)
		}
	}
}

func (rqst *Request) executePath(path string, interrupt func()) {
	if rqst != nil {
		rqst.initPath = path
		rqst.AddPath(path)
		rqst.processPaths(true)
	}
}

func (rqst *Request) ismediaExt(ext string) bool {
	ext = filepath.Ext(ext)
	return ext == ".mp4"
}

func (rqst *Request) executeRW(interrupt func()) (err error) {
	if rqst != nil {
		rqst.prms = parameters.NewParameters()
		//if rqststdio := newrequeststdio(rqst); rqststdio != nil {
		//	func() {
		//		defer rqststdio.dispose()
		//		err = rqststdio.executeStdIO()
		//	}()
		//	rqst.wrapup()
		//}
		//rqst.wrapup()
		rqst.AddPath(rqst.initPath)
		rqst.processPaths(true)
	}
	return
}

type rqstdbms struct {
	dbms *database.DBMS
	rqst *Request
}

func (rstdbms *rqstdbms) Query(a interface{}, qryargs ...interface{}) (reader *database.Reader) {
	if len(qryargs) == 0 {
		qryargs = []interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}
	} else {
		qryargs = append([]interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}, qryargs...)
	}
	reader = rstdbms.dbms.Query(a, qryargs...)
	return
}

func (rstdbms *rqstdbms) Execute(a interface{}, excargs ...interface{}) (exctr *database.Executor) {
	if len(excargs) == 0 {
		excargs = []interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}
	} else {
		excargs = append([]interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}, excargs...)
	}
	exctr = rstdbms.dbms.Execute(a, excargs...)
	return
}

func (rstdbms *rqstdbms) InOutS(in interface{}, ioargs ...interface{}) (out string, err error) {
	if len(ioargs) == 0 {
		ioargs = []interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}
	} else {
		ioargs = append([]interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}, ioargs...)
	}
	out, err = rstdbms.dbms.InOutS(in, ioargs...)
	return
}

func (rstdbms *rqstdbms) InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error) {
	if len(ioargs) == 0 {
		ioargs = []interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}
	} else {
		ioargs = append([]interface{}{rstdbms.rqst.atv, rstdbms.rqst.prms}, ioargs...)
	}
	err = rstdbms.dbms.InOut(in, out, ioargs...)
	return
}

func (rqtdbms *rqstdbms) Connection(alias string) (cn *database.Connection) {
	cn = rqtdbms.dbms.Connection(alias)
	return
}

func (rqtdbms *rqstdbms) Connections() (cns []string) {
	cns = rqtdbms.dbms.Connections()
	return
}

func (rqtdbms *rqstdbms) RegisterConnection(alias string, driver string, datasource string, a ...interface{}) (registered bool) {
	registered = rqtdbms.dbms.RegisterConnection(alias, driver, datasource, a...)
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

func (rqst *Request) ExecutePath(path string) (err error) {
	actn := newAction(rqst, path)
	defer func() {
		actn.Close()
		actn = nil
	}()
	err = executeAction(actn)
	return
}

func (rqst *Request) invokeAtv() {
	if rqst.atv == nil {
		rqst.atv = active.NewActive()
		nmspce := ""
		if rqst.atv != nil {
			nmspce = rqst.atv.Namespace
			if nmspce != "" {
				nmspce = nmspce + "."
			}
		}
		rqst.objmap[nmspce+"request"] = rqst
		rqst.objmap[nmspce+"channel"] = rqst.chnl
		rqst.objmap[nmspce+"caching"] = rqst.mphndlr
		rqst.objmap[nmspce+"caching"] = rqst.mphndlr

		rqst.objmap[nmspce+"mqtting"] = rqst.chnl.MQTT()

		if rqst.mqttmsg != nil {
			rqst.objmap[nmspce+"mqttmsg"] = rqst.mqttmsg
		}
		if rqst.mqttevent != nil {
			rqst.objmap[nmspce+"mqttevent"] = rqst.mqttmsg
		}
		rqst.dbms = &rqstdbms{rqst: rqst, dbms: database.GLOBALDBMS()}
		rqst.objmap[nmspce+"dbms"] = rqst.dbms
		rqst.objmap[nmspce+"resourcing"] = resources.GLOBALRSNG()
		rqst.objmap[nmspce+"newrqstbuffer"] = func() (buff *iorw.Buffer) {
			buff = iorw.NewBuffer()
			buff.OnClose = rqst.removeBuffer
			rqst.intrnbuffs[buff] = buff
			return
		}
		rqst.objmap[nmspce+"newcommand"] = func(execpath string, execargs ...string) (cmd *osprc.Command, err error) {
			cmd, err = osprc.NewCommand(execpath, execargs...)
			if err == nil && cmd != nil {
				cmd.OnClose = rqst.removeCommand
				rqst.cmnds[cmd.PrcID()] = cmd
			}
			return
		}
		rqst.objmap[nmspce+"action"] = func() *Action {
			return rqst.lstexctngactng
		}
		rqst.objmap[nmspce+"webing"] = rqst.WebClient()

		fstls := fsutils.NewFSUtils()
		rqst.objmap[nmspce+"_fsutils"] = fstls

		for cobjk, cobj := range rqst.chnl.objmap {
			rqst.objmap[cobjk] = cobj
		}

		if len(rqst.objmap) > 0 {
			rqst.atv.ImportGlobals(rqst.objmap)
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
}

func (rqst *Request) removeCommand(cmdprcid int) {
	if len(rqst.cmnds) > 0 {
		if cmd, cmdok := rqst.cmnds[cmdprcid]; cmdok && cmd != nil {
			rqst.cmnds[cmdprcid] = nil
			delete(rqst.cmnds, cmdprcid)
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
