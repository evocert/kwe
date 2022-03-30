//go:build !mqtt
// +build !mqtt

package main

import (
	"os"

	"github.com/evocert/kwe/api"
	"github.com/evocert/kwe/channeling"
	"github.com/evocert/kwe/channeling/channelingapi"

	"github.com/evocert/kwe/listen"

	"github.com/evocert/kwe/requesting"
	_ "github.com/evocert/kwe/requirejs/html"

	_ "github.com/evocert/kwe/imports"
	"github.com/evocert/kwe/service"
	_ "github.com/evocert/kwe/webactions"
)

func main() {
	var serveRequest func(a ...interface{}) error = nil
	var lstnr listen.ListenerAPI = nil
	lstnr = listen.NewListener(func(ra requesting.RequestAPI) error {
		return serveRequest(ra, lstnr)
	})

	serveRequest = func(a ...interface{}) (err error) {
		channeling.ExecuteSession(nil, a...)
		return
	}

	api.FAFExecute = func(ssn channelingapi.SessionAPI, a ...interface{}) (err error) {
		if rqst := requesting.NewRequest(nil, a...); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, lstnr)
			}
		}
		return
	}

	service.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) error {
		return serveRequest(rqst, lstnr)
	}
	service.RunService(os.Args...)
}
