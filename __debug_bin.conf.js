<@
dbms.RegisterConnection("avon","oracle","oracle://SYSTEM:N%40N61ng%40@localhost/XE");
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.RegisterEndpoint("/bs","https://cdn.jsdelivr.net/npm/bootstrap");
resourcing.RegisterEndpoint("/master/kweutils","https://raw.githubusercontent.com/evocert/kweutils/main/");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths().join("\r\n"));
println(resourcing.RegisteredRootPaths().join("\r\n"));
channel.Listener().Listen(":1040");
 @>