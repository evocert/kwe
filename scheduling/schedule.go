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
	From       time.Time
	To         time.Time
	schdlrs    *Schedules
	frstactn   *schdlaction
	lstactn    *schdlaction
	actnslck   *sync.Mutex
	StartArgs  []interface{}
	OnError    func(*Schedules, *Schedule, error)
	OnStart    func(a ...interface{}) error
	StopArgs   []interface{}
	OnStop     func(a ...interface{}) error
	OnShutdown func() error
	Seconds    int64
	intrvl     time.Duration
	running    bool
	prcng      bool
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
	var frm time.Time = time.Now()
	frm = time.Date(frm.Year(), frm.Month(), frm.Day(), 0, 0, 0, 0, frm.Location())
	var to time.Time = frm.Add(time.Hour * 24)
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
	}

	schdl = &Schedule{schdlrs: schdlrs,
		wg:         &sync.WaitGroup{},
		intrvl:     time.Microsecond * 500,
		frstactn:   nil,
		lstactn:    nil,
		running:    false,
		prcng:      false,
		once:       once,
		actnslck:   &sync.Mutex{},
		OnStart:    start,
		StartArgs:  startargs,
		OnStop:     stop,
		StopArgs:   stopargs,
		OnShutdown: shutdown,
		Seconds:    seconds,
		From:       frm,
		To:         to,
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

func (schdl *Schedule) execute() (err error) {
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
	return
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
	//var strtprcng, endprcng time.Time
	var errprcng error = nil
	var nxttrggrstmp, frmstmp, tostmp time.Time
	frmstmp = schdl.From
	tostmp = schdl.To
	nxttrggrstmp = frmstmp
	var intrvl time.Duration = schdl.intrvl
	//var scnds int64 = schdl.Seconds
	var calcnxttrggr = func() (cantrggr bool) {
		tmpNow := time.Now()
		tmpfrm := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), frmstmp.Hour(), frmstmp.Minute(), frmstmp.Second(), frmstmp.Nanosecond(), frmstmp.Location())
		tmpto := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tostmp.Hour(), tostmp.Minute(), tostmp.Second(), tostmp.Nanosecond(), tostmp.Location())

		if tmpNow.After(tmpfrm) && tmpNow.Before(tmpto) {
			if nxttrggrstmp.Before(tmpfrm) {
				nxttrggrstmp = time.Date(tmpfrm.Year(), tmpfrm.Month(), tmpfrm.Day(), tmpfrm.Hour(), tmpfrm.Minute(), tmpfrm.Second(), tmpfrm.Nanosecond(), tmpfrm.Location())
				cantrggr = true
			} else {
				//var tmdif = int64(tmpNow.Sub(tmpfrm))
				//var secdif = int64(time.Duration(scnds) * time.Nanosecond)

				//var nxttrggrdif = int64(nxttrggrstmp.Sub(tmpfrm))

			}

		} else {
			cantrggr = false
		}
		return
	}
	for schdl.running {
		if calcnxttrggr() {
			//if !schdl.prcng {
			//schdl.prcng = true
			_, _, errprcng = schdl.process()
			if errprcng != nil {
				if schdl.OnError != nil {
					schdl.OnError(schdl.schdlrs, schdl, errprcng)
					//endprcng = time.Now()
				}
			}
			schdl.prcng = false
			if schdl.once {
				break
			}
			//}
		} else {
			time.Sleep(intrvl)
		}
	}
	schdl.wg.Done()
}

func (schdl *Schedule) process() (strtprcng, endprcng time.Time, err error) {
	if schdl != nil {
		strtprcng = time.Now()
		func() {
			defer func() {
				if rv := recover(); rv != nil {
					err = fmt.Errorf("%v", rv)
				}
			}()
			err = schdl.execute()
		}()
		endprcng = time.Now()
	}
	return
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
