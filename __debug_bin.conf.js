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

kwe.fs().mkdir("/gendocs","./gendocs");
try {
	eval(kwe.send("/gendocs/gendocs.js").readAll());
} catch(e){
	console.log(e.message);
}*/

kwe.dbms().registerConnection("b1","kwesqlite",":memory:");

var a=kwe.dbms().query({"alias":"b1","query":"select 42"});
println(kwemethods(a));
console.log(a.json());


//Capture bible
/*try {
	
	var translationsr=kwe.Send("https://getbible.net/v2/translations");

	var translationline="";
	kwe.FS().MKDIR("/bible/"+translationkey,"")
	kwe.FS().SET("/bible/translations.json",kwe.Send("https://getbible.net/v2/translations.json"));
	while((translationline=translationsr.Readln())!==""){
		if(!(translationline=translationline.trim()).startsWith("#")&&translationline!=="") {
			var translationkey=translationline.split('\t')[5];
			console.log(translationkey);
			
			kwe.FS().SET("/bible/"+translationkey+"/books.json",kwe.Send("https://getbible.net/v2/"+translationkey+"/books.json"));
			for(var i=1;i<=66;i++){
				kwe.FS().SET("/bible/"+translationkey+"/"+i+".json",kwe.Send("https://getbible.net/v2/"+translationkey+"/"+i+".json"));
			}
		}
		if (translationkey==="aov")	{
			break;
		}	
	}
	
} catch(e){
	console.log(e.message);
}*/

kwe.fs().mkdir("/movies","D:/movies");

try {
	//eval(kwe.Send("/movies/schedule.js").ReadAll());
} catch(e){
	console.log(e.message);
}
