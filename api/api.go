package api

import (
	"io"

	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/env"
	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/requesting"
)

//Sessioning

type SessionAPI interface {
	MQTTManager() mqtt.MQTTManagerAPI
	MQTTEvent() mqtt.MqttEvent
	MQTTMessage() mqtt.Message
	DBMS() database.DBMSAPI
	Parameters() parameters.ParametersAPI
	In() requesting.RequestAPI
	Out() requesting.ResponseAPI
	Send(string, ...interface{}) (iorw.Reader, error)
	SendRecieve(string, ...interface{}) (iorw.PrinterReader, error)
	SessionSend(string, ...interface{}) (iorw.Reader, error)
	SessionSendRecieve(string, ...interface{}) (iorw.PrinterReader, error)
	FS() *fsutils.FSUtils
	SessionFS() *fsutils.FSUtils
	FSUTILS() fsutils.FSUtils
	Caching() caching.MapAPI
	Close() error
	Execute(...interface{}) error
	Bind(nxtpth ...string) error
	FAFExecute(...interface{}) error
	Env() env.EnvAPI
	Listen(network string, addr ...string) (err error)
	Shutdown(addr ...string) (err error)
	UnCertifyAddr(...string)
	CertifyAddr(string, string, ...string) error
	Path() PathAPI
	Active(...interface{}) *active.Active
	Command(execpath string, execargs ...string) (cmd *osprc.Command, err error)
}

type PathAPI interface {
	Path() string
	Ext() string
	PathRoot() string
	Args() []interface{}
}

//Scheduling

type ScheduleAPI interface {
	Schedules() SchedulesAPI
	Session() SessionAPI
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

var FAFExecute func(ssn SessionAPI, a ...interface{}) (err error) = nil
