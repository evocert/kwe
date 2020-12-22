package resources

import (
	"bufio"
	"io"
	"strings"
)

//ResourcingManager - struct
type ResourcingManager struct {
	rsngpaths     map[string]string
	rsngrootpaths map[string]*ResourcingEndpoint
}

//RemovePathResource - Remove Endpoint Resource via path
func (rscngmngr *ResourcingManager) RemovePathResource(path string) (rmvd bool) {
	if path != "" {
		path = strings.Replace(path, "\\", "/", -1)
		if rune(path[0]) != '/' {
			path = "/" + path
		}
		if path == "/" {
			return
		}
		var rspthFound = ""

		for rsgnpath := range rscngmngr.rsngpaths {
			if len(rsgnpath) > len(rspthFound) && strings.HasPrefix(path, rsgnpath) {
				if len(rsgnpath) > len(rspthFound) {
					rspthFound = rsgnpath
				}
			}
		}
		if len(rspthFound) > 0 {
			rmvd = rscngmngr.rsngrootpaths[rscngmngr.rsngpaths[rspthFound]].RemoveResource(path[len(rspthFound):])
		}
	}
	return
}

//EndpointViaRootPath return ResourcingEndpoint via root path
func (rscngmngr *ResourcingManager) EndpointViaRootPath(rootpath string) (rsngendpt *ResourcingEndpoint) {
	if rootpath != "" {
		rsngendpt = rscngmngr.rsngrootpaths[rootpath]
	}
	return
}

//EndpointViaPath return ResourcingEndpoint via path
func (rscngmngr *ResourcingManager) EndpointViaPath(path string) (rsngendpt *ResourcingEndpoint) {
	if path != "" {
		if endpntpth, endpntpthok := rscngmngr.rsngpaths[path]; endpntpthok {
			rsngendpt = rscngmngr.rsngrootpaths[endpntpth]
		}
	}
	return
}

//EndpointResource - Endpoint embedded resource via path
func (rscngmngr *ResourcingManager) EndpointResource(path string) (epntrs interface{}) {
	if path != "" {
		path = strings.Replace(path, "\\", "/", -1)
		if rune(path[0]) != '/' {
			path = "/" + path
		}
		if path == "/" {
			return
		}
		var rspthFound = ""

		for rsgnpath := range rscngmngr.rsngpaths {
			if len(rsgnpath) > len(rspthFound) && strings.HasPrefix(path, rsgnpath) {
				if len(rsgnpath) > len(rspthFound) {
					rspthFound = rsgnpath
				}
			}
		}
		if len(rspthFound) > 0 {
			epntrs = rscngmngr.rsngrootpaths[rscngmngr.rsngpaths[rspthFound]].Resource(path[len(rspthFound):])
		}
	}
	return
}

