package resources

import (
	"archive/zip"
	"net/url"
	"os"
	"strings"
)

//ResourcingEndpoint - struct
type ResourcingEndpoint struct {
	path        string
	epnttype    string
	islocal     bool
	isRemote    bool
	host        string
	schema      string
	querystring string
	rsngmngr    *ResourcingManager
}

func (rscngepnt *ResourcingEndpoint) dispose() {
	if rscngepnt != nil {
		if rscngepnt.rsngmngr != nil {
			delete(rscngepnt.rsngmngr.rsngendpntspaths, rscngepnt.path)
			rscngepnt.rsngmngr = nil
		}
		rscngepnt = nil
	}
}

func (rscngepnt *ResourcingEndpoint) findRS(path string) (rs *Resource) {
	if path != "" {
		if path = strings.TrimSpace(strings.Replace(path, "\\", "/", -1)); path != "" {
			if rscngepnt.islocal {
				var tmppath = ""
				var tmppaths = strings.Split(path, "/")
				for pn, ps := range tmppaths {
					if pn < len(tmppaths)-1 {
						if fi, fierr := os.Stat(rscngepnt.path + tmppath + ps + ".zip"); fierr == nil && !fi.IsDir() {
							var testpath = strings.Join(tmppaths[pn+1:len(tmppaths)], "/")
							if testpath != "" {
								if r, err := zip.OpenReader(rscngepnt.path + tmppath + ps + ".zip"); err == nil {
									for _, f := range r.File {
										if f.Name == testpath {
											if rc, rcerr := f.Open(); rcerr == nil {
												rs = newRS(rscngepnt, path, rc)
											} else {

											}
											return
										}
									}
								}
							}
							break
						} else {
							tmppath = tmppath + ps + "/"
						}
					} else {
						break
					}
				}
				if fi, fierr := os.Stat(rscngepnt.path + path); fierr == nil && !fi.IsDir() {
					if f, ferr := os.Open(rscngepnt.path + path); ferr == nil && f != nil {
						rs = newRS(rscngepnt, path, f)
					}
				}
			} else if rscngepnt.isRemote {

			}
		}
	}
	return
}

func nextResourcingEndpoint(rsngmngr *ResourcingManager, path string, a ...interface{}) (rsngepnt *ResourcingEndpoint, rsngepntpath string) {
	rsngepntpath = path
	if rsngepntpath != "" {
		rsngepntpath = strings.Replace(strings.TrimSpace(rsngepntpath), "\\", "/", -1)
		if fi, fierr := os.Stat(rsngepntpath); fierr == nil {
			if rsngepntpath != "/" && rune(rsngepntpath[len(rsngepntpath)-1]) != '/' {
				rsngepntpath = rsngepntpath + "/"
			}
			if fi.IsDir() {
				rsngepnt = &ResourcingEndpoint{rsngmngr: rsngmngr, islocal: true, isRemote: false, host: "", schema: "", querystring: "", path: rsngepntpath}
			}
		} else {
			_, err := url.ParseRequestURI(rsngepntpath)
			if err == nil {
				u, err := url.Parse(rsngepntpath)
				if err == nil && u.Scheme != "" && u.Host != "" {
					var querystring = ""
					if u.RawQuery == "" {
						querystring = ""
					} else {
						querystring = u.RawQuery
					}
					rsngepnt = &ResourcingEndpoint{rsngmngr: rsngmngr, islocal: false, isRemote: true, host: u.Host, schema: u.Scheme, querystring: querystring, path: path}
				}
			}
		}
	}
	return
}
