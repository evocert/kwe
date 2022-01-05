function buildgo(goosgoarchsarr,codepath,outputpath,upxpath) {
    if (goosgoarchsarr !==undefined && typeof goosgoarchsarr === "object" && Array.isArray(goosgoarchsarr)) {
        
        if ((codepath!==undefined && typeof codepath === "string" && codepath!=="") && (outputpath!==undefined && typeof outputpath === "string" && outputpath!=="")) {
            var cmd=kwe.command("cmd");
            goosgoarchsarr.forEach((goosarch)=>{
                var goosarcharr=goosarch.trim().split("/");
                var goos=goosarcharr[0].trim();
                var goarch=goosarcharr[1].trim();
                if (goarch==="" || goos==="" || goos==="ios") return;
                console.log(`${goos}: ${goarch}`);
                try {
                    cmd.setReadTimeout(10000,100);
                    cmd.readAll();
                    cmd.println("SET GOOS="+goos);
                    cmd.println("SET GOARCH="+goarch);
                    if (goos==="ios") {
                        cmd.println("SET CGO_ENABLED=1");
                    }
                    var binpath=outputpath+"";
                    cmd.println(`go build -v -ldflags "-w -s" -o ${binpath}_${goos}_${goarch}${(goos==="windows"?".exe":(goos==="js"?".wasm":(goos==="darwin"?".dmg":"")))} ${codepath}`);
                    cmd.println("echo finit");
                    for(var ln = cmd.readln();!ln.endsWith("finit");ln= cmd.readln()){
                        if (ln!=="") {
                            console.log(ln);
                        }
                    }
                } catch (error) {
                    console.log("error[",goos,":",goarch,"]:",error.toString());
                }
                
                console.log("done build - ",goos,":",goarch);	
            });
            cmd.close();
            console.log("done - build");
        }

        if ((upxpath!==undefined && typeof upxpath === "string" && upxpath!=="") && (outputpath!==undefined && typeof outputpath === "string" && outputpath!=="")){
            var cmd=kwe.command("cmd");
            goosgoarchsarr.forEach((goosarch)=>{
                var goosarcharr=goosarch.trim().split("/");
                var goos=goosarcharr[0].trim();
                var goarch=goosarcharr[1].trim();
                if (goos==="" || goarch==="") return;
                console.log(`${goos}: ${goarch}`);
                
                try {
                    cmd.setReadTimeout(10000,100);
                    cmd.readAll();
                    var binpath=outputpath+"";
                    cmd.println(`${upxpath} ${binpath}_${goos}_${goarch}${(goos==="windows"?".exe":(goos==="js"?".wasm":(goos==="darwin"?".dmg":"")))}`)
                    cmd.println("echo finit");
                    for(var ln = cmd.readln();!ln.endsWith("finit");ln= cmd.readln()){
                        if (ln!=="") {
                            console.log(ln);
                        }
                    }
                } catch (error) {
                    console.log("error[",goos,":",goarch,"]:",error.toString());
                }
                console.log("done upx - ",goos,":",goarch);	
            });
            cmd.close();
            console.log("done - upx");
        }
    }
}