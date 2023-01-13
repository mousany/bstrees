package rb

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
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

// func PropertyCheck[T ordered.Ordered](root *node.RBNode[T]) (uint, error) {
// 	if root == nil {
// 		return uint(1), nil
// 	}
// 	left, right := root.Left, root.Right
// 	if root.Red() {
// 		if node.IsRed(left) || node.IsRed(right) {
// 			return 0, errors.ErrViolatedRedBlackTree
// 		}
// 	}
// 	leftHeight, leftOk := PropertyCheck(left)
// 	rightHeight, rightOk := PropertyCheck(right)

// 	if (left != nil && left.Value > root.Value) || (right != nil && right.Value < root.Value) {
// 		return 0, errors.ErrViolatedRedBlackTree
// 	}

// 	if leftOk == nil && rightOk == nil {
// 		if leftHeight != rightHeight {
// 			return 0, errors.ErrViolatedRedBlackTree
// 		}
// 		if root.Red() {
// 			return leftHeight, nil
// 		}
// 		return leftHeight + 1, nil
// 	}
// 	return 0, errors.ErrViolatedRedBlackTree
// }

// func (thisTree *RBTree[T]) PropertyCheck() error {
// 	_, err := PropertyCheck(thisTree.Root)
// 	return err
// }
