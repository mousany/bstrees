package anderson

import "bstrees/internal/node"

func singleRotate[T node.Ordered](direction bool, root *andersonTreeNode[T]) node.Noded[T] {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	return save
}

func skew[T node.Ordered](root *andersonTreeNode[T]) node.Noded[T] {
	if root.Left().(*andersonTreeNode[T]) == nil || root.Left().(*andersonTreeNode[T]).Level() != root.Level() {
		return root
	}
	return singleRotate(true, root)
}

func split[T node.Ordered](root *andersonTreeNode[T]) node.Noded[T] {
	if root.Right().(*andersonTreeNode[T]) == nil ||
		root.Right().Right().(*andersonTreeNode[T]) == nil ||
		root.Right().Right().(*andersonTreeNode[T]).Level() != root.Level() {
		return root
	}
	root = singleRotate(false, root).(*andersonTreeNode[T])
	root.SetLevel(root.Level() + 1)
	return root
}
