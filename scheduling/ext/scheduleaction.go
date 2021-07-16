package ext

type ScheduleActionAPI interface {
	Schedule() ScheduleAPI
}

type ScheduleAction struct {
	schdl ScheduleAPI
}

func (schdlactn *ScheduleAction) Schedule() (schdl ScheduleAPI) {
	if schdlactn != nil {
		schdl = schdlactn.schdl
	}
	return
}
