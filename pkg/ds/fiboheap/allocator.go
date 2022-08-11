package fiboheap

import "github.com/PieterD/brot/pkg/ds"

type Allocator[T any] interface {
	NewNode(degree uint) *Node[T]
	FreeNode(n *Node[T])
}

type niceAllocator[T any] struct {
	less ds.Less[T]
	raw  Allocator[T]
}

func newNiceAllocator[T any](less ds.Less[T], raw Allocator[T]) *niceAllocator[T] {
	return &niceAllocator[T]{
		less: less,
		raw:  raw,
	}
}

func (a *niceAllocator[T]) FreeNode(n *Node[T]) {
	a.raw.FreeNode(n)
}

func (a *niceAllocator[T]) NewNode(v T, children ...*Node[T]) *Node[T] {
	n := a.raw.NewNode(uint(len(children)))
	n.Value = v
	copy(n.Children, children)
	return n
}

func (a *niceAllocator[T]) NewTuple(v1, v2 T) *Node[T] {
	if a.less(v2, v1) {
		v2, v1 = v1, v2
	}
	return a.NewNode(v1, a.NewNode(v2))
}

func (a *niceAllocator[T]) CombineNode(n *Node[T], newChild *Node[T]) *Node[T] {
	if n.Degree() != newChild.Degree() {
		panic("degree mismatch")
	}
	if a.less(newChild.Value, n.Value) {
		newChild, n = n, newChild
	}
	newN := a.raw.NewNode(n.Degree() + 1)
	newN.Value = n.Value
	copy(newN.Children, n.Children)
	newN.Children[len(newN.Children)-1] = newChild
	a.raw.FreeNode(n)
	return newN
}
