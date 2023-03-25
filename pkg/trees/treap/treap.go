package treap

import (
	"bstrees/internal/order"
	"bstrees/pkg/errors"
)

type TreapTree[T order.Ordered] struct {
	root *treapTreeNode[T]
}

func New[T order.Ordered]() *TreapTree[T] {
	return &TreapTree[T]{root: nil}
}

func kth[T order.Ordered](root *treapTreeNode[T], k uint) *treapTreeNode[T] {
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

func find[T order.Ordered](root *treapTreeNode[T], value T) *treapTreeNode[T] {
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

func (t *TreapTree[T]) Contains(value T) bool {
	return find(t.root, value) != nil
}

func insert[T order.Ordered](root *treapTreeNode[T], value T) *treapTreeNode[T] {
	if root == nil {
		return newTreapTreeNode(value)
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

func (t *TreapTree[T]) Insert(value T) {
	t.root = insert(t.root, value)
}

func delete[T order.Ordered](root *treapTreeNode[T], value T) *treapTreeNode[T] {
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

func (t *TreapTree[T]) Delete(value T) {
	t.root = delete(t.root, value)
}

func (t *TreapTree[T]) Kth(k uint) (T, error) {
	result := kth(t.root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.value, nil
}

func (t *TreapTree[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *TreapTree[T]) Empty() bool {
	return t.root == nil
}

func (t *TreapTree[T]) Clear() {
	t.root = nil
}

func rank[T order.Ordered](root *treapTreeNode[T], value T) uint {
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

func (t *TreapTree[T]) Rank(value T) uint {
	return rank(t.root, value)
}

func prev[T order.Ordered](root *treapTreeNode[T], value T) *treapTreeNode[T] {
	var result *treapTreeNode[T] = nil
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

func (t *TreapTree[T]) Prev(value T) (T, error) {
	result := prev(t.root, value)
	if result == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return result.value, nil
}

func next[T order.Ordered](root *treapTreeNode[T], value T) *treapTreeNode[T] {
	var result *treapTreeNode[T] = nil
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

func (t *TreapTree[T]) Next(value T) (T, error) {
	result := next(t.root, value)
	if result == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return result.value, nil
}
