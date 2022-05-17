package contexting

import (
	"sync"
	"sync/atomic"
	"time"
)

type EventContextRead interface {
	Read(grpalias, topicalias string, a ...interface{})
}

type EventContextWrite interface {
	Print(grpalias, topicalias string, a ...interface{})
}

type Context struct {
	serial     int64
	EventRead  EventContextRead
	EventWrite EventContextWrite
}

func newContextReader(serial int64) (ctxrdr *ContextReader) {
	ctxrdr = &ContextReader{serial: nextContextReaderSerial()}
	func() {
		cntxtrdrsrlslck.Lock()
		defer cntxtrdrsrlslck.Unlock()
		cntxtrdrsrls[ctxrdr.serial] = ctxrdr
	}()
	return
}

func newContextWriter(serial int64) (ctxwtr *ContextWriter) {
	ctxwtr = &ContextWriter{serial: nextContextWriterSerial()}
	func() {
		cntxtwtrsrlslck.Lock()
		defer cntxtwtrsrlslck.Unlock()
		cntxtwtrsrls[ctxwtr.serial] = ctxwtr
	}()
	return
}

func (ctx *Context) Reader() (ctxrdr *ContextReader) {
	if ctx != nil && ctx.EventRead != nil {
		ctxrdr = newContextReader(ctx.serial)
		ctxrdr.evtCtxRd = ctx.EventRead
	}
	return
}

func (ctx *Context) Writer() (ctxwtr *ContextWriter) {
	if ctx != nil && ctx.EventWrite != nil {
		ctxwtr = newContextWriter(ctx.serial)
		ctxwtr.evtCtxWt = ctx.EventWrite
	}
	return
}

var lastContextserial int64 = time.Now().UnixNano()

func nextContextSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lastContextserial, atomic.LoadInt64(&lastContextserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

type ContextReader struct {
	serial   int64
	evtCtxRd EventContextRead
}

func (ctxrdr *ContextReader) Seek(offset int64, whence int) (n int64, err error) {

	return
}

func (ctxrdr *ContextReader) SetMaxRead(maxlen int64) (err error) {

	return
}

func (ctxrdr *ContextReader) Read(p []byte) (n int, err error) {

	return
}

func (ctxrdr *ContextReader) ReadRune() (r rune, size int, err error) {

	return
}

func (ctxrdr *ContextReader) Readln() (line string, err error) {

	return
}

func (ctxrdr *ContextReader) Readlines() (lines []string, err error) {

	return
}

func (ctxrdr *ContextReader) ReadAll() (all string, err error) {

	return
}

func (ctxrdr *ContextReader) Close() (err error) {

	return
}

type ContextWriter struct {
	serial   int64
	evtCtxWt EventContextWrite
}

func (ctxwtr *ContextWriter) Write(p []byte) (n int, err error) {

	return
}

func (ctxwtr *ContextWriter) Print(a ...interface{}) (err error) {

	return
}

func (ctxwtr *ContextWriter) Flush() (err error) {

	return
}

func (ctxwtr *ContextWriter) Close() (err error) {

	return
}

var lastContextReaderserial int64 = time.Now().UnixNano()

func nextContextReaderSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lastContextReaderserial, atomic.LoadInt64(&lastContextReaderserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

var lastContextWriterserial int64 = time.Now().UnixNano()

func nextContextWriterSerial() (nxsrl int64) {
	for {
		if nxsrl = time.Now().UnixNano(); atomic.CompareAndSwapInt64(&lastContextWriterserial, atomic.LoadInt64(&lastContextWriterserial), nxsrl) {
			break
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return
}

var cntxtsrls map[int64]*Context = map[int64]*Context{}
var cntxtsrlslck *sync.RWMutex = &sync.RWMutex{}

var cntxtrdrsrls map[int64]*ContextReader = map[int64]*ContextReader{}
var cntxtrdrsrlslck *sync.RWMutex = &sync.RWMutex{}

var cntxtwtrsrls map[int64]*ContextWriter = map[int64]*ContextWriter{}
var cntxtwtrsrlslck *sync.RWMutex = &sync.RWMutex{}

func NewContext() (cntxt *Context) {
	var serial = nextContextSerial()

	cntxt = &Context{serial: serial}

	func() {
		cntxtsrlslck.Lock()
		defer cntxtsrlslck.Unlock()
		cntxtsrls[serial] = cntxt
	}()
	return
}

func ContextBySerial(serial int64) (cntx *Context) {
	func() {
		cntxtsrlslck.RLock()
		defer cntxtsrlslck.RUnlock()
		cntx = cntxtsrls[serial]
	}()
	return
}
