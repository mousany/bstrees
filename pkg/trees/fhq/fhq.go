package fhq

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
	"github.com/yanglinshu/bstrees/v2/pkg/errors"
)

type FHQTree[T order.Number] struct {
	root *fhqTreeNode[T]
}

func New[T order.Number]() *FHQTree[T] {
	return &FHQTree[T]{root: nil}
}

func find[T order.Number](root *fhqTreeNode[T], value T) *fhqTreeNode[T] {
	for root != nil {
		if value < root.value {
			root = root.left
		} else if root.value < value {
			root = root.right
		} else {
			return root
		}
	}
	return nil
}

func kth[T order.Number](root *fhqTreeNode[T], k uint) *fhqTreeNode[T] {
	for root != nil {
		leftSize := uint(0)
		if root.left != nil {
			leftSize = root.left.size
		}
		if leftSize+1 == k {
			return root
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.right
		} else {
			root = root.left
		}
	}
	return nil
}

func (t *FHQTree[T]) Insert(value T) {
	left, right := split(t.root, value)
	t.root = merge(merge(left, newFHQTreeNode(value)), right)
}

func (t *FHQTree[T]) Delete(value T) {
	left, right := split(t.root, value)
	left, mid := split(left, value-1)
	if mid != nil {
		mid = merge(mid.left, mid.right)
	}
	t.root = merge(merge(left, mid), right)
}

func (t *FHQTree[T]) Contains(value T) bool {
	return find(t.root, value) != nil
}

func (t *FHQTree[T]) Rank(value T) uint {
	left, right := split(t.root, value-1)
	defer func() {
		t.root = merge(left, right)
	}()
	if left == nil {
		return 1
	}
	return left.size + 1
}

func (t *FHQTree[T]) Kth(k uint) (T, error) {
	result := kth(t.root, k)
	if result == nil {
		return T(0), errors.ErrOutOfRange
	}
	return result.value, nil
}

func (t *FHQTree[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *FHQTree[T]) Empty() bool {
	return t.root == nil
}

func (t *FHQTree[T]) Clear() {
	t.root = nil
}

func (t *FHQTree[T]) Prev(value T) (T, error) {
	left, right := split(t.root, value-1)
	defer func() {
		t.root = merge(left, right)
	}()
	result := kth(left, left.size)
	if result == nil {
		return T(0), errors.ErrNoPrevValue
	}
	return result.value, nil
}

func (t *FHQTree[T]) Next(value T) (T, error) {
	left, right := split(t.root, value)
	defer func() {
		t.root = merge(left, right)
	}()
	result := kth(right, 1)
	if result == nil {
		return T(0), errors.ErrNoNextValue
	}
	return result.value, nil
}
