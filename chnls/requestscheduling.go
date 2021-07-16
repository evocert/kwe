package chnls

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/scheduling/ext"
)

//Schedule refer to scheduling.ScheduleHandler - StartedSchedule()
func (rqst *Request) Schedule() *ext.Schedule {
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
	err = rqst.Close()
	return
}

//RequestScheduleAction - struct implementing scheduling.ActionHandler and wrapping *Request
type RequestScheduleAction struct {
	*ext.ScheduleAction
	rqst    *Request
	atvfunc func(goja.FunctionCall) goja.Value
}

func (rqst *Request) executeSchdlRequest(a ...interface{}) (err error) {
	for len(a) > 0 {
		d := a[0]
		a = a[1:]
		if s, sok := d.(string); sok {
			if s != "" {
				rqst.AddPath(s)
				rqst.processPaths(false)
			}
		} else if args, argsok := d.([]interface{}); argsok {
			if len(args) > 0 {
				a = append(a, args...)
			}
		} else if rqstmp, rqstmpok := d.(map[string]interface{}); rqstmpok {
			for rqstk, rqstv := range rqstmp {
				if rqstk == "path" {
					a = append(a, rqstv)
				} else if rqstk == "paths" {
					a = append(a, rqstv)
				}
			}
		}
	}
	return
}

func (rqst *Request) executeSchdlDbms(a ...interface{}) (err error) {
	if len(a) > 0 {
		err = database.GLOBALDBMS().InOut(a[0], nil, a[1:]...)
	}
	return
}

func (rqst *Request) executeSchdlCommand(a ...interface{}) (err error) {

	return
}

func (rqst *Request) executeSchdlScript(a ...interface{}) (err error) {
	fmt.Print(a...)
	return
}

func (rqst *Request) executeScheduleAction(a ...interface{}) (err error) {
	for len(a) > 1 {
		if cmd, cmdok := a[0].(string); cmdok && cmd != "" {
			a = a[1:]
			var cmdfnctoexec func(...interface{}) error = nil
			if cmd == "dbms" {
				cmdfnctoexec = rqst.executeSchdlDbms
			} else if cmd == "request" {
				cmdfnctoexec = rqst.executeSchdlRequest
			} else if cmd == "command" {
				cmdfnctoexec = rqst.executeSchdlCommand
			} else if cmd == "script" {
				cmdfnctoexec = rqst.executeSchdlScript
			}
			if cmdfnctoexec != nil {
				if cmdmap, cmdmapok := a[0].(map[string]interface{}); cmdmapok && len(cmdmap) > 0 {
					cmdfnctoexec(cmdmap)
				} else if cmdargs, cmdargsok := a[0].([]interface{}); cmdargsok && len(cmdargs) > 0 {
					cmdfnctoexec(cmdargs)
				} else if cmdarg, cmdargok := a[0].(interface{}); cmdargok {
					if cmdarg != nil {
						cmdfnctoexec(cmdarg)
					}
				} else {
					break
				}
				a = a[1:]
			} else {
				break
			}
		} else {
			break
		}
	}

	return
}

//PrepActionArgs refer to scheduling.ScheduleHandler - PrepActionArgs()
func (rqst *Request) PrepActionArgs(a ...interface{}) (preppedargs []interface{}, err error) {
	if al := len(a); al > 0 {
		ai := 0
		regatvfnc := func(atvfnc func(goja.FunctionCall) goja.Value) bool {
			if atvfnc != nil {
				var prppdatvfnc ext.FuncArgsErrHandle = nil
				prppdatvfnc = func(args ...interface{}) (rserr error) {
					rqst.invokeAtv()
					if rslt := rqst.atv.InvokeFunction(atvfnc, args...); rslt != nil {
						if dne, dneok := rslt.(bool); dneok && dne {
							rserr = fmt.Errorf("DONE")
						}
					}
					return
				}
				a[ai] = nil
				a[ai] = prppdatvfnc
				return true
			}
			return false
		}
		for ai < al {
			d := a[ai]
			if sfnc, sfncok := d.(string); sfncok {
				if sfnc != "" {
					if strings.HasPrefix(sfnc, "function(") {
						rqst.atv.InvokeVM(func(vm *goja.Runtime) (vmerr error) {
							atvfncval, _ := vm.RunString("(" + sfnc + ")")
							var atvfncref func(goja.FunctionCall) goja.Value = nil
							vm.ExportTo(atvfncval, &atvfncref)
							if !regatvfnc(atvfncref) {

							}
							return
						})
					}
				}
			}
			if atvfnc, atvfcnok := d.(func(goja.FunctionCall) goja.Value); atvfcnok {
				if rqst.prntrqst != nil && rqst.prntrqst.atv != nil {
					rqst.prntrqst.atv.InvokeVM(func(vm *goja.Runtime) (vmerr error) {
						/*atvfncval, _ := vm.RunString("(" + sfnc + ")")
						var atvfncref func(goja.FunctionCall) goja.Value = nil
						vm.ExportTo(atvfncval, &atvfncref)
						if !regatvfnc(atvfncref) {

						}*/
						vm.Set("xfnc", atvfnc)
						if fncv, _ := vm.RunString("xfnc.toString()"); fncv != nil {
							fncs := ""
							vm.ExportTo(fncv, &fncs)
							if fncs != "" {

							}
						}
						return
					})
				}
				if !regatvfnc(atvfnc) {
					//a[ai] = nil
					//a[ai] = prppdatvfnc
				}
			} else if rqstactnmap, rqstactnmapok := d.(map[string]interface{}); rqstactnmapok {
				ignore := false
				for rqstmk, rqstmv := range rqstactnmap {
					if rqstmk != "" && strings.Contains("|request|dbms|command|script|", "|"+rqstmk+"|") {
						if !ignore {
							ignore = true
						}
						a[ai] = ext.FuncArgsErrHandle(rqst.executeScheduleAction)
						tmpa := a[ai+1:]
						a = append(append(a[:ai+1], []interface{}{rqstmk, rqstmv}), tmpa...)
						al = len(a)
						ai++
						ai++
					}
				}
				if ignore {

					continue
				}
			}
			ai++
		}
		preppedargs = a[:]
	}
	return
}

//OnExecuteAction implementioan that is called by *scheduling.ScheduleAction ExecuteAction()
func (rqstschldactn *RequestScheduleAction) OnExecuteAction(a ...interface{}) (result bool, err error) {

	return
}
