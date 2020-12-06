package database

import (
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/web"
)

//EndPoint - struct
type EndPoint struct {
	datasource string
	args       []interface{}
	clnt       *web.Client
}

func newEndPoint(datasource string, a ...interface{}) (endpnt *EndPoint) {
	endpnt = &EndPoint{datasource: datasource, args: a, clnt: web.NewClient()}
	return
}

func (endpnt *EndPoint) query(exctr *Executor, forrows bool, out io.Writer, iorags ...interface{}) (err error) {
	pi, pw := io.Pipe()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	func() {
		defer func() {
			pi.Close()
			pi = nil
			wg = nil
		}()
		go func() {
			defer func() {
				pw.Close()
			}()
			wg.Done()
			encw := json.NewEncoder(pw)
			rqstmpstngs := map[string]interface{}{}
			if len(exctr.mappedArgs) > 0 {
				for kmp, vmp := range exctr.mappedArgs {
					rqstmpstngs[kmp] = vmp
				}
			}
			if forrows {
				rqstmpstngs["query"] = exctr.stmnt
			} else {
				rqstmpstngs["execute"] = exctr.stmnt
			}

			rqstmp := map[string]interface{}{"1": rqstmpstngs}
			encw.Encode(&rqstmp)
			encw = nil
			rqstmp = nil
		}()

		if strings.HasPrefix(endpnt.datasource, "http://") || strings.HasPrefix(endpnt.datasource, "https://") {
			func() {
				var rspheaders = map[string]string{}
				var rqstheaders = map[string]string{}
				rqstheaders["Content-Type"] = "application/json"
				endpnt.clnt.Send(endpnt.datasource, rqstheaders, rspheaders, pi, out)
				rqstheaders = nil
				rspheaders = nil
			}()
		}
	}()
	return
}
