<@ 
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths())
println(resourcing.RegisteredRootPaths())
channel.Listener().Listen(":1030");
 @>