package chnls

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/resources"
)

func internalNewRequest(chnl *Channel, prntrqst *Request, rdr func() io.Reader, wtr func() io.Writer, httpw func() http.ResponseWriter, httpr func() *http.Request, httpflshr func() http.Flusher, rqstsettings map[string]interface{}, a ...interface{}) (rqst *Request, interrupt func()) {
	var ai = 0
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
		}
		ai++
	}
	if rqstsettings == nil {
		rqstsettings = map[string]interface{}{}
	}
	rqst = &Request{prntrqst: prntrqst, chnl: chnl, isFirstRequest: true, mimetype: "", zpw: nil, Interrupted: false, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, rqstw: wtr, httpw: httpw, rqstr: rdr, httpr: httpr, settings: rqstsettings, rsngpthsref: map[string]*resources.ResourcingPath{}, actnslst: enumeration.NewList(), actns: []*Action{}, args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{} /*, embeddedResources: map[string]interface{}{}*/, activecns: map[string]*database.Connection{}, cmnds: map[int]*osprc.Command{},
		initPath:      "",
		mediarqst:     false,
		rqstoffset:    -1,
		rqstendoffset: -1,
		rqstoffsetmax: -1,
		rqstmaxsize:   -1}
	//rqst.invokeAtv()
	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}

	interrupt = func() {
		rqst.Interrupt()
	}
	return
}

func internalExecuteRequest(rqst *Request, interrupt func()) {
	lck := &sync.Mutex{}
	lck.Lock()
	defer lck.Unlock()
	var bgrndctnx context.Context = nil
	httpr, httpw, rqstw, rqstr := rqst.httpr(), rqst.httpw(), rqst.rqstw(), rqst.rqstr()
	if httpr != nil && httpw != nil {
		rqst.prtcl = httpr.Proto
		rqst.prtclmethod = httpr.Method
		bgrndctnx = httpr.Context()
	} else if rqstw != nil && rqstr != nil {
		bgrndctnx = context.Background()
	}
	func() {
		isCancelled := false
		ctx, cancel := context.WithCancel(bgrndctnx)
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
			if httpr != nil && httpw != nil {
				rqst.executeHTTP(interrupt)
			} else if rqstr != nil && rqstw != nil {
				if rwerr := rqst.executeRW(interrupt); rwerr != nil {
					fmt.Println(rwerr)
				}
			}
		}()
		<-ctx.Done()
		if ctxerr := ctx.Err(); ctxerr != nil {
			if !isCancelled {
				if interrupt != nil {
					interrupt()
				}
			}
		}
	}()
}

func processingRequestIO(chnl *Channel, prntrqst *Request, rdr func() io.Reader, wtr func() io.Writer, httpw func() http.ResponseWriter, httpflshr func() http.Flusher, httpr func() *http.Request, a ...interface{}) {
	var rqstsettings map[string]interface{} = nil
	var ai = 0
	var excrqst *Request = nil
	var interrupt func() = nil
	if wtr == nil && httpw != nil {
		wtr = func() io.Writer { return httpw() }
	}

	if httpw != nil && httpflshr == nil {
		if flshr, flshrok := httpw().(http.Flusher); flshrok {
			httpflshr = func() http.Flusher { return flshr }
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
		} else if prnstrq, prntrqok := a[ai].(*Request); prntrqok {
			if prntrqst == nil {
				prntrqst = prnstrq
			}
			a = a[1:]
			continue
		} else if rstngs, rstngsok := a[ai].(map[string]interface{}); rstngsok {
			if rstngs != nil {
				if rqstsettings == nil {
					rqstsettings = rstngs
				} else {
					for k, v := range rstngs {
						rqstsettings[k] = v
					}
				}
			}
			a = a[1:]
			continue
		}
		ai++
	}
	if prntrqst == nil {
		excrqst, interrupt = internalNewRequest(chnl, prntrqst, rdr, wtr, httpw, httpr, httpflshr, rqstsettings, a...)
	} else {
		excrqst = prntrqst
	}
	if excrqst != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
				excrqst.Close()
			}()
			internalExecuteRequest(excrqst, interrupt)
			/*internalExecuteRequest(excrqst, interrupt,
			func() io.Writer {
				if wtr != nil {
					return wtr()
				}
				return nil
			}(),
			func() io.Reader {
				if rdr != nil {
					return rdr()
				}
				return nil
			}(),
			func() *http.Request {
				if httpr != nil {
					return httpr()
				}
				return nil
			}(),
			func() http.ResponseWriter {
				if httpw != nil {
					return httpw()
				}
				return nil
			}())*/
		}()
		excrqst = nil
	}
	if httpflshr != nil {
		httpflshr().Flush()
	}
}
