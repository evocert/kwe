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
	mp     *Map
	dspsng bool
	crntmp *Map
}

func mapHandlerFinalize(mphndlr *MapHandler) {
	runtime.SetFinalizer(mphndlr, nil)
	if mphndlr != nil {
		mphndlr.Close()
		runtime.SetFinalizer(mphndlr, nil)
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
		mphndlr = &MapHandler{mp: mp, intern: false, dspsng: false, rnble: rnble, crntmp: nil}
	} else {
		mphndlr = &MapHandler{mp: NewMap(), intern: true, dspsng: false, rnble: rnble, crntmp: nil}
	}
	runtime.SetFinalizer(mphndlr, mapHandlerFinalize)
	return
}

func (mphndlr *MapHandler) Reset() {
	if mphndlr != nil && mphndlr.mp != nil {
		mphndlr.crntmp = mphndlr.mp
	}
}

func (mphndlr *MapHandler) currentmp() (crntmp *Map) {
	if mphndlr != nil {
		if mphndlr.crntmp != nil {
			crntmp = mphndlr.crntmp
		} else {
			crntmp = mphndlr.mp
		}
	}
	return
}

func (mphndlr *MapHandler) Keys(k ...interface{}) (ks []interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		ks = mp.Keys(k...)
	}
	return
}

func (mphndlr *MapHandler) String() (s string) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		s = mapString(mp, mphndlr)
	}
	return
}

func (mphndlr *MapHandler) Reader() (mprdr *iorw.EOFCloseSeekReader) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		mprdr = mapReader(mp, mphndlr)
	} else {
		mprdr = iorw.NewEOFCloseSeekReader(nil)
	}
	return
}

func (mphndlr *MapHandler) Fprint(w io.Writer) (err error) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		err = mapFprint(mp, mphndlr, w)
	}
	return
}

func (mphndlr *MapHandler) Values(k ...interface{}) (vs []interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		vs = mp.Values(k...)
	}
	return
}

func (mphndlr *MapHandler) ReplaceKey(a ...interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		mapReplaceKey(mp, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Remove(a ...interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		mapRemove(false, mp, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Push(k interface{}, a ...interface{}) (length int) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				length = mapPush(mp, mphndlr, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			length = mapPush(mp, mphndlr, a...)
		}
	}
	return
}

func (mphndlr *MapHandler) Shift(k interface{}, a ...interface{}) (length int) {
	length = -1
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				length = mapShift(mp, mphndlr, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			length = mapShift(mp, mphndlr, a...)
		}
	}
	return
}

func (mphndlr *MapHandler) Pop(k interface{}, a ...interface{}) (pop interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				pop = mapPop(mp, mphndlr, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			pop = mapPop(mp, mphndlr, a...)
		}
	}
	return
}

func (mphndlr *MapHandler) Unshift(k interface{}, a ...interface{}) (unshift interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				unshift = mapUnshift(mp, mphndlr, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			unshift = mapUnshift(mp, mphndlr, a...)
		}
	}
	return
}

func (mphndlr *MapHandler) Put(k interface{}, a ...interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		if len(a) == 0 {
			if _, mpsok := k.(map[string]interface{}); mpsok {
				a = []interface{}{k}
			} else if _, mpiok := k.(map[interface{}]interface{}); mpsok || mpiok {
				a = []interface{}{k}
			} else if a != nil {
				a = append([]interface{}{k}, []interface{}{a})
			}
		} else {
			a = append([]interface{}{k}, a...)
		}
		mapPut(mp, mphndlr, a...)
	}
}

func (mphndlr *MapHandler) Find(ks ...interface{}) (vfound interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		vfound = mapFind(mp, mphndlr, ks...)
	}
	return
}

func (mphndlr *MapHandler) Size() (size int) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		size = mapSize(mp, mphndlr)
	}
	return size
}

func (mphndlr *MapHandler) ValueByIndex(index int64) (v interface{}) {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		v = mapValueByIndex(mp, mphndlr, index)
	}
	return
}

func (mphndlr *MapHandler) Clear() {
	if mp := mphndlr.currentmp(); mphndlr != nil && mp != nil {
		mapClear(mp, mphndlr)
	}
}

