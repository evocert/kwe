package http

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/parameters"
)

type remoting interface {
	LocalAddr() string
	RemoteAddr() string
}

type Request struct {
	ctx         context.Context
	httpr       *http.Request
	rdr         io.Reader
	rqstr       iorw.Reader
	path        string
	startrdng   bool
	lcladdr     string
	rmtaddr     string
	prtclmethod string
	prtcl       string
}

var ConnContextKey = "http-conn"

func NewRequest(path string, r interface{}) (rqst *Request) {
	var ctx context.Context = nil
	var httpr *http.Request = nil
	var rdr io.Reader = nil
	var rmtngaddr remoting = nil
	if r != nil {
		if httpr, _ = r.(*http.Request); httpr == nil {
			rdr, _ = r.(io.Reader)
		}
	}
	if httpr != nil {
		path = httpr.URL.Path
		if rdr == nil {
			rdr = httpr.Body
		}
		ctx = httpr.Context()
	} else {
		ctx = context.Background()
	}
	if rdr != nil {
		rmtngaddr, _ = rdr.(remoting)
	}
	rqst = &Request{path: path, rdr: rdr, httpr: httpr, ctx: ctx, rqstr: iorw.NewEOFCloseSeekReader(rdr, false)}

	if ctxv := ctx.Value(ConnContextKey); ctxv != nil {
		if cnctn, _ := ctxv.(net.Conn); cnctn != nil {
			rqst.rmtaddr = cnctn.RemoteAddr().String()
			rqst.lcladdr = cnctn.LocalAddr().String()
		} else if rmtngaddr != nil {
			rqst.rmtaddr = rmtngaddr.RemoteAddr()
			rqst.lcladdr = rmtngaddr.LocalAddr()
		}
	} else if rmtngaddr != nil {
		rqst.rmtaddr = rmtngaddr.RemoteAddr()
		rqst.lcladdr = rmtngaddr.LocalAddr()
	}
	if httpr != nil {
		rqst.prtcl = httpr.Proto
		rqst.prtclmethod = httpr.Method
	}
	return
}

func (rqst *Request) IsValid() (bool, error) {
	if rqst != nil && rqst.ctx != nil {
		select {
		case <-rqst.ctx.Done():
			return rqst.ctx.Err() == nil, rqst.ctx.Err()
		default:
		}
	}
	return true, nil
}

func (rqst *Request) Proto() (proto string) {
	if rqst != nil {
		proto = rqst.prtcl
	}
	return
}

func (rqst *Request) Method() (method string) {
	if rqst != nil {
		method = rqst.prtclmethod
	}
	return
}

func (rqst *Request) Path() (path string) {
	if rqst != nil {
		path = rqst.path
	}
	return
}

func (rqst *Request) LoadParameters(prms *parameters.Parameters) {
	if rqst != nil && rqst.httpr != nil && prms != nil {
		parameters.LoadParametersFromHTTPRequest(prms, rqst.httpr)
	}
}

func (rqst *Request) Headers() (headers []string) {
	if rqst != nil && rqst.httpr != nil {
		headers = []string{}
		if hdr := rqst.httpr.Header; hdr != nil {
			for k := range hdr {
				headers = append(headers, k)
			}
		}
	}
	return
}

func (rqst *Request) Header(header string) (value string) {
	if rqst != nil && rqst.httpr != nil {
		value = rqst.httpr.Header.Get(header)
	}
	return
}

func (rqst *Request) StartedReading() (err error) {
	if rqst != nil {
		if !rqst.startrdng {

			rqst.startrdng = true
		}
	}
	return
}

func (rqst *Request) Read(p []byte) (n int, err error) {
	if rqst != nil {
		if pl := len(p); pl > 0 {
			if _, err = rqst.IsValid(); err == nil {
				if !rqst.startrdng {
					if err = rqst.StartedReading(); err != nil {
						return
					}
				}
				n, err = rqst.internRead(p)
			}
		}
	}
	return
}

func (rqst *Request) internRead(p []byte) (n int, err error) {
	if rqst != nil && rqst.rqstr != nil {
		n, err = rqst.rqstr.Read(p)
	}
	return
}

func (rqst *Request) SetMaxRead(max int64) (err error) {
	if rqst != nil && rqst.rqstr != nil {
		err = rqst.rqstr.SetMaxRead(max)
	}
	return
}

func (rqst *Request) Seek(offset int64, whence int) (n int64, err error) {
	if rqst != nil && rqst.rqstr != nil {
		n, err = rqst.rqstr.Seek(offset, whence)
	}
	return
}

func (rqst *Request) ReadRune() (r rune, size int, err error) {
	if rqst != nil && rqst.rqstr != nil {
		r, size, err = rqst.rqstr.ReadRune()
	}
	return
}

func (rqst *Request) RemoteAddr() (addr string) {
	if rqst != nil {
		addr = rqst.rmtaddr
	}
	return
}

func (rqst *Request) LocalAddr() (addr string) {
	if rqst != nil {
		addr = rqst.lcladdr
	}
	return
}

func (rqst *Request) Readln() (ln string, err error) {
	if rqst != nil {
		ln, err = iorw.ReadLine(rqst)
	}
	return
}

func (rqst *Request) Readlines() (lines []string, err error) {
	if rqst != nil {
		lines, err = iorw.ReadLines(rqst)
	}
	return
}

func (rqst *Request) ReadAll() (s string, err error) {
	if rqst != nil {
		s, err = iorw.ReaderToString(rqst)
	}
	return
}

func (rqst *Request) Close() (err error) {
	if rqst != nil {
		if rqst.httpr != nil {
			rqst.httpr = nil
		}
		if rqst.rdr != nil {
			if clsr, _ := rqst.rdr.(io.Closer); clsr != nil {
				clsr.Close()
			}
			rqst.rdr = nil
		}
		if rqst.rqstr != nil {
			rqst.rqstr = nil
		}
		rqst = nil
	}
	return
}
