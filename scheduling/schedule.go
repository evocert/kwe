package scheduling

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
)

//ScheduleHandler - interface
type ScheduleHandler interface {
	StartedSchedule(...interface{}) error
	StoppedSchedule(...interface{}) error
	ShutdownSchedule() error
	PrepActionArgs(...interface{}) ([]interface{}, error)
	Schedule() *Schedule
}

//Schedule - struct
type Schedule struct {
	actnmde      scheduleactiontype
	initstart    bool
	schdlid      string
	once         bool
	schdlhndlr   ScheduleHandler
	From         time.Time
	To           time.Time
	schdlrs      *Schedules
	actns        *enumeration.Chain
	initactns    *enumeration.Chain
	wrapupactns  *enumeration.Chain
	actnslck     *sync.Mutex
	StartArgs    []interface{}
	OnError      func(*Schedules, *Schedule, error)
	OnStart      func(a ...interface{}) error
	StopArgs     []interface{}
	OnStop       func(a ...interface{}) error
	OnShutdown   func() error
	Milliseconds int64
	Seconds      int64
	Minutes      int64
	Hours        int64
	prcintrvl    int64
	intrvl       time.Duration
	running      bool
	wg           *sync.WaitGroup
}

func newSchedule(schdlrs *Schedules, a ...interface{}) (schdl *Schedule) {
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
	if len(a) > 1 {
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
	}

	schdl = &Schedule{
		initstart:    true,
		actnmde:      schdlactninit,
		schdlrs:      schdlrs,
		wg:           &sync.WaitGroup{},
		intrvl:       time.Microsecond * 2,
		prcintrvl:    0,
		actns:        nil,
		initactns:    nil,
		wrapupactns:  nil,
		running:      false,
		once:         once,
		actnslck:     &sync.Mutex{},
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
	}
	schdl.actns = enumeration.NewChain()
	schdl.initactns = enumeration.NewChain()
	schdl.wrapupactns = enumeration.NewChain()
	return
}

type FuncArgsErrHandle func(...interface{}) error
type FuncArgsHandle func(...interface{})
type FuncErrHandle func() error
type FuncHandle func(...interface{})

type scheduleactiontype int

const (
	schdlactnmain scheduleactiontype = iota
	schdlactninit
	schdlactnwrapup
)

//AddAction - add action(s) to *Schedule
func (schdl *Schedule) AddAction(a ...interface{}) (err error) {
	err = schdl.internalAction(schdlactnmain, a...)
	return
}

//AddInitAction - add action(s) to *Schedule that will be execute initially
func (schdl *Schedule) AddInitAction(a ...interface{}) (err error) {
	err = schdl.internalAction(schdlactninit, a...)
	return
}

//AddWrapupAction - add action(s) to *Schedule that will be execute when there are no more
// main list fo action(s) to execute
func (schdl *Schedule) AddWrapupAction(a ...interface{}) (err error) {
	err = schdl.internalAction(schdlactnwrapup, a...)
	return
}

func (schdl *Schedule) internalAction(actntpe scheduleactiontype, a ...interface{}) (err error) {
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
		addactns(schdl, actntpe, vldactions...)
	}
	return
}

func (schdl *Schedule) doLink(lnk *enumeration.Link) (done bool, err error) {
	if d := lnk.Value(); d != nil {
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
						lnk.Done()
					}
				} else {
					if schdl.actnmde == schdlactninit || schdl.actnmde == schdlactnwrapup {
						lnk.Done()
					}
				}
			}()
		}
	}
	return
}

func (schdl *Schedule) errDoLink(lnk *enumeration.Link, err error) (done bool) {
	if schdl.actnmde != schdlactnmain {
		done = true
	}
	return
}

func (schdl *Schedule) doneLink(lnk *enumeration.Link) (err error) {

	return
}

func (schdl *Schedule) errDoneLink(lnk *enumeration.Link, err error) (done bool) {
	if schdl.actnmde != schdlactnmain {
		done = true
	}
	return
}

func (schdl *Schedule) executeMain() (done bool, err error) {
	if actnsl := schdl.actns.Size(); actnsl > 0 {
		schdl.actns.Do(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink)
	}
	if schdl.actns.Size() == 0 || schdl.once {
		schdl.actnmde = schdlactnwrapup
		if actnsl := schdl.wrapupactns.Size(); actnsl > 0 {
			schdl.wrapupactns.Do(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink)
		}
		if done = (schdl.actns.Size() == 0 || schdl.once); !done {
			schdl.actnmde = schdlactnmain
		}
	}
	return
}

