package fhq

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

type fhqTreapNode[T constraints.Ordered] struct {
	value  T
	left   *fhqTreapNode[T]
	right  *fhqTreapNode[T]
	weight uint // Random weight
	size   uint // Size of subtree, unnecessary if you don't need kth element
}

func newFHQTreapNode[T constraints.Ordered](value T) *fhqTreapNode[T] {
	return &fhqTreapNode[T]{value: value, left: nil, right: nil, weight: uint(rand.Uint32()), size: 1}
}

func (n *fhqTreapNode[T]) Update() {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}
