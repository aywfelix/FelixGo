package utils

import (
	"sync"

	. "github.com/aywfelix/felixgo/container/list"
)

type IObjPool interface {
	GetObj() interface{}
	FreeObj(obj interface{})
	UseCount() int
	Count() int
	Clear()
}

type ObjPool struct {
	size  int
	count int
	deque *Deque
	New   func() interface{}
	sync.RWMutex
	incrStep int
}

const incrStep = 15

func NewObjPool(newFunc func() interface{}, defaultCount int, incrStep int) *ObjPool {
	pool := &ObjPool{
		size:     0,
		count:    defaultCount,
		deque:    NewDeque(),
		New:      newFunc,
		incrStep: incrStep,
	}

	for i := 0; i < defaultCount; i++ {
		o := pool.New()
		pool.deque.PushBack(o)
	}
	return pool
}

func (p *ObjPool) expand() {
	if p.New == nil {
		return
	}
	p.count += incrStep
	for i := 0; i < incrStep; i++ {
		o := p.New()
		p.deque.PushBack(o)
	}
}

func (p *ObjPool) GetObj() interface{} {
	if p.deque.Empty() {
		p.Lock()
		p.expand()
		defer p.Unlock()
	}
	return p.deque.PopFront()
}

func (p *ObjPool) FreeObj(obj interface{}) {
	if obj != nil {
		p.deque.PushBack(obj)
	}
}

func (p *ObjPool) UseCount() int {
	p.Lock()
	defer p.Unlock()
	return p.count - p.deque.Size()
}

func (p *ObjPool) Count() int {
	p.RLock()
	defer p.RUnlock()
	return p.count
}

func (p *ObjPool) Clear() {
	p.New = nil
	p.deque = nil
}
