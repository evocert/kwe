package caching

import (
	"encoding/json"
	"io"
	"runtime"
	"sync"

	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
)

type MapHandler struct {
	rnble  active.Runtime
	intern bool
	*Map
	dspsng bool
}

func mapHandlerFinalize(mphndlr *MapHandler) {
	if mphndlr != nil {
		mphndlr.Close()
		mphndlr = nil
	}
}

func NewMapHandler(a ...interface{}) (mphndlr *MapHandler) {
	var mp *Map = nil
	var rnble active.Runtime = nil
	for _, d := range a {
		if mp == nil {
			mp, _ = d.(*Map)
		}
		if rnble == nil {
			rnble, _ = d.(active.Runtime)
		}
	}

	if mp != nil {
		mphndlr = &MapHandler{Map: mp, intern: false, dspsng: false, rnble: rnble}
	} else {
		mphndlr = &MapHandler{Map: NewMap(), intern: true, dspsng: false, rnble: rnble}
	}
	runtime.SetFinalizer(mphndlr, mapHandlerFinalize)
	return
}

func (mphndlr *MapHandler) Keys(k ...interface{}) (ks []interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		ks = mphndlr.Map.Keys(k...)
	}
	return
}

func (mphndlr *MapHandler) String() (s string) {
	if mphndlr != nil && mphndlr.Map != nil {
		s = mapString(mphndlr.Map, mphndlr)
	}
	return
}

func (mphndlr *MapHandler) Reader() (mprdr *iorw.EOFCloseSeekReader) {
	if mphndlr != nil && mphndlr.Map != nil {
		mprdr = mapReader(mphndlr.Map, mphndlr)
	} else {
		mprdr = iorw.NewEOFCloseSeekReader(nil)
	}
	return
}

func (mphndlr *MapHandler) Values(k ...interface{}) (vs []interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		vs = mphndlr.Map.Values(k...)
	}
	return
}

