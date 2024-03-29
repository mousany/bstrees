package scapegoat

import "golang.org/x/exp/constraints"

func toSlice[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T]) []*scapeGoatTreeNode[T] {
	if root == nil {
		return []*scapeGoatTreeNode[T]{}
	}
	if root.active() {
		defer func() {
			root.left = nil
			root.right = nil
			root.size = 1
			root.weight = 1
		}()
		return append(append(toSlice(root.left), root), toSlice(root.right)...)
	} else {
		return append(toSlice(root.left), toSlice(root.right)...)
	}
}

func fromSlice[T constraints.Integer | constraints.Float](slice []*scapeGoatTreeNode[T]) *scapeGoatTreeNode[T] {
	if len(slice) == 0 {
		return nil
	}
	mid := len(slice) / 2
	root := slice[mid]
	root.left = fromSlice(slice[:mid])
	root.right = fromSlice(slice[mid+1:])
	root.update()
	return root
}

func reconstruct[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T]) *scapeGoatTreeNode[T] {
	return fromSlice(toSlice(root))
}

func imbalance[T constraints.Integer | constraints.Float](root *scapeGoatTreeNode[T], alpha float64) bool {
	if root == nil {
		return false
	}
	if root.left != nil && root.left.weight > uint(alpha*float64(root.weight)) {
		return true
	}
	if root.right != nil && root.right.weight > uint(alpha*float64(root.weight)) {
		return true
	}
	return false
}
