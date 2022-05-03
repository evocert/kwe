package listen

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/evocert/kwe/requesting"
	"github.com/evocert/kwe/security"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type listenwrap struct {
	close func() error
	addr  func() net.Addr
	accpt func() (net.Conn, error)
}

func (lstnwrp *listenwrap) Addr() (addr net.Addr) {
	if lstnwrp != nil && lstnwrp.addr != nil {
		addr = lstnwrp.addr()
	}
	return
}

func (lstnwrp *listenwrap) Accept() (conn net.Conn, err error) {
	if lstnwrp != nil && lstnwrp.accpt != nil {
		conn, err = lstnwrp.accpt()
	}
	return
}

func (lstnwrp *listenwrap) Close() (err error) {
	if lstnwrp != nil && lstnwrp.close != nil {
		err = lstnwrp.close()
	}
	return
}

type listen struct {
	lstnr    *Listener
	ln       net.Listener
	tcpln    *net.TCPListener
	network  string
	addr     string
	lstnwrap *listenwrap
}

type connwrap struct {
	conn   net.Conn
	lstn   *listen
	lck    *sync.RWMutex
	strtdr bool
	strtdw bool
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (connwrp *connwrap) Read(b []byte) (n int, err error) {
	if connwrp != nil && connwrp.conn != nil {
		if connwrp != nil && connwrp.conn != nil {
			if bl := len(b); bl > 0 {
				rn := 0
				maxrl := 4096
				if bl < maxrl {
					maxrl = bl
				}
				for n < bl {
					if !connwrp.strtdr {
						connwrp.strtdr = true
						if err = connwrp.conn.SetReadDeadline(time.Now().Add(time.Duration(maxrl+5) * time.Millisecond)); err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
							err = nil
						}
					} else {
						if err = connwrp.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
							err = nil
						}
					}
					if err == nil {
						rn, err = connwrp.conn.Read(b[n : n+maxrl])
						n += rn
						if (bl - n) <= maxrl {
							maxrl = bl - n
						}
						if err != nil {
							if errors.Is(err, os.ErrDeadlineExceeded) {
								err = nil
							}
							break
						}
					}
				}
			}
		}
	}
	return
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (connwrp *connwrap) Write(b []byte) (n int, err error) {
	if connwrp != nil && connwrp.conn != nil {
		if bl := len(b); bl > 0 {
			wn := 0
			maxwl := 4096
			if bl < maxwl {
				maxwl = bl
			}
			for n < bl {
				if !connwrp.strtdw {
					connwrp.strtdw = true
				}
				if err == nil {
					wn, err = connwrp.conn.Write(b[n : n+maxwl])
					n += wn
					if (bl - n) <= maxwl {
						maxwl = bl - n
					}
					if err != nil {
						break
					}
				}
			}
		}
	}
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (connwrp *connwrap) Close() (err error) {
	if connwrp != nil {
		connwrp.lck.Lock()
		defer connwrp.lck.Unlock()
		if connwrp.conn != nil {
			err = connwrp.conn.Close()
		}
		if connwrp.lstn != nil {
			connwrp.lstn = nil
		}
		connwrp = nil
	}
	return
}

// LocalAddr returns the local network address.
func (connwrp *connwrap) LocalAddr() (addr net.Addr) {
	if connwrp != nil && connwrp.conn != nil {
		addr = connwrp.conn.LocalAddr()
	}
	return
}

// RemoteAddr returns the remote network address.
func (connwrp *connwrap) RemoteAddr() (addr net.Addr) {
	if connwrp != nil && connwrp.conn != nil {
		addr = connwrp.conn.RemoteAddr()
	}
	return
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (connwrp *connwrap) SetDeadline(t time.Time) (err error) {

	return
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (connwrp *connwrap) SetReadDeadline(t time.Time) (err error) {

	return
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (connwrp *connwrap) SetWriteDeadline(t time.Time) (err error) {

	return
}

func (lstn *listen) shutdown() {
	if lstn != nil {
		if lstn.lstnwrap != nil {
			lstn.lstnwrap.Close()
			lstn.ln = nil
		}
		if lstn.ln != nil {
			lstn.ln.Close()
			lstn.ln = nil
		}
		if lstn.lstnr != nil {
			delete(lstn.lstnr.lstnrs, lstn.addr)
		}
	}
}

func (lstn *listen) start() {
	go func() {
		var h2s = &http2.Server{}
		/*var srv = &http.Server{Handler: h2c.NewHandler(lstn.lstnr, h2s), ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, ConnContxtKey, c)
		}}
		srv.Serve(lstn.lstnwrap)*/
		internalServe(lstn.lstnwrap, h2c.NewHandler(lstn.lstnr, h2s))
	}()
}

type connkeyapi string

var ConnContxtKey connkeyapi = "http-con"

type ListenerAPI interface {
	Listen(string, ...string) error
	Shutdown(...string) error
	UnCertifyAddr(...string)
	CertifyAddr(string, string, ...string) error
	CasAddr(int64, int64, ...string) error
}

type Listener struct {
	lstnrs       map[string]*listen
	tlscertsconf map[string]*tls.Config
	ServeRequest func(requesting.RequestAPI) error
}

func NewListener(srvrqst ...func(ra requesting.RequestAPI) error) (lstnr *Listener) {
	lstnr = &Listener{lstnrs: map[string]*listen{}, tlscertsconf: map[string]*tls.Config{}}
	if len(srvrqst) == 1 && srvrqst[0] != nil {
		lstnr.ServeRequest = srvrqst[0]
	}
	return
}

func (lstnr *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	func() {
		if rqst := requesting.NewRequest(nil, r.URL.Path, r, w); rqst != nil {
			defer rqst.Close()
			if lstnr.ServeRequest != nil {
				lstnr.ServeRequest(rqst)
			}
		}
	}()
}

func (lstnr *Listener) CasAddr(caserial int64, certserial int64, addr ...string) (err error) {
	if lstnr != nil && caserial > 0 && certserial > 0 {
		if addrl := len(addr); addrl > 0 {
			if ca := security.GLOBALCAS().CA(caserial); ca != nil {
				if cert := ca.Certificate(certserial); cert != nil {
					for _, adr := range addr {
						if cnfg := lstnr.tlscertsconf[adr]; cnfg != nil {
							cnfg = nil
						}
						lstnr.tlscertsconf[adr] = cert.ServerTLSConf()
					}
				}
			}
		}
	}
	return
}

func (lstnr *Listener) CertifyAddr(servercert string, serverkey string, addr ...string) (err error) {
	if lstnr != nil && serverkey != "" && servercert != "" {
		if caserial, cacerialerr := strconv.ParseInt(servercert, 10, 64); cacerialerr == nil {
			if certserial, certserialerr := strconv.ParseInt(serverkey, 10, 64); certserialerr == nil {
				err = lstnr.CasAddr(caserial, certserial, addr...)
			}
		} else {
			cer, cererr := tls.X509KeyPair([]byte(servercert), []byte(serverkey))
			if cererr != nil {
				err = cererr
			} else if addrl := len(addr); addrl > 0 {
				for _, adr := range addr {
					if cnfg := lstnr.tlscertsconf[adr]; cnfg != nil {
						cnfg = nil
					}
					lstnr.tlscertsconf[adr] = &tls.Config{Certificates: []tls.Certificate{cer}}
				}
			}
		}
	}
	return
}

func (lstnr *Listener) UnCertifyAddr(addr ...string) {
	if lstnr != nil {
		if addrl := len(addr); addrl > 0 {
			for _, adr := range addr {
				if cnfg := lstnr.tlscertsconf[adr]; cnfg != nil {
					lstnr.tlscertsconf[adr] = nil
				}
			}
		}
	}
}

func (lstnr *Listener) accepts(addr string, cn net.Conn) (preppedcn net.Conn, err error) {
	if lstnr != nil {
		if adrcn := lstnr.tlscertsconf[addr]; adrcn != nil {
			preppedcn = tls.Server(cn, adrcn)
		} else {
			preppedcn = cn
		}
	}
	return
}

func (lstnr *Listener) Listen(network string, addr ...string) (err error) {
	if lstnr != nil {
		if len(addr) > 0 {
			for _, adr := range addr {
				if lstn := lstnr.lstnrs[adr]; lstn == nil {
					if lstn, err = newlisten(lstnr, network, adr); lstn != nil {
						lstnr.lstnrs[adr] = lstn
						lstn.start()
					}
				}
			}
		}
	}
	return
}

func (lstnr *Listener) Shutdown(addr ...string) (err error) {
	if lstnr != nil {
		if len(addr) > 0 {
			var addrsfound []string = nil
			for _, adr := range addr {
				if lstn := lstnr.lstnrs[adr]; lstn != nil {
					addrsfound = append(addrsfound, adr)
				}
			}
			if len(addrsfound) > 0 {
				for addri := range addrsfound {
					lstnr.lstnrs[addrsfound[addri]].shutdown()
					delete(lstnr.lstnrs, addrsfound[addri])
				}
			}
		}
	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (lstn *listen) Close() (err error) {
	if lstn != nil {
		if lstn.ln != nil {
			err = lstn.ln.Close()
		}
	}
	return
}

// Addr returns the listener's network address.
func (lstn *listen) Addr() (addr net.Addr) {
	if lstn != nil {
		if lstn.ln != nil {
			addr = lstn.ln.Addr()
		}
	}
	return
}

func newlisten(lstnr *Listener, network string, addr string) (lstn *listen, err error) {
	if ln, lnerr := net.Listen(network, addr); lnerr == nil && ln != nil {
		lstn = &listen{lstnr: lstnr, network: network, addr: addr, ln: ln}
		if tcpln, _ := ln.(*net.TCPListener); tcpln != nil {
			lstn.tcpln = tcpln
		}
		lstn.lstnwrap = &listenwrap{
			close: func() error {
				return lstn.ln.Close()
			},
			addr: func() net.Addr {
				return lstn.ln.Addr()
			},
			accpt: func() (cn net.Conn, err error) {
				if lstn.tcpln != nil {
					if tc, tcerr := lstn.tcpln.AcceptTCP(); tcerr == nil {
						tc.SetNoDelay(true)

						tc.SetReadBuffer(64 * 1024)
						tc.SetWriteBuffer(64 * 1024)
						if tcerr = tc.SetKeepAlive(true); tcerr == nil {
							if tcerr = tc.SetKeepAlivePeriod(time.Second * 60); tcerr == nil {
								if tcerr = tc.SetLinger(-1); tcerr == nil {
									cn = tc
									cn, err = lstnr.accepts(lstn.addr, cn)
								} else {
									err = tcerr
								}
							} else {
								err = tcerr
							}
						} else {
							err = tcerr
							tc.Close()
						}
					}
				} else if lstn.ln != nil {
					if cn, err = lstn.ln.Accept(); err == nil {
						cn, err = lstnr.accepts(lstn.addr, cn)
					}
				}
				return
			}}
	} else if lnerr != nil {
		err = lnerr
	}
	return
}
