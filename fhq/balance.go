package fhq

import (
	"golang.org/x/exp/constraints"
)

func merge[T constraints.Integer | constraints.Float](left *fhqTreapNode[T], right *fhqTreapNode[T]) *fhqTreapNode[T] {
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

func split[T constraints.Integer | constraints.Float](root *fhqTreapNode[T], key T) (*fhqTreapNode[T], *fhqTreapNode[T]) {
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
