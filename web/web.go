package web

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

//Client - struct
type Client struct {
	httpclient *http.Client
}

//NewClient - instance
func NewClient() (clnt *Client) {
	clnt = &Client{httpclient: &http.Client{}}
	return
}

//Send - Client send
func (clnt *Client) Send(rqstpath string, rqstheaders map[string]string, rspheaders map[string]string, a ...interface{}) (rspr io.Reader, err error) {
	if strings.HasPrefix(rqstpath, "http:") || strings.HasPrefix(rqstpath, "https://") {
		var method = "GET"
		var r io.Reader = nil
		var w io.Writer = nil
		var aok bool = false
		var ai = 0

		for ai < len(a) {
			d := a[ai]
			if r == nil {
				if r, aok = d.(io.Reader); aok {
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
			method = "POST"
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
				if rspheaders != nil {
					for rsph, rsphv := range resp.Header {
						rspheaders[rsph] = strings.Join(rsphv, ";")
					}
				}
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
			} else {
				err = resperr
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
