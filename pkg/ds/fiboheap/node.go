package fiboheap

type Node[T any] struct {
	Value     T
	MoreRoots *Node[T]
	Children  []*Node[T]
}

func (n *Node[T]) Degree() uint {
	return uint(len(n.Children))
}

func (n *Node[T]) Clean() {
	children := n.Children
	*n = Node[T]{}
	for i := range children {
		children[i] = nil
	}
	n.Children = children
}
