ssn.fs().mkdir("/","C:/projects/inovo/avon/avonone");
ssn.fs().mkdir("/kweauth","C:/github/kweauth");
ssn.fs().mkdir("kwetl","C:/GitHub/kwetl");
ssn.fs().mkdir("kweslnk","C:/projects/slnks");
ssn.fs().mkdir("collect","C:/projects/collect");
ssn.env().set("kwetl-path","/kwetl");
ssn.dbms().register("avonone","sqlserver","server=localhost; database=ONER; user id=ONER; password=ONER");