func (mphndlr *MapHandler) ReplaceKey(a ...interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		mapReplaceKey(mphndlr.Map, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Remove(a ...interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		mapRemove(false, mphndlr.Map, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Push(k interface{}, a ...interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		if len(a) == 0 {
			a = []interface{}{k}
		} else {
			a = append([]interface{}{k}, a...)
		}
		mapPush(mphndlr.Map, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Find(ks ...interface{}) (vfound interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		vfound = mapFind(mphndlr.Map, mphndlr, ks...)
	}
	return
}

func (mphndlr *MapHandler) Size() (size int) {
	if mphndlr != nil && mphndlr.Map != nil {
		size = mapSize(mphndlr.Map, mphndlr)
	}
	return size
}

func (mphndlr *MapHandler) ValueByIndex(index int64) (v interface{}) {
	if mphndlr != nil && mphndlr.Map != nil {
		v = mapValueByIndex(mphndlr.Map, mphndlr, index)
	}
	return
}

func (mphndlr *MapHandler) Clear() {
	if mphndlr != nil && mphndlr.Map != nil {
		mapClear(mphndlr.Map, mphndlr)
	}
}

func (mphndlr *MapHandler) Close() {
	if mphndlr != nil {
		if mphndlr.Map != nil {
			if mphndlr.intern {
				mapClose(mphndlr.Map, mphndlr)
				mphndlr.Map.Close()
			}
			mphndlr.Map = nil
		}
		mphndlr = nil
	}
}

type mapAction int

const (
	actnunknown mapAction = iota
	actnnone
	actnput
	actnpush
	actnremove
	actnreplacekey
	actnclear
	actnclose
	actnfind
	actnread
)

type Map struct {
	lck     *sync.RWMutex
	keys    *enumeration.List
	kvndm   map[*enumeration.Node]*enumeration.Node
	vkndm   map[*enumeration.Node]*enumeration.Node
	values  *enumeration.List
	lstactn mapAction
}

func mapFinalize(mp *Map) {
	if mp != nil {
		mp.Close()
		mp = nil
	}
}

//NewMap return instance of *Map
func NewMap() (mp *Map) {
	mp = &Map{
		lck:     &sync.RWMutex{},
		keys:    enumeration.NewList(true),
		kvndm:   map[*enumeration.Node]*enumeration.Node{},
		values:  enumeration.NewList(),
		vkndm:   map[*enumeration.Node]*enumeration.Node{},
		lstactn: actnnone}
	runtime.SetFinalizer(mp, mapFinalize)
	return
}

func (mp *Map) lastAction(nxtactn ...mapAction) mapAction {
	if mp != nil {
		if len(nxtactn) == 1 {
			mp.lck.Lock()
			defer mp.lck.Unlock()
			if lstactn, nxtctn := mp.lstactn, nxtactn[0]; nxtctn != lstactn {
				if (lstactn == actnclear || lstactn == actnclose) && nxtctn == actnnone {
					mp.lstactn = nxtactn[0]
				} else if lstactn != actnclear && lstactn != actnclose {
					mp.lstactn = nxtactn[0]
				}
			}
			return mp.lstactn
		}
		mp.lck.RLock()
		defer mp.lck.RUnlock()
		return mp.lstactn
	}
	return actnunknown
}

func (mp *Map) Handler() (mphndlr *MapHandler) {
	mphndlr = NewMapHandler(mp)
	return
}

func (mp *Map) Size() (size int) {
	if mp != nil {
		size = mapSize(mp, nil)
	}
	return size
}

func mapSize(mp *Map, mphndlr *MapHandler) (size int) {
	if mp != nil {
		func() {
			if mphndlr != nil {
				mp.lck.RLock()
				defer mp.lck.RUnlock()
			}
			if mp.keys != nil {
				size = mp.keys.Length()
			}
		}()
	}
	return
}

func (mp *Map) String() (s string) {
	s = mapString(mp, nil)
	return
}

func mapString(mp *Map, mphndlr *MapHandler) (s string) {
	if mp != nil {
		s, _ = iorw.ReaderToString(mapReader(mp, mphndlr))
	}
	return
}

func (mp *Map) Reader() (mprdr *iorw.EOFCloseSeekReader) {
	mprdr = mapReader(mp, nil)
	return
}

func mapReader(mp *Map, mphndlr *MapHandler) (mprdr *iorw.EOFCloseSeekReader) {
	var rdr io.Reader = nil
	if mp != nil {
		if lstactn := mp.lastAction(actnread); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnread {
			func() {
				defer mp.lastAction(actnnone)
				wg := &sync.WaitGroup{}
				wg.Add(1)
				pi, pw := io.Pipe()
				go func() {
					defer pw.Close()
					ks := mp.Keys()
					vs := mp.Values()
					wg.Done()
					jsnencd := json.NewEncoder(pw)
					iorw.Fprint(pw, "{")
					if ksl, vsl := len(ks), len(vs); ksl > 0 && ksl == vsl {
						for kn, k := range ks {
							jsnencd.Encode(k)
							iorw.Fprint(pw, ":")
							if vmp, vmpok := vs[kn].(*Map); vmpok {
								if mphndlr == nil {
									iorw.Fprint(pw, vmp)
								} else {
									iorw.Fprint(pw, mapReader(vmp, mphndlr))
								}
							} else {
								jsnencd.Encode(vs[kn])
							}
						}
					}
					iorw.Fprint(pw, "}")
				}()
				wg.Wait()
				rdr = pi
			}()
		}
	}
	mprdr = iorw.NewEOFCloseSeekReader(rdr)
	return
}

func (mp *Map) Find(ks ...interface{}) (vfound interface{}) {
	vfound = mapFind(mp, nil, ks...)
	return
}

func mapFind(mp *Map, mphndlr *MapHandler, ks ...interface{}) (vfound interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnfind); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnfind {
			func() {
				defer mp.lastAction(actnnone)
				var lkpmp *Map = mp
				if ksl := len(ks); ksl > 0 {
					for kn, k := range ks {
						func() {
							if mphndlr != nil {
								lkpmp.lck.RLock()
								defer lkpmp.lck.RUnlock()
							}
							if knde := lkpmp.keys.ValueNode(k); knde != nil {
								if vnde := lkpmp.kvndm[knde]; vnde != nil {
									if vl := vnde.Value(); vl != nil {
										if vmp, vmpok := vl.(*Map); vmpok {
											if (kn + 1) == ksl {
												vfound = vmp
											} else {
												lkpmp = vmp
											}
										} else if (kn + 1) == ksl {
											vfound = vl
										}
									} else {
										return
									}
								} else {
									return
								}
							} else {
								return
							}
						}()
					}
				}
			}()
		}
	}
	return
}

func (mp *Map) Keys(k ...interface{}) (ks []interface{}) {
	if mp != nil && mp.keys != nil {
		mp.keys.Do(func(knde *enumeration.Node, val interface{}) bool {
			if ks == nil {
				ks = []interface{}{}
			}
			ks = append(ks, val)
			return false
		}, nil, nil)
	}
	return
}

func (mp *Map) Values(k ...interface{}) (vs []interface{}) {
	if mp != nil && mp.values != nil {
		mp.values.Do(func(knde *enumeration.Node, val interface{}) bool {
			if vs == nil {
				vs = []interface{}{}
			}
			vs = append(vs, val)
			return false
		}, nil, nil)
	}
	return
}

func (mp *Map) Remove(a ...interface{}) {
	mapRemove(false, mp, nil, a...)
}

func mapRemove(forceRemove bool, mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if len(a) > 0 {
			if lstactn := mp.lastAction(actnremove); (forceRemove && (lstactn == actnclear || lstactn == actnclose)) || (!forceRemove && (!(lstactn == actnclear || lstactn == actnclose) && lstactn == actnremove)) {
				defer mp.lastAction(actnnone)
				if keys := mp.keys; keys.Length() > 0 {
					for _, d := range a {
						func() {
							if mphndlr != nil {
								mp.lck.Lock()
								defer mp.lck.Unlock()
							}
							if knde := keys.ValueNode(d); knde != nil {
								k := knde.Value()
								vnde := mp.kvndm[knde]
								v := vnde.Value()
								knde.Dispose(
									//REMOVED KEY
									func(nde *enumeration.Node, val interface{}) {
										if nde == knde {
											delete(mp.kvndm, knde)
										}
									},
									//DISPOSED KEY
									func(nde *enumeration.Node, val interface{}) {
										vnde.Dispose(
											//REMOVED VALUE
											func(nde *enumeration.Node, val interface{}) {
												if vnde == nde {
													delete(mp.vkndm, vnde)
												}
											},
											//DISPOSED VALUE
											func(nde *enumeration.Node, val interface{}) {

											})
									})
								if k != nil {

								}
								if v != nil {

								}
							}
						}()
					}
				}
			}
		}
	}
}

func (mp *Map) Clear() {
	mapClear(mp, nil)
}

func (mp *Map) Close() {
	if mp != nil {
		mapClose(mp, nil)
		mp = nil
	}
}

func mapClear(mp *Map, mphndlr *MapHandler) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnclear); lstactn == actnclear || lstactn == actnclose {
			func() {
				defer mp.lastAction(actnnone)
				mapRemove(true, mp, mphndlr, mp.Keys()...)
			}()
		}
	}
}

