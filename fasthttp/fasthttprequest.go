package fasthttp

import (
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/parameters"
	"github.com/valyala/fasthttp"
)

type FastHttpRequest struct {
	ftshttprqstor *FasttHttpRequestor
	qryargsdst    *fasthttp.Args
	postargsdst   *fasthttp.Args

	headers    map[string]string
	rqstr      iorw.Reader
	path       string
	prtcl      string
	prtclmthd  string
	strtdrding bool
	lcladdr    string
	rmtaddr    string
}

func iniFastHttpRequest(fsthttprqstor *FasttHttpRequestor) (fsthttprqst *FastHttpRequest) {
	fsthttprqst = &FastHttpRequest{ftshttprqstor: fsthttprqstor, headers: map[string]string{}}
	var rqstr iorw.Reader = nil
	if ctx := fsthttprqstor.fstctx; ctx != nil {
		fsthttprqst.lcladdr = ctx.Conn().LocalAddr().String()
		fsthttprqst.rmtaddr = ctx.Conn().RemoteAddr().String()
		fsthttprqst.path = string(ctx.Request.URI().Path())
		fsthttprqst.prtcl = string(ctx.Request.Header.Protocol())
		fsthttprqst.prtclmthd = string(ctx.Request.Header.Method())

		ctx.Request.Header.VisitAll(func(key, value []byte) {
			fsthttprqst.headers[string(key)] = string(value)
		})
		rqstr = iorw.NewEOFCloseSeekReader(ctx.RequestBodyStream(), false)
	} else {
		rqstr = iorw.NewEOFCloseSeekReader(nil)
	}
	fsthttprqst.rqstr = rqstr
	return
}

func (fsthttprqst *FastHttpRequest) IsValid() (valid bool, err error) {
	if fsthttprqst != nil && fsthttprqst.ftshttprqstor != nil {
		select {
		case <-fsthttprqst.ftshttprqstor.ctxvalid.Done():
			valid, err = false, fsthttprqst.ftshttprqstor.ctxvalid.Err()
		default:
			valid = true
		}
	}
	return
}
func (fsthttprqst *FastHttpRequest) Proto() (prcl string) {
	if fsthttprqst != nil {
		prcl = fsthttprqst.prtcl
	}
	return
}

func (fsthttprqst *FastHttpRequest) Method() (prtclmthd string) {
	if fsthttprqst != nil {
		prtclmthd = fsthttprqst.prtclmthd
	}
	return
}

func (fsthttprqst *FastHttpRequest) Path() (path string) {
	if fsthttprqst != nil {
		path = fsthttprqst.path
	}
	return
}

func (fsthttprqst *FastHttpRequest) LoadParameters(prms *parameters.Parameters) {
	if fsthttprqst != nil && fsthttprqst.ftshttprqstor != nil {
		if ctx := fsthttprqst.ftshttprqstor.fstctx; ctx != nil {
			fsthttprqst.qryargsdst = fasthttp.AcquireArgs()
			ctx.QueryArgs().CopyTo(fsthttprqst.qryargsdst)
			fsthttprqst.qryargsdst.VisitAll(func(key, value []byte) {
				prms.SetParameter(string(key), false, string(value))
			})
			fsthttprqst.postargsdst = fasthttp.AcquireArgs()
			ctx.PostArgs().CopyTo(fsthttprqst.postargsdst)
			fsthttprqst.postargsdst.VisitAll(func(key, value []byte) {
				prms.SetParameter(string(key), false, string(value))
			})
			if frm, frmerr := ctx.MultipartForm(); frmerr == nil {
				if frmvals := frm.Value; len(frmvals) > 0 {
					for frmk, frmv := range frmvals {
						for _, fv := range frmv {
							prms.SetParameter(frmk, false, fv)
						}
					}
				}
			}
		}
	}
}

