package anderson

import (
	"github.com/yanglinshu/bstrees/v2"
	"golang.org/x/exp/constraints"
)

type AndersonTree[T constraints.Ordered] struct {
	root *andersonTreeNode[T]
}

func New[T constraints.Ordered]() AndersonTree[T] {
	return AndersonTree[T]{root: nil}
}

func insert[T constraints.Ordered](root *andersonTreeNode[T], value T) *andersonTreeNode[T] {
	if root == nil {
		return newAndersonTreeNode(value, 1)
	}
	if value < root.value {
		root.left = insert(root.left, value)
	} else {
		root.right = insert(root.right, value)
	}
	root.update()
	root = skew(root)
	root = split(root)
	return root
}

func delete[T constraints.Ordered](root *andersonTreeNode[T], value T) *andersonTreeNode[T] {
	if root == nil {
		return nil
	}
	if value < root.value {
		root.left = delete(root.left, value)
	} else if value > root.value {
		root.right = delete(root.right, value)
	} else {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		} else {
			minNode := at(root.right, 1)
			root.value = minNode.value
			root.right = delete(root.right, minNode.value)
		}
	}
	root.update()
	if (root.left != nil && root.left.level < root.level-1) ||
		(root.right != nil && root.right.level < root.level-1) {
		root.level -= 1
		if root.right != nil && root.right.level > root.level {
			root.right.level = root.level
		}
		root = skew(root)
		root = split(root)
	}
	return root
}

func at[T constraints.Ordered](root *andersonTreeNode[T], k uint) *andersonTreeNode[T] {
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

func (t *AndersonTree[T]) Insert(value T) {
	t.root = insert(t.root, value)
}

func (t *AndersonTree[T]) Delete(value T) {
	t.root = delete(t.root, value)
}

func (t *AndersonTree[T]) At(k uint) (T, error) {
	root := at(t.root, k)
	if root == nil {
		return T(rune(0)), bstrees.ErrIndexIsOutOfRange
	}
	return root.value, nil
}

func (t *AndersonTree[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *AndersonTree[T]) Empty() bool {
	return t.root == nil
}

func (t *AndersonTree[T]) Clear() {
	t.root = nil
}

func search[T constraints.Ordered](root *andersonTreeNode[T], value T) *andersonTreeNode[T] {
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

func (t *AndersonTree[T]) Contains(value T) bool {
	return search(t.root, value) != nil
}

func index[T constraints.Ordered](root *andersonTreeNode[T], value T) uint {
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

func (t *AndersonTree[T]) Index(value T) uint {
	return index(t.root, value)
}

func predecessor[T constraints.Ordered](root *andersonTreeNode[T], value T) *andersonTreeNode[T] {
	var prev *andersonTreeNode[T] = nil
	for root != nil {
		if root.value < value {
			prev = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return prev
}

func (t *AndersonTree[T]) Predecessor(value T) (T, error) {
	prev := predecessor(t.root, value)
	if prev == nil {
		return T(rune(0)), bstrees.ErrPredecessorDoesNotExist
	}
	return prev.value, nil
}

func successor[T constraints.Ordered](root *andersonTreeNode[T], value T) *andersonTreeNode[T] {
	var next *andersonTreeNode[T] = nil
	for root != nil {
		if root.value > value {
			next = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return next
}

func (t *AndersonTree[T]) Successor(value T) (T, error) {
	prev := successor(t.root, value)
	if prev == nil {
		return T(rune(0)), bstrees.ErrSuccessorDoesNotExist
	}
	return prev.value, nil
}
