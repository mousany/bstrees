package rbtree

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/rbtree/node"
	"bstrees/pkg/trait/ordered"
)

type RBTree[T ordered.Ordered] struct {
	Root *node.RBNode[T]
}

func New[T ordered.Ordered]() *RBTree[T] {
	return &RBTree[T]{Root: nil}
}

func LeftRotate[T ordered.Ordered](root *node.RBNode[T]) *node.RBNode[T] {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate[T ordered.Ordered](root *node.RBNode[T]) *node.RBNode[T] {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Kth[T ordered.Ordered](root *node.RBNode[T], k uint32) (T, error) {
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

func Insert[T ordered.Ordered](root *node.RBNode[T], value T) {

}

func (thisTree *RBTree[T]) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *RBTree[T]) Kth(k uint32) (T, error) {
	return Kth(thisTree.Root, k)
}

func (thisTree *RBTree[T]) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *RBTree[T]) Clear() {
	thisTree.Root = nil
}

func Rank[T ordered.Ordered](root *node.RBNode[T], value T) uint32 {
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

func (thisTree *RBTree[T]) Rank(value T) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var prev *node.RBNode[T] = nil
	for root != nil {
		if root.Value < value {
			prev = root
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return prev
}

func (thisTree *RBTree[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var next *node.RBNode[T] = nil
	for root != nil {
		if root.Value > value {
			next = root
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return next
}

func (thisTree *RBTree[T]) Next(value T) (T, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value, nil
}
