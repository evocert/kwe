package service

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evocert/kwe/env"
	"github.com/evocert/kwe/requesting"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/serving"
)

//Service Service
type Service struct {
	*serving.Service
	brkrfnc func(exenme string, exealias string, args ...string)
}

//NewService NewService
func NewService(name string, displayName string, description string, brokerfunc ...interface{}) (nwesrvs *Service, err error) {
	nwesrvs = &Service{}
	var srv, svrerr = serving.NewService(name, displayName, description, func(srvs *serving.Service, args ...string) {
		nwesrvs.startService(args...)
	}, func(srvs *serving.Service, args ...string) {
		nwesrvs.runService(args...)
	}, func(srvs *serving.Service, args ...string) {
		nwesrvs.stopService(args...)
	})
	glblenv := env.Env()
	glblenv.Set("APP-NAME", srv.ServiceName())
	glblenv.Set("APP-DISPLAY-NAME", srv.ServiceDisplayName())
	glblenv.Set("APP-DESCRIPTION", srv.ServiceDescription())
	if len(brokerfunc) == 1 {
		if brfnc, brfcnok := brokerfunc[0].(func(exenme string, exealias string, args ...string)); brfcnok {
			nwesrvs.brkrfnc = brfnc
		}
	}
	if svrerr == nil {
		nwesrvs.Service = srv
	} else {
		err = svrerr
		nwesrvs = nil
	}
	return
}

func (srvs *Service) startService(args ...string) {
	var defaultroot = "./"
	if srvs.IsService() {
		defaultroot = strings.Replace(srvs.ServiceExeFolder(), "\\", "/", -1)
	}
	resources.GLOBALRSNG().RegisterEndpoint("/", defaultroot)
	var out io.Writer = nil
	var in io.Reader = nil
	var conflabel = "conf"
	if srvs.IsConsole() || srvs.IsBroker() {
		out = os.Stdout
		if srvs.IsBroker() {
			in = os.Stdin
			conflabel = "broker"
		}
	}

	if ServeRequest != nil {
		func() {
			rqst := requesting.NewRequest(nil, "/active:"+srvs.ServiceName()+"."+conflabel+".js", in, out)
			if rqst != nil {
				func() {
					defer rqst.Close()
					ServeRequest(rqst)
				}()
			}
		}()

		func() {
			rqst := requesting.NewRequest(nil, "/active:"+srvs.ServiceName()+".init.js", in, out)
			if rqst != nil {
				func() {
					defer rqst.Close()
					ServeRequest(rqst)
				}()
			}
		}()
	}
}

func (srvs *Service) runService(args ...string) {
	if srvs.IsConsole() {
		cancelChan := make(chan os.Signal, 2)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
		env.WrapupCall(func() {
			cancelChan <- syscall.SIGTERM
			cancelChan <- syscall.SIGINT
		})
		<-cancelChan
	} else if srvs.IsBroker() {
		if srvs.brkrfnc != nil {
			srvs.brkrfnc(srvs.ServiceExeName(), srvs.ServiceName(), args...)
		}
	}
	if ServeRequest != nil {
		var out io.Writer = nil
		var in io.Reader = nil
		if srvs.IsConsole() || srvs.IsBroker() {
			out = os.Stdout
			if srvs.IsBroker() {
				in = os.Stdin
			}
		}

		func() {
			rqst := requesting.NewRequest(nil, "/active:"+srvs.ServiceName()+".finit.js", in, out)
			if rqst != nil {
				defer rqst.Close()
				ServeRequest(rqst)
			}
		}()
	}
}

func (srvs *Service) stopService(args ...string) {
	var out io.Writer = nil
	var in io.Reader = nil
	var conflabel = "final"
	if srvs.IsConsole() || srvs.IsBroker() {
		out = os.Stdout
		if srvs.IsBroker() {
			in = os.Stdin
			conflabel = "broker.final"
		}
	}
	if srvs.IsService() {
		if ServeRequest != nil {

			rqst := requesting.NewRequest(nil, "/active:"+srvs.ServiceName()+"."+conflabel+".js", in, out)
			if rqst != nil {
				defer rqst.Close()
				ServeRequest(rqst)
			}
		}
		env.ShutdownEnvironment()
	}
}

//RunService - startup Service pasing args...string
func RunService(args ...string) {
	if len(args) == 0 {
		args = os.Args
	}
	var srvs, err = NewService("", "", "", RunBroker)
	if err == nil {
		err = srvs.Execute(args...)
	}
	if err != nil {
		println(err)
	}
}

//RunBroker - RunBroker command as request in global channel
func RunBroker(exename string, exealias string, args ...string) {
	if ServeRequest != nil {
		rqst := requesting.NewRequest(nil, "/", os.Stdin, os.Stdout)
		if rqst != nil {
			defer rqst.Close()
			ServeRequest(rqst)
		}
	}
	//chnls.GLOBALCHNL().ServeReaderWriter("/", os.Stdout, os.Stdin)
}

var ServeRequest func(rqst requesting.RequestAPI, a ...interface{}) (err error) = nil
