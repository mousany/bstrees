package rb

import (
	"github.com/yanglinshu/bstrees/v2/internal/order"
	"github.com/yanglinshu/bstrees/v2/pkg/errors"
)

type RBTree[T order.Ordered] struct {
	root *rbTreeNode[T]
}

func New[T order.Ordered]() *RBTree[T] {
	return &RBTree[T]{root: nil}
}

func kth[T order.Ordered](root *rbTreeNode[T], k uint) *rbTreeNode[T] {
	for root != nil {
		leftSize := uint(0)
		if root.left != nil {
			leftSize = root.left.size
		}
		if leftSize+1 == k {
			return root
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.right
		} else {
			root = root.left
		}
	}
	return nil
}

// https://archive.ph/EJTsz, Eternally Confuzzled's Blog
func insert[T order.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	if root == nil {
		root = newRBTreeNode(value)
	} else {
		superRoot := newRBTreeNode(T(rune(0))) // Head in Eternally Confuzzled's paper
		superRoot.right = root

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
				parent.setChild(direction, child)
				ok = true
			} else {
				// Update size
				child.size += 1
				if isRed(child.left) && isRed(child.right) {
					// Color flip
					child.color = red
					child.left.color = black
					child.right.color = black
				}
			}

			if isRed(child) && isRed(parent) {
				// Fix red violation
				direction2 := greatGrandParent.right == grandParent
				if child == parent.child(lastDirection) {
					greatGrandParent.setChild(direction2, singleRotate(grandParent, !lastDirection))
					// When performing a single rotation to grandparent, child is not affected.
					// So when grandparent(old) and parent(old) is updated, there are all +1ed.
				} else {
					greatGrandParent.setChild(direction2, doubleRotate(grandParent, !lastDirection))
					if !ok {
						// When performing a double rotation to grandparent, child is affected.
						// So we need to update child(now grandParent)'s size. But there is no need we insert is done.
						greatGrandParent.child(direction2).size += 1
					}
				}
			}

			lastDirection = direction
			direction = child.value < value
			if grandParent != nil {
				greatGrandParent = grandParent
			}

			grandParent = parent
			parent = child
			child = child.child(direction)
		}

		// Update root
		root = superRoot.right
	}
	root.color = black
	return root
}

func (t *RBTree[T]) Insert(value T) {
	t.root = insert(t.root, value)
}

func find[T order.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	for root != nil {
		if root.value == value {
			return root
		} else if root.value < value {
			root = root.right
		} else {
			root = root.left
		}
	}
	return nil
}

func (t *RBTree[T]) Contains(value T) bool {
	return find(t.root, value) != nil
}

func delete[T order.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	if root == nil || find(root, value) == nil {
		return root
	}
	superRoot := newRBTreeNode(T(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.right = root

	var child *rbTreeNode[T] = superRoot // Q in Eternally Confuzzled's paper
	var parent *rbTreeNode[T] = nil      // P in Eternally Confuzzled's paper
	var grandParent *rbTreeNode[T] = nil // G in Eternally Confuzzled's paper
	var target *rbTreeNode[T] = nil      // F in Eternally Confuzzled's paper
	direction := true

	// Search and push a red down
	for child.child(direction) != nil {
		lastDirection := direction

		grandParent = parent
		parent = child
		child = child.child(direction)
		direction = child.value < value

		// Update size
		child.size -= 1

		// Save the target node
		if child.value == value {
			target = child
		}

		// Push the red node down
		if !isRed(child) && !isRed(child.child(direction)) {
			if isRed(child.child(!direction)) {
				parent.setChild(lastDirection, singleRotate(child, direction))
				parent = parent.child(lastDirection)

				// When performing a single rotation to child, child is affected.
				// So we need to update child and sibling(now parent)'s size.
				child.size -= 1
				parent.Update()
			} else if !isRed(child.child(!direction)) {
				sibling := parent.child(!lastDirection)
				if sibling != nil {
					if !isRed(sibling.child(!lastDirection)) && !isRed(sibling.child(lastDirection)) {
						// Color flip
						parent.color = black
						sibling.color = red
						child.color = red
					} else {
						direction2 := grandParent.right == parent
						if isRed(sibling.child(lastDirection)) {
							grandParent.setChild(direction2, doubleRotate(parent, lastDirection))
						} else if isRed(sibling.child(!lastDirection)) {
							grandParent.setChild(direction2, singleRotate(parent, lastDirection))
						}

						// When performing a rotation to parent, child is not affected.
						// So all nodes on the path are -1ed.

						// // Update Size
						// parent.Update()
						// grandParent.Child(direction2).Update()

						// Ensure correct coloring
						child.color = red
						grandParent.child(direction2).color = red
						grandParent.child(direction2).left.color = black
						grandParent.child(direction2).right.color = black
					}
				}
			}
		}
	}

	// Replace and remove the target node
	if target != nil {
		target.value = child.value
		parent.setChild(parent.right == child, child.child(child.left == nil))
	}

	// Update root and make it black
	root = superRoot.right
	if root != nil {
		root.color = black
	}
	return root
}

func (t *RBTree[T]) Delete(value T) {
	t.root = delete(t.root, value)
}

func (t *RBTree[T]) Size() uint {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func (t *RBTree[T]) Kth(k uint) (T, error) {
	result := kth(t.root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.value, nil
}

func (t *RBTree[T]) Empty() bool {
	return t.root == nil
}

func (t *RBTree[T]) Clear() {
	t.root = nil
}

func rank[T order.Ordered](root *rbTreeNode[T], value T) uint {
	rank := uint(0)
	for root != nil {
		if root.value < value {
			rank += 1
			if root.left != nil {
				rank += root.left.size
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return rank + 1
}

func (t *RBTree[T]) Rank(value T) uint {
	return rank(t.root, value)
}

func prev[T order.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	var prev *rbTreeNode[T] = nil
	for root != nil {
		if root.value < value {
			prev = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return prev
}

func (t *RBTree[T]) Prev(value T) (T, error) {
	prev := prev(t.root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.value, nil
}

func next[T order.Ordered](root *rbTreeNode[T], value T) *rbTreeNode[T] {
	var next *rbTreeNode[T] = nil
	for root != nil {
		if root.value > value {
			next = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return next
}

func (t *RBTree[T]) Next(value T) (T, error) {
	next := next(t.root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.value, nil
}
