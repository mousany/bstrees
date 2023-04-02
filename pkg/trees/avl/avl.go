package avl

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
	"github.com/yanglinshu/bstrees/v2/pkg/errors"
)

type AVLTree[T order.Ordered] struct {
	root *avlTreeNode[T]
}

func New[T order.Ordered]() *AVLTree[T] {
	return &AVLTree[T]{root: nil}
}

func kth[T order.Ordered](root *avlTreeNode[T], k uint) *avlTreeNode[T] {
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

func insert[T order.Ordered](root *avlTreeNode[T], value T) *avlTreeNode[T] {
	if root == nil {
		return newAVLTreeNode(value)
	}
	if value < root.value {
		root.left = insert(root.left, value)
	} else {
		root.right = insert(root.right, value)
	}
	root.update()
	return balance(root)
}

func (t *AVLTree[T]) Insert(value T) {
	t.root = insert(t.root, value)
}

func delete[T order.Ordered](root *avlTreeNode[T], value T) *avlTreeNode[T] {
	if root == nil {
		return nil
	}
	if value < root.value {
		root.left = delete(root.left, value)
	} else if root.value < value {
		root.right = delete(root.right, value)
	} else {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		} else {
			minNode := kth(root.right, 1) // root.right is not nil, so this will not fail
			root.value = minNode.value
			root.right = delete(root.right, minNode.value)
		}
	}
	root.update()
	return balance(root)
}

func find[T order.Ordered](root *avlTreeNode[T], value T) *avlTreeNode[T] {
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

func (t *AVLTree[T]) Delete(value T) {
	t.root = delete(t.root, value)
}

func (t *AVLTree[T]) Contains(value T) bool {
	return find(t.root, value) != nil
}

func (t *AVLTree[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *AVLTree[T]) Height() int {
	if t.root == nil {
		return -1
	}
	return t.root.height
}

func (t *AVLTree[T]) Kth(k uint) (T, error) {
	result := kth(t.root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.value, nil
}

func (t *AVLTree[T]) Empty() bool {
	return t.root == nil
}

func (t *AVLTree[T]) Clear() {
	t.root = nil
}

func rank[T order.Ordered](root *avlTreeNode[T], value T) uint {
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

func (t *AVLTree[T]) Rank(value T) uint {
	return rank(t.root, value)
}

func prev[T order.Ordered](root *avlTreeNode[T], value T) *avlTreeNode[T] {
	var result *avlTreeNode[T] = nil
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

func (t *AVLTree[T]) Prev(value T) (T, error) {
	prev := prev(t.root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.value, nil
}

func next[T order.Ordered](root *avlTreeNode[T], value T) *avlTreeNode[T] {
	var result *avlTreeNode[T] = nil
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

func (t *AVLTree[T]) Next(value T) (T, error) {
	next := next(t.root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.value, nil
}
