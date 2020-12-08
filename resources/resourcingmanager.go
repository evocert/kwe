package resources

import "strings"

//ResourcingManager - struct
type ResourcingManager struct {
	rsngendpnts      map[string]string
	rsngendpntspaths map[string]*ResourcingEndpoint
}

//RegisterEndpoint - register ResourcingEndPoint
func (rscngmngr *ResourcingManager) RegisterEndpoint(epntpath string, path string, prms ...interface{}) {
	if epntpath != "" && path != "" {
		if _, rsngepntok := rscngmngr.rsngendpnts[epntpath]; !rsngepntok {
			if newrsngepnt, newrsngepntpath := nextResourcingEndpoint(rscngmngr, path, prms...); newrsngepnt != nil {
				rsngepnt, rsngepntok := rscngmngr.rsngendpntspaths[newrsngepntpath]
				if rsngepntok {
					if rsngepnt != newrsngepnt {
						rsngepnt.dispose()
						rscngmngr.rsngendpntspaths[newrsngepntpath] = newrsngepnt
						rscngmngr.rsngendpnts[epntpath] = newrsngepntpath
					}
				} else {
					rscngmngr.rsngendpntspaths[newrsngepntpath] = newrsngepnt
					rscngmngr.rsngendpnts[epntpath] = newrsngepntpath
				}
			}
		} else {
			if rscngmngr.rsngendpnts[epntpath] != path {
				if newrsngepnt, newrsngepntpath := nextResourcingEndpoint(rscngmngr, path, prms...); newrsngepnt != nil {
					rsngepnt, rsngepntok := rscngmngr.rsngendpntspaths[newrsngepntpath]
					if rsngepntok {
						if rsngepnt != newrsngepnt {
							rsngepnt.dispose()
							rscngmngr.rsngendpntspaths[newrsngepntpath] = newrsngepnt
							rscngmngr.rsngendpnts[epntpath] = newrsngepntpath
						}
					} else {
						rscngmngr.rsngendpntspaths[newrsngepntpath] = newrsngepnt
						rscngmngr.rsngendpnts[epntpath] = newrsngepntpath
					}
				}
			}
		}
	}
}

//FindRS - find Resource
func (rscngmngr *ResourcingManager) FindRS(path string) (rs *Resource) {
	if path != "" {
		path = strings.Replace(path, "\\", "/", -1)
		if rune(path[0]) != '/' {
			path = "/" + path
		}
		if path == "/" {
			return
		}
		var rspthFound = ""

		for rsgnpath := range rscngmngr.rsngendpnts {
			if len(rsgnpath) > len(rspthFound) && strings.HasPrefix(path, rsgnpath) {
				if len(rsgnpath) > len(rspthFound) {
					rspthFound = rsgnpath
				}
			}
		}
		if len(rspthFound) > 0 {
			rs = rscngmngr.rsngendpntspaths[rscngmngr.rsngendpnts[rspthFound]].findRS(path[len(rspthFound):])
		}
	}
	return
}

//NewResourcingManager - instance
func NewResourcingManager() (rscngmngr *ResourcingManager) {
	rscngmngr = &ResourcingManager{rsngendpntspaths: map[string]*ResourcingEndpoint{}, rsngendpnts: map[string]string{}}
	return
}

var glbrscngmngr *ResourcingManager

//GLOBALRSNGMANAGER - GLOBAL ResourcingManager for app
func GLOBALRSNGMANAGER() *ResourcingManager {
	return glbrscngmngr
}

func init() {
	if glbrscngmngr == nil {
		glbrscngmngr = NewResourcingManager()
		glbrscngmngr.RegisterEndpoint("/", "./")
	}
}
