package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/resources"
)

func main() {
	database.GLOBALDBMS().RegisterConnection("mydb2", "mysql", "mysql:1234!qwer!QWER@tcp(154.0.161.242)/test")
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	if args := os.Args; len(args) == 3 {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", args[1])
		listen.Listening().Listen(args[2], false)
	} else {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./")
		listen.Listening().Listen(":1002", false)
	}
	database.GLOBALDBMS()
	<-cancelChan
	os.Exit(0)
}
