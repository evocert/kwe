package main

import (
	"os"

	//_ "github.com/evocert/kwe/database/db2"

	_ "github.com/evocert/kwe/alertify"
	_ "github.com/evocert/kwe/fonts/material"

	_ "github.com/evocert/kwe/typescript"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	service.RunService(os.Args...)
}

/*requesting.DefaultRequestInvoker = httpr.RequestInvoker
requesting.DefaultResponseInvoker = httpr.ResponseInvoker
resources.GLOBALRSNG().RegisterEndpoint("/", "./")
resources.GLOBALRSNG().FS().MKDIR("/tools", "")
resources.GLOBALRSNG().FS().SET("/tools/telnet.js", mytelnet)
chnls.GLOBALCHNL().ServeRequest("/active:config.js", nil, nil)
go http.ListenAndServe(":1002", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	chnls.GLOBALCHNL().ServeRequest(r, httpr.RequestInvoker, w, httpr.ResponseInvoker)
}))

// Listen for incoming connections.
l, err := net.Listen("tcp", "0.0.0.0:1234")
if err != nil {
	fmt.Println("Error listening:", err.Error())
	os.Exit(1)
}
// Close the listener when the application closes.
defer l.Close()
fmt.Println("Listening on " + "0.0.0.0:1234")
for {
	// Listen for an incoming connection.
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	// Handle connections in a new goroutine.
	go func(cn net.Conn) {
		defer cn.Close()
		chnls.GLOBALCHNL().ServeRequest("/tools/telnet.js", cn, cn)
	}(conn)
}*/
