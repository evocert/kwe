package listen

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/evocert/kwe/iorw"
	http2 "golang.org/x/net/http2"
	h2c "golang.org/x/net/http2/h2c"
)

type lstnrserver struct {
	h2s  *http2.Server
	srvr *http.Server
	addr string
}

type ListnerHandler struct {
	ln             net.Listener
	backlog        int
	backcon        chan *ConnHandler
	backLogStarted bool
	lck            *sync.Mutex
}

type ConnHandler struct {
	con      *net.TCPConn
	maxread  int64
	maxwrite int64
}

func newConnHandler(con net.Conn) (cnhdnlr *ConnHandler) {
	if tcpcn, tcpcnok := con.(*net.TCPConn); tcpcnok {
		tcpcn.SetLinger(-1)
		tcpcn.SetReadBuffer(65536)
		tcpcn.SetWriteBuffer(65536)
		cnhdnlr = &ConnHandler{con: tcpcn, maxread: 0, maxwrite: 0}
	}
	return cnhdnlr
}

func (cnhdnlr *ConnHandler) Read(b []byte) (n int, err error) {
	if cnhdnlr.con != nil {
		n, err = cnhdnlr.con.Read(b)
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

func (cnhdnlr *ConnHandler) Write(b []byte) (n int, err error) {
	n, err = cnhdnlr.con.Write(b)
	return
}

func (cnhdnlr *ConnHandler) Readln() (ln string, err error) {
	ln, err = iorw.ReadLine(cnhdnlr)
	return
}

func (cnhndlr *ConnHandler) Readlines() (lines []string, err error) {
	lines, err = iorw.ReadLines(cnhndlr)
	return
}

func (cnhndlr *ConnHandler) ReadAll() (s string, err error) {
	s, err = iorw.ReaderToString(cnhndlr)
	return
}

func (cnhndlr *ConnHandler) Print(a ...interface{}) {
	iorw.Fprint(cnhndlr)
}

func (cnhndlr *ConnHandler) Println(a ...interface{}) {
	iorw.Fprintln(cnhndlr)
}

func (cnhndlr *ConnHandler) Close() (err error) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			err = cnhndlr.con.Close()
			cnhndlr.con = nil
		}
		cnhndlr = nil
	}
	return
}

func (cnhndlr *ConnHandler) LocalAddr() (addr net.Addr) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			addr = cnhndlr.con.LocalAddr()
		}
	}
	return
}

func (cnhndlr *ConnHandler) RemoteAddr() (addr net.Addr) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			addr = cnhndlr.con.RemoteAddr()
		}
	}
	return
}

func (cnhndlr *ConnHandler) SetDeadline(t time.Time) (err error) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			err = cnhndlr.con.SetDeadline(t)
		}
	}
	return
}

func (cnhndlr *ConnHandler) SetReadDeadline(t time.Time) (err error) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			err = cnhndlr.con.SetReadDeadline(t)
		}
	}
	return
}

func (cnhndlr *ConnHandler) SetWriteDeadline(t time.Time) (err error) {
	if cnhndlr != nil {
		if cnhndlr.con != nil {
			err = cnhndlr.con.SetWriteDeadline(t)
		}
	}
	return
}