func mapClose(mp *Map, mphndlr *MapHandler) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnclose); lstactn == actnclose {
			func() {
				defer mp.lastAction(actnnone)
				mapRemove(true, mp, mphndlr, mp.Keys()...)
				if mp.keys != nil {
					mp.keys.Dispose(func(n *enumeration.Node, i interface{}) {}, func(n *enumeration.Node, i interface{}) {})
					mp.keys = nil
				}
				if mp.values != nil {
					mp.values.Dispose(func(n *enumeration.Node, i interface{}) {}, func(n *enumeration.Node, i interface{}) {})
					mp.values = nil
				}
			}()
		}
		mp = nil
	}
}

func (mp *Map) ValueByIndex(index int64) (v interface{}) {
	if mp != nil {
		v = mapValueByIndex(mp, nil, index)
	}
	return
}

func mapValueByIndex(mp *Map, mphndlr *MapHandler, index int64) (v interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnfind); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnfind {
			func() {
				if mphndlr != nil {
					mp.lck.RLock()
					defer mp.lck.RUnlock()
				}
				if mp.keys != nil {
					index++
					stri := int64(0)
					for _, vnde := range mp.kvndm {
						stri++
						if stri >= index {
							v = vnde.Value()
							break
						}
					}
				}
			}()
		}
	}
	return
}

