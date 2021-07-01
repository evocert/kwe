package caching

import (
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
)

type Map struct {
	lck      *sync.RWMutex
	lckedlvl int
	imp      map[interface{}]interface{}
	valid    bool
}

func (mp *Map) Lock() (lcked bool) {
	mp.lck.Lock()
	return
}

func (mp *Map) RLock() (lcked bool) {
	mp.lck.RLock()
	return
}

func (mp *Map) Unlock() (lcked bool) {
	mp.lck.Unlock()
	return
}

func (mp *Map) RUnlock() (lcked bool) {
	mp.lck.RUnlock()
	return
}

func (mp *Map) Keys(ks ...interface{}) (keys []interface{}) {
	if mp != nil {
		func() {
			if len(ks) == 0 {
				func() {
					mp.lck.RLock()
					defer mp.lck.RUnlock()
					keys = mapKeys(mp, nil)
				}()
			} else {
				var lkpmp *Map = nil
				func() {
					mp.lck.RLock()
					defer mp.lck.RUnlock()
					if lkv := mapFind(mp, nil, ks...); lkv != nil {
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
	return
}

func mapKeys(mp *Map, mphndlr *MapHandler) (keys []interface{}) {
	if mp != nil {
		if mp != nil && len(mp.imp) > 0 {
			keys := make([]interface{}, len(mp.imp))
			keysi := 0
			for k := range mp.imp {
				keys[keysi] = k
				keysi++
			}
		}
	}
	return
}

func (mp *Map) Values(ks ...interface{}) (values []interface{}) {
	if mp != nil {
		func() {
			if len(ks) == 0 {
				func() {
					mp.lck.RLock()
					defer mp.lck.RUnlock()
					values = mapValues(mp, nil)
				}()
			} else {
				var lkpmp *Map = nil
				func() {
					mp.lck.RLock()
					defer mp.lck.RUnlock()
					if lkv := mapFind(mp, nil, ks...); lkv != nil {
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
	return
}

func mapValues(mp *Map, mphndlr *MapHandler) (values []interface{}) {
	if mp != nil {
		if mp != nil && len(mp.imp) > 0 {
			values := make([]interface{}, len(mp.imp))
			valuesi := 0
			for _, v := range mp.imp {
				values[valuesi] = v
				valuesi++
			}
		}
	}
	return
}

func (mp *Map) Put(name interface{}, a ...interface{}) {
	if name != nil {
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
			if validMap(mp) {
				mp.Lock()
				defer mp.Unlock()
				mapPut(mp, nil, true, a...)
			}
		}()
	}
}

func mapPut(mp *Map, mphndlr *MapHandler, candispose bool, a ...interface{}) {
	if mp != nil && len(a) > 0 {
		for {
			if al := len(a); al > 0 {
				if al%2 == 0 {
					k := a[0]
					if k != nil {
						a = a[1:]
						v := a[0]
						func() {
							if v != nil {
								if vmp, vmpok := v.(map[interface{}]interface{}); vmpok {
									vnmp := NewMap()
									vnmp.Put(vmp)
									v = vnmp
								} else if vmp, vmpok := v.(map[string]interface{}); vmpok {
									vnmp := NewMap()
									vnmp.Put(vmp)
									v = vnmp
								}
							}
							if prvv, hask := mp.imp[k]; hask {
								delete(mp.imp, k)
								if candispose && prvv != nil {
									disposeMapVal(mp, mphndlr, prvv)
								}
							}
							mp.imp[k] = v
						}()
						a = a[1:]
					}
				} else {
					if m, mok := a[0].(map[string]interface{}); mok && len(m) > 0 {
						a = a[1:]
						for k, v := range m {
							a = append(a, k, v)
						}
					} else if mi, miok := a[0].(map[interface{}]interface{}); miok && len(mi) > 0 {
						a = a[1:]
						for k, v := range mi {
							a = append(a, k, v)
						}
					} else {
						a = a[1:]
					}
				}
			} else {
				break
			}
		}
	}
}

func disposeMapVal(mp *Map, mphndlr *MapHandler, valdispose interface{}) {
	if valdispose != nil {
		if arr, arrok := valdispose.([]interface{}); arrok && len(arr) > 0 {
			for _, d := range arr {
				disposeMapVal(mp, mphndlr, d)
			}
		} else if buf, bufok := valdispose.(*iorw.Buffer); bufok && buf != nil {
			buf.Close()
			buf = nil
		} else if mpvk, mpvkok := valdispose.(*Map); mpvkok && mpvk != nil {
			mpvk.Close(mphndlr)
			mpvk = nil
		}
	}
}

func (mp *Map) Close(ks ...interface{}) {
	var mphndlr *MapHandler = nil
	if len(ks) >= 1 {
		mphndlr, _ = ks[0].(*MapHandler)
		ks = ks[1:]
	}
	if validMap(mp) {
		if len(ks) == 0 {
			func() {
				mp.Lock()
				defer mp.Unlock()
				if len(mp.imp) > 0 {
					ks := make([]interface{}, len(mp.imp))
					ksi := 0
					for k := range mp.imp {
						ks[ksi] = k
						ksi++
					}
					mapRemove(mp, mphndlr, ks...)
				}
				mp.imp = nil
			}()
			mp.valid = false
			mp = nil
		} else {
			var lkpmp *Map = nil
			var prvlkmp *Map = mp
			if len(ks) > 1 {
				if v := mp.Find(ks[:len(ks)-1]); v != nil {
					prvlkmp, _ = v.(*Map)
					v = nil
				}
			}
			if v := mp.Find(ks...); v != nil {
				lkpmp, _ = v.(*Map)
				v = nil
			}
			if lkpmp != nil && prvlkmp != nil {
				func() {
					prvlkmp.Lock()
					defer prvlkmp.Unlock()
					mapRemove(prvlkmp, mphndlr, ks[len(ks)-1])
				}()
				lkpmp = nil
				prvlkmp = nil
			}
		}
	}
}

func validMap(mp *Map) (valid bool) {
	if mp != nil {
		func() {
			mp.RLock()
			defer mp.RUnlock()
			valid = mp.valid
		}()
	}
	return
}

func (mp *Map) Reset(ks ...interface{}) {
}

func (mp *Map) Clear(ks ...interface{}) {
	var mphndlr *MapHandler = nil
	if len(ks) >= 1 {
		mphndlr, _ = ks[0].(*MapHandler)
		ks = ks[1:]
	}
	if validMap(mp) {
		if len(ks) == 0 {
			func() {
				mp.Lock()
				defer mp.Unlock()
				if len(mp.imp) > 0 {
					ks := make([]interface{}, len(mp.imp))
					ksi := 0
					for k := range mp.imp {
						ks[ksi] = k
						ksi++
					}
					mapRemove(mp, mphndlr, ks...)
				}
			}()
		} else {
			var lkpmp *Map = nil
			if v := mp.Find(ks...); v != nil {
				lkpmp, _ = v.(*Map)
				v = nil
			}
			if lkpmp != nil {
				lkpmp.Clear(mphndlr)
			}
		}
	}
}

func (mp *Map) Remove(name ...interface{}) {
	if validMap(mp) && len(name) > 0 {
		func() {
			if validMap(mp) {
				mp.Lock()
				defer mp.Unlock()
				mapRemove(mp, nil, name...)
			}
		}()
	}
}

func mapRemove(mp *Map, mphndlr *MapHandler, name ...interface{}) {
	if len(name) > 0 {
		if len(mp.imp) > 0 {
			for _, nme := range name {
				if nmv, nmok := mp.imp[nme]; nmok {
					if nmv != nil {
						disposeMapVal(mp, mphndlr, nmv)
					}
					delete(mp.imp, nme)
				}
			}
		}
	}
}

func (mp *Map) Find(ks ...interface{}) (value interface{}) {
	if len(ks) > 0 {
		func() {
			mp.RLock()
			defer mp.RUnlock()
			value = mapFind(mp, nil, ks...)
		}()
	}
	return
}

func mapFind(mp *Map, mphndlr *MapHandler, ks ...interface{}) (vfound interface{}) {
	if mp != nil {
		var lkpmp *Map = mp
		if ksl := len(ks); ksl > 0 {
			ksi := 0
			var subfind = func() (valf interface{}, found bool) {
				if lkpmp != mp {
					crntmp := lkpmp
					crntmp.RLock()
					defer crntmp.RUnlock()
				}
				ksi++
				if valf, found = lkpmp.imp[ks[ksi-1]]; found {
					if ksi < ksl {
						found = false
						if vmp, vmpok := valf.(*Map); vmpok && vmp != nil {
							lkpmp = vmp
						} else {
							ksi = ksl
						}
					}
				}
				return
			}
			for ksi < ksl {
				if valf, found := subfind(); found {
					vfound = valf
					break
				}
			}
		}
	}
	return
}

func (mp *Map) Push(k interface{}, a ...interface{}) (length int) {
	length = -1
	if validMap(mp) {
		func() {
			if len(a) == 0 {
				if a != nil {
					a = append([]interface{}{k}, a)
					length = mapPush(mp, nil, false, a...)
				}
			} else {
				a = append([]interface{}{k}, a...)
				length = mapPush(mp, nil, false, a...)
			}
		}()
	}
	return
}

func mapPush(mp *Map, mphndlr *MapHandler, focusCurrentMap bool, a ...interface{}) (length int) {
	length = -1
	if al := len(a); al >= 2 {
		ks := a[:al-1]
		a := a[al-1]
		if validMap(mp) {
			func() {
				mp.RLock()
				if v := mapFind(mp, mphndlr, ks...); v != nil {
					mp.RUnlock()
					if arr, arrok := v.([]interface{}); arrok {
						func() {
							mp.Lock()
							defer mp.Unlock()
							if varr, varrok := a.([]interface{}); varrok && len(varr) > 0 {
								arr = append(arr, varr)
							} else {
								var lkpmp *Map = nil
								if a != nil {
									if vmp, _ := a.(map[interface{}]interface{}); vmp != nil {
										lkpmp = NewMap()
										lkpmp.Put(vmp)
										a = lkpmp
									} else if vmp, _ := a.(map[string]interface{}); vmp != nil {
										lkpmp = NewMap()
										lkpmp.Put(vmp)
										a = lkpmp
									}
								}
								arr = append(arr, a)
								if lkpmp != nil && focusCurrentMap && mphndlr != nil {
									mphndlr.crntmp = lkpmp
								}
							}
							v = arr
							mapPut(mp, mphndlr, false, append(ks, v)...)
							length = len(arr)
						}()
					}
				} else {
					mp.RUnlock()
				}
			}()
		}
	}
	return length
}

func (mp *Map) Shift(k interface{}, a ...interface{}) (length int) {
	length = -1
	if validMap(mp) {
		func() {
			if len(a) == 0 {
				if a != nil {
					a = append([]interface{}{k}, a)
					length = mapShift(mp, nil, a...)
				}
			} else {
				a = append([]interface{}{k}, a...)
				length = mapShift(mp, nil, a...)
			}
		}()
	}
	return
}

func mapShift(mp *Map, mphndlr *MapHandler, a ...interface{}) (length int) {
	length = -1
	if al := len(a); al >= 2 {
		ks := a[:al-1]
		a := a[al-1]
		if validMap(mp) {
			func() {
				mp.RLock()
				if v := mapFind(mp, mphndlr, ks...); v != nil {
					mp.RUnlock()
					if arr, arrok := v.([]interface{}); arrok {
						func() {
							mp.Lock()
							defer mp.Unlock()
							if varr, varrok := a.([]interface{}); varrok && len(varr) > 0 {
								arr = append([]interface{}{varr}, arr...)
							} else {
								var lkpmp *Map = nil
								if a != nil {
									if vmp, _ := a.(map[interface{}]interface{}); vmp != nil {
										lkpmp = NewMap()
										lkpmp.Put(vmp)
										a = lkpmp
									} else if vmp, _ := a.(map[string]interface{}); vmp != nil {
										lkpmp = NewMap()
										lkpmp.Put(vmp)
										a = lkpmp
									}
								}
								arr = append([]interface{}{a}, arr...)
							}
							v = arr
							mapPut(mp, mphndlr, false, append(ks, v)...)
							length = len(arr)
						}()
					}
				} else {
					mp.RUnlock()
				}
			}()
		}
	}
	return length
}

func (mp *Map) Pop(k interface{}, a ...interface{}) (pop interface{}) {
	if validMap(mp) {
		func() {
			mp.Lock()
			mp.Unlock()
			a = append([]interface{}{k}, a...)
			pop = mapPop(mp, nil, a...)
		}()
	}
	return
}

func mapPop(mp *Map, mphndlr *MapHandler, a ...interface{}) (pop interface{}) {
	if v := mapFind(mp, mphndlr, a...); v != nil {
		if arr, arrok := v.([]interface{}); arrok && len(arr) > 0 {
			pop = arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			if arr == nil {
				arr = []interface{}{}
			}
			if len(a) > 1 {
				if lkpmp, _ := mapFind(mp, mphndlr, a[:len(a)-1]...).(*Map); lkpmp != nil {
					mapPut(lkpmp, mphndlr, false, append(a[:len(a)-1], arr)...)
				}
			} else {
				mapPut(mp, mphndlr, false, append(a[:len(a)-1], arr)...)
			}
		}
	}
	return
}

func (mp *Map) Unshift(k interface{}, a ...interface{}) (unshift interface{}) {
	if validMap(mp) {
		func() {
			mp.Lock()
			mp.Unlock()
			a = append([]interface{}{k}, a...)
			unshift = mapUnshift(mp, nil, a...)
		}()
	}
	return
}

func mapUnshift(mp *Map, mphndlr *MapHandler, a ...interface{}) (unshift interface{}) {
	if v := mapFind(mp, mphndlr, a...); v != nil {
		if arr, arrok := v.([]interface{}); arrok && len(arr) > 0 {
			unshift = arr[0]
			arr = arr[1:]
			if arr == nil {
				arr = []interface{}{}
			}
			if len(a) > 1 {
				if lkpmp, _ := mapFind(mp, mphndlr, a[:len(a)-1]...).(*Map); lkpmp != nil {
					mapPut(lkpmp, mphndlr, false, append(a[:len(a)-1], arr)...)
				}
			} else {
				mapPut(mp, mphndlr, false, append(a[:len(a)-1], arr)...)
			}
		}
	}
	return
}

func (mp *Map) At(k interface{}, a ...interface{}) (arv interface{}) {
	if validMap(mp) {
		func() {
			mp.RLock()
			defer mp.RUnlock()
			arv = mapAt(mp, nil, a...)
		}()
	}
	return
}

func mapAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (arv interface{}) {

	return
}

func (mp *Map) FocusAt(k interface{}, a ...interface{}) (focused bool) {
	return
}

func mapFocusAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (focused bool) {
	if mphndlr != nil && validMap(mp) {
		if av := mapAt(mp, mphndlr, a...); av != nil {
			if lkpmp, _ := av.(*Map); lkpmp != nil {
				if mphndlr.mp == lkpmp {
					mphndlr.crntmp = nil
				} else {
					mphndlr.crntmp = lkpmp
				}
			}
		}
	}
	return
}

func (mp *Map) Fprint(w io.Writer, ks ...interface{}) (err error) {
	if w != nil {
		if validMap(mp) {
			if len(ks) == 0 {
				func() {
					mp.RLock()
					defer mp.RUnlock()
					err = mapFPrint(mp, nil, w)
				}()
			} else {
				var lkpmv *Map = nil
				if v := mp.Find(ks...); v != nil {
					lkpmv, _ = v.(*Map)
					v = nil
				}
				if lkpmv != nil {
					func() {
						lkpmv.RLock()
						defer lkpmv.RUnlock()
						err = mapFPrint(lkpmv, nil, w)
					}()
				} else {
					iorw.Fprint(w, "{}")
				}
			}
		} else {
			iorw.Fprint(w, "null")
		}
	}
	return
}

func mapFPrint(mp *Map, mphndlr *MapHandler, w io.Writer) (err error) {
	if w != nil {
		if mp != nil {
			iorw.Fprint(w, "{")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "")
			enc.SetEscapeHTML(false)
			lnimp := len(mp.imp)
			for k, v := range mp.imp {
				enc.Encode(k)
				iorw.Fprint(w, ":")
				writeMapVal(w, enc, mphndlr, v)
				lnimp--
				if lnimp > 0 {
					iorw.Fprint(w, ",")
				}
			}
			iorw.Fprint(w, "}")
		}
	}
	return
}

func writeMapVal(w io.Writer, enc *json.Encoder, mphndlr *MapHandler, v interface{}) (err error) {
	if v == nil {
		iorw.Fprint(w, "null")
	} else {
		if arr, arrok := v.([]interface{}); arrok {
			iorw.Fprint(w, "[")
			if al := len(arr); al > 0 {
				for an, a := range arr {
					writeMapVal(w, enc, mphndlr, a)
					if an < al-1 {
						iorw.Fprint(w, ",")
					}
				}
			}
			iorw.Fprint(w, "]")
		} else if vmp, vmpok := v.(*Map); vmpok {
			func() {
				if validMap(vmp) {
					vmp.RLock()
					defer vmp.RUnlock()
					err = mapFPrint(vmp, mphndlr, w)
				} else {
					iorw.Fprint(w, "{}")
				}
			}()
		} else if buf, bufok := v.(*iorw.Buffer); bufok {
			enc.Encode(buf.String())
		} else {
			enc.Encode(v)
		}
	}
	return
}

func (mp *Map) String(ks ...interface{}) (s string) {
	mprdr := MapReader(mp, ks...)
	s, _ = mprdr.ReadAll()
	mprdr.Close()
	mprdr = nil
	if s != "" {
		s = strings.Replace(s, "\n", " ", -1)
	}
	return
}

func (mp *Map) Focus(ks ...interface{}) (focused bool) {
	return
}

func (mp *Map) Handler() (mphndlr *MapHandler) {
	if validMap(mp) {
		mphndlr = NewMapHandler(mp)
	}
	return
}

func NewMap() (mp *Map) {
	mp = &Map{lck: &sync.RWMutex{}, lckedlvl: 0, valid: true, imp: map[interface{}]interface{}{}}
	return
}

func GLOBALMAP() *Map {
	return glbmp
}

func GLOBALMAPHANDLER() *MapHandler {
	return glbmphndlr
}

var glbmphndlr *MapHandler = nil

var glbmp *Map = nil

func init() {
	glbmp = NewMap()
	glbmphndlr = NewMapHandler(glbmp)
}
