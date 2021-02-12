package chnls

import (
	"github.com/dop251/goja"
	"github.com/evocert/kwe/scheduling"
)

//Schedule refer to scheduling.ScheduleHandler - StartedSchedule()
func (rqst *Request) Schedule() *scheduling.Schedule {
	return rqst.schdl
}

//StartedSchedule refer to scheduling.ScheduleHandler - StartedSchedule()
func (rqst *Request) StartedSchedule(a ...interface{}) (err error) {
	return
}

//StoppedSchedule refer to scheduling.ScheduleHandler - StoppedSchedule()
func (rqst *Request) StoppedSchedule(a ...interface{}) (err error) {
	return
}

//ShutdownSchedule refer to scheduling.ScheduleHandler - ShutdownSchedule()
func (rqst *Request) ShutdownSchedule() (err error) {
	return
}

//RequestScheduleAction - struct implementing scheduling.ActionHandler and wrapping *Request
type RequestScheduleAction struct {
	*scheduling.ScheduleAction
	rqst    *Request
	atvfunc func(goja.FunctionCall) goja.Value
}

//NewScheduleAction refer to scheduling.ScheduleHandler - NewScheduleAction()
func (rqst *Request) NewScheduleAction(a ...interface{}) (actnhndlr scheduling.ActionHandler) {
	var atvfunc func(goja.FunctionCall) goja.Value = nil
	if al := len(a); al > 0 {
		ai := 0
		for ai < al {
			d := a[ai]
			if atvfnc, atvfcnok := d.(func(goja.FunctionCall) goja.Value); atvfcnok {
				if atvfnc != nil {
					//rqst.atv.InvokeFunction(atvfnc)
					atvfunc = atvfnc
				}
				a = append(a[:ai], a[ai+1:])
			} else {
				ai++
			}
		}
		if atvfunc != nil {

		}
	}
	//var schdlactn *scheduling.ScheduleAction = scheduling.NewScheduleAction()
	return
}

//OnExecuteAction implementioan that is called by *scheduling.ScheduleAction ExecuteAction()
func (rqstschldactn *RequestScheduleAction) OnExecuteAction(a ...interface{}) (result bool, err error) {

	return
}
