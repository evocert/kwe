package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/requesting"
)

type Response struct {
	rqst        requesting.RequestAPI
	httpw       http.ResponseWriter
	httpstatus  int
	httpflshr   http.Flusher
	wtr         io.Writer
	startwrtng  bool
	contenttype string
	//wbuf        *iorw.Buffer
	//wbufr       *iorw.BuffReader
	wbytes []byte
	wbytei int
}

func NewResponse(w interface{}, a ...requesting.RequestAPI) (rspns *Response) {
	var rqst requesting.RequestAPI = nil
	if len(a) > 0 {
		for _, d := range a {
			if d != nil {
				if rqstrd := d.(requesting.RequestAPI); rqstrd != nil {
					if rqst == nil {
						rqst = rqstrd
					}
				}
			}
		}
	}
	var wtr io.Writer
	var httpw http.ResponseWriter

	if httpw, _ = w.(http.ResponseWriter); httpw != nil {
		wtr = httpw
	} else {
		wtr, _ = w.(io.Writer)
	}

	var httpflshr http.Flusher = nil
	if httpflshr == nil && wtr != nil {
		httpflshr, _ = wtr.(http.Flusher)
	}
	rspns = &Response{rqst: rqst, wtr: wtr, httpw: httpw, httpflshr: httpflshr, httpstatus: 200, wbytei: 0, wbytes: make([]byte, 1024*1024)}
	//rspns.wbuf.MaxLenToWrite = 1024 * 1024
	//rspns.wbuf.OnMaxWritten = func(maxwritten int64) bool {
	//	rspns.internFlush()
	//	return true
	//}
	//rspns.wbufr = rspns.wbuf.Reader()

	return
}

func (rspns *Response) Request() (rqst requesting.RequestAPI) {
	if rspns != nil {
		rqst = rspns.rqst
	}
	return
}

func (rspns *Response) IsValid() (valid bool, err error) {
	if rspns != nil {
		if rspns.rqst != nil {
			valid, err = rspns.rqst.IsValid()
		} else {
			valid = true
		}
	} else {
		valid = true
	}
	return valid, err
}

func (rspns *Response) Headers() (headers []string) {
	if rspns != nil && rspns.httpw != nil {
		headers = []string{}
		if hdr := rspns.httpw.Header(); hdr != nil {
			for k := range hdr {
				headers = append(headers, k)
			}
		}
	}
	return
}

func (rspns *Response) Header(header string) (value string) {
	if rspns != nil && rspns.httpw != nil {
		value = rspns.httpw.Header().Get(header)
	}
	return
}

func (rspns *Response) SetHeader(header string, value string) {
	if rspns != nil && rspns.httpw != nil {
		rspns.httpw.Header().Set(header, value)
	}
}

func (rspns *Response) SetStatus(status int) {
	if rspns != nil && rspns.httpw != nil {
		rspns.httpstatus = status
	}
}

func (rspns *Response) ContentType() (contenttype string) {
	if rspns != nil {
		contenttype = rspns.contenttype
	}
	return
}

func (rspns *Response) SetContentType(contenttype string) {
	if rspns != nil && rspns.httpw != nil {
		rspns.contenttype = contenttype
	}
}

func (rspns *Response) StartedWriting(wrapup ...bool) (err error) {
	if rspns.startwrtng {
		return
	}
	rspns.startwrtng = true
	defer func() {
		if rv := recover(); rv != nil {
			err = fmt.Errorf("%v", rv)
		}
	}()
	if rspns != nil && rspns.httpw != nil {
		if len(wrapup) == 1 && wrapup[0] && rspns.httpw != nil {
			rspns.httpw.Header().Set("Content-Length", "0")
		}
		if hdr := rspns.httpw.Header(); hdr != nil {
			if hdr.Get("Content-Type") == "" {
				hdr.Set("Content-Type", rspns.contenttype)
			}
			if cntntl := hdr.Get("Content-Length"); cntntl != "" {
				if cntntl != "0" {
					hdr.Del("Content-Length")
				}
			}
			hdr.Set("Cache-Control", "no-cache")
			hdr.Set("Expires", time.Now().Format(http.TimeFormat))
			hdr.Set("Connection", "close")
		}
		if rspns.httpw != nil {
			func() {
				defer func() {
					if rv := recover(); rv != nil {
						err = fmt.Errorf("%v", rv)
					}
				}()
				if valid, _ := rspns.IsValid(); valid {
					rspns.httpw.WriteHeader(rspns.httpstatus)
				}
			}()
		}
	}
	return
}

