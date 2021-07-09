package ext

import (
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
)

type SchedulesAPI interface {
	Register(string, ...interface{}) error
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

type Schedules struct {
	schdls    map[string]*Schedule
	schdlsref map[*Schedule]string
	schdlslck *sync.RWMutex
}

func NewSchedules() (schdls *Schedules) {
	schdls = &Schedules{schdlslck: &sync.RWMutex{}, schdls: map[string]*Schedule{}, schdlsref: map[*Schedule]string{}}
	return
}

func (schdls *Schedules) Register(schdlid string, a ...interface{}) (err error) {
	if schdls != nil {
		if schdlid = strings.TrimSpace(schdlid); schdlid != "" {
			var schdl *Schedule = nil
			func() {
				schdls.schdlslck.RLock()
				defer schdls.schdlslck.RUnlock()
				schdl = schdls.schdls[schdlid]
			}()
			func() {
				if schdl == nil {
					schdls.schdlslck.Lock()
					defer schdls.schdlslck.Unlock()
					a = append([]interface{}{schdls}, a...)
					if schdl = NewSchedule(a...); schdl != nil {
						schdls.schdls[schdlid] = schdl
						schdls.schdlsref[schdl] = schdlid
						schdl.schdlid = schdlid
					}
				}
			}()
		}
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
