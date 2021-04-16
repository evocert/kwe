package enumeration

type Node struct {
	lst *List
	val interface{}
}

type listaction int

const (
	insertAfter listaction = iota
	insertBefore
)

func (nde *Node) InsertAfter(val interface{}) {
	if nde != nil && nde.lst != nil {
		internalAdd(nde.lst, nde, insertAfter, val)
	}
}

func (nde *Node) InsertBefore(val interface{}) {
	if nde != nil && nde.lst != nil {
		internalAdd(nde.lst, nde, insertBefore, val)
	}
}

func (nde *Node) Set(val interface{}) {
	if nde != nil {
		if nde.val != val {
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
	reversemap map[*Node]*Node
	forwardmap map[*Node]*Node
}

func NewList() (lst *List) {
	lst = &List{head: nil, tail: nil, reversemap: map[*Node]*Node{}, forwardmap: map[*Node]*Node{}}
	return
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

func (lst *List) Do(RemovingNode func(*Node, interface{}) bool,
	RemovedNode func(*Node, interface{}),
	DisposingNode func(*Node, interface{})) {
	if lst.head != nil && lst.tail != nil {

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
				}
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

func (lst *List) Add(val ...interface{}) {
	lst.InsertAfter(nil, val...)
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
			if eventRemoved != nil {
				eventRemoved(nde, nde.val)
			}
		}
		if eventDisposing != nil {
			eventDisposing(nde, nde.val)
		}
	}
}

func internalAdd(lst *List, ndefrm *Node, action listaction, val ...interface{}) {
	if len(val) > 0 {

		var inbfadd = func(nde *Node, prvnde *Node, nxtNde *Node, nval interface{}) {
			if nwndo := lst.newNode(nval); nwndo != nil {
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
				ndefrm = nwndo
			}
		}

		if ndefrm == nil {
			for _, vl := range val {
				inbfadd(ndefrm, lst.reversemap[ndefrm], lst.forwardmap[ndefrm], vl)
			}
		}
	}
}

func (lst *List) InsertBefore(ndefrm *Node, val ...interface{}) {
	internalAdd(lst, ndefrm, insertBefore, val...)
}

func (lst *List) InsertAfter(ndefrm *Node, val ...interface{}) {
	internalAdd(lst, ndefrm, insertAfter, val...)
}
