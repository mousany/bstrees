package splay

import "golang.org/x/exp/constraints"

func leftRotate[T constraints.Ordered](root *splayNode[T]) *splayNode[T] {
	right := root.right
	root.setChild(right.left, true)
	right.setChild(root, false)
	root.update()
	right.update()
	return right
}

func rightRotate[T constraints.Ordered](root *splayNode[T]) *splayNode[T] {
	left := root.left
	root.setChild(left.right, false)
	left.setChild(root, true)
	root.update()
	left.update()
	return left
}

// Rotate root to its parent
// After this operation, parent will be the child of root
func rotateToParent[T constraints.Ordered](root *splayNode[T]) {
	grandParent := root.parent.parent
	if root == root.parent.left {
		// root is left child
		root = rightRotate(root.parent)
	} else {
		// root is right child
		root = leftRotate(root.parent)
	}
	if grandParent != nil {
		if grandParent.left == root.parent {
			grandParent.setChild(root, false)
			grandParent.update()
		} else {
			grandParent.setChild(root, true)
			grandParent.update()
		}
	}
}

// Rotate root to target
// After this operation, target will be the child of root
func splayRotate[T constraints.Ordered](root, target *splayNode[T]) {
	targetParent := target.parent
	for root.parent != targetParent {
		parent := root.parent
		grandParent := parent.parent
		direction := root == parent.left
		grandDirection := parent == grandParent.left
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
