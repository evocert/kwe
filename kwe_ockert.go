package main

import (
	"os"
	"os/signal"
	runtime "runtime"
	"runtime/debug"
	"syscall"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/database/postgres"
)

func main() {
	debug.SetGCPercent(25)
	runtime.GOMAXPROCS(runtime.NumCPU() * 10)
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	//resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./")
	database.GLOBALDBMS().RegisterConnection("mydb","postgres","user=postgres password=1234!@#$qwerQWER dbname=postgres sslmode=disable")
//"postgres","user=postgres password=1234!@#$qwerQWER dbname=postgres sslmode=disable"
	if args := os.Args; len(args) == 3 {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", args[1])
		listen.Listening().Listen(args[2], false)
	} else {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./www")
		listen.Listening().Listen(":80", false)
	}
	<-cancelChan
	os.Exit(0)
}
