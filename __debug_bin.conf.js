<@
dbms.RegisterConnection("avon","oracle","oracle://SYSTEM:N%40N61ng%40@localhost/XE");
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.RegisterEndpoint("/bs","https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta2/dist");
resourcing.RegisterEndpoint("/master/kweutils","https://raw.githubusercontent.com/evocert/kweutils/main/");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
channel.Listener().Listen(":1040");
//println(resourcing.FindRSString("/test-this.html"))
/*println(resourcing.RegisteredPaths().join("\r\n"));
println(resourcing.RegisteredRootPaths().join("\r\n"));
channel.Listener().Listen(":1040");
var test1=channel.Schedules().RegisterSchedule("test1",{"Seconds":20},request);
for (var i=0;i<20;i++){
    test1.AddAction((i+""), function() {
        println("test action ",(i+""));
    });
}*/
 @>