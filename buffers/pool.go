package buffers

import (
	"bytes"
	"container/list"
	"time"
)

type buf struct {
	when time.Time
	buf  *bytes.Buffer
}

type Pool struct {
	Get  chan *bytes.Buffer
	Give chan *bytes.Buffer

	q       *list.List
	timeout time.Duration
}

func NewPool(timeout time.Duration) *Pool {
	p := &Pool{
		Get:  make(chan *bytes.Buffer),
		Give: make(chan *bytes.Buffer),

		q:       new(list.List),
		timeout: timeout,
	}
	go p.run()
	return p
}

func (p *Pool) run() {
	for {
		if p.q.Len() == 0 {
			p.q.PushFront(buf{when: time.Now(), buf: bytes.NewBuffer(nil)})
		}

		e := p.q.Front()
		timeout := time.NewTimer(p.timeout)
		select {
		case b := <-p.Give:
			p.q.PushFront(buf{when: time.Now(), buf: b})

		case p.Get <- e.Value.(buf).buf:
			timeout.Stop()
			p.q.Remove(e)

		case <-timeout.C:
			e := p.q.Front()
			for e != nil {
				n := e.Next()
				if time.Since(e.Value.(buf).when) > p.timeout {
					p.q.Remove(e)
					e.Value = nil
				}
				e = n
			}
		}
	}
}
