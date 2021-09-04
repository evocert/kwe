package ext

import (
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

type SchedulesHandler interface {
	NewSchedule(*Schedule, ...interface{}) ScheduleHandler
	Schedules() *Schedules
}

type SchedulesAPI interface {
	Handler() SchedulesHandler
	Register(string, ...interface{}) error
	Get(string) ScheduleAPI
	Unregister(string) error
	Exists(string) bool
	Start(string, ...interface{}) error
	Stop(string) error
	Ammend(string, ...interface{}) error
	Shutdown(string) error
	InOut(io.Reader, io.Writer) error
	Fprint(io.Writer) error
	Reader() iorw.Reader
	ActiveSCHEDULING(active.Runtime, ...parameters.ParametersAPI) ActiveSchedulesAPI
}

type ActiveSchedulesAPI interface {
	Register(string, ...interface{}) error
	Get(string) ScheduleAPI
	Unregister(string) error
	Exists(string) bool
	Start(string, ...interface{}) error
	Stop(string) error
	Ammend(string, ...interface{}) error
	Shutdown(string) error
	InOut(io.Reader, io.Writer) error
	Fprint(io.Writer) error
	Reader() iorw.Reader
}

type ActiveSchedules struct {
	schdls SchedulesAPI
	rntme  active.Runtime
	prms   parameters.ParametersAPI
}

func newActiveSchedules(schdls SchedulesAPI, rntme active.Runtime, prms ...parameters.ParametersAPI) (atvschdls *ActiveSchedules) {
	if rntme != nil && schdls != nil {
		atvschdls = &ActiveSchedules{rntme: rntme, schdls: schdls}
		if len(prms) > 0 && prms[0] != nil {
			atvschdls.prms = prms[0]
		}
	}
	return
}

func (atvschdls *ActiveSchedules) Register(schdlid string, a ...interface{}) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		a = append([]interface{}{atvschdls.rntme}, a...)
		err = atvschdls.schdls.Register(schdlid, a...)
	}
	return
}

func (atvschdls *ActiveSchedules) Get(schdlid string) (schdl ScheduleAPI) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		schdl = atvschdls.schdls.Get(schdlid)
	}
	return
}

func (atvschdls *ActiveSchedules) Unregister(schdlid string) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		err = atvschdls.schdls.Unregister(schdlid)
	}
	return
}

func (atvschdls *ActiveSchedules) Exists(schdlid string) (exists bool) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		exists = atvschdls.schdls.Exists(schdlid)
	}
	return
}

func (atvschdls *ActiveSchedules) Start(schdlid string, a ...interface{}) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		a = append([]interface{}{atvschdls.rntme}, a...)
		err = atvschdls.schdls.Start(schdlid, a...)
	}
	return
}

func (atvschdls *ActiveSchedules) Stop(schdlid string) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		err = atvschdls.schdls.Stop(schdlid)
	}
	return
}

func (atvschdls *ActiveSchedules) Ammend(schdlid string, a ...interface{}) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		a = append([]interface{}{atvschdls.rntme}, a...)
		err = atvschdls.schdls.Ammend(schdlid, a...)
	}
	return
}

func (atvschdls *ActiveSchedules) Shutdown(schdlid string) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		err = atvschdls.schdls.Shutdown(schdlid)
	}
	return
}

func (atvschdls *ActiveSchedules) InOut(in io.Reader, out io.Writer) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		err = atvschdls.schdls.InOut(in, out)
	}
	return
}

func (atvschdls *ActiveSchedules) Fprint(w io.Writer) (err error) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		err = atvschdls.schdls.Fprint(w)
	}
	return
}

func (atvschdls *ActiveSchedules) Reader() (rdr iorw.Reader) {
	if atvschdls != nil && atvschdls.rntme != nil && atvschdls.schdls != nil {
		rdr = atvschdls.schdls.Reader()
	}
	return
}

type Schedules struct {
	schdls      map[string]*Schedule
	schdlsref   map[*Schedule]string
	schdlslck   *sync.RWMutex
	schdlshndlr SchedulesHandler
}

func NewSchedules(schdlshndlr ...SchedulesHandler) (schdls *Schedules) {
	schdls = &Schedules{schdlslck: &sync.RWMutex{}, schdls: map[string]*Schedule{}, schdlsref: map[*Schedule]string{}}
	if len(schdlshndlr) == 1 && schdlshndlr[0] != nil {
		schdls.schdlshndlr = schdlshndlr[0]
	}
	return
}

func (schdls *Schedules) ActiveSCHEDULING(rntme active.Runtime, prms ...parameters.ParametersAPI) (atvschdlsapi ActiveSchedulesAPI) {
	atvschdlsapi = newActiveSchedules(schdls, rntme, prms...)
	return
}

func (schdls *Schedules) Handler() (schdlshndlr SchedulesHandler) {
	if schdls != nil {
		schdlshndlr = schdls.schdlshndlr
	}
	return
}

