package chnls

import (
	"io"
	"net/http"
	"os"

	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/ws"
	"github.com/gorilla/websocket"
)

/*Channel -
 */
type Channel struct {
	rqsts  map[*Request]*Request
	objmap map[string]interface{}
	lstnr  *listen.Listener
}

//Listener - *listen.Listener listener for Channel
func (chnl *Channel) Listener() *listen.Listener {
	if chnl.lstnr == nil {
		chnl.lstnr = listen.NewListener(chnl)
	}
	return chnl.lstnr
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
	func() {
		if wsrw := ws.NewReaderWriter(wscon); wsrw != nil {
			defer wsrw.Close()
			//var rruns = make([]rune, 1024)
			//var rrunsi = 0

			//var tmpbuf = iorw.NewBuffer()
			//for {
			//	if wsrw.CanWrite() {
			//		io.Copy(tmpbuf, wsrw)
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

			//	}
			//	if wsrw.CanRead() {
			//		wsrw.Print(tmpbuf.String())
			//		tmpbuf.Clear()
			//		wsrw.Flush()
			//	} else {
			//		break
			//	}
			//}
			//tmpbuf.Close()
			//tmpbuf = nil
			chnl.ServeRW(wsrw, wsrw, a...)

		}
	}()
}

//ServeRW - serve Reader Writer
func (chnl *Channel) ServeRW(r io.Reader, w io.Writer, a ...interface{}) {
	if len(a) > 0 {
		a = append([]interface{}{r, w}, a...)
	} else {
		a = []interface{}{r, w}
	}
	if rqst, interrupt := newRequest(chnl, a...); rqst != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
				rqst.Close()

			}()
			rqst.execute(interrupt)
		}()
		rqst = nil
	}
}

//Stdio - os.Stdout, os.Stdin
func (chnl *Channel) Stdio(out *os.File, in *os.File, err *os.File, a ...interface{}) {
	chnl.ServeRW(in, out, a...)
}

//DefaultServeHTTP - helper to perform dummy ServeHttp request on channel
func (chnl *Channel) DefaultServeHTTP(w io.Writer, method string, url string, body io.Reader) {
	if rhttp, rhttperr := http.NewRequest(method, url, body); rhttperr == nil {
		if rhttp != nil {
			var whttp = NewResponse(w, rhttp)
			chnl.ServeHTTP(whttp, rhttp)
		}
	}
}

var gblchnl *Channel

func GLOBALCHNL() *Channel {
	if gblchnl == nil {
		gblchnl = NewChannel()
	}
	return gblchnl
}
