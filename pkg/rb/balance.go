package rb

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
	"bstrees/pkg/errors"
)

func singleRotate[T node.Ordered](direction bool, root *rbTreeNode[T]) node.Noded[T] {
	root = tree.SingleRotate(direction, node.Noded[T](root)).(*rbTreeNode[T])
	root.SetColor(black)
	root.Child(direction).(*rbTreeNode[T]).SetColor(red)
	return root
}

func doubleRotate[T node.Ordered](direction bool, root *rbTreeNode[T]) node.Noded[T] {
	root.SetChild(!direction, singleRotate(!direction, root.Child(!direction).(*rbTreeNode[T])))
	return singleRotate(direction, root)
}

func CheckRBProperty[T node.Ordered](root *rbTreeNode[T]) (uint, error) {
	if root == nil {
		return uint(1), nil
	}
	left, right := root.Left().(*rbTreeNode[T]), root.Right().(*rbTreeNode[T])
	if root.IsRed() {
		if left.IsRed() || right.IsRed() {
			return 0, errors.ErrViolatedRBProperty
		}
	}
	leftHeight, leftOk := CheckRBProperty(left)
	rightHeight, rightOk := CheckRBProperty(right)

	if leftOk == nil && rightOk == nil {
		if leftHeight != rightHeight {
			return 0, errors.ErrViolatedRBProperty
		}
		if root.IsRed() {
			return leftHeight, nil
		}
		return leftHeight + 1, nil
	}
	return 0, errors.ErrViolatedRBProperty
}
