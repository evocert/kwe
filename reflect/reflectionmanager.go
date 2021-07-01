package reflect

import (
	"sync"
)

type ReflectManager struct {
	rlctrs   map[string]*Reflector
	rlctonrs map[interface{}]*Reflector

	lck *sync.RWMutex
}

func rflctcall(rfltmngr *ReflectManager, refname string, callname string, args ...interface{}) (rval interface{}, err error) {
	if rfltmngr != nil {
		func() {
			rfltmngr.lck.RLock()
			defer rfltmngr.lck.RUnlock()
			if rflctr := rfltmngr.rlctrs[refname]; rflctr != nil {
				rval, err = call(rflctr, callname, args...)
			}
		}()
	}
	return
}

func (rfltmngr *ReflectManager) Register(owner interface{}, refname string) {
	if rfltmngr != nil && owner != nil {
		func() {
			rfltmngr.lck.Lock()
			defer rfltmngr.lck.RUnlock()
			var rflctr *Reflector = nil
			if _, rflctonrok := rfltmngr.rlctonrs[owner]; !rflctonrok {
				if rflctr = NewReflector(owner); rflctr != nil {
					rfltmngr.rlctonrs[owner] = rflctr
				}
			} else {
				rflctr = rfltmngr.rlctonrs[owner]
			}
			if rflctr != nil {
				rfltmngr.rlctrs[refname] = rflctr
			}
		}()
	}
}
