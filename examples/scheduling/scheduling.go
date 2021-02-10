package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/evocert/kwe/scheduling"
)

func main() {
	if test1 := scheduling.GLOBALSCHEDULES().RegisterSchedule("test1", map[string]interface{}{"Seconds": 2}); test1 != nil {
		for _, sn := range strings.Split("1,2,3,4,5,6,7,8,9,0", ",") {
			if sn != "" {
				test1.AddAction([]interface{}{sn}, func(a ...interface{}) {
					fmt.Println("test action ", a, " ", time.Now())
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