func (fsthttprqst *FastHttpRequest) Headers() (headers []string) {
	if hrdsl := len(fsthttprqst.headers); hrdsl > 0 {
		headers = make([]string, hrdsl)
		hdrsi := 0
		for hdr := range fsthttprqst.headers {
			headers[hdrsi] = hdr
			hdrsi++
		}
	}
	return
}

func (fsthttprqst *FastHttpRequest) Header(header string) (hdr string) {
	if header != "" && fsthttprqst != nil && len(fsthttprqst.headers) > 0 {
		hdr = fsthttprqst.headers[header]
	}
	return
}

func (fsthttprqst *FastHttpRequest) RemoteAddr() (rmtaddr string) {
	if fsthttprqst != nil {
		rmtaddr = fsthttprqst.rmtaddr
	}
	return
}

func (fsthttprqst *FastHttpRequest) LocalAddr() (lcladdr string) {
	if fsthttprqst != nil {
		lcladdr = fsthttprqst.lcladdr
	}
	return
}

func (fsthttprqst *FastHttpRequest) StartedReading() (err error) {
	if fsthttprqst != nil {
		if !fsthttprqst.strtdrding {

			fsthttprqst.strtdrding = true
		}
	}
	return
}

func (fsthttprqst *FastHttpRequest) Readln() (ln string, err error) {
	if fsthttprqst != nil {
		ln, err = iorw.ReadLine(fsthttprqst)
	}
	return
}

func (fsthttprqst *FastHttpRequest) Readlines() (lines []string, err error) {
	if fsthttprqst != nil {
		lines, err = iorw.ReadLines(fsthttprqst)
	}
	return
}

func (fsthttprqst *FastHttpRequest) ReadAll() (s string, err error) {
	if fsthttprqst != nil {
		s, err = iorw.ReaderToString(fsthttprqst)
	}
	return
}

func (fsthttprqst *FastHttpRequest) Read(p []byte) (n int, err error) {
	if fsthttprqst != nil {
		if pl := len(p); pl > 0 {
			if _, err = fsthttprqst.IsValid(); err == nil {
				if !fsthttprqst.strtdrding {
					if err = fsthttprqst.StartedReading(); err != nil {
						return
					}
				}
				n, err = fsthttprqst.internRead(p)
			}
		}
	}
	return
}

func (fsthttprqst *FastHttpRequest) internRead(p []byte) (n int, err error) {
	if fsthttprqst != nil && fsthttprqst.rqstr != nil {
		n, err = fsthttprqst.rqstr.Read(p)
	}
	return
}

func (fsthttprqst *FastHttpRequest) SetMaxRead(max int64) (err error) {
	if fsthttprqst.rqstr != nil {
		err = fsthttprqst.rqstr.SetMaxRead(max)
	}
	return
}

func (fsthttprqst *FastHttpRequest) Seek(offset int64, whence int) (n int64, err error) {
	if fsthttprqst.rqstr != nil {
		n, err = fsthttprqst.rqstr.Seek(offset, whence)
	}
	return
}

func (fsthttprqst *FastHttpRequest) ReadRune() (r rune, size int, err error) {
	if fsthttprqst.rqstr != nil {
		r, size, err = fsthttprqst.rqstr.ReadRune()
	}
	return
}

func (fsthttprqst *FastHttpRequest) Close() (err error) {
	if fsthttprqst != nil {
		if fsthttprqst.ftshttprqstor != nil {
			fsthttprqst.ftshttprqstor = nil
		}
		if fsthttprqst.qryargsdst != nil {
			fasthttp.ReleaseArgs(fsthttprqst.qryargsdst)
			fsthttprqst.qryargsdst = nil
		}
		if fsthttprqst.postargsdst != nil {
			fasthttp.ReleaseArgs(fsthttprqst.postargsdst)
			fsthttprqst.postargsdst = nil
		}
		if fsthttprqst.headers != nil {
			fsthttprqst.headers = nil
		}
		fsthttprqst = nil
	}
	return
}
