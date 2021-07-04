package caching

import (
	"context"
	"io"
	"runtime"
	"strings"

	"github.com/evocert/kwe/iorw"
)

type MapHandler struct {
	mp       *Map
	crntmp   *Map
	internal bool
}

func finalizeMapHandler(mphndlr *MapHandler) {
	runtime.SetFinalizer(mphndlr, nil)
	mphndlr.Close()
	mphndlr = nil
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
		mphndlr = &MapHandler{mp: mp, internal: internal}
		//runtime.SetFinalizer(mphndlr, finalizeMapHandler)
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

func (mphndlr *MapHandler) Reset(ks ...interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) == 0 {
				if mphndlr.mp != crntmp {
					mphndlr.crntmp = nil
				}
			} else {
				mphndlr.Focus(ks...)
			}
		}
	}
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

func (mphndlr *MapHandler) Exist(ks ...interface{}) (exist bool) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) > 0 {
				func() {
					crntmp.RLock()
					defer crntmp.RUnlock()
					exist = mapExist(crntmp, mphndlr, ks...)
				}()
			}
		}
	}
	return
}

func (mphndlr *MapHandler) Find(ks ...interface{}) (value interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			if len(ks) > 0 {
				func() {
					crntmp.RLock()
					defer crntmp.RUnlock()
					value = mapFind(crntmp, mphndlr, ks...)
				}()
			}
		}
	}
	return
}

func (mphndlr *MapHandler) Keys(ks ...interface{}) (keys []interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(ks) == 0 {
					func() {
						crntmp.lck.RLock()
						defer crntmp.lck.RUnlock()
						keys = mapKeys(crntmp, mphndlr)
					}()
				} else {
					var lkpmp *Map = nil
					func() {
						crntmp.lck.RLock()
						defer crntmp.lck.RUnlock()
						if lkv := mapFind(crntmp, mphndlr, ks...); lkv != nil {
							lkpmp, _ = lkv.(*Map)
						}
					}()
					if lkpmp != nil {
						func() {
							lkpmp.lck.RLock()
							defer lkpmp.lck.RUnlock()
							keys = mapKeys(lkpmp, nil)
						}()
					}
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Values(ks ...interface{}) (values []interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				if len(ks) == 0 {
					func() {
						crntmp.lck.RLock()
						defer crntmp.lck.RUnlock()
						values = mapValues(crntmp, mphndlr)
					}()
				} else {
					var lkpmp *Map = nil
					func() {
						crntmp.lck.RLock()
						defer crntmp.lck.RUnlock()
						if lkv := mapFind(crntmp, mphndlr, ks...); lkv != nil {
							lkpmp, _ = lkv.(*Map)
						}
					}()
					if lkpmp != nil {
						func() {
							lkpmp.lck.RLock()
							defer lkpmp.lck.RUnlock()
							values = mapValues(lkpmp, nil)
						}()
					}
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Put(name interface{}, a ...interface{}) {
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
				func() {
					if validMap(crntmp) {
						crntmp.Lock()
						defer crntmp.Unlock()
						mapPut(crntmp, mphndlr, true, a...)
					}
				}()
			}
		}
	}
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

func (mphndlr *MapHandler) At(k interface{}, a ...interface{}) (arv interface{}) {
	if mphndlr != nil {
		if crntmp := mphndlr.currentMap(); crntmp != nil {
			func() {
				crntmp.RLock()
				defer crntmp.RUnlock()
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
				crntmp.RLock()
				defer crntmp.RUnlock()
				if len(a) == 0 {
					if a != nil {
						a = append([]interface{}{k}, a)
						focused = mapFocusAt(crntmp, mphndlr, a...)
					}
				} else {
					a = append([]interface{}{k}, a...)
					focused = mapFocusAt(crntmp, mphndlr, a...)
				}
			}()
		}
	}
	return
}

func (mphndlr *MapHandler) Clear(ks ...interface{}) {
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
				}
			}
		}
	}
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

func (mphndlr *MapHandler) Close(ks ...interface{}) {
	if mphndlr != nil {
		kl := len(ks)
		if kl == 0 {
			if crntmp := mphndlr.mp; validMap(crntmp) {
				crntmp.Close(mphndlr)
				mphndlr.mp = nil
				mphndlr.crntmp = nil
				crntmp = nil
				mphndlr = nil
			}
		} else {
			if crntmp := mphndlr.currentMap(); validMap(crntmp) {
				ks = append([]interface{}{mphndlr}, ks...)
				crntmp.Close(ks...)
			}
		}
	}
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
