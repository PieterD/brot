package fiboheap

import (
	"github.com/PieterD/brot/pkg/ds"
)

type Heap[T any] struct {
	less       ds.Less[T]
	alloc      *niceAllocator[T]
	part       *degreePartitioner[T]
	fixerLimit uint
}

func NewHeap[T any](less ds.Less[T], alloc Allocator[T], fixerLimit uint) *Heap[T] {
	h := &Heap[T]{
		less:       less,
		alloc:      newNiceAllocator[T](less, alloc),
		part:       newDegreePartitioner[T](less),
		fixerLimit: fixerLimit,
	}
	return h
}

func (h *Heap[T]) Pop() (T, bool) {
	var z T
	n := h.part.PopMin()
	if n == nil {
		return z, false
	}
	v := n.Value
	for _, child := range n.Children {
		h.addNode(child, h.fixerLimit)
	}
	return v, true
}

func (h *Heap[T]) GetMin() (T, bool) {
	return h.part.GetMin()
}

func (h *Heap[T]) Add(values ...T) {
	if len(values) == 0 {
		return
	}
	more := func() bool {
		return len(values) > 0
	}
	pop := func() T {
		v := values[0]
		values = values[1:]
		return v
	}
	for len(values) > 0 {
		v := pop()
		if combined := h.combineWithD0(v); combined != nil {
			h.addNode(combined, h.fixerLimit)
			continue
		}
		if !more() {
			h.addNode(h.alloc.NewNode(v), h.fixerLimit)
			break
		}
		v2 := pop()
		h.addNode(h.alloc.NewTuple(v, v2), h.fixerLimit)
	}
}

func (h *Heap[T]) addNode(n *Node[T], fixerLimit uint) {
	n2 := h.part.Pop(false, n.Degree())
	if n2 == nil || fixerLimit == 0 {
		h.part.Add(n)
		return
	}
	combined := h.alloc.CombineNode(n, n2)
	h.addNode(combined, fixerLimit-1)
}

func (h *Heap[T]) combineWithD0(v T) *Node[T] {
	n := h.part.Pop(false, 0)
	if n == nil {
		return nil
	}
	if h.less(v, n.Value) {
		newN := h.alloc.NewNode(v, n)
		return newN
	}
	newN := h.alloc.NewNode(n.Value, h.alloc.NewNode(v))
	h.alloc.FreeNode(n)
	return newN
}
