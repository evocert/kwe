package api

import (
	"io"

	"github.com/evocert/kwe/channeling/channelingapi"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
)

//Sessioning

//Scheduling
type ScheduleAPI interface {
	Schedules() SchedulesAPI
	Session() channelingapi.SessionAPI
	Start(...interface{}) error
	AddAction(...interface{}) error
	AddInitAction(...interface{}) error
	AddWrapupAction(...interface{}) error
	Stop() error
	Shutdown() error
	Active(...interface{}) *active.Active
}

type SchedulesAPI interface {
	//Handler() SchedulesHandler
	Register(string, ...interface{}) (ScheduleAPI, error)
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
	ActiveSCHEDULING(active.Runtime) ActiveSchedulesAPI
}

type ActiveSchedulesAPI interface {
	Register(string, ...interface{}) (ScheduleAPI, error)
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
	Dispose()
}

var FAFExecute func(ssn channelingapi.SessionAPI, a ...interface{}) (err error) = nil
