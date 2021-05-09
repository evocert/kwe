<@ 
dbms.RegisterConnection("avonone","sqlserver","sqlserver://COLLECTIONS:COLLECTIONS@127.0.0.1");
resourcing.RegisterEndpoint("/","D:/projects/inovo/clients/Avon/avonone");
resourcing.RegisterEndpoint("/inovoone","D:/projects/inovo/inovoone/one");
resourcing.RegisterEndpoint("/jquery","https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0");
resourcing.RegisterEndpoint("/webactions","C:/GitHub/kwe/webactions");
channel.Listener().Listen(":1037");
console.Log("Setting up schedule...");
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
sch0.Start();
@>