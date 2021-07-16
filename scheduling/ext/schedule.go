package ext

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/enumeration"
)

type ScheduleAPI interface {
	Schedules() SchedulesAPI
	Start(...interface{}) error
	Stop() error
	Shutdown() error
}

type ScheduleHandler interface {
	StartedSchedule(...interface{}) error
	StoppedSchedule(...interface{}) error
	ShutdownSchedule() error
	PrepActionArgs(...interface{}) ([]interface{}, error)
	Schedule() *Schedule
}

type FuncArgsErrHandle func(...interface{}) error
type FuncArgsHandle func(...interface{})
type FuncErrHandle func() error
type FuncHandle func(...interface{})

type scheduleactionsection int

const (
	schdlactnmain scheduleactionsection = iota
	schdlactninit
	schdlactnwrapup
)

type Schedule struct {
	actnmde        scheduleactionsection
	initstart      bool
	schdlid        string
	once           bool
	schdls         SchedulesAPI
	schdlhndlr     ScheduleHandler
	From           time.Time
	To             time.Time
	initactns      *enumeration.List
	lckinitactns   *sync.RWMutex
	actns          *enumeration.List
	lckactns       *sync.RWMutex
	wrapupactns    *enumeration.List
	lckwrapupactns *sync.RWMutex
	StartArgs      []interface{}
	OnError        func(SchedulesAPI, *Schedule, error)
	OnStart        func(a ...interface{}) error
	StopArgs       []interface{}
	OnStop         func(a ...interface{}) error
	OnShutdown     func() error
	Milliseconds   int64
	Seconds        int64
	Minutes        int64
	Hours          int64
	prcintrvl      int64
	intrvl         time.Duration
	running        bool
	wg             *sync.WaitGroup
}

func NewSchedule(a ...interface{}) (schdl *Schedule) {
	var schdls SchedulesAPI = nil

	var start func(a ...interface{}) error = nil
	var startargs []interface{} = nil
	var stop func(a ...interface{}) error = nil
	var stopargs []interface{} = nil
	var shutdown func() error = nil
	var milliseconds int64 = 0
	var seconds int64 = 0
	var minutes int64 = 0
	var hours int64 = 0
	var once = false
	var frm time.Time = time.Now()
	frm = time.Date(frm.Year(), frm.Month(), frm.Day(), 0, 0, 0, 0, frm.Location())
	var to time.Time = frm.Add(time.Hour * 24)

	if al := len(a); al > 0 {
		ai := 0
		for ai < al {
			d := a[ai]
			if dschdls, dschdlsok := d.(SchedulesAPI); dschdlsok {
				if dschdls != nil && schdls == nil {
					schdls = dschdls
				}
				a = append(a[:ai], a[ai+1:]...)
				al--
				ai++
				continue
			} else if dmp, dmpok := a[0].(map[string]interface{}); dmpok {
				for stngk, stngv := range dmp {
					if strings.ToLower(stngk) == "start" {
						start, _ = stngv.(func(...interface{}) error)
					} else if strings.ToLower(stngk) == "startargs" {
						startargs, _ = stngv.([]interface{})
					} else if strings.ToLower(stngk) == "stop" {
						stop, _ = stngv.(func(...interface{}) error)
					} else if strings.ToLower(stngk) == "stopargs" {
						stopargs, _ = stngv.([]interface{})
					} else if strings.ToLower(stngk) == "shutdown" {
						shutdown, _ = stngv.(func() error)
					} else if strings.ToLower(stngk) == "milliseconds" && hours == 0 && seconds == 0 && minutes == 0 {
						if scint, scintok := stngv.(int); scintok {
							milliseconds = int64(scint)
						} else {
							milliseconds, _ = stngv.(int64)
						}
					} else if strings.ToLower(stngk) == "seconds" && hours == 0 && milliseconds == 0 && minutes == 0 {
						if scint, scintok := stngv.(int); scintok {
							seconds = int64(scint)
						} else {
							seconds, _ = stngv.(int64)
						}
					} else if strings.ToLower(stngk) == "minutes" && hours == 0 && seconds == 0 && milliseconds == 0 {
						if scint, scintok := stngv.(int); scintok {
							minutes = int64(scint)
						} else {
							minutes, _ = stngv.(int64)
						}
					} else if strings.ToLower(stngk) == "hours" && milliseconds == 0 && seconds == 0 && minutes == 0 {
						if scint, scintok := stngv.(int); scintok {
							hours = int64(scint)
						} else {
							hours, _ = stngv.(int64)
						}
					} else if strings.ToLower(stngk) == "once" {
						once, _ = stngv.(bool)
					} else if strings.ToLower(stngk) == "from" {
						if tmpstp, _ := stngv.(string); tmpstp != "" {
							if tmptme, tmptmeerr := time.Parse(time.RFC3339, strings.Replace(tmpstp, " ", "T", -1)); tmptmeerr == nil {
								frm = tmptme
							}
						} else if tmpstmt, tmpstmptok := stngv.(time.Time); tmpstmptok {
							frm = tmpstmt
						}
					} else if strings.ToLower(stngk) == "to" {
						if tmpstp, _ := stngv.(string); tmpstp != "" {
							if tmptme, tmptmeerr := time.Parse(time.RFC3339, strings.Replace(tmpstp, " ", "T", -1)); tmptmeerr == nil {
								to = tmptme
							}
						} else if tmpstmt, tmpstmptok := stngv.(time.Time); tmpstmptok {
							to = tmpstmt
						}
					}
				}
			}
			ai++
		}
	}

	schdl = &Schedule{
		wg:           &sync.WaitGroup{},
		initstart:    true,
		actnmde:      schdlactninit,
		schdls:       schdls,
		once:         once,
		OnStart:      start,
		StartArgs:    startargs,
		OnStop:       stop,
		StopArgs:     stopargs,
		OnShutdown:   shutdown,
		Milliseconds: milliseconds,
		Seconds:      seconds,
		Minutes:      minutes,
		Hours:        hours,
		From:         frm,
		To:           to,
		initactns:    enumeration.NewList(true), lckinitactns: &sync.RWMutex{},
		actns: enumeration.NewList(true), lckactns: &sync.RWMutex{},
		wrapupactns: enumeration.NewList(true), lckwrapupactns: &sync.RWMutex{}}

	if schdls != nil {
		if schdls.Handler() != nil {
			schdl.schdlhndlr = schdls.Handler().NewSchedule(schdl, a...)
			if schdl.OnStart == nil {
				schdl.OnStart = schdl.schdlhndlr.StartedSchedule
			}
			if schdl.OnStop == nil {
				schdl.OnStop = schdl.schdlhndlr.StoppedSchedule
			}
			if schdl.OnShutdown == nil {
				schdl.OnShutdown = schdl.schdlhndlr.ShutdownSchedule
			}
		}
	}
	return
}

