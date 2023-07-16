package splay

import "golang.org/x/exp/constraints"

type splayNode[T constraints.Ordered] struct {
	value  T
	left   *splayNode[T]
	right  *splayNode[T]
	parent *splayNode[T]
	size   uint
	rec    uint // This field is Splay only
	// Because Splay operation will scatter nodes with the same value
	// While traditional BST search mechanics is too slow on Splay
}

func newSplayNode[T constraints.Ordered](value T) *splayNode[T] {
	return &splayNode[T]{
		value:  value,
		left:   nil,
		right:  nil,
		parent: nil,
		size:   1,
		rec:    1,
	}
}

func (n *splayNode[T]) update() {
	n.size = n.rec
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func (n *splayNode[T]) setChild(child *splayNode[T], direction bool) {
	if direction {
		n.right = child
	} else {
		n.left = child
	}
	if child != nil {
		child.parent = n
	}
}
