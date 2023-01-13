package avltree

import (
	"bstrees/internal/math"
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

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
			root.SetChild(grandFatherDirection, tree.SingleRotate(!fatherDirection, node.Noded[T](father)))
		}
		return tree.SingleRotate(grandFatherDirection, node.Noded[T](root))
	}
	return root
}
