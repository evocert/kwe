package resources

//ResourcingManager - struct
type ResourcingManager struct {
	rsngendpnts map[string]*ResourcingEndpoint
}

//FindRS - find Resource
func (rsrcngmngr *ResourcingManager) FindRS(path string) (rs *Resource) {

	return
}

//NewResourcingManager - instance
func NewResourcingManager() (rsrcngmngr *ResourcingManager) {
	rsrcngmngr = &ResourcingManager{rsngendpnts: map[string]*ResourcingEndpoint{}}
	return
}

var glbrsrcngmngr *ResourcingManager

func init() {
	if glbrsrcngmngr == nil {
		glbrsrcngmngr = NewResourcingManager()
	}
}
