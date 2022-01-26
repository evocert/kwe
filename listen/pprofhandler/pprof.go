package pprofhandler

import (
	"net/http"
	"net/http/pprof"
	rtp "runtime/pprof"
	"strings"

	"github.com/evocert/kwe/listen"
)

var (
	cmdline = pprof.Cmdline
	profile = pprof.Profile
	symbol  = pprof.Symbol
	trace   = pprof.Trace
	index   = pprof.Index
)

func HandlePprof(w http.ResponseWriter, r *http.Request) (hndldpprof bool) {
	path := r.URL.Path
	if hndldpprof = strings.HasPrefix(path, "/debug/pprof/"); hndldpprof {
		w.Header().Set("Content-Type", "text/html")
		if strings.HasPrefix(path, "/debug/pprof/cmdline") {
			cmdline(w, r)
		} else if strings.HasPrefix(path, "/debug/pprof/profile") {
			profile(w, r)
		} else if strings.HasPrefix(path, "/debug/pprof/symbol") {
			symbol(w, r)
		} else if strings.HasPrefix(path, "/debug/pprof/trace") {
			trace(w, r)
		} else {
			for _, v := range rtp.Profiles() {
				ppName := v.Name()
				if strings.HasPrefix(path, "/debug/pprof/"+ppName) {
					pprof.Handler(ppName).ServeHTTP(w, r)
					return
				}
			}
			index(w, r)
		}
	}
	return
}

func init() {
	listen.DefaultHandlePprof = HandlePprof
}
