package scheduling

import (
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
)

//SchedulesAPI - interface
type SchedulesAPI interface {
	NewSchedule(*Schedule, ...interface{}) ScheduleAPI
	Schedules() *Schedules
}

//Schedules - struct
type Schedules struct {
	schdls      map[string]*Schedule
	schdlshndlr SchedulesAPI
	lck         *sync.RWMutex
}

//NewSchedules instance
func NewSchedules(schdlshndlr SchedulesAPI) (schdls *Schedules) {
	schdls = &Schedules{schdlshndlr: schdlshndlr, schdls: map[string]*Schedule{}, lck: &sync.RWMutex{}}
	return
}

//InOutS - OO{ in io.Reader -> out string } loop till no input
func (schdls *Schedules) InOutS(in interface{}, ioargs ...interface{}) (out string, err error) {
	var buff = iorw.NewBuffer()
	defer buff.Close()
	err = schdls.InOut(in, buff)
	out = buff.String()
	return
}

func (schdls *Schedules) inMapOut(mpin map[string]interface{}, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if mpl := len(mpin); mpl > 0 {
		if out != nil {
			hasoutput = true
			iorw.Fprint(out, "{")
		}
		var dfltalias string = ""
		var dfltschdl, crntschdl *Schedule = nil, nil
		if aliasv, aliasok := mpin["alias"]; aliasok {
			if dfltalias == "" {
				if aliass, aliassok := aliasv.(string); aliassok {
					dfltalias = aliass
					if dfcnok, dfschdl := schdls.ScheduleExists(aliass); dfcnok {
						dfltschdl = dfschdl
					}
				}
			}
			mpl--
			delete(mpin, "alias")
		}
		for mk, mv := range mpin {
			mpl--
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "\""+mk+"\":")
			}
			if mvp, mvpok := mv.(map[string]interface{}); mvpok {
				crntschdl = nil
				if dalias, daliasok := mvp["alias"]; daliasok && dalias != nil {
					if salias, saliasok := dalias.(string); saliasok && salias != "" {
						if shok, sh := schdls.ScheduleExists(salias); shok {
							crntschdl = sh
						}
					}
					if crntschdl == nil {
						if out != nil {
							hasoutput = true
							iorw.Fprint(out, "{\"error\":\"alias does not exist\"}")
						}
					}
				} else {
					if crntschdl == nil && dfltschdl != nil {
						crntschdl = dfltschdl
					} else {
						if dfltalias != "" {
							if out != nil {
								if out != nil {
									hasoutput = true
									iorw.Fprint(out, "{\"error\":\"default alias does not exist\"}")
								}
							}
						} else {
							if out != nil {
								if out != nil {
									hasoutput = true
									iorw.Fprint(out, "{\"error\":\no alias\"}")
								}
							}
						}
					}
				}

				if crntschdl != nil {
					hasoutput, err = crntschdl.inMapOut(mvp, out)
				}
			} else {
				if out != nil {
					if out != nil {
						hasoutput = true

						iorw.Fprint(out, "{\"error\":\"invalid requestt\"}")
					}
				}
			}
			if mpl >= 1 {
				if out != nil {
					hasoutput = true
					iorw.Fprint(out, ",")
				}
			}
		}
		if out != nil {
			hasoutput = true
			iorw.Fprint(out, "}")
		}
	}
	return
}

func (schdls *Schedules) inReaderOut(ri io.Reader, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if ri != nil {
		func() {
			var buff = iorw.NewBuffer()
			defer buff.Close()
			buffl, bufferr := io.Copy(buff, ri)
			if bufferr == nil || bufferr == io.EOF {
				if buffl > 0 {
					func() {
						var buffr = buff.Reader()
						defer func() {
							buffr.Close()
						}()
						d := json.NewDecoder(buffr)
						rqstmp := map[string]interface{}{}
						if jsnerr := d.Decode(&rqstmp); jsnerr == nil {
							if len(rqstmp) > 0 {
								hasoutput, err = schdls.inMapOut(rqstmp, out, ioargs...)
							}
						} else {
							err = jsnerr
						}
					}()
				}
			}
		}()
	}
	return
}

//InOut - OO{ in io.Reader -> out io.Writer } loop till no input
func (schdls *Schedules) InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error) {
	if in != nil {
		var hasoutput = false
		if mp, mpok := in.(map[string]interface{}); mpok {
			hasoutput, err = schdls.inMapOut(mp, out, ioargs...)
		} else if ri, riok := in.(io.Reader); riok && ri != nil {
			hasoutput, err = schdls.inReaderOut(ri, out, ioargs...)
		} else if si, siok := in.(string); siok && si != "" {
			hasoutput, err = schdls.inReaderOut(strings.NewReader(si), out, ioargs...)
		}
		if !hasoutput {
			if out != nil {
				if err != nil {
					iorw.Fprint(out, "{\"error\":\""+err.Error()+"\"}")
				} else {
					iorw.Fprint(out, "{}")
				}
			}
		}
	} else {
		if out != nil {
			iorw.Fprint(out, "{}")
		}
	}
	return
}

//Get - Scheduler by schdlname
func (schdls *Schedules) Get(schdlname string) (schdl *Schedule) {
	if schdls != nil && schdlname != "" {
		func() {
			schdls.lck.RLock()
			defer schdls.lck.RUnlock()
			schdl = schdls.schdls[schdlname]
		}()
	}
	return
}

//Exists - return true if Scheduler by schdlname exists
func (schdls *Schedules) Exists(schdlname string) (exist bool) {
	if schdls != nil && schdlname != "" {
		func() {
			schdls.lck.RLock()
			defer schdls.lck.RUnlock()
			_, exist = schdls.schdls[schdlname]
		}()
	}
	return
}

func (schdls *Schedules) UnregisterSchedule(schdlname string, a ...interface{}) {
	if schdls != nil {
		if schdlname = strings.TrimSpace(schdlname); schdlname != "" {
			func() {
				if schdl := schdls.Get(schdlname); schdl != nil {
					schdl.Shutdown()
					schdl = nil
				}
			}()
		}
	}
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
	return
}

//scheduleExists returns true if *Schedule with schdlid string exists
func (schdls *Schedules) ScheduleExists(scdhkid string) (schdlexist bool, schdl *Schedule) {
	if schdls != nil {
		schdl = schdls.Get(scdhkid)
		schdlexist = schdl != nil
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

//Schedules return []*Schedule of schdlid(s) ..string

func (schdls *Schedules) Schedules(schdlids ...string) (scls []*Schedule) {
	if len(schdls.schdls) > 0 {
		func() {
			schdls.lck.RLock()
			defer schdls.lck.RUnlock()
			if len(schdlids) > 0 {
				for _, schdlid := range schdlids {
					if schdl, schdlok := schdls.schdls[schdlid]; schdlok {
						if scls == nil {
							scls = []*Schedule{}
						}
						scls = append(scls, schdl)
					}
				}
			} else {
				for _, schdl := range schdls.schdls {
					if scls == nil {
						scls = []*Schedule{}
					}
					scls = append(scls, schdl)
				}
			}
		}()
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
