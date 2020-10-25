package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/resources"
)

func main() {
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	resources.GLOBALRSNGMANAGER().RegisterEndpoint("/", "D:/mystuff/bcoring")
	listen.Listening().Listen(":1002", false)
	<-cancelChan
	os.Exit(0)
}
