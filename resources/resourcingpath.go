package resources

//ResourcingPath - struct
type ResourcingPath struct {
	Path     string
	rsngmngr *ResourcingManager
}

//NewResourcingPath - instance
func NewResourcingPath(path string, rsngmngr *ResourcingManager) (rsngpth *ResourcingPath) {
	if rsngmngr == nil {
		rsngmngr = glbrscngmngr
	}
	rsngpth = &ResourcingPath{Path: path, rsngmngr: rsngmngr}
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
func (rsngpth *ResourcingPath) ResourceHandler() (rshndlr *ResourceHandler) {
	if rsngpth != nil && rsngpth.rsngmngr != nil {
		if rs := rsngpth.rsngmngr.FindRS(rsngpth.Path); rs != nil {
			rshndlr = newResourceHandler(rs)
		}
	}
	return
}
