ssn.fs().mkdir("/movies","D:/movies");
ssn.fs().mkdir("/kweauth","C:/github/kweauth");
ssn.fs().mkdir("kwetl","C:/GitHub/kwetl");
ssn.fs().mkdir("kweslnk","C:/projects/slnks");
ssn.fs().mkdir("collect","C:/projects/collect");
ssn.env().set("kwetl-path","/kwetl");
ssn.dbms().register("tsql","sqlserver","server=localhost;user id=PTOOLS; password=PTOOLS");