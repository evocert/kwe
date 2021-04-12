package chnls

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/scheduling"
	"github.com/evocert/kwe/ws"
	"github.com/gorilla/websocket"
)

/*Channel -
 */
type Channel struct {
	rqsts     map[*Request]*Request
	objmap    map[string]interface{}
	lstnr     *listen.Listener
	schdls    *scheduling.Schedules
	chnlrqsts chan func()
}

//Listener - *listen.Listener listener for Channel
func (chnl *Channel) Listener() *listen.Listener {
	if chnl.lstnr == nil {
		chnl.lstnr = listen.NewListener(chnl)
	}
	return chnl.lstnr
}

//Schedules - *scheduling.Schedules schedules for Channel
func (chnl *Channel) Schedules() *scheduling.Schedules {
	if chnl.schdls == nil {
		chnl.schdls = scheduling.NewSchedules(chnl)
	}
	return chnl.schdls
}

//NewSchedule - implement scheduling.ScheduleHandler NewScheduler()
func (chnl *Channel) NewSchedule(schdl *scheduling.Schedule, a ...interface{}) (scdhlhndlr scheduling.ScheduleHandler) {
	if al := len(a); al > 0 {
		ai := 0
		atvprntmap := map[string]interface{}{}
		for ai < al {
			d := a[ai]
			if rqst, rqstok := d.(*Request); rqstok {
				if rqst.atv != nil {
					rqst.atv.ExtractGlobals(atvprntmap)
				}
				//a = append(a[:ai], a[ai+1:])
				//al--
				ai++
			} else {
				ai++
			}
		}

		if scdhlrqst, _ := newRequest(chnl, nil, nil, a...); scdhlrqst != nil {
			scdhlrqst.schdl = schdl
			lclglbs := map[string]interface{}{}
			scdhlrqst.atv.ExtractGlobals(lclglbs)
			if len(atvprntmap) > 0 {
				for k := range atvprntmap {
					if len(atvprntmap) > 0 {
						if _, katvok := scdhlrqst.objmap[k]; katvok {
							atvprntmap[k] = nil
							delete(atvprntmap, k)
						} else if _, klclok := lclglbs[k]; klclok {
							atvprntmap[k] = nil
							delete(atvprntmap, k)
						}
					}
				}
			}

			scdhlrqst.atv.ImportGlobals(atvprntmap)

			scdhlhndlr = scdhlrqst
		}
	}
	return
}

//NewChannel - instance
func NewChannel() (chnl *Channel) {
	nrcpu := runtime.NumCPU()
	chnl = &Channel{rqsts: map[*Request]*Request{}, objmap: map[string]interface{}{}, chnlrqsts: make(chan func(), nrcpu*nrcpu)}
	runtime.GOMAXPROCS(nrcpu * nrcpu)
	go func() {
		for fnc := range chnl.chnlrqsts {
			go fnc()
		}
	}()
	return
}

func (chnl *Channel) nextRequest() (rqst *Request, interrupt func()) {
	rqst, interrupt = newRequest(chnl, nil, nil)
	return
}

func (chnl *Channel) internalServeHTTP(w http.ResponseWriter, r *http.Request, a ...interface{}) {
	if conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024); err == nil {
		chnl.ServeWS(conn, a...)
	} else {
		a = append([]interface{}{r}, a...)
		chnl.ServeRW(r.Body, w, a...)
	}
}

//ServeHTTP - refer http.Handler
func (chnl *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chnl.internalServeHTTP(w, r)
}

//ServeWS - server websocket Connection
func (chnl *Channel) ServeWS(wscon *websocket.Conn, a ...interface{}) {
	func() {
		if wsrw := ws.NewReaderWriter(wscon); wsrw != nil {
			defer wsrw.Close()

			/*var tmpbuf = iorw.NewBuffer()
			for {
				if _, rerr := io.Copy(tmpbuf, wsrw); rerr != nil && rerr != io.EOF {
					break
				}
				if tmpbuf.Size() > 0 {
					fmt.Println(tmpbuf)
				}
				io.Copy(wsrw, strings.NewReader(tmpbuf.String()))
				wsrw.Print(tmpbuf)
				tmpbuf.Clear()
			}
			tmpbuf.Close()
			tmpbuf = nil*/
			chnl.ServeRW(wsrw, wsrw, a...)
		}
	}()
}

//ServeRW - serve Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	if rqst, interrupt := newRequest(chnl, r, w, a...); rqst != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
				rqst.Close()
			}()
			rqst.execute(interrupt)
		}()
		rqst = nil
	}
}

//Stdio - os.Stdout, os.Stdin
func (chnl *Channel) Stdio(out *os.File, in *os.File, err *os.File, a ...interface{}) {
	chnl.ServeRW(in, out, a...)
}

//DefaultServeHTTP - helper to perform dummy ServeHttp request on channel
func (chnl *Channel) DefaultServeHTTP(w io.Writer, method string, url string, body io.Reader, a ...interface{}) {
	if rhttp, rhttperr := http.NewRequest(method, url, body); rhttperr == nil {
		if rhttp != nil {
			var whttp = NewResponse(w, rhttp)
			chnl.internalServeHTTP(whttp, rhttp, a...)
		}
	}
}

//DefaultServeRW - helper to perform dummy ServeRW request on channel
func (chnl *Channel) DefaultServeRW(w io.Writer, url string, r io.Reader, a ...interface{}) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	chnl.chnlrqsts <- func() {
		defer wg.Done()
		var method = "GET"
		if r != nil {
			method = "POST"
		}
		if rhttp, rhttperr := http.NewRequest(method, url, r); rhttperr == nil {
			if rhttp != nil {
				var whttp = NewResponse(w, rhttp)
				whttp.canWriteHeader = false
				chnl.internalServeHTTP(whttp, rhttp, a...)
			}
		}
	}
	wg.Wait()
	wg = nil
}

var gblchnl *Channel

//GLOBALCHNL - Global app *Channel
func GLOBALCHNL() *Channel {
	if gblchnl == nil {
		gblchnl = NewChannel()
	}
	return gblchnl
}