func (mphndlr *MapHandler) Close() {
	if mphndlr != nil {
		if mphndlr.crntmp != nil {
			mphndlr.crntmp = nil
		}
		if mphndlr.mp != nil {
			if mphndlr.intern {
				mapClose(mphndlr.mp, mphndlr)
				mphndlr.mp.Close()
			}
			mphndlr.mp = nil
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
	actnpop
	atcnshift
	actnunshift
	actnreplacekey
	actnclear
	actnclose
	actnfind
	actnread
	actnwrite
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
		runtime.SetFinalizer(mp, nil)
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
			return func() mapAction {
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
			}()
		} else {
			return func() mapAction {
				mp.lck.RLock()
				defer mp.lck.RUnlock()
				return mp.lstactn
			}()
		}
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

func (mp *Map) Fprint(w io.Writer) (err error) {
	if mp != nil {
		err = mapFprint(mp, nil, w)
	}
	return
}

func mapFprint(mp *Map, mphndlr *MapHandler, w io.Writer) (err error) {
	if mp != nil && w != nil {
		if lstactn := mp.lastAction(actnwrite); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnwrite {
			func() {
				defer mp.lastAction(actnnone)
				encodeMap(w, nil, mp, mphndlr)
			}()
		}
	}
	return
}

func (mp *Map) Reader() (mprdr *iorw.EOFCloseSeekReader) {
	mprdr = mapReader(mp, nil)
	return
}

func encodeMapAVal(w io.Writer, jsnenc *json.Encoder, mp *Map, mphndlr *MapHandler, val interface{}) {
	if val != nil {
		if vmp, vmpok := val.(*Map); vmpok {
			encodeMap(w, jsnenc, vmp, mphndlr)
		} else {
			if varr, varrok := val.([]interface{}); varrok {
				iorw.Fprint(w, "[")
				for vn, va := range varr {
					encodeMapVal(w, jsnenc, mp, mphndlr, va, vn == len(varr)-1)
				}
				iorw.Fprint(w, "]")
			} else {
				jsnenc.Encode(val)
			}
		}
	} else {
		iorw.Fprint(w, "null")
	}
}

func encodeMapVal(w io.Writer, jsnenc *json.Encoder, mp *Map, mphndlr *MapHandler, val interface{}, isLastVal bool) {
	encodeMapAVal(w, jsnenc, mp, mphndlr, val)
	if !isLastVal {
		iorw.Fprint(w, ",")
	}
}

func encodeMap(w io.Writer, jsnenc *json.Encoder, mp *Map, mphndlr *MapHandler) {
	if jsnenc == nil {
		jsnenc = json.NewEncoder(w)
	}
	iorw.Fprint(w, "{")
	var nxtkh *enumeration.Node = nil
	if kh := mp.keys.Head(); kh != nil {
		nxtkh = kh
		for nxtkh != nil {
			jsnenc.Encode(kh.Value())
			iorw.Fprint(w, ":")
			nxtkh = kh.Next()
			encodeMapVal(w, jsnenc, mp, mphndlr, mp.kvndm[kh].Value(), nxtkh == nil)
			kh = nxtkh
		}
	}
	iorw.Fprint(w, "}")
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
					wg.Done()
					encodeMap(pw, nil, mp, mphndlr)
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
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if lstactn := mp.lastAction(actnfind); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnfind {
			func() {
				defer mp.lastAction(actnnone)
				var lkpmp *Map = mp
				if ksl := len(ks); ksl > 0 {
					for kn, k := range ks {
						if !func() bool {
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
												if mphndlr != nil {
													mphndlr.crntmp = vmp
													vfound = mphndlr
												}
												return false
											} else {
												lkpmp = vmp
											}
										} else if (kn + 1) == ksl {
											vfound = vl
											return false
										}
									} else {
										return false
									}
								} else {
									return false
								}
							} else {
								return false
							}
							return true
						}() {
							break
						}
					}
				}
			}()
		}
	}
	return
}

func (mp *Map) Keys(k ...interface{}) (ks []interface{}) {
	if mp != nil && mp.keys != nil {
		if kh := mp.keys.Head(); kh != nil {
			for kh != nil {
				if ks == nil {
					ks = []interface{}{}
				}
				ks = append(ks, kh.Value())
				kh = kh.Next()
			}
		}
	}
	return
}

func (mp *Map) Values(k ...interface{}) (vs []interface{}) {
	if mp != nil && mp.values != nil {
		if kh := mp.keys.Head(); kh != nil {
			for kh != nil {
				if vs == nil {
					vs = []interface{}{}
				}
				vs = append(vs, mp.kvndm[kh].Value())
				kh = kh.Next()
			}
		}
	}
	return
}

func (mp *Map) Remove(a ...interface{}) {
	mapRemove(false, mp, nil, a...)
}