func (schdls *Schedules) Register(schdlid string, a ...interface{}) (err error) {
	if schdls != nil {
		if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
			var schdl *Schedule = nil
			if func() bool {
				schdls.schdlslck.RLock()
				defer schdls.schdlslck.RUnlock()
				schdl = schdls.schdls[schdlid]
				return schdl == nil
			}() {
				func() {
					schdls.schdlslck.Lock()
					defer schdls.schdlslck.Unlock()
					var schdlactions map[string][]interface{} = nil
					var schdlactionsok bool = false
					ai := 0
					for {
						if al := len(a); ai < al {
							d := a[ai]
							ai++
							if schdlactions, schdlactionsok = d.(map[string][]interface{}); schdlactionsok {
								a = append(a[:ai], a[ai:]...)
								ai--
							}
						} else {
							break
						}
					}
					a = append([]interface{}{schdls}, a...)
					if schdl = NewSchedule(a...); schdl != nil {
						schdls.schdls[schdlid] = schdl
						schdls.schdlsref[schdl] = schdlid
						schdl.schdlid = schdlid
						if len(schdlactions) > 0 {
							for schdlactntpe, actns := range schdlactions {
								if len(actns) > 0 {
									if schdlactntpe = strings.ToLower(schdlactntpe); schdlactntpe == "init" {
										schdl.AddInitAction(actns...)
									} else if schdlactntpe == "main" {
										schdl.AddAction(actns...)
									} else if schdlactntpe == "wrapup" {
										schdl.AddWrapupAction(actns...)
									}
								}
							}
						}
					}
				}()
			}
		}
	}
	return
}

//Get - Scheduler by schdlname
func (schdls *Schedules) Get(schdlname string) (schdl ScheduleAPI) {
	if schdls != nil && schdlname != "" {
		func() {
			schdls.schdlslck.RLock()
			defer schdls.schdlslck.RUnlock()
			schdl = schdls.schdls[schdlname]
		}()
	}
	return
}

func (schdls *Schedules) Unregister(schdlid string) (err error) {
	if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
		var schdl *Schedule = nil
		func() {
			schdls.schdlslck.RLock()
			defer schdls.schdlslck.RUnlock()
			schdl = schdls.schdls[schdlid]
		}()
		func() {
			if schdl != nil {
				schdl.Shutdown()
				schdl = nil
			}
		}()
	}
	return
}

func (schdls *Schedules) Start(schdlid string, a ...interface{}) (err error) {
	if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
		var schdl *Schedule = nil
		func() {
			schdls.schdlslck.RLock()
			defer schdls.schdlslck.RUnlock()
			schdl = schdls.schdls[schdlid]
		}()
		func() {
			if schdl != nil {
				err = schdl.Start(a...)
			}
		}()
	}
	return
}

func (schdls *Schedules) Stop(schdlid string) (err error) {
	if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
		var schdl *Schedule = nil
		func() {
			schdls.schdlslck.RLock()
			defer schdls.schdlslck.RUnlock()
			schdl = schdls.schdls[schdlid]
		}()
		func() {
			if schdl != nil {
				err = schdl.Stop()
			}
		}()
	}
	return
}

func (schdls *Schedules) Exists(schdlid string) (exist bool) {
	if schdls != nil {
		if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
			func() {
				schdls.schdlslck.RLock()
				defer schdls.schdlslck.RUnlock()
				_, exist = schdls.schdls[schdlid]
			}()
		}
	}
	return
}

func (schdls *Schedules) Shutdown(schdlid string) (err error) {
	if schdls != nil {
		if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
			var schdl *Schedule = nil
			func() {
				schdls.schdlslck.RLock()
				defer schdls.schdlslck.RUnlock()
				schdl = schdls.schdls[schdlid]
			}()
			if schdl != nil {
				err = schdl.Shutdown()
			}
		}
	}
	return
}

func (schdls *Schedules) Ammend(schdlid string, a ...interface{}) (err error) {

	return
}

func (schdls *Schedules) InOut(r io.Reader, w io.Writer) (err error) {

	return
}

func (schdls *Schedules) Fprint(w io.Writer) (err error) {
	return
}

func (schdls *Schedules) Reader() (r iorw.Reader) {

	return
}

func (schdls *Schedules) removeSchedule(schdl *Schedule) {
	if schdls != nil {
		schdls.schdlslck.Lock()
		defer schdls.schdlslck.Unlock()
		if schdlid := schdls.schdlsref[schdl]; schdlid != "" {
			delete(schdls.schdlsref, schdl)
			if schdls.schdls[schdlid] == schdl {
				schdls.schdls[schdlid] = nil
			}
			delete(schdls.schdls, schdlid)
		}
	}
}

var glblschdls *Schedules = nil

func GLOBALSCHEDULES() *Schedules {
	return glblschdls
}

func init() {
	if glblschdls == nil {
		glblschdls = NewSchedules()
	}
}
