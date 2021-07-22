package chnls

import (
	"runtime"

	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/requesting"
)

func internalNewRequest(initPath string, chnl *Channel, prntrqst *Request, rqstrw requesting.RequestorAPI, a ...interface{}) (rqst *Request, interrupt func()) {
	defer runtime.GC()
	var rqstsettings map[string]interface{} = nil
	var ai = 0

	var mqttmsg mqtt.Message = nil
	var mqttevent mqtt.MqttEvent = nil
	var aok = false

	for ai < len(a) {
		if da, daok := a[ai].([]interface{}); daok {
			if al := len(da); al > 0 {
				a = append(da, a[1:]...)
				ai = 0
			} else {
				a = a[1:]
			}
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

	rqst = &Request{rqstrw: rqstrw, prntrqst: prntrqst, chnl: chnl, isFirstRequest: true, Interrupted: false, settings: rqstsettings, actnslst: enumeration.NewList(), args: make([]interface{}, len(a)), objmap: map[string]interface{}{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}, activecns: map[string]*database.Connection{}, cmnds: map[int]*osprc.Command{},
		initPath:  initPath,
		mphndlr:   caching.GLOBALMAP().Handler(),
		mqttmsg:   mqttmsg,
		mqttevent: mqttevent}
	if len(rqst.args) > 0 {
		copy(rqst.args[:], a[:])
	}

	interrupt = func() {
		rqst.Interrupt()
	}
	return
}

/*func internalExecuteRequest(rqst *Request, interrupt func()) {
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
		notify := func() func() <-chan bool {
			var clsntfy http.CloseNotifier = nil
			if httpw != nil {
				clsntfy, _ = httpw.(http.CloseNotifier)
			} else if rqstw != nil {
				clsntfy, _ = rqstw.(http.CloseNotifier)
			}
			if clsntfy != nil {
				return clsntfy.CloseNotify
			}
			return nil
		}()
		isCancelled := false
		ctx, cancel := context.WithCancel(bgrndctnx)
		go func() {
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
		if notify != nil {
			select {
			case <-notify():
				if interrupt != nil {
					interrupt()
					interrupt = nil
				}
			case <-ctx.Done():
				if ctxerr := ctx.Err(); ctxerr != nil {
					if !isCancelled {
						if interrupt != nil {
							interrupt()
							interrupt = nil
						}
					}
				}
			}
		} else {
			<-ctx.Done()
			if ctxerr := ctx.Err(); ctxerr != nil {
				if !isCancelled {
					if interrupt != nil {
						interrupt()
						interrupt = nil
					}
				}
			}
		}
	}()
}*/

func internalExecuteRequest(rqst *Request, interrupt func()) {
	rqst.executeNow(interrupt)
}

func processingRequestIO(initpath string, chnl *Channel, prntrqst *Request, rqstrw requesting.RequestorAPI, a ...interface{}) {
	var excrqst *Request = nil
	var interrupt func() = nil
	if prntrqst == nil {
		excrqst, interrupt = internalNewRequest(initpath, chnl, prntrqst, rqstrw, a...)
	} else {
		excrqst = prntrqst
	}
	if excrqst != nil {
		func() {
			if prntrqst != excrqst {
				defer excrqst.Close()
			}
			internalExecuteRequest(excrqst, func() {
				if interrupt != nil {
					interrupt()
				}
			})
		}()
		excrqst = nil
	}
}
