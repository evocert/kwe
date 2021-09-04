package caching

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

//ActiveHandler - struct
type ActiveHandler struct {
	mphndlr *MapHandler
	rntme   active.Runtime
	prms    parameters.ParametersAPI
}

func (atvmphndlr *ActiveHandler) Keys(ks ...interface{}) (keys []interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		keys = atvmphndlr.mphndlr.Keys(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Values(ks ...interface{}) (values []interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		values = atvmphndlr.mphndlr.Values(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) IsMap(ks ...interface{}) (ismap bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		ismap = atvmphndlr.mphndlr.IsMap(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Exists(ks ...interface{}) (exists bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		exists = atvmphndlr.mphndlr.Exists(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Find(ks ...interface{}) (val interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		val = atvmphndlr.mphndlr.Find(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Put(ks interface{}, a ...interface{}) (putit bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		putit = atvmphndlr.mphndlr.Put(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Remove(ks ...interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		atvmphndlr.mphndlr.Remove(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Fprint(w io.Writer, ks ...interface{}) (err error) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		err = atvmphndlr.mphndlr.Fprint(w, ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) String(ks ...interface{}) (s string) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		s = atvmphndlr.mphndlr.String(ks...)
	} else {
		s = "null"
	}
	return
}

func (atvmphndlr *ActiveHandler) Focus(ks ...interface{}) (focused bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		focused = atvmphndlr.mphndlr.Focus(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Reset(ks ...interface{}) (isreset bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		isreset = atvmphndlr.mphndlr.Reset(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Clear(ks ...interface{}) (cleared bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		cleared = atvmphndlr.mphndlr.Clear(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Close(ks ...interface{}) (closed bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		closed = atvmphndlr.mphndlr.Close(ks...)
		if len(ks) == 0 {
			if atvmphndlr != nil {
				func() {
					if atvmphndlr.mphndlr != nil {
						atvmphndlr.mphndlr = nil
					}
					if atvmphndlr.rntme != nil {
						atvmphndlr.rntme = nil
					}
				}()
				atvmphndlr = nil
			}
		}
	}
	return
}

//Array
func (atvmphndlr *ActiveHandler) IsMapAt(ks ...interface{}) (ismap bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		ismap = atvmphndlr.mphndlr.IsMap(ks...)
	}
	return
}

func (atvmphndlr *ActiveHandler) ExistsAt(ks interface{}, a ...interface{}) (existsast bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		existsast = atvmphndlr.mphndlr.ExistsAt(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Push(ks interface{}, a ...interface{}) (size int) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		size = atvmphndlr.mphndlr.Push(ks, a...)
	} else {
		size = -1
	}
	return
}

func (atvmphndlr *ActiveHandler) Pop(ks interface{}, a ...interface{}) (pop interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		pop = atvmphndlr.mphndlr.Pop(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) Shift(ks interface{}, a ...interface{}) (size int) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		size = atvmphndlr.mphndlr.Shift(ks, a...)
	} else {
		size = -1
	}
	return
}

func (atvmphndlr *ActiveHandler) Unshift(ks interface{}, a ...interface{}) (unshift interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		unshift = atvmphndlr.mphndlr.Unshift(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) At(ks interface{}, a ...interface{}) (val interface{}) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		val = atvmphndlr.mphndlr.At(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) FocusAt(ks interface{}, a ...interface{}) (focused bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		focused = atvmphndlr.mphndlr.FocusAt(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) ClearAt(ks interface{}, a ...interface{}) (clearedat bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		clearedat = atvmphndlr.mphndlr.ClearAt(ks, a...)
	}
	return
}

func (atvmphndlr *ActiveHandler) CloseAt(ks interface{}, a ...interface{}) (closedat bool) {
	if atvmphndlr != nil && atvmphndlr.mphndlr != nil && atvmphndlr.rntme != nil {
		closedat = atvmphndlr.mphndlr.CloseAt(ks, a...)
	}
	return
}

//ActiveHandler return ActiveHandler for active.Runtime
func (mphndlr *MapHandler) ActiveHandler(rntme active.Runtime, prms ...parameters.ParametersAPI) *ActiveHandler {
	return newActiveHandler(mphndlr, rntme, prms...)
}

func newActiveHandler(mphndlr *MapHandler, rntme active.Runtime, prms ...parameters.ParametersAPI) (atvmphnldr *ActiveHandler) {
	if mphndlr != nil && rntme != nil {
		atvmphnldr = &ActiveHandler{mphndlr: mphndlr, rntme: rntme}
		if len(prms) > 0 && prms[0] != nil {
			atvmphnldr.prms = prms[0]
		}
	}
	return
}

type MapHandler struct {
	mp       *Map
	crntmp   *Map
	internal bool
}

var handlerPool *sync.Pool = nil

func initMapHandler() (mphndlr *MapHandler) {
	mphndlr = &MapHandler{}
	return
}

func newHandler(mp *Map, internal bool) (mphndlr *MapHandler) {
	if v := handlerPool.Get(); v != nil {
		mphndlr = v.(*MapHandler)
	} else {
		mphndlr = initMapHandler()
	}
	mphndlr.mp = mp
	mphndlr.internal = internal
	return
}

func clearHandler(mphndlr *MapHandler) {
	if mphndlr != nil {
		mphndlr.mp = nil
		mphndlr.crntmp = nil
	}
}

func putHandler(mphndlr *MapHandler) {
	clearHandler(mphndlr)
	handlerPool.Put(mphndlr)
}

func NewMapHandler(a ...interface{}) (mphndlr *MapHandler) {
	var mp *Map = nil
	var internal bool = false
	if len(a) >= 1 {
		if mp, _ = a[0].(*Map); mp != nil {
			a = a[1:]
		}
		if mp == nil {
			mp = NewMap()
			internal = true
		}
	}
	if mp != nil {
		if len(a) > 0 {
			mp.Put(a[0], a[1:])
		}
		mphndlr = initMapHandler()
		mphndlr.mp = mp
		mphndlr.internal = internal
	}
	return
}

func (mphndlr *MapHandler) currentMap() (crntmp *Map) {
	if mphndlr != nil {
		if crntmp = mphndlr.crntmp; crntmp == nil {
			if validMap(mphndlr.mp) {
				crntmp = mphndlr.mp
			}
		} else if !validMap(crntmp) {
			crntmp = nil
		}
	}
	return
}

func (mphndlr *MapHandler) Reset(ks ...interface{}) (reset bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) == 0 {
				if mphndlr.mp != crntmp {
					mphndlr.crntmp = nil
				}
				reset = true
			} else {
				reset = mphndlr.Focus(ks...)
			}
		}
	}
	return
}

//NewBuffer helper that returns instance of *iorw.Buffer
func (mphndlr *MapHandler) NewBuffer() (buf *iorw.Buffer) {
	buf = iorw.NewBuffer()
	return
}

func (mphndlr *MapHandler) Focus(ks ...interface{}) (focused bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) > 0 {
				func() {
					crntmp.RLock()
					if value := mapFind(crntmp, mphndlr, ks...); value != nil {
						crntmp.RUnlock()
						if lkpmp, _ := value.(*Map); lkpmp != nil && validMap(lkpmp) {
							if crntmp != lkpmp {
								if mphndlr.mp != lkpmp {
									mphndlr.crntmp = lkpmp
								} else {
									mphndlr.crntmp = nil
								}
							}
							focused = true
						}
					} else {
						crntmp.RUnlock()
					}
				}()
			}
		}
	}
	return
}

func (mphndlr *MapHandler) IsMap(ks ...interface{}) (ismap bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			ismap = mapIsMap(crntmp, mphndlr, ks...)
		}
	}
	return
}

func (mphndlr *MapHandler) Exists(ks ...interface{}) (exist bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			exist = mapExists(crntmp, mphndlr, ks...)
		}
	}
	return
}

func (mphndlr *MapHandler) Find(ks ...interface{}) (value interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			value = mapFind(crntmp, mphndlr, ks...)
		}
	}
	return
}

func (mphndlr *MapHandler) Keys(ks ...interface{}) (keys []interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			keys = mapKeys(crntmp, mphndlr, ks...)
		}
	}
	return
}

func (mphndlr *MapHandler) Values(ks ...interface{}) (values []interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			values = mapValues(crntmp, mphndlr, ks...)
		}
	}
	return
}

