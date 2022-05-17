package contexting

import (
	"sync"
	"sync/atomic"
	"time"
)

var lastGroupserial int64 = time.Now().UnixNano()

func nextGroupSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lastGroupserial, atomic.LoadInt64(&lastGroupserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

type ContextGroup struct {
	serial int64
}

func newContexctGroup() (ctxgrp *ContextGroup) {
	ctxgrp = &ContextGroup{serial: nextGroupSerial()}
	func() {
		ctxgrpslck.Lock()
		defer ctxgrpslck.Unlock()
		ctxgrps[ctxgrp.serial] = ctxgrp
	}()
	return
}

func (cntxnggrp *ContextGroup) Print(topic string, a ...interface{}) (err error) {

	return
}

func (cntxnggrp *ContextGroup) Println(topic string, a ...interface{}) (err error) {

	return
}

var ctxgrps map[int64]*ContextGroup = map[int64]*ContextGroup{}
var ctxgrpslck *sync.RWMutex = &sync.RWMutex{}

func groupBySerial(serial int64) (ctxtpc *ContextTopic) {
	ctxtpcslck.RLock()
	defer ctxtpcslck.RUnlock()
	ctxtpc = ctxtpcs[serial]
	return
}

func groupSerialByAlias(alias string) (serial int64) {
	if GroupExists(alias) {
		func() {
			ctxgrpaliaslck.RLock()
			defer ctxgrpaliaslck.RUnlock()
			serial = ctxgrpalias[alias]
		}()
	}
	return
}

func GroupByAlias(alias string) (ctxgrp *ContextGroup) {
	if serial := groupSerialByAlias(alias); serial > 0 {
		func() {
			ctxgrpslck.RLock()
			defer ctxgrpslck.RUnlock()
			ctxgrp = ctxgrps[serial]
		}()
	}
	return
}

func GroupExists(alias string) (exists bool) {
	if alias != "" {
		func() {
			ctxgrpaliaslck.RLock()
			defer ctxgrpaliaslck.RUnlock()
			_, exists = ctxgrpalias[alias]
		}()
	}
	return
}

func RegisterGroup(alias string) {

}

var ctxgrpalias map[string]int64 = map[string]int64{}
var ctxgrpaliaslck *sync.RWMutex = &sync.RWMutex{}
