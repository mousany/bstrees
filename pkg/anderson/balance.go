package anderson

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

func skew[T node.Ordered](root *andersonTreeNode[T]) node.Noded[T] {
	if root.Left().(*andersonTreeNode[T]) == nil || root.Left().(*andersonTreeNode[T]).Level() != root.Level() {
		return root
	}
	return tree.SingleRotate(true, node.Noded[T](root))
}

func split[T node.Ordered](root *andersonTreeNode[T]) node.Noded[T] {
	if root.Right().(*andersonTreeNode[T]) == nil ||
		root.Right().Right().(*andersonTreeNode[T]) == nil ||
		root.Right().Right().(*andersonTreeNode[T]).Level() != root.Level() {
		return root
	}
	root = tree.SingleRotate(false, node.Noded[T](root)).(*andersonTreeNode[T])
	root.SetLevel(root.Level() + 1)
	return root
}
