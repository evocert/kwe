package chnls

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/requesting"
	httprqstng "github.com/evocert/kwe/requesting/http"
	scheduling "github.com/evocert/kwe/scheduling/ext"
	"github.com/evocert/kwe/ws"
)

/*Channel -
 */
type Channel struct {
	objmap   map[string]interface{}
	lstnr    *listen.Listener
	schdls   *scheduling.Schedules
	mqttmngr *mqtt.MQTTManager
}

//Listener - *listen.Listener listener for Channel
func (chnl *Channel) Listener() *listen.Listener {
	if chnl.lstnr == nil {
		chnl.lstnr = listen.NewListener(chnl)
	}
	return chnl.lstnr
}

func (chnl *Channel) MQTT() *mqtt.MQTTManager {
	if chnl.mqttmngr == nil {
		chnl.mqttmngr = mqtt.NewMQTTManager(func(message mqtt.Message) {
			processingRequestIO(message.TopicPath(), chnl, nil, nil, nil, []interface{}{message}...)
		})
	}
	return chnl.mqttmngr
}

//Schedules - *scheduling.Schedules schedules for Channel
func (chnl *Channel) Schedules() *scheduling.Schedules {
	if chnl.schdls == nil {
		chnl.schdls = scheduling.NewSchedules(chnl)
	}
	return chnl.schdls
}

//NewSchedule - implement scheduling.ScheduleAPI NewScheduler()
func (chnl *Channel) NewSchedule(schdl *scheduling.Schedule, a ...interface{}) (scdhlhndlr scheduling.ScheduleHandler) {
	if al := len(a); al > 0 {
		ai := 0
		var prntrqst *Request = nil
		atvprntmap := map[string]interface{}{}
		inipath := "/"
		for ai < al {
			d := a[ai]
			if rqst, rqstok := d.(*Request); rqstok {
				prntrqst = rqst
				if rqst.atv != nil {
					rqst.atv.ExtractGlobals(atvprntmap)
				}
				ai++
			} else if ipth, ipthok := d.(string); ipthok {
				if ipth != "" {
					inipath = ipth
				}
				ai++
			} else {
				ai++
			}
		}
		if scdhlrqst, _ := internalNewRequest(inipath, chnl, prntrqst, nil, a...); scdhlrqst != nil {
			scdhlrqst.schdl = schdl
			lclglbs := map[string]interface{}{}
			if len(atvprntmap) > 0 {
				scdhlrqst.invokeAtv()
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
			}
			scdhlhndlr = scdhlrqst
		}
	}
	return
}

//NewChannel - instance
func NewChannel() (chnl *Channel) {
	chnl = &Channel{objmap: map[string]interface{}{}}
	return
}

func (chnl *Channel) internalServePath(path string, a ...interface{}) {
	processingRequestIO(path, chnl, nil, nil, nil, a...)
}

func (chnl *Channel) internalServeHTTP(w http.ResponseWriter, r *http.Request, a ...interface{}) {
	inirspath := r.URL.Path
	cnctn := r.Context().Value(listen.ConnContextKey)
	if wsrw, wsrwerr := ws.NewServerReaderWriter(w, r); wsrw != nil && wsrwerr == nil {
		go func() {
			defer wsrw.Close()
			if cnctn != nil {
				a = append([]interface{}{cnctn}, a...)
			}
			rqstr := httprqstng.NewRequest(inirspath, wsrw)
			processingRequestIO(inirspath, chnl, nil, rqstr, httprqstng.NewResponse(wsrw, rqstr), a...)
		}()
	} else {
		if cnctn != nil {
			a = append([]interface{}{cnctn}, a...)
		}
		rqstr := httprqstng.NewRequest(inirspath, r)
		processingRequestIO(inirspath, chnl, nil, rqstr, httprqstng.NewResponse(w, rqstr), a...)
	}
}

//ServeHTTP - refer http.Handler
func (chnl *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chnl.internalServeHTTP(w, r)
}

//Serve - refer requesting.RequestorHandler
func (chnl *Channel) ServeReaderWriter(path string, w io.Writer, r io.Reader) (err error) {
	rqst := httprqstng.NewRequest(path, r)
	rspns := httprqstng.NewResponse(w, rqst)
	defer func() {
		rqst.Close()
		rspns.Close()
	}()
	err = chnl.Serve(path, rspns.Request(), rspns)
	return
}

//Serve - refer requesting.RequestorHandler
func (chnl *Channel) Serve(path string, rqst requesting.RequestAPI, rspns requesting.ResponseAPI) (err error) {
	processingRequestIO(path, chnl, nil, rqst, rspns)
	return
}

//ServeRW - serve Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	initpath := "/"
	if al := len(a); al > 0 {
		for an, d := range a {
			if pthi, pthiok := d.(string); pthiok {
				if pthi != "" {
					initpath = pthi
				}
				a = append(a[:an], a[an+1:]...)
				break
			}
		}
	}
	rqstr := httprqstng.NewRequest(initpath, r)
	processingRequestIO(initpath, chnl, nil, rqstr, httprqstng.NewResponse(w, rqstr), a...)
}

//Stdio - os.Stdout, os.Stdin
func (chnl *Channel) Stdio(out *os.File, in *os.File, err *os.File, a ...interface{}) {
	chnl.ServeRW(in, out, a...)
}

//Send inline request
func (chnl *Channel) Send(path string, a ...interface{}) (rspr iorw.Reader, err error) {
	if chnl != nil {
		ctx, cancel := context.WithCancel(context.Background())
		pi, pw := io.Pipe()
		go func() {
			defer func() {
				pw.Close()
			}()
			cancel()
			processingRequestIO(path, chnl, nil, nil, httprqstng.NewResponse(pw, nil))
		}()
		ctx.Done()
		rspr = iorw.NewEOFCloseSeekReader(pi, true)
	}
	return
}

func (chnl *Channel) DefaultServePath(path string, a ...interface{}) {
	cntxt := context.Background()
	go func() {
		_, cncl := context.WithCancel(cntxt)
		defer cncl()
		chnl.internalServePath(path, a...)
	}()
	<-cntxt.Done()
}

var gblchnl *Channel

//GLOBALCHNL - Global app *Channel
func GLOBALCHNL() *Channel {
	if gblchnl == nil {
		gblchnl = NewChannel()
		if gblchnl.mqttmngr == nil {
			gblchnl.mqttmngr = mqtt.GLOBALMQTTMANAGER()
			if gblchnl.mqttmngr.MqttMessaging == nil {
				gblchnl.mqttmngr.MqttMessaging = func(message mqtt.Message) {
					processingRequestIO(message.TopicPath(), gblchnl, nil, nil, nil, []interface{}{message}...)
				}
			}
			if gblchnl.mqttmngr.MqttEventing == nil {
				gblchnl.mqttmngr.MqttEventing = func(event mqtt.MqttEvent) {
					processingRequestIO(event.EventPath(), gblchnl, nil, nil, nil, []interface{}{event}...)
				}
			}
		}
		gblchnl.Listener()
	}
	return gblchnl
}
