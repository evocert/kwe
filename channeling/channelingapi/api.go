package channelingapi

import (
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
	"github.com/evocert/kwe/security"
)

type SessionAPI interface {
	RegisterInterval(string, int64, ...string) bool
	CheckInterval(string) bool
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
	CAS() *security.CAS
	FSUTILS() fsutils.FSUtils
	Caching() caching.MapAPI
	Close() error
	Execute(...interface{}) error
	Bind(nxtpth ...string) error
	Join(nxtpth ...string) error
	Faf(nxtpth ...string) error
	//FafJoin(nxtpth ...string) error
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
	QueryString() string
	Parameters() parameters.ParametersAPI
}

//var FAFExecute func(ssn SessionAPI, a ...interface{}) (err error) = nil
