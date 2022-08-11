package ds

type Less[T any] func(l, r T) bool

type Heap[T any] interface {
	GetMin() (T, bool)
	Pop() (T, bool)
	Add(...T)
}

type Tree[T any] interface {
}
