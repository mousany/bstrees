package rbtree

import "bstrees/internal/node"

func SingleRotate[T node.Ordered](root *rbTreeNode[T], direction bool) node.Noded[T] {
	save := root.Child(!direction).(*rbTreeNode[T])
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	root.SetColor(Red)
	save.SetColor(Black)
	return save
}

func DoubleRotate[T node.Ordered](root *rbTreeNode[T], direction bool) node.Noded[T] {
	root.SetChild(!direction, SingleRotate(root.Child(!direction).(*rbTreeNode[T]), !direction))
	return SingleRotate(root, direction)
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
