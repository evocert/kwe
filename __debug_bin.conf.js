<@ 
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths())
println(resourcing.EndpointViaPath("/").Dirs())
println(resourcing.EndpointViaPath("/").Files())
println(resourcing.RegisteredRootPaths())
parseEval("println(\"from parseval\");")
channel.Listener().Listen(":1040");
require(["https://cdn.jsdelivr.net/npm/@babel/standalone@7.12.12/babel.min.js"]);
 @> 