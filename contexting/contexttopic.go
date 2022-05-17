package contexting

import (
	"sync"
	"sync/atomic"
	"time"
)

type ContextTopic struct {
	serial int64
}

func (ctxtpc *ContextTopic) Reader() (ctxtpcrdr *ContextTopicReader) {
	if ctxtpc != nil {
		ctxtpcrdr = newContextTopicReader(ctxtpc.serial)
	}
	return
}

func (ctxtcp *ContextTopic) Writer() (ctxtcpwtr *ContextTopicWriter) {
	if ctxtcp != nil {
		ctxtcpwtr = newContextTopicWriter(ctxtcp.serial)
	}
	return
}

func newContextTopic() (ctxtpc *ContextTopic) {
	ctxtpc = &ContextTopic{serial: nextTopicSerial()}
	func() {
		ctxtpcslck.Lock()
		defer ctxtpcslck.Unlock()
		ctxtpcs[ctxtpc.serial] = ctxtpc
	}()
	return
}

var lasttopicserial int64 = time.Now().UnixNano()

func nextTopicSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lasttopicserial, atomic.LoadInt64(&lasttopicserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

var lasttopicReaderserial int64 = time.Now().UnixNano()

func nextTopicReaderSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lasttopicReaderserial, atomic.LoadInt64(&lasttopicReaderserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

var lasttopicWriterserial int64 = time.Now().UnixNano()

func nextTopicWriterSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lasttopicWriterserial, atomic.LoadInt64(&lasttopicWriterserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

func newContextTopicReader(tpcserial int64) (ctxtpcrdr *ContextTopicReader) {
	ctxtpcrdr = &ContextTopicReader{serial: nextTopicReaderSerial()}
	func() {
		ctxtpcrdrslck.Lock()
		defer ctxtpcrdrslck.Unlock()
		ctxtpcrdrs[ctxtpcrdr.serial] = ctxtpcrdr
	}()
	return
}

func newContextTopicWriter(tpcserial int64) (ctxtpcwtr *ContextTopicWriter) {
	ctxtpcwtr = &ContextTopicWriter{serial: nextTopicWriterSerial()}
	func() {
		ctxtpcwtrslck.Lock()
		defer ctxtpcwtrslck.Unlock()
		ctxtpcwtrs[ctxtpcwtr.serial] = ctxtpcwtr
	}()
	return
}

type ContextTopicWriter struct {
	serial int64
}

type ContextTopicReader struct {
	serial int64
}

func (ctxtpcwtr *ContextTopicWriter) Print(a ...interface{}) (err error) {

	return
}

func (ctxtpcwtr *ContextTopicWriter) Write(p []byte) (n int, err error) {

	return
}

func (ctxtpcrdr *ContextTopicReader) Read(p []byte) (err error) {

	return
}

func (ctxtpcrdr *ContextTopicReader) Readln() (ln string, err error) {

	return
}

func (ctxtpcrdr *ContextTopicReader) Readlines() (lines []string, err error) {

	return
}

func (ctxtpcrdr *ContextTopicReader) ReadAll() (all string, err []error) {

	return
}

//Topic Reader instances
var ctxtpcrdrs map[int64]*ContextTopicReader = map[int64]*ContextTopicReader{}
var ctxtpcrdrslck *sync.RWMutex = &sync.RWMutex{}

//TopicWriter instances
var ctxtpcwtrs map[int64]*ContextTopicWriter = map[int64]*ContextTopicWriter{}
var ctxtpcwtrslck *sync.RWMutex = &sync.RWMutex{}

//Topic Instances
var ctxtpcs map[int64]*ContextTopic = map[int64]*ContextTopic{}
var ctxtpcslck *sync.RWMutex = &sync.RWMutex{}

var ctxtpcalias map[string]int64 = map[string]int64{}
var ctxtpcaliaslck *sync.RWMutex = &sync.RWMutex{}

func TopicExists(alias string) (exists bool) {
	if alias != "" {
		func() {
			ctxgrpaliaslck.RLock()
			defer ctxgrpaliaslck.RUnlock()
			_, exists = ctxgrpalias[alias]
		}()
	}
	return
}

func topicBySerial(serial int64) (ctxtpc *ContextTopic) {
	ctxtpcslck.RLock()
	defer ctxtpcslck.RUnlock()
	ctxtpc = ctxtpcs[serial]
	return
}

func topicSerialByAlias(alias string) (serial int64) {
	if TopicExists(alias) {
		func() {
			ctxtpcaliaslck.RLock()
			defer ctxtpcaliaslck.RUnlock()
			serial = ctxtpcalias[alias]
		}()
	}
	return
}

func TopicByAlias(alias string) (ctxtpc *ContextTopic) {
	if serial := topicSerialByAlias(alias); serial > 0 {
		ctxtpc = topicBySerial(serial)
	}
	return
}

func RegisterTopic(alias string) {
	if alias != "" {
		if !TopicExists(alias) {
			func() {
				var ctxtpc = newContextTopic()
				ctxtpcalias[alias] = ctxtpc.serial
			}()
		}
	}
}
