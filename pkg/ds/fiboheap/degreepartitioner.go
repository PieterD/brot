package fiboheap

import "github.com/PieterD/brot/pkg/ds"

type degreePartitioner[T any] struct {
	less      ds.Less[T]
	degrees   []*rootCollection[T]
	minDegree int
}

func newDegreePartitioner[T any](less ds.Less[T]) *degreePartitioner[T] {
	return &degreePartitioner[T]{
		less:      less,
		minDegree: -1,
	}
}

func (p *degreePartitioner[T]) mustRootCollection(degree uint) *rootCollection[T] {
	for uint(len(p.degrees)) <= degree {
		cur := uint(len(p.degrees))
		p.degrees = append(p.degrees, newRootCollection[T](p.less, cur))
	}
	return p.degrees[degree]
}

func (p *degreePartitioner[T]) Add(n *Node[T]) {
	degree := n.Degree()
	rootColl := p.mustRootCollection(degree)
	updateMin := rootColl.Push(n)
	if updateMin {
		p.updateMin(degree)
	}
}

func (p *degreePartitioner[T]) Pop(includeMin bool, degree uint) *Node[T] {
	rootColl := p.mustRootCollection(degree)
	n := rootColl.PopHigh()
	if n != nil {
		return n
	}
	if !includeMin {
		return nil
	}
	n = p.PopMin()
	return n
}

func (p *degreePartitioner[T]) PopMin() *Node[T] {
	if p.minDegree == -1 {
		return nil
	}
	rootColl := p.mustRootCollection(uint(p.minDegree))
	n := rootColl.PopMin()
	if n == nil {
		panic("minDegree collection PopMin returned nil")
	}
	p.minDegree = p.findMinDegree()
	return n
}

func (p *degreePartitioner[T]) GetMin() (T, bool) {
	var z T
	if p.minDegree == -1 {
		return z, false
	}
	return p.degrees[p.minDegree].GetMin()
}

func (p *degreePartitioner[T]) findMinDegree() int {
	var min T
	minDegree := -1
	for degree, rootColl := range p.degrees {
		v, ok := rootColl.GetMin()
		if !ok {
			continue
		}
		if minDegree == -1 {
			min = v
			minDegree = degree
			continue
		}
		if p.less(v, min) {
			v, min = min, v
			minDegree = degree
		}
	}
	return minDegree
}

func (p *degreePartitioner[T]) updateMin(to uint) {
	if p.minDegree < 0 {
		_, newExists := p.degrees[to].GetMin()
		if !newExists {
			return
		}
		p.minDegree = int(to)
		return
	}
	if uint(p.minDegree) == to {
		return
	}
	curValue, curExists := p.degrees[p.minDegree].GetMin()
	if !curExists {
		p.minDegree = -1
	}
	newValue, newExists := p.degrees[to].GetMin()
	if !newExists {
		return
	}
	if curExists {
		if p.less(newValue, curValue) {
			p.minDegree = int(to)
		}
		return
	}
	// newExists && !curExists
	p.minDegree = int(to)
}
