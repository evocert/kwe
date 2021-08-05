package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/evocert/kwe/alertify"
	"github.com/evocert/kwe/chnls"
	rhttp "github.com/evocert/kwe/http"
)

type RawConn struct {
	rwlstnr        *RawListener
	valid          bool
	startedReading bool
	startedWriting bool
	doneRead       bool
	doneWrite      bool
	conn           net.Conn
	rbytes         []byte
	rbytesi        int
	rbytesl        int
	wbytes         []byte
	wbytesl        int
	wbytesi        int
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (rwcon *RawConn) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if rwcon != nil {
			for rwcon.valid && n < pl {
				if rwcon.rbytesl == 0 || (rwcon.rbytesl > 0 && rwcon.rbytesl == rwcon.rbytesi) {
					if rwcon.rbytesl > 0 {
						rwcon.rbytesl = 0
					}
					if rwcon.rbytesi > 0 {
						rwcon.rbytesi = 0
					}
					if rwcon.doneRead {
						break
					}
					cn := 0
					var cnerr error = nil
					if !rwcon.startedReading {
						rwcon.startedReading = true
						cn, cnerr = rwcon.conn.Read(rwcon.rbytes)
						if cnerr == nil {
							rwcon.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
						}
					} else {
						cn, cnerr = rwcon.conn.Read(rwcon.rbytes)
					}
					if cn > 0 {
						rwcon.rbytesl += cn
					}
					rwcon.doneRead = cnerr != nil
					if rwcon.doneRead {
						if cnerr != io.EOF {
							if !errors.Is(cnerr, os.ErrDeadlineExceeded) {
								rwcon.valid = cnerr != io.EOF
								return
							}
						}
					}
				}

				if rwcon.rbytesl > 0 {
					if tl := (rwcon.rbytesl - rwcon.rbytesi); (pl - n) >= tl {
						cl := copy(p[n:n+tl], rwcon.rbytes[rwcon.rbytesi:rwcon.rbytesi+tl])
						n += cl
						rwcon.rbytesi += cl
					} else if tl := (pl - n); (rwcon.rbytesl - rwcon.rbytesi) > tl {
						cl := copy(p[n:n+tl], rwcon.rbytes[rwcon.rbytesi:rwcon.rbytesi+tl])
						n += cl
						rwcon.rbytesi += cl
					}
				} else {
					break
				}
			}
		}
		if n == 0 {
			if err == nil {
				err = io.EOF
			}
		}
	}
	return
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (rwcon *RawConn) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if rwcon != nil {
			for rwcon.valid && n < pl {
				if rwcon.wbytesl = len(rwcon.wbytes); rwcon.wbytesl > 0 {
					if tl := (rwcon.wbytesl - rwcon.wbytesi); (pl - n) >= tl {
						cl := copy(rwcon.wbytes[rwcon.wbytesi:rwcon.wbytesi+tl], p[n:n+tl])
						n += cl
						rwcon.wbytesi += cl
					} else if tl := (pl - n); (rwcon.wbytesl - rwcon.wbytesi) > tl {
						cl := copy(rwcon.wbytes[rwcon.wbytesi:rwcon.wbytesi+tl], p[n:n+tl])
						n += cl
						rwcon.wbytesi += cl
					}
					if rwcon.wbytesl == rwcon.wbytesi {
						if flserr := rwcon.Flush(); flserr != nil {
							err = flserr
							break
						}
					}
				} else {
					break
				}
			}
		}
	}
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (rwcon *RawConn) Close() (err error) {
	if rwcon != nil {
		/*if rwcon.conn != nil {
			err = rwcon.conn.Close()
			rwcon.conn = nil
		}
		if rwcon.rwlstnr != nil {
			if rwcon.rwlstnr.rwchnlspool != nil {
				rwcon.rwlstnr.rwchnlspool.Put(rwcon)
			}
		}*/
		rwcon.Flush()
		rwcon.Dispose()
		rwcon = nil
	}
	return
}