func (schdl *Schedule) executeInit() (nextactns bool, err error) {
	if schdl != nil {
		if schdl.actnmde == schdlactninit && schdl.initstart {
			schdl.initstart = false
			if actnsl := schdl.initactns.Size(); actnsl > 0 {
				schdl.initactns.Do(schdl.doLink, schdl.errDoLink, schdl.doneLink, schdl.errDoneLink)
			}
			if actnsl := schdl.initactns.Size(); actnsl == 0 {
				schdl.actnmde = schdlactnmain
			}
		}
	}
	return
}

func (schdl *Schedule) execute() (done bool, err error) {
	if schdl != nil {
		var nextactns bool = false
		if nextactns = (schdl.actnmde == schdlactninit && schdl.initstart); nextactns {
			nextactns, err = schdl.executeInit()
		}
		if (!nextactns || nextactns) && (schdl.actnmde == schdlactnmain) {
			done, err = schdl.executeMain()
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
						if prcngdone, errprcng = schdl.process(); errprcng != nil {
							if schdl.OnError != nil {
								schdl.OnError(schdl.schdlrs, schdl, errprcng)
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

func (schdl *Schedule) process() (done bool, err error) {
	if schdl != nil {
		func() {
			defer func() {
				if rv := recover(); rv != nil {
					err = fmt.Errorf("%v", rv)
					fmt.Println(err.Error())
				}
			}()
			done, err = schdl.execute()
		}()
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
	} else {
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
			schdl.schdlrs.removeSchedule(schdl)
			schdl.schdlrs = nil
		}
		schdl = nil
	}
	return
}

func (schdl *Schedule) inMapOut(mpin map[string]interface{}, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if schdl != nil {
		var enc *json.Encoder = nil
		if mpl := len(mpin); mpl > 0 {
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "{")
			}
			for mk, mv := range mpin {
				mpl--
				if out != nil {
					hasoutput = true
					iorw.Fprint(out, "\""+mk+"\":[")
				}
				if mvarr, mvarrok := mv.([]interface{}); mvarrok {
					if mvarrl := len(mvarr); mvarrl > 0 {
						for _, mvmvarrv := range mvarr {
							mvarrl--
							if mvp, mvpok := mvmvarrv.(map[string]interface{}); mvpok {
								if len(mvp) > 0 {
									for mk, mv := range mvp {
										if actnsargs, atcnsargsok := mv.([]interface{}); atcnsargsok {
											if len(actnsargs) > 0 {
												if actnk := strings.ToLower(mk); actnk != "" && strings.HasPrefix(actnk, "action-") && (actnk[len("action-"):] == "" || strings.Contains("|wrapup|init|main|", "|"+actnk[len("action-"):]+"|")) {
													iorw.Fprint(out, "{\""+"action-"+actnk+"\":")
													var actnerr error = nil
													if actnk = actnk[len("action-"):]; actnk == "" || actnk == "main" {
														actnerr = schdl.AddAction(actnsargs...)
													} else if actnk == "init" {
														actnerr = schdl.AddInitAction(actnsargs...)
													} else if actnk == "wrapup" {
														actnerr = schdl.AddWrapupAction(actnsargs...)
													}
													if out != nil {
														hasoutput = true
														if actnerr != nil {
															if enc == nil {
																enc = json.NewEncoder(out)
															}
															iorw.Fprint(out, "{\"error\":")
															enc.Encode(err.Error())
															iorw.Fprint(out, "}")
															err = nil
														} else {
															iorw.Fprint(out, "{}")
														}
													}
													iorw.Fprint(out, "}")
												}
											}
										}
									}
								} else {
									if out != nil {
										hasoutput = true
										//jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("no request"))
										//io.Copy(out, jsnrdr)
										//jsnrdr = nil
										iorw.Fprint(out, "{}")
									}
								}
							} else if schdlargs, schdlargsok := mvmvarrv.([]interface{}); schdlargsok && len(schdlargs) > 0 {
								scdhdlsargsl := len(schdlargs)
								for _, schdla := range schdlargs {
									if schdlas, schdlasok := schdla.(string); schdlasok && schdlas != "" {
										if schdlas = strings.ToLower(schdlas); schdlas != "" && strings.Contains("|start|stop|shutdown|", "|"+schdlas+"|") {

										}
									} else {
										if out != nil {
											hasoutput = true
											//jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("invalid request"))
											//io.Copy(out, jsnrdr)
											//jsnrdr = nil
											iorw.Fprint(out, "{\"error\":\"invalid request\"}")
										}
									}
									scdhdlsargsl--
									if scdhdlsargsl >= 1 {
										if out != nil {
											hasoutput = true
											iorw.Fprint(out, ",")
										}
									}
								}

							} else {
								if out != nil {
									hasoutput = true
									//jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("invalid request"))
									//io.Copy(out, jsnrdr)
									//jsnrdr = nil
									iorw.Fprint(out, "{\"error\":\"invalid request\"}")
								}
							}
						}
						if mvarrl >= 1 {
							if out != nil {
								hasoutput = true
								iorw.Fprint(out, ",")
							}
						}
					}
				}
				if out != nil {
					hasoutput = true
					iorw.Fprint(out, "]")
				}
				if mpl >= 1 {
					if out != nil {
						hasoutput = true
						iorw.Fprint(out, ",")
					}
				}
			}
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "}")
			}
		}
	}
	return
}

