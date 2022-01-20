package listen

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"

	//"net/http/pprof"
	"strings"

	"github.com/evocert/kwe/requesting"
)

type responseWriter struct {
	chunkedWriter io.WriteCloser
	wroteHeader   bool
	bufw          *bufio.Writer
	statusCode    int
	header        http.Header
	orgwtr        io.Writer
	req           *http.Request
}

func newResponseWriter(req *http.Request, conn net.Conn) *responseWriter {
	var orgwtr io.Writer = conn
	return &responseWriter{
		header: http.Header{}, bufw: bufio.NewWriter(orgwtr), orgwtr: orgwtr, statusCode: 200, req: req}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Flush() {
	if w != nil && w.bufw != nil {
		w.bufw.Flush()
	}
}

func (w *responseWriter) Hijack() (con net.Conn, bufrw *bufio.ReadWriter, err error) {
	if w != nil {
		if con, _ = w.orgwtr.(net.Conn); con != nil {
			bufrw = bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con))
		} else {
			err = fmt.Errorf("unable to hijack connection")
		}
	} else {
		err = fmt.Errorf("unable to hijack connection")
	}
	return con, bufrw, err
}

func (w *responseWriter) Close() (err error) {
	if w != nil {
		w.Flush()
		if !w.wroteHeader {
			w.writeHeader()
		}
		if w.bufw != nil {
			w.bufw = nil
		}
		if w.req != nil {
			w.req = nil
		}
		if w.chunkedWriter != nil {
			w.chunkedWriter.Close()
			w.chunkedWriter = nil
		}
		if w.orgwtr != nil {
			if clswtr := w.orgwtr.(io.Closer); clswtr != nil {
				clswtr.Close()
			}
			w.orgwtr = nil
		}
	}
	return err
}

func (w *responseWriter) writeHeader() {
	if !w.wroteHeader {
		w.wroteHeader = true
		if w.bufw != nil {
			if w.req != nil {
				protoHeaderLine := fmt.Sprintf("%s %d %s\r\n", w.req.Proto, w.statusCode, http.StatusText(w.statusCode))
				fmt.Fprint(w.bufw, protoHeaderLine)
				ischunked := false
				if len(w.header) > 0 {
					for hdr, hdv := range w.header {
						fmt.Fprintln(w.bufw, hdr+": "+strings.Join(hdv, ";"))
					}
				}
				fmt.Fprintln(w.bufw)
				w.Flush()
				if ischunked {
					w.chunkedWriter = httputil.NewChunkedWriter(w.orgwtr)
					w.bufw.Reset(w.chunkedWriter)
				}
			}
		}
	}
}

func (w *responseWriter) Write(b []byte) (n int, err error) {
	if bl := len(b); bl > 0 {
		if w != nil && w.bufw != nil {
			w.writeHeader()
			n, err = w.bufw.Write(b[:bl])
		}
	}
	return n, err
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w != nil {
		if !w.wroteHeader {
			w.statusCode = statusCode
		}
	}
}

/*
http.HandleFunc("/debug/pprof/", Index)
	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	http.HandleFunc("/debug/pprof/profile", Profile)
	http.HandleFunc("/debug/pprof/symbol", Symbol)
	http.HandleFunc("/debug/pprof/trace", Trace)
*/

func internalServe(ln net.Listener, httpHnflr http.Handler) {
	if ln != nil {
		go func() {
			for {
				var conn, connerr = ln.Accept()
				if connerr != nil {
					break
				}

				if conn != nil {
					go func() {
						defer conn.Close()
						if req, reqerr := http.ReadRequest(bufio.NewReader(conn)); reqerr == nil {
							if req != nil {
								if httpHnflr != nil {
									if w := newResponseWriter(req, conn); w != nil {
										func() {
											defer w.Close()
											ctx := context.WithValue(req.Context(), requesting.ConnContextKey, conn)
											req = req.WithContext(ctx)
											/*if req.URL.Path == "/debug/pprof/" {
												pprof.Index(w, req)
											} else if strings.HasPrefix(req.URL.Path, "/debug/pprof/cmdline") {
												pprof.Cmdline(w, req)
											} else if strings.HasPrefix(req.URL.Path, "/debug/pprof/profile") {
												pprof.Profile(w, req)
											} else if strings.HasPrefix(req.URL.Path, "/debug/pprof/symbol") {
												pprof.Symbol(w, req)
											} else if strings.HasPrefix(req.URL.Path, "/debug/pprof/trace") {
												pprof.Trace(w, req)
											} else {
												httpHnflr.ServeHTTP(w, req)
											}*/
											httpHnflr.ServeHTTP(w, req)
										}()
									}
								}
							}
						}
					}()
				}
			}
		}()
	}
}
