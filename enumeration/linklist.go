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

func (nde *Node) Dispose() {
	if nde != nil {
		if nde.lst != nil {
			diposeNode(nde.lst, nde)
		}
	}
}

type List struct {
	head          *Node
	tail          *Node
	reversemap    map[*Node]*Node
	forwardmap    map[*Node]*Node
	RemovingNode  func(*Node) bool
	RemovedNode   func(*Node)
	DisposingNode func(*Node)
}

func NewList() (lst *List) {
	lst = &List{}
	return
}

func (lst *List) newNode(val interface{}) (nde *Node) {
	nde = &Node{lst: lst, val: val}
	return
}

func (lst *List) Add(val ...interface{}) {
	lst.InsertAfter(nil, val...)
}

func diposeNode(lst *List, nde *Node) {
	if nde != nil && lst != nil {
		if nde.lst == lst {
			canRemove := true
			if lst.RemovingNode != nil {
				canRemove = lst.RemovingNode(nde)
			}
			if canRemove {
				pvrnde := lst.reversemap[nde]
				nxtnde := lst.forwardmap[nde]

				if nde == lst.head && nde == lst.head {

				} else if nde == lst.tail && nde != lst.head {

				} else if nde == lst.head && nde != lst.tail {

				}
				if pvrnde != nil {
					lst.forwardmap[pvrnde] = nxtnde
				}
				if nxtnde != nil {
					lst.reversemap[nxtnde] = pvrnde
				}
				delete(lst.forwardmap, nde)
				delete(lst.reversemap, nde)
				if lst.RemovedNode != nil {
					lst.RemovingNode(nde)
				}
			}
		}
		if lst.DisposingNode != nil {
			lst.DisposingNode(nde)
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
