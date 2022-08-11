package poolalloc

import (
	"github.com/PieterD/brot/pkg/ds/fiboheap"
	"sync"
)

type PoolAlloc[T any] struct {
	poolsByDegree []*sync.Pool
}

func newPool[T any](degree uint) *sync.Pool {
	return &sync.Pool{New: func() any {
		var children []*fiboheap.Node[T]
		if degree > 0 {
			children = make([]*fiboheap.Node[T], degree, degree)
		}
		return &fiboheap.Node[T]{
			Children: children,
		}
	}}
}

func New[T any]() *PoolAlloc[T] {
	return &PoolAlloc[T]{
		poolsByDegree: []*sync.Pool{
			newPool[T](0),
			newPool[T](1),
			newPool[T](2),
			newPool[T](3),
		},
	}
}

func (p *PoolAlloc[T]) mustPool(degree uint) {
	size := uint(len(p.poolsByDegree))
	if degree < size {
		return
	}
	for uint(len(p.poolsByDegree)) < degree {
		p.poolsByDegree = append(p.poolsByDegree, newPool[T](uint(len(p.poolsByDegree))))
	}
}

func (p *PoolAlloc[T]) NewNode(degree uint) *fiboheap.Node[T] {
	p.mustPool(degree)
	n := p.poolsByDegree[degree].Get().(*fiboheap.Node[T])
	n.Clean()
	return n
}

func (p *PoolAlloc[T]) FreeNode(n *fiboheap.Node[T]) {
	degree := n.Degree()
	p.mustPool(degree)
	n.Clean()
	p.poolsByDegree[degree].Put(n)
}

var _ fiboheap.Allocator[int] = &PoolAlloc[int]{}
