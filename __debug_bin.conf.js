<@
import "bla";
async function f1() {
    var x = await resolveAfter2Seconds(10);
    console.log(x); // 10
  }
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
println(resourcing.RegisteredPaths())
println(resourcing.EndpointViaPath("/").Dirs())
println(resourcing.EndpointViaPath("/").Files())
println(resourcing.RegisteredRootPaths())
parseEval("println(\"from parseval\");")
channel.Listener().Listen(":1040");

 @> 