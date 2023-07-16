package treap

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

type treapNode[T constraints.Ordered] struct {
	value  T
	left   *treapNode[T]
	right  *treapNode[T]
	weight uint // Random weight
	size   uint // Size of subtree, unnecessary if you don't need kth element
}

func newTreapNode[T constraints.Ordered](value T) *treapNode[T] {
	return &treapNode[T]{value: value, left: nil, right: nil, weight: uint(rand.Uint32()), size: 1}
}

func (n *treapNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}
