package enumeration

type Node struct {
	lst *List
	val interface{}
}

func (nde *Node) Value() interface{} {
	return nde.val
}

func (nde *Node) Next() *Node {
	if nde != nil && nde.lst != nil {
		return nextNode(nde.lst, nde)
	}
	return nil
}

func (nde *Node) Previous() *Node {
	if nde != nil && nde.lst != nil {
		return previousNode(nde.lst, nde)
	}
	return nil
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
		internalInput(nde.lst, mdfying, mdfied, nde, insertAfter, val)
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
		internalInput(nde.lst, mdfying, mdfied, nde, insertBefore, val)
	}
}

func (nde *Node) Set(val interface{}, forceset ...bool) {
	if nde != nil {
		if (len(forceset) == 1 && forceset[0]) || ((nde.val == nil && val != nil) || (val == nil && nde.val != nil) || (nde.val != val)) {
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
			disposeNode(nde.lst, nde, eventRemoved, disposingRemoving)
			nde.lst = nil
		}
		nde.val = nil
		nde = nil
	}
}

type List struct {
	doingnde   *Node
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

func previousNode(lst *List, nde *Node) (prvnde *Node) {
	func() {
		if lst != nil && len(lst.reversemap) > 0 {
			prvnde = lst.reversemap[nde]
		}
	}()
	return
}

func nextNode(lst *List, nde *Node) (nxtnde *Node) {
	func() {
		if lst != nil && len(lst.forwardmap) > 0 {
			nxtnde = lst.forwardmap[nde]
		}
	}()
	return
}

func (lst *List) CurrentDoing() *Node {
	return lst.doingnde
}

func (lst *List) Head() *Node {
	return lst.head
}

func (lst *List) Tail() *Node {
	return lst.tail
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

func (lst *List) Iterate(dolnk func(*Node, interface{}) (bool, error), errdolnk func(*Node, interface{}, error) bool, donelnk func(*Node) error, errdonelnk func(*Node, error) bool, disposelnk func(nde *Node), eventRemoved func(nde *Node, val interface{}), disposingRemoving func(nde *Node, val interface{})) {
	if lst != nil && lst.head != nil && lst.tail != nil {
		var done = false
		var err error = nil
		func() {
			for e := lst.Head(); e != nil && !done; done, e = false, e.Next() {
				if e == nil {
					continue
				}
				func() {
					if dolnk != nil {
						if done, err = dolnk(e, e.Value()); done || err != nil {
							if done && err == nil {
								if donelnk != nil {
									if err = donelnk(e); err != nil {
										if errdonelnk != nil {
											done = errdonelnk(e, err)
										}
									}
								}
							} else if !done && errdonelnk != nil {
								done = errdonelnk(e, err)
							}
						} else if err != nil {
							if errdolnk != nil {
								if done = errdolnk(e, e.Value(), err); done {
									if donelnk != nil {
										if err = donelnk(e); err != nil {
											if errdonelnk != nil {
												done = errdonelnk(e, err)
											}
										}
									}
								}
							}
						}
						if done {
							if disposelnk == nil {
								e.Dispose(eventRemoved, disposingRemoving)
							} else {
								disposelnk(e)
							}
						}
					}
				}()
			}
		}()
	}
}

func (lst *List) Do(RemovingNode func(*Node, interface{}) (bool, error),
	ErrRemovingNode func(*Node, interface{}, error) bool,
	RemovedNode func(*Node, interface{}),
	DisposingNode func(*Node, interface{})) {
	if lst.head != nil && lst.tail != nil {
		lst.doingnde = lst.head
		nxtnde := lst.doingnde
		for nxtnde != nil {
			if RemovingNode != nil {
				dne, err := RemovingNode(nxtnde, nxtnde.val)
				if err == nil {
					if dne {
						lst.doingnde = nxtnde
						nxtnde = lst.forwardmap[lst.doingnde]
						lst.doingnde.Dispose(RemovedNode, DisposingNode)
						lst.doingnde = nil
					} else {
						nxtnde = lst.forwardmap[lst.doingnde]
						lst.doingnde = nxtnde
					}
				} else {
					if ErrRemovingNode == nil {
						lst.doingnde = nxtnde
						nxtnde = lst.forwardmap[lst.doingnde]
						lst.doingnde.Dispose(RemovedNode, DisposingNode)
						lst.doingnde = nil
					} else if ErrRemovingNode != nil {
						if dne = ErrRemovingNode(nxtnde, nxtnde.val, err); dne {
							lst.doingnde = nxtnde
							nxtnde = lst.forwardmap[lst.doingnde]
							lst.doingnde.Dispose(RemovedNode, DisposingNode)
							lst.doingnde = nil
						} else {
							nxtnde = lst.forwardmap[lst.doingnde]
							lst.doingnde = nxtnde
						}
					}
				}
			} else {
				nxtnde = lst.forwardmap[lst.doingnde]
				lst.doingnde = nxtnde
			}
		}
	}
}

func (lst *List) DoReverse(RemovingNode func(*Node, interface{}) bool,
	RemovedNode func(*Node, interface{}),
	DisposingNode func(*Node, interface{})) {
	if lst.head != nil && lst.tail != nil {
		crntnde := lst.tail
		prvnde := crntnde
		for prvnde != nil {
			if RemovingNode != nil && RemovingNode(prvnde, prvnde.val) {
				crntnde = prvnde
				prvnde = lst.reversemap[crntnde]
				crntnde.Dispose(RemovedNode, DisposingNode)
				crntnde = nil
			} else {
				prvnde = lst.reversemap[crntnde]
				crntnde = prvnde
			}
		}
	}
}

func (lst *List) Dispose(eventRemoving func(*Node, interface{}), eventDisposing func(*Node, interface{})) {
	if lst != nil {
		if lst.forwardmap != nil && lst.reversemap != nil {
			for lst.head != nil {
				disposeNode(lst, lst.head, eventRemoving, eventDisposing)
			}
			lst.forwardmap = nil
			lst.reversemap = nil
		}
		lst = nil
	}
}

func (lst *List) Push(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), val ...interface{}) int {
	lst.InsertAfter(mdfying, mdfied, lst.tail, val...)
	return lst.Length()
}

func (lst *List) Shift(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), val ...interface{}) int {
	lst.InsertBefore(mdfying, mdfied, lst.head, val...)
	return lst.Length()
}

