package requesting

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/ws"
)

type contexting interface {
	Context() context.Context
}

type addressing interface {
	LocalAddr() string
	RemoteAddr() string
}

type Request struct {
	ctx         context.Context
	path        string
	prms        parameters.ParametersAPI
	prtcl       string
	prtclmthd   string
	headers     map[string]string
	rdr         io.Reader
	rqstr       iorw.Reader
	rmtaddr     string
	lcladdr     string
	rspns       ResponseAPI
	rangeType   string
	rangeOffset int64
}

var ConnContextKey = "http-con"

func NewRequest(rdr io.Reader, a ...interface{}) (rqstapi RequestAPI) {
	var rqst *Request = nil
	var path string = ""
	var prtcl string = ""
	var prtclmthd string = ""
	var prtclrangetype = ""
	var prtclrange = ""
	var prtclrangeoffset = int64(-1)
	var headers = map[string]string{}
	var ctx context.Context = nil
	var addrsng addressing = nil
	var prms parameters.ParametersAPI = nil
	var httpr *http.Request = nil
	var httpw http.ResponseWriter = nil
	var rmtaddr = ""
	var lcladdr = ""
	var wtr io.Writer = nil
	for _, d := range a {
		if d != nil {
			if pathd, _ := d.(string); pathd != "" {
				if path == "" {
					path = pathd
				}
			} else if rdrd, _ := d.(io.Reader); rdrd != nil {
				if rdr == nil {
					rdr = rdrd
				}
			} else if prmsd, _ := d.(parameters.ParametersAPI); prmsd != nil {
				if prms == nil {
					prms = prmsd
				}
			} else if httpwd, _ := d.(http.ResponseWriter); httpwd != nil {
				if httpw == nil {
					httpw = httpwd
					if wtr == nil {
						wtr = httpw
					}
				}
			} else if httprd, _ := d.(*http.Request); httprd != nil {
				if httpr == nil {
					httpr = httprd
				}
				if ctx == nil {
					ctx = httprd.Context()
					if ctxv := ctx.Value(ConnContextKey); ctxv != nil {
						if cnctn, _ := ctxv.(net.Conn); cnctn != nil {
							rmtaddr = cnctn.RemoteAddr().String()
							lcladdr = cnctn.LocalAddr().String()
						}
					}
				}
				if httprd.Body != nil {
					if rdr == nil {
						rdr = httprd.Body
					}
					prtcl = httprd.Proto
					prtclmthd = httprd.Method
				}
				for hdrk, hdrv := range httprd.Header {
					headers[hdrk] = strings.Join(hdrv, "")
					if hdrk == "Range" {
						if prtclrange = headers[hdrk]; strings.Index(prtclrange, "=") > 0 {
							if prtclrangetype = prtclrange[:strings.Index(prtclrange, "=")]; prtclrange != "" {
								if prtclrange = prtclrange[strings.Index(prtclrange, "=")+1:]; strings.Index(prtclrange, "-") > 0 {
									prtclrangeoffset, _ = strconv.ParseInt(prtclrange[:strings.Index(prtclrange, "-")], 10, 64)
								}
							}
						}
					}
				}
			} else {
				if cntxtng, _ := d.(contexting); cntxtng != nil && ctx == nil {
					ctx = cntxtng.Context()
				}
				if addrsngd, _ := d.(addressing); addrsngd != nil && addrsng == nil {
					addrsng = addrsngd
				}
			}
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if rdr != nil && addrsng == nil {
		addrsng, _ = rdr.(addressing)
	}
	if httpr != nil {
		if prms == nil {
			prms = parameters.NewParameters()
		}
		parameters.LoadParametersFromHTTPRequest(prms, httpr)
		if httpw != nil {
			if wsrw, wsrwerr := ws.NewServerReaderWriter(httpw, httpr); wsrw != nil && wsrwerr == nil {
				rdr = wsrw
				wtr = wsrw
			}
		}
	}

	if addrsng != nil {
		rmtaddr = addrsng.RemoteAddr()
		lcladdr = addrsng.LocalAddr()
	}

	rqst = &Request{rdr: rdr, ctx: ctx, rqstr: iorw.NewEOFCloseSeekReader(rdr, false), lcladdr: lcladdr, rmtaddr: rmtaddr, path: path, prtcl: prtcl, prtclmthd: prtclmthd, headers: headers, prms: prms, rangeType: prtclrangetype, rangeOffset: prtclrangeoffset}

	rqst.rspns = NewResponse(wtr, rqst, httpw)

	rqstapi = rqst
	return
}

func (rqst *Request) RangeType() (rangetype string) {
	if rqst != nil {
		rangetype = rqst.rangeType
	}
	return
}

func (rqst *Request) RangeOffset() (rangeoffset int64) {
	if rqst != nil {
		rangeoffset = rqst.rangeOffset
	}
	return
}

func (rqst *Request) Proto() (prtcl string) {
	if rqst != nil {
		prtcl = rqst.prtcl
	}
	return
}

func (rqst *Request) Method() (prtclmthd string) {
	if rqst != nil {
		prtclmthd = rqst.prtclmthd
	}
	return
}

func (rqst *Request) Path() (path string) {
	if rqst != nil {
		path = rqst.path
	}
	return
}

func (rqst *Request) Headers() (headers []string) {
	if rqst != nil && len(rqst.headers) > 0 {
		headers = make([]string, len(rqst.headers))
		headersi := 0
		for hdr := range rqst.headers {
			headers[headersi] = hdr
			headersi++
		}
	}
	return
}

func (rqst *Request) Header(header string) (value string) {
	if header != "" && rqst != nil && len(rqst.headers) > 0 {
		value = rqst.headers[header]
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

func (rqst *Request) Parameters() (prms parameters.ParametersAPI) {
	if rqst != nil {
		prms = rqst.prms
	}
	return
}

func (rqst *Request) Read(p []byte) (n int, err error) {
	if rqst != nil && rqst.rdr != nil {
		n, err = rqst.rdr.Read(p)
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

func (rqst *Request) ReadRune() (r rune, size int, err error) {
	if rqst != nil && rqst.rdr != nil && rqst.rqstr != nil {
		r, size, err = rqst.rqstr.ReadRune()
	}
	return
}
func (rqst *Request) Readln() (ln string, err error) {
	if rqst != nil && rqst.rdr != nil && rqst.rqstr != nil {
		ln, err = rqst.rqstr.Readln()
	}
	return
}
func (rqst *Request) ReadLines() (lines []string, err error) {
	if rqst != nil && rqst.rdr != nil && rqst.rqstr != nil {
		lines, err = rqst.rqstr.Readlines()
	}
	return
}
func (rqst *Request) ReadAll() (all string, err error) {
	if rqst != nil && rqst.rdr != nil && rqst.rqstr != nil {
		all, err = rqst.rqstr.ReadAll()
	}
	return
}

func (rqst *Request) IsValid() (bool, error) {
	if rqst != nil && rqst.ctx != nil {
		select {
		case <-rqst.ctx.Done():
			return rqst.ctx.Err() == nil, rqst.ctx.Err()
		default:
			return true, nil
		}
	}
	return false, nil
}

func (rqst *Request) Response() (rspns ResponseAPI) {
	if rqst != nil {
		rspns = rqst.rspns
	}
	return
}

func (rqst *Request) Close() (err error) {
	if rqst != nil {
		if rqst.headers != nil {
			rqst.headers = nil
		}
		if rqst.rdr != nil {
			if eofr, _ := rqst.rdr.(*iorw.EOFCloseSeekReader); eofr != nil {
				eofr.CanClose = true
				eofr.Close()
			} else if clsr, _ := rqst.rdr.(io.Closer); clsr != nil {
				clsr.Close()
			}
			rqst.rdr = nil
		}
		if rqst.rqstr != nil {
			rqst.rqstr = nil
		}
		if rqst.prms != nil {
			rqst.prms.CleanupParameters()
			rqst.prms = nil
		}
		if rqst.rspns != nil {
			rqst.rspns.Close()
			rqst.rspns = nil
		}
		rqst = nil
	}
	return
}
