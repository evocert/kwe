package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/evocert/kwe/chnls"
	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/web"
	"github.com/evocert/kwe/ws"
)

func main() {
	go func() {
		b := false
		lck := &sync.Mutex{}
		for {
			time.Sleep(time.Second * 10)
			func() {
				lck.Lock()
				defer lck.Unlock()
				if !b {
					b = true
					go func() {
						runtime.GC()
						b = false
					}()
				}
			}()
		}
	}()
	//database.GLOBALDBMS().RegisterConnection("mydb", "postgres", "user=postgres password=1234!@#$qwerQWER host=skullquake.dedicated.co.za port=5432 dbname=postgres sslmode=disable")
	database.GLOBALDBMS().RegisterConnection("psg", "postgres", "user=postgres password=n@n61ng@ dbname=postgres sslmode=disable host=127.0.0.1 port=5433")
	database.GLOBALDBMS().RegisterConnection("psgrmt", "remote", "http://127.0.0.1:1002/dbms-psg/.json")
	cancelChan := make(chan os.Signal, 2)
	// catch SIGTERM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	if args := os.Args; len(args) == 3 {
		resources.GLOBALRSNG().RegisterEndpoint("/", args[1])
		chnls.GLOBALCHNL().Listener().Listen(args[2])
	} else if args := os.Args; len(args) == 2 {
		chnls.GLOBALCHNL().Listener().Listen(args[1])
	} else {
		resources.GLOBALRSNG().RegisterEndpoint("/", "./")
		resources.GLOBALRSNG().RegisterEndpoint("/cdn/", "https://code.jquery.com/")
		resources.GLOBALRSNG().RegisterEndpoint("/mem/", "")
		resources.GLOBALRSNG().RegisterEndpoint("/dojo/", "https://ajax.googleapis.com/ajax/libs/dojo/1.14.1/dojo/")
		if f, ferr := os.Open("./trythis.html"); ferr == nil {
			resources.GLOBALRSNG().FS().SET("/mem/uiop/string.txt", f)
			resources.GLOBALRSNG().FS().SET("/dojo/uiop/string.html", f)
		}
		chnls.GLOBALCHNL().Listener().Listen(":1002")
	}
	buff := iorw.NewBuffer()
	cnlt := web.NewClient()
	fmt.Println(web.DefaultClient.SendRespondString("http://skullquake.dedicated.co.za/", nil))
	if rw, rwerr := cnlt.SendReceive("ws://127.0.0.1:1002"); rwerr == nil {
		if rw != nil {
			if wsrw, wsrwok := rw.(*ws.ReaderWriter); wsrwok {
				wsrw.Println("#!dbms")
				wsrw.Println(`{"alias":"psg","5555":{"query":"select * from test.tbltest"},"1234":{"query":"select * from test.tbltest"}}`)
				wsrw.Println("#!commit")
				buff.Print(wsrw)
				//wsrw.Println("#!exit")
				wsrw.Close()
				fmt.Println(buff)
				buff.Clear()
			}
		}
	} else {
		fmt.Println(rwerr.Error())
	}
	cnlt.Send("http://127.0.0.1:1002/dbms/.json",
		map[string]string{"Content-Type": "application/json"},
		nil,
		strings.NewReader(`{"alias":"psg","5555":{"query":"select * from test.tbltest"},"1234":{"query":"select * from test.tbltest"}}`),
		buff)
	fmt.Println(buff)
	go func() {
		for {
			time.Sleep(time.Second * 30)
			runtime.GC()
		}
	}()
	<-cancelChan
	os.Exit(0)
}
