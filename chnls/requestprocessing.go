package chnls

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"

	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
)

func internalNewRequest(initPath string, chnl *Channel, prntrqst *Request, rdr io.Reader, wtr io.Writer, httpw http.ResponseWriter, httpr *http.Request, httpflshr http.Flusher, a ...interface{}) (rqst *Request, interrupt func()) {
	defer runtime.GC()
	var rqstsettings map[string]interface{} = nil
	var ai = 0
	var rqstr iorw.Reader = nil
	var remoteHost = ""
	var localHost = ""
	var mqttmsg mqtt.Message = nil
	var mqttevent mqtt.MqttEvent = nil
	var aok = false
	if rdr != nil {
		if rqstr, _ = rdr.(iorw.Reader); rqstr == nil {
			rqstr = iorw.NewEOFCloseSeekReader(rdr)
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
		} else if cnctn, cnctnok := a[ai].(*listen.ConnHandler); cnctnok {
			if cnctn != nil {
				remoteHost = cnctn.RemoteAddr().String()
				localHost = cnctn.LocalAddr().String()
			}
			a = a[1:]
			continue
		} else if cnctn, cnctnok := a[ai].(net.Conn); cnctnok {
			if cnctn != nil {
				remoteHost = cnctn.RemoteAddr().String()
				localHost = cnctn.LocalAddr().String()
			}
			a = a[1:]
			continue
		} else if mqttmsg, aok = a[ai].(mqtt.Message); aok {
			a = a[1:]
			continue
		} else if mqttevent, aok = a[ai].(mqtt.MqttEvent); aok {
			a = a[1:]
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

	rqst = &Request{prntrqst: prntrqst, chnl: chnl, rmtHost: remoteHost, lclHost: localHost, isFirstRequest: true, mimetype: "", zpw: nil, Interrupted: false, startedWriting: false, wbytes: make([]byte, 8192), wbytesi: 0, flshr: httpflshr, rqstw: wtr, httpw: httpw, rqstr: rqstr, httpr: httpr, settings: rqstsettings, actnslst: enumeration.NewList(), args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}, activecns: map[string]*database.Connection{}, cmnds: map[int]*osprc.Command{},
		initPath:      initPath,
		mphndlr:       caching.GLOBALMAP().Handler(),
		mediarqst:     false,
		rqstoffset:    -1,
		rqstendoffset: -1,
		rqstoffsetmax: -1,
		rqstmaxsize:   -1,
		mqttmsg:       mqttmsg,
		mqttevent:     mqttevent}
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
	} else {
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

		func() {
			defer func() {
				if r := recover(); r != nil {

				}
				isCancelled = true
				cancel()
			}()
			if httpr != nil && httpw != nil {
				rqst.executeHTTP(interrupt)
			} else if (rqstr == nil || rqstw == nil) || (rqstr != nil && rqstw != nil) {
				if rwerr := rqst.executeRW(interrupt); rwerr != nil {
					fmt.Println(rwerr)
				}
			} else {
				rqst.executePath("", interrupt)
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

func processingRequestIO(initpath string, chnl *Channel, prntrqst *Request, rdr io.Reader, wtr io.Writer, httpw http.ResponseWriter, httpflshr http.Flusher, httpr *http.Request, a ...interface{}) {
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
	if prntrqst == nil {
		excrqst, interrupt = internalNewRequest(initpath, chnl, prntrqst, rdr, wtr, httpw, httpr, httpflshr, a...)
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
}
