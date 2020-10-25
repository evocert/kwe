package main

import (
	"os"
	"os/signal"
	runtime "runtime"
	"runtime/debug"
	"syscall"

	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/resources"
)

func main() {
	debug.SetGCPercent(25)
	runtime.GOMAXPROCS(runtime.NumCPU() * 10)
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./")
	if args := os.Args; len(args) == 3 {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", args[1])
		listen.Listening().Listen(args[2], false)
	} else {
		resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "./")
		listen.Listening().Listen(":1002", false)
	}
	<-cancelChan
	os.Exit(0)
}
