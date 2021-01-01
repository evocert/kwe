<@
import "bla";
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.MapEndpointResource("/","test-this.html","<\u0040println(\"test\");\u0040>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths())
println(resourcing.EndpointViaPath("/").Dirs())
println(resourcing.EndpointViaPath("/").Files())
println(resourcing.RegisteredRootPaths())
parseEval("println(\"from parseval\");")
channel.Listener().Listen(":1040");
 @> 