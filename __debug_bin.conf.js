kwe.caching.Put("a",{"b":[1,2,{"f":[7,8,9,18]},3,{"g":[17,18,19,118]},4]});
console.Log(kwe.caching.String());
if (kwe.caching.ExistsAt("a","b",2)) {
	console.Log(kwe.caching.At("a","b",2).String());
	kwe.caching.ClearAt("a","b",2)
	console.Log(kwe.caching.At("a","b",2).String());
}
console.Log(kwe.caching.String());
if (kwe.caching.ExistsAt("a","b",2)) {
	kwe.caching.CloseAt("a","b",2)
	console.Log(kwe.caching.At("a","b",2));
	console.Log(kwe.caching.String());
}

kwe.fs.MKDIR("/gendocs","C:/GitHub/kwe/gendocs");
try {
	eval(kwe.send("/gendocs/gendocs.js").ReadAll());
} catch(e){
	console.Log(e.message);
}

kwe.fs.MKDIR("/movies","D:/movies");

try {
	eval(kwe.send("/movies/schedule.js").ReadAll());
} catch(e){
	console.Log(e.message);
}

kwe.listen(":3332",":3331");

