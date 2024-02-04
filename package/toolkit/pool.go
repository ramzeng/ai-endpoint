package toolkit

import "sync"

type BytesBufferPool struct {
	defaultSize int
	pool        sync.Pool
}

func (p *BytesBufferPool) Get() []byte {
	return p.pool.Get().([]byte)
}

func (p *BytesBufferPool) Put(bs []byte) {
	p.pool.Put(bs)
}

func NewBytesBufferPool(defaultSize int) *BytesBufferPool {
	p := &BytesBufferPool{
		defaultSize: defaultSize,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, defaultSize)
			},
		},
	}
	return p
}