func (mphndlr *MapHandler) Put(name interface{}, a ...interface{}) (putit bool) {
	if name != nil {
		if mphndlr != nil {
			if crntmp := mphndlr.currentMap(); crntmp != nil {
				if len(a) == 0 {
					if mpkv, mpkvok := name.(map[interface{}]interface{}); mpkvok && len(mpkv) > 0 {
						a = []interface{}{mpkv}
					} else {
						if mpkv, mpkvok := name.(map[string]interface{}); mpkvok && len(mpkv) > 0 {
							a = []interface{}{mpkv}
						}
					}
				} else {
					a = append([]interface{}{name}, a...)
				}
				putit = mapPut(crntmp, mphndlr, true, a...)
			}
		}
	}
	return
}

func (mphndlr *MapHandler) Pop(k interface{}, a ...interface{}) (pop interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				crntmp.Lock()
				crntmp.Unlock()
				a = append([]interface{}{k}, a...)
				pop = mapPop(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Unshift(k interface{}, a ...interface{}) (unshift interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				crntmp.Lock()
				crntmp.Unlock()
				a = append([]interface{}{k}, a...)
				unshift = mapUnshift(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) IsMapAt(k interface{}, a ...interface{}) (ismap bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					if a != nil {
						a = append([]interface{}{k}, a)
						ismap = mapIsMapAt(crntmp, mphndlr, a...)
					}
				} else {
					a = append([]interface{}{k}, a...)
					ismap = mapIsMapAt(crntmp, mphndlr, a...)
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) ExistsAt(k interface{}, a ...interface{}) (exists bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					a = append([]interface{}{k}, a)
				} else {
					a = append([]interface{}{k}, a...)
				}
				exists = mapExistsAt(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) ClearAt(k interface{}, a ...interface{}) (cleared bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(a) == 0 {
				a = append([]interface{}{k}, a)
			} else {
				a = append([]interface{}{k}, a...)
			}
			cleared = mapClearAt(crntmp, mphndlr, a...)
		}
	}
	return
}

func (mphndlr *MapHandler) CloseAt(k interface{}, a ...interface{}) (closed bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					a = append([]interface{}{k}, a)
				} else {
					a = append([]interface{}{k}, a...)
				}
				closed = mapCloseAt(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) At(k interface{}, a ...interface{}) (arv interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					a = append([]interface{}{k}, a)
				} else {
					a = append([]interface{}{k}, a...)
				}
				arv = mapAt(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) FocusAt(k interface{}, a ...interface{}) (focused bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					a = []interface{}{k}
				} else {
					a = append([]interface{}{k}, a...)
				}
				focused = mapFocusAt(crntmp, mphndlr, a...)
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Clear(ks ...interface{}) (cleared bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) == 0 {
				func() {
					crntmp.Lock()
					defer crntmp.Unlock()
					if len(crntmp.imp) > 0 {
						ks := make([]interface{}, len(crntmp.imp))
						ksi := 0
						for k := range crntmp.imp {
							ks[ksi] = k
							ksi++
						}
						mapRemove(crntmp, mphndlr, ks...)
					}
				}()
				cleared = true
			} else {
				var lkpmp *Map = nil
				func() {
					crntmp.RLock()
					defer crntmp.RUnlock()
					if v := mapFind(crntmp, mphndlr, ks...); v != nil {
						lkpmp, _ = v.(*Map)
						v = nil
					}
				}()
				if lkpmp != nil {
					lkpmp.Clear(mphndlr)
					cleared = true
				}
			}
		}
	}
	return
}

