package splay

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/splay/node"
	"bstrees/pkg/trait/ordered"
)

type Splay[T ordered.Ordered] struct {
	Root *node.SplayNode[T]
}

func New[T ordered.Ordered]() *Splay[T] {
	return &Splay[T]{Root: nil}
}

func LeftRotate[T ordered.Ordered](root *node.SplayNode[T]) *node.SplayNode[T] {
	right := root.Right
	root.SetChild(right.Left, true)
	right.SetChild(root, false)
	root.Update()
	right.Update()
	return right
}

func RightRotate[T ordered.Ordered](root *node.SplayNode[T]) *node.SplayNode[T] {
	left := root.Left
	root.SetChild(left.Right, false)
	left.SetChild(root, true)
	root.Update()
	left.Update()
	return left
}

// Rotate root to its parent
// After this operation, parent will be the child of root
func RotateToParent[T ordered.Ordered](root *node.SplayNode[T]) {
	grandParent := root.Parent.Parent
	if root == root.Parent.Left {
		// root is left child
		root = RightRotate(root.Parent)
	} else {
		// root is right child
		root = LeftRotate(root.Parent)
	}
	if grandParent != nil {
		if grandParent.Left == root.Parent {
			grandParent.SetChild(root, false)
			grandParent.Update()
		} else {
			grandParent.SetChild(root, true)
			grandParent.Update()
		}
	}
}

// Rotate root to target
// After this operation, target will be the child of root
func SplayRotate[T ordered.Ordered](root, target *node.SplayNode[T]) {
	targetParent := target.Parent
	for root.Parent != targetParent {
		parent := root.Parent
		grandParent := parent.Parent
		direction := root == parent.Left
		grandDirection := parent == grandParent.Left
		if parent == target {
			// root is the child of target
			RotateToParent(root)
		} else if direction == grandDirection {
			// zig-zig
			RotateToParent(parent)
			RotateToParent(root)
		} else {
			// zig-zag
			RotateToParent(root)
			RotateToParent(root)
		}
	}
}

func Find[T ordered.Ordered](root *node.SplayNode[T], value T) *node.SplayNode[T] {
	for p := root; p != nil; {
		if p.Value == value {
			return p
		} else if value < p.Value {
			p = p.Left
		} else {
			p = p.Right
		}
	}
	return nil
}

func Kth[T ordered.Ordered](root *node.SplayNode[T], k uint32) *node.SplayNode[T] {
	for p := root; p != nil; {
		leftSize := uint32(0)
		if p.Left != nil {
			leftSize = p.Left.Size
		}
		if leftSize+1 == k {
			// SplayRotate(p, root)
			return p
		} else if leftSize+1 < k {
			k -= leftSize + 1
			p = p.Right
		} else {
			p = p.Left
		}
	}
	return nil
}

func Insert[T ordered.Ordered](root *node.SplayNode[T], value T) *node.SplayNode[T] {
	if root == nil {
		return node.New(value)
	} else {
		superRoot := node.New(T(rune(0)))
		superRoot.SetChild(root, true)

		for p := root; p != nil; {
			p.Size += 1
			if value < p.Value {
				if p.Left == nil {
					p.SetChild(node.New(value), false)
					SplayRotate(p.Left, root)
					break
				} else {
					p = p.Left
				}
			} else {
				if p.Right == nil {
					p.SetChild(node.New(value), true)
					SplayRotate(p.Right, root)
					break
				} else {
					p = p.Right
				}
			}
		}

		superRoot.Right.Parent = nil
		return superRoot.Right
	}
}

func Delete[T ordered.Ordered](root *node.SplayNode[T], value T) *node.SplayNode[T] {
	if root == nil {
		return nil
	}
	superRoot := node.New(T(rune(0)))
	superRoot.SetChild(root, true)
	p := Find(root, value)
	if p == nil {
		root.Parent = nil
		return root
	}
	SplayRotate(p, root)
	if p.Left == nil && p.Right == nil {
		superRoot.SetChild(nil, true)
	} else if p.Left == nil {
		superRoot.SetChild(p.Right, true)
	} else if p.Right == nil {
		superRoot.SetChild(p.Left, true)
	} else {
		maxLeft := p.Left
		for maxLeft.Right != nil {
			maxLeft.Size -= 1
			maxLeft = maxLeft.Right
		}
		SplayRotate(maxLeft, superRoot.Right)
		maxLeft.SetChild(p.Right, true)
		superRoot.SetChild(maxLeft, true)
		superRoot.Right.Update()
	}
	superRoot.Right.Parent = nil
	return superRoot.Right
}

func (thisTree *Splay[T]) Insert(value T) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func (thisTree *Splay[T]) Delete(value T) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *Splay[T]) Contains(value T) bool {
	return Find(thisTree.Root, value) != nil
}

func (thisTree *Splay[T]) Kth(k uint32) (T, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.Value, nil
}

func (thisTree *Splay[T]) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *Splay[T]) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *Splay[T]) Clear() {
	thisTree.Root = nil
}

func (thisTree *Splay[T]) Rank(value T) uint32 {
	rank := uint32(0)
	for p := thisTree.Root; p != nil; {
		if value < p.Value {
			p = p.Left
		} else if value > p.Value {
			rank += 1
			if p.Left != nil {
				rank += p.Left.Size
			}
			p = p.Right
		} else {
			if p.Left != nil {
				rank += p.Left.Size
			}
			break
		}
	}
	return rank + 1
}

func Prev[T ordered.Ordered](root *node.SplayNode[T], value T) *node.SplayNode[T] {
	var result *node.SplayNode[T]
	for p := root; p != nil; {
		if value > p.Value {
			result = p
			p = p.Right
		} else {
			p = p.Left
		}
	}
	return result
}

func (thisTree *Splay[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T ordered.Ordered](root *node.SplayNode[T], value T) *node.SplayNode[T] {
	var result *node.SplayNode[T]
	for p := root; p != nil; {
		if value < p.Value {
			result = p
			p = p.Left
		} else {
			p = p.Right
		}
	}
	return result
}

func (thisTree *Splay[T]) Next(value T) (T, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value, nil
}
