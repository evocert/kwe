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
kwe.fs().mkdir("/gendocs","./gendocs");
try {
	eval(kwe.send("/gendocs/gendocs.js").readAll());
} catch(e){
	console.log(e.message);
}

kwe.dbms().registerConnection("b1","kwesqlite",":memory:",{"max-open-cons":1});
var a=kwe.dbms().query({"alias":"b1","query":"select @@prm1@@"},{"prm1":42});
console.log(JSON.stringify(kwemethods(a)));
console.log(a.json());
kwe.dbms().unregisterConnection("b1");

kwe.fs().mkdir("/movies","D:/movies");


var goosgoarch={"android":"arm",
"darwin":"386",
"darwin":"amd64",
"darwin":"arm",
"darwin":"arm64",
"dragonfly":"amd64",
"freebsd":"386",
"freebsd":"amd64",
"freebsd":"arm",
"linux":"386",
"linux":"amd64",
"linux":"arm",
"linux":"arm64",
"linux":"ppc64",
"linux":"ppc64le",
"linux":"mips",
"linux":"mipsle",
"linux":"mips64",
"linux":"mips64le",
"netbsd":"386",
"netbsd":"amd64",
"netbsd":"arm",
"openbsd":"386",
"openbsd":"amd64",
"openbsd":"arm",
"plan9":"386",
"plan9":"amd64",
"solaris":"amd64",
"windows":"386",
"windows":"amd64"};


for (const [key, value] of Object.entries(goosgoarch)) {
	sleep(10);
	var cmd=kwe.command("cmd");
	cmd.readAll();
	console.log(`${key}: ${value}`);
	var s="";
	s+=("SET GOOS="+key+"\n");
	s+=("SET GOARCH="+value+"\n");
	s+=(`go build -ldflags "-w -s" -o D:/movies/kwebuilds/scripts/buildbin/kwe_`+key+`_`+value+(key==="windows"?".exe":"")+` C:/GitHub/kwe/kwe.go`+"\n");
	cmd.print(s);
	cmd.println("echo finit");
	for(var ln = cmd.readln();!ln.endsWith("finit");ln= cmd.readln()){
		if (ln!=="") {
			console.log(ln);
		}
	}
	console.log("kla-",key,":",value);
	cmd.close();
}

console.log(s);

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
