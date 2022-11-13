package node

import (
	"bstrees/pkg/trait/ordered"
	"bstrees/pkg/util/vendor"
)

type AVLNode[T ordered.Ordered] struct {
	Value  T
	Left   *AVLNode[T]
	Right  *AVLNode[T]
	Height int32  // Height of the node
	Size   uint32 // Size of subtree, unnecessary if you don't need kth element
}

func New[T ordered.Ordered](value T) *AVLNode[T] {
	return &AVLNode[T]{Value: value, Left: nil, Right: nil, Height: 0, Size: 1}
}

func (thisNode *AVLNode[T]) Update() {
	thisNode.Height = 0
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Height = vendor.Max(thisNode.Height, thisNode.Left.Height+1)
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Height = vendor.Max(thisNode.Height, thisNode.Right.Height+1)
		thisNode.Size += thisNode.Right.Size
	}
}
