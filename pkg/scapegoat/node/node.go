package node

import "bstrees/pkg/trait/number"

type NodeState bool

const (
	Inactive NodeState = false
	Active   NodeState = true
)

type ScapeGoatNode[T number.Number] struct {
	Value  T
	Left   *ScapeGoatNode[T]
	Right  *ScapeGoatNode[T]
	State  NodeState
	Size   uint32 // Number of active nodes in the subtree
	Weight uint32 // Number of nodes in the subtree
}

func New[T number.Number](value T) *ScapeGoatNode[T] {
	return &ScapeGoatNode[T]{
		Value:  value,
		Left:   nil,
		Right:  nil,
		State:  Active,
		Size:   1,
		Weight: 1,
	}
}

func (root *ScapeGoatNode[T]) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *ScapeGoatNode[T]) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *ScapeGoatNode[T]) Inactive() bool {
	return root.State == Inactive
}

func (root *ScapeGoatNode[T]) Active() bool {
	return root.State == Active
}

func (root *ScapeGoatNode[T]) Update() {
	if root.Active() {
		root.Size = 1
	} else {
		root.Size = 0
	}
	root.Weight = 1
	if root.Left != nil {
		root.Size += root.Left.Size
		root.Weight += root.Left.Weight
	}
	if root.Right != nil {
		root.Size += root.Right.Size
		root.Weight += root.Right.Weight
	}
}

func (root *ScapeGoatNode[T]) Deactivate() {
	root.State = Inactive
	root.Update()
}
