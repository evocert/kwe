package ext

import "github.com/evocert/kwe/api"

type ScheduleActionAPI interface {
	Schedule() api.ScheduleAPI
}

type ScheduleAction struct {
	schdl api.ScheduleAPI
}

func (schdlactn *ScheduleAction) Schedule() (schdl api.ScheduleAPI) {
	if schdlactn != nil {
		schdl = schdlactn.schdl
	}
	return
}
