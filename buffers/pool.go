package buffers

import (
	"bytes"
	"sync"
)

type Pool struct {
	m   sync.Mutex
	buf chan *bytes.Buffer
}

func NewPool(size int) *Pool {
	p := &Pool{
		buf: make(chan *bytes.Buffer, size),
	}
	return p
}

func (p *Pool) Get() *bytes.Buffer {
	select {
	case b := <-p.buf:
		b.Reset()
		return b
	default:
		p.m.Lock()
		b := bytes.NewBuffer(nil)
		p.m.Unlock()
		return b
	}
}

func (p *Pool) Return(b *bytes.Buffer) {
	select {
	case p.buf <- b:
	default:
	}
}
