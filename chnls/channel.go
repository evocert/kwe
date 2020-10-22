package chnls

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

/*Channel -
 */
type Channel struct {
	rqsts map[*Request]*Request
}

//NewChannel - instance
func NewChannel() (chnl *Channel) {
	chnl = &Channel{rqsts: map[*Request]*Request{}}
	return
}

func (chnl *Channel) nextRequest() (rqst *Request) {
	rqst = newRequest(chnl)
	return
}

//ServeHTTP - refer http.Handler
func (chnl *Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024); err == nil {
		chnl.ServeWS(conn, r)
	} else {
		chnl.ServeRW(r.Body, w, r)
	}
}

//ServeWS - server websocket Connection
func (chnl *Channel) ServeWS(wscon *websocket.Conn, a ...interface{}) {
	if wsrw := NewWsReaderWriter(wscon); wsrw != nil {
		chnl.ServeRW(wsrw, wsrw, a...)
		wsrw.Close()
	}
}

//ServeRW - Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	if rqst := newRequest(chnl, r, w, a); rqst != nil {
		var dne = make(chan bool, 1)
		go func(d chan<- bool) {
			defer func() {
				if r := recover(); r != nil {
					//fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
				}
				rqst.Close()
				d <- true

			}()
			rqst.execute()
		}(dne)
		<-dne
		rqst = nil
	}
}

//Stdio - os.Stdout, os.Stdin
func (chnl *Channel) Stdio(out *os.File, in *os.File, err *os.File, a ...interface{}) {
	chnl.ServeRW(in, out, a...)
}
