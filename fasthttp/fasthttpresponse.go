package fasthttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/requesting"
)

type FastHttpResponse struct {
	fsthttprqstor *FasttHttpRequestor
	headers       map[string]string
	startwrtng    bool
	contenttype   string
	httpstatus    int
}

func iniFastHttpResponse(fsthttprqstor *FasttHttpRequestor) (fsthttprspns *FastHttpResponse) {
	fsthttprspns = &FastHttpResponse{fsthttprqstor: fsthttprqstor, headers: map[string]string{}, httpstatus: 200}

	return
}

func (fsthttprspns *FastHttpResponse) Request() (rqst requesting.RequestAPI) {
	if fsthttprspns != nil && fsthttprspns.fsthttprqstor != nil {
		rqst = fsthttprspns.fsthttprqstor.fsthttprqst
	}
	return
}

func (fsthttprspns *FastHttpResponse) IsValid() (valid bool, err error) {
	if fsthttprspns != nil && fsthttprspns.fsthttprqstor != nil {
		select {
		case <-fsthttprspns.fsthttprqstor.ctxvalid.Done():
			valid, err = false, fsthttprspns.fsthttprqstor.ctxvalid.Err()
		default:
			valid = true
		}
	}
	return
}

func (fsthttprspns *FastHttpResponse) Headers() (headers []string) {
	if hrdsl := len(fsthttprspns.headers); hrdsl > 0 {
		headers = make([]string, hrdsl)
		hdrsi := 0
		for hdr := range fsthttprspns.headers {
			headers[hdrsi] = hdr
			hdrsi++
		}
	}
	return
}

func (fsthttprspns *FastHttpResponse) Header(header string) (hdr string) {
	if header != "" && fsthttprspns != nil && len(fsthttprspns.headers) > 0 {
		hdr = fsthttprspns.headers[header]
	}
	return
}

func (fsthttprspns *FastHttpResponse) SetHeader(name string, value string) {
	if name != "" && fsthttprspns != nil {
		if name == "Content-Type" {
			fsthttprspns.contenttype = value
		}
		fsthttprspns.headers[name] = value
	}
}

func (fsthttprspns *FastHttpResponse) SetContentType(contentype string) {
	if contentype != "" && fsthttprspns != nil {
		fsthttprspns.contenttype = contentype
	}
}

func (fsthttprspns *FastHttpResponse) ContentType() (contenttype string) {
	if fsthttprspns != nil {
		contenttype = fsthttprspns.contenttype
	}
	return
}

func (fsthttprspns *FastHttpResponse) SetStatus(status int) {
	if fsthttprspns != nil {
		fsthttprspns.httpstatus = status
	}
}

func (fsthttprspns *FastHttpResponse) StartedWriting(wrapup ...bool) (err error) {
	if fsthttprspns.startwrtng {
		return
	}
	fsthttprspns.startwrtng = true
	defer func() {
		if rv := recover(); rv != nil {
			err = fmt.Errorf("%v", rv)
		}
	}()
	if fsthttprspns != nil {
		if len(wrapup) == 1 && wrapup[0] {
			fsthttprspns.SetHeader("Content-Length", "0")
		}
		if fsthttprspns.Header("Content-Type") == "" {
			fsthttprspns.SetHeader("Content-Type", fsthttprspns.contenttype)
		}
		if cntntl := fsthttprspns.Header("Content-Length"); cntntl != "" {
			if cntntl != "0" {
				delete(fsthttprspns.headers, "Content-Length")
			}
		}
		fsthttprspns.SetHeader("Cache-Control", "no-cache")
		fsthttprspns.SetHeader("Expires", time.Now().Format(http.TimeFormat))
		fsthttprspns.SetHeader("Connection", "close")

		if fsthttprspns.fsthttprqstor != nil {
			if fstctx := fsthttprspns.fsthttprqstor.fstctx; fstctx != nil {
				func() {
					defer func() {
						if rv := recover(); rv != nil {
							err = fmt.Errorf("%v", rv)
						}
					}()
					//if valid, _ := fsthttprspns.IsValid(); valid {
					for hdr, hdrv := range fsthttprspns.headers {
						fstctx.Response.Header.Add(hdr, hdrv)
					}
					//}
					fstctx.Response.Header.SetStatusCode(fsthttprspns.httpstatus)
				}()
			}
		}
	}
	return
}

func (fsthttprspns *FastHttpResponse) Print(a ...interface{}) {
	if fsthttprspns != nil {
		iorw.Fprint(fsthttprspns, a...)
	}
}

func (fsthttprspns *FastHttpResponse) Println(a ...interface{}) {
	if fsthttprspns != nil {
		iorw.Fprintln(fsthttprspns, a...)
	}
}

func (fsthttprspns *FastHttpResponse) Write(p []byte) (n int, err error) {
	if fsthttprspns != nil {
		if pl := len(p); pl > 0 {
			n, err = fsthttprspns.internWrite(p)
		}
	}
	return
}

func (fsthttprspns *FastHttpResponse) internWrite(p []byte) (n int, err error) {
	if fsthttprspns != nil && fsthttprspns.fsthttprqstor != nil {
		if pl := len(p); pl > 0 {
			if fstctx := fsthttprspns.fsthttprqstor.fstctx; fstctx != nil {
				if !fsthttprspns.startwrtng {
					if err = fsthttprspns.StartedWriting(); err != nil {
						return
					}
				}
				n, err = fstctx.Write(p)
			}
		}
	}
	return
}

func (fsthttprspns *FastHttpResponse) Flush() {
	if !fsthttprspns.startwrtng {
		if err := fsthttprspns.StartedWriting(); err != nil {
			return
		}
	}
}

func (fsthttprspns *FastHttpResponse) Close() (err error) {
	if fsthttprspns != nil {
		if fsthttprspns.headers != nil {
			fsthttprspns.headers = nil
		}
		if fsthttprspns.fsthttprqstor != nil {
			fsthttprspns.fsthttprqstor = nil
		}
		fsthttprspns = nil
	}
	return
}
