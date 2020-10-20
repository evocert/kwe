package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/evocert/kwe/listen"
)

func main() {
	cancelChan := make(chan os.Signal, 2)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	listen.Listening().Listen(":1002", false)

	<-cancelChan
	os.Exit(0)
}
