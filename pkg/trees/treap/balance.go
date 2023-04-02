package treap

import "github.com/yanglinshu/bstrees/v2/internal/order"

func leftRotate[T order.Ordered](root *treapTreeNode[T]) *treapTreeNode[T] {
	right := root.right
	root.right = right.left
	right.left = root
	root.Update()
	right.Update()
	return right
}

func rightRotate[T order.Ordered](root *treapTreeNode[T]) *treapTreeNode[T] {
	left := root.left
	root.left = left.right
	left.right = root
	root.Update()
	left.Update()
	return left
}
