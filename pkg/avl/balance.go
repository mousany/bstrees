package avltree

import (
	"bstrees/internal/math"
	"bstrees/internal/node"
)

func singleRotate[T node.Ordered](direction bool, root *avlTreeNode[T]) node.Noded[T] {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	return save
}

func balance[T node.Ordered](root *avlTreeNode[T]) node.Noded[T] {
	leftHeight := root.Left().(*avlTreeNode[T]).Height()
	rightHeight := root.Right().(*avlTreeNode[T]).Height()
	if math.Abs(leftHeight-rightHeight) > 1 {
		grandFatherDirection := leftHeight < rightHeight+1
		father := root.Child(grandFatherDirection).(*avlTreeNode[T])
		fatherLeftHeight := father.Left().(*avlTreeNode[T]).Height()
		fatherRightHeight := father.Right().(*avlTreeNode[T]).Height()
		fatherDirection := fatherLeftHeight < fatherRightHeight+1
		if grandFatherDirection != fatherDirection {
			root.SetChild(grandFatherDirection, singleRotate(!fatherDirection, father))
		}
		return singleRotate(!grandFatherDirection, root)
	}
	return root
}