func (rwcon *RawConn) Flush() (err error) {
	if rwcon != nil && rwcon.conn != nil && rwcon.wbytesi > 0 {
		cn := 0
		wi := 0
		wl := rwcon.wbytesi
		rwcon.wbytesi = 0
		var cnerr error = nil
		for cnerr == nil && wi < wl {
			cn, cnerr = rwcon.conn.Write(rwcon.wbytes[wi : wi+(wl-wi)])
			if cn > 0 {
				wi += cn
			}
			if !rwcon.startedWriting {
				rwcon.startedWriting = true
				if cnerr == nil {
					//rwcon.conn.SetWriteDeadline(time.Now().Add(200 * time.Millisecond))
				}
			}
		}
	}
	return
}

func (rwcon *RawConn) Dispose() (err error) {
	if rwcon != nil {
		if rwcon.conn != nil {

			err = rwcon.conn.Close()
			rwcon.conn = nil
		}
		if rwcon.rwlstnr != nil {
			rwcon.rwlstnr = nil
		}
		rwcon = nil
	}
	return
}

// LocalAddr returns the local network address.
func (rwcon *RawConn) LocalAddr() (lcladdr net.Addr) {
	if rwcon != nil && rwcon.conn != nil {
		lcladdr = rwcon.conn.LocalAddr()
	}
	return
}

// RemoteAddr returns the remote network address.
func (rwcon *RawConn) RemoteAddr() (rmtaddr net.Addr) {
	if rwcon != nil && rwcon.conn != nil {
		rmtaddr = rwcon.conn.RemoteAddr()
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
func (rwcon *RawConn) SetDeadline(t time.Time) (err error) {

	return
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (rwcon *RawConn) SetReadDeadline(t time.Time) (err error) {

	return
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (rwcon *RawConn) SetWriteDeadline(t time.Time) (err error) {

	return
}

func newRawConn(rwlstnr *RawListener, conn net.Conn, readsize int, writesize int) (rwconn *RawConn) {
	rwconn = &RawConn{valid: true, rwlstnr: rwlstnr, conn: conn, rbytes: make([]byte, readsize), wbytes: make([]byte, writesize), rbytesi: 0, wbytesi: 0}
	if rwconn.rwlstnr != rwlstnr {
		if rwconn.rwlstnr != nil {
			rwconn.rwlstnr = nil
		}
		rwconn.rwlstnr = rwlstnr
	}
	rwconn.conn = conn
	return
}

type RawListener struct {
	network   string
	addr      string
	started   bool
	lstnr     net.Listener
	lstnraddr net.Addr
	rwchnls   chan *RawConn
}

func NewRawListener(network string, addr string) (rwlstnr *RawListener) {
	rwlstnr = &RawListener{network: network, addr: addr, rwchnls: make(chan *RawConn)}
	return
}

func (rwlstnr *RawListener) startListening() (err error) {
	if !rwlstnr.started {
		rwlstnr.lstnr, err = net.Listen(rwlstnr.network, rwlstnr.addr)
		if rwlstnr.started = err == nil; rwlstnr.started {
			go rwlstnr.accepting(rwlstnr.lstnr)
		}
	}
	return
}

func (rwlstnr *RawListener) accepting(listener net.Listener) (err error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			break
		}
		go func() {
			rwlstnr.rwchnls <- newRawConn(rwlstnr, conn, 4096, 4096)
		}()
	}
	return
}

// Accept waits for and returns the next connection to the listener.
func (rwlstnr *RawListener) Accept() (conn net.Conn, err error) {
	if rwlstnr != nil {
		conn, err = <-rwlstnr.rwchnls, nil
	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (rwlstnr *RawListener) Close() (err error) {

	return
}

// Addr returns the listener's network address.
func (rwlstnr *RawListener) Addr() (addr net.Addr) {
	if rwlstnr != nil {
		addr = rwlstnr.lstnraddr
	}
	return
}

func main() {
	chnls := chnls.GLOBALCHNL()
	if rwlstnr := NewRawListener("tcp", ":11223"); rwlstnr != nil {
		if err := rwlstnr.startListening(); err == nil {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGQUIT)
			srv := &http.Server{Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					/*w.Header().Set("Content-Type", func() (s string) {
						s, _ = mimes.FindMimeType(r.URL.Path, "text/plain")
						return
					}())
					w.Header().Set("Connection", "close")
					iorw.Fprintln(w, alertify.AlertifyJS())*/
					chnls.ServeRequest(r, w, rhttp.RequestInvoker, rhttp.ResponseInvoker)
				})}
			srv.Serve(rwlstnr)
			<-sigs
		} else {
			fmt.Print(err.Error())
		}
	}
}
