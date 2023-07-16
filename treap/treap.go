package treap

import (
	"github.com/yanglinshu/bstrees/v2"
	"golang.org/x/exp/constraints"
)

type Treap[T constraints.Ordered] struct {
	root *treapNode[T]
}

func New[T constraints.Ordered]() *Treap[T] {
	return &Treap[T]{root: nil}
}

func at[T constraints.Ordered](root *treapNode[T], k uint) *treapNode[T] {
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

func search[T constraints.Ordered](root *treapNode[T], value T) *treapNode[T] {
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

func (t *Treap[T]) Contains(value T) bool {
	return search(t.root, value) != nil
}

func insert[T constraints.Ordered](root *treapNode[T], value T) *treapNode[T] {
	if root == nil {
		return newTreapNode(value)
	}
	if root.value <= value {
		root.right = insert(root.right, value)
		if root.right.weight < root.weight {
			root = leftRotate(root)
		}
	} else {
		root.left = insert(root.left, value)
		if root.left.weight < root.weight {
			root = rightRotate(root)
		}
	}
	root.Update()
	return root
}

func (t *Treap[T]) Insert(value T) {
	t.root = insert(t.root, value)
}

func delete[T constraints.Ordered](root *treapNode[T], value T) *treapNode[T] {
	if root == nil {
		return nil
	}
	if root.value == value {
		if root.left == nil {
			return root.right
		}
		if root.right == nil {
			return root.left
		}
		if root.left.weight < root.right.weight {
			root = rightRotate(root)
			root.right = delete(root.right, value)
		} else {
			root = leftRotate(root)
			root.left = delete(root.left, value)
		}
	} else if root.value < value {
		root.right = delete(root.right, value)
	} else {
		root.left = delete(root.left, value)
	}
	root.Update()
	return root
}

func (t *Treap[T]) Delete(value T) {
	t.root = delete(t.root, value)
}

func (t *Treap[T]) At(k uint) (T, error) {
	result := at(t.root, k)
	if result == nil {
		return T(rune(0)), bstrees.ErrIndexIsOutOfRange
	}
	return result.value, nil
}

func (t *Treap[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *Treap[T]) Empty() bool {
	return t.root == nil
}

func (t *Treap[T]) Clear() {
	t.root = nil
}

func index[T constraints.Ordered](root *treapNode[T], value T) uint {
	rank := uint(0)
	for root != nil {
		if root.value < value {
			rank += 1
			if root.left != nil {
				rank += root.left.size
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return rank + 1
}

func (t *Treap[T]) Index(value T) uint {
	return index(t.root, value)
}

func predecessor[T constraints.Ordered](root *treapNode[T], value T) *treapNode[T] {
	var result *treapNode[T] = nil
	for root != nil {
		if root.value < value {
			result = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return result
}

func (t *Treap[T]) Predecessor(value T) (T, error) {
	result := predecessor(t.root, value)
	if result == nil {
		return T(rune(0)), bstrees.ErrPredecessorDoesNotExist
	}
	return result.value, nil
}

func successor[T constraints.Ordered](root *treapNode[T], value T) *treapNode[T] {
	var result *treapNode[T] = nil
	for root != nil {
		if root.value > value {
			result = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return result
}

func (t *Treap[T]) Successor(value T) (T, error) {
	result := successor(t.root, value)
	if result == nil {
		return T(rune(0)), bstrees.ErrSuccessorDoesNotExist
	}
	return result.value, nil
}
