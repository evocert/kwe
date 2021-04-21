package web

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"

	"github.com/evocert/kwe/ws"

	"github.com/gorilla/websocket"
)

type ClientHandle struct {
	*Client
	SendReceive       func(rqstpath string, a ...interface{}) (rw ReaderWriter, err error)
	SendRespondString func(rqstpath string, method string, rqstheaders map[string]string, a ...interface{}) (rspstr string, err error)
	Send              func(rqstpath string, method string, rqstheaders map[string]string, a ...interface{}) (rspr io.Reader, err error)
	Close             func()
}

//Client - struct
type Client struct {
	httpclient *http.Client
}

//NewClient - instance
func NewClient() (clnt *Client) {
	clnt = &Client{httpclient: &http.Client{}}
	return
}

//Close *Client
func (clnt *Client) Close() {
	if clnt != nil {
		if clnt.httpclient != nil {
			clnt.httpclient.CloseIdleConnections()
			clnt = nil
		}
		clnt = nil
	}
}

//ReaderWriter interface
type ReaderWriter interface {
	io.ReadWriteCloser
	Flush() error
}

//SendReceive return ReaderWriter that implement io.Reader,io.Writer
func (clnt *Client) SendReceive(rqstpath string, a ...interface{}) (rw ReaderWriter, err error) {
	if strings.HasPrefix(rqstpath, "ws:") || strings.HasPrefix(rqstpath, "wss://") {
		if c, _, cerr := websocket.DefaultDialer.Dial(rqstpath, nil); cerr == nil {
			rw = ws.NewReaderWriter(c)
		} else {
			err = cerr
		}
	}
	return
}

//SendRespondString - Client Send but return response as string
func (clnt *Client) SendRespondString(rqstpath string, rqstheaders map[string]string /* rspheaders map[string]string,*/, a ...interface{}) (rspstr string, err error) {
	var rspr io.Reader = nil
	rspstr = ""
	if rspr, err = clnt.Send(rqstpath, rqstheaders /* rspheaders,*/, a...); err == nil {
		if rspr != nil {
			rspstr, err = iorw.ReaderToString(rspr)
		}
	}
	return
}

//Send - Client send
func (clnt *Client) Send(rqstpath string, rqstheaders map[string]string /*rspheaders map[string]string,*/, a ...interface{}) (rspr io.Reader, err error) {
	if strings.HasPrefix(rqstpath, "http:") || strings.HasPrefix(rqstpath, "https://") {
		var method string = ""
		var r io.Reader = nil
		var w io.Writer = nil
		var aok bool = false
		var ai = 0
		var rntme active.Runtime = nil
		var onsucess interface{} = nil
		var onerror interface{} = nil
		for ai < len(a) {
			d := a[ai]
			if r == nil {
				if rs, rsok := d.(string); rsok {
					if rs != "" {
						r = strings.NewReader(rs)
					}
					if ai < len(a)-1 {
						a = append(a[:ai], a[ai+1:]...)
						continue
					} else {
						a = append(a[:ai], a[ai+1:]...)
						break
					}
				} else if r, aok = d.(io.Reader); aok {
					if ai < len(a)-1 {
						a = append(a[:ai], a[ai+1:]...)
						continue
					} else {
						a = append(a[:ai], a[ai+1:]...)
						break
					}
				} else if rntme, aok = d.(active.Runtime); aok {
					if ai < len(a)-1 {
						a = append(a[:ai], a[ai+1:]...)
						continue
					} else {
						a = append(a[:ai], a[ai+1:]...)
						break
					}
				} else if mp, mpok := d.(map[string]interface{}); mpok {
					if mp != nil && len(mp) > 0 {
						for k, v := range mp {
							if k == "sucess" {
								if onsucess == nil {
									onsucess = v
								}
							} else if k == "error" {
								if onerror == nil {
									onerror = v
								}
							} else if k == "method" {
								if mthd, _ := v.(string); mthd != "" && method == "" {
									method = strings.ToUpper(mthd)
								}
							}
						}
					}
					if ai < len(a)-1 {
						a = append(a[:ai], a[ai+1:]...)
						continue
					} else {
						a = append(a[:ai], a[ai+1:]...)
						break
					}
				}
			}
			if w == nil {
				if w, aok = d.(io.Writer); aok {
					if ai < len(a)-1 {
						a = append(a[:ai], a[ai+1:]...)
						continue
					} else {
						a = append(a[:ai], a[ai+1:]...)
						break
					}
				}
			}
			ai++
		}
		if r != nil {
			if method == "" || method == "GET" {
				method = "POST"
			}
		} else if method == "" {
			method = "GET"
		}
		var rqst, rqsterr = http.NewRequest(method, rqstpath, r)
		if rqsterr == nil {
			if len(rqstheaders) > 0 {
				for hdk, hdv := range rqstheaders {
					rqst.Header.Add(hdk, hdv)
				}
			}

			var resp, resperr = clnt.Do(rqst)
			if resperr == nil {
				if rntme != nil && onsucess != nil {
					rntme.InvokeFunction(onsucess, method, resp)
				}
				if scde := resp.StatusCode; scde >= 200 && scde < 300 {
					if respbdy := resp.Body; respbdy != nil {
						if w != nil {
							wg := &sync.WaitGroup{}
							wg.Add(1)
							pi, pw := io.Pipe()
							go func() {
								defer func() {
									pw.Close()
								}()
								wg.Done()
								if w != nil {
									io.Copy(pw, respbdy)
								}
							}()
							wg.Wait()
							io.Copy(w, pi)
						} else if rspr == nil {
							rspr = respbdy
						}
					}
				}
			} else {
				err = resperr
				if rntme != nil && onerror != nil {
					rntme.InvokeFunction(onerror, err, onerror, method, resp)
				}
			}
		} else {
			if rntme != nil && onerror != nil {
				rntme.InvokeFunction(onerror, err, onerror)
			}
		}
	}
	return
}

//Do - refer tp http.Client Do interface
func (clnt *Client) Do(rqst *http.Request) (rspnse *http.Response, err error) {
	rspnse, err = clnt.httpclient.Do(rqst)
	return
}

//DefaultClient  - default global web Client
var DefaultClient *Client

func init() {
	DefaultClient = NewClient()
}
