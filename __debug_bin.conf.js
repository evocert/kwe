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

var email=ssn.readMail(`Received: from HMC-DAG02SVR04.he.businessgateway.co.za (192.168.200.113) by
HMC-DAG02SVR02.he.businessgateway.co.za (192.168.200.111) with Microsoft SMTP
Server (TLS) id 15.0.1497.32 via Mailbox Transport; Sun, 22 May 2022 07:16:00
+0200
Received: from HMC-DAG05SVR01.he.businessgateway.co.za (192.168.200.142) by
HMC-DAG02SVR04.he.businessgateway.co.za (192.168.200.128) with Microsoft SMTP
Server (TLS) id 15.0.1497.32; Sun, 22 May 2022 07:15:59 +0200
Received: from securemail-pl-mx25.synaq.com (196.35.198.137) by
HMC-DAG05SVR01.he.businessgateway.co.za (192.168.200.142) with Microsoft SMTP
Server (TLS) id 15.0.1497.32 via Frontend Transport; Sun, 22 May 2022
07:15:59 +0200
ARC-Seal: i=2; cv=pass; a=rsa-sha256; d=synaq.com; s=securemail; t=1653196559;
	b=BMuJAGFmBmRsihdx7VEG1j5hVzmU+PyMfPwgTzhJ6Ah3WCVf3MYU5B+Jb28yCBKl7KgJDFeIFC
	 QRqAE/Evm83Y74KXZ5ou0U2jjOBhXl7eqgMx7HrutrDb1g6i646F5mBAp0UAolBXpE4S6ondZY
	 hgF9m5MQJSYZk7W3ikC5MQoQtiJdMQZyL967AyekG/rF8+GtaY1So65MZyQ8VSgXvE/gO4PgMW
	 U5Uci3jT0L/XzQLSTvvIHqGK396iqC+frsnSSjIGogrb3joOxIV6TOjlfwodaGENdLLJzCsgS6
	 tJniHSttr1kL3JaNf7q0lSVBHPZmX9hp36EQMMLkULKmbg==;
ARC-Authentication-Results: i=2; synaq.com;
   iprev=pass (mail-mw2nam12on2069.outbound.protection.outlook.com) smtp.remote-ip=40.107.244.69;
   spf=pass smtp.mailfrom=avon.com;
   dkim=pass header.d=avonglobal.onmicrosoft.com header.s=selector2-avonglobal-onmicrosoft-com header.a=rsa-sha256;
   dkim=pass header.d=gmail.com header.s=20210112 header.a=rsa-sha256;
   dmarc=pass header.from=avon.com;
   arc=pass (i=1) header.s=arcselector9901 arc.oldest-pass=1 smtp.remote-ip=40.107.244.69
ARC-Message-Signature: i=2; a=rsa-sha256; c=relaxed; d=synaq.com; s=securemail; t=1653196559;
   bh=2wRFM4ujm2tD0546GDKJDUL0JOTX1jgzbYVETNiRad4=;
   h=Content-Type:To:Subject:Message-ID:Date:From:MIME-Version:DKIM-Signature:
	 Resent-From:DKIM-Signature;
   b=CxGFZjQbNQ5rp6yNzpqgS3d7qtF943nGbQ3slcgnjqZMwFMiSx7hJTVsMWEHyqabTa7JlLZLyh
	 40bk+DTuFNLqQKVLS40+PpoWIzuwmsG+0AYSbRt95POPt+rUJVBpBhvILTpJkmvJ80nY+53PR2
	 Fwy3pBdFx22FEmu3qHPgz79nyYxQwHMj9gi1t/DNpFQ2kPBm/puDWplbglpo18CAMeGw47fieM
	 TdKg16PJ1Cwe2sIUOIXjORpttu/zlVhhlIHtLHdRH+S9b734LTv6mT+x8qbULCu0csgHBJk+2d
	 9X5l6v1BjKpJ1s6Fa98HbyCgVk6Va/rMqNrnYp6MJDKRyA==;
Authentication-Results: synaq.com;
   iprev=pass (mail-mw2nam12on2069.outbound.protection.outlook.com) smtp.remote-ip=40.107.244.69;
   spf=pass smtp.mailfrom=avon.com;
   dkim=pass header.d=avonglobal.onmicrosoft.com header.s=selector2-avonglobal-onmicrosoft-com header.a=rsa-sha256;
   dkim=pass header.d=gmail.com header.s=20210112 header.a=rsa-sha256;
   dmarc=pass header.from=avon.com;
   arc=pass (i=1) header.s=arcselector9901 arc.oldest-pass=1 smtp.remote-ip=40.107.244.69
Received: from mail-mw2nam12on2069.outbound.protection.outlook.com ([40.107.244.69] helo=NAM12-MW2-obe.outbound.protection.outlook.com)
   by securemail-pl-mx25.synaq.com with esmtps  (TLS1.2) tls TLS_RSA_WITH_AES_256_GCM_SHA384
   (Exim 4.95)
   (envelope-from <queries@avon.com>)
   id 1nsdwE-000b99-BC
   for queries@onlyfromavon.co.za;
   Sun, 22 May 2022 07:15:50 +0200
ARC-Seal: i=1; a=rsa-sha256; s=arcselector9901; d=microsoft.com; cv=none;
b=O7aCz0v1hMmr6GXuGuZRr93gFBjaFAWR6hd5HQxi8wLEGXZI0g4OuWTZRN6VsyGpop+k1jay2CLyHhZFlANF2YNKlPaYxiPgLWfMiQiRzKm3nBYBZW3adaFYWVId09tovEzbrBz+93EERDbncZvoQH2VVw9Y2tTJWzVHs8nOxDfoV8mGiiNmJejP12ZwkOEOwznenVAlH0awILdj/ojHZPzHaAZykP85WB8Jd0mekUdTW6XZMjKN0xIQa2kdfyHAAkenFX3hYFqwFfDnQb+cxAnCfh+uukImW6F5BMQRfVSgix5Z7lmZV0W6Cc6JBahvQKq3D+Pj67JlHqc/2IE9vw==
ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=microsoft.com;
s=arcselector9901;
h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-AntiSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Exchange-AntiSpam-MessageData-1;
bh=2wRFM4ujm2tD0546GDKJDUL0JOTX1jgzbYVETNiRad4=;
b=kGb/rOCW9Hm2q3tUpaYZAbBQyO/+5JvkhsYsI1jSS4KcO7aAk5R/Q7hA2eRW0INCqqCQ9vF5bEp1Xftg+O8H0XxKXgP3Br5ExE4vNyMdn8XAF+76ftLQsTNSmMssSzmrkyozhmx7GBa578DQDq2zDWKjdIZ5fIcIsM2JJ4J5IGjKoX1VOwhbyfKYEk3/Lqim2N7M537AFQwyXqDtt4rqgK+fp5EWbXxQoPaIdLQXGFYTk3u/gO78dZNicyYSMZdA0VzA3Q4ho54QIZykc3MCfhCB0p+F1CPmW/fyGYNKoapYLTJtPQnnMXCDs/yEPucRTK8xPg0hxShCA22jwKCsCQ==
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass (sender ip is
209.85.160.43) smtp.rcpttodomain=avon.com smtp.mailfrom=gmail.com; dmarc=pass
(p=none sp=quarantine pct=100) action=none header.from=gmail.com; dkim=pass
(signature was verified) header.d=gmail.com; arc=none (0)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
d=avonglobal.onmicrosoft.com; s=selector2-avonglobal-onmicrosoft-com;
h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck;
bh=2wRFM4ujm2tD0546GDKJDUL0JOTX1jgzbYVETNiRad4=;
b=SUhT7x738HbP+ufg2vX5kn+sEHjefRy4PfMBDtqek958R4AbTk+QuDnJDQlgOG/vGYkACB3cjXf5m0zKrxuES/BiYy0kcKQieuupT1gtV9vo1jcYgKYhshv5ObqsoqmKaltR5Xi5I3Vlcu8FJ3hwp8Ao6C3Rb3Q9HrkGYSDq0HA=
Resent-From: <queries@avon.com>
Received: from BN0PR08CA0012.namprd08.prod.outlook.com (2603:10b6:408:142::34)
by BN3PR07MB2691.namprd07.prod.outlook.com (2a01:111:e400:7bb9::12) with
Microsoft SMTP Server (version=TLS1_2,
cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.5273.13; Sun, 22 May
2022 05:15:41 +0000
Received: from BN8NAM11FT045.eop-nam11.prod.protection.outlook.com
(2603:10b6:408:142:cafe::41) by BN0PR08CA0012.outlook.office365.com
(2603:10b6:408:142::34) with Microsoft SMTP Server (version=TLS1_2,
cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.5273.15 via Frontend
Transport; Sun, 22 May 2022 05:15:41 +0000
Authentication-Results: spf=pass (sender IP is 209.85.160.43)
smtp.mailfrom=gmail.com; dkim=pass (signature was verified)
header.d=gmail.com;dmarc=pass action=none header.from=gmail.com;
Received-SPF: Pass (protection.outlook.com: domain of gmail.com designates
209.85.160.43 as permitted sender) receiver=protection.outlook.com;
client-ip=209.85.160.43; helo=mail-oa1-f43.google.com; pr=C
Received: from mail-oa1-f43.google.com (209.85.160.43) by
BN8NAM11FT045.mail.protection.outlook.com (10.13.177.47) with Microsoft SMTP
Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id
15.20.5273.14 via Frontend Transport; Sun, 22 May 2022 05:15:41 +0000
Received: by mail-oa1-f43.google.com with SMTP id 586e51a60fabf-edf3b6b0f2so14605652fac.9
	   for <queries@avon.com>; Sat, 21 May 2022 22:15:41 -0700 (PDT)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
	   d=gmail.com; s=20210112;
	   h=mime-version:from:date:message-id:subject:to;
	   bh=2wRFM4ujm2tD0546GDKJDUL0JOTX1jgzbYVETNiRad4=;
	   b=HH0sKjLnTPd2sWDhr2cjyhgnCT+1ACrefr6s5/2XuOAsv9SoPovKYP3KkIrTPt5pqM
		oIMJvl6xtfWPSzmC9M5JsEW1EtmzTUXDj34+ieBn1kZaocDRijoi2sJrOxYn/PMsBdG/
		GmnkjoNQvOPMs37J3Vi2LjvSnxUi7QZEj1SsN/i3Ji4CYcLLnD685+3PxK1XkLpll5ob
		8QGWQm2X53frAEIOwMOXsOn1sAbwQN7fw3h3aXAS8g/JCj1W3V/uifaqHOpcbVFaRSrS
		VPuXGoRYk2US0LIbWQay+1TkN+H9Z5ahwmOnah2BaV4gf8+aLAIPxthSFT0UPHZKn7pK
		CSsg==
X-Google-DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
	   d=1e100.net; s=20210112;
	   h=x-gm-message-state:mime-version:from:date:message-id:subject:to;
	   bh=2wRFM4ujm2tD0546GDKJDUL0JOTX1jgzbYVETNiRad4=;
	   b=GtJWYjG8tg/UXbA4ooFdn/0+80Em+aOeGAAFIsa2CQQvaCMSRKeQG/93Ofbt4bWnUA
		hIARsUFEUp1IezSl0NthLuAQ05JwIlBhdbDsUy8kLmvaahTJ8bQ71vXHURnLUWZ2NQum
		x3YSl9O8fjZqVS5LqXO0CSufsO/jx/9YfnjrxUfA6l+YhTperh77rzuUcYbW5wnG+bct
		3b3J6CGoMn8IeazNl4w3tDG20YeXX5OoEy9cUUe2rT2iynUCmSInEr7QbOnlmOj1VGT/
		1aDOOa0YvqnsLamnWaIUl5JKzmfPXHsPkpnoTWlgeb4MDv/LmVnroD+2//0OaLHoJNPB
		tXCw==
X-Gm-Message-State: AOAM53000Xp6L86Eso4N6O3y/knPm4PpHiQwdNHiw6YTX4xyt2hp/Hbf
   YYdXF7sUmiR6/0p2sSyiYD6XOqWvKVIzT1pKEIaG8aX4
X-Google-Smtp-Source: ABdhPJywnooPzS/y+4sSJ/ZblHMDwoZuIzKpG6lPAOIzczKv6MX9jTPxW42DnodMLdZUa3SoJnZIt/eAeV7jDQX5RR4=
X-Received: by 2002:a05:6870:15c3:b0:ed:9980:db99 with SMTP id
k3-20020a05687015c300b000ed9980db99mr9837939oad.154.1653196540687; Sat, 21
May 2022 22:15:40 -0700 (PDT)
From: Nonkie Mngomezulu <nonkiemngomezulu56@gmail.com>
Date: Sun, 22 May 2022 07:13:23 +0200
Message-ID: <CADr+mSCJdCUdcnPpMC_D1RDSc7Wrbk7x6aH3n_djufZuXtFy_w@mail.gmail.com>
Subject: Worried
To: queries@avon.com
Content-Type: multipart/alternative; boundary="0000000000008d3ea305df92d0d7"
X-EOPAttributedMessage: 0
X-EOPTenantAttributedMessage: e7f5a0d7-8564-45b1-bdf4-5a28b6136195:0
X-MS-PublicTrafficType: Email
X-MS-Office365-Filtering-Correlation-Id: c2d8449d-08a7-43f0-4988-08da3bb21daf
X-MS-TrafficTypeDiagnostic: BN3PR07MB2691:EE_
X-LD-Processed: e7f5a0d7-8564-45b1-bdf4-5a28b6136195,ExtFwd
X-Microsoft-Antispam-PRVS: <BN3PR07MB2691432E4B0CD5863C75F5FFC5D59@BN3PR07MB2691.namprd07.prod.outlook.com>
X-MS-Exchange-SenderADCheck: 0
X-MS-Exchange-AntiSpam-Relay: 0
X-Microsoft-Antispam: BCL:0;
X-Microsoft-Antispam-Message-Info: rPWErzU1n2M+9sgmkry6fsSCabe/YnRNDatPFf4Z/M+o7JLAY9dY7T1sfgEQYMKSomHw7W123g9+Gzv/G0nJbvJttKr7MtKll/AhpOkpSqYFXRv5lTaSKXukm5KBn8uP/PKtqyRwsYJKGCvakm8wgEuWbup5movacfUuBB3dTSsiEJPV+VJLUUhJiw7cG5ee8a0/TUaa4ZwbCpyAbMfN/WGowSO6DO8zEa8ANAIP512hYTZ8AMfNS3STm2/9a22n4bPqhwe53DmfGvjdobIDI/eXScOC9RttBOKbRzioDl0zRm1rEEXqVzJDgFL6MP+0ilfEi148TffZh76AoJKNJKMDuOTD/u5sEwuHXjwbfWungSoNuYq9REquT54RRcwir/fRs1GKAdE/QVJ6iGELbM84ZUgV9zahmrbkYfDNCMqy2ZnGkQHhfqNvtk48UQd7XT3vc/p/2sE8Hf+bD97M/lEKYaDoX8EVHwnDBr4GJXG/vEGEKiloimX+ZTp9Q84UTWrKIU7mfx9NJtQMENrxc2Nzf3EoWrRxyEGuZhDEJK97WyCE3MpK63zFq2km/R+JW6evpRUAQNrrLGG+3wuY1CY+N7GQJ3+Y6axbuLygZAJ4iYL94lzHsB4XFLmHuw4Z
X-Forefront-Antispam-Report: CIP:209.85.160.43;CTRY:US;LANG:en;SCL:1;SRV:;IPV:NLI;SFV:NSPM;H:mail-oa1-f43.google.com;PTR:mail-oa1-f43.google.com;CAT:NONE;SFS:(13230001)(84050400002)(8676002)(5660300002)(70586007)(316002)(42186006)(76482006)(55446002)(6666004)(7116003)(86362001)(68406010)(34206002)(508600001)(7636003)(33964004)(356005)(7596003)(26005)(73392003)(82202003)(336012)(83380400001)(3480700007)(2906002);DIR:OUT;SFP:1101;
X-ExternalRecipientOutboundConnectors: e7f5a0d7-8564-45b1-bdf4-5a28b6136195
X-MS-Exchange-ForwardingLoop: queries@avon.com;e7f5a0d7-8564-45b1-bdf4-5a28b6136195
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 22 May 2022 05:15:41.4434
(UTC)
X-MS-Exchange-CrossTenant-Network-Message-Id: c2d8449d-08a7-43f0-4988-08da3bb21daf
X-MS-Exchange-CrossTenant-Id: e7f5a0d7-8564-45b1-bdf4-5a28b6136195
X-MS-Exchange-CrossTenant-AuthSource: BN8NAM11FT045.eop-nam11.prod.protection.outlook.com
X-MS-Exchange-CrossTenant-AuthAs: Anonymous
X-MS-Exchange-CrossTenant-FromEntityHeader: Internet
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BN3PR07MB2691
X-IS-SYNAQ-MX: mail-mw2nam12on2069.outbound.protection.outlook.com ([40.107.244.69] helo=NAM12-MW2-obe.outbound.protection.outlook.com)
Authentication-Results: accept;
   iprev=pass (mail-mw2nam12on2069.outbound.protection.outlook.com) smtp.remote-ip=40.107.244.69;
   spf=pass smtp.mailfrom=avon.com;
   dkim=pass header.d=avonglobal.onmicrosoft.com header.s=selector2-avonglobal-onmicrosoft-com header.a=rsa-sha256;
   dkim=pass header.d=gmail.com header.s=20210112 header.a=rsa-sha256;
   dmarc=pass header.from=avon.com dmarc.action=accept
X-SYNAQ-Pinpoint-Information: Please contact SYNAQ for more information
X-SYNAQ-Pinpoint-ID: 1nsdwE-000b99-BC
X-SYNAQ-Pinpoint: No virus infections found
X-SYNAQ-Pinpoint-SpamCheck: not spam (whitelisted), SpamAssassin (not cached,
   score=5.807, required 5, BAYES_99 4.00, DEAR_SOMETHING 0.99,
   DKIM_SIGNED 0.10, DKIM_VALID -0.10, DKIM_VALID_AU -0.10,
   DMARC_ACCEPT -0.05, HEADER_FROM_DIFFERENT_DOMAINS 0.25,
   HTML_MESSAGE 0.00, RCVD_IN_MSPIKE_H2 -0.00, SPF_HELO_SOFTFAIL 0.73,
   SPF_PASS -0.00, T_SCC_BODY_TEXT_LINE -0.01)
X-Pinpoint-From: queries@avon.com
X-Spam-Flag: NO
Return-Path: queries@avon.com
X-MS-Exchange-Organization-Network-Message-Id: 62769f77-9cac-4886-4223-08da3bb22884
X-MS-Exchange-Organization-AVStamp-Enterprise: 1.0
X-MS-Exchange-Organization-AuthSource: HMC-DAG05SVR01.he.businessgateway.co.za
X-MS-Exchange-Organization-AuthAs: Anonymous
MIME-Version: 1.0

--0000000000008d3ea305df92d0d7
Content-Type: text/plain; charset="UTF-8"

Dear sir/ ma'am please check the order number of #7458646 went to the wrong
person it did not arrive to me and please cancel the balance on the account
number of 122470378 so that i can use my credit to enter my orders please
answer me ASAP THANK YOU

--0000000000008d3ea305df92d0d7
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Dutf-8"><d=
iv dir=3D"auto">Dear sir/ ma'am please check the order number of #7458646 w=
ent to the wrong person it did not arrive to me and please cancel the balan=
ce on the account number of 122470378 so that i can use my credit to enter =
my orders please answer me ASAP THANK YOU</div>

--0000000000008d3ea305df92d0d7--`);

if (email!==undefined&&email!==null){
	for(var mpart=email.nextPart();mpart!==undefined&&mpart!==null;mpart=email.nextPart()) {
		var filename=mpart.filename();
		if (filename===""&&(mpart.isContentType("text/plain")||mpart.isContentType("text/html"))) {
			//console.log(mpart.body.readAll());
			break;
		}
	}
	email.close();
}

var CUSTOMERREFKEY='122480904';
var COMMENTERID='1122';
var COMMENT='test comment';
var ariescomment=`agentid:${COMMENTERID},note:${COMMENT}`;
                                var CUTSOMERREFKEY='122480904';
								try {
                                    ssn.send(`http://127.0.0.1:1030/AvonOneAriesIntegration/AriesNotes?account=${CUSTOMERREFKEY}`,{"method":"POST","headers":{"Content-Type":"application/json"},"body":JSON.stringify({"note":ariescomment})});    
                                } catch (error) {
                                    console.log(error.toString());
                                }
