[@include test@]
ssn.fs().mkdir("/gendocs","./gendocs");
try {
	eval(ssn.send("/gendocs/gendocs.js").readAll());
} catch(e){
	console.log(e.message);
}

var testxml=`<comments><data><CUSTOMERID>1107782</CUSTOMERID><COMMENT>Koos Comment test 2</COMMENT><COMMENTERID>0</COMMENTERID><COMMENTERNAME></COMMENTERNAME><CONTACTID>30428541</CONTACTID><CLASSSECTIONID>133912</CLASSSECTIONID><CLASSSECTIONNAME>Contact Centre</CLASSSECTIONNAME><CLASSCODE>133912</CLASSCODE><CLASSCODENAME>Balance&#x0D;
</CLASSCODENAME><DateTime>2022-04-27T17:07:36.860</DateTime></data><data><CUSTOMERID>1107782</CUSTOMERID><COMMENT>Koos Test  Avon Oneview comment testing at multiple account level</COMMENT><COMMENTERID>0</COMMENTERID><COMMENTERNAME></COMMENTERNAME><CONTACTID>30428541</CONTACTID><CLASSSECTIONID>133912</CLASSSECTIONID><CLASSSECTIONNAME>Contact Centre</CLASSSECTIONNAME><CLASSCODE>133912</CLASSCODE><CLASSCODENAME>Balance&#x0D;
</CLASSCODENAME><DateTime>2022-04-27T17:06:50.103</DateTime></data><data><CUSTOMERID>1107782</CUSTOMERID><COMMENT>Koos  Br 01 Test</COMMENT><COMMENTERID>0</COMMENTERID><COMMENTERNAME></COMMENTERNAME><CONTACTID>30428541</CONTACTID><CLASSSECTIONID>133912</CLASSSECTIONID><CLASSSECTIONNAME>Contact Centre</CLASSSECTIONNAME><CLASSCODE>133912</CLASSCODE><CLASSCODENAME>Balance&#x0D;
</CLASSCODENAME><DateTime>2022-04-27T16:58:42.383</DateTime></data></comments>`;

var cmntrec=ssn.dbms().query({"xml":{"data":testxml}});
if (cmntrec!=undefined && cmntrec!=null){
	console.log(cmntrec.json());
	/*print(`[`);
	var cmntcnt=0;
	while(cmntrec.next()){
		print(`{"data":`);
		print(JSON.stringify(cmntrec.dataMap()));
		print("}");
		if (cmntrec.isMore()){
			print(",");
		}
	}
	print("]");*/
}

ssn.dbms().register("test","sqlserver","server=LAPTOP-LPIKRBBA; database=ONER; user id=ONER; password=ONER;");
var exec=ssn.dbms().execute({"alias":"test","query":`EXECUTE  [ONER].[spLOGCOMMENT] 
@OWNERREFID = @OWNERREFID@
,@OWNERREFKEY =''
,@COMMENT=''
,@COMMENTERID = 0
,@COMMENTERNAME = ''
,@REFID = 0
,@REFKEY = ''
,@CLASSSECTIONID = 0
,@CLASSCODE = @CLASSCODE@`},{"OWNERREFID":0,"CLASSCODE":0});

if (exec!==undefined && exec!=null){
	console.log(exec.stmnt);
}

ssn.dbms().register("avon","oracle","oracle://SYSTEM:60N61ng0@localhost:1521/XE?MIN POOL SIZE=0&DECR POOL SIZE=1");
var reccomments=ssn.dbms().query({"alias":"avon","query":"select sysdate from dual"});
if(reccomments!==undefined && reccomments!==null){
	var cols=reccomments.columns();
	
	console.log("[")
	while(reccomments.next()){
		console.log("{")
		var data=reccomments.data();
		console.log(data);
		cols.forEach((col,coln)=>{
			console.log(`"${col}":${JSON.stringify(data[coln])}`);
			if(coln<cols.length-1){
				console.log(",");
			}
		});                                
		console.log("}");
		if(reccomments.isMore()){
			console.log(",");
		}
	}
	
	console.log("]");
} else {
	print("[]");
}
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
/*var jsn=jsonsax({"startobj":(jsx,key)=>{

},"startarr":(jsx,key)=>{

},"endarr":(jsx)=>{

},"startobj":(jsx,key)=>{

},"endobj":(jsx)=>{

},"setkeyval":(jsx,key,val,vtype)=>{

},"appendarr":(jsx,val,vtype)=>{

}},
`{"obj1":null}`);

try {
	while(jsn.next());
} catch (error) {
	console.log(error);
}


ssn.dbms().register("test","sqlserver","server=LAPTOP-LPIKRBBA; database=ONER; user id=ONER; password=ONER;");

rectest=ssn.dbms().query({"alias":"test","query":"select * from ONER.TESTXML"});

if (rectest!==undefined&&rectest!==null){
	while(rectest.next()){
		var recxml=ssn.dbms().query({"xml":{"data":rectest.data()[1]}});
		var recjson=ssn.dbms().query({"json":{"data":recxml.json()}});//JSON.stringify([{"ID":1,"CONTENT":"TEST"},{"ID":2,"CONTENT":"WEST"},{"ID":3,"CONTENT":"GEST"}])}});
		console.log(recjson.json());
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
*/
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
ssn.fs().mkdir("/","C:/projects/inovo/avon/avonone");


