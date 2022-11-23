package rbtree

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/rbtree/node"
	"bstrees/pkg/trait/ordered"
)

type RBTree[T ordered.Ordered] struct {
	Root *node.RBNode[T]
}

func New[T ordered.Ordered]() *RBTree[T] {
	return &RBTree[T]{Root: nil}
}

func Kth[T ordered.Ordered](root *node.RBNode[T], k uint32) *node.RBNode[T] {
	for root != nil {
		leftSize := uint32(0)
		if root.Left != nil {
			leftSize = root.Left.Size
		}
		if leftSize+1 == k {
			return root
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return nil
}

func SingleRotate[T ordered.Ordered](root *node.RBNode[T], direction bool) *node.RBNode[T] {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	root.Color = node.Red
	save.Color = node.Black
	return save
}

func DoubleRotate[T ordered.Ordered](root *node.RBNode[T], direction bool) *node.RBNode[T] {
	root.SetChild(!direction, SingleRotate(root.Child(!direction), !direction))
	return SingleRotate(root, direction)
}

// https://archive.ph/EJTsz, Eternally Confuzzled's Blog
func Insert[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	if root == nil {
		root = node.New(value)
	} else {
		superRoot := node.New(T(rune(0))) // Head in Eternally Confuzzled's paper
		superRoot.Right = root

		var child *node.RBNode[T] = root                 // Q in Eternally Confuzzled's paper
		var parent *node.RBNode[T] = nil                 // P in Eternally Confuzzled's paper
		var grandParent *node.RBNode[T] = nil            // G in Eternally Confuzzled's paper
		var greatGrandParent *node.RBNode[T] = superRoot // T in Eternally Confuzzled's paper

		var direction bool = false
		var lastDirection bool = false

		// Search down
		for ok := false; !ok; {
			if child == nil {
				// Insert new node at the bottom
				child = node.New(value)
				parent.SetChild(direction, child)
				ok = true
			} else {
				// Update size
				child.Size += 1
				if node.IsRed(child.Left) && node.IsRed(child.Right) {
					// Color flip
					child.Color = node.Red
					child.Left.Color = node.Black
					child.Right.Color = node.Black
				}
			}

			if node.IsRed(child) && node.IsRed(parent) {
				// Fix red violation
				direction2 := greatGrandParent.Right == grandParent
				if child == parent.Child(lastDirection) {
					greatGrandParent.SetChild(direction2, SingleRotate(grandParent, !lastDirection))
					// When performing a single rotation to grandparent, child is not affected.
					// So when grandparent(old) and parent(old) is updated, there are all +1ed.
				} else {
					greatGrandParent.SetChild(direction2, DoubleRotate(grandParent, !lastDirection))
					if !ok {
						// When performing a double rotation to grandparent, child is affected.
						// So we need to update child(now grandParent)'s size. But there is no need we insert is done.
						greatGrandParent.Child(direction2).Size += 1
					}
				}
			}

			lastDirection = direction
			direction = child.Value < value
			if grandParent != nil {
				greatGrandParent = grandParent
			}

			grandParent = parent
			parent = child
			child = child.Child(direction)
		}

		// Update root
		root = superRoot.Right
	}
	root.Color = node.Black
	return root
}

func (thisTree *RBTree[T]) Insert(value T) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func Find[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	for root != nil {
		if root.Value == value {
			return root
		} else if root.Value < value {
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return nil
}

func (thisTree *RBTree[T]) Contains(value T) bool {
	return Find(thisTree.Root, value) != nil
}

func Delete[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	if root == nil || Find(root, value) == nil {
		return root
	}
	superRoot := node.New(T(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.Right = root

	var child *node.RBNode[T] = superRoot // Q in Eternally Confuzzled's paper
	var parent *node.RBNode[T] = nil      // P in Eternally Confuzzled's paper
	var grandParent *node.RBNode[T] = nil // G in Eternally Confuzzled's paper
	var target *node.RBNode[T] = nil      // F in Eternally Confuzzled's paper
	direction := true

	// Search and push a red down
	for child.Child(direction) != nil {
		lastDirection := direction

		grandParent = parent
		parent = child
		child = child.Child(direction)
		direction = child.Value < value

		// Update size
		child.Size -= 1

		// Save the target node
		if child.Value == value {
			target = child
		}

		// Push the red node down
		if !node.IsRed(child) && !node.IsRed(child.Child(direction)) {
			if node.IsRed(child.Child(!direction)) {
				parent.SetChild(lastDirection, SingleRotate(child, direction))
				parent = parent.Child(lastDirection)

				// When performing a single rotation to child, child is affected.
				// So we need to update child and sibling(now parent)'s size.
				child.Size -= 1
				parent.Update()
			} else if !node.IsRed(child.Child(!direction)) {
				sibling := parent.Child(!lastDirection)
				if sibling != nil {
					if !node.IsRed(sibling.Child(!lastDirection)) && !node.IsRed(sibling.Child(lastDirection)) {
						// Color flip
						parent.Color = node.Black
						sibling.Color = node.Red
						child.Color = node.Red
					} else {
						direction2 := grandParent.Right == parent
						if node.IsRed(sibling.Child(lastDirection)) {
							grandParent.SetChild(direction2, DoubleRotate(parent, lastDirection))
						} else if node.IsRed(sibling.Child(!lastDirection)) {
							grandParent.SetChild(direction2, SingleRotate(parent, lastDirection))
						}

						// When performing a rotation to parent, child is not affected.
						// So all nodes on the path are -1ed.

						// // Update Size
						// parent.Update()
						// grandParent.Child(direction2).Update()

						// Ensure correct coloring
						child.Color = node.Red
						grandParent.Child(direction2).Color = node.Red
						grandParent.Child(direction2).Left.Color = node.Black
						grandParent.Child(direction2).Right.Color = node.Black
					}
				}
			}
		}
	}

	// Replace and remove the target node
	if target != nil {
		target.Value = child.Value
		parent.SetChild(parent.Right == child, child.Child(child.Left == nil))
	}

	// Update root and make it black
	root = superRoot.Right
	if root != nil {
		root.Color = node.Black
	}
	return root
}

func (thisTree *RBTree[T]) Delete(value T) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *RBTree[T]) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *RBTree[T]) Kth(k uint32) (T, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.Value, nil
}

func (thisTree *RBTree[T]) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *RBTree[T]) Clear() {
	thisTree.Root = nil
}

func Rank[T ordered.Ordered](root *node.RBNode[T], value T) uint32 {
	rank := uint32(0)
	for root != nil {
		if root.Value < value {
			rank += 1
			if root.Left != nil {
				rank += root.Left.Size
			}
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return rank + 1
}

func (thisTree *RBTree[T]) Rank(value T) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var prev *node.RBNode[T] = nil
	for root != nil {
		if root.Value < value {
			prev = root
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return prev
}

func (thisTree *RBTree[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var next *node.RBNode[T] = nil
	for root != nil {
		if root.Value > value {
			next = root
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return next
}

func (thisTree *RBTree[T]) Next(value T) (T, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value, nil
}

// func Print[T ordered.Ordered](root *node.RBNode[T]) string {
// 	if root == nil {
// 		return "null"
// 	}
// 	var buffer bytes.Buffer
// 	buffer.WriteString(fmt.Sprintf("[%v", root.Value))
// 	// if root.Red() {
// 	// 	buffer.WriteString("1, ")
// 	// } else {
// 	// 	buffer.WriteString("0, ")
// 	// }
// 	buffer.WriteString(fmt.Sprintf(".%v, ", root.Size))
// 	// buffer.WriteString(fmt.Sprintf("%v, ", root.Size))
// 	buffer.WriteString(fmt.Sprintf("%v, ", Print(root.Left)))
// 	buffer.WriteString(fmt.Sprintf("%v]", Print(root.Right)))
// 	return buffer.String()
// }

// func (thisTree *RBTree[T]) Print() {
// 	fmt.Println(Print(thisTree.Root))
// }

// func PropertyCheck[T ordered.Ordered](root *node.RBNode[T]) (uint, error) {
// 	if root == nil {
// 		return uint(1), nil
// 	}
// 	left, right := root.Left, root.Right
// 	if root.Red() {
// 		if node.IsRed(left) || node.IsRed(right) {
// 			return 0, errors.ErrViolatedRedBlackTree
// 		}
// 	}
// 	leftHeight, leftOk := PropertyCheck(left)
// 	rightHeight, rightOk := PropertyCheck(right)

// 	if (left != nil && left.Value > root.Value) || (right != nil && right.Value < root.Value) {
// 		return 0, errors.ErrViolatedRedBlackTree
// 	}

// 	if leftOk == nil && rightOk == nil {
// 		if leftHeight != rightHeight {
// 			return 0, errors.ErrViolatedRedBlackTree
// 		}
// 		if root.Red() {
// 			return leftHeight, nil
// 		}
// 		return leftHeight + 1, nil
// 	}
// 	return 0, errors.ErrViolatedRedBlackTree
// }

// func (thisTree *RBTree[T]) PropertyCheck() error {
// 	_, err := PropertyCheck(thisTree.Root)
// 	return err
// }
