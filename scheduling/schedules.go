package scheduling

import (
	"strings"
)

//Schedules - struct
type Schedules struct {
	schdls map[string]*Schedule
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
				schdl = newSchedule(schdls, a...)
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
	schdls = &Schedules{schdls: map[string]*Schedule{}}
}
