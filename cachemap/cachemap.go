package cachemap

import (
	"context"
	"sync"
)

type CacheMapHandle map[string]interface {
	/*Add(key string, a ...interface{})
	Get(key ...string) map[string]interface{}
	GetMap(key ...string) []CacheMapHandle
	Keys(key ...string) []string
	Value(key ...string) []interface{}*/
}

//CacheMap struct
type CacheMap struct {
	CacheMapHandle
	prtnmp *CacheMap
	topmp  *CacheMap
	lck    *sync.RWMutex
}

var glbcachedmap *CacheMap

//NewCacheMap return instance of new *CacheMap
//prntmap *CacheMap - parent *CacheMap
//a ...interface{} - initial arguments
func NewCacheMap(prntmap *CacheMap, a ...interface{}) (chmp *CacheMap) {
	var topmp *CacheMap = nil
	if prntmap != nil {
		if prntmap.topmp != nil {
			topmp = prntmap.topmp
		} else {
			topmp = prntmap
		}
	}
	chmp = &CacheMap{CacheMapHandle: CacheMapHandle{}, lck: &sync.RWMutex{}, prtnmp: prntmap, topmp: topmp}
	if topmp == nil {
		chmp.topmp = chmp
	}
	if len(a) > 0 {
		cntxt := context.Background()
		chmp.addInternal(cntxt, false, a...)
		<-cntxt.Done()
	}
	return
}

//Add
func (chdmp *CacheMap) Add(key string, a ...interface{}) {
	cntxt, cnl := context.WithCancel(context.Background())
	if len(a) > 0 {
		chdmp.addInternal(cntxt, true, key, a)
	} else {
		chdmp.addInternal(cntxt, true, key, map[string]interface{}{})
	}
	<-cntxt.Done()
	cnl()
}

func (chdmp *CacheMap) addInternal(prntcntxt context.Context, canLock bool, a ...interface{}) {
	if len(a) > 0 {
		var cntxtcnl context.CancelFunc = nil
		if prntcntxt != nil {
			_, cntxtcnl = context.WithCancel(prntcntxt)
		}
		var unlck bool = false
		defer func() {
			if unlck {
				chdmp.lck.Unlock()
			}
			if cntxtcnl != nil {
				cntxtcnl()
				cntxtcnl = nil
			}
		}()

		lckChdmp := func() {
			if canLock && !unlck {
				unlck = true
				chdmp.lck.Lock()
			}
		}

		addkeval := func(ky string, vl interface{}) {
			if mp, mpok := vl.(map[string]interface{}); mpok {
				nxtchdmp := NewCacheMap(chdmp, mp)
				lckChdmp()
				chdmp.CacheMapHandle[ky] = nxtchdmp
			} else if arr, arrok := vl.([]interface{}); arrok {
				lckChdmp()
				chdmp.CacheMapHandle[ky] = arr
			} else if fnc, fncok := vl.(func(...interface{}) (interface{}, error)); fncok {
				lckChdmp()
				chdmp.CacheMapHandle[ky] = fnc
			} else {
				lckChdmp()
				chdmp.CacheMapHandle[ky] = vl
			}
		}
		for {
			if al := len(a); al > 0 {
				if mp, mpok := a[0].(map[string]interface{}); mpok && len(mp) > 0 {
					for k, v := range mp {
						addkeval(k, v)
					}
				} else if al%2 == 0 {
					if k, kok := a[0].(string); kok && k != "" {
						a = a[1:]
						addkeval(k, a[0])
						a = a[1:]
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
}

func init() {
	if glbcachedmap == nil {
		glbcachedmap = NewCacheMap(nil)
	}
}
