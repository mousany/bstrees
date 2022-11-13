package avltree

import (
	"bstrees/pkg/avltree/node"
	"bstrees/pkg/errors"
	"bstrees/pkg/trait/ordered"
)

type AVLTree[T ordered.Ordered] struct {
	Root *node.AVLNode[T]
}

func New[T ordered.Ordered]() *AVLTree[T] {
	return &AVLTree[T]{Root: nil}
}

func LeftRotate[T ordered.Ordered](root *node.AVLNode[T]) *node.AVLNode[T] {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate[T ordered.Ordered](root *node.AVLNode[T]) *node.AVLNode[T] {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Balance[T ordered.Ordered](root *node.AVLNode[T]) *node.AVLNode[T] {
	leftHeight := int32(-1)
	if root.Left != nil {
		leftHeight = root.Left.Height
	}
	rightHeight := int32(-1)
	if root.Right != nil {
		rightHeight = root.Right.Height
	}
	if leftHeight > rightHeight+1 {
		left := root.Left
		leftLeftHeight := int32(-1)
		if left.Left != nil {
			leftLeftHeight = left.Left.Height
		}
		leftRightHeight := int32(-1)
		if left.Right != nil {
			leftRightHeight = left.Right.Height
		}
		if leftLeftHeight < leftRightHeight {
			root.Left = LeftRotate(left)
		}
		ret := RightRotate(root)
		return ret
	} else if rightHeight > leftHeight+1 {
		right := root.Right
		rightLeftHeight := int32(-1)
		if right.Left != nil {
			rightLeftHeight = right.Left.Height
		}
		rightRightHeight := int32(-1)
		if right.Right != nil {
			rightRightHeight = right.Right.Height
		}
		if rightRightHeight < rightLeftHeight {
			root.Right = RightRotate(right)
		}
		return LeftRotate(root)
	}
	return root
}

func Kth[T ordered.Ordered](root *node.AVLNode[T], k uint32) (T, error) {
	for root != nil {
		leftSize := uint32(0)
		if root.Left != nil {
			leftSize = root.Left.Size
		}
		if leftSize+1 == k {
			return root.Value, nil
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return T(rune(0)), errors.ErrOutOfRange
}

func Insert[T ordered.Ordered](root *node.AVLNode[T], value T) *node.AVLNode[T] {
	if root == nil {
		return node.New(value)
	}
	if value < root.Value {
		root.Left = Insert(root.Left, value)
	} else {
		root.Right = Insert(root.Right, value)
	}
	root.Update()
	return Balance(root)
}

func (thisTree *AVLTree[T]) Insert(value T) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func Delete[T ordered.Ordered](root *node.AVLNode[T], value T) *node.AVLNode[T] {
	if root == nil {
		return nil
	}
	if value < root.Value {
		root.Left = Delete(root.Left, value)
	} else if root.Value < value {
		root.Right = Delete(root.Right, value)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		} else {
			min, _ := Kth(root.Right, 1) // root.Right is not nil, so this will not fail
			root.Value = min
			root.Right = Delete(root.Right, min)
		}
	}
	root.Update()
	return Balance(root)
}

func (thisTree *AVLTree[T]) Delete(value T) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *AVLTree[T]) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *AVLTree[T]) Height() int32 {
	if thisTree.Root == nil {
		return -1
	}
	return thisTree.Root.Height
}

func (thisTree *AVLTree[T]) Kth(k uint32) (T, error) {
	return Kth(thisTree.Root, k)
}

func (thisTree *AVLTree[T]) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *AVLTree[T]) Clear() {
	thisTree.Root = nil
}

func Rank[T ordered.Ordered](root *node.AVLNode[T], value T) uint32 {
	rank := uint32(0)
	for root != nil {
		if root.Value < value {
			rank += 1
			if root.Left != nil {
				rank += root.Left.Size
			}
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return rank + 1
}

func (thisTree *AVLTree[T]) Rank(value T) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev[T ordered.Ordered](root *node.AVLNode[T], value T) *node.AVLNode[T] {
	var result *node.AVLNode[T] = nil
	for root != nil {
		if root.Value < value {
			result = root
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return result
}

func (thisTree *AVLTree[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T ordered.Ordered](root *node.AVLNode[T], value T) *node.AVLNode[T] {
	var result *node.AVLNode[T] = nil
	for root != nil {
		if root.Value > value {
			result = root
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return result
}

func (thisTree *AVLTree[T]) Next(value T) (T, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value, nil
}
