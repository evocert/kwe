package caching

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
)

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
