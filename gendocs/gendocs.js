var base="C:/GitHub/kwe";
var root=base+"";
var strgen="";
strgen+=("[");
var filters=_fsutils.FIND(root).filter(function(e){
    return e.AbsolutePath().endsWith(".go");
})
var filtersl=filters.length;
filters.forEach(function(e,ei) {
    strgen+=(JSON.stringify(e.AbsolutePath().substring(base.length+1)));
    if(ei<filtersl-1) {
        strgen+=(",");
    }
}.bind(this));
request.ResponseHeader().Set("Content-Type","application/json");
strgen+=("]");
_fsutils.SET(base+"/gendocs/codelist.json",strgen);