package node

import "bstrees/pkg/trait/ordered"

type SplayNode[T ordered.Ordered] struct {
	Value  T
	Left   *SplayNode[T]
	Right  *SplayNode[T]
	Parent *SplayNode[T]
	Size   uint32
}

func New[T ordered.Ordered](value T) *SplayNode[T] {
	return &SplayNode[T]{Value: value, Left: nil, Right: nil, Parent: nil, Size: 1}
}

func (root *SplayNode[T]) Update() {
	root.Size = 1
	if root.Left != nil {
		root.Size += root.Left.Size
	}
	if root.Right != nil {
		root.Size += root.Right.Size
	}
}

func (root *SplayNode[T]) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *SplayNode[T]) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *SplayNode[T]) SetChild(child *SplayNode[T], direction bool) {
	if direction {
		root.Right = child
	} else {
		root.Left = child
	}
	if child != nil {
		child.Parent = root
	}
}

func (root *SplayNode[T]) Child(direction bool) *SplayNode[T] {
	if direction {
		return root.Right
	}
	return root.Left
}
