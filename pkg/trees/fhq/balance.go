package fhq

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
)

func merge[T order.Number](left *fhqTreeNode[T], right *fhqTreeNode[T]) *fhqTreeNode[T] {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.weight < right.weight {
		left.right = merge(left.right, right)
		left.Update()
		return left
	} else {
		right.left = merge(left, right.left)
		right.Update()
		return right
	}
}

func split[T order.Number](root *fhqTreeNode[T], key T) (*fhqTreeNode[T], *fhqTreeNode[T]) {
	if root == nil {
		return nil, nil
	}
	if root.value <= key {
		left, right := split(root.right, key)
		root.right = left
		root.Update()
		return root, right
	} else {
		left, right := split(root.left, key)
		root.left = right
		root.Update()
		return left, root
	}
}