func (mphndlr *MapHandler) Remove(name ...interface{}) {
	if len(name) > 0 {
		if mphndlr != nil {
			if crntmp := mphndlr.currentMap(); crntmp != nil {
				func() {
					crntmp.Lock()
					defer crntmp.Unlock()
					mapRemove(crntmp, mphndlr, name...)
				}()
			}
		}
	}
}

func (mphndlr *MapHandler) Fprint(w io.Writer, ks ...interface{}) (err error) {
	if w != nil {
		if mphndlr != nil {
			if crntmp := mphndlr.currentMap(); validMap(crntmp) {
				if len(ks) == 0 {
					func() {
						crntmp.RLock()
						defer crntmp.RUnlock()
						err = mapFPrint(crntmp, mphndlr, w)
					}()
				} else {
					var lkpmv *Map = nil
					func() {
						crntmp.RLock()
						defer crntmp.RUnlock()
						if v := mapFind(crntmp, mphndlr, ks...); v != nil {
							lkpmv, _ = v.(*Map)
							v = nil
						}
					}()
					if lkpmv != nil {
						func() {
							lkpmv.RLock()
							defer lkpmv.RUnlock()
							err = mapFPrint(lkpmv, mphndlr, w)
						}()
					} else {
						iorw.Fprint(w, "{}")
					}
				}
			} else {
				iorw.Fprint(w, "null")
			}
		} else {
			iorw.Fprint(w, "null")
		}
	}
	return
}

