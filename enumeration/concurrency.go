package enumeration

type Pool struct {
	pool    chan interface{}
	New     func() interface{}
	Dispose func(interface{})
}

// Borrow a Client from the pool.
func (p *Pool) Borrow() (c interface{}) {
	select {
	case c = <-p.pool:
	default:
		if p.New != nil {
			c = p.New()
		}
	}
	return c
}

// Put returns a Client to the pool.
func (p *Pool) Put(c interface{}) (returned bool) {
	select {
	case p.pool <- c:
		returned = true
	default:
		if p.Dispose != nil {
			p.Dispose(c)
		}
	}
	return
}

// Put returns a Client to the pool.
func (p *Pool) Close() {
	if p != nil {

	}
}

func NewPool(max int, New func() interface{}, Dispose func(interface{})) (pool *Pool) {
	pool = &Pool{New: New, Dispose: Dispose, pool: make(chan interface{})}
	return
}
