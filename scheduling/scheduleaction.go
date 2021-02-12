package scheduling

//ScheduleAction - struct core implementation of ActionHandler
type ScheduleAction struct {
	schdl     *Schedule
	schdls    *Schedules
	settings  map[string]interface{}
	actnargs  []interface{}
	OnExecute func(...interface{}) (bool, error)
}

//ExecuteAction implementation of ActionHandler
func (schdlactn *ScheduleAction) ExecuteAction(a ...interface{}) (result bool, err error) {
	if schdlactn.OnExecute != nil {
		result, err = schdlactn.OnExecute(a...)
	}
	return
}

//NewScheduleAction return new ScheduleAction instance
func NewScheduleAction(schdl *Schedule, settings map[string]interface{}, args ...interface{}) (schdlactn *ScheduleAction) {
	schdlactn = &ScheduleAction{schdl: schdl, schdls: schdls, settings: map[string]interface{}{}, actnargs: args}
	return
}
