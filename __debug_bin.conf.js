kwe.Caching().Put("a",{"b":[1,2,{"f":[7,8,9,18]},3,{"g":[17,18,19,118]},4]});
console.log(kwe.Caching().String());
if (kwe.Caching().ExistsAt("a","b",2)) {
	console.log(kwe.Caching().At("a","b",2).String());
	kwe.Caching().ClearAt("a","b",2)
	console.log(kwe.Caching().At("a","b",2).String());
}
console.log(kwe.Caching().String());
if (kwe.Caching().ExistsAt("a","b",2)) {
	kwe.Caching().CloseAt("a","b",2)
	console.log(kwe.Caching().At("a","b",2));
	console.log(kwe.Caching().String());
}

kwe.FS().MKDIR("/kweauth","C:/github/kweauth");

kwe.FS().MKDIR("/gendocs","./gendocs");
try {
	eval(kwe.Send("/gendocs/gendocs.js").ReadAll());
} catch(e){
	console.log(e.message);
}

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

kwe.FS().MKDIR("/movies","D:/movies");

try {
	//eval(kwe.Send("/movies/schedule.js").ReadAll());
} catch(e){
	console.log(e.message);
}
