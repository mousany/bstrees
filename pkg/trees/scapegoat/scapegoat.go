package scapegoat

import (
	"bstrees/internal/order"
	"bstrees/pkg/errors"
)

type ScapeGoatTree[T order.Number] struct {
	root  *scapeGoatTreeNode[T]
	alpha float64
}

func New[T order.Number](alpha float64) *ScapeGoatTree[T] {
	return &ScapeGoatTree[T]{
		root:  nil,
		alpha: alpha,
	}
}

func insert[T order.Number](root *scapeGoatTreeNode[T], value T, alpha float64) *scapeGoatTreeNode[T] {
	if root == nil {
		return newScapeGoatTreeNode(value)
	}
	if value < root.value {
		root.left = insert(root.left, value, alpha)
	} else {
		root.right = insert(root.right, value, alpha)
	}
	root.update()
	if imbalance(root, alpha) {
		return reconstruct(root)
	}
	return root
}

func (t *ScapeGoatTree[T]) Insert(value T) {
	t.root = insert(t.root, value, t.alpha)
}

func rank[T order.Number](root *scapeGoatTreeNode[T], value T) uint {
	result := uint(0)
	for root != nil {
		if root.value >= value {
			root = root.left
		} else {
			if root.left != nil {
				result += root.left.size
			}
			if root.active() {
				result += 1
			}
			root = root.right
		}
	}
	return result + 1
}

func (t *ScapeGoatTree[T]) Rank(value T) uint {
	return rank(t.root, value)
}

func kth[T order.Number](root *scapeGoatTreeNode[T], k uint) *scapeGoatTreeNode[T] {
	var result *scapeGoatTreeNode[T] = nil
	for root != nil {
		leftSize := uint(0)
		if root.left != nil {
			leftSize = root.left.size
		}
		if root.active() && leftSize+1 == k {
			result = root
			break
		} else if leftSize >= k {
			root = root.left
		} else {
			k -= leftSize
			if root.active() {
				k -= 1
			}
			root = root.right
		}
	}
	return result
}

func (t *ScapeGoatTree[T]) Kth(k uint) (T, error) {
	result := kth(t.root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.value, nil
}

func find[T order.Number](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	if root == nil {
		return nil
	}
	if root.value == value {
		if root.active() {
			return root
		} else {
			if result := find(root.left, value); result != nil {
				return result
			} else if result := find(root.right, value); result != nil {
				return result
			}
			return nil
		}
	} else if root.value > value {
		return find(root.left, value)
	} else {
		return find(root.right, value)
	}
}

func (t *ScapeGoatTree[T]) Find(value T) bool {
	return find(t.root, value) != nil
}

func delete[T order.Number](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	if root == nil {
		return nil
	}
	if root.value == value {
		if root.active() {
			root.deactivate()
			return root
		} else {
			if result := delete(root.left, value); result != nil {
				root.size -= 1
				return result
			} else if result := delete(root.right, value); result != nil {
				root.size -= 1
				return result
			}
			return nil
		}
	} else if root.value > value {
		if result := delete(root.left, value); result != nil {
			root.size -= 1
			return result
		}
	} else {
		if result := delete(root.right, value); result != nil {
			root.size -= 1
			return result
		}
	}
	return nil
}

func (t *ScapeGoatTree[T]) Delete(value T) {
	target := find(t.root, value)
	if target != nil {
		delete(t.root, value)
	}
}

func (t *ScapeGoatTree[T]) Clear() {
	t.root = nil
}

func (t *ScapeGoatTree[T]) Size() uint {
	return t.root.size
}

func (t *ScapeGoatTree[T]) Empty() bool {
	return t.root == nil
}

func prev[T order.Number](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	return kth(root, rank(root, value)-1)
}

func (t *ScapeGoatTree[T]) Prev(value T) (T, error) {
	prev := prev(t.root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.value, nil
}

func next[T order.Number](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	return kth(root, rank(root, value+1))
}

func (t *ScapeGoatTree[T]) Next(value T) (T, error) {
	next := next(t.root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.value, nil
}