// Accept waits for and returns the next connection to the listener.
func (lstnhndlr *ListnerHandler) Accept() (con net.Conn, err error) {
	/*if lstnhndlr.backlog > 0 {
		if !lstnhndlr.backLogStarted {
			func() {
				lstnhndlr.lck.Lock()
				defer lstnhndlr.lck.Unlock()
				if !lstnhndlr.backLogStarted {
					lstnhndlr.backLogStarted = true
					go func() {
						for {
							if con, err = lstnhndlr.ln.Accept(); err == nil {
								lstnhndlr.backcon <- newConnHandler(con)
							} else {
								time.Sleep(10 * time.Millisecond)
							}
						}
					}()
				}
			}()
		}

		doneChecking := false
		for !doneChecking {
			select {
			case con = <-lstnhndlr.backcon:
				doneChecking = true
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}

	} else {
		if con, err = lstnhndlr.ln.Accept(); err == nil {
			con = newConnHandler(con)
		}
	}*/
	if con, err = lstnhndlr.ln.Accept(); err == nil {
		if tcpcn, tcpcnok := con.(*net.TCPConn); tcpcnok {
			tcpcn.SetLinger(-1)
			tcpcn.SetReadBuffer(8192)
			tcpcn.SetWriteBuffer(8192)
			//tcpcn.SetNoDelay(false)
			tcpcn.SetKeepAlive(true)
			con = tcpcn
		}
		//con = newConnHandler(con)
	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (lstnhndlr *ListnerHandler) Close() (err error) {
	err = lstnhndlr.ln.Close()
	return
}

// Addr returns the listener's network address.
func (lstnhndlr *ListnerHandler) Addr() (addr net.Addr) {
	addr = lstnhndlr.ln.Addr()
	return
}

func (lstnrsrvr *lstnrserver) startListening(lstnr *Listener, backlog ...int) {
	if ln, err := net.Listen("tcp", lstnrsrvr.srvr.Addr); err == nil {
		go func() {
			if err := lstnrsrvr.srvr.Serve(&ListnerHandler{ln: ln, lck: &sync.Mutex{}}); err != nil && err != http.ErrServerClosed {
				fmt.Printf("error: Failed to serve HTTP: %v", err.Error())
			}
		}()
	} else {
		fmt.Println("error:", err.Error())
	}
}

func (lstnrsrvr *lstnrserver) stopListening(lstnr *Listener) {
	ctx := context.Background()
	if err := lstnrsrvr.srvr.Shutdown(ctx); err != nil {
		fmt.Println("Error closing server at ", lstnrsrvr.srvr.Addr, " ", err.Error())
	} else {
		fmt.Println("Closed server at ", lstnrsrvr.srvr.Addr)
	}
}

type contextKey struct {
	key string
}

var ConnContextKey = &contextKey{"http-conn"}

func newlstnrserver(hndlr http.Handler, addr string, unencrypted bool) (lstnrsrvr *lstnrserver) {
	var h2s = &http2.Server{}
	var srvr = &http.Server{Addr: addr, ConnState: func(cn net.Conn, cnstate http.ConnState) {
		switch cnstate {
		case http.StateNew:
		case http.StateIdle:
			cn.SetReadDeadline(time.Now())
		case http.StateClosed, http.StateHijacked:
			cn.SetReadDeadline(time.Now())
		}
	}, ConnContext: func(ctx context.Context, c net.Conn) context.Context {
		return context.WithValue(ctx, ConnContextKey, c)
	}, Handler: h2c.NewHandler(hndlr, h2s), ReadHeaderTimeout: time.Millisecond * 2000}

	lstnrsrvr = &lstnrserver{srvr: srvr, h2s: h2s, addr: addr}
	return
}

//Listener - struct
type Listener struct {
	hndlr        http.Handler
	lstnrservers map[string]*lstnrserver
	dne          chan bool
}

//NewListener - instance
func NewListener(hndlr http.Handler) (lstnr *Listener) {
	lstnr = &Listener{dne: make(chan bool, 1), hndlr: hndlr, lstnrservers: map[string]*lstnrserver{}}
	return
}

//ServeHTTP - refer http.Handler
func (lstnr *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if lstnr.hndlr == nil {
		if DefaultHttpHandler != nil {
			DefaultHttpHandler.ServeHTTP(w, r)
		}
	} else {
		lstnr.hndlr.ServeHTTP(w, r)
	}
}

//Shutdown - Listener gracefully
func (lstnr *Listener) Shutdown() {
	if lstnr != nil {

	}
}

//WaitOnShutdown - wait on internal done chan to close
func (lstnr *Listener) WaitOnShutdown() {
	if lstnr != nil {
		if lstnr.dne != nil {
			go func(dne chan bool) {
				defer close(dne)

				for {
					select {
					case d := <-dne:
						if d || !d {
							return
						}
					}
				}
			}(lstnr.dne)
			<-lstnr.dne
		}
	}
}

//Listen - on addr and indicate if ish2c
func (lstnr *Listener) Listen(addr string, ish2c ...bool) {
	if _, lstok := lstnr.lstnrservers[addr]; !lstok {
		var lstnrsrvr = newlstnrserver(lstnr, addr, len(ish2c) == 1 && ish2c[0])
		lstnr.lstnrservers[addr] = lstnrsrvr
		lstnrsrvr.startListening(lstnr, runtime.NumCPU()*15)
	}
}

//DefaultHttpHandler - DefaultHttpHandler
var DefaultHttpHandler http.Handler = nil

var glblstnr *Listener

//Listening - Global Listening
func Listening() (lstnr *Listener) {
	lstnr = glblstnr
	return
}

func init() {
	if glblstnr == nil {
		glblstnr = NewListener(DefaultHttpHandler)
	}
}
