package node

import "bstrees/pkg/trait/ordered"

type RBColor bool

const (
	Red   RBColor = true
	Black RBColor = false
)

type RBNode[T ordered.Ordered] struct {
	Value  T
	Left   *RBNode[T]
	Right  *RBNode[T]
	Color  RBColor
	Size   uint32     // Size of subtree, unnecessary if you don't need kth element
	// Father *RBNode[T] // Not necessary, but easier to implement
}

func New[T ordered.Ordered](value T) *RBNode[T] {
	return &RBNode[T]{Value: value, Left: nil, Right: nil, Color: Red, Size: 1}
}

func (thisNode *RBNode[T]) Update() {
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Size += thisNode.Right.Size
	}
}

func (thisNode *RBNode[T]) Red() bool {
	return thisNode.Color == Red
}

func (thisNode *RBNode[T]) Black() bool {
	return thisNode.Color == Black
}

func (thisNode *RBNode[T]) Full() bool {
	return thisNode.Left != nil && thisNode.Right != nil
}