func mapRemove(forceRemove bool, mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if len(a) > 0 {
			if lstactn := mp.lastAction(actnremove); (forceRemove && (lstactn == actnclear || lstactn == actnclose)) || (!forceRemove && (!(lstactn == actnclear || lstactn == actnclose) && lstactn == actnremove)) {
				defer mp.lastAction(actnnone)
				if keys := mp.keys; keys.Length() > 0 {
					for _, d := range a {
						if d != nil {
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
										if vmp, _ := v.(*Map); vmp != nil {
											vmp.Close()
											vmp = nil
										}
									}
								}
							}()
						}
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
		mp = mphndlr.currentmp()
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
		mp = mphndlr.mp
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
					mp.values.Dispose(func(n *enumeration.Node, i interface{}) {}, func(n *enumeration.Node, i interface{}) {
						if i != nil {
							if vmp, _ := i.(*Map); vmp != nil {
								vmp.Close()
								vmp = nil
							}
						}
					})
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
		mp = mphndlr.currentmp()
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
					if mp.values != nil && mp.values.Length() > 0 {
						for vnxt := mp.values.Head(); vnxt != nil; vnxt = vnxt.Next() {
							stri++
							if stri >= index {
								v = vnxt.Value()
								break
							}
						}
					}
				}
			}()
		}
	}
	return
}

func (mp *Map) Shift(k interface{}, a ...interface{}) (length int) {
	length = -1
	if mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				length = mapShift(mp, nil, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			length = mapShift(mp, nil, a...)
		}
	}
	return length
}

func mapShift(mp *Map, mphndlr *MapHandler, a ...interface{}) (length int) {
	length = -1
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if al := len(a); al > 1 {
			ks := a[:al-1]
			arrv := a[al-1]
			if lstactn := mp.lastAction(atcnshift); !(lstactn == actnclear || lstactn == actnclose) && lstactn == atcnshift {
				func() {
					defer mp.lastAction(actnnone)
					var lkpmp *Map = mp
					if ksl := len(ks); ksl > 0 {
						for kn, k := range a {
							if !func() bool {
								if knde := lkpmp.keys.ValueNode(k); knde != nil {
									if vnde := lkpmp.kvndm[knde]; vnde != nil {
										if vl := vnde.Value(); vl != nil {
											if (kn + 2) <= ksl {
												if vmp, vmpok := vl.(*Map); vmpok {
													lkpmp = vmp
												} else {
													return false
												}
											} else if (kn+1) == ksl && lkpmp != nil {
												func() {
													if mphndlr != nil {
														lkpmp.lck.Lock()
														defer lkpmp.lck.Unlock()
													}
													if arv, arrvok := vl.([]interface{}); arrvok {
														arv = append([]interface{}{arrv}, arv...)
														vnde.Set(arv, true)
														length = len(arv)
													}
												}()
												return false
											}
										} else {
											return false
										}
									} else {
										return false
									}
								} else {
									return false
								}
								return true
							}() {
								break
							}
						}
					}
				}()
			}
		}
	}
	return length
}

func (mp *Map) Push(k interface{}, a ...interface{}) (length int) {
	length = -1
	if mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				length = mapPush(mp, nil, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			length = mapPush(mp, nil, a...)
		}
	}
	return length
}

func mapPush(mp *Map, mphndlr *MapHandler, a ...interface{}) (length int) {
	length = -1
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if al := len(a); al > 1 {
			ks := a[:al-1]
			arrv := a[al-1]
			if lstactn := mp.lastAction(actnpush); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnpush {
				func() {
					defer mp.lastAction(actnnone)
					var lkpmp *Map = mp
					if ksl := len(ks); ksl > 0 {
						for kn, k := range ks {
							if !func() bool {
								if knde := lkpmp.keys.ValueNode(k); knde != nil {
									if vnde := lkpmp.kvndm[knde]; vnde != nil {
										if vl := vnde.Value(); vl != nil {
											if (kn + 2) <= ksl {
												if vmp, vmpok := vl.(*Map); vmpok {
													lkpmp = vmp
												} else {
													return false
												}
											} else if (kn+1) == ksl && lkpmp != nil {
												func() {
													if mphndlr != nil {
														lkpmp.lck.Lock()
														defer lkpmp.lck.Unlock()
													}
													if arv, arrvok := vl.([]interface{}); arrvok {
														arv = append(arv, arrv)
														vnde.Set(arv, true)
														length = len(arv)
													}
												}()
												return false
											}
										} else {
											return false
										}
									} else {
										return false
									}
								} else {
									return false
								}
								return true
							}() {
								break
							}
						}
					}
				}()
			}
		}
	}
	return length
}

func (mp *Map) Pop(k interface{}, a ...interface{}) (pop interface{}) {
	if mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				pop = mapPop(mp, nil, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			pop = mapPop(mp, nil, a...)
		}
	}
	return
}

func mapPop(mp *Map, mphndlr *MapHandler, a ...interface{}) (pop interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if al := len(a); al > 1 {
			ks := a
			if lstactn := mp.lastAction(actnpop); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnpop {
				func() {
					defer mp.lastAction(actnnone)
					var lkpmp *Map = mp
					if ksl := len(ks); ksl > 0 {
						for kn, k := range ks {
							if !func() bool {
								if knde := lkpmp.keys.ValueNode(k); knde != nil {
									if vnde := lkpmp.kvndm[knde]; vnde != nil {
										if vl := vnde.Value(); vl != nil {
											if (kn + 1) <= ksl {
												if vmp, vmpok := vl.(*Map); vmpok {
													lkpmp = vmp
												} else {
													return true
												}
											}
											if (kn+1) == ksl && lkpmp != nil {
												func() {
													if mphndlr != nil {
														lkpmp.lck.Lock()
														defer lkpmp.lck.Unlock()
													}
													if arv, arrvok := vl.([]interface{}); arrvok {
														if len(arv) > 0 {
															pop = arv[len(arv)-1]
															arv = arv[:len(arv)-1]
															if arv == nil {
																arv = []interface{}{}
															}
															vnde.Set(arv, true)
														}
													}
												}()
												return true
											}
										} else {
											return false
										}
									} else {
										return false
									}
								} else {
									return false
								}
								return true
							}() {
								break
							}
						}
					}
				}()
			}
		}
	}
	return
}

func (mp *Map) Unshift(k interface{}, a ...interface{}) (unshift interface{}) {
	if mp != nil {
		if len(a) == 0 {
			if a != nil {
				a = append([]interface{}{k}, a)
				unshift = mapUnshift(mp, nil, a...)
			}
		} else {
			a = append([]interface{}{k}, a...)
			unshift = mapUnshift(mp, nil, a...)
		}
	}
	return
}

func mapUnshift(mp *Map, mphndlr *MapHandler, a ...interface{}) (unshift interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
	}
	if mp != nil {
		if al := len(a); al > 1 {
			ks := a
			if lstactn := mp.lastAction(actnunshift); !(lstactn == actnclear || lstactn == actnclose) && lstactn == actnunshift {
				func() {
					defer mp.lastAction(actnnone)
					var lkpmp *Map = mp
					if ksl := len(ks); ksl > 0 {
						for kn, k := range ks {
							if !func() bool {
								if knde := lkpmp.keys.ValueNode(k); knde != nil {
									if vnde := lkpmp.kvndm[knde]; vnde != nil {
										if vl := vnde.Value(); vl != nil {
											if (kn + 1) <= ksl {
												if vmp, vmpok := vl.(*Map); vmpok {
													lkpmp = vmp
												} else {
													return false
												}
											}
											if (kn+1) == ksl && lkpmp != nil {
												func() {
													if mphndlr != nil {
														lkpmp.lck.Lock()
														defer lkpmp.lck.Unlock()
													}
													if arv, arrvok := vl.([]interface{}); arrvok {
														if len(arv) > 0 {
															unshift = arv[0]
															arv = arv[1:]
															if arv == nil {
																arv = []interface{}{}
															}
															vnde.Set(arv, true)
														}
													}
												}()
												return false
											}
										} else {
											return false
										}
									} else {
										return false
									}
								} else {
									return false
								}
								return true
							}() {
								break
							}
						}
					}
				}()
			}
		}
	}
	return
}

func (mp *Map) Put(k interface{}, a ...interface{}) {
	if mp != nil {
		if len(a) == 0 {
			if _, mpsok := k.(map[string]interface{}); mpsok {
				a = []interface{}{k}
			} else if _, mpiok := k.(map[interface{}]interface{}); mpsok || mpiok {
				a = []interface{}{k}
			} else if a != nil {
				a = append([]interface{}{k}, []interface{}{a})
			}
		} else {
			a = append([]interface{}{k}, a...)
		}
		mapPut(mp, nil, a...)
	}
}

func mapPut(mp *Map, mphndlr *MapHandler, a ...interface{}) {
	if mp == nil && mphndlr != nil {
		mp = mphndlr.currentmp()
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
		mp = mphndlr.currentmp()
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
