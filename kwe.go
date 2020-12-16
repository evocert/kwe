package main

import (
	"os"

	"github.com/evocert/kwe/service"
)

func main() {
	/*debug.SetGCPercent(25)
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
	os.Exit(0)*/
	service.RunService(os.Args...)
}
