package fhq

import (
	"github.com/yanglinshu/bstrees/v2"
	"golang.org/x/exp/constraints"
)

type FHQTreap[T constraints.Integer | constraints.Float] struct {
	root *fhqTreapNode[T]
}

func New[T constraints.Integer | constraints.Float]() *FHQTreap[T] {
	return &FHQTreap[T]{root: nil}
}

func search[T constraints.Integer | constraints.Float](root *fhqTreapNode[T], value T) *fhqTreapNode[T] {
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

func At[T constraints.Integer | constraints.Float](root *fhqTreapNode[T], k uint) *fhqTreapNode[T] {
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

func (t *FHQTreap[T]) Insert(value T) {
	left, right := split(t.root, value)
	t.root = merge(merge(left, newFHQTreapNode(value)), right)
}

func (t *FHQTreap[T]) Delete(value T) {
	left, right := split(t.root, value)
	left, mid := split(left, value-1)
	if mid != nil {
		mid = merge(mid.left, mid.right)
	}
	t.root = merge(merge(left, mid), right)
}

func (t *FHQTreap[T]) Contains(value T) bool {
	return search(t.root, value) != nil
}

func (t *FHQTreap[T]) Index(value T) uint {
	left, right := split(t.root, value-1)
	defer func() {
		t.root = merge(left, right)
	}()
	if left == nil {
		return 1
	}
	return left.size + 1
}

func (t *FHQTreap[T]) At(k uint) (T, error) {
	result := At(t.root, k)
	if result == nil {
		return T(0), bstrees.ErrIndexIsOutOfRange
	}
	return result.value, nil
}

func (t *FHQTreap[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *FHQTreap[T]) Empty() bool {
	return t.root == nil
}

func (t *FHQTreap[T]) Clear() {
	t.root = nil
}

func (t *FHQTreap[T]) Predecessor(value T) (T, error) {
	left, right := split(t.root, value-1)
	defer func() {
		t.root = merge(left, right)
	}()
	result := At(left, left.size)
	if result == nil {
		return T(0), bstrees.ErrPredecessorDoesNotExist
	}
	return result.value, nil
}

func (t *FHQTreap[T]) Successor(value T) (T, error) {
	left, right := split(t.root, value)
	defer func() {
		t.root = merge(left, right)
	}()
	result := At(right, 1)
	if result == nil {
		return T(0), bstrees.ErrSuccessorDoesNotExist
	}
	return result.value, nil
}
