<@
dbms.RegisterConnection("avon","oracle","oracle://SYSTEM:N%40N61ng%40@localhost/XE");
resourcing.RegisterEndpoint("/","./");
//resourcing.RegisterEndpoint("/cdnjs","https://cdnjs.cloudflare.com/ajax/libs/");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths().join("\r\n"));
println(resourcing.RegisteredRootPaths().join("\r\n"));
channel.Listener().Listen(":1040");
 @>