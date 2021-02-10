package scheduling

import (
	"strings"
)

//SchedulesHandler - interface
type SchedulesHandler interface {
	NewSchedule(...interface{}) ScheduleHandler
}

//Schedules - struct
type Schedules struct {
	schdls      map[string]*Schedule
	schdlshndlr SchedulesHandler
}

//NewSchedules instance
func NewSchedules(schdlshndlr SchedulesHandler) (schdls *Schedules) {
	schdls = &Schedules{schdlshndlr: schdlshndlr, schdls: map[string]*Schedule{}}
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
			if _, schdlok := schdls.schdls[schdlname]; !schdlok {
				if schdls.schdlshndlr != nil {
					schdl = newSchedule(schdls, schdls.schdlshndlr.NewSchedule(a...), a...)
				} else {
					schdl = newSchedule(schdls, nil, a...)
				}
				schdls.schdls[schdlname] = schdl
			}
		}
	}
	return
}

var schdls *Schedules

//GLOBALSCHEDULES - Global *Schedules instance
func GLOBALSCHEDULES() *Schedules {
	return schdls
}

func init() {
	schdls = NewSchedules(nil)
}
