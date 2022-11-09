package rotateless

import (
	"bstrees/pkg/treap/node"
	"bstrees/pkg/type/number"
	"errors"
)

type RotatelessTreap[T number.Number] struct {
	Root *node.TreapNode[T]
}

func New[T number.Number]() *RotatelessTreap[T] {
	return &RotatelessTreap[T]{Root: nil}
}

func Merge[T number.Number](left *node.TreapNode[T], right *node.TreapNode[T]) *node.TreapNode[T] {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.Weight < right.Weight {
		left.Right = Merge(left.Right, right)
		left.Update()
		return left
	} else {
		right.Left = Merge(left, right.Left)
		right.Update()
		return right
	}
}

func Split[T number.Number](root *node.TreapNode[T], key T) (*node.TreapNode[T], *node.TreapNode[T]) {
	if root == nil {
		return nil, nil
	}
	if root.Value <= key {
		left, right := Split(root.Right, key)
		root.Right = left
		root.Update()
		return root, right
	} else {
		left, right := Split(root.Left, key)
		root.Left = right
		root.Update()
		return left, root
	}
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

func (this *RotatelessTreap[T]) Insert(value T) {
	left, right := Split(this.Root, value)
	this.Root = Merge(Merge(left, node.New(value)), right)
}

func (this *RotatelessTreap[T]) Delete(value T) {
	left, right := Split(this.Root, value)
	left, mid := Split(left, value-1)
	if mid != nil {
		mid = Merge(mid.Left, mid.Right)
	}
	this.Root = Merge(Merge(left, mid), right)
}

func (this *RotatelessTreap[T]) Rank(value T) uint32 {
	left, right := Split(this.Root, value-1)
	defer func() {
		this.Root = Merge(left, right)
	}()
	if left == nil {
		return 1
	}
	return left.Size + 1
}

func (this *RotatelessTreap[T]) Kth(k uint32) (T, error) {
	result := Kth(this.Root, k)
	if result == nil {
		return T(0), errors.New("k is out of range")
	}
	return result.Value, nil
}

func (this *RotatelessTreap[T]) Size() uint32 {
	if this.Root == nil {
		return 0
	}
	return this.Root.Size
}

func (this *RotatelessTreap[T]) Empty() bool {
	return this.Root == nil
}

func (this *RotatelessTreap[T]) Clear() {
	this.Root = nil
}

func (this *RotatelessTreap[T]) Prev(value T) (T, error) {
	left, right := Split(this.Root, value-1)
	defer func() {
		this.Root = Merge(left, right)
	}()
	result := Kth(left, left.Size)
	if result == nil {
		return T(0), errors.New("no prev")
	}
	return result.Value, nil
}

func (this *RotatelessTreap[T]) Next(value T) (T, error) {
	left, right := Split(this.Root, value)
	defer func() {
		this.Root = Merge(left, right)
	}()
	result := Kth(right, 1)
	if result == nil {
		return T(0), errors.New("no next")
	}
	return result.Value, nil
}