func (rspns *Response) Write(p []byte) (n int, err error) {
	if rspns != nil {
		if pl := len(p); pl > 0 {
			n, err = rspns.internWrite(p)
		}
	}
	return
}

func (rspns *Response) internFlush() (err error) {
	if rspns != nil {
		if rspns.wbytei > 0 {
			if !rspns.startwrtng {
				if err = rspns.StartedWriting(); err != nil {
					return
				}
			}
			p := rspns.wbytes[:rspns.wbytei]
			rspns.wbytei = 0
			if httpw := rspns.httpw; httpw != nil {
				if httpw != nil {
					_, err = httpw.Write(p)
				} else if rqstw := rspns.wtr; rqstw != nil {
					_, err = rqstw.Write(p)
				}
			} else if rqstw := rspns.wtr; rqstw != nil {
				_, err = rqstw.Write(p)
			}
		}
	}
	return
}

func (rspns *Response) internWrite(p []byte) (n int, err error) {
	if rspns != nil {
		if pl, wl := len(p), len(rspns.wbytes); pl > 0 {
			for n < pl {
				if tl := (wl - rspns.wbytei); (pl - n) >= tl {
					cpl := copy(rspns.wbytes[rspns.wbytei:rspns.wbytei+(tl)], p[n:n+tl])
					n += cpl
					rspns.wbytei += cpl
				} else if tl := (pl - n); tl < (wl - rspns.wbytei) {
					cpl := copy(rspns.wbytes[rspns.wbytei:rspns.wbytei+(tl)], p[n:n+tl])
					n += cpl
					rspns.wbytei += cpl
				}
				if rspns.wbytei == wl {
					rspns.internFlush()
				}
			}
		}
	}
	return
}

func (rspns *Response) Flush() {
	if rspns != nil {
		rspns.internFlush()
	}
	if rspns != nil && rspns.httpflshr != nil {
		rspns.httpflshr.Flush()
	}
}

func (rspns *Response) Print(a ...interface{}) {
	if rspns != nil {
		iorw.Fprint(rspns, a...)
	}
}

func (rspns *Response) Println(a ...interface{}) {
	if rspns != nil {
		iorw.Fprintln(rspns, a...)
	}
}

func (rspns *Response) Close() (err error) {
	if rspns != nil {
		if rspns.httpflshr != nil {
			rspns.httpflshr = nil
		}
		if rspns.httpw != nil {
			rspns.httpw = nil
		}
		//if rspns.wbufr != nil {
		//	rspns.wbufr.Close()
		//	rspns.wbufr = nil
		//}
		//if rspns.wbuf != nil {
		//	rspns.wbuf.Close()
		//	rspns.wbuf = nil
		//}
		if rspns.wtr != nil {
			if cls, _ := rspns.wtr.(io.Closer); cls != nil {
				cls.Close()
			}
			rspns.wtr = nil
		}
		if rspns.rqst != nil {
			rspns.rqst = nil
		}
		rspns = nil
	}
	return
}

var ResponseInvoker requesting.ResponseInvokerFunc = nil

func init() {
	ResponseInvoker = func(w interface{}, a ...requesting.RequestAPI) requesting.ResponseAPI {
		return NewResponse(w, a...)
	}
}
