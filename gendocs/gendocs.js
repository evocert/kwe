var base="C:/GitHub/kwe";
var root=base+"";

var prep={};
var files={};
var crntpreplvl=0;
var preplvl=-1;
var preppath="";
var currentpath="";

function prepFile(preppath,preppathlvl,src){
    if (src!==undefined) {
        var srclines=src.Readlines();
        var ln=srclines.length===0?"":srclines[0].trim();
        srclines=srclines.slice(1);
        var packagenme="";
        var crntpckg=null;
        if (ln.startsWith("package ")) {
            if((ln=ln.substring("package ".length).trim())!==""){
                packagenme=ln;
                if (typeof files[preppath+"/"]["package"]!=="string") {
                    files[preppath+"/"]["package"]=packagenme;
                }

                if (files[preppath+"/"]["package"]===packagenme) {
                    //if(preppath.split("/").length===1&&preppath===""){
                    //    if (typeof prep[packagenme]!=="object") {
                    //        prep[packagenme]={};
                    //        crntpckg=prep[packagenme];
                    //    } 
                    //} else 
                    if(preppath.split("/").length>0){
                        var prvpckg=null;
                        preppath.split("/").forEach(function(pth,pthi,pths) {
                            var tstpath=pths.slice(0,pthi).join("/");
                            println(tstpath+"/");
                            //var tstpckg=files[tstpath+"/"]["package"];
                            try {
                                if (typeof files[tstpath+"/"]["package"]==="string") {
                                    var tstpckg=files[tstpath+"/"]["package"];
                                    if(pthi===0) {
                                        if(typeof (prvpckg=prep[tstpckg])!=="object"){
                                            prvpckg=(prep[tstpckg]={});
                                        }
                                    } else  if(typeof prvpckg === "object" && typeof prvpckg[tstpckg]!=="object"){
                                        prvpckg=(prvpckg[tstpckg]={});
                                    }
                                }
                                crntpckg=prvpckg;
                            } catch (e) {
                                println("e:",e.toString());
                            }
                        });
                    }

                    if (crntpckg!==null && typeof crntpckg==="object") {
                        var data=null;
                        if(typeof crntpckg["_"]!=="object") {
                            data=(crntpckg["_"]={});   
                        } else {
                            data=crntpckg["_"];
                        }
                        var skiplines=0;
                        srclines.forEach(function(srcln,srcli,srclns){
                            if (skiplines>0) {
                                skiplines--;
                                return;
                            }
                            if((srcln=srcln.trim())!=="") {
                                if ((srcln.startsWith("type ") && srcln.endsWith("{"))||(srcln.startsWith("func ") && srcln.endsWith("{"))) {
                                    
                                    var comments=[];
                                    var contents=[];
                                    var cmnti=srcli-1;
                                    var type="";
                                    var typename="";
                                    var typebase="";
                                    var args=[];
                                    var returntypes=[];
                                    var typeowner="";
                                    if (srcln.startsWith("type ") && srcln.endsWith("{") && (srcln=srcln.substring("type ".length,srcln.lastIndexOf("{")).trim())!=="") {
                                        type="type";
                                        typename=srcln.trim();
                                        typename=typename.substring(0,typename.indexOf(" ")).trim();
                                        typebase=srcln.substring(srcln.indexOf(" ")).trim();
                                        typebase=typebase.substring(0,typebase.length-1).trim();
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
                                    while(cmnti>=0) {
                                        if ((srcln=srclines[cmnti].trim()).startsWith("//")) {
                                            srcln=srcln.substring("//".length).trim();
                                            cmnti--; 
                                                
                                            comments.unshift(srcln);
                                            if (srcln.startsWith(typename+" ")){
                                                break;
                                            } else {
                                                comments=[]
                                                break;
                                            }
                                        } else {
                                            break;
                                        }
                                    }              
                                    while ((srcln=srclns[srcli+skiplines+1].trim())!=="}") {
                                        contents.push(srcln);
                                        skiplines++;
                                    }
                                    if(srcln==="}") {
                                        skiplines++;
                                        if(type=="func") {
                                            if (typeof data[typename+""] !=="object") {
                                                var objfnc={"type":type,"owner":typeowner,"parameters":args,"comments":comments.slice(0),"result":returntypes.slice(0)};
                                                data[typename+""]=objfnc;
                                            }
                                        } else if (type!="") {
                                            if (typeof data[typename+""] !=="object") {
                                                var objtpe={"type":type,"base":typebase,"comments":comments.slice(0)};
                                                data[typename+""]=objtpe
                                               
                                            }
                                        }
                                    } else {
                                        contents=[];
                                        commends=[];
                                        skiplines=0;
                                    }
                                }
                            }
                        });
                    }
                }
            }
        }
    }
}


_fsutils.FIND(root).filter(function(e){
    return e.AbsolutePath().endsWith(".go");
}).sort(function(a,b){
    var apath=a.AbsolutePath().substring(base.length+1);
    //var alvl=apath.split("/").length;
    var bpath=b.AbsolutePath().substring(base.length+1);
    //var blvl=bpath.split("/").length;
    return apath.split("/").length-bpath.split("/").length || // sort by length, if equal then
         apath.localeCompare(bpath);    // sort by dictionary order
}).forEach(function(e,ei,elines) {
    currentpath=e.AbsolutePath().substring(base.length+1);
    preppath=currentpath.lastIndexOf("/")>0?currentpath.substring(0,currentpath.lastIndexOf("/")+1):"/";
    var prepfiles=typeof files[preppath]==="object"?files[preppath]:(files[preppath]={"list":[]}); 
    prepfiles.list.push(currentpath.lastIndexOf("/")>0?currentpath.substring(currentpath.lastIndexOf("/")+1):currentpath);   
    if(preppath.startsWith("iorw/")){
        prepFile(preppath.substring(0,preppath.length-1),preppath==="/"?0:preppath.substring(0,preppath.length-1).split("/").length,_fsutils.CAT(root+"/"+currentpath));
    }
}.bind(this));
request.ResponseHeader().Set("Content-Type","application/json");
print(JSON.stringify(prep));