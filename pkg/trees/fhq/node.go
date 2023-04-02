package fhq

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
	"math/rand"
)

type fhqTreeNode[T order.Ordered] struct {
	value  T
	left   *fhqTreeNode[T]
	right  *fhqTreeNode[T]
	weight uint // Random weight
	size   uint // Size of subtree, unnecessary if you don't need kth element
}

func newFHQTreeNode[T order.Ordered](value T) *fhqTreeNode[T] {
	return &fhqTreeNode[T]{value: value, left: nil, right: nil, weight: uint(rand.Uint32()), size: 1}
}

func (n *fhqTreeNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}