//AddAction - add action(s) to *Schedule
func (schdl *Schedule) AddAction(a ...interface{}) (err error) {
	err = internalAction(schdl, schdlactnmain, a...)
	return
}

//AddInitAction - add action(s) to *Schedule that will be execute initially
func (schdl *Schedule) AddInitAction(a ...interface{}) (err error) {
	err = internalAction(schdl, schdlactninit, a...)
	return
}

//AddWrapupAction - add action(s) to *Schedule that will be execute when there are no more
// main list fo action(s) to execute
func (schdl *Schedule) AddWrapupAction(a ...interface{}) (err error) {
	err = internalAction(schdl, schdlactnwrapup, a...)
	return
}

func internalAction(schdl *Schedule, actntpe scheduleactionsection, a ...interface{}) (err error) {
	var lstargs []interface{} = nil
	var lstactn func(...interface{}) error = nil
	var al = 0
	var vldactions = []*schdlaction{}
	var cactn func(...interface{}) error = nil
	if schdl.schdlhndlr != nil && len(a) > 0 {
		if preppedargs, preppederr := schdl.schdlhndlr.PrepActionArgs(a...); preppederr == nil {
			if len(preppedargs) > 0 {
				a = preppedargs[:]
			}
		} else {
			err = preppederr
			return
		}
	}
	for {
		if al = len(a); al > 0 {
			d := a[0]
			a = a[1:]
			if args, argsok := d.([]interface{}); argsok {
				if al > 1 {
					d = a[0]
					if atcne, actneok := d.(func(...interface{}) error); actneok {
						vldactions = append(vldactions, newSchdlAction(schdl, actntpe, atcne, args...))
						a = a[1:]
						lstargs = nil
					} else if actn, actnok := d.(func(...interface{})); actnok {
						vldactions = append(vldactions, newSchdlAction(schdl, actntpe, func(fna ...interface{}) (fnerr error) {
							func() {
								defer func() {
									if rv := recover(); rv != nil {
										fnerr = fmt.Errorf("%v", rv)
									}
								}()
								actn(fna...)
							}()
							return fnerr
						}, args...))
						a = a[1:]
						lstargs = nil
					} else {
						lstargs = args[:]
					}
				} else {
					if lstactn != nil {
						vldactions = append(vldactions, newSchdlAction(schdl, actntpe, lstactn, args...))
					}
					break
				}
			} else {
				if cactn != nil {
					cactn = nil
				}
				d = interface{}(d)
				if actnae, actnaeok := d.(FuncArgsErrHandle); actnaeok {
					cactn = actnae
				} else if actna, actnaok := d.(FuncArgsHandle); actnaok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							actna(fna...)
						}()
						return fnerr
					}
				} else if actne, actneok := d.(FuncErrHandle); actneok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							fnerr = actne()
						}()
						return fnerr
					}
				} else if actn, actnok := d.(FuncHandle); actnok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							actn()
						}()
						return fnerr
					}
				} else if actnae, actnaeok := d.(func(...interface{}) error); actnaeok {
					cactn = actnae
				} else if actna, actnaok := d.(func(...interface{})); actnaok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							actna(fna...)
						}()
						return fnerr
					}
				} else if actne, actneok := d.(func() error); actneok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							fnerr = actne()
						}()
						return fnerr
					}
				} else if actn, actnok := d.(func()); actnok {
					cactn = func(fna ...interface{}) (fnerr error) {
						func() {
							defer func() {
								if rv := recover(); rv != nil {
									fnerr = fmt.Errorf("%v", rv)
								}
							}()
							actn()
						}()
						return fnerr
					}
				}
				if cactn != nil {
					if al > 1 {
						if lstargs != nil {
							vldactions = append(vldactions, newSchdlAction(schdl, actntpe, cactn, lstargs...))
							lstargs = nil
						} else {
							d = a[0]
							if args, argsok := d.([]interface{}); argsok {
								vldactions = append(vldactions, newSchdlAction(schdl, actntpe, cactn, args...))
								a = a[1:]
							} else {
								lstactn = cactn
								a = a[1:]
							}
						}
					} else {
						vldactions = append(vldactions, newSchdlAction(schdl, actntpe, cactn, lstargs...))
						break
					}
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	if len(vldactions) > 0 {
		addactns(schdl, actntpe, vldactions...)
	}
	return
}

