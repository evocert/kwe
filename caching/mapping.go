package caching

type MapManager struct {
	maps map[string]*Map
}

func NewMapManager() (mpmmngr *MapManager) {
	mpmmngr = &MapManager{maps: map[string]*Map{}}
	return
}
