package caching

import (
	"sync"

	"github.com/evocert/kwe/enumeration"
)

type MapHandler struct {
	*Map
}

func NewMapHandler(mp ...*Map) (mphndlr *MapHandler) {
	if len(mp) == 1 && mp[0] != nil {
		mphndlr = &MapHandler{Map: mp[0]}
	} else {
		mphndlr = &MapHandler{Map: NewMap()}
	}
	return
}

func (mphndlr *MapHandler) Put(a ...interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		mapPut(mphndlr.Map, mphndlr, a...)
	}
}

type Map struct {
	lck    *sync.RWMutex
	keys   *enumeration.List
	kvndm  map[*enumeration.Node]*enumeration.Node
	vkndm  map[*enumeration.Node]*enumeration.Node
	values *enumeration.List
}

//NewMap return instance of *Map
func NewMap() (mp *Map) {
	mp = &Map{
		lck:    &sync.RWMutex{},
		keys:   enumeration.NewList(true),
		kvndm:  map[*enumeration.Node]*enumeration.Node{},
		values: enumeration.NewList(),
		vkndm:  map[*enumeration.Node]*enumeration.Node{}}
	return
}

func (mp *Map) Handler() (mphndlr *MapHandler) {
	mphndlr = NewMapHandler(mp)
	return
}

func (mp *Map) Find(ks ...interface{}) (vs []interface{}) {
	return
}

func (mp *Map) Keys() (ks []interface{}) {
	mp.keys.Do(func(knde *enumeration.Node, val interface{}) bool {
		if ks == nil {
			ks = []interface{}{}
		}
		ks = append(ks, val)
		return false
	}, nil, nil)
	return
}

func (mp *Map) Values() (vs []interface{}) {
	mp.values.Do(func(knde *enumeration.Node, val interface{}) bool {
		if vs == nil {
			vs = []interface{}{}
		}
		vs = append(vs, val)
		return false
	}, nil, nil)
	return
}

func (mp *Map) Put(a ...interface{}) {
	mapPut(mp, nil, a...)
}

func mapPut(mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		keys := mp.keys
		values := mp.values
		for {
			if al := len(a); al > 0 {
				if al%2 == 0 {
					k := a[0]
					if k != nil {
						a = a[1:]
						v := a[0]
						func() {
							if mphndlr != nil {
								mp.lck.Lock()
								defer mp.lck.Unlock()
							}
							keys.Add(k)
							values.Add(v)
						}()
						a = a[1:]
					}
				} else {
					if m, mok := a[0].(map[string]interface{}); mok && len(m) > 0 {
						a = a[1:]
						for k, v := range m {
							a = append(a, k, v)
						}
					} else if mi, miok := a[0].(map[interface{}]interface{}); miok && len(mi) > 0 {
						a = a[1:]
						for k, v := range mi {
							a = append(a, k, v)
						}
					} else {
						a = a[1:]
					}
				}
			} else {
				break
			}
		}
	}
}

var glbmphndlr *MapHandler = nil

func init() {
	glbmphndlr = NewMapHandler(NewMap())
}
