package fiboheap

import "github.com/PieterD/brot/pkg/ds"

type rootCollection[T any] struct {
	degree    uint
	less      ds.Less[T]
	rootCount uint
	minRoot   *Node[T]
	moreRoots *Node[T]
}

func newRootCollection[T any](less ds.Less[T], degree uint) *rootCollection[T] {
	return &rootCollection[T]{
		degree: degree,
		less:   less,
	}
}

func (c *rootCollection[T]) GetMin() (T, bool) {
	var z T
	if c.minRoot == nil {
		return z, false
	}
	return c.minRoot.Value, true
}

func (c *rootCollection[T]) Push(n *Node[T]) (updateMin bool) {
	c.rootCount++
	if c.minRoot == nil {
		c.minRoot = n
		return true
	}
	if c.less(n.Value, c.minRoot.Value) {
		n, c.minRoot = c.minRoot, n
		updateMin = true
	}
	n.MoreRoots = c.moreRoots
	c.moreRoots = n
	return updateMin
}

func (c *rootCollection[T]) PopHigh() *Node[T] {
	n := c.moreRoots
	if n == nil {
		return nil
	}
	c.rootCount--
	c.moreRoots = n.MoreRoots
	n.MoreRoots = nil
	return n
}

func (c *rootCollection[T]) PopMin() *Node[T] {
	if c.minRoot == nil {
		return nil
	}
	c.rootCount--
	min := c.minRoot
	c.minRoot = nil
	c.minRoot = c.PopHigh()
	if c.minRoot == nil {
		return min
	}
	var newRoot *Node[T]
	for {
		high := c.PopHigh()
		if high == nil {
			break
		}
		if c.less(high.Value, c.minRoot.Value) {
			high, c.minRoot = c.minRoot, high
		}
		high.MoreRoots = newRoot
		newRoot = high
	}
	c.moreRoots = newRoot
	return min
}
