package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/web"
)

func main() {

	//database.GLOBALDBMS().RegisterConnection("mydb", "postgres", "user=postgres password=1234!@#$qwerQWER host=skullquake.dedicated.co.za port=5432 dbname=postgres sslmode=disable")
	database.GLOBALDBMS().RegisterConnection("psg", "postgres", "user=postgres password=n@n61ng@ dbname=postgres sslmode=disable host=127.0.0.1 port=5433")
	database.GLOBALDBMS().RegisterConnection("psgrmt", "remote", "http://127.0.0.1:1002/dbms-psg/.json")
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	if args := os.Args; len(args) == 3 {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", args[1])
		listen.Listening().Listen(args[2], false)
	} else if args := os.Args; len(args) == 2 {
		listen.Listening().Listen(args[1], false)
	} else {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./")
		listen.Listening().Listen(":1002", false)
	}
	buff := iorw.NewBuffer()
	cnlt := &web.Client{}
	cnlt.Send("http://127.0.0.1:1002/dbms/.json",
		map[string]string{"Content-Type": "application/json"},
		nil,
		strings.NewReader(`{"alias":"psg","5555":{"query":"select * from test.tbltest"},"1234":{"query":"select * from test.tbltest"}}`),
		buff,
	)
	fmt.Println(buff)
	<-cancelChan
	os.Exit(0)
}
