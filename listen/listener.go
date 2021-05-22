package listen

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
	ln net.Listener
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
	n, err = cnhdnlr.con.Read(b)
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

func (cnhdnlr *ConnHandler) Readlines() (lines []string, err error) {
	lines, err = iorw.ReadLines(cnhdnlr)
	return
}

func (cnhdnlr *ConnHandler) ReadAll() (s string, err error) {
	s, err = iorw.ReaderToString(cnhdnlr)
	return
}

func (cnhdnlr *ConnHandler) Print(a ...interface{}) {
	iorw.Fprint(cnhdnlr)
}

func (cnhdnlr *ConnHandler) Println(a ...interface{}) {
	iorw.Fprintln(cnhdnlr)
}

func (cnhndlr *ConnHandler) Close() (err error) {
	err = cnhndlr.con.Close()
	return
}

func (cnhdnlr *ConnHandler) LocalAddr() (addr net.Addr) {
	addr = cnhdnlr.con.LocalAddr()
	return
}

func (cnhdnlr *ConnHandler) RemoteAddr() (addr net.Addr) {
	addr = cnhdnlr.con.RemoteAddr()
	return
}

func (cnhdnlr *ConnHandler) SetDeadline(t time.Time) (err error) {
	err = cnhdnlr.con.SetDeadline(t)
	return
}

func (cnhdnlr *ConnHandler) SetReadDeadline(t time.Time) (err error) {
	err = cnhdnlr.con.SetReadDeadline(t)
	return
}

func (cnhdnlr *ConnHandler) SetWriteDeadline(t time.Time) (err error) {
	err = cnhdnlr.con.SetWriteDeadline(t)
	return
}

// Accept waits for and returns the next connection to the listener.
func (lstnhndlr *ListnerHandler) Accept() (con net.Conn, err error) {
	if con, err = lstnhndlr.ln.Accept(); err == nil {
		con = newConnHandler(con)
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

func (lstnrsrvr *lstnrserver) startListening(lstnr *Listener) {
	//go func() {
	if ln, err := net.Listen("tcp", lstnrsrvr.srvr.Addr); err == nil {
		nxtln := &ListnerHandler{ln: ln}
		go func() {
			if err := lstnrsrvr.srvr.Serve(nxtln); err != nil && err != http.ErrServerClosed {
				fmt.Println("error: Failed to serve HTTP: %v", err.Error())
			}
		}()
	} else {
		fmt.Println("error:", err.Error())
	}
	/*if err := lstnrsrvr.srvr.ListenAndServe(); err != nil {
		lstnr.lstnrservers[lstnrsrvr.addr] = nil
		delete(lstnr.lstnrservers, lstnrsrvr.addr)
	}*/
	//}()
}

func (lstnrsrvr *lstnrserver) stopListening(lstnr *Listener) {
	ctx := context.Background()
	if err := lstnrsrvr.srvr.Shutdown(ctx); err != nil {
		fmt.Println("Error closing server at ", lstnrsrvr.srvr.Addr, " ", err.Error())
	}
	fmt.Println("Closed server at :%d", lstnrsrvr.srvr.Addr)
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
		lstnrsrvr.startListening(lstnr)
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
