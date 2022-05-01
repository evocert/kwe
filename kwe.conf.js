ssn.fs().mkdir("/kweauth","C:/github/kweauth");
ssn.fs().mkdir("kwetl","C:/GitHub/kwetl");
ssn.fs().mkdir("kweslnk","C:/projects/slnks");
ssn.fs().mkdir("collect","C:/projects/collect");
ssn.fs().mkdir("/oner","C:/projects/oner")
ssn.fs().mkdir("/tvseries","C:/tvseries");
ssn.fs().mkdir("/movies","C:/movies");
ssn.fs().mkdir("/music","C:/music");
ssn.env().set("kwetl-path","/kwetl");
ssn.dbms().register("avonone","sqlserver","server=localhost; database=ONER; user id=ONER; password=ONER");
ssn.dbms().register("oner-api","sqlserver","server=localhost; database=ONER; user id=ONER; password=ONER");