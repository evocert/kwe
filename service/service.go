package service

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	//runtime "runtime"

	//"github.com/efjoubert/lnksys/network"
	//"github.com/efjoubert/lnksys/network"
	"github.com/evocert/kwe/chnls"
	"github.com/evocert/kwe/env"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/serving"
)

//LnkService LnkService
type LnkService struct {
	*serving.Service
	brkrfnc func(exenme string, exealias string, args ...string)
}

//NewLnkService NewLnkService
func NewLnkService(name string, displayName string, description string, brokerfunc ...interface{}) (lnksrvs *LnkService, err error) {
	lnksrvs = &LnkService{}
	var srv, svrerr = serving.NewService(name, displayName, description, func(srvs *serving.Service, args ...string) {
		lnksrvs.startLnkService(args...)
	}, func(srvs *serving.Service, args ...string) {
		lnksrvs.runLnkService(args...)
	}, func(srvs *serving.Service, args ...string) {
		lnksrvs.stopLnkService(args...)
	})
	if len(brokerfunc) == 1 {
		if brfnc, brfcnok := brokerfunc[0].(func(exenme string, exealias string, args ...string)); brfcnok {
			lnksrvs.brkrfnc = brfnc
		}
	}
	if svrerr == nil {
		lnksrvs.Service = srv
	} else {
		err = svrerr
		lnksrvs = nil
	}
	return
}

func (lnksrvs *LnkService) startLnkService(args ...string) {
	var defaultroot = "./"
	if lnksrvs.IsService() {
		defaultroot = strings.Replace(lnksrvs.ServiceExeFolder(), "\\", "/", -1)
	}
	//network.MapRoots("/", defaultroot, "resources/", "./resources", "apps/", "./apps")
	resources.GLOBALRSNG().RegisterEndpoint("/", defaultroot)
	var out io.Writer = nil
	var in io.Reader = nil
	var conflabel = "conf"
	if lnksrvs.IsConsole() || lnksrvs.IsBroker() {
		out = os.Stdout
		if lnksrvs.IsBroker() {
			in = os.Stdin
			conflabel = "broker"
		}
	}
	chnls.GLOBALCHNL().DefaultServeRW(out, "/"+lnksrvs.ServiceName()+"."+conflabel+".js", in)
	//network.DefaultServeHttp(nil, "GET", "/@"+lnksrvs.ServiceName()+".conf@.js", nil)
}

func (lnksrvs *LnkService) runLnkService(args ...string) {
	if lnksrvs.IsConsole() {
		cancelChan := make(chan os.Signal, 2)
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
		env.WrapupCall(func() {
			cancelChan <- syscall.SIGTERM
			cancelChan <- syscall.SIGINT
		})
		<-cancelChan
	} else if lnksrvs.IsBroker() {
		if lnksrvs.brkrfnc != nil {
			lnksrvs.brkrfnc(lnksrvs.ServiceExeName(), lnksrvs.ServiceName(), args...)
		}
	}
}

func (lnksrvs *LnkService) stopLnkService(args ...string) {
	if lnksrvs.IsService() {
		env.ShutdownEnvironment()
	}
}

//RunService - startup Service pasing args...string
func RunService(args ...string) {
	//runtimedbg.SetGCPercent(33)
	//runtime.GOMAXPROCS(runtime.NumCPU() * 10)
	if len(args) == 0 {
		args = os.Args
	}
	var lnksrvs, err = NewLnkService("", "", "", RunBroker)
	if err == nil {
		err = lnksrvs.Execute(args...)
	}
	if err != nil {
		println(err)
	}
}

//RunBroker - RunBroker command as request in global channel
func RunBroker(exename string, exealias string, args ...string) {
	chnls.GLOBALCHNL().DefaultServeRW(os.Stdout, "/", os.Stdin)
}
