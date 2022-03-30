package channeling

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/channeling/channelingapi"
)

type sessionschedule struct {
	ssn          channelingapi.SessionAPI
	initstart    bool
	prcintrvl    int64
	intrvl       time.Duration
	running      bool
	wg           *sync.WaitGroup
	Milliseconds int64
	Seconds      int64
	Minutes      int64
	Hours        int64
	From         time.Time
	To           time.Time
	Once         bool
	InitPath     string
	StartArgs    []interface{}
	OnError      func(*sessionschedule, error)
	OnStart      func(a ...interface{}) error
	StopArgs     []interface{}
	OnStop       func(a ...interface{}) error
	OnShutdown   func() error
}

func (ssnschdl *sessionschedule) Start() {
	if ssnschdl != nil {

	}
}

func (ssnschdl *sessionschedule) Stop() {
	if ssnschdl != nil {

	}
}

func (ssnschdl *sessionschedule) Shutdown() (err error) {
	if ssnschdl != nil {

	}
	return
}

func ticking(schdl *sessionschedule) {
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
							/*if schdl.OnError != nil {
								schdl.OnError(schdl.schdls, schdl, errprcng)
							}*/
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

func process(schdl *sessionschedule) (done bool, err error) {
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

func execute(schdl *sessionschedule) (done bool, err error) {
	if schdl != nil {
		/*var nextactns bool = false
		if nextactns = (schdl.actnmde == schdlactninit && schdl.initstart); nextactns {
			nextactns, err = executeInit(schdl)
		}
		if (!nextactns || nextactns) && (schdl.actnmde == schdlactnmain) {
			done, err = executeMain(schdl)
		}*/
	}
	return
}

var ssnschdlspool = &sync.Pool{
	New: func() interface{} {
		return newSsnSchdl()
	},
}

func newSsnSchdl(a ...interface{}) (ssnschdl *sessionschedule) {
	ssnschdl = &sessionschedule{}
	setupSsnSchedule(ssnschdl, nil, a...)
	return
}

func setupSsnSchedule(ssnschdl *sessionschedule, ssn channelingapi.SessionAPI, a ...interface{}) {
	if ssnschdl != nil {
		var initPath string = ""
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
				if d != nil {
					if dssn, _ := d.(channelingapi.SessionAPI); dssn != nil {
						if ssn == nil {
							ssn = dssn
						}
					} else if initpathd, _ := d.(string); initpathd != "" {
						if initPath == "" {
							initPath = initpathd
						}
					} else if dmp, _ := d.(map[string]interface{}); dmp != nil {
						for stngk, stngv := range dmp {
							if strings.ToLower(stngk) == "milliseconds" && hours == 0 && seconds == 0 && minutes == 0 {
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
				ai++
			}
		}
		if milliseconds > 0 || seconds > 0 || hours > 0 {
			ssnschdl.ssn = ssn
			ssnschdl.Milliseconds = milliseconds
			ssnschdl.Seconds = seconds
			ssnschdl.Hours = hours
			ssnschdl.From = frm
			ssnschdl.To = to
			ssnschdl.Once = once
			ssnschdl.InitPath = initPath
		}
	}
}

func ScheduleSession(a ...interface{}) (ssnschdl *sessionschedule, err error) {
	if ssn := NewSession(a...); ssn != nil {
		if ssnschdl, _ := ssnschdlspool.Get().(*sessionschedule); ssnschdl != nil {
			setupSsnSchedule(ssnschdl, ssn, a...)
		} else {
			ssnschdl = newSsnSchdl(a...)
		}
		ssnschdl.ssn = ssn
	}
	return
}
