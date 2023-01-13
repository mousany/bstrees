package tree

import "bstrees/internal/node"

func SingleRotate[T node.Ordered](direction bool, root node.Noded[T]) node.Noded[T] {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	return save
}

func LeftRotate[T node.Ordered](root node.Noded[T]) node.Noded[T] {
	return SingleRotate(true, root)
}

func RightRotate[T node.Ordered](root node.Noded[T]) node.Noded[T] {
	return SingleRotate(false, root)
}
