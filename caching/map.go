package caching

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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
		keys = mapKeys(mp, nil, ks...)
	}
	return
}

func keysFound(ksmp *Map) (keys []interface{}) {
	if validMap(ksmp) {
		ksmp.RLock()
		defer ksmp.RUnlock()
		if ksmp != nil && len(ksmp.imp) > 0 {
			keys = make([]interface{}, len(ksmp.imp))
			keysi := 0
			for _, v := range ksmp.imp {
				keys[keysi] = v
				keysi++
			}
		}
	}
	return
}

func mapKeys(mp *Map, mphndlr *MapHandler, ks ...interface{}) (keys []interface{}) {
	if len(ks) == 0 {
		keys = keysFound(mp)
	} else {
		var lkpmp *Map = nil
		func() {
			if lkv := mapFind(mp, mphndlr, ks...); lkv != nil {
				lkpmp, _ = lkv.(*Map)
			}
		}()
		if lkpmp != nil {
			keys = keysFound(lkpmp)
		}
	}
	return
}

func (mp *Map) Values(ks ...interface{}) (values []interface{}) {
	if validMap(mp) {
		values = mapValues(mp, nil, ks...)
	}
	return
}

func valuesFound(ksmp *Map) (values []interface{}) {
	if validMap(ksmp) {
		ksmp.RLock()
		defer ksmp.RUnlock()
		if ksmp != nil && len(ksmp.imp) > 0 {
			values = make([]interface{}, len(ksmp.imp))
			valsi := 0
			for _, v := range ksmp.imp {
				values[valsi] = v
				valsi++
			}
		}
	}
	return
}

func mapValues(mp *Map, mphndlr *MapHandler, ks ...interface{}) (values []interface{}) {
	if len(ks) == 0 {
		values = valuesFound(mp)
	} else {
		var lkpmp *Map = nil
		func() {
			if lkv := mapFind(mp, nil, ks...); lkv != nil {
				lkpmp, _ = lkv.(*Map)
			}
		}()
		if lkpmp != nil {
			values = valuesFound(lkpmp)
		}
	}
	return
}

func (mp *Map) Put(name interface{}, a ...interface{}) (putit bool) {
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
		if validMap(mp) {
			putit = mapPut(mp, nil, true, a...)
		}
	}
	return
}

func mapPut(mp *Map, mphndlr *MapHandler, candispose bool, a ...interface{}) (putit bool) {
	if validMap(mp) && len(a) > 0 {
		func() {
			mp.Lock()
			defer mp.Unlock()
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
								mp.imp[k] = prepNewVal(v)
								putit = true
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
		}()
	}
	return
}

func prepNewVal(v interface{}) (vn interface{}) {
	if v != nil {
		if m, mok := v.(map[string]interface{}); mok && len(m) > 0 {
			vmp := NewMap()
			for k, v := range m {
				vmp.Put(k, v)
			}
			vn = vmp
		} else if arv, arvok := v.([]interface{}); arvok {
			for an, av := range arv {
				arv[an] = prepNewVal(av)
			}
			vn = arv
		} else {
			vn = v
		}
	}
	return
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

func (mp *Map) Close(ks ...interface{}) (closed bool) {
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
			closed = true
		} else {
			var lkpmp *Map = nil
			var prvlkmp *Map = mp
			if len(ks) > 1 {
				func() {
					mp.Lock()
					defer mp.Unlock()
					if v := mapFind(mp, mphndlr, ks[:len(ks)-1]); v != nil {
						prvlkmp, _ = v.(*Map)
						v = nil
					}
				}()
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
				closed = true
			}
		}
	}
	return
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

func (mp *Map) Reset(ks ...interface{}) bool {
	return false
}

func (mp *Map) Clear(ks ...interface{}) (cleared bool) {
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
					ks = make([]interface{}, len(mp.imp))
					ksi := 0
					for k := range mp.imp {
						ks[ksi] = k
						ksi++
					}
					mapRemove(mp, mphndlr, ks...)
				}
			}()
			cleared = true
		} else {
			var lkpmp *Map = nil
			if v := mp.Find(ks...); v != nil {
				lkpmp, _ = v.(*Map)
				v = nil
			}
			if lkpmp != nil {
				cleared = lkpmp.Clear(mphndlr)
			}
		}
	}
	return
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

func (mp *Map) IsMap(ks ...interface{}) (ismap bool) {
	if len(ks) > 0 {
		func() {
			ismap = mapIsMap(mp, nil, ks...)
		}()
	}
	return
}

