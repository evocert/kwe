<@ 
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.RegisterEndpoint("/jquery","https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0");
resourcing.RegisterEndpoint("/webactions","C:/GitHub/kwe/webactions");
resourcing.RegisterEndpoint("/etl","https://raw.githubusercontent.com/evocert/kwetl/main/src");
resourcing.RegisterEndpoint("/etl/wspace","D:/projects/system/etl/workspace/");
request.AddPath("/etl/init.js")
resourcing.RegisterEndpoint("/movies","D:/movies");
resourcing.RegisterEndpoint("/testthis","D:/projects/testthis");
channel.Listener().Listen(":1030");
resourcing.FS().MKDIR("/mem");
resourcing.FS().SET("mem/test.html",`<div>mem <@ print("test"); @></div>`);
resourcing.FS().SET("mem/index.html",`<html><body><#test/></body></html>`);
resourcing.FS().SET("mem/index.ht",`<html><body></body></html>`);
 @>