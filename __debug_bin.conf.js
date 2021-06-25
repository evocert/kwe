<@
resourcing.RegisterEndpoint("/gendocs","./gendocs");
//println(channel.Send("/active:gendocs/gendocs.js"));
resourcing.RegisterEndpoint("/","D:/projects/inovo/clients/Avon/one/system");
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
@>