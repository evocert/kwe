resourcing.RegisterEndpoint("/gendocs","C:/GitHub/kwe/gendocs");
resourcing.RegisterEndpoint("/webcrawler","D:/projects/system/kweexamples-main/kweexamples-main/src/webcrawler");
channel.Send("/gendocs/active:gendocs.js");
resourcing.RegisterEndpoint("/","D:/projects/inovo/clients/Avon/one/system");
resourcing.RegisterEndpoint("/movies","D:/movies");
resourcing.RegisterEndpoint("/mqtt","C:/Users/User/Downloads/mqtt");

dbms.RegisterConnection("avonone","sqlserver","sqlserver://COLLECTIONS:COLLECTIONS@127.0.0.1");
dbms.RegisterConnection("avononeremote","remote","ws://127.0.0.1:1038/dbms-avonone/.json");
resourcing.RegisterEndpoint("/controls","D:/projects/system/controls");
console.Log("Start Service");
channel.Listener().Listen(":1044");
console.Log("Started Service");
/*dbms.RegisterConnection("avonone","sqlserver","sqlserver://COLLECTIONS:COLLECTIONS@127.0.0.1");
resourcing.RegisterEndpoint("/","D:/projects/inovo/clients/Avon/avonone");
resourcing.RegisterEndpoint("/inovoone","D:/projects/inovo/inovoone/one");
resourcing.RegisterEndpoint("/jquery","https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0");
resourcing.RegisterEndpoint("/webactions","C:/GitHub/kwe/webactions");
channel.Listener().Listen(":1038");*/
/*for (var i = 0; i < 1; i++) {
	var mqttclid = "mqtt"+(i+1);
	mqtting.RegisterConnection(mqttclid, {"broker":"ws://skullquake.dedicated.co.za", "port": 8080, "user": "emqx", "password":"public"})
	try {
	mqtting.Connect(mqttclid);
	 }catch(e){
		println("handling error;");
		println(e.toString()+";");
	}
}*/
/*console.Log("Setting up schedule...");
var schid="mySchedule";
sch0=channel
.Schedules()
.RegisterSchedule(
    schid,
    {
	"Seconds":1
    },
    request
);
sch0.AddAction(
	(function(args){
		console.Log("Executing "+args.schid);
	}).toString(),
	[{schid:schid}]
);    
sch0.Start();*/

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
