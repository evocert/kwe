<@

var k="a3";
var arr1=[];
caching.Clear();
caching.Put(k,[arr1]);
arr1=caching.Find(k);
caching.Push(k,[[8]]);
caching.Push(k,[8]);
caching.Shift(k,18);
caching.Fprint(request);
println(caching.At(k,[[1,0]]));
//caching.Push(k,8)
caching.Put("a4","hjkhjhkhj");
caching.Put("a3",{"d1":89989,"d2":{"f1":6}});
var obj=(caching.ValueByIndex(0));
caching.Find("a3").Find("d2");
caching.Reset();
caching.Fprint(request);



dbms.RegisterConnection("avonone","sqlserver","sqlserver://COLLECTIONS:COLLECTIONS@127.0.0.1");
resourcing.RegisterEndpoint("/","D:/projects/inovo/clients/Avon/avonone");
resourcing.RegisterEndpoint("/inovoone","D:/projects/inovo/inovoone/one");
resourcing.RegisterEndpoint("/jquery","https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0");
resourcing.RegisterEndpoint("/webactions","C:/GitHub/kwe/webactions");
channel.Listener().Listen(":1038");

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