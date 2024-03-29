package scapegoat

import (
	"github.com/yanglinshu/bstrees/v2"
	"golang.org/x/exp/constraints"
)

type ScapeGoatTree[T constraints.Integer | constraints.Float] struct {
	root  *scapeGoatTreeNode[T]
	alpha float64
}

func New[T constraints.Integer | constraints.Float](alpha float64) *ScapeGoatTree[T] {
	return &ScapeGoatTree[T]{
		root:  nil,
		alpha: alpha,
	}
}

func insert[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T, alpha float64) *scapeGoatTreeNode[T] {
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

func index[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T) uint {
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

func (t *ScapeGoatTree[T]) Index(value T) uint {
	return index(t.root, value)
}

func at[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], k uint) *scapeGoatTreeNode[T] {
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

func (t *ScapeGoatTree[T]) At(k uint) (T, error) {
	result := at(t.root, k)
	if result == nil {
		return T(rune(0)), bstrees.ErrIndexIsOutOfRange
	}
	return result.value, nil
}

func search[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	if root == nil {
		return nil
	}
	if root.value == value {
		if root.active() {
			return root
		} else {
			if result := search(root.left, value); result != nil {
				return result
			} else if result := search(root.right, value); result != nil {
				return result
			}
			return nil
		}
	} else if root.value > value {
		return search(root.left, value)
	} else {
		return search(root.right, value)
	}
}

func (t *ScapeGoatTree[T]) Contains(value T) bool {
	return search(t.root, value) != nil
}

func delete[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
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
	target := search(t.root, value)
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

func predecessor[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	return at(root, index(root, value)-1)
}

func (t *ScapeGoatTree[T]) Predecessor(value T) (T, error) {
	prev := predecessor(t.root, value)
	if prev == nil {
		return T(rune(0)), bstrees.ErrPredecessorDoesNotExist
	}
	return prev.value, nil
}

func successor[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], value T) *scapeGoatTreeNode[T] {
	return at(root, index(root, value+1))
}

func (t *ScapeGoatTree[T]) Successor(value T) (T, error) {
	next := successor(t.root, value)
	if next == nil {
		return T(rune(0)), bstrees.ErrSuccessorDoesNotExist
	}
	return next.value, nil
}
