package chnls

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/scheduling"
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
		chnl.mqttmngr = mqtt.NewMQTTManager(func(topic mqtt.Topic, message mqtt.Message, mqttcn *mqtt.MQTTConnection, mqttmngr *mqtt.MQTTManager) {
			processingRequestIO(chnl, nil, nil, nil, nil, nil, nil, []interface{}{topic.TopicPath(), message, mqttcn, mqttmngr}...)
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

//NewSchedule - implement scheduling.ScheduleHandler NewScheduler()
func (chnl *Channel) NewSchedule(schdl *scheduling.Schedule, a ...interface{}) (scdhlhndlr scheduling.ScheduleHandler) {
	if al := len(a); al > 0 {
		ai := 0
		var prntrqst *Request = nil
		atvprntmap := map[string]interface{}{}
		for ai < al {
			d := a[ai]
			if rqst, rqstok := d.(*Request); rqstok {
				prntrqst = rqst
				if rqst.atv != nil {
					rqst.atv.ExtractGlobals(atvprntmap)
				}
				ai++
			} else {
				ai++
			}
		}
		if scdhlrqst, _ := internalNewRequest(chnl, prntrqst, nil, nil, nil, nil, nil, a...); scdhlrqst != nil {
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
	/*if al := len(a); al > 0 {
		ai := 0
		atvprntmap := map[string]interface{}{}
		for ai < al {
			d := a[ai]
			if rqst, rqstok := d.(*Request); rqstok {
				if rqst.atv != nil {
					rqst.atv.ExtractGlobals(atvprntmap)
				}
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
	}*/
	return
}

//NewChannel - instance
func NewChannel() (chnl *Channel) {
	chnl = &Channel{objmap: map[string]interface{}{}}
	return
}

func (chnl *Channel) internalServePath(path string, a ...interface{}) {
	inirspath := path
	a = append([]interface{}{inirspath}, a...)
	processingRequestIO(chnl, nil, nil, nil, nil, nil, nil, a...)
}

func (chnl *Channel) internalServeHTTP(w http.ResponseWriter, r *http.Request, a ...interface{}) {
	inirspath := r.URL.Path
	cnctn := r.Context().Value(listen.ConnContextKey)
	if wsrw, wsrwerr := ws.NewServerReaderWriter(w, r); wsrw != nil && wsrwerr == nil {
		go func() {
			defer wsrw.Close()
			if cnctn == nil {
				a = append([]interface{}{inirspath}, a...)
			} else {
				a = append([]interface{}{inirspath, cnctn}, a...)
			}
			processingRequestIO(chnl, nil, wsrw, wsrw, nil, nil, nil, a...)
		}()
	} else {
		if cnctn != nil {
			a = append([]interface{}{cnctn}, a...)
		}
		processingRequestIO(chnl, nil, r.Body, w, w, nil, r, a...)
	}
}

//ServeHTTP - refer http.Handler
func (chnl *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chnl.internalServeHTTP(w, r)
}

//ServeRW - serve Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	processingRequestIO(chnl, nil, r, w, nil, nil, nil, a...)
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

func (chnl *Channel) DefaultServePath(path string, a ...interface{}) {
	cntxt := context.Background()
	go func() {
		_, cncl := context.WithCancel(cntxt)
		defer cncl()
		chnl.internalServePath(path, a...)
	}()
	<-cntxt.Done()
}

//DefaultServeRW - helper to perform dummy ServeRW request on channel
func (chnl *Channel) DefaultServeRW(w io.Writer, url string, r io.Reader, a ...interface{}) {
	cntxt := context.Background()
	go func() {
		_, cncl := context.WithCancel(cntxt)
		defer cncl()

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
	}()
	<-cntxt.Done()
}

var gblchnl *Channel

//GLOBALCHNL - Global app *Channel
func GLOBALCHNL() *Channel {
	if gblchnl == nil {
		gblchnl = NewChannel()
		gblchnl.MQTT()
		gblchnl.Listener()
	}
	return gblchnl
}
