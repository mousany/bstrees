package avl

import (
	"math"

	"golang.org/x/exp/constraints"
)

type avlTreeNode[T constraints.Ordered] struct {
	value  T
	left   *avlTreeNode[T]
	right  *avlTreeNode[T]
	height int  // Height of the node
	size   uint // Size of subtree, unnecessary if you don't need kth element
}

func newAVLTreeNode[T constraints.Ordered](value T) *avlTreeNode[T] {
	return &avlTreeNode[T]{value: value, left: nil, right: nil, height: 0, size: 1}
}

func (n *avlTreeNode[T]) update() {
	n.height = 0
	n.size = 1
	if n.left != nil {
		n.height = int(math.Max(float64(n.height), float64(n.left.height+1)))
		n.size += n.left.size
	}
	if n.right != nil {
		n.height = int(math.Max(float64(n.height), float64(n.right.height+1)))
		n.size += n.right.size
	}
}
