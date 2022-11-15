package rbtree

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/rbtree/node"
	"bstrees/pkg/trait/ordered"
	"bytes"
	"fmt"
)

type RBTree[T ordered.Ordered] struct {
	Root *node.RBNode[T]
}

func New[T ordered.Ordered]() *RBTree[T] {
	return &RBTree[T]{Root: nil}
}

func LeftRotate[T ordered.Ordered](root *node.RBNode[T]) *node.RBNode[T] {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate[T ordered.Ordered](root *node.RBNode[T]) *node.RBNode[T] {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Kth[T ordered.Ordered](root *node.RBNode[T], k uint32) (T, error) {
	for root != nil {
		leftSize := uint32(0)
		if root.Left != nil {
			leftSize = root.Left.Size
		}
		if leftSize+1 == k {
			return root.Value, nil
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return T(rune(0)), errors.ErrOutOfRange
}

func Reorient[T ordered.Ordered](charles, william, louis *node.RBNode[T]) *node.RBNode[T] {
	if william == charles.Left {
		charles.Color = node.Red
		if william.Right == louis {
			charles.Left = LeftRotate(william)
		}
		charles.Left.Color = node.Black
		return RightRotate(charles)
	} else {
		charles.Color = node.Red
		if william.Left == louis {
			charles.Right = RightRotate(william)
		}
		charles.Right.Color = node.Black
		return LeftRotate(charles)
	}
}

func FlipColor[T ordered.Ordered](root *node.RBNode[T]) {
	if root != nil {
		root.Color = node.RBColor(!root.Color)
		if root.Left != nil {
			root.Left.Color = node.RBColor(!root.Left.Color)
		}
		if root.Right != nil {
			root.Right.Color = node.RBColor(!root.Right.Color)
		}
	}
}

func Insert[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var header *node.RBNode[T] = node.New(T(rune(0)))
	header.Right = root
	header.Update()

	var elizabeth, charles *node.RBNode[T]
	var william *node.RBNode[T] = header
	var louis *node.RBNode[T] = root
	// See Queen Elizabeth II's family tree for reference.
	for louis != nil {
		if louis.Full() && louis.Left.Red() && louis.Right.Red() {
			FlipColor(louis)
		} // Rotation cannot happen on level 2, so not need to check elizabeth.
		if william != nil && charles != nil && louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
				elizabeth.Update()
				if elizabeth.Left == louis {
					if value < louis.Value {
						louis, william, charles = william, louis, elizabeth
					} else {
						louis, william, charles = charles, louis, elizabeth
					}
				} else {
					charles = elizabeth
				}
			} else {
				elizabeth.Right = Reorient(charles, william, louis)
				elizabeth.Update()
				if elizabeth.Right == louis {
					if value < louis.Value {
						louis, william, charles = charles, louis, elizabeth
					} else {
						louis, william, charles = william, louis, elizabeth
					}
				} else {
					charles = elizabeth
				}
			}
		}
		elizabeth, charles, william = charles, william, louis
		if value < louis.Value {
			louis = louis.Left
		} else {
			louis = louis.Right
		}
	}
	if charles == nil { // The tree is empty.
		header.Right = node.New(value)
		header.Right.Color = node.Black
		return header.Right
	}
	// fmt.Println(william, charles, elizabeth)
	louis = node.New(value)
	if value < william.Value {
		william.Left = louis
		william.Update()
		// Rotation cannot happen on level 2, so not need to check elizabeth.
		if louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
				elizabeth.Update()
			} else {
				elizabeth.Right = Reorient(charles, william, louis)
				elizabeth.Update()
			}
		}
	} else {
		william.Right = louis
		william.Update()
		// Rotation cannot happen on level 2, so not need to check elizabeth.
		if louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
				elizabeth.Update()
			} else {
				elizabeth.Right = Reorient(charles, william, louis)
				elizabeth.Update()
			}
		}
	}

	for louis != nil {
		louis.Update()
		louis = louis.Father
	}

	if header.Right.Red() {
		header.Right.Color = node.Black
	}
	header.Right.Father = nil
	return header.Right
}

func (tree *RBTree[T]) Insert(value T) {
	tree.Root = Insert(tree.Root, value)
}

func Sibling[T ordered.Ordered](william, louis *node.RBNode[T]) *node.RBNode[T] {
	if william.Left == louis {
		return william.Right
	}
	return william.Left
}

func WeakBlack[T ordered.Ordered](root *node.RBNode[T]) bool {
	return root == nil || root.Black()
}

func StrongRed[T ordered.Ordered](root *node.RBNode[T]) bool {
	return root != nil && root.Red()
}

