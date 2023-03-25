package anderson

import "bstrees/internal/order"

type andersonTreeNode[T order.Ordered] struct {
	value T
	left  *andersonTreeNode[T]
	right *andersonTreeNode[T]
	size  uint
	level uint
}

func newAndersonTreeNode[T order.Ordered](value T, level uint) *andersonTreeNode[T] {
	return &andersonTreeNode[T]{
		value: value,
		left:  nil,
		right: nil,
		size:  1,
		level: level,
	}
}

func (n *andersonTreeNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func (n *andersonTreeNode[T]) Leaf() bool {
	return n.left == nil && n.right == nil
}

func (n *andersonTreeNode[T]) Full() bool {
	return n.left != nil && n.right != nil
}

func (n *andersonTreeNode[T]) Child(direction bool) *andersonTreeNode[T] {
	if direction {
		return n.right
	}
	return n.left
}

func (n *andersonTreeNode[T]) SetChild(child *andersonTreeNode[T], direction bool) {
	if direction {
		n.right = child
	} else {
		n.left = child
	}
}
