package anderson

import "bstrees/internal/order"

func leftRotate[T order.Ordered](root *andersonTreeNode[T]) *andersonTreeNode[T] {
	right := root.right
	root.right = right.left
	right.left = root
	root.Update()
	right.Update()
	return right
}

func rightRotate[T order.Ordered](root *andersonTreeNode[T]) *andersonTreeNode[T] {
	left := root.left
	root.left = left.right
	left.right = root
	root.Update()
	left.Update()
	return left
}

func skew[T order.Ordered](root *andersonTreeNode[T]) *andersonTreeNode[T] {
	if root.left == nil || root.left.level != root.level {
		return root
	}
	return rightRotate(root)
}

func split[T order.Ordered](root *andersonTreeNode[T]) *andersonTreeNode[T] {
	if root.right == nil || root.right.right == nil || root.right.right.level != root.level {
		return root
	}
	root = leftRotate(root)
	root.right.level += 1
	return root
}
