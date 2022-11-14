package rbtree

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/rbtree/node"
	"bstrees/pkg/trait/ordered"
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

func Reorient[T ordered.Ordered](grandpa, father, me *node.RBNode[T]) (*node.RBNode[T], *node.RBNode[T]) {
	if grandpa.Left == father {
		grandpa.Color = node.Red
		if father.Right == me {
			father = LeftRotate(father)
		}
		father.Color = node.Black
		_ = RightRotate(grandpa)
	} else {
		grandpa.Color = node.Red
		if father.Left == me {
			father = RightRotate(father)
		}
		father.Color = node.Black
		_ = LeftRotate(grandpa)
	}
	return father, me
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
	grandpa_ptr, father_ptr, me_ptr := (**node.RBNode[T])(nil), (**node.RBNode[T])(nil), &root
	for *me_ptr != nil {
		if (*me_ptr).Full() && (*me_ptr).Left.Red() && (*me_ptr).Right.Red() {
			FlipColor(*me_ptr)
		}
		if grandpa_ptr != nil && father_ptr != nil && (*father_ptr).Red() && (*me_ptr).Red() {
			father, me := Reorient(*grandpa_ptr, *father_ptr, *me_ptr)
			*grandpa_ptr = father
			father_ptr = grandpa_ptr
			if father.Left == me {
				me_ptr = &father.Left
			} else {
				me_ptr = &father.Right
			}
		}
		grandpa_ptr = father_ptr
		father_ptr = me_ptr
		if value < (*me_ptr).Value {
			me_ptr = &(*me_ptr).Left
		} else {
			me_ptr = &(*me_ptr).Right
		}
	}
	*me_ptr = node.New(value)
	if grandpa_ptr != nil && father_ptr != nil {
		// fmt.Println((*grandpa_ptr).Value, (*grandpa_ptr).Color, (*father_ptr).Value, (*father_ptr).Color, (*me_ptr).Value, (*me_ptr).Color)
		if (*father_ptr).Red() && (*me_ptr).Red() {
			father, _ := Reorient(*grandpa_ptr, *father_ptr, *me_ptr)
			*grandpa_ptr = father
		}
	}
	if root.Red() {
		root.Color = node.Black
	}
	return root
}

func (tree *RBTree[T]) Insert(value T) {
	tree.Root = Insert(tree.Root, value)
}

func Delete[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	grandpa_ptr, father_ptr, me_ptr := (**node.RBNode[T])(nil), (**node.RBNode[T])(nil), &root
	for *me_ptr != nil && (*me_ptr).Value != value {
		if (*me_ptr).Full() && (*me_ptr).Left.Red() && (*me_ptr).Right.Red() {
			FlipColor(*me_ptr)
		}
		if grandpa_ptr != nil && father_ptr != nil && (*father_ptr).Red() && (*me_ptr).Red() {
			father, me := Reorient(*grandpa_ptr, *father_ptr, *me_ptr)
			*grandpa_ptr = father
			father_ptr = grandpa_ptr
			if father.Left == me {
				me_ptr = &father.Left
			} else {
				me_ptr = &father.Right
			}
		}
		grandpa_ptr = father_ptr
		father_ptr = me_ptr
		if value < (*me_ptr).Value {
			me_ptr = &(*me_ptr).Left
		} else {
			me_ptr = &(*me_ptr).Right
		}
	}
	if *me_ptr != nil {
		if (*me_ptr).Left == nil && (*me_ptr).Right == nil {
			*me_ptr = nil
		} else if (*me_ptr).Left == nil {
			*me_ptr = (*me_ptr).Right
			(*me_ptr).Color = node.Black
		} else if (*me_ptr).Right == nil {
			*me_ptr = (*me_ptr).Left
			(*me_ptr).Color = node.Black
		} else {
			// find the min of right subtree
			min, _ := Kth((*me_ptr).Right, 1) // guaranteed to be not nil
			(*me_ptr).Value = min
			(*me_ptr).Right = Delete((*me_ptr).Right, min)
		}
	}
	return root
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

func Print[T ordered.Ordered](root *node.RBNode[T]) {
	if root == nil {
		return
	}
	color := "red"
	if root.Color == node.Black {
		color = "black"
	}
	fmt.Println(root.Value, color)
	fmt.Println("Left: ")
	Print(root.Left)
	fmt.Println("Right: ")
	Print(root.Right)
}

func (thisTree *RBTree[T]) Print() {
	Print(thisTree.Root)
}
