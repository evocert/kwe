package main

import (
	"os"

	_ "github.com/evocert/kwe/alertify"
	_ "github.com/evocert/kwe/babylon"
	_ "github.com/evocert/kwe/bootstrap"
	"github.com/evocert/kwe/channeling"
	_ "github.com/evocert/kwe/crypto"
	_ "github.com/evocert/kwe/datatables"
	_ "github.com/evocert/kwe/datepicker"
	_ "github.com/evocert/kwe/fonts/material"
	_ "github.com/evocert/kwe/fonts/robotov27latin"
	_ "github.com/evocert/kwe/goldenlayout"
	_ "github.com/evocert/kwe/jqueryui"
	_ "github.com/evocert/kwe/jspanel"
	"github.com/evocert/kwe/listen"
	_ "github.com/evocert/kwe/luxon"
	"github.com/evocert/kwe/mqtt"
	_ "github.com/evocert/kwe/raphael"
	"github.com/evocert/kwe/requesting"
	_ "github.com/evocert/kwe/requirejs/html"
	scheduling "github.com/evocert/kwe/scheduling/ext"
	"github.com/evocert/kwe/service"

	_ "github.com/evocert/kwe/sip"

	_ "github.com/evocert/kwe/typescript"

	_ "github.com/evocert/kwe/database/kwesqlite"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/sqlite"

	//To use ora import use go 1.6+
	//_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"

	_ "github.com/evocert/kwe/webactions"
)

func main() {
	var serveRequest func(a ...interface{}) error = nil
	var glblschdls scheduling.SchedulesAPI = scheduling.GLOBALSCHEDULES()
	var lstnr listen.ListenerAPI = nil
	var mqttmngr mqtt.MQTTManagerAPI = nil
	lstnr = listen.NewListener(func(ra requesting.RequestAPI) error {
		return serveRequest(ra, lstnr, mqttmngr, glblschdls)
	})

	serveRequest = func(a ...interface{}) (err error) {
		channeling.ExecuteSession(nil, a...)
		return
	}

	mqttmngr = mqtt.NewMQTTManager(mqtt.MqttEventing(func(event mqtt.MqttEvent) {
		if rqst := requesting.NewRequest(nil, event.EventPath()); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, event, mqttmngr, glblschdls, lstnr)
			}
		}
	}), mqtt.MqttMessaging(func(message mqtt.Message) {
		if rqst := requesting.NewRequest(nil, message.TopicPath()); rqst != nil {
			defer rqst.Close()
			if serveRequest != nil {
				serveRequest(rqst, message, mqttmngr, glblschdls, lstnr)
			}
		}
	}))

	/*gblrsngfs := resources.GLOBALRSNG().FS()
	gblrsngfs.MKDIR("/testws", "")
	gblrsngfs.SET("/testws/test.js", `kwe.out().print("hello there");`)

	go func() {
		time.Sleep(time.Second * 10)
		clnt := web.NewClient()
		if rw, rwerr := clnt.SendReceive("ws://localhost/testws/test.js"); rw != nil {
			fmt.Println(rw.ReadAll())
		} else if rwerr != nil {
			fmt.Println(rwerr)
		}

	}()*/

	service.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) error {
		return serveRequest(rqst, mqttmngr, glblschdls, lstnr)
	}
	service.RunService(os.Args...)
}
