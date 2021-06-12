package main

import (
	"os"

	//_ "github.com/evocert/kwe/database/db2"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	/*web := web.NewClient()

	soapsend := `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<CreateJustineVideoLink xmlns="http://tempuri.org/">
		  <p>w5e7y1h8s9</p>
		  <a>5659678</a>
		  <s>Bronze</s>
		  <u>https://www.youtube.com/embed/oTOZ0mFsBOk</u>
		</CreateJustineVideoLink>
	  </soap:Body>
	</soap:Envelope>`

	if rdr, err := web.Send("http://api.msgl.ink/smsportalws.asmx", map[string]string{
		"Host":           "api.msgl.ink",
		"Content-Type":   "text/xml; charset=utf-8",
		"Content-Length": fmt.Sprintf("%d", len(soapsend)),
		"SOAPAction":     "http://tempuri.org/CreateJustineVideoLink"}, soapsend); err == nil {
		iorw.Fprintln(os.Stdout, rdr)
	} else {
		fmt.Println(err.Error())
	}*/

	/*mqtt.GLOBALMQTTMANAGER().ActivateTopic("controls/test")

	for i := 0; i < 2; i++ {
		mqttclid := fmt.Sprintf("mqtt%d", (i + 1))
		mqtt.GLOBALMQTTMANAGER().RegisterConnection(mqttclid, "broker", "skullquake.dedicated.co.za", "port", 1883, "user", "emqx", "password", "public", "autoack", true)

		if err := mqtt.GLOBALMQTTMANAGER().Connect(mqttclid); err != nil {
			fmt.Println(err.Error())
		}
	}

	for i := 0; i < 1; i++ {
		mqttclid := fmt.Sprintf("mqtt%d", (i + 1))

		if err := mqtt.GLOBALMQTTMANAGER().Subscribe(mqttclid, "controls/test", 0); err != nil {
			fmt.Println(err.Error())
		}
	}

	go func() {
		time.Sleep(10 * time.Second)
		//for i := 0; i < 1; i++ {
		mqttclid := fmt.Sprintf("mqtt%d", (2))
		if err := mqtt.GLOBALMQTTMANAGER().Publish(mqttclid, "controls/test", 0, false, "hello there"); err != nil {
			fmt.Println(err.Error())
		}
		//}
	}()*/
	service.RunService(os.Args...)

}
