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
func (clnt *Client) Send(rqstpath string, rqstheaders map[string]string, rspheaders map[string]string, r io.Reader, w io.Writer, a ...interface{}) (err error) {
	if strings.HasPrefix(rqstpath, "http:") || strings.HasPrefix(rqstpath, "https://") {
		var method = "GET"

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
				wg := &sync.WaitGroup{}
				wg.Add(1)
				pi, pw := io.Pipe()
				go func() {
					defer func() {
						pw.Close()
					}()
					wg.Done()
					if w != nil {
						if respbdy := resp.Body; respbdy != nil {
							io.Copy(pw, respbdy)
						}
					}
				}()
				wg.Wait()
				io.Copy(w, pi)
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
