/*kwe.caching().put("a",{"b":[1,2,{"f":[7,8,9,18]},3,{"g":[17,18,19,118]},4]});
console.log(kwe.caching().string());
if (kwe.caching().existsAt("a","b",2)) {
	console.log(kwe.caching().at("a","b",2).string());
	kwe.caching().clearAt("a","b",2)
	console.log(kwe.caching().at("a","b",2).string());
}
console.log(kwe.caching().string());
if (kwe.caching().existsAt("a","b",2)) {
	kwe.caching().closeAt("a","b",2)
	console.log(kwe.caching().at("a","b",2));
	console.log(kwe.caching().string());
}

kwe.fs().mkdir("/kweauth","C:/github/kweauth");
*/

kwe.fs().mkdir("/kweauth","C:/github/kweauth");
kwe.fs().mkdir("kwetl","C:/GitHub/kwetl");
kwe.env().set("kwetl-path","/kwetl");

kwe.fs().mkdir("/gendocs","./gendocs");
try {
	eval(kwe.send("/gendocs/gendocs.js").readAll());
} catch(e){
	console.log(e.message);
}

kwe.dbms().register("b1","kwesqlite",":memory:");

kwe.dbms().execute({"alias":"b1","query":`CREATE TABLE t1(
	v1 INT,
    v2 VARCHAR(1000)
);`});

kwe.dbms().execute({"alias":"b1","query":`CREATE TABLE t2(
	v1 INT,
    v2 VARCHAR(1000)
);`});

for (let index = 0; index < 100000; index++) {
	var vt="b"+index;
	kwe.dbms().execute({"alias":"b1","query":`insert into t1 (v1,v2) values(@@p1@@,@@p2@@)`},{"p1":index+1,"p2":vt});
}
console.log(JSON.stringify(kwe.dbms().info("b2")));

var a=kwe.dbms().query({"alias":"b1","query":"select v1 as p1,v2 as p2 from t1","exec":[
	{"alias":"b1","query":"insert into t2 (v1,v2) values(@@p1@@,@@p2@@);"}
]},{"prm1":42});

//console.log(JSON.stringify(kwemethods(a)));
console.log(a.json());
//kwe.dbms().unregister("b1");

kwe.dbms().register("postgres_","postgres","user=postgres password=n@n61ng@ host=localhost port=5433 dbname=postgres sslmode=disable")

kwe.fs().mkdir("/movies","D:/movies");

kwe.fs().mkdir("/kwehyg","http://skullquake.dedicated.co.za:3001/ockert/kwehyg/raw/master/");
kwe.fs().mkdir("/kwetest","http://skullquake.dedicated.co.za:3001/ockert/kwedt/raw/master/");
kwe.fs().mkdir("/materialfonts","https://cdn.jsdelivr.net/npm/@mdi/font@6.5.95");

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

//eval(kwe.fs().cat("/kwe/build/js/build.js").readAll());
//buildgo(goosgoarchsarr,`C:/GitHub/kwe/kwe.go`,`C:/GitHub/kwe/build/kwe`,`C:/GitHub/kwe/build/upx`);

/*var cmd=kwe.command("cmd");
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
		cmd.println(`go build -v -ldflags "-w -s" -o C:/GitHub/kwe/build/kwe_`+key+`_`+value+(key==="windows"?".exe":(key==="js"?".wasm":(key==="darwin"?".dmg":"")))+` C:/GitHub/kwe/kwe.go`);
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

var cmd=kwe.command("cmd");
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
	//eval(kwe.Send("/movies/schedule.js").ReadAll());
} catch(e){
	console.log(e.message);
}
