package ext

type ScheduleActionAPI interface {
	Schedule() ScheduleAPI
}

type ScheduleAction struct {
	schdl ScheduleAPI
}

type schdlaction struct {
	crntschdlactn ScheduleActionAPI
	schdl         *Schedule
	prvactn       *schdlaction
	nxtactn       *schdlaction
	args          []interface{}
	actn          func(...interface{}) error
	valid         bool
}

func (schdlactn *ScheduleAction) Schedule() (schdl ScheduleAPI) {
	if schdlactn != nil {
		schdl = schdlactn.schdl
	}
	return
}
