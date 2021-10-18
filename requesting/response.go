package requesting

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/ws"
)

type flusherapi interface {
	Flush()
}

type Response struct {
	headers      map[string]string
	wtr          io.Writer
	flshr        flusherapi
	startw       bool
	status       int
	startWriting func() error
}

func NewResponse(wtr io.Writer, a ...interface{}) (rspnsapi ResponseAPI) {
	var httpw http.ResponseWriter = nil
	var rqstapi RequestAPI = nil
	var flshr flusherapi = nil
	if len(a) > 0 {
		for _, d := range a {
			if d != nil {
				if httpwd, _ := d.(http.ResponseWriter); httpwd != nil {
					if httpw == nil {
						httpw = httpwd
					}
				} else if rqstapid, _ := d.(RequestAPI); rqstapid != nil {
					if rqstapi == nil {
						rqstapi = rqstapid
					}
				}
			}
		}
	}
	if httpw == nil && wtr != nil {
		httpw, _ = wtr.(http.ResponseWriter)
	}
	if httpw != nil {
		if wtr == httpw {
			if rspnsapi != nil {
				if rqst, _ := rqstapi.(*Request); rqst != nil {
					if rqst.rdr != nil {
						if ws, _ := rqst.rdr.(*ws.ReaderWriter); ws != nil {
							wtr = ws
						}
					}
				}
			}
		} else {
			if rqstapi != nil {
				if rqst, _ := rqstapi.(*Request); rqst != nil {
					if rqst.rdr != nil {
						if ws, _ := rqst.rdr.(*ws.ReaderWriter); ws != nil {
							wtr = ws
						}
					}
				}
			}
		}
		if wtr == nil {
			wtr = httpw
		}
	}
	if wtr != nil {
		flshr, _ = wtr.(flusherapi)
		var rspns = &Response{wtr: wtr, flshr: flshr, headers: make(map[string]string), status: 200}
		if httpw != nil {
			rspns.startWriting = func() (err error) {
				if !rspns.startw {
					rspns.SetHeader("Content-Length", "0")
				}
				if rspns.Header("Content-Length") == "" {
					rspns.SetHeader("Connection", "close")
				}
				//MULTIMEDIA support for HTTP 1.1
				if rspns.startw {
					if strings.Contains(rspns.Header("Content-Type"), "video/") || strings.Contains(rspns.Header("Content-Type"), "audio/") {
						if rspns.Header("Content-Range") == "" {
							rspns.SetHeader("Accept-Ranges", "bytes")
						} else {
							rspns.SetStatus(206)
							rspns.SetHeader("Connection", "keep-alive")
						}
					}
				}
				if len(rspns.headers) > 0 {
					for hdr, hdv := range rspns.headers {
						httpw.Header().Set(hdr, hdv)
					}
				}
				httpw.WriteHeader(rspns.status)
				time.Sleep(5)
				if flshr != nil {
					flshr.Flush()
				}
				return
			}
		}
		rspnsapi = rspns
	}
	return
}

func (rspns *Response) IsValid() (valid bool, err error) {
	if rspns != nil {
		valid, err = true, nil
	}
	return
}

func (rspns *Response) Headers() (headers []string) {
	if rspns != nil && len(rspns.headers) > 0 {
		headers = make([]string, len(rspns.headers))
		headersi := 0
		for header := range rspns.headers {
			headers[headersi] = header
			headersi++
		}
	}
	return
}

func (rspns *Response) Header(header string) (value string) {
	if header != "" && len(rspns.headers) > 0 {
		value = rspns.headers[header]
	}
	return
}

func (rspns *Response) SetHeader(header string, value string) {
	if rspns != nil && header != "" && value != "" {
		rspns.headers[header] = value
	}
}

func (rspns *Response) SetStatus(status int) {
	if rspns != nil {
		rspns.status = status
	}
}

func (rspns *Response) Print(a ...interface{}) {
	if rspns != nil && rspns.wtr != nil {
		iorw.Fprint(rspns, a...)
	}
}

func (rspns *Response) Println(a ...interface{}) {
	if rspns != nil && rspns.wtr != nil {
		iorw.Fprintln(rspns, a...)
	}
}

func (rspns *Response) Write(p []byte) (n int, err error) {
	if rspns != nil && rspns.wtr != nil {
		if pl := len(p); pl > 0 {
			if !rspns.startw {
				rspns.startw = true
				if rspns.startWriting != nil {
					if err = rspns.startWriting(); err != nil {
						return
					}
					rspns.startWriting = nil
					rspns.Flush()
				}
			}
			n, err = rspns.wtr.Write(p)
		}
	}
	return
}

func (rspns *Response) Flush() {
	if rspns != nil && rspns.flshr != nil {
		rspns.flshr.Flush()
	}
}

func (rspns *Response) Close() (err error) {
	if rspns != nil {
		if rspns.startWriting != nil {
			if !rspns.startw {
				err = rspns.startWriting()
			}
			rspns.startWriting = nil
		}
		if rspns.headers != nil {
			rspns.headers = nil
		}
		if rspns.flshr != nil {
			rspns.flshr = nil
		}
		if rspns.wtr != nil {
			rspns.wtr = nil
		}
	}
	return
}
