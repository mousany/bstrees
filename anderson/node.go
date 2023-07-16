package anderson

import "golang.org/x/exp/constraints"

type andersonTreeNode[T constraints.Ordered] struct {
	value T
	left  *andersonTreeNode[T]
	right *andersonTreeNode[T]
	size  uint
	level uint
}

func newAndersonTreeNode[T constraints.Ordered](value T, level uint) *andersonTreeNode[T] {
	return &andersonTreeNode[T]{
		value: value,
		left:  nil,
		right: nil,
		size:  1,
		level: level,
	}
}

func (n *andersonTreeNode[T]) update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}
