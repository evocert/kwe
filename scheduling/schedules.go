package scheduling

import (
	"strings"
	"sync"
)

//SchedulesHandler - interface
type SchedulesHandler interface {
	NewSchedule(*Schedule, ...interface{}) ScheduleHandler
	Schedules() *Schedules
}

//Schedules - struct
type Schedules struct {
	schdls      map[string]*Schedule
	schdlshndlr SchedulesHandler
	lck         *sync.Mutex
}

//NewSchedules instance
func NewSchedules(schdlshndlr SchedulesHandler) (schdls *Schedules) {
	schdls = &Schedules{schdlshndlr: schdlshndlr, schdls: map[string]*Schedule{}, lck: &sync.Mutex{}}
	return
}

//Get - Scheduler by schdlname
func (schdls *Schedules) Get(schdlname string) (schdl *Schedule) {
	if schdlname != "" {
		schdl, _ = schdls.schdls[schdlname]
	}
	return
}

//RegisterSchedule - If schedule  with same name do not exists
// will the schedule be registered
func (schdls *Schedules) RegisterSchedule(schdlname string, a ...interface{}) (schdl *Schedule) {
	if schdls != nil {
		if schdlname = strings.TrimSpace(schdlname); schdlname != "" {
			func() {
				defer schdls.lck.Unlock()
				schdls.lck.Lock()
				if _, schdlok := schdls.schdls[schdlname]; !schdlok {
					schdl = newSchedule(schdls, a...)
					if schdls.schdlshndlr != nil {
						schdl.schdlhndlr = schdls.schdlshndlr.NewSchedule(schdl, a...)
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
					schdls.schdls[schdlname] = schdl
					schdl.schdlid = schdlname
				}
			}()
		}
	}
	return
}

func (schdls *Schedules) removeSchedule(schdl *Schedule) {
	if schdls != nil && schdl != nil {
		func() {
			defer schdls.lck.Unlock()
			schdls.lck.Lock()
			if _, schdlok := schdls.schdls[schdl.schdlid]; schdlok {
				schdls.schdls[schdl.schdlid] = nil
				delete(schdls.schdls, schdl.schdlid)
			}
		}()
	}
}

var schdls *Schedules

//GLOBALSCHEDULES - Global *Schedules instance
func GLOBALSCHEDULES() *Schedules {
	return schdls
}

func init() {
	schdls = NewSchedules(nil)
}