func (mphndlr *MapHandler) Shift(k interface{}, a ...interface{}) (length int) {
	length = -1
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					if a != nil {
						a = append([]interface{}{k}, a)
						length = mapShift(crntmp, mphndlr, a...)
					}
				} else {
					a = append([]interface{}{k}, a...)
					length = mapShift(crntmp, mphndlr, a...)
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Push(k interface{}, a ...interface{}) (length int) {
	length = -1
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(a) == 0 {
					if a != nil {
						a = append([]interface{}{k}, a)
						length = mapPush(crntmp, mphndlr, true, a...)
					}
				} else {
					a = append([]interface{}{k}, a...)
					length = mapPush(crntmp, mphndlr, true, a...)
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) String(ks ...interface{}) (s string) {
	if mphndlr != nil {
		pi, pw := io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func() {
			defer func() {
				pw.Close()
			}()
			ctxcancel()
			mphndlr.Fprint(pw, ks...)
		}()
		<-ctx.Done()
		s, _ = iorw.ReaderToString(pi)
		if s != "" {
			s = strings.Replace(s, "\n", " ", -1)
		}
	}
	return
}

func (mphndlr *MapHandler) Close(ks ...interface{}) (closed bool) {
	if mphndlr != nil {
		kl := len(ks)
		if kl == 0 {
			if crntmp := mphndlr.mp; validMap(crntmp) {
				if mphndlr.internal {
					crntmp.Close(mphndlr)
				}
				clearHandler(mphndlr)
				crntmp = nil
				mphndlr = nil
			}
			closed = true
		} else {
			if crntmp := mphndlr.currentMap(); validMap(crntmp) {
				ks = append([]interface{}{mphndlr}, ks...)
				crntmp.Close(ks...)
			}
			closed = true
		}
	}
	return
}

func (mphndlr *MapHandler) Reader(ks ...interface{}) io.Reader {
	return MapReader(mphndlr, ks...)
}

/*
Remove(...interface{}) interface{}
Find(...interface{}) interface{}
FPrint(io.Writer, ...interface{}) error
String(...interface{}) string
Clear()
Close()
*/