func addactns(schdl *Schedule, actntpe scheduleactionsection, schdlactns ...*schdlaction) {
	for len(schdlactns) > 0 {
		schlactn := schdlactns[0]
		addactn(schdl, actntpe, schlactn)
		schdlactns = schdlactns[1:]
	}
}

func addactn(schdl *Schedule, actntpe scheduleactionsection, schdlactn *schdlaction) {
	if schdl != nil {
		if schdlactn != nil {
			switch actntpe {
			case schdlactnmain:
				if schdl.actns != nil {
					func() {
						schdl.lckactns.Lock()
						defer schdl.lckactns.Unlock()
						schdl.actns.Push(nil, nil, schdlactn)
					}()
				}
			case schdlactninit:
				if schdl.initactns != nil {
					func() {
						schdl.lckinitactns.Lock()
						defer schdl.lckinitactns.Unlock()
						schdl.initactns.Push(nil, nil, schdlactn)
					}()
				}
			case schdlactnwrapup:
				if schdl.wrapupactns != nil {
					func() {
						schdl.lckwrapupactns.Lock()
						defer schdl.lckwrapupactns.Unlock()
						schdl.wrapupactns.Push(nil, nil, schdlactn)
					}()
				}
			}
		}
	}
}

/*func removeactns(schdl *Schedule, schdlactns ...*schdlaction) {
	for len(schdlactns) > 0 {
		removeactn(schdl, schdlactns[0])
		schdlactns = schdlactns[1:]
	}
}*/

func removeactn(schdl *Schedule, schdlactn *schdlaction) {
	if schdlactn != nil && schdl != nil {
		var rmvctncall = func(actnlst *enumeration.List, actnslck *sync.RWMutex) {
			actnslck.Lock()
			defer actnslck.Unlock()
			actnlst.ValueNode(schdlactn).Dispose(nil, nil)
		}
		switch schdlactn.actnsctn {
		case schdlactninit:
			rmvctncall(schdl.initactns, schdl.lckinitactns)
		case schdlactnmain:
			rmvctncall(schdl.actns, schdl.lckactns)
		case schdlactnwrapup:
			rmvctncall(schdl.wrapupactns, schdl.lckwrapupactns)
		}
		schdlactn = nil
	}
}

type schdlaction struct {
	crntschdlactn ScheduleActionAPI
	schdl         *Schedule
	actnsctn      scheduleactionsection
	args          []interface{}
	actn          func(...interface{}) error
	valid         bool
}

