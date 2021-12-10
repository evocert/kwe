package channeling

import (
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/api"
)

type sessionschedule struct {
	ssn          api.SessionAPI
	Milliseconds int64
	Seconds      int64
	Minutes      int64
	Hours        int64
	From         time.Time
	To           time.Time
	Once         bool
	InitPath     string
}

func (ssnschdl *sessionschedule) Start() {
	if ssnschdl != nil {

	}
}

func (ssnschdl *sessionschedule) Stop() {
	if ssnschdl != nil {

	}
}

func (ssnschdl *sessionschedule) Shutdown() {
	if ssnschdl != nil {

	}
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

func setupSsnSchedule(ssnschdl *sessionschedule, ssn api.SessionAPI, a ...interface{}) {
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
					if dssn, _ := d.(api.SessionAPI); dssn != nil {
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