try {
	//eval(ssn.Send("/movies/schedule.js").ReadAll());
} catch(e){
	console.log(e.message);
}
ssn.dbms().register("avonone","sqlserver","server=localhost; database=ONER; user id=ONER; password=ONER");
ssn.dbms().register("oner-api","sqlserver","server=localhost; database=ONER; user id=ONER; password=ONER");
ssn.fs().mkdir("/oner","C:/projects/mystuff/oner");
ssn.listen("tcp",":1030");
ssn.listen("tcp",":1032");
ssn.cas().register(20,{"orginization":"bla"});
ssn.cas().ca(20).register(30);
ssn.fs().mkdir("/tvseries","C:/tvseries");
ssn.fs().mkdir("/movies","C:/movies");
ssn.fs().mkdir("/music","C:/music");
ssn.fs().mkdir("/oner","C:/projects/oner");

ssn.certifyAddr(20,30,":1032");
ssn.listen("tcp",":3336");
ssn.certifyAddr(`-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIUHZM13kT8YYe3WBwkcaf6LCGcjPswDQYJKoZIhvcNAQEL
BQAwYjELMAkGA1UEBhMCWkExEDAOBgNVBAgMB0dBVVRFTkcxETAPBgNVBAcMCFBy
ZXRvcmlhMQ4wDAYDVQQKDAVJbm92bzEOMAwGA1UECwwFSW5vdm8xDjAMBgNVBAMM
BUlub3ZvMB4XDTIyMDUxNjEzMjcwM1oXDTIzMDUxNjEzMjcwM1owYjELMAkGA1UE
BhMCWkExEDAOBgNVBAgMB0dBVVRFTkcxETAPBgNVBAcMCFByZXRvcmlhMQ4wDAYD
VQQKDAVJbm92bzEOMAwGA1UECwwFSW5vdm8xDjAMBgNVBAMMBUlub3ZvMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnNSbzrsDZXh4gET28qqxAEd9YAA5
vx+Cr/BTmz+cok6oJo7UZ9OVyXHwSVopmTuQZaiMPUBlT0LX51NLArsI7vWW3x70
u9N/L7K/nMZUKtEpvwi4R5gVmnVlfHH/XdkZN1VIehp1jDe8Bj9eoWRzKrjj0eC6
UT4Wi3Oa2MKj21ETp0/aDO1E52UKD+FNcXLH0ZpdlkPCBvCSsYd7J/6f43W3Mg+q
7IGNftgG0eGRFKKb/7RyWkUlwg8Vi/+xzLh1AZjCztHiP6Jep1SaVj8p9+jDuNsL
TEfJhT/jrWHkp+YBnu6xZuYlK3Z8/urkTDgY029qJnsqlWfMGsHQtL75/wIDAQAB
o1MwUTAdBgNVHQ4EFgQU8L4rbj1Bhn4QRqps1krM4wmKpIIwHwYDVR0jBBgwFoAU
8L4rbj1Bhn4QRqps1krM4wmKpIIwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEAkxzHuvRrgqtUPf9fEipD2EukoUyEyQG4G3uTp5wB2xm+YrvC7N19
a7RdXHS+3DmXQBLkO+P9mg0fUTOh+4LhUPUjOgAfDQ9dAyTH/3qVnPhScKSr6rbR
XbmpPDWmC1JNRARhH/IcY3kY3IdRdkZWnn2kJbz+E7CQkXEY+W0LdfZQYDoXfJaO
/85NN8ryf1zKsIxJWov3z3//Zg/kzEe8iRTbngkywVdCw00fwErBRcM5K5AYZlkW
flE3F7P+o5jtrgPo9m4w3E1H7Nwd3WbjnQmhE8HEqGjIJDmwhf8uNvo+o7PbZ01E
Aqogutte2JNnmVfaAUxk6NkEw0xWndgchg==
-----END CERTIFICATE-----
`,`-----BEGIN PRIVATE KEY-----
MIIEuwIBADANBgkqhkiG9w0BAQEFAASCBKUwggShAgEAAoIBAQCc1JvOuwNleHiA
RPbyqrEAR31gADm/H4Kv8FObP5yiTqgmjtRn05XJcfBJWimZO5BlqIw9QGVPQtfn
U0sCuwju9ZbfHvS7038vsr+cxlQq0Sm/CLhHmBWadWV8cf9d2Rk3VUh6GnWMN7wG
P16hZHMquOPR4LpRPhaLc5rYwqPbUROnT9oM7UTnZQoP4U1xcsfRml2WQ8IG8JKx
h3sn/p/jdbcyD6rsgY1+2AbR4ZEUopv/tHJaRSXCDxWL/7HMuHUBmMLO0eI/ol6n
VJpWPyn36MO42wtMR8mFP+OtYeSn5gGe7rFm5iUrdnz+6uRMOBjTb2omeyqVZ8wa
wdC0vvn/AgMBAAECgf86z3+eBoM/ie2mLDZuyZOWhzh8x5jgEvDvCTBRB4m1U8m6
q9T7Gl0RLajt2OHAlJWRiaMNVRiooGhWVuXKIFk5Qt9QzEr6JFWNjXpNUBI//C+r
c5mnP2DaiyuDfzxD9SV/mnuTTljGPCBGN31FCGYnny3PhnZYAPBzWua2YkcP9sCI
nYPiWc0a61JjLEsUiaaSItfiqxMPKNn9g7QNZ7Ymtk0wq3N0tTdUaHpblHXPE4IB
ugQgBw4vamPTbPP2Po0JRbqj291JvWgoa978/VbxQVsV8fvBiIi4M4b3tHiOd9X2
eUFpi/mHshQUbreWthiSCrjOSN3IOVcLjlvuxQECgYEAt1PrTfqiBFGUp1q3gn41
tcfxDkCsjng13RkfwRcDrugOVSvHxOD5dBsTd/Ez+UV+0rLXR9GFzk6igKpD81Qb
EBEGjWLQqWIF0+0M/2RJpNMh7AFii3m8eeyF1ewBIYeDrn82kOsYO5PYSYyTWnWJ
UVU1iMOrm+rSu9MI+oq3LQ0CgYEA2v+/FWf8F16PPzktT5uwdGo630bkpEjFnTIY
iziw53+bgkkGMGBw+UoaZ4CGFBsK8RQ9Gu323DeAYZI5qsVsBtu1gFJv4HlVWwAS
MdpUpVMwts6q+0kqiZbED/jCQm3hR07k4MmqNrh/PrZQFHLfnr3bxHV1w9JXZ4//
wjGv+DsCgYAikyklo8c6mUg359wEOFlY10SXM4tXs0Q1Vq+ucvC24/0QAxnB/8wM
Ia8iR9NNh9XLVv9TBCkAJ8RuD66RDaOs/AkIUUKZL2t59JMm93sMIuWa5Qf41hS9
yeXT2pa8BBrJpiRcYHpJgjCgbmq7/L7RIAjgqkaLOVZVJg/jcJXrYQKBgQC51IgV
pv3/6opNELyFL8xEjJvWOLbtPJ8LK6YuBPYACoUvwb1RsZRLxPiw2Rts1iGrvgu3
3TM7XJFAui7a3Nk03Jyf/dPXO74VEPNfgC+Rdg0BIk9uGYDR7bADCYf1jH2735NR
t24LLvUyhste/rcIYXypsS4z8zmdtFBHPZhfXwKBgF/41kLBOzCTA40N3JSjmGAe
6WFpWs+PUkcThL87KW3y4aMV0b3htOJgii6lNMaU7aG4obZ9nEzlQDERYzTbuQ5W
UxAnla0cPaZWvzuYXqphs1jjRLTB3lSLPr/rup40BWbW0hzG13PFZoTNyLF+pALc
25kPoEQPh7whHhhjEfg+
-----END PRIVATE KEY-----
`,":3336");

//ssn.certifyAddr("20","30",":1030")
