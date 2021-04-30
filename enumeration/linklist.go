package enumeration

type Node struct {
	lst *List
	val interface{}
}

func (nde *Node) Value() interface{} {
	return nde.val
}

type listaction int

const (
	insertAfter listaction = iota
	insertBefore
)

func (lstactn listaction) String() string {
	if lstactn == insertAfter {
		return "InsertAfter"
	} else if lstactn == insertBefore {
		return "InsertBefore"
	}
	return ""
}

func (nde *Node) InsertAfter(val interface{}, a ...interface{}) {
	if nde != nil && nde.lst != nil {
		var mdfying func(interface{})
		var mdfied func(bool, bool, int, *Node, interface{})
		if al := len(a); al > 0 && al <= 2 {
			for _, d := range a {
				if mdfied == nil {
					mdfied, _ = d.(func(bool, bool, int, *Node, interface{}))
				}
				if mdfying == nil {
					mdfying, _ = d.(func(interface{}))
				}
			}
		}
		internalAdd(nde.lst, mdfying, mdfied, nde, insertAfter, val)
	}
}

func (nde *Node) InsertBefore(val interface{}, a ...interface{}) {
	if nde != nil && nde.lst != nil {
		var mdfying func(interface{})
		var mdfied func(bool, bool, int, *Node, interface{})
		if al := len(a); al > 0 && al <= 2 {
			for _, d := range a {
				if mdfied == nil {
					mdfied, _ = d.(func(bool, bool, int, *Node, interface{}))
				}
				if mdfying == nil {
					mdfying, _ = d.(func(interface{}))
				}
			}
		}
		internalAdd(nde.lst, mdfying, mdfied, nde, insertBefore, val)
	}
}

func (nde *Node) Set(val interface{}) {
	if nde != nil {
		if (nde.val == nil && val != nil) || (val == nil && nde.val != nil) || (nde.val != val) {
			if nde.lst.distinct {
				delete(nde.lst.vnds, nde.val)
				nde.lst.vnds[val] = nde
			}
			nde.val = val
		}
	}
}

func (nde *Node) Dispose(eventRemoved func(nde *Node, val interface{}), disposingRemoving func(nde *Node, val interface{})) {
	if nde != nil {
		if nde.lst != nil {
			diposeNode(nde.lst, nde, eventRemoved, disposingRemoving)
			nde.lst = nil
		}
		nde.val = nil
		nde = nil
	}
}

type List struct {
	head       *Node
	tail       *Node
	vnds       map[interface{}]*Node
	reversemap map[*Node]*Node
	forwardmap map[*Node]*Node
	distinct   bool
}

func NewList(distinct ...bool) (lst *List) {
	lst = &List{head: nil, tail: nil, reversemap: map[*Node]*Node{}, forwardmap: map[*Node]*Node{}, distinct: len(distinct) == 1 && distinct[0]}
	if lst.distinct {
		lst.vnds = map[interface{}]*Node{}
	}
	return
}

func (lst *List) IsDistinct() bool {
	return lst.distinct
}

func (lst *List) newNode(val interface{}) (nde *Node) {
	nde = &Node{lst: lst, val: val}
	return
}

func (lst *List) Length() (ln int) {
	if lst != nil {
		ln = len(lst.forwardmap)
	}
	return
}

func (lst *List) ValueNode(val interface{}) (nde *Node) {
	if lst.distinct && val != nil {
		nde = lst.vnds[val]
	}
	return
}

func (lst *List) Do(RemovingNode func(*Node, interface{}) bool,
	RemovedNode func(*Node, interface{}),
	DisposingNode func(*Node, interface{})) {
	if lst.head != nil && lst.tail != nil {
		crntnde := lst.head
		nxtnde := crntnde
		for nxtnde != nil {
			if RemovingNode != nil && RemovingNode(nxtnde, nxtnde.val) {
				crntnde = nxtnde
				nxtnde = lst.forwardmap[crntnde]
				crntnde.Dispose(RemovedNode, DisposingNode)
				crntnde = nil
			} else {
				nxtnde = lst.forwardmap[crntnde]
				crntnde = nxtnde
			}
		}
	}
}

func (lst *List) DoReverse(RemovingNode func(*Node) bool,
	RemovedNode func(*Node),
	DisposingNode func(*Node)) {
}

