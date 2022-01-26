package channeling

import (
	"sync"
	"sync/atomic"

	"github.com/evocert/kwe/api"
)

type Sessions struct {
	ssnscnt         uint32
	ssnslck         *sync.RWMutex
	mnssns          *Sessions
	ssns            map[*Session]*Session
	InitiateSession func(a ...interface{}) (ssn api.SessionAPI)
}

func NewSessions(mnssns *Sessions) (ssns *Sessions) {
	ssns = &Sessions{ssnscnt: 0, ssnslck: &sync.RWMutex{}, ssns: map[*Session]*Session{}, mnssns: mnssns}
	return
}

func (ssns *Sessions) CloseSession(ssn *Session) {
	closeSession(ssns, ssn)
}

func (ssns *Sessions) InvokeSession(a ...interface{}) (ssnapi api.SessionAPI) {
	if ssns != nil {
		ssnapi = InvokeSession(a...)
		if ssn, _ := ssnapi.(*Session); ssn != nil {
			appendSession(ssns, ssn)
		}
	}
	return
}

func closeSession(ssns *Sessions, ssn ...*Session) {
	if ssnsl := len(ssn); ssns != nil && ssnsl > 0 {
		var ssnclsecnt uint32 = 0
		var ssnrmv []bool = make([]bool, ssnsl)
		var exists = false
		func() {
			ssns.ssnslck.RLock()
			defer ssns.ssnslck.RUnlock()
			for ssnn := range ssn {
				sn := ssn[ssnn]
				exists = false
				_, exists = ssns.ssns[sn]
				if ssnrmv[ssnn] = exists; exists {
					ssnclsecnt++
				}
			}
		}()
		if ssnclsecnt > 0 {
			var didlck = false
			defer func() {
				if didlck {
					ssns.ssnslck.Unlock()
				}
			}()
			for sn := range ssnrmv {
				if ssnrmv[sn] {
					if !didlck {
						didlck = true
						ssns.ssnslck.Lock()
					}
					delete(ssns.ssns, ssn[sn])
				}
			}
			atomic.AddUint32(&ssns.ssnscnt, ^uint32(ssnclsecnt))
		}
	}
}

func appendSession(ssns *Sessions, ssn ...*Session) {
	if ssnsl := len(ssn); ssns != nil && ssnsl > 0 {
		var ssnappndcnt uint32 = 0
		var ssnrmv []bool = make([]bool, ssnsl)
		var exists = false
		func() {
			ssns.ssnslck.RLock()
			defer ssns.ssnslck.RUnlock()
			for ssnn := range ssn {
				sn := ssn[ssnn]
				exists = false
				_, exists = ssns.ssns[sn]
				if ssnrmv[ssnn] = exists; !exists {
					ssnappndcnt++
				}
			}
		}()
		if ssnappndcnt > 0 {
			var didlck = false
			defer func() {
				if didlck {
					ssns.ssnslck.Unlock()
				}
			}()
			for sn := range ssnrmv {
				if ssnrmv[sn] {
					if !didlck {
						didlck = true
						ssns.ssnslck.Lock()
					}
					sssn := ssn[sn]
					ssns.ssns[sssn] = sssn
				}
			}
			atomic.AddUint32(&ssns.ssnscnt, ssnappndcnt)
		}
	}
}

func (ssns *Sessions) Close() {
	if ssns != nil {
		if ssns.ssns != nil {
			for ssn := range ssns.ssns {
				closeSession(ssns, ssn)
			}
			ssns.ssns = nil
		}
		ssns = nil
	}
}