func baseMapFind(mp *Map, mphndlr *MapHandler, ks ...interface{}) (found bool, vfound interface{}) {
	var lkpmp *Map = mp
	if ksl := len(ks); ksl > 0 && validMap(lkpmp) {
		ksi := 0
		var subfind = func() (valf interface{}, found bool) {
			var currentMap *Map = nil
			if validMap(lkpmp) {
				currentMap = lkpmp
				currentMap.RLock()
				defer currentMap.RUnlock()
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
		var vfnd interface{} = nil
		var fnd bool = false
		for ksi < ksl {
			if vfnd, fnd = subfind(); fnd {
				found, vfound = fnd, vfnd
				break
			}
		}
	}
	return
}

func mapIsMap(mp *Map, mphndlr *MapHandler, ks ...interface{}) (ismap bool) {
	if found, vfound := baseMapFind(mp, mphndlr, ks...); found && vfound != nil {
		_, ismap = vfound.(*Map)
	}
	return
}

func (mp *Map) Exists(ks ...interface{}) (exist bool) {
	exist = mapExists(mp, nil, ks...)
	return
}

func mapExists(mp *Map, mphndlr *MapHandler, ks ...interface{}) (exist bool) {
	if found, _ := baseMapFind(mp, mphndlr, ks...); found {
		exist = found
	}
	return
}

func (mp *Map) Find(ks ...interface{}) (value interface{}) {
	value = mapFind(mp, nil, ks...)
	return
}

func mapFind(mp *Map, mphndlr *MapHandler, ks ...interface{}) (vfound interface{}) {
	_, vfound = baseMapFind(mp, mphndlr, ks...)
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

func (mp *Map) IsMapAt(k interface{}, a ...interface{}) (ismap bool) {
	if len(a) == 0 {
		if a != nil {
			a = append([]interface{}{k}, a)
		}
	} else {
		a = append([]interface{}{k}, a...)
	}
	ismap = mapIsMapAt(mp, nil, a...)
	return
}

func mapIsMapAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (ismap bool) {
	if validMap(mp) && len(a) > 1 {
		var lkpmp *Map = mp
		ks := a[0 : len(a)-1]
		a = a[len(a)-1:]
		var arrv []interface{} = nil
		if arv, atargsok := a[0].([]interface{}); atargsok {
			arrv = arv[:]
		} else {
			arrv = []interface{}{a[0]}
		}
		if ksl := len(ks); ksl > 0 {
			ksi := 0
			var subfind = func() (valf interface{}, found bool) {
				var currentMap *Map = nil
				if lkpmp != nil {
					currentMap = lkpmp
					currentMap.RLock()
					defer currentMap.RUnlock()
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
					if arv, arrvok := valf.([]interface{}); arrvok {
						for an, ad := range arrv {
							if adi, aierr := strconv.ParseInt(fmt.Sprint(ad), 0, 64); aierr == nil && adi > -1 {
								if ai := int(adi); ai > -1 && ai < len(arv) {
									if (an + 1) < len(arrv) {
										if arv, arrvok = arv[ai].([]interface{}); arrvok {
											continue
										} else {
											break
										}
									} else {
										if av := arv[ai]; av != nil {
											_, ismap = av.(*Map)
										}
										break
									}
								} else {
									break
								}

							} else {
								break
							}
						}
					}
					break
				}
			}
		}
	}
	return
}

func (mp *Map) ExistsAt(k interface{}, a ...interface{}) (exist bool) {
	if validMap(mp) {
		if len(a) == 0 {
			a = append([]interface{}{k}, a)
		} else {
			a = append([]interface{}{k}, a...)
		}
		exist = mapExistsAt(mp, nil, a...)
	}
	return
}

func mapExistsAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (exists bool) {
	if mp != nil && len(a) > 1 {
		var lkpmp *Map = mp
		ks := a[0 : len(a)-1]
		a = a[len(a)-1:]
		var arrv []interface{} = nil
		if arv, atargsok := a[0].([]interface{}); atargsok {
			arrv = arv[:]
		} else {
			arrv = []interface{}{a[0]}
		}
		if ksl := len(ks); ksl > 0 {
			ksi := 0
			var subfind = func() (valf interface{}, found bool) {
				var currentMap *Map = nil
				if lkpmp != nil {
					currentMap = lkpmp
					currentMap.RLock()
					defer currentMap.RUnlock()
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
					if arv, arrvok := valf.([]interface{}); arrvok {
						for an, ad := range arrv {
							if adi, aierr := strconv.ParseInt(fmt.Sprint(ad), 0, 64); aierr == nil && adi > -1 {
								if ai := int(adi); ai > -1 && ai < len(arv) {
									if (an + 1) < len(arrv) {
										if arv, arrvok = arv[ai].([]interface{}); arrvok {
											continue
										} else {
											break
										}
									} else {
										exists = true
										break
									}
								} else {
									break
								}

							} else {
								break
							}
						}
					}
					break
				}
			}
		}
	}
	return
}

func (mp *Map) At(k interface{}, a ...interface{}) (arv interface{}) {
	if validMap(mp) {
		func() {
			if len(a) == 0 {
				a = append([]interface{}{k}, a)
			} else {
				a = append([]interface{}{k}, a...)
			}
			arv = mapAt(mp, nil, a...)
		}()
	}
	return
}

func mapAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (av interface{}) {
	if validMap(mp) && len(a) > 1 {
		var lkpmp *Map = mp
		ks := a[0 : len(a)-1]
		a = a[len(a)-1:]
		var arrv []interface{} = nil
		if arv, atargsok := a[0].([]interface{}); atargsok {
			arrv = arv[:]
		} else {
			arrv = []interface{}{a[0]}
		}
		if ksl := len(ks); ksl > 0 {
			ksi := 0
			var subfind = func() (valf interface{}, found bool) {
				var currentMap *Map = nil
				if validMap(lkpmp) {
					currentMap = lkpmp
					currentMap.RLock()
					defer currentMap.RUnlock()
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
					if arv, arrvok := valf.([]interface{}); arrvok {
						for an, ad := range arrv {
							if adi, aierr := strconv.ParseInt(fmt.Sprint(ad), 0, 64); aierr == nil && adi > -1 {
								if ai := int(adi); ai > -1 && ai < len(arv) {
									if (an + 1) < len(arrv) {
										if arv, arrvok = arv[ai].([]interface{}); arrvok {
											continue
										} else {
											break
										}
									} else {
										av = arv[ai]
										break
									}
								} else {
									break
								}

							} else {
								break
							}
						}
					}
					break
				}
			}
		}
	}
	return
}

func (mp *Map) FocusAt(k interface{}, a ...interface{}) (focused bool) {
	return
}

func (mp *Map) ClearAt(k interface{}, a ...interface{}) (cleared bool) {
	if len(a) == 0 {
		a = append([]interface{}{k}, a)
	} else {
		a = append([]interface{}{k}, a...)
	}
	cleared = mapClearAt(mp, nil, a...)

	return
}

func mapClearAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (cleared bool) {
	if validMap(mp) && len(a) > 1 {
		var lkpmp *Map = nil
		func() {
			if v := mapAt(mp, mphndlr, a...); v != nil {
				lkpmp, _ = v.(*Map)
			}
		}()
		if lkpmp != nil {
			lkpmp.Clear(mphndlr)
			cleared = true
		}
	}
	return
}

func (mp *Map) CloseAt(k interface{}, a ...interface{}) (closed bool) {
	if len(a) == 0 {
		a = append([]interface{}{k}, a)
	} else {
		a = append([]interface{}{k}, a...)
	}
	closed = mapCloseAt(mp, nil, a...)
	return
}

func mapCloseAt(mp *Map, mphndlr *MapHandler, a ...interface{}) (closed bool) {
	if validMap(mp) && len(a) > 1 {
		var lkpmp *Map = mp
		ks := a[0 : len(a)-1]
		a = a[len(a)-1:]
		var arrv []interface{} = nil
		if arv, atargsok := a[0].([]interface{}); atargsok {
			arrv = arv[:]
		} else {
			arrv = []interface{}{a[0]}
		}
		if ksl := len(ks); ksl > 0 {
			ksi := 0
			var subfind = func() (valf interface{}, found bool) {
				var currentMap *Map = nil
				if validMap(lkpmp) {
					currentMap = lkpmp
					currentMap.RLock()
					defer currentMap.RUnlock()
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
					if arv, arrvok := valf.([]interface{}); arrvok {
						for an, ad := range arrv {
							if adi, aierr := strconv.ParseInt(fmt.Sprint(ad), 0, 64); aierr == nil && adi > -1 {
								if ai := int(adi); ai > -1 && ai < len(arv) {
									if (an + 1) < len(arrv) {
										if arv, arrvok = arv[ai].([]interface{}); arrvok {
											continue
										} else {
											break
										}
									} else {
										if av := arv[ai]; av != nil {
											disposeMapVal(mp, nil, av)
											arv[ai] = nil
											var narv []interface{} = nil
											if ai > 0 {
												narv = append(arv[:ai], arv[ai+1:]...)
											} else {
												narv = []interface{}{}
											}
											lkpmp.imp[ks[ksi-1]] = narv
										}
										break
									}
								} else {
									break
								}
							} else {
								break
							}
						}
					}
					break
				}
			}
		}
	}
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

func NewMap(a ...interface{}) (mp *Map) {
	mp = &Map{lck: &sync.RWMutex{}, lckedlvl: 0, valid: true, imp: map[interface{}]interface{}{}}
	if len(a) > 0 {
		mapPut(mp, nil, false, a...)
	}
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
	handlerPool = &sync.Pool{New: func() interface{} {
		return initMapHandler()
	}}
	glbmp = NewMap()
	glbmphndlr = NewMapHandler(glbmp)
}
