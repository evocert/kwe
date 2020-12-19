package resources

import (
	"archive/zip"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/web"
)

//ResourcingEndpoint - struct
type ResourcingEndpoint struct {
	lck               *sync.Mutex
	path              string
	epnttype          string
	islocal           bool
	isRemote          bool
	isEmbedded        bool
	remoteHeaders     map[string]string
	host              string
	schema            string
	querystring       string
	embeddedResources map[string]interface{}
	rsngmngr          *ResourcingManager
}

func (rscngepnt *ResourcingEndpoint) dispose() {
	if rscngepnt != nil {
		if rscngepnt.rsngmngr != nil {
			rsendpath := rscngepnt.path
			delete(rscngepnt.rsngmngr.rsngendpntspaths, rsendpath)
			for rspth, rsndpth := range rscngepnt.rsngmngr.rsngendpnts {
				if rsndpth == rsendpath {
					delete(rscngepnt.rsngmngr.rsngendpnts, rspth)
				}
			}
			rscngepnt.rsngmngr = nil
		}
		rscngepnt = nil
		if rscngepnt.embeddedResources != nil {
			for embk := range rscngepnt.embeddedResources {
				rscngepnt.RemoveResource(embk)
			}
			rscngepnt.embeddedResources = nil
		}
		rscngepnt = nil
	}
}

func (rscngepnt *ResourcingEndpoint) findRS(path string) (rs *Resource, err error) {
	if path != "" {
		func() {
			rscngepnt.lck.Lock()
			defer rscngepnt.lck.Unlock()
			if path = strings.TrimSpace(strings.Replace(path, "\\", "/", -1)); path != "" {
				if embdrs, embdrsok := rscngepnt.embeddedResources[path]; embdrsok {
					if buff, buffok := embdrs.(*iorw.Buffer); buffok {
						rs = newRS(rscngepnt, path, buff.Reader())
					} else if funcr, funcrok := embdrs.(func() io.Reader); funcrok {
						rs = newRS(rscngepnt, path, funcr())
					}
				} else if rscngepnt.islocal {
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
													err = rcerr
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
						} else if ferr != nil {
							err = ferr
						}
					}
				} else if rscngepnt.isRemote {
					prms := map[string]interface{}{}
					if rscngepnt.querystring != "" {

					}
					remoteHeaders := map[string]string{}
					mimetype, _ := mimes.FindMimeType(path, "text/plain")
					var rqstr io.Reader = nil
					if mimetype == "application/json" {
						if len(prms) > 0 {

						}
					}
					remoteHeaders["Content-Type"] = mimetype
					if r, rerr := web.DefaultClient.Send(rscngepnt.schema+"://"+rscngepnt.host+rscngepnt.path+path, remoteHeaders, nil, rqstr); rerr == nil {
						rs = newRS(rscngepnt, path, r)
					} else if rerr != nil {
						err = rerr
					}
				}
			}
		}()
	}
	return
}

//RemoveResource - remove inline resource - true if found and removed and false if not exists
func (rscngepnt *ResourcingEndpoint) RemoveResource(path string) (rmvd bool) {
	if path != "" {
		if rs, rsok := rscngepnt.embeddedResources[path]; rsok {
			rmvd = rsok
			delete(rscngepnt.embeddedResources, path)
			if rs != nil {
				if bf, bfok := rs.(*iorw.Buffer); bfok && bf != nil {
					bf.Close()
					bf = nil
				}
			}
		}
	}
	return
}

//Resource - return mapped resource interface{} by path
func (rscngepnt *ResourcingEndpoint) Resource(path string) (rs interface{}) {
	if path != "" {
		rs, _ = rscngepnt.embeddedResources[path]
	}
	return
}

//Resources list of embedded resource paths
func (rscngepnt *ResourcingEndpoint) Resources() (rsrs []string) {
	if lrsrs := len(rscngepnt.embeddedResources); lrsrs > 0 {
		rsrs = make([]string, lrsrs)
		rsrsi := 0
		for rsrsk := range rscngepnt.embeddedResources {
			rsrs[rsrsi] = rsrsk
			rsrsi++
		}
	}
	return
}

//MapResource - inline resource -  can be either func() io.Reader, *iorw.Buffer
func (rscngepnt *ResourcingEndpoint) MapResource(path string, resource interface{}) {
	if path != "" && resource != nil {
		var validResource = false
		var strng = ""
		var isReader = false
		var r io.Reader = nil
		var isBuffer = false
		var buff *iorw.Buffer = nil

		if strng, validResource = resource.(string); !validResource {
			if _, validResource = resource.(func() io.Reader); !validResource {
				if buff, validResource = resource.(*iorw.Buffer); !validResource {
					if r, validResource = resource.(io.Reader); validResource {
						validResource = (r != nil)
					}
					isReader = validResource
				} else {
					isBuffer = true
				}
			}
		} else {
			if strng != "" {
				r = strings.NewReader(strng)
				isReader = true
			} else {
				validResource = false
			}
		}
		if validResource {
			if isReader {
				buff := iorw.NewBuffer()
				io.Copy(buff, r)
				resource = buff
			}
			if _, resourceok := rscngepnt.embeddedResources[path]; resourceok && rscngepnt.embeddedResources[path] != resource {
				if rscngepnt.embeddedResources[path] != nil {
					if isBuffer {
						if buff, isBuffer = rscngepnt.embeddedResources[path].(*iorw.Buffer); isBuffer {
							buff.Close()
							buff = nil
						}
						rscngepnt.embeddedResources[path] = resource
					}
				} else {
					rscngepnt.embeddedResources[path] = resource
				}
			} else {
				rscngepnt.embeddedResources[path] = resource
			}
		}
	}
}

func nextResourcingEndpoint(rsngmngr *ResourcingManager, path string, a ...interface{}) (rsngepnt *ResourcingEndpoint, rsngepntpath string) {
	rsngepntpath = path
	if rsngepntpath != "" {
		rsngepntpath = strings.Replace(strings.TrimSpace(rsngepntpath), "\\", "/", -1)
		if strings.HasPrefix(rsngepntpath, "http://") || strings.HasPrefix(rsngepntpath, "https://") {
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
					path = u.Path
					rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, islocal: false, isRemote: true, embeddedResources: map[string]interface{}{}, host: u.Host, schema: u.Scheme, querystring: querystring, path: path}
				}
			}
		} else {
			if fi, fierr := os.Stat(rsngepntpath); fierr == nil {
				if rsngepntpath != "/" && rune(rsngepntpath[len(rsngepntpath)-1]) != '/' {
					rsngepntpath = rsngepntpath + "/"
				}
				if fi.IsDir() {
					rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, islocal: true, isRemote: false, isEmbedded: false, embeddedResources: map[string]interface{}{}, host: "", schema: "", querystring: "", path: rsngepntpath}
				}
			}
		}
	} else {
		rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, islocal: false, isRemote: false, isEmbedded: true, embeddedResources: map[string]interface{}{}, host: "", schema: "", querystring: "", path: ""}
	}
	return
}
