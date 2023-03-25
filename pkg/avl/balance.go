package avl

import (
	"bstrees/internal/order"
)

func leftRotate[T order.Ordered](root *avlTreeNode[T]) *avlTreeNode[T] {
	right := root.right
	root.right = right.left
	right.left = root
	root.Update()
	right.Update()
	return right
}

func rightRotate[T order.Ordered](root *avlTreeNode[T]) *avlTreeNode[T] {
	left := root.left
	root.left = left.right
	left.right = root
	root.Update()
	left.Update()
	return left
}

func balance[T order.Ordered](root *avlTreeNode[T]) *avlTreeNode[T] {
	leftHeight := -1
	if root.left != nil {
		leftHeight = root.left.height
	}
	rightHeight := -1
	if root.right != nil {
		rightHeight = root.right.height
	}
	if leftHeight > rightHeight+1 {
		left := root.left
		leftLeftHeight := -1
		if left.left != nil {
			leftLeftHeight = left.left.height
		}
		leftRightHeight := -1
		if left.right != nil {
			leftRightHeight = left.right.height
		}
		if leftLeftHeight < leftRightHeight {
			root.left = leftRotate(left)
		}
		ret := rightRotate(root)
		return ret
	} else if rightHeight > leftHeight+1 {
		right := root.right
		rightLeftHeight := -1
		if right.left != nil {
			rightLeftHeight = right.left.height
		}
		rightRightHeight := -1
		if right.right != nil {
			rightRightHeight = right.right.height
		}
		if rightRightHeight < rightLeftHeight {
			root.right = rightRotate(right)
		}
		return leftRotate(root)
	}
	return root
}
