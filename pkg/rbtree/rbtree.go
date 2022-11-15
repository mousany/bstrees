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
		// Rotation cannot happen on level 2, so not need to check elizabeth.
		if louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
			} else {
				elizabeth.Right = Reorient(charles, william, louis)
			}
		}
	} else {
		william.Right = louis
		// Rotation cannot happen on level 2, so not need to check elizabeth.
		if louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
			} else {
				elizabeth.Right = Reorient(charles, william, louis)
			}
		}
	}
	if header.Right.Red() {
		header.Right.Color = node.Black
	}
	return header.Right
}

func (tree *RBTree[T]) Insert(value T) {
	tree.Root = Insert(tree.Root, value)
}

func Delete[T ordered.Ordered](root *node.RBNode[T], value T) *node.RBNode[T] {
	var header *node.RBNode[T] = node.New(T(rune(0)))
	header.Right = root

	var elizabeth, charles *node.RBNode[T]
	var william *node.RBNode[T] = header
	var louis *node.RBNode[T] = root
	// See Queen Elizabeth II's family tree for reference.
	for louis != nil && value != louis.Value {
		if louis.Full() && louis.Left.Red() && louis.Right.Red() {
			FlipColor(louis)
		} // Rotation cannot happen on level 2, so not need to check elizabeth.
		if william != nil && charles != nil && louis.Red() && william.Red() {
			if elizabeth.Left == charles {
				elizabeth.Left = Reorient(charles, william, louis)
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
	if louis == nil {
		if header.Right.Red() {
			header.Right.Color = node.Black
		}
		return header.Right
	} else {
		if louis.Full() {
			min, _ := Kth(louis.Right, 1) // guaranteed to exist
			louis.Value = min
			louis.Right = Delete(louis.Right, min)
		} else if louis.Left == nil {
			if louis.Right == nil {
				if william.Left == louis {
					william.Left = nil
				} else {
					william.Right = nil
				}
			} else {
				if william.Left == louis {
					william.Left = louis.Right
				} else {
					william.Right = louis.Right
				}
			}
		} else {
			if william.Left == louis {
				william.Left = louis.Left
			} else {
				william.Right = louis.Left
			}
		}
	}
	if header.Right.Red() {
		header.Right.Color = node.Black
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
	buffer.WriteString(fmt.Sprintf("%v, ", Print(root.Left)))
	buffer.WriteString(fmt.Sprintf("%v]", Print(root.Right)))
	return buffer.String()
}

func (thisTree *RBTree[T]) Print() {
	fmt.Println(Print(thisTree.Root))
}
