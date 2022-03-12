/*ssn.caching().put("a",{"b":[1,2,{"f":[7,8,9,18]},3,{"g":[17,18,19,118]},4]});
console.log(ssn.caching().string());
if (ssn.caching().existsAt("a","b",2)) {
	console.log(ssn.caching().at("a","b",2).string());
	ssn.caching().clearAt("a","b",2)
	console.log(ssn.caching().at("a","b",2).string());
}
console.log(ssn.caching().string());
if (ssn.caching().existsAt("a","b",2)) {
	ssn.caching().closeAt("a","b",2)
	console.log(ssn.caching().at("a","b",2));
	console.log(ssn.caching().string());
}

	ssn.fs().mkdir("/kweauth","C:/github/kweauth");
*/
ssn.dbms().register("test","sqlserver","server=LAPTOP-LPIKRBBA; database=ONER; user id=ONER; password=ONER;");

rectest=ssn.dbms().query({"alias":"test","query":"select * from ONER.TESTXML"});

if (rectest!==undefined&&rectest!==null){
	while(rectest.next()){
		var recxml=ssn.dbms().query({"xml":{"data":rectest.data()[1]}});
		console.log(recxml.json());
	}
}

ssn.dbms().register("avon","oracle","oracle://SYSTEM:60N61ng0@localhost:1521/XE?MIN POOL SIZE=0&DECR POOL SIZE=1");
var rec=ssn.dbms().query({alias:"avon",query:"select sysdate as d from dual"});
rec.next();
ssn.fs().mkdir("/playground","C:/projects/playground");
ssn.fs().mkdir("/crm","https://demo.1crmcloud.com/");

ssn.fs().mkdir("/kweauth","C:/github/kweauth");
ssn.fs().mkdir("kwetl","C:/GitHub/kwetl");
ssn.env().set("kwetl-path","/kwetl");

ssn.fs().mkdir("/gendocs","./gendocs");
try {
	eval(ssn.send("/gendocs/gendocs.js").readAll());
} catch(e){
	console.log(e.message);
}

ssn.dbms().register("b1","kwesqlite",":memory:");

ssn.dbms().execute({"alias":"b1","query":`CREATE TABLE t1 (
	v1 INT,
    v2 VARCHAR(1000)
);`});

ssn.dbms().execute({"alias":"b1","query":`CREATE TABLE t2 (
	v1 INT,
    v2 VARCHAR(1000)
);`});

for (let index = 0; index < 1; index++) {
	var vt="b"+index;
	ssn.dbms().execute({"alias":"b1","query":`insert into t1 (v1,v2) values(@p1@,@p2@)`},{"p1":index+1,"p2":vt});
}
console.log(JSON.stringify(ssn.dbms().info("b2")));

var a=ssn.dbms().query({"alias":"b1","query":"select v1 as p1,v2 as p2 from t1"},{"prm1":42});

//console.log(JSON.stringify(kwemethods(a)));
console.log(a.json());
//ssn.dbms().unregister("b1");

ssn.dbms().register("postgres_","postgres","user=postgres password=n@n61ng@ host=localhost port=5433 dbname=postgres sslmode=disable")

ssn.fs().mkdir("/movies","D:/movies");
/*
ssn.fs().mkdir("/kwehyg","http://skullquake.dedicated.co.za:3001/ockert/kwehyg/raw/master/");
ssn.fs().mkdir("/kwetest","http://skullquake.dedicated.co.za:3001/ockert/kwedt/raw/master/");
ssn.fs().mkdir("/materialfonts","https://cdn.jsdelivr.net/npm/@mdi/font@6.5.95");
*/

var goosgoarchsarr=`aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/amd64
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
illumos/amd64
ios/amd64
ios/arm64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/mips64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
windows/386
windows/amd64
windows/arm
windows/arm64`.split("\n");

//eval(ssn.fsutils().cat("./build/build.js").readAll());
//buildgo(goosgoarchsarr,`C:/GitHub/kwe/kwe.go`,`C:/GitHub/kwe/build/kwe`,`C:/GitHub/kwe/build/upx`);

/*var cmd=ssn.command("cmd");
goosgoarchsarr.forEach((goosarch)=>{
	var keyvalue=goosarch.trim().split("/");
	var key=keyvalue[0].trim();
	var value=keyvalue[1].trim();
	if (value==="" || key==="") return;
	console.log(`${key}: ${value}`);
	
	try {
		cmd.setReadTimeout(10000,100);
		cmd.readAll();
		cmd.println("SET GOOS="+key);
		cmd.println("SET GOARCH="+value);
		if (key==="ios") {
			cmd.println("SET CGO_ENABLED=1");
		}
		cmd.println(`go build -v -ldflags "-w -s" -o C:/GitHub/kwe/build/kwe_`+key+`_`+value+(key==="windows"?".exe":(key==="js"?".wasm":(key==="darwin"?".dmg":"")))+` C:/GitHub/kwe/ssn.go`);
		cmd.println("echo finit");
		for(var ln = cmd.readln();!ln.endsWith("finit");ln= cmd.readln()){
			if (ln!=="") {
				console.log(ln);
			}
		}
	} catch (error) {
		console.log("error[",key,":",value,"]:",error.toString());
	}
	
	console.log("done -",key,":",value);	
});
cmd.close();
console.log("done - build");

var cmd=ssn.command("cmd");
goosgoarchsarr.forEach((goosarch)=>{
	var keyvalue=goosarch.trim().split("/");
	var key=keyvalue[0].trim();
	var value=keyvalue[1].trim();
	if (value==="" || key==="") return;
	console.log(`${key}: ${value}`);
	
	try {
		cmd.setReadTimeout(10000,100);
		cmd.readAll();
		cmd.println(`C:/GitHub/kwe/build/upx C:/GitHub/kwe/build/kwe_`+key+`_`+value+(key==="windows"?".exe":(key==="js"?".wasm":(key==="darwin"?".dmg":""))))
		cmd.println("echo finit");
		for(var ln = cmd.readln();!ln.endsWith("finit");ln= cmd.readln()){
			if (ln!=="") {
				console.log(ln);
			}
		}
	} catch (error) {
		console.log("error[",key,":",value,"]:",error.toString());
	}
	
	console.log("done -",key,":",value);	
});
cmd.close();
console.log("done - all");	
*/
/*goosgoarch.array.forEach(element => {
	var goosv=element;
	var goarchv=goosgoarch[goosv];
	console.log(goosv,":",goarchv);
});*/

//cmd.setReadTimeout(1000*60*5);


try {
	//eval(ssn.Send("/movies/schedule.js").ReadAll());
} catch(e){
	console.log(e.message);
}
ssn.listen("tcp",":80");