func (mp *Map) Push(k interface{}, a ...interface{}) {
	if len(a) == 0 {
		a = []interface{}{k}
	} else {
		a = append([]interface{}{k}, a...)
	}
	mapPush(mp, nil, a...)
}

func mapPush(mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnpush); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnpush {
			func() {
				defer mp.lastAction(actnnone)
				keys := mp.keys
				values := mp.values
				for {
					if al := len(a); al > 0 {
						if al%2 == 0 {
							k := a[0]
							if k != nil {
								a = a[1:]
								v := a[0]
								func() {
									var knde *enumeration.Node = nil
									var kndecngd bool
									var vldky bool
									var vnde *enumeration.Node = nil
									var vndecngd bool
									var vldv bool
									if mphndlr != nil {
										mp.lck.Lock()
										defer mp.lck.Unlock()
									}
									keys.Add(nil, func(cngd bool, valvld bool, idx int, n *enumeration.Node, i interface{}) {
										kndecngd = cngd
										knde = n
										vldky = valvld
									}, k)
									if v != nil && vldky {
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
									if vldky && kndecngd {
										values.Add(nil, func(cngd bool, valvld bool, idx int, n *enumeration.Node, i interface{}) {
											if vndecngd = cngd; kndecngd {
												vnde = n
											}
											vldv = valvld
										}, v)
										if kndecngd && vndecngd && vldky && vldv {
											mp.kvndm[knde] = vnde
											mp.vkndm[vnde] = knde
										}
									} else if vldky && !kndecngd {
										vnde = mp.kvndm[knde]
										vnde.Set(v)
									}
								}()
								a = a[1:]
							}
						}
					} else {
						break
					}
				}
			}()
		}
	}
}

func (mp *Map) Put(k interface{}, a ...interface{}) {
	if len(a) == 0 {
		a = []interface{}{k}
	} else {
		a = append([]interface{}{k}, a...)
	}
	mapPut(mp, nil, a...)
}

func mapPut(mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnput); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnput {
			func() {
				defer mp.lastAction(actnnone)
				keys := mp.keys
				values := mp.values
				for {
					if al := len(a); al > 0 {
						if al%2 == 0 {
							k := a[0]
							if k != nil {
								a = a[1:]
								v := a[0]
								func() {
									var knde *enumeration.Node = nil
									var kndecngd bool
									var vldky bool
									var vnde *enumeration.Node = nil
									var vndecngd bool
									var vldv bool
									if mphndlr != nil {
										mp.lck.Lock()
										defer mp.lck.Unlock()
									}
									keys.Add(nil, func(cngd bool, valvld bool, idx int, n *enumeration.Node, i interface{}) {
										kndecngd = cngd
										knde = n
										vldky = valvld
									}, k)
									if v != nil && vldky {
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
									if vldky && kndecngd {
										values.Add(nil, func(cngd bool, valvld bool, idx int, n *enumeration.Node, i interface{}) {
											if vndecngd = cngd; kndecngd {
												vnde = n
											}
											vldv = valvld
										}, v)
										if kndecngd && vndecngd && vldky && vldv {
											mp.kvndm[knde] = vnde
											mp.vkndm[vnde] = knde
										}
									} else if vldky && !kndecngd {
										vnde = mp.kvndm[knde]
										vnde.Set(v)
									}
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
	}
}

func (mp *Map) ReplaceKey(a ...interface{}) {
	mapReplaceKey(mp, nil, a...)
}

func mapReplaceKey(mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.Map
	}
	if mp != nil {
		keys := mp.keys
		for {
			if al := len(a); al > 0 {
				if al%2 == 0 {
					k := a[0]
					if k != nil {
						a = a[1:]
						if v := a[0]; v != nil {
							a = a[1:]
							func() {
								if mphndlr != nil {
									mp.lck.Lock()
									defer mp.lck.Unlock()
								}
								if knde := keys.ValueNode(k); knde != nil {
									if knde.Value() != v {
										knde.Set(v)
									}
								}
							}()
						} else {
							break
						}
					} else {
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
}

func GLOBALMAP() *MapHandler {
	return glbmphndlr
}

var glbmphndlr *MapHandler = nil

func init() {
	glbmphndlr = NewMapHandler(NewMap())
}
