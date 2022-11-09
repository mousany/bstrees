package vanilla

import (
	"bstrees/pkg/treap/node"
	"bstrees/pkg/type/number"
	"errors"
)

type VanillaTreap[T number.Number] struct {
	Root *node.TreapNode[T]
}

func New[T number.Number]() *VanillaTreap[T] {
	return &VanillaTreap[T]{Root: nil}
}

func LeftRotate[T number.Number](root *node.TreapNode[T]) *node.TreapNode[T] {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate[T number.Number](root *node.TreapNode[T]) *node.TreapNode[T] {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Kth[T number.Number](root *node.TreapNode[T], k uint32) *node.TreapNode[T] {
	for root != nil {
		leftSize := uint32(0)
		if root.Left != nil {
			leftSize = root.Left.Size
		}
		if leftSize+1 == k {
			return root
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return nil
}

func Insert[T number.Number](root *node.TreapNode[T], value T) *node.TreapNode[T] {
	if root == nil {
		return node.New(value)
	}
	if root.Value <= value {
		root.Right = Insert(root.Right, value)
		if root.Right.Weight < root.Weight {
			root = LeftRotate(root)
		}
	} else {
		root.Left = Insert(root.Left, value)
		if root.Left.Weight < root.Weight {
			root = RightRotate(root)
		}
	}
	root.Update()
	return root
}

func (this *VanillaTreap[T]) Insert(value T) {
	this.Root = Insert(this.Root, value)
}

func Delete[T number.Number](root *node.TreapNode[T], value T) *node.TreapNode[T] {
	if root == nil {
		return nil
	}
	if root.Value == value {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		if root.Left.Weight < root.Right.Weight {
			root = RightRotate(root)
			root.Right = Delete(root.Right, value)
		} else {
			root = LeftRotate(root)
			root.Left = Delete(root.Left, value)
		}
	} else if root.Value < value {
		root.Right = Delete(root.Right, value)
	} else {
		root.Left = Delete(root.Left, value)
	}
	root.Update()
	return root
}

func (this *VanillaTreap[T]) Delete(value T) {
	this.Root = Delete(this.Root, value)
}

func (this *VanillaTreap[T]) Kth(k uint32) (T, error) {
	result := Kth(this.Root, k)
	if result == nil {
		return T(0), errors.New("k is out of range")
	}
	return result.Value, nil
}

func (this *VanillaTreap[T]) Size() uint32 {
	if this.Root == nil {
		return 0
	}
	return this.Root.Size
}

func (this *VanillaTreap[T]) Empty() bool {
	return this.Root == nil
}

func (this *VanillaTreap[T]) Clear() {
	this.Root = nil
}

func Rank[T number.Number](root *node.TreapNode[T], value T) uint32 {
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
	return rank
}

func (this *VanillaTreap[T]) Rank(value T) uint32 {
	return Rank(this.Root, value)
}

func Prev[T number.Number](root *node.TreapNode[T], value T) *node.TreapNode[T] {
	var result *node.TreapNode[T] = nil
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

func (this *VanillaTreap[T]) Prev(value T) (T, error) {
	result := Prev(this.Root, value)
	if result == nil {
		return T(0), errors.New("no prev value")
	}
	return result.Value, nil
}

func Next[T number.Number](root *node.TreapNode[T], value T) *node.TreapNode[T] {
	var result *node.TreapNode[T] = nil
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

func (this *VanillaTreap[T]) Next(value T) (T, error) {
	result := Next(this.Root, value)
	if result == nil {
		return T(0), errors.New("no next value")
	}
	return result.Value, nil
}
