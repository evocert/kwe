<@ 
//dbms.RegisterConnection("avon","oracle","oracle://SYSTEM:N%40N61ng%40@localhost/XE");
//resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.RegisterEndpoint("/etl","https://raw.githubusercontent.com/evocert/kwetl/main/src");
resourcing.RegisterEndpoint("/etl/wspace","D:/projects/system/etl/workspace/");
request.AddPath("/etl/init.js")
channel.Listener().Listen(":1030");
 @>