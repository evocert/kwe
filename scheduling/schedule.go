package scheduling

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

//Schedule - struct
type Schedule struct {
	schdlid    string
	once       bool
	schdlrs    *Schedules
	frstactn   *schdlaction
	lstactn    *schdlaction
	actnslck   *sync.Mutex
	StartArgs  []interface{}
	OnStart    func(a ...interface{}) error
	StopArgs   []interface{}
	OnStop     func(a ...interface{}) error
	OnShutdown func() error
	Seconds    int64
	intrvl     time.Duration
	running    bool
	wg         *sync.WaitGroup
}

func newSchedule(schdlrs *Schedules, a ...interface{}) (schdl *Schedule) {
	var start func(a ...interface{}) error = nil
	var startargs []interface{} = nil
	var stop func(a ...interface{}) error = nil
	var stopargs []interface{} = nil
	var shutdown func() error = nil
	var seconds int64 = 10
	var once = false
	if len(a) == 1 {
		if dmp, dmpok := a[0].(map[string]interface{}); dmpok {
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
				} else if strings.ToLower(stngk) == "seconds" {
					seconds, _ = stngv.(int64)
				} else if strings.ToLower(stngk) == "once" {
					once, _ = stngv.(bool)
				}
			}
		}
	}
	schdl = &Schedule{schdlrs: schdlrs,
		wg:         &sync.WaitGroup{},
		intrvl:     time.Second * 5,
		frstactn:   nil,
		lstactn:    nil,
		once:       once,
		actnslck:   &sync.Mutex{},
		OnStart:    start,
		StartArgs:  startargs,
		OnStop:     stop,
		StopArgs:   stopargs,
		OnShutdown: shutdown,
		Seconds:    seconds,
	}
	return
}

//AddAction - add action to *Schedule
func (schdl *Schedule) AddAction(a ...interface{}) (err error) {
	var lstargs []interface{} = nil
	var lstactn func(...interface{}) error = nil
	var al = 0
	var vldactions = []*schdlaction{}
	var cactn func(...interface{}) error = nil
	for {
		if al = len(a); al > 0 {
			d := a[0]
			a = a[1:]
			if args, argsok := d.([]interface{}); argsok {
				if al > 1 {
					d = a[0]
					if atcne, actneok := d.(func(...interface{}) error); actneok {
						vldactions = append(vldactions, newSchdlAction(schdl, atcne, args...))
						a = a[1:]
						lstargs = nil
					} else if actn, actnok := d.(func(...interface{})); actnok {
						vldactions = append(vldactions, newSchdlAction(schdl, func(fna ...interface{}) (fnerr error) {
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
						vldactions = append(vldactions, newSchdlAction(schdl, lstactn, args...))
					}
					break
				}
			} else {
				if cactn != nil {
					cactn = nil
				}
				if actnae, actnaeok := d.(func(...interface{}) error); actnaeok {
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
							vldactions = append(vldactions, newSchdlAction(schdl, cactn, lstargs...))
							lstargs = nil
						} else {
							d = a[0]
							if args, argsok := d.([]interface{}); argsok {
								vldactions = append(vldactions, newSchdlAction(schdl, cactn, args...))
								a = a[1:]
							} else {
								lstactn = cactn
								a = a[1:]
							}
						}
					} else {
						vldactions = append(vldactions, newSchdlAction(schdl, cactn, lstargs...))
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
		addactns(schdl, vldactions...)
	}

	return
}

func (schdl *Schedule) Execute() {
	if schdl != nil {
		actn := schdl.frstactn
		for actn != nil {
			if actn.valid {
				if actnerr := actn.actn(actn.args...); actnerr != nil {
					if strings.ToLower(actnerr.Error()) == "done" {
						actn.valid = false
						removeactn(schdl, actn)
					}
				}
			}
			actn = actn.nxtactn
		}
	}
}

//Start - Schedule
func (schdl *Schedule) Start() (err error) {
	if schdl != nil {
		if !schdl.running {
			schdl.wg.Add(1)
			go func() {
				defer schdl.wg.Done()
				if schdl.OnStart != nil {
					err = schdl.OnStart(schdl.StartArgs...)
				}
				if err == nil {

					schdl.running = true
				}
			}()
			schdl.wg.Wait()
			if schdl.running {
				schdl.wg.Add(1)
				go schdl.ticking()
				schdl.wg.Wait()
			}
		}
	}
	return
}

func (schdl *Schedule) ticking() {
	schdl.wg.Done()

	for schdl.running {
		time.Sleep(schdl.intrvl)
		schdl.process()
	}
	schdl.wg.Done()
}

func (schdl *Schedule) process() {
	if schdl != nil {
		schdl.Execute()
	}
}

//Stop - Schedule
func (schdl *Schedule) Stop() (err error) {
	if schdl.running {
		schdl.wg.Add(1)
		schdl.running = false
		schdl.wg.Wait()
		if schdl.OnStop != nil {
			err = schdl.OnStop(schdl.StopArgs...)
		}
	}
	return
}

//Shutdown - Schedule
//after this Schedule is destroyed adn not accessable anymore
func (schdl *Schedule) Shutdown() (err error) {
	if schdl != nil {
		err = schdl.Stop()
		if schdl.OnShutdown != nil {
			err = schdl.OnShutdown()
		}
		if schdl.schdlrs != nil {
			if _, schdlok := schdl.schdlrs.schdls[schdl.schdlid]; schdlok {
				schdl.schdlrs.schdls[schdl.schdlid] = nil
				delete(schdl.schdlrs.schdls, schdl.schdlid)
			}
			schdl.schdlrs = nil
		}
		schdl = nil
	}
	return
}

func addactns(schdl *Schedule, schdlactns ...*schdlaction) {
	for len(schdlactns) > 0 {

		addactn(schdl, schdlactns[0])
		schdlactns = schdlactns[1:]
	}
}

func addactn(schdl *Schedule, schdlactn *schdlaction) {
	if schdl != nil {
		if schdlactn != nil {
			if schdl.frstactn == nil && schdl.lstactn == nil {
				schdl.frstactn = schdlactn
				schdl.lstactn = schdlactn
			} else if schdl.frstactn != nil && schdl.lstactn != nil {
				schdlactn.prvactn = schdl.lstactn
				schdl.lstactn.nxtactn = schdlactn
				schdl.lstactn = schdlactn
			}
		}
	}
}

func removeactns(schdl *Schedule, schdlactns ...*schdlaction) {
	for len(schdlactns) > 0 {
		removeactn(schdl, schdlactns[0])
		schdlactns = schdlactns[1:]
	}
}

func removeactn(schdl *Schedule, schdlactn *schdlaction) {
	if schdlactn != nil {
		if schdl != nil {

		}
	}
}

type schdlaction struct {
	schdl   *Schedule
	prvactn *schdlaction
	nxtactn *schdlaction
	args    []interface{}
	actn    func(...interface{}) error
	valid   bool
}

func newSchdlAction(schdl *Schedule, actn func(...interface{}) error, a ...interface{}) (scdhlactn *schdlaction) {
	scdhlactn = &schdlaction{schdl: schdl, prvactn: nil, nxtactn: nil, actn: actn, args: a, valid: true}
	return
}

func (scdhlctn *schdlaction) dispose() {
	if scdhlctn != nil {
		if scdhlctn.schdl != nil {
			removeactn(scdhlctn.schdl, scdhlctn)
		}
		scdhlctn = nil
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