func newSchdlAction(schdl *Schedule, actnsctn scheduleactionsection, actn func(...interface{}) error, a ...interface{}) (scdhlactn *schdlaction) {
	scdhlactn = &schdlaction{schdl: schdl, actnsctn: actnsctn,
		actn: actn, args: a, valid: true}
	return
}

func (schdlctn *schdlaction) dispose() {
	if schdlctn != nil {
		if schdlctn.schdl != nil {
			removeactn(schdlctn.schdl, schdlctn)
		}
		if schdlctn.schdl != nil {
			schdlctn.schdl = nil
		}
		if schdlctn.actn != nil {
			schdlctn.actn = nil
		}
		if schdlctn.args != nil {
			schdlctn.args = nil
		}
		if schdlctn.crntschdlactn != nil {
			schdlctn.crntschdlactn = nil
		}
		schdlctn = nil
	}
}

func (scdhlctn *schdlaction) execute() (err error) {
	if scdhlctn != nil {
		err = scdhlctn.actn(scdhlctn.args...)
		if err != nil && strings.ToLower(err.Error()) == "done" {
			scdhlctn.valid = false
		}
	}
	return
}

func (schdl *Schedule) Schedules() (schdls SchedulesAPI) {
	if schdl != nil {
		schdls = schdl.schdls
	}
	return
}

func (schdl *Schedule) Start(a ...interface{}) (err error) {
	if schdl != nil {
		if !schdl.running {
			ctx, ctxcancel := context.WithCancel(context.Background())

			go func() {
				defer ctxcancel()
				if schdl.OnStart != nil {
					err = schdl.OnStart(schdl.StartArgs...)
				}
				if err == nil {
					schdl.running = true
				}
			}()
			<-ctx.Done()
			if schdl.running {
				schdl.wg.Add(1)
				go ticking(schdl)
				schdl.wg.Wait()
			}
		}
	}
	return
}

func ticking(schdl *Schedule) {
	schdl.wg.Done()
	tckwg := &sync.WaitGroup{}
	var errprcng error = nil
	var prcngdone bool = false
	var nxttrggrstmp, frmstmp, tostmp time.Time
	frmstmp = schdl.From
	tostmp = schdl.To
	tostmp = time.Date(tostmp.Year(), tostmp.Month(), tostmp.Day(), tostmp.Hour(), tostmp.Minute(), tostmp.Second(), tostmp.Nanosecond()-1, tostmp.Location())
	nxttrggrstmp = frmstmp.Add(time.Nanosecond * 0)
	var intrvl time.Duration = schdl.intrvl
	var recheck bool = false
	var calcnxttrggr = func() (cantrggr bool) {
		tmpNow := time.Now()
		tmpfrm := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), frmstmp.Hour(), frmstmp.Minute(), frmstmp.Second(), frmstmp.Nanosecond(), frmstmp.Location())
		tmpto := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tostmp.Hour(), tostmp.Minute(), tostmp.Second(), tostmp.Nanosecond(), tostmp.Location())
		if tmpNow.After(tmpfrm) && tmpNow.Before(tmpto) {
			if cantrggr {
				cantrggr = false
			}
			if nxttrggrstmp.Before(tmpfrm) || recheck {
				if recheck {
					recheck = false
				}
				nxttrggrstmp = time.Date(tmpfrm.Year(), tmpfrm.Month(), tmpfrm.Day(), tmpfrm.Hour(), tmpfrm.Minute(), tmpfrm.Second(), tmpfrm.Nanosecond(), tmpfrm.Location())
			}
			if secdif := int64(time.Duration(schdl.prcintrvl)); secdif > 0 {
				if tmdif := int64(tmpNow.Sub(tmpfrm)); tmdif > 0 {
					tf := (tmdif / secdif)
					if tmpfrm.Add(time.Duration(tf * secdif)).Before(tmpNow) {
						if nxttrggrstmp.Before(tmpfrm.Add(time.Duration(tf * secdif)).Add(time.Nanosecond * (1))) {
							nxttrggrstmp = tmpfrm.Add(time.Duration((tf + 1) * secdif))
							cantrggr = true
						}
					}
				}
			}
		} else {
			cantrggr = false
		}
		return
	}

	var crntprcintrvl = func() int64 {
		if schdl.Milliseconds > 0 {
			return schdl.Milliseconds * int64(time.Millisecond)
		} else if schdl.Seconds > 0 {
			return schdl.Seconds * int64(time.Second)
		} else if schdl.Minutes > 0 {
			return schdl.Minutes * int64(time.Minute)
		} else if schdl.Hours > 0 {
			return schdl.Hours * int64(time.Hour)
		}
		return 0
	}

	for schdl.running {
		if cnrtsec := crntprcintrvl(); cnrtsec > int64(intrvl) {
			if schdl.prcintrvl != cnrtsec {
				schdl.prcintrvl = cnrtsec
				recheck = true
			}
			if calcnxttrggr() {
				func() {
					defer tckwg.Wait()
					tckwg.Add(1)
					go func() {
						defer tckwg.Done()
						if prcngdone, errprcng = process(schdl); errprcng != nil {
							if schdl.OnError != nil {
								schdl.OnError(schdl.schdls, schdl, errprcng)
							}
						}
					}()
				}()
				if prcngdone {
					schdl.running = false
				}
			} else {
				time.Sleep(intrvl)
			}
		} else {
			time.Sleep(intrvl)
		}
	}
	if schdl != nil && (prcngdone) {
		errprcng = schdl.Shutdown()
	}
}

