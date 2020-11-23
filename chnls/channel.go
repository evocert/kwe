package chnls

import (
	"io"
	"net/http"
	"os"

	"github.com/evocert/kwe/iorw"
	"github.com/gorilla/websocket"
)

/*Channel -
 */
type Channel struct {
	rqsts  map[*Request]*Request
	objmap map[string]interface{}
}

//NewChannel - instance
func NewChannel() (chnl *Channel) {
	chnl = &Channel{rqsts: map[*Request]*Request{}, objmap: map[string]interface{}{}}
	return
}

func (chnl *Channel) nextRequest() (rqst *Request, interrupt func()) {
	rqst, interrupt = newRequest(chnl)
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
		//var rruns = make([]rune, 1024)
		//var rrunsi = 0
		var tmpbuf = iorw.NewBuffer()
		for {
			if wsrw.CanWrite() {
				io.Copy(tmpbuf, wsrw)
				/*for {
					rn, size, rnerr := wsrw.ReadRune()
					if size > 0 {
						rruns[rrunsi] = rn
						rrunsi++
						if rrunsi == len(rruns) {
							//fmt.Print(string(rruns[:rrunsi]))
							rrunsi = 0
						}
					}
					if rnerr != nil {
						break
					}
				}
				if rrunsi > 0 {
					//fmt.Print(string(rruns[:rrunsi]))
					rrunsi = 0
				}*/

			}
			if wsrw.CanRead() {
				wsrw.Print(tmpbuf.String())
				tmpbuf.Clear()
				wsrw.Flush()
			} else {
				break
			}
		}
		tmpbuf.Close()
		tmpbuf = nil
		//chnl.ServeRW(wsrw, wsrw, a...)
		wsrw.Close()
	}
}

//ServeRW - serve Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	if rqst, interrupt := newRequest(chnl, r, w, a); rqst != nil {
		//var dne = make(chan bool, 1)
		//go func(d chan<- bool) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					//fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
				}
				rqst.Close()
				//d <- true

			}()
			rqst.execute(interrupt)
		}()
		//}(dne)
		//<-dne
		rqst = nil
	}
}

//Stdio - os.Stdout, os.Stdin
func (chnl *Channel) Stdio(out *os.File, in *os.File, err *os.File, a ...interface{}) {
	chnl.ServeRW(in, out, a...)
}
