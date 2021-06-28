var base="C:/GitHub/kwe";
var root=base+"";

var lgntscopeglobal="/g/";
var lgntscopelocal="/p/";
var lgntinherits="==";
var lgntreturns="<<";
var lgntdatapaths="/";
var lgntpackagepaths="/-/";
var lgntcomments="/../";

//type
var lgntfields="/f/";
var lgntmethods="/m/";

var prep={
        "/legent":{
            "/scope-global":lgntscopeglobal,
            "/scope-local":lgntscopelocal,
            "/inherits":lgntinherits,
            "/returns":lgntreturns,
            "/comments":lgntcomments,
            "/fields":lgntfields,
            "/methods":lgntmethods,
            "/package-paths":lgntpackagepaths
        }};
        prep[lgntpackagepaths]=[];
var pathspckgs={};
var files={};
var crntpreplvl=0;
var preplvl=-1;
var preppath="";
var currentpath="";
var lastpreppath=".";

_fsutils.FIND(root).filter(function(e){
    return e.AbsolutePath().endsWith(".go");
}).sort(function(a,b){
    var apath=a.AbsolutePath().substring(base.length+1);
    var bpath=b.AbsolutePath().substring(base.length+1);
    return apath.split("/").length-bpath.split("/").length || // sort by length, if equal then
         apath.localeCompare(bpath);    // sort by dictionary order
}).forEach(function(e,ei,elines) {
    currentpath=e.AbsolutePath().substring(base.length+1);
    preppath=currentpath.lastIndexOf("/")>0?currentpath.substring(0,currentpath.lastIndexOf("/")+1):"/";
    if(lastpreppath!==preppath) {
        lastpreppath=preppath;
    }
    var comments={};
    var commentsrti=-1;
    var lastcommentsrti=-1;
    var commentslbl="";
    var cmnts=[];
    var pkgobj;
    var data;
    var package="";
    var typedefln="";
    var typedeflni=-1;

    var srclines=_fsutils.CAT(root+"/"+currentpath).Readlines();
    srclines.forEach(function(ln,lni,lns){
        if ((ln=ln.trim())!==""){
            if(package==="" && ln.startsWith("package ")&&(ln=ln.substring("package ".length).trim())!=="") {
                package=ln;
                pkgobj=null;
                if(pathspckgs[lastpreppath]===undefined) {
                    pathspckgs[lastpreppath]={"package":package,"package-path":""};
                }
                var pckgthsfound=[];
                (lastpreppath.endsWith("/")?lastpreppath.substring(0,lastpreppath.length-1):lastpreppath).split("/").forEach(function(tstpth,tstpthi,tstpths){
                    var testpath= tstpths.slice(0,tstpthi+1).join("/")+"/";
                    var testpathspckgs=pathspckgs[testpath];
                    var crntpckg="";
                    if(testpathspckgs===undefined) {
                        pathspckgs[testpath]={"package":tstpth,"package-path":""};
                        testpathspckgs=pathspckgs[testpath]; 
                    }
                    crntpckg=testpathspckgs["package"];
                    pckgthsfound.push(crntpckg);
                    if (testpathspckgs["package-path"]!==pckgthsfound.join(".")){
                        testpathspckgs["package-path"]=pckgthsfound.join(".")
                    }
                });
                var prvtstpkg;
                
                pckgthsfound.forEach(function(pkgnm,pkgnmi) {
                    if(pkgnmi==0) {
                        if (prep[pkgnm]===undefined){
                            prvtstpkg=(prep[pkgnm]={});
                        } else {
                            prvtstpkg=prep[pkgnm];
                        }
                    } else if (prvtstpkg!==undefined && typeof prvtstpkg==="object") {
                        if (prvtstpkg[pkgnm]===undefined){
                            prvtstpkg[pkgnm]={}
                        } 
                        prvtstpkg=prvtstpkg[pkgnm];
                    }
                    if(pkgnmi===(pckgthsfound.length-1)&&prvtstpkg!==undefined && typeof prvtstpkg==="object") {
                        pkgobj=prvtstpkg;
                    }
                });
                prep[lgntpackagepaths].push(pckgthsfound.join("."));
            } else if (package!=="" && pkgobj!==null && typeof pkgobj==="object") {
                if(typedeflni==-1) {
                    if ((ln.startsWith("type ") && ln.endsWith("{"))||(ln.startsWith("func ") && ln.endsWith("{"))) {
                        typedeflni=lni;
                        typedefln=ln;
                    } else {
                        typedefln="";
                    }
                } else if(typedeflni>-1&&ln==="}"){
                    if (pkgobj[lgntdatapaths]===undefined) {
                        pkgobj[lgntdatapaths]={};
                        pkgobj[lgntdatapaths][lgntscopeglobal]={}
                        pkgobj[lgntdatapaths][lgntscopelocal]={}
                    }
                    data=pkgobj[lgntdatapaths];
                    var cmnti=typedeflni;
                    var typelni=typedeflni;
                    typedeflni=-1;
                    var typecmnts=[];
                    cmnti--;
                    while(cmnti>=0&&(ln=lns[cmnti].trim()).startsWith("//")){
                        typecmnts.push(ln.substring("//".length));
                        cmnti--;
                    }

                    var type="";
                    var typename="";
                    var typebase="";
                    var args=[];
                    var returntypes=[];
                    var typeowner="";
                    var srcln=typedefln;
                    if (srcln.startsWith("type ") && srcln.endsWith("{") && (srcln=srcln.substring("type ".length,srcln.lastIndexOf("{")).trim())!=="") {
                        type="type";
                        typename=srcln.trim();
                        typename=typename.substring(0,typename.indexOf(" ")).trim();
                        typebase=srcln.substring(srcln.indexOf(" ")).trim();
                        typelni++;
                    } else if (srcln.startsWith("func ") && srcln.endsWith("{") && (srcln=srcln.substring("func ".length,srcln.lastIndexOf("{")).trim())!==""&& srcln.indexOf("(")>=0) {
                        type="func";
                        if ((typename=srcln.substring(0,srcln.indexOf("(")).trim())==="") {
                            srcln=srcln.substring(srcln.indexOf("(")+1).trim();
                            typebase=srcln.substring(0,srcln.indexOf(")")).trim();
                            typeowner=typebase.substring(typebase.indexOf(" ")).trim();
                            srcln=srcln.substring(srcln.indexOf(")")+1).trim();
                        }
                        if ((typename=srcln.substring(0,srcln.indexOf("(")).trim())!=="") {
                            srcln=srcln.substring(srcln.indexOf("(")+1).trim();
                        }
                        if(typename!=="") {
                            if (srcln.indexOf(")")) {
                                var argscount=0;
                                srcln.substring(0,srcln.indexOf(")")).trim().split(",").forEach(function(arg,argi,rgs){
                                    arg=arg.trim();
                                    if(arg.indexOf(" ")>0) {
                                        var argsnme=arg.substring(0,arg.indexOf(" ")).trim();
                                        var argtype=arg.substring(arg.indexOf(" ")+1).trim();
                                        var obj={};
                                        obj[argsnme]=argtype;
                                        args.push(obj);
                                    } else {
                                        var argsnme=arg.trim();
                                        argi++;
                                        while(argi<rgs.length-1 && ((arg=rgs[argi]).trim()).indexOf(" ")<=0) {
                                            argi++;
                                            continue;
                                        }
                                        if(arg.indexOf(" ")>0) {
                                            var argtype=arg.substring(arg.indexOf(" ")+1).trim();
                                            var obj={};
                                            obj[argsnme]=argtype;
                                            args.push(obj);
                                        } else {
                                            var obj={};
                                            obj[argscount+""]=argtype;
                                            argscount++;
                                            args.push(obj);
                                        }
                                    }
                                });
                                var rsltcnount=0;
                                ((srcln=srcln.substring(srcln.indexOf(")")+1).trim())?srcln.startsWith("(")&&srcln.endsWith(")")?srcln.substring(1,srcln.length-1):srcln:"").split(",").forEach(function(arg,argi,rgs){
                                    arg=arg.trim();
                                    if(arg.indexOf(" ")>0) {
                                        var argsnme=arg.substring(0,arg.indexOf(" ")).trim();
                                        var argtype=arg.substring(arg.indexOf(" ")+1).trim();
                                        var obj={};
                                        obj[argsnme]=argtype;
                                        returntypes.push(obj);
                                    } else {
                                        var argsnme=arg.trim();
                                        argi++;
                                        while(argi<rgs.length-1 && ((arg=rgs[argi]).trim()).indexOf(" ")<=0) {
                                            argi++;
                                            continue;
                                        }
                                        if(arg.indexOf(" ")>0) {
                                            var argtype=arg.substring(arg.indexOf(" ")+1).trim();
                                            var obj={};
                                            obj[argsnme]=argtype;
                                            returntypes.push(obj);
                                        } else {
                                            var obj={};
                                            obj[rsltcnount+""]=argtype;
                                            rsltcnount++;
                                            returntypes.push(obj);
                                        }
                                    }
                                });
                            }
                        }
                    }
                    if(type=="func") {
                        var datascpobj;
                        if ((typename.charAt(0)+"").toLowerCase()===(typename.charAt(0)+"")) {
                            datascpobj=data[lgntscopelocal];
                        } else {
                            datascpobj=data[lgntscopeglobal];
                        }
                        var objfnc={"type":type,"owner":typeowner,"parameters":args};
                        objfnc[lgntcomments]=typecmnts.slice(0)
                        objfnc[lgntreturns]=returntypes.slice(0);

                        if (typeowner!=="") {
                            var ownerref=typeowner;
                            if((ownerref=ownerref.startsWith("*")?ownerref.substring(1).trim():ownerref.trim())!==""){
                                if ((ownerref.charAt(0)+"").toLowerCase()===(ownerref.charAt(0)+"")) {
                                    datascpobj=data[lgntscopelocal];
                                } else {
                                    datascpobj=data[lgntscopeglobal];
                                }
                                if(datascpobj[ownerref+""]!==undefined && typeof datascpobj[ownerref+""]==="object") {
                                    if(datascpobj[ownerref+""][lgntmethods]===undefined){
                                        datascpobj[ownerref+""][lgntmethods]={};
                                        datascpobj[ownerref+""][lgntmethods][lgntscopeglobal]={};
                                        datascpobj[ownerref+""][lgntmethods][lgntscopelocal]={};
                                    }
                                    if ((typename.charAt(0)+"").toLowerCase()===(typename.charAt(0)+"")) {
                                        datascpobj[ownerref+""][lgntmethods][lgntscopelocal][typename+""]=objfnc;
                                    } else {
                                        datascpobj[ownerref+""][lgntmethods][lgntscopeglobal][typename+""]=objfnc;
                                    }
                                }
                            }
                        } else {
                            datascpobj[typename+""]=objfnc;
                        }
                    } else if (type!="") {
                        var datascpobj;
                        if ((typename.charAt(0)+"").toLowerCase()===(typename.charAt(0)+"")) {
                            datascpobj=data[lgntscopelocal];
                        } else {
                            datascpobj=data[lgntscopeglobal];
                        }
                        if (typeof datascpobj[typename+""] !=="object") {
                            var members={};
                            members[lgntscopeglobal]={};
                            members[lgntscopelocal]={};
                            var inherits=[];
                            if(typelni<lni){
                                while(typelni<lni) {
                                    if((ln=lns[typelni].trim())!==""){
                                        if(!ln.startsWith("//")) {
                                            if(ln.indexOf(" ")>0) {
                                                var memnme=ln.substring(0,ln.indexOf(" ")).trim();
                                                ln=ln.substring(ln.indexOf(" ")+1).trim();
                                                if(memnme!==""&&ln!==""){
                                                    while(ln.indexOf("  ")>0) {
                                                        ln=ln.replace("  "," ");
                                                    }
                                                    if(ln.indexOf("//")>-1){
                                                        ln=ln.substring(0,ln.indexOf("//")).trim();
                                                    }
                                                    
                                                    if ((memnme.charAt(0)+"").toLowerCase()===(memnme.charAt(0)+"")) {
                                                        members[lgntscopelocal][memnme]=ln;
                                                    } else {
                                                        members[lgntscopeglobal][memnme]=ln;
                                                    }
                                                }
                                            }else{
                                                inherits.push(ln);
                                            }
                                        }
                                    }
                                    typelni++;
                                }
                            }
                            var objtpe={"type":type,"base":typebase};
                            
                            objtpe[lgntfields]=members;
                            objtpe[lgntinherits]=inherits.slice(0);
                            objtpe[lgntcomments]=typecmnts.slice(0);
                            datascpobj[typename+""]=objtpe;               
                        }
                    }
                }
            } 
        }
    });
}.bind(this));
request.ResponseHeader().Set("Content-Type","application/json");
_fsutils.SET("./gendocs/codedefs.json",JSON.stringify(prep))
print(JSON.stringify(prep));