package splay

import "bstrees/internal/order"

type splayTreeNode[T order.Ordered] struct {
	value  T
	left   *splayTreeNode[T]
	right  *splayTreeNode[T]
	parent *splayTreeNode[T]
	size   uint
	rec    uint // This field is Splay only
	// Because Splay operation will scatter nodes with the same value
	// While traditional BST search mechanics is too slow on Splay
}

func newSplayTreeNode[T order.Ordered](value T) *splayTreeNode[T] {
	return &splayTreeNode[T]{
		value:  value,
		left:   nil,
		right:  nil,
		parent: nil,
		size:   1,
		rec:    1,
	}
}

func (n *splayTreeNode[T]) update() {
	n.size = n.rec
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func (n *splayTreeNode[T]) setChild(child *splayTreeNode[T], direction bool) {
	if direction {
		n.right = child
	} else {
		n.left = child
	}
	if child != nil {
		child.parent = n
	}
}

func (n *splayTreeNode[T]) child(direction bool) *splayTreeNode[T] {
	if direction {
		return n.right
	}
	return n.left
}
