package channeling

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/evocert/kwe/api"
)

type Sessions struct {
	ssnscnt         uint32
	ssnslck         *sync.RWMutex
	mnssns          *Sessions
	ssns            map[*Session]*Session
	InitiateSession func(a ...interface{}) (ssn api.SessionAPI)
	ssnsserial      int64
}

func NewSessions(mnssns *Sessions) (ssns *Sessions) {
	ssns = &Sessions{ssnscnt: 0, ssnslck: &sync.RWMutex{}, ssns: map[*Session]*Session{}, mnssns: mnssns, ssnsserial: nextserial()}
	return
}

var nxtsetiallck sync.RWMutex

func nextserial() (nxsrl int64) {
	func() {
		defer nxtsetiallck.RUnlock()
		nxtsetiallck.RLock()
		time.Sleep(1 * time.Nanosecond)
		nxsrl = time.Now().UnixNano()
	}()
	return
}

func (ssns *Sessions) Serial() int64 {
	return ssns.ssnsserial
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
			atomic.CompareAndSwapUint32(&ssns.ssnscnt, ssns.ssnscnt, ssns.ssnscnt-ssnclsecnt)
		}
	}
}

func appendSession(ssns *Sessions, ssn ...*Session) {
	if ssnsl := len(ssn); ssns != nil && ssnsl > 0 {
		var ssnappndcnt uint32 = 0
		var ssnappnd []bool = make([]bool, ssnsl)
		var exists = false
		func() {
			ssns.ssnslck.RLock()
			defer ssns.ssnslck.RUnlock()
			for ssnn := range ssn {
				sn := ssn[ssnn]
				exists = false
				_, exists = ssns.ssns[sn]
				if ssnappnd[ssnn] = !exists; !exists {
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
			for sn := range ssnappnd {
				if ssnappnd[sn] {
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

func (ssns *Sessions) Wait(milsecs ...int) {
	func() {
		var ssnlen = func() (ssnl int) {
			ssns.ssnslck.RLock()
			defer ssns.ssnslck.RUnlock()
			ssnl = len(ssns.ssns)
			return
		}
		for ssnlen() > 0 {
			if len(milsecs) > 0 && milsecs[0] > 5 {
				time.Sleep(time.Duration(milsecs[0]) * time.Millisecond)
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func (ssns *Sessions) Close() {
	if ssns != nil {

		ssns = nil
	}
}