func process(schdl *Schedule) (done bool, err error) {
	if schdl != nil {
		func() {
			defer func() {
				if rv := recover(); rv != nil {
					err = fmt.Errorf("%v", rv)
				}
			}()
			done, err = execute(schdl)
		}()
	}
	return
}

func execute(schdl *Schedule) (done bool, err error) {
	if schdl != nil {
		var nextactns bool = false
		if nextactns = (schdl.actnmde == schdlactninit && schdl.initstart); nextactns {
			nextactns, err = executeInit(schdl)
		}
		if (!nextactns || nextactns) && (schdl.actnmde == schdlactnmain) {
			done, err = executeMain(schdl)
		}
	}
	return
}

func (schdl *Schedule) doLink(lnk *enumeration.Node, d interface{}) (done bool, err error) {
	if d != nil {
		if schdlactn, schdlactnok := d.(*schdlaction); schdlactnok {
			func() {
				defer func() {
					if rv := recover(); rv != nil {
						err = fmt.Errorf("%v", rv)
					}
				}()
				if err = schdlactn.actn(schdlactn.args...); err != nil {
					if strings.ToLower(err.Error()) == "done" {
						schdlactn.valid = false
						err = nil
						done = true
					}
				} else {
					if schdl.actnmde == schdlactninit || schdl.actnmde == schdlactnwrapup {
						done = true
					}
				}
			}()
		}
	}
	return
}

func (schdl *Schedule) errDoLink(lnk *enumeration.Node, d interface{}, err error) (done bool) {
	if schdl.actnmde != schdlactnmain {
		done = true
	}
	return
}

func (schdl *Schedule) doneLink(lnk *enumeration.Node) (err error) {

	return
}

func (schdl *Schedule) disposeLink(lnk *enumeration.Node) {
	if schdl != nil && lnk != nil {
		if schdlactn, _ := lnk.Value().(*schdlaction); schdlactn != nil {
			schdlactn.dispose()
		} else {
			lnk.Dispose(nil, nil)
		}
	}
}

func (schdl *Schedule) errDoneLink(lnk *enumeration.Node, err error) (done bool) {
	if schdl.actnmde != schdlactnmain {
		done = true
	}
	return
}

func executeMain(schdl *Schedule) (done bool, err error) {
	if actnsl := schdl.actns.Length(); actnsl > 0 {
		schdl.actns.Iterate(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink, schdl.disposeLink, nil, nil)
	}
	if schdl.actns.Length() == 0 || schdl.once {
		schdl.actnmde = schdlactnwrapup
		if actnsl := schdl.wrapupactns.Length(); actnsl > 0 {
			schdl.wrapupactns.Iterate(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink, schdl.disposeLink, nil, nil)
		}
		if done = (schdl.actns.Length() == 0 || schdl.once); !done {
			schdl.actnmde = schdlactnmain
		}
	}
	return
}

func executeInit(schdl *Schedule) (nextactns bool, err error) {
	if schdl != nil {
		if schdl.actnmde == schdlactninit && schdl.initstart {
			schdl.initstart = false
			if actnsl := schdl.initactns.Length(); actnsl > 0 {
				schdl.initactns.Iterate(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink, schdl.disposeLink, nil, nil)
			}
			if actnsl := schdl.initactns.Length(); actnsl == 0 {
				schdl.actnmde = schdlactnmain
			}
		}
	}
	return
}

func (schdl *Schedule) Stop() (err error) {

	return
}

func (schdl *Schedule) Shutdown() (err error) {
	if schdl != nil {
		if schdl.schdls != nil {
			if schdls, _ := schdl.schdls.(*Schedules); schdls != nil {
				schdls.removeSchedule(schdl)
			}
		}
	}
	return
}