func (lst *List) Dispose(eventRemoving func(*Node, interface{}), eventDisposing func(*Node, interface{})) {
	if lst != nil {
		if lst.forwardmap != nil && lst.reversemap != nil {
			for lst.head != nil {
				diposeNode(lst, lst.head, eventRemoving, eventDisposing)
			}
			lst.forwardmap = nil
			lst.reversemap = nil
		}
		lst = nil
	}
}

//func (lst *List) DoAdd(Adding func())

func (lst *List) Add(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), val ...interface{}) {
	lst.InsertAfter(mdfying, mdfied, lst.tail, val...)
}

func diposeNode(lst *List, nde *Node, eventRemoved func(*Node, interface{}), eventDisposing func(*Node, interface{})) {
	if nde != nil && lst != nil {
		if nde.lst == lst {
			pvrnde := lst.reversemap[nde]
			nxtnde := lst.forwardmap[nde]
			if nde == lst.head && nde == lst.tail {
				lst.head = nil
				lst.tail = nil
			} else if nde == lst.tail && nde != lst.head {
				lst.tail = pvrnde
			} else if nde == lst.head && nde != lst.tail {
				lst.head = nxtnde
			} else if nde != lst.head && nde != lst.tail {

			}
			if pvrnde != nil {
				lst.forwardmap[pvrnde] = nxtnde
			}
			if nxtnde != nil {
				lst.reversemap[nxtnde] = pvrnde
			}
			delete(lst.forwardmap, nde)
			delete(lst.reversemap, nde)
			if lst.distinct {
				delete(lst.vnds, nde.val)
			}
			if eventRemoved != nil {
				eventRemoved(nde, nde.val)
			}
		}
		if eventDisposing != nil {
			eventDisposing(nde, nde.val)
		}
	}
}

func internalAdd(lst *List, modifying func(value interface{}), modified func(changed bool, validval bool, vindex int, node *Node, value interface{}), ndefrm *Node, action listaction, val ...interface{}) {
	if len(val) > 0 {
		var inbfadd = func(nde *Node, prvnde *Node, nxtNde *Node, nvali int, nval interface{}) {
			if nwndo := lst.newNode(nval); nwndo != nil {
				if nde == nil {
					if lst.head == nil && lst.tail == nil {
						lst.head = nwndo
						lst.tail = nwndo
						lst.forwardmap[nwndo] = nil
						lst.reversemap[nwndo] = nil
					} else if lst.head != nil && lst.tail != nil {
						lst.forwardmap[lst.tail] = nwndo
						lst.reversemap[nwndo] = lst.tail
						lst.forwardmap[nwndo] = nil
						lst.tail = nwndo
					}
				} else if action == insertAfter {
					if nde == lst.tail {
						lst.forwardmap[lst.tail] = nwndo
						lst.reversemap[nwndo] = lst.tail
						lst.forwardmap[nwndo] = nil
						lst.tail = nwndo
					} else if prvnde != nil && nxtNde != nil {
						lst.forwardmap[lst.tail] = nwndo
						lst.reversemap[nwndo] = lst.tail
						lst.forwardmap[nwndo] = nil
					}
				} else if action == insertBefore {

				}
				if lst.distinct {
					lst.vnds[nval] = nwndo
				}
				if modified != nil {
					modified(true, true, nvali, nwndo, nval)
				}
				ndefrm = nwndo
			}
		}

		for vli, vl := range val {
			if lst.distinct {
				if vl == nil {
					if modified != nil {
						modified(false, false, vli, nil, val)
					}
					continue
				}
				if vlnde, vlok := lst.vnds[vl]; vlok {
					modified(false, true, vli, vlnde, val)
					continue
				}
			}
			inbfadd(ndefrm, lst.reversemap[ndefrm], lst.forwardmap[ndefrm], vli, vl)
		}
	}
}

func (lst *List) InsertBefore(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), ndefrm *Node, val ...interface{}) {
	internalAdd(lst, mdfying, mdfied, ndefrm, insertBefore, val...)
}

func (lst *List) InsertAfter(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), ndefrm *Node, val ...interface{}) {
	internalAdd(lst, mdfying, mdfied, ndefrm, insertAfter, val...)
}
