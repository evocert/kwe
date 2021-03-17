package resources

import (
	"archive/zip"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/web"
)

//ResourcingEndpoint - struct
type ResourcingEndpoint struct {
	lck               *sync.Mutex
	fsutils           *fsutils.FSUtils
	path              string
	epnttype          string
	isLocal           bool
	isRemote          bool
	isEmbedded        bool
	remoteHeaders     map[string]string
	host              string
	schema            string
	querystring       string
	embeddedResources map[string]interface{}
	rsngmngr          *ResourcingManager
}

//FS return fsutils.FSUtils implementation for *ResourcingEndPoint
func (rscngepnt *ResourcingEndpoint) FS() *fsutils.FSUtils {
	if rscngepnt.fsutils == nil {
		rscngepnt.fsutils = &fsutils.FSUtils{
			FIND: func(path ...string) (finfos []fsutils.FileInfo) {
				finfos, _ = rscngepnt.fsfind(path...)
				return
			}, LS: func(path ...string) (finfos []fsutils.FileInfo) {
				return
			}, MKDIR: func(path string) bool {
				return rscngepnt.fsmkdir(path)
			}, MKDIRALL: func(path string) bool {
				return rscngepnt.fsmkdirall(path)
			}, RM: func(path string) bool {
				return rscngepnt.fsrm(path)
			}, MV: func(path string, destpath string) bool {
				return rscngepnt.fsmv(path, destpath)
			}, TOUCH: func(path string) bool {
				return rscngepnt.fstouch(path)
			}, CAT: func(path string) string {
				return rscngepnt.fscat(path)
			}, SET: func(path string, a ...interface{}) bool {
				return rscngepnt.fsset(path, a...)
			}, APPEND: func(path string, a ...interface{}) bool {
				return rscngepnt.fsappend(path, a...)
			},
		}
	}
	return rscngepnt.fsutils
}

func (rscngepnt *ResourcingEndpoint) fsappend(path string, a ...interface{}) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && strings.LastIndex(path, ".") > 0 && (strings.LastIndex(path, "/") == -1 || strings.LastIndex(path, ".") > strings.LastIndex(path, "/")) {
			if err := fsutils.APPEND(rscngepnt.path+path, a...); err == nil {
				return true
			}
		}
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsset(path string, a ...interface{}) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && strings.LastIndex(path, ".") > 0 && (strings.LastIndex(path, "/") == -1 || strings.LastIndex(path, ".") > strings.LastIndex(path, "/")) {
			if err := fsutils.SET(rscngepnt.path+path, a...); err == nil {
				return true
			}
		}
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fscat(path string) (s string) {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && strings.LastIndex(path, ".") > 0 && (strings.LastIndex(path, "/") == -1 || strings.LastIndex(path, ".") > strings.LastIndex(path, "/")) {
			s, _ = fsutils.CAT(rscngepnt.path + path)
		}
	}
	return s
}

