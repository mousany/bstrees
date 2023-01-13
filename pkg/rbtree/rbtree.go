package rbtree

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

type RBTree[T node.Ordered] struct {
	*tree.BaseTree[T]
}

func New[T node.Ordered]() *RBTree[T] {
	tree := &RBTree[T]{tree.New[T]()}
	tree.SetRoot((*rbTreeNode[T])(nil))
	return tree
}

// https://archive.ph/EJTsz, Eternally Confuzzled's Blog
func Insert[T node.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	if root == nil {
		root = newRBTreeNode(value)
		root.SetColor(Black)
		return root
	}
	superRoot := newRBTreeNode(T(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.SetRight(root)

	var child *rbTreeNode[T] = root                 // Q in Eternally Confuzzled's paper
	var parent *rbTreeNode[T] = nil                 // P in Eternally Confuzzled's paper
	var grandParent *rbTreeNode[T] = nil            // G in Eternally Confuzzled's paper
	var greatGrandParent *rbTreeNode[T] = superRoot // T in Eternally Confuzzled's paper

	var direction bool = false
	var lastDirection bool = false

	// Search down
	for ok := false; !ok; {
		if child == nil {
			// Insert new node at the bottom
			child = newRBTreeNode(value)
			parent.SetChild(direction, child)
			ok = true
		} else {
			// Update size
			child.SetSize(child.Size() + 1)
			if child.Left().(*rbTreeNode[T]).IsRed() && child.Right().(*rbTreeNode[T]).IsRed() {
				// Color flip
				child.SetColor(Red)
				child.Left().(*rbTreeNode[T]).SetColor(Black)
				child.Right().(*rbTreeNode[T]).SetColor(Black)
			}
		}

		if child.IsRed() && parent.IsRed() {
			// Fix red violation
			direction2 := greatGrandParent.Right().(*rbTreeNode[T]) == grandParent
			if child == parent.Child(lastDirection).(*rbTreeNode[T]) {
				greatGrandParent.SetChild(direction2, SingleRotate(grandParent, !lastDirection))
				// When performing a single rotation to grandparent, child is not affected.
				// So when grandparent(old) and parent(old) is updated, there are all +1ed.
			} else {
				greatGrandParent.SetChild(direction2, DoubleRotate(grandParent, !lastDirection))
				if !ok {
					// When performing a double rotation to grandparent, child is affected.
					// So we need to update child(now grandParent)'s size. But there is no need we insert is done.
					greatGrandParent.Child(direction2).SetSize(greatGrandParent.Child(direction2).Size() + 1)
				}
			}
		}

		lastDirection = direction
		direction = child.Value() < value
		if grandParent != nil {
			greatGrandParent = grandParent
		}

		grandParent = parent
		parent = child
		child = child.Child(direction).(*rbTreeNode[T])
	}

	// Update root
	root = superRoot.Right().(*rbTreeNode[T])

	root.SetColor(Black)
	return root
}

func (tree *RBTree[T]) Insert(value T) {
	tree.SetRoot(Insert(tree.Root().(*rbTreeNode[T]), value))
}

func Delete[T node.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	if root == nil || tree.Find(node.Noded[T](root), value).IsNil() {
		return root
	}
	superRoot := newRBTreeNode(T(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.SetRight(root)

	var child *rbTreeNode[T] = superRoot // Q in Eternally Confuzzled's paper
	var parent *rbTreeNode[T] = nil      // P in Eternally Confuzzled's paper
	var grandParent *rbTreeNode[T] = nil // G in Eternally Confuzzled's paper
	var target *rbTreeNode[T] = nil      // F in Eternally Confuzzled's paper
	direction := true

	// Search and push a red down
	for !child.Child(direction).IsNil() {
		lastDirection := direction

		grandParent = parent
		parent = child
		child = child.Child(direction).(*rbTreeNode[T])
		direction = child.Value() < value

		// Update size
		child.SetSize(child.Size() - 1)

		// Save the target node
		if child.Value() == value {
			target = child
		}

		// Push the red node down
		if !child.IsRed() && !child.Child(direction).(*rbTreeNode[T]).IsRed() {
			if child.Child(!direction).(*rbTreeNode[T]).IsRed() {
				parent.SetChild(lastDirection, SingleRotate(child, direction))
				parent = parent.Child(lastDirection).(*rbTreeNode[T])

				// When performing a single rotation to child, child is affected.
				// So we need to update child and sibling(now parent)'s size.
				child.SetSize(child.Size() - 1)
				parent.Update()
			} else if !child.Child(!direction).(*rbTreeNode[T]).IsRed() {
				sibling := parent.Child(!lastDirection).(*rbTreeNode[T])
				if sibling != nil {
					if !sibling.Child(!lastDirection).(*rbTreeNode[T]).IsRed() && !sibling.Child(lastDirection).(*rbTreeNode[T]).IsRed() {
						// Color flip
						parent.SetColor(Black)
						sibling.SetColor(Red)
						child.SetColor(Red)
					} else {
						direction2 := grandParent.Right().(*rbTreeNode[T]) == parent
						if sibling.Child(lastDirection).(*rbTreeNode[T]).IsRed() {
							grandParent.SetChild(direction2, DoubleRotate(parent, lastDirection))
						} else if sibling.Child(!lastDirection).(*rbTreeNode[T]).IsRed() {
							grandParent.SetChild(direction2, SingleRotate(parent, lastDirection))
						}

						// When performing a rotation to parent, child is not affected.
						// So all nodes on the path are -1ed.

						// // Update Size
						// parent.Update()
						// grandParent.Child(direction2).Update()

						// Ensure correct coloring
						child.SetColor(Red)
						grandParent.Child(direction2).(*rbTreeNode[T]).SetColor(Red)
						grandParent.Child(direction2).Left().(*rbTreeNode[T]).SetColor(Black)
						grandParent.Child(direction2).Right().(*rbTreeNode[T]).SetColor(Black)
					}
				}
			}
		}
	}

	// Replace and remove the target node
	if target != nil {
		target.SetValue(child.Value())
		parent.SetChild(parent.Right().(*rbTreeNode[T]) == child, child.Child(child.Left().(*rbTreeNode[T]) == nil))
	}

	// Update root and make it black
	root = superRoot.Right().(*rbTreeNode[T])
	if root != nil {
		root.SetColor(Black)
	}
	return root
}

func (tree *RBTree[T]) Delete(value T) {
	tree.SetRoot(Delete(tree.Root().(*rbTreeNode[T]), value))
}