func (schdl *Schedule) inReaderOut(ri io.Reader, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if ri != nil {
		func() {
			var buff = iorw.NewBuffer()
			defer buff.Close()
			buffl, bufferr := io.Copy(buff, ri)
			if bufferr == nil || bufferr == io.EOF {
				if buffl > 0 {
					func() {
						var buffr = buff.Reader()
						defer func() {
							buffr.Close()
						}()
						d := json.NewDecoder(buffr)
						rqstmp := map[string]interface{}{}
						if jsnerr := d.Decode(&rqstmp); jsnerr == nil {
							if len(rqstmp) > 0 {
								hasoutput, err = schdl.inMapOut(rqstmp, out, ioargs...)
							}
						} else {
							err = jsnerr
						}
					}()
				}
			}
		}()
	}
	return
}

//InOut - OO{ in io.Reader -> out io.Writer } loop till no input
func (schdl *Schedule) InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error) {
	if in != nil {
		var hasoutput = false
		if mp, mpok := in.(map[string]interface{}); mpok {
			hasoutput, err = schdl.inMapOut(mp, out, ioargs...)
		} else if ri, riok := in.(io.Reader); riok && ri != nil {
			hasoutput, err = schdl.inReaderOut(ri, out, ioargs...)
		} else if si, siok := in.(string); siok && si != "" {
			hasoutput, err = schdl.inReaderOut(strings.NewReader(si), out, ioargs...)
		}
		if !hasoutput {
			if out != nil {
				if err != nil {
					iorw.Fprint(out, "{\"error\":\""+err.Error()+"\"}")
				} else {
					iorw.Fprint(out, "{}")
				}
			}
		}
	} else {
		if out != nil {
			iorw.Fprint(out, "{}")
		}
	}
	return
}

func addactns(schdl *Schedule, actntpe scheduleactiontype, schdlactns ...*schdlaction) {

	for len(schdlactns) > 0 {
		schlactn := schdlactns[0]
		addactn(schdl, actntpe, schlactn)
		schdlactns = schdlactns[1:]
	}
}

func addactn(schdl *Schedule, actntpe scheduleactiontype, schdlactn *schdlaction) {
	if schdl != nil {
		if schdlactn != nil {
			switch actntpe {
			case schdlactnmain:
				if schdl.actns != nil {
					schdl.actns.Add(schdlactn)
				}
			case schdlactninit:
				if schdl.initactns != nil {
					schdl.initactns.Add(schdlactn)
				}
			case schdlactnwrapup:
				if schdl.wrapupactns != nil {
					schdl.wrapupactns.Add(schdlactn)
				}
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
			var prvactn *schdlaction = schdlactn.prvactn
			var nxtactn *schdlaction = schdlactn.nxtactn
			if prvactn != nil {
				if nxtactn != nil {
					prvactn.nxtactn = nxtactn
				} else {
					prvactn.nxtactn = nil
				}
			}
			if nxtactn != nil {
				nxtactn.prvactn = prvactn
			}
		}
		schdlactn = nil
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
		if scdhlctn.schdl != nil {
			scdhlctn.schdl = nil
		}
		if scdhlctn.actn != nil {
			scdhlctn.actn = nil
		}
		if scdhlctn.args != nil {
			scdhlctn.args = nil
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
