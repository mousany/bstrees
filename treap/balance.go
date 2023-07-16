package treap

import "golang.org/x/exp/constraints"

func leftRotate[T constraints.Ordered](root *treapNode[T]) *treapNode[T] {
	right := root.right
	root.right = right.left
	right.left = root
	root.Update()
	right.Update()
	return right
}

func rightRotate[T constraints.Ordered](root *treapNode[T]) *treapNode[T] {
	left := root.left
	root.left = left.right
	left.right = root
	root.Update()
	left.Update()
	return left
}
