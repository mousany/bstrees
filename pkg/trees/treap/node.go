package treap

import (
	"bstrees/internal/order"
	"math/rand"
)

type treapTreeNode[T order.Ordered] struct {
	value  T
	left   *treapTreeNode[T]
	right  *treapTreeNode[T]
	weight uint // Random weight
	size   uint // Size of subtree, unnecessary if you don't need kth element
}

func newTreapTreeNode[T order.Ordered](value T) *treapTreeNode[T] {
	return &treapTreeNode[T]{value: value, left: nil, right: nil, weight: uint(rand.Uint32()), size: 1}
}

func (n *treapTreeNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}