//MapEndpointResources map multiple endpoint embedable resources
func (rscngmngr *ResourcingManager) MapEndpointResources(a ...interface{}) {
	var epntpath string = ""
	var path string = ""
	var resource interface{} = nil
	var epntok = false
	for {
		if al := len(a); al >= 3 {
			if epntpath, epntok = a[0].(string); epntok && epntpath != "" {
				if path, epntok = a[1].(string); epntok {
					resource = a[2]
					a = a[3:]
					if _, epntpathok := rscngmngr.rsngpaths[epntpath]; epntpathok {
						if rscngepnt := rscngmngr.rsngrootpaths[rscngmngr.rsngpaths[epntpath]]; rscngepnt != nil {
							rscngepnt.MapResource(path, resource)
						}
					}
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

//MapEndpointResource - inline resource -  can be either func() io.Reader, *iorw.Buffer
func (rscngmngr *ResourcingManager) MapEndpointResource(epntpath string, path string, resource interface{}) {
	if epntpath != "" && path != "" {
		if _, epntpathok := rscngmngr.rsngpaths[epntpath]; epntpathok {
			if rscngepnt := rscngmngr.rsngrootpaths[rscngmngr.rsngpaths[epntpath]]; rscngepnt != nil {
				rscngepnt.MapResource(path, resource)
			}
		}
	}
}

//UnregisterPaths unregister multiple paths
func (rscngmngr *ResourcingManager) UnregisterPaths(path ...string) {
	if len(path) > 0 {
		for _, pth := range path {
			if pth != "" {
				if pndpth, pthok := rscngmngr.rsngpaths[pth]; pthok {
					delete(rscngmngr.rsngpaths, pth)
					fndEndPtsh := false
					for _, ptepth := range rscngmngr.rsngpaths {
						if ptepth == pndpth {
							fndEndPtsh = true
							break
						}
					}
					if !fndEndPtsh {
						if rspnt := rscngmngr.rsngrootpaths[pndpth]; rspnt != nil {
							rspnt.dispose()
							rspnt = nil
						}
					}
				}
			}
		}
	}
}

var emptypaths []string = make([]string, 0)

//RegisteredRootPaths return registered rootpaths
func (rscngmngr *ResourcingManager) RegisteredRootPaths() (paths []string) {
	if rscngmngr != nil {
		if ln := len(rscngmngr.rsngrootpaths); ln > 0 {
			paths = make([]string, ln)
			pi := 0
			for pth := range rscngmngr.rsngrootpaths {
				paths[pi] = pth
				pi++
			}
			return paths
		}
	}
	return emptypaths
}

//RegisteredPaths return registered paths
func (rscngmngr *ResourcingManager) RegisteredPaths() (paths []string) {
	if rscngmngr != nil {
		if ln := len(rscngmngr.rsngpaths); ln > 0 {
			paths = make([]string, ln)
			pi := 0
			for pth := range rscngmngr.rsngpaths {
				paths[pi] = pth
				pi++
			}
			return paths
		}
	}
	return emptypaths
}

//UnregisterPath - register path string
func (rscngmngr *ResourcingManager) UnregisterPath(path string) (rmvd bool) {
	if path != "" {
		if pndpth, pthok := rscngmngr.rsngpaths[path]; pthok {
			delete(rscngmngr.rsngpaths, path)
			fndEndPtsh := false
			for _, ptepth := range rscngmngr.rsngpaths {
				if ptepth == pndpth {
					fndEndPtsh = true
					break
				}
			}
			if !fndEndPtsh {
				if rspnt := rscngmngr.rsngrootpaths[pndpth]; rspnt != nil {
					rspnt.dispose()
					rspnt = nil
				}
			}
		}
	}
	return
}

//UnregisterRootPaths unregister multiple RootPaths and their ResourcingEndPoints
func (rscngmngr *ResourcingManager) UnregisterRootPaths(epntpath ...string) {
	if len(epntpath) > 0 {
		for _, epth := range epntpath {
			if epth != "" {
				if rsndpt := rscngmngr.rsngrootpaths[epth]; rsndpt != nil {
					rsndpt.dispose()
				}
			}
		}
	}
}

//UnregisterRootPath unregister RootPath and dispose the ResourcingEndPoint
func (rscngmngr *ResourcingManager) UnregisterRootPath(epntpath string) (rmvd bool) {
	if epntpath != "" {
		if rsndpt := rscngmngr.rsngrootpaths[epntpath]; rsndpt != nil {
			rsndpt.dispose()
		}
	}
	return
}

//RegisterEndpoints register multiple Endpoints
func (rscngmngr *ResourcingManager) RegisterEndpoints(args ...interface{}) {
	var a []interface{} = nil
	var epntpath string = ""
	var path string = ""
	var argok = false
	for {
		if argsl := len(args); argsl >= 2 {
			if epntpath, argok = args[0].(string); argok {
				if path, argok = args[1].(string); argok {
					if argsl >= 3 {
						if a, argok = args[2].([]interface{}); argok {
							rscngmngr.RegisterEndpoint(epntpath, path, a...)
							args = args[3:]
						} else if argsl > 3 {
							rscngmngr.RegisterEndpoint(epntpath, path)
							args = args[2:]
						} else {
							break
						}
					} else {
						rscngmngr.RegisterEndpoint(epntpath, path)
						args = args[2:]
					}
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

//RegisterEndpoint - register ResourcingEndPoint
func (rscngmngr *ResourcingManager) RegisterEndpoint(path string, rootpath string, prms ...interface{}) {
	if path != "" {
		if _, rsngepntok := rscngmngr.rsngpaths[path]; !rsngepntok {
			if newrsngepnt, newrsngepntpath := nextResourcingEndpoint(rscngmngr, rootpath, prms...); newrsngepnt != nil {
				rsngepnt, rsngepntok := rscngmngr.rsngrootpaths[newrsngepntpath]
				if rsngepntok {
					if rsngepnt != newrsngepnt {
						rsngepnt.dispose()
						rscngmngr.rsngrootpaths[newrsngepntpath] = newrsngepnt
						rscngmngr.rsngpaths[path] = newrsngepntpath
					}
				} else {
					rscngmngr.rsngrootpaths[newrsngepntpath] = newrsngepnt
					rscngmngr.rsngpaths[path] = newrsngepntpath
				}
			}
		} else {
			if rscngmngr.rsngpaths[path] != rootpath {
				if newrsngepnt, newrsngepntpath := nextResourcingEndpoint(rscngmngr, rootpath, prms...); newrsngepnt != nil {
					rsngepnt, rsngepntok := rscngmngr.rsngrootpaths[newrsngepntpath]
					if rsngepntok {
						if rsngepnt != newrsngepnt {
							rsngepnt.dispose()
							rscngmngr.rsngrootpaths[newrsngepntpath] = newrsngepnt
							rscngmngr.rsngpaths[path] = newrsngepntpath
						}
					} else {
						rscngmngr.rsngrootpaths[newrsngepntpath] = newrsngepnt
						rscngmngr.rsngpaths[path] = newrsngepntpath
					}
				}
			}
		}
	}
}

//FindRSString - find Resource
func (rscngmngr *ResourcingManager) FindRSString(path string) (s string, err error) {
	if rs, rserr := rscngmngr.FindRS(path); rs != nil /*&& rs.isText*/ {
		func() {
			defer rs.Close()
			p := make([]rune, 1024)
			pi := 0
			buf := bufio.NewReader(rs)
			for {
				r, size, rerr := buf.ReadRune()
				if size > 0 {
					p[pi] = r
					pi++
					if pi == len(p) {
						pi = 0
						s += string(p[:])
					}
				}
				if rerr != nil {
					if rerr == io.EOF {
						rerr = nil
					} else {
						err = rerr
					}
					break
				}
			}
			if pi > 0 {
				s += string(p[:pi])
			}
		}()
	} else if rserr != nil {
		err = rserr
	}
	return
}

//FindRS - find Resource
func (rscngmngr *ResourcingManager) FindRS(path string) (rs *Resource, err error) {
	if path != "" {
		path = strings.Replace(path, "\\", "/", -1)
		if rune(path[0]) != '/' {
			path = "/" + path
		}
		if path == "/" {
			return
		}
		var rspthFound = ""

		for rsgnpath := range rscngmngr.rsngpaths {
			if len(rsgnpath) > len(rspthFound) && strings.HasPrefix(path, rsgnpath) {
				if len(rsgnpath) > len(rspthFound) {
					rspthFound = rsgnpath
				}
			}
		}
		if len(rspthFound) > 0 {
			rs, err = rscngmngr.rsngrootpaths[rscngmngr.rsngpaths[rspthFound]].findRS(path[len(rspthFound):])
		}
	}
	return
}

//NewResourcingManager - instance
func NewResourcingManager() (rscngmngr *ResourcingManager) {
	rscngmngr = &ResourcingManager{rsngrootpaths: map[string]*ResourcingEndpoint{}, rsngpaths: map[string]string{}}
	return
}

var glbrscngmngr *ResourcingManager

//GLOBALRSNG - GLOBAL Resourcing for app
func GLOBALRSNG() *ResourcingManager {
	return glbrscngmngr
}

func init() {
	if glbrscngmngr == nil {
		glbrscngmngr = NewResourcingManager()
		glbrscngmngr.RegisterEndpoint("/", "./")
	}
}
