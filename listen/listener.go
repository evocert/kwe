package listen

import (
	"net/http"

	"github.com/evocert/kwe/chnls"
	http2 "golang.org/x/net/http2"
	h2c "golang.org/x/net/http2/h2c"
)

type lstnrserver struct {
	h2s  *http2.Server
	srvr *http.Server
	addr string
}

func (lstnrsrvr *lstnrserver) startListening(lstnr *Listener) {
	go func() {
		if err := lstnrsrvr.srvr.ListenAndServe(); err != nil {
			lstnr.lstnrservers[lstnrsrvr.addr] = nil
			delete(lstnr.lstnrservers, lstnrsrvr.addr)
		}
	}()
}

func newlstnrserver(hndlr http.Handler, addr string, unencrypted bool) (lstnrsrvr *lstnrserver) {
	var h2s = &http2.Server{}
	var srvr = &http.Server{Addr: addr, Handler: h2c.NewHandler(hndlr, h2s)}
	lstnrsrvr = &lstnrserver{srvr: srvr, h2s: h2s, addr: addr}
	return
}

//Listener - struct
type Listener struct {
	chnl         *chnls.Channel
	lstnrservers map[string]*lstnrserver
	dne          chan bool
}

//NewListener - instance
func NewListener(chnl *chnls.Channel) (lstnr *Listener) {
	if chnl == nil {
		chnl = chnls.NewChannel()
	}
	lstnr = &Listener{dne: make(chan bool, 1), chnl: chnl, lstnrservers: map[string]*lstnrserver{}}
	return
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
func (lstnr *Listener) Listen(addr string, ish2c bool) {
	if _, lstok := lstnr.lstnrservers[addr]; !lstok {
		var lstnrsrvr = newlstnrserver(lstnr.chnl, addr, ish2c)
		lstnr.lstnrservers[addr] = lstnrsrvr
		lstnrsrvr.startListening(lstnr)
	}
}

var glblstnr *Listener

//Listening - Global Listening
func Listening() (lstnr *Listener) {
	lstnr = glblstnr
	return
}

func init() {
	if glblstnr == nil {
		glblstnr = NewListener(nil)
	}
}
