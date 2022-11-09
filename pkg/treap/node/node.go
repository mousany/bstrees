package node

import (
	"bstrees/pkg/trait/ordered"
	"math/rand"
)

type TreapNode[T ordered.Ordered] struct {
	Value  T
	Left   *TreapNode[T]
	Right  *TreapNode[T]
	Weight uint32 // Random weight
	Size   uint32 // Size of subtree, unnecessary if you don't need kth element
}

func New[T ordered.Ordered](value T) *TreapNode[T] {
	return &TreapNode[T]{Value: value, Left: nil, Right: nil, Weight: rand.Uint32(), Size: 1}
}

func (thisNode *TreapNode[T]) Update() {
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Size += thisNode.Right.Size
	}
}
