package splay

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
	"bstrees/pkg/errors"
)

// Rotate root to its parent
// After this operation, parent will be the child of root
func rotateToParent[T node.Ordered](root *splayTreeNode[T]) {
	grandParent := root.Parent().Parent()
	parentDirection := root == root.Parent().Left().(*splayTreeNode[T])
	root = tree.SingleRotate(parentDirection, node.Noded[T](root.Parent())).(*splayTreeNode[T])
	if grandParent != nil {
		grandParentDirection := root.Parent() == grandParent.Left().(*splayTreeNode[T])
		grandParent.SetChild(!grandParentDirection, node.Noded[T](root))
	}
}

// Rotate root to target
// After this operation, target will be the child of root
func splayRotate[T node.Ordered](root, target *splayTreeNode[T]) {
	targetParent := target.Parent()
	for root.Parent() != targetParent {
		parent := root.Parent()
		grandParent := parent.Parent()
		direction := root == parent.Left().(*splayTreeNode[T])
		grandDirection := parent == grandParent.Left().(*splayTreeNode[T])
		if parent == target {
			// root is the child of target
			rotateToParent(root)
		} else if direction == grandDirection {
			// zig-zig
			rotateToParent(parent)
			rotateToParent(root)
		} else {
			// zig-zag
			rotateToParent(root)
			rotateToParent(root)
		}
	}
}

func CheckBSTProperty[T node.Ordered](root node.Noded[T]) error {
	if root.IsNil() {
		return nil
	}
	left, right := root.Left(), root.Right()
	if !left.IsNil() && left.Value() > root.Value() {
		return errors.ErrViolatedBSTProperty
	}
	if !right.IsNil() && right.Value() < root.Value() {
		return errors.ErrViolatedBSTProperty
	}
	leftSize := uint(0)
	if !left.IsNil() {
		leftSize = left.Size()
	}
	rightSize := uint(0)
	if !right.IsNil() {
		rightSize = right.Size()
	}
	if leftSize+rightSize+root.(*splayTreeNode[T]).Rec() != root.Size() {
		return errors.ErrViolatedBSTProperty
	}

	if err := CheckBSTProperty(left); err != nil {
		return err
	}
	if err := CheckBSTProperty(right); err != nil {
		return err
	}
	return nil
}
