package main

import (
	"os"

	"github.com/evocert/kwe/api"
	"github.com/evocert/kwe/channeling"

	"github.com/evocert/kwe/listen"

	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/requesting"
	_ "github.com/evocert/kwe/requirejs/html"

	"github.com/evocert/kwe/service"	
	_ "github.com/evocert/kwe/webactions"
)

func main() {
	var serveRequest func(a ...interface{}) error = nil
	var lstnr listen.ListenerAPI = nil
	var mqttmngr mqtt.MQTTManagerAPI = nil
	lstnr = listen.NewListener(func(ra requesting.RequestAPI) error {
		return serveRequest(ra, lstnr, mqttmngr)
	})

	serveRequest = func(a ...interface{}) (err error) {
		channeling.ExecuteSession(nil, a...)
		return
	}

	mqttmngr = mqtt.NewMQTTManager(mqtt.MqttEventing(func(event mqtt.MqttEvent) {
		if rqst := requesting.NewRequest(nil, event.EventPath()); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, event, mqttmngr, lstnr)
			}
		}
	}), mqtt.MqttMessaging(func(message mqtt.Message) {
		if rqst := requesting.NewRequest(nil, message.TopicPath()); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, message, mqttmngr, lstnr)
			}
		}
	}))

	api.FAFExecute = func(ssn api.SessionAPI, a ...interface{}) (err error) {
		if rqst := requesting.NewRequest(nil, a...); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, lstnr, mqttmngr)
			}
		}
		return
	}

	service.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) error {
		return serveRequest(rqst, lstnr, mqttmngr)
	}
	service.RunService(os.Args...)
}
