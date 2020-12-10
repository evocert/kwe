package resources

import "strings"

//ResourcingPath - struct
type ResourcingPath struct {
	Path       string
	LookupPath string
	rsngmngr   *ResourcingManager
}

//NewResourcingPath - instance
func NewResourcingPath(path string, rsngmngr *ResourcingManager) (rsngpth *ResourcingPath) {
	if rsngmngr == nil {
		rsngmngr = glbrscngmngr
	}
	var lkppath = path
	var lkpi = strings.Index(path, "@")
	var lkpli = strings.LastIndex(path, "@")
	if lkpi >= 0 && lkpi < lkpli {
		lkppath = path[:lkpi] + path[lkpi+1:lkpli] + path[lkpli+1:]
	}
	rsngpth = &ResourcingPath{Path: path, LookupPath: lkppath, rsngmngr: rsngmngr}
	return
}

//Close - refer to io.Closer
func (rsngpth *ResourcingPath) Close() (err error) {
	if rsngpth != nil {
		if rsngpth.rsngmngr != nil {
			rsngpth.rsngmngr = nil
		}
		rsngpth = nil
	}
	return
}

//ResourceHandler - instance of Resource Handler
func (rsngpth *ResourcingPath) ResourceHandler(altpath ...string) (rshndlr *ResourceHandler) {
	if rsngpth != nil && rsngpth.rsngmngr != nil {
		if len(altpath) > 0 && altpath[0] != "" {
			if rs, rserr := rsngpth.rsngmngr.FindRS(altpath[0]); rs != nil && rserr == nil {
				rshndlr = newResourceHandler(rs)
			}
		} else if rs, rserr := rsngpth.rsngmngr.FindRS(rsngpth.LookupPath); rs != nil && rserr == nil {
			rshndlr = newResourceHandler(rs)
		}
	}
	return
}
