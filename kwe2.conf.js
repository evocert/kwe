<@
console.Log("START");
console.Log("Connecting Databases...");
dbms.RegisterConnection("pg","postgres","user=postgres password=1234!@#$qwerQWER dbname=postgres sslmode=disable");
dbms.RegisterConnection("mysql","mysql","mysql:1234!qwer!QWER@tcp(localhost)/test");
console.Log("Registering Endpoints...");
resourcing.RegisterEndpoint("/","./www");
resourcing.RegisterEndpoint("/master","https://raw.githubusercontent.com/evocert/kweutils/main/");
//resourcing.RegisterEndpoint("/","C:/tmp/www/kweutils");
//resourcing.RegisterEndpoint("/reqtstzip","C:/tmp/www/requirejsexcercise.zip");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
//println(resourcing.FindRSString("/test-this.html"))
console.Log(resourcing.RegisteredPaths())
console.Log(resourcing.RegisteredRootPaths())
console.Log("Starting listener...");
channel.Listener().Listen(":80");
console.Log("END");
@>