func Delete[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var header *node.RBNode[T] = node.New(T(rune(0)))
	header.Right = root
	header.Update()

	var charles *node.RBNode[T]
	var william *node.RBNode[T] = header
	var louis *node.RBNode[T] = root
	// See Queen Elizabeth II's family tree for reference.

	isCase2B := false

	for louis != nil && louis.Value != value {
		if !isCase2B {
			if charles == nil { // Step 1, Louis is the root.
				if WeakBlack(louis.Left) && WeakBlack(louis.Right) {
					// Case 1-A, Louis is the root and has two black children.
					louis.Color = node.Red
				} else { // Case 1-B, Louis is the root and has one red child.
					isCase2B = true
				}
			} else { // Step 2, the main case.
				if WeakBlack(louis.Left) && WeakBlack(louis.Right) {
					// Case 2-A, Louis has two black children.
					charlotte := Sibling(charles, william)
					if WeakBlack(charlotte.Left) && WeakBlack(charlotte.Right) {
						// Case 2-A1, Charlotte has two black children.
						FlipColor(william)
					} else if StrongRed(charlotte.Left) {
						// Case 2-A2, Charlotte has one red child.
						louis.Color = node.Red
						if charles.Left == william {
							charles.Left = Reorient(william, charlotte, charlotte.Left)
							charles.Update()
							FlipColor(charles.Left)
						} else {
							charles.Right = Reorient(william, charlotte, charlotte.Left)
							charles.Update()
							FlipColor(charles.Right)
						}
					} else if StrongRed(charlotte.Right) {
						// Case 2-A3, Charlotte has one red child.
						louis.Color = node.Red
						if charles.Left == william {
							charles.Left = Reorient(william, charlotte, charlotte.Right)
							charles.Update()
							FlipColor(charles.Left)
						} else {
							charles.Right = Reorient(william, charlotte, charlotte.Right)
							charles.Update()
							FlipColor(charles.Right)
						}
					}
				} else { // Case 2-B, Louis has one red child.
					isCase2B = true
				}
			}
		} else { // Case 2B: Louis has at least one red child.
			if louis.Black() { // If louis is black, then perform a rotation.
				charlette := Sibling(charles, william)
				if charles.Left == william {
					if charlette == william.Left {
						charles.Left = RightRotate(william)
						charles.Update()
					} else {
						charles.Left = LeftRotate(william)
						charles.Update()
					}
				} else {
					if charlette == william.Left {
						charles.Right = RightRotate(william)
						charles.Update()
					} else {
						charles.Right = LeftRotate(william)
						charles.Update()
					}
				}
				charlette.Color = node.Black
				william.Color = node.Red
				louis, william = william, charlette
				isCase2B = false
			}
		}

		charles, william = william, louis
		if value < louis.Value {
			louis = louis.Left
		} else {
			louis = louis.Right
		}
	}
	if louis != nil {
		if louis.Leaf() {
			if william.Left == louis {
				william.Left = nil
			} else {
				william.Right = nil
			}
		} else if louis.Left != nil && louis.Right == nil {
			leftMax, _ := Kth(louis.Left, louis.Left.Size)
			louis.Value = leftMax
			louis.Left = Delete(louis.Left, leftMax)
			louis.Update()
		} else {
			rightMin, _ := Kth(louis.Right, 1)
			louis.Value = rightMin
			louis.Right = Delete(louis.Right, rightMin)
			louis.Update()
		}
	}

	for louis != nil {
		louis.Update()
		louis = louis.Father
	}

	if header.Right != nil {
		if header.Right.Red() {
			header.Right.Color = node.Black
		}
		header.Right.Father = nil
	}

	return header.Right
}

func (tree *RBTree[T]) Delete(value T) {
	tree.Root = Delete(tree.Root, value)
}

func (thisTree *RBTree[T]) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *RBTree[T]) Kth(k uint32) (T, error) {
	return Kth(thisTree.Root, k)
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

func Print[T ordered.Ordered](root *node.RBNode[T]) string {
	if root == nil {
		return "null"
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[%v", root.Value))
	if root.Red() {
		buffer.WriteString("1, ")
	} else {
		buffer.WriteString("0, ")
	}
	// buffer.WriteString(fmt.Sprintf("%v, ", root.Size))
	buffer.WriteString(fmt.Sprintf("%v, ", Print(root.Left)))
	buffer.WriteString(fmt.Sprintf("%v]", Print(root.Right)))
	return buffer.String()
}

func (thisTree *RBTree[T]) Print() {
	fmt.Println(Print(thisTree.Root))
}
