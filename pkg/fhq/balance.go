package fhq

import (
	"bstrees/internal/node"
	"bstrees/pkg/treap"
)

func merge[T node.Number](left *treap.TreapTreeNode[T], right *treap.TreapTreeNode[T]) *treap.TreapTreeNode[T] {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.Weight() < right.Weight() {
		left.SetRight(merge(left.Right().(*treap.TreapTreeNode[T]), right))
		left.Update()
		return left
	} else {
		right.SetLeft(merge(left, right.Left().(*treap.TreapTreeNode[T])))
		right.Update()
		return right
	}
}

func split[T node.Number](root *treap.TreapTreeNode[T], key T) (*treap.TreapTreeNode[T], *treap.TreapTreeNode[T]) {
	if root == nil {
		return nil, nil
	}
	if root.Value() <= key {
		left, right := split(root.Right().(*treap.TreapTreeNode[T]), key)
		root.SetRight(left)
		root.Update()
		return root, right
	} else {
		left, right := split(root.Left().(*treap.TreapTreeNode[T]), key)
		root.SetLeft(right)
		root.Update()
		return left, root
	}
}
