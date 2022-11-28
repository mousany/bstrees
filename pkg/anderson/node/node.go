package node

import "bstrees/pkg/trait/ordered"

type AndersonNode[T ordered.Ordered] struct {
	Value T
	Left  *AndersonNode[T]
	Right *AndersonNode[T]
	Size  uint32
	Level uint32
}

func New[T ordered.Ordered](value T, level uint32) *AndersonNode[T] {
	return &AndersonNode[T]{
		Value: value,
		Left:  nil,
		Right: nil,
		Size:  uint32(1),
		Level: level,
	}
}

func (root *AndersonNode[T]) Update() {
	root.Size = uint32(1)
	if root.Left != nil {
		root.Size += root.Left.Size
	}
	if root.Right != nil {
		root.Size += root.Right.Size
	}
}

func (root *AndersonNode[T]) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *AndersonNode[T]) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *AndersonNode[T]) SetChild(child *AndersonNode[T], direction bool) {
	if direction {
		root.Right = child
	} else {
		root.Left = child
	}
}

func (root *AndersonNode[T]) Child(direction bool) *AndersonNode[T] {
	if direction {
		return root.Right
	}
	return root.Left
}
