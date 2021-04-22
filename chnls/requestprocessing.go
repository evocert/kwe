package chnls

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/osprc"
)

func internalNewRequest(chnl *Channel, prntrqst *Request, rdr io.Reader, wtr io.Writer, httpw http.ResponseWriter, httpr *http.Request, httpflshr http.Flusher, a ...interface{}) (rqst *Request, interrupt func()) {
	var rqstsettings map[string]interface{} = nil
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
	rqst = &Request{prntrqst: prntrqst, chnl: chnl, isFirstRequest: true, mimetype: "", zpw: nil, Interrupted: false, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, rqstw: wtr, httpw: httpw, rqstr: rdr, httpr: httpr, settings: rqstsettings, actnslst: enumeration.NewList(), args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{} /*, embeddedResources: map[string]interface{}{}*/, activecns: map[string]*database.Connection{}, cmnds: map[int]*osprc.Command{},
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
	var bgrndctnx context.Context = nil
	httpr, httpw, rqstw, rqstr := rqst.httpr, rqst.httpw, rqst.rqstw, rqst.rqstr
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

func processingRequestIO(chnl *Channel, prntrqst *Request, rdr io.Reader, wtr io.Writer, httpw http.ResponseWriter, httpflshr http.Flusher, httpr *http.Request, a ...interface{}) {
	//func() {
	//var rqstsettings map[string]interface{} = nil
	//var ai = 0
	var excrqst *Request = nil
	var interrupt func() = nil
	if wtr == nil && httpw != nil {
		wtr = httpw
	}

	if httpw != nil && httpflshr == nil {
		if flshr, flshrok := httpw.(http.Flusher); flshrok {
			httpflshr = flshr
		}
	}

	/*for ai < len(a) {
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
	}*/
	if prntrqst == nil {
		excrqst, interrupt = internalNewRequest(chnl, prntrqst, rdr, wtr, httpw, httpr, httpflshr, a...)
	} else {
		excrqst = prntrqst
	}
	if excrqst != nil {
		internalExecuteRequest(excrqst, func() {
			if interrupt != nil {
				interrupt()
			}
		})
		excrqst = nil
	}
	//}()
}