func (rscngepnt *ResourcingEndpoint) fstouch(path string) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && strings.LastIndex(path, ".") > 0 && (strings.LastIndex(path, "/") == -1 || strings.LastIndex(path, ".") > strings.LastIndex(path, "/")) {
			if err := fsutils.TOUCH(rscngepnt.path + path); err != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsmv(path string, destpath string) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" {
			if destpath = strings.Replace(strings.TrimSpace(destpath), "\\", "/", -1); destpath != "" {
				if err := fsutils.MV(rscngepnt.path+path, rscngepnt.path+destpath); err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsrm(path string) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" {
			if err := fsutils.RM(rscngepnt.path + path); err != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsmkdirall(path string) bool {
	if rscngepnt.isLocal {

		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && (strings.LastIndex(path, ".") == -1 || strings.LastIndex(path, ".") < strings.LastIndex(path, "/")) {
			lklpath := rscngepnt.path + strings.TrimSpace(strings.Replace(path, "\\", "/", -1))
			if strings.LastIndex(lklpath, "/") > 0 && strings.HasSuffix(lklpath, "/") {
				lklpath = lklpath[:len(lklpath)-1]
			}
			if err := fsutils.MKDIRALL(rscngepnt.path + lklpath); err != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsmkdir(path string) bool {
	if rscngepnt.isLocal {
		if path = strings.Replace(strings.TrimSpace(path), "\\", "/", -1); path != "" && (strings.LastIndex(path, ".") == -1 || strings.LastIndex(path, ".") < strings.LastIndex(path, "/")) {
			lklpath := rscngepnt.path + strings.TrimSpace(strings.Replace(path, "\\", "/", -1))
			if strings.LastIndex(lklpath, "/") > 0 && strings.HasSuffix(lklpath, "/") {
				lklpath = lklpath[:len(lklpath)-1]
			}
			if err := fsutils.MKDIR(rscngepnt.path + lklpath); err != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (rscngepnt *ResourcingEndpoint) fsls(path ...string) (finfos []fsutils.FileInfo) {
	if rscngepnt.isLocal {
		lklpath := rscngepnt.path + strings.TrimSpace(strings.Replace(path[0], "\\", "/", -1))
		if strings.LastIndex(lklpath, "/") > 0 && strings.HasSuffix(lklpath, "/") {
			lklpath = lklpath[:len(lklpath)-1]
		}
		if len(path) == 1 {
			finfos, _ = fsutils.LS(lklpath, strings.TrimSpace(strings.Replace(path[0], "\\", "/", -1)))
		} else if len(path) == 2 {
			finfos, _ = fsutils.LS(lklpath, strings.TrimSpace(strings.Replace(path[1], "\\", "/", -1)))
		}
	}
	return
}

func (rscngepnt *ResourcingEndpoint) fsfind(path ...string) (finfos []fsutils.FileInfo, err error) {
	if rscngepnt.isLocal {
		lklpath := rscngepnt.path + strings.TrimSpace(strings.Replace(path[0], "\\", "/", -1))
		if strings.LastIndex(lklpath, "/") > 0 && strings.HasSuffix(lklpath, "/") {
			lklpath = lklpath[:len(lklpath)-1]
		}
		if len(path) == 1 {
			finfos, _ = fsutils.FIND(lklpath, strings.TrimSpace(strings.Replace(path[0], "\\", "/", -1)))
		} else if len(path) == 2 {
			finfos, _ = fsutils.FIND(lklpath, strings.TrimSpace(strings.Replace(path[1], "\\", "/", -1)))
		}
	}
	return
}

func (rscngepnt *ResourcingEndpoint) dispose() {
	if rscngepnt != nil {
		if rscngepnt.rsngmngr != nil {
			rsendpath := rscngepnt.path
			delete(rscngepnt.rsngmngr.rsngrootpaths, rsendpath)
			for rspth, rsndpth := range rscngepnt.rsngmngr.rsngpaths {
				if rsndpth == rsendpath {
					delete(rscngepnt.rsngmngr.rsngpaths, rspth)
				}
			}
			rscngepnt.rsngmngr = nil
		}
		if rscngepnt.embeddedResources != nil {
			for embk := range rscngepnt.embeddedResources {
				rscngepnt.RemoveResource(embk)
			}
			rscngepnt.embeddedResources = nil
		}
		if rscngepnt.fsutils != nil {
			rscngepnt.fsutils = nil
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
				} else if rscngepnt.isLocal {
					var tmppath = ""
					var tmppaths = strings.Split(path, "/")
					for pn, ps := range tmppaths {
						if tmpl := len(tmppaths); pn < tmpl-1 {
							if fi, fierr := os.Stat(rscngepnt.path + tmppath + ps + ".zip"); fierr == nil && !fi.IsDir() {
								var testpath = strings.Join(tmppaths[pn+1:tmpl], "/")
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
					if r, rerr := web.DefaultClient.Send(rscngepnt.schema+"://"+strings.Replace(rscngepnt.host+rscngepnt.path+path, "//", "/", -1), remoteHeaders, nil, rqstr); rerr == nil {
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
			rscngepnt.embeddedResources[path] = nil
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

//Dirs return list of directories of a local endpoint
func (rscngepnt *ResourcingEndpoint) Dirs(lkppath ...string) (dirs []string, err error) {
	dirs = []string{}
	if len(lkppath) == 0 {
		err = filepath.Walk(rscngepnt.path, func(path string, info os.FileInfo, err error) error {
			path = strings.Replace(path, "\\", "/", -1)
			if info.IsDir() {
				if path != rscngepnt.path {
					dirs = append(dirs, path[len(rscngepnt.path):])
				}
			} else if strings.HasSuffix(path, ".zip") {
				dirs = append(dirs, path[len(rscngepnt.path):len(path)-len(".zip")])
			}
			return nil
		})
	} else {
		for _, lkp := range lkppath {
			err = filepath.Walk(rscngepnt.path+lkp, func(path string, info os.FileInfo, err error) error {
				path = strings.Replace(path, "\\", "/", -1)
				if info.IsDir() {
					if path != rscngepnt.path {
						dirs = append(dirs, path[len(rscngepnt.path+lkp):])
					}
				} else if strings.HasSuffix(path, ".zip") {
					dirs = append(dirs, path[len(rscngepnt.path+lkp):len(path)-len(".zip")])
				}
				return nil
			})
		}
	}
	if err != nil {

	}
	return
}

//Files return list of files of a local endpoint
func (rscngepnt *ResourcingEndpoint) Files(lkppath ...string) (files []string, err error) {
	files = []string{}
	if len(lkppath) == 0 {
		err = filepath.Walk(rscngepnt.path, func(path string, info os.FileInfo, err error) error {
			path = strings.Replace(path, "\\", "/", -1)
			if !info.IsDir() {
				if !strings.HasSuffix(path, ".zip") && path != rscngepnt.path {
					files = append(files, path[len(rscngepnt.path):])
				}
			}
			return nil
		})
	} else {
		for _, lkp := range lkppath {
			err = filepath.Walk(rscngepnt.path+lkp, func(path string, info os.FileInfo, err error) error {
				path = strings.Replace(path, "\\", "/", -1)
				if !info.IsDir() {
					if path != rscngepnt.path {
						files = append(files, path[len(rscngepnt.path+lkp):])
					}
				}
				return nil
			})
		}
	}
	if err != nil {

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
					rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, isLocal: false, isRemote: true, embeddedResources: map[string]interface{}{}, host: u.Host, schema: u.Scheme, querystring: querystring, path: path}
				}
			}
		} else {
			if fi, fierr := os.Stat(rsngepntpath); fierr == nil {
				if rsngepntpath != "/" && rune(rsngepntpath[len(rsngepntpath)-1]) != '/' {
					rsngepntpath = rsngepntpath + "/"
				}
				if fi.IsDir() {
					rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, isLocal: true, isRemote: false, isEmbedded: false, embeddedResources: map[string]interface{}{}, host: "", schema: "", querystring: "", path: rsngepntpath}
				}
			}
		}
	} else {
		rsngepnt = &ResourcingEndpoint{lck: &sync.Mutex{}, rsngmngr: rsngmngr, isLocal: false, isRemote: false, isEmbedded: true, embeddedResources: map[string]interface{}{}, host: "", schema: "", querystring: "", path: ""}
	}
	return
}
