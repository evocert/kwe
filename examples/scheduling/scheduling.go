package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evocert/kwe/scheduling"
)

func main() {
	if test1 := scheduling.GLOBALSCHEDULES().RegisterSchedule("test1", map[string]interface{}{"Seconds": 20}); test1 != nil {
		for _, sn := range strings.Split("1", ",") {
			if sn != "" {
				test1.AddAction([]interface{}{sn}, func(a ...interface{}) {
					fmt.Println("test action ", a)
				})
			}
		}
		test1.Start()
	}

	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	<-cancelChan
	os.Exit(0)
}
