<@
dbms.RegisterConnection("avon","oracle","oracle://SYSTEM:N%40N61ng%40@localhost/XE");
resourcing.RegisterEndpoint("/","D:/projects/system");
resourcing.RegisterEndpoint("/bs","https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta2/dist");
resourcing.RegisterEndpoint("/master/kweutils","https://raw.githubusercontent.com/evocert/kweutils/main/");
resourcing.MapEndpointResource("/","test-this.html","<h1>test this</h1>");
channel.Listener().Listen(":1040");
//println(resourcing.FindRSString("/test-this.html"))
//println(resourcing.RegisteredPaths().join("\r\n"));
//println(resourcing.RegisteredRootPaths().join("\r\n"));
var fis = _fsutils.LS("D:/projects/system/bootstrap/css","bla");
    if (fis!==undefined) {
        fis.forEach(function(fi){
            //println(fi.JSON())
        });
        //println(_fsutils.FINFOPATHSJSON(fis));
    }
var cntdone=1;
for (var j=0;j<cntdone/1;j++){
    var test1=channel.Schedules().RegisterSchedule("test"+(j+1),{"Seconds":2},request);
    test1.AddInitAction(`function() {
        console.Log("test init");
    }`);
    for (var i=0;i<(cntdone*2);i++){
        test1.AddAction(`function() {
            //console.Log("test this"+_fsutils.FINFOPATHSJSON(fis));
            //cntdone--;
            return true;
        }`);
    }
    //test1.AddAction({"request":{"path":"/test/this.js"}});
    test1.AddWrapupAction(`function() {
        console.Log("test wrapup");
    }`);
    test1.Start();
}
 @>