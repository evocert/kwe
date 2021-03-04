package enumeration

import (
	"strings"
	"sync"
)

//Chain struct
type Chain struct {
	head       *Link
	tale       *Link
	reversemap map[*Link]*Link
	forwardmap map[*Link]*Link
	//API
	DoLink        func(*Link) (bool, error)
	DoneLink      func(*Link) error
	ErrorDoneLink func(*Link, error) bool
	ErrorDoLink   func(*Link, error) bool
}

//NewChain isntance
func NewChain(settings ...map[string]interface{}) (chn *Chain) {
	chn = &Chain{head: nil, tale: nil, reversemap: map[*Link]*Link{}, forwardmap: map[*Link]*Link{},
		DoLink: nil, DoneLink: nil, ErrorDoneLink: nil, ErrorDoLink: nil}
	if len(settings) == 1 && settings[0] != nil && len(settings[0]) > 0 {
		for k, d := range settings[0] {
			if d != nil {
				if strings.ToLower(k) == "dolink" {
					if chn.DoLink == nil {
						if dolnk, dolnkok := d.(func(*Link) (bool, error)); dolnkok && dolnk != nil {
							chn.DoLink = dolnk
						}
					}
				} else if strings.ToLower(k) == "errordolink" {
					if chn.ErrorDoLink == nil {
						if errdolnk, errdolnkok := d.(func(*Link, error) bool); errdolnkok && errdolnk != nil {
							chn.ErrorDoLink = errdolnk
						}
					}
				} else if strings.ToLower(k) == "donelink" {
					if chn.DoneLink == nil {
						if donelnk, donelnkok := d.(func(*Link) error); donelnkok && donelnk != nil {
							chn.DoneLink = donelnk
						}
					}
				} else if strings.ToLower(k) == "errordonelink" {
					if chn.ErrorDoneLink == nil {
						if errdonelnk, errdonelnkok := d.(func(*Link, error) bool); errdonelnkok && errdonelnk != nil {
							chn.ErrorDoneLink = errdonelnk
						}
					}
				}
			}
		}
	}
	return
}

//Link struct
type Link struct {
	chn   *Chain
	value interface{}
	done  bool
	//API
	DoLink        func(*Link) (bool, error)
	DoneLink      func(*Link) error
	ErrorDoneLink func(*Link, error) bool
	ErrorDoLink   func(*Link, error) bool
	Removed       func(*Link)
}

//Done  set lnk.done to true
func (lnk *Link) Done() {
	if lnk != nil {
		if !lnk.done {
			lnk.done = true
		}
	}
}

//Value of *Link
func (lnk *Link) Value() interface{} {
	return lnk.value
}

//Next *Link
func (lnk *Link) Next() *Link {
	return lnk.chn.forwardmap[lnk]
}

//NextValue *Link value interface{}
func (lnk *Link) NextValue() interface{} {
	if nxt := lnk.Next(); nxt != nil {
		return nxt.value
	}
	return nil
}

//Prev *Link
func (lnk *Link) Prev() *Link {
	return lnk.chn.reversemap[lnk]
}

//PrevValue *Link value interface{}
func (lnk *Link) PrevValue() interface{} {
	if prv := lnk.Prev(); prv != nil {
		return prv.value
	}
	return nil
}

//Chain that *Link belongs to
func (lnk *Link) Chain() *Chain {
	return lnk.chn
}

type FuncDoLink func(*Link) (bool, error)
type FuncErrorDoLink func(*Link, error) bool
type FuncDoneLink func(*Link) error
type FuncErrorDoneLink func(*Link, error) bool

type FuncRemoved func(*Link)

func (chn *Chain) newLink(value interface{}, settings ...map[string]interface{}) (lnk *Link) {
	if chn != nil {
		lnk = &Link{chn: chn, value: value,
			DoLink: nil, DoneLink: nil, ErrorDoneLink: nil, ErrorDoLink: nil}
		if len(settings) == 1 && settings[0] != nil && len(settings[0]) > 0 {
			for k, d := range settings[0] {
				if d != nil {
					if strings.ToLower(k) == "dolink" {
						if lnk.DoLink == nil {
							if dolnk, dolnkok := d.(FuncDoLink); dolnkok && dolnk != nil {
								lnk.DoLink = dolnk
							}
						}
					} else if strings.ToLower(k) == "errordolink" {
						if lnk.ErrorDoLink == nil {
							if errdolnk, errdolnkok := d.(FuncErrorDoLink); errdolnkok && errdolnk != nil {
								lnk.ErrorDoLink = errdolnk
							}
						}
					} else if strings.ToLower(k) == "donelink" {
						if lnk.DoneLink == nil {
							if donelnk, donelnkok := d.(FuncDoneLink); donelnkok && donelnk != nil {
								lnk.DoneLink = donelnk
							}
						}
					} else if strings.ToLower(k) == "errordonelink" {
						if lnk.ErrorDoneLink == nil {
							if errdonelnk, errdonelnkok := d.(FuncErrorDoneLink); errdonelnkok && errdonelnk != nil {
								lnk.ErrorDoneLink = errdonelnk
							}
						}
					} else if strings.ToLower(k) == "removed" {
						if lnk.Removed == nil {
							if removedlnk, removedlnkok := d.(FuncRemoved); removedlnkok && removedlnk != nil {
								lnk.Removed = removedlnk
							}
						}
					}
				}
			}
		}
	}
	return
}

//Add value interface{} and return *Link in *Chain that represents value interface{}
func (chn *Chain) Add(value interface{}, settings ...map[string]interface{}) (lnk *Link) {
	if lnk = chn.newLink(value, settings...); lnk != nil {
		if chn.head == nil && chn.tale == nil {
			chn.head = lnk
			chn.tale = lnk
			chn.forwardmap[lnk] = nil
			chn.reversemap[lnk] = nil
		} else if chn.head != nil && chn.tale != nil {
			chn.forwardmap[chn.tale] = lnk
			chn.reversemap[lnk] = chn.tale
			chn.forwardmap[lnk] = nil
			chn.tale = lnk
		}
	}
	return
}

type chainmode int

const (
	cminsertbefore chainmode = iota
	cminsertafter
	cmadd
	cmremove
)

func insert(cm chainmode, chn *Chain, lnk *Link, slnks ...*Link) {
	if cm == cminsertbefore || cm == cminsertafter {
		if sl := len(slnks); sl > 0 {
			if chn != nil && lnk != nil && lnk.chn == chn {
				var crntlnk *Link = lnk
				var insertlnkfnc func(*Link)
				switch cm {
				case cminsertbefore:
					insertlnkfnc = func(slnk *Link) {
						prvlnk := chn.reversemap[crntlnk]
						chn.reversemap[crntlnk] = slnk
						chn.forwardmap[slnk] = crntlnk
						chn.reversemap[slnk] = prvlnk
						if prvlnk == nil {
							chn.head = slnk
						} else {
							chn.forwardmap[prvlnk] = slnk
						}
						crntlnk = slnk
					}
				case cminsertafter:
					insertlnkfnc = func(slnk *Link) {
						nxtlnk := chn.forwardmap[crntlnk]
						chn.forwardmap[crntlnk] = slnk
						chn.reversemap[slnk] = crntlnk
						chn.forwardmap[slnk] = nxtlnk
						if nxtlnk == nil {
							chn.tale = slnk
						} else {
							chn.reversemap[nxtlnk] = slnk
						}
						crntlnk = slnk
					}
				}
				for sn := range slnks {
					switch cm {
					case cminsertbefore:
						insertlnkfnc(slnks[sl-(sn+1)])
					case cminsertafter:
						insertlnkfnc(slnks[sn])
					}
				}
			}
		}
	}
}

//InsertBefore insert values ...interface{} before lnk*Link
func (chn *Chain) InsertBefore(settings map[string]interface{}, lnk *Link, values ...interface{}) bool {
	if vl := len(values); vl > 0 {
		if lnk != nil && chn != nil && lnk.chn == chn {
			slnks := make([]*Link, vl)
			lnksettings := map[string]interface{}{}
			if lnk.DoLink != nil {
				lnksettings["dolink"] = lnk.DoLink
			}
			if lnk.ErrorDoLink != nil {
				lnksettings["errordolink"] = lnk.ErrorDoLink
			}
			if lnk.ErrorDoneLink != nil {
				lnksettings["errordonelink"] = lnk.ErrorDoneLink
			}
			if lnk.DoneLink != nil {
				lnksettings["donelink"] = lnk.DoneLink
			}
			for vn, val := range values {
				if settings != nil {
					slnks[vn] = chn.newLink(val, settings)
				} else {
					slnks[vn] = chn.newLink(val, lnksettings)
				}
			}
			insert(cminsertbefore, chn, lnk, slnks...)
			slnks = nil
			values = nil
			return true
		}
	}
	return false
}

//InsertAfter insert values ...interface{} after lnk*Link
func (chn *Chain) InsertAfter(settings map[string]interface{}, lnk *Link, values ...interface{}) bool {
	if vl := len(values); vl > 0 {
		if lnk != nil && chn != nil && lnk.chn == chn {
			slnks := make([]*Link, vl)
			lnksettings := map[string]interface{}{}
			if lnk.DoLink != nil {
				lnksettings["dolink"] = lnk.DoLink
			}
			if lnk.ErrorDoLink != nil {
				lnksettings["errordolink"] = lnk.ErrorDoLink
			}
			if lnk.ErrorDoneLink != nil {
				lnksettings["errordonelink"] = lnk.ErrorDoneLink
			}
			if lnk.DoneLink != nil {
				lnksettings["donelink"] = lnk.DoneLink
			}
			for vn, val := range values {
				if settings != nil {
					slnks[vn] = chn.newLink(val, settings)
				} else {
					slnks[vn] = chn.newLink(val, lnksettings)
				}
			}
			insert(cminsertafter, chn, lnk, slnks...)
			slnks = nil
			values = nil
			return true
		}
	}
	return false
}

//Remove ...*Link from *Chain - note *Link is also disposed
func (chn *Chain) Remove(link ...*Link) (rmvd bool) {
	if len(link) > 0 {
		for _, lnk := range link {
			if lnk != nil && chn != nil && lnk.chn == chn {
				nxtlnk := lnk.Next()
				prvlnk := lnk.Prev()
				delete(chn.forwardmap, lnk)
				delete(chn.reversemap, lnk)
				if chn.head == lnk {
					chn.head = nxtlnk
				}
				if chn.tale == lnk {
					chn.tale = prvlnk
				}
				if nxtlnk != nil {
					chn.reversemap[nxtlnk] = prvlnk
				}
				if prvlnk != nil {
					chn.forwardmap[prvlnk] = nxtlnk
				}
				if lnk.Removed != nil {
					lnk.Removed(lnk)
				}
				if !rmvd {
					rmvd = true
				}
			}
		}
	}
	return
}

//Back of chain
func (chn *Chain) Back() (lnk *Link) {
	lnk = chn.tale
	return
}

//Front of chain
func (chn *Chain) Front() (lnk *Link) {
	lnk = chn.head
	return
}

//Size of chain
func (chn *Chain) Size() int {
	if chn != nil {
		if len(chn.forwardmap) == len(chn.reversemap) {
			return len(chn.forwardmap)
		}
	}
	return 0
}

//Do dolnk (func(*Link) (boolm,error)) iterate over *Chain
// iterate until end or dolnk return true or error
func (chn *Chain) Do(dolnk FuncDoLink, errdolnk FuncErrorDoLink, donelnk FuncDoneLink, errdonelnk FuncErrorDoneLink) (diddo bool) {
	if chn != nil && chn.head != nil && chn.tale != nil {
		var done = false
		var err error = nil
		wg := &sync.WaitGroup{}
		cnt := 0
		func() {
			for e := chn.Front(); e != nil && !done; e = e.Next() {
				cnt++
				func() {
					wg.Add(1)
					defer wg.Wait()
					go func() {
						defer wg.Done()
						if done, err = dolnk(e); done || err != nil {
							if done && err == nil {
								if donelnk != nil {
									if err = donelnk(e); err != nil {
										if errdonelnk != nil {
											if errdonelnk(e, err) {

											}
										}
									}
								}
							}
							if !done && errdonelnk != nil {
								done = errdonelnk(e, err)
							}
						} else if err != nil {
							if errdolnk != nil {
								if done = errdolnk(e, err); done {
									if donelnk != nil {
										if err = donelnk(e); err != nil {
											if errdonelnk != nil {
												if errdonelnk(e, err) {

												}
											}
										}
									}
								}
							}
						}
					}()
				}()
			}
			diddo = true
		}()
	}
	return false
}