func (lst *List) Pop() (val interface{}) {
	if lst.tail != nil {
		lst.tail.Dispose(nil, func(nde *Node, v interface{}) {
			val = v
		})
	}
	return
}

func (lst *List) Unshift() (val interface{}) {
	if lst.head != nil {
		lst.head.Dispose(nil, func(nde *Node, v interface{}) {
			val = v
		})
	}
	return
}

func disposeNode(lst *List, nde *Node, eventRemoved func(*Node, interface{}), eventDisposing func(*Node, interface{})) {
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

func internalInput(lst *List, modifying func(value interface{}), modified func(changed bool, validval bool, vindex int, node *Node, value interface{}), ndefrm *Node, action listaction, val ...interface{}) {
	if len(val) > 0 {
		var inbfinput = func(nde *Node, prvnde *Node, nxtNde *Node, nvali int, nval interface{}) {
			if nwndo := lst.newNode(nval); nwndo != nil {
				if nde == nil {
					if lst.head == nil && lst.tail == nil {
						lst.head = nwndo
						lst.tail = nwndo
						lst.forwardmap[nwndo] = nil
						lst.reversemap[nwndo] = nil
					} else if action == insertAfter && lst.head != nil && lst.tail != nil {
						lst.forwardmap[lst.tail] = nwndo
						lst.reversemap[nwndo] = lst.tail
						lst.forwardmap[nwndo] = nil
						lst.tail = nwndo
					} else if action == insertBefore && lst.head != nil && lst.tail != nil {
						lst.reversemap[lst.head] = nwndo
						lst.forwardmap[nwndo] = lst.head
						lst.reversemap[nwndo] = nil
						lst.head = nwndo
					}
				} else if action == insertAfter {
					if nde == lst.tail {
						lst.forwardmap[lst.tail] = nwndo
						lst.reversemap[nwndo] = lst.tail
						lst.forwardmap[nwndo] = nil
						lst.tail = nwndo
					} else if prvnde != nil && nxtNde != nil {
						lst.forwardmap[nde] = nwndo
						lst.reversemap[nwndo] = nde
						lst.forwardmap[nwndo] = nxtNde
						lst.reversemap[nxtNde] = nwndo
					} else if nde == lst.head {
						lst.forwardmap[lst.head] = nwndo
						lst.reversemap[nxtNde] = nwndo
						lst.forwardmap[nwndo] = nxtNde
						//lst.tail = nwndo
					}
				} else if action == insertBefore {
					if nde == lst.head {
						lst.reversemap[lst.head] = nwndo
						lst.forwardmap[nwndo] = lst.head
						lst.reversemap[nwndo] = nil
						lst.head = nwndo
					} else if prvnde != nil && nxtNde != nil {
						lst.reversemap[nde] = nwndo
						lst.forwardmap[nwndo] = nde
						lst.reversemap[nwndo] = prvnde
						lst.forwardmap[prvnde] = nwndo
					} else if nde == lst.tail {
						lst.reversemap[lst.tail] = nwndo
						lst.forwardmap[prvnde] = nwndo
						lst.reversemap[nwndo] = prvnde
						//lst.tail = nwndo
					}
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
			inbfinput(ndefrm, lst.reversemap[ndefrm], lst.forwardmap[ndefrm], vli, vl)
		}
	}
}

func (lst *List) InsertBefore(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), ndefrm *Node, val ...interface{}) {
	internalInput(lst, mdfying, mdfied, ndefrm, insertBefore, val...)
}

func (lst *List) InsertAfter(mdfying func(interface{}), mdfied func(bool, bool, int, *Node, interface{}), ndefrm *Node, val ...interface{}) {
	internalInput(lst, mdfying, mdfied, ndefrm, insertAfter, val...)
}
