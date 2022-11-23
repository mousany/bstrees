package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type SplayNode struct {
	Value  int
	Left   *SplayNode
	Right  *SplayNode
	Parent *SplayNode
	Size   uint32
}

func NewNode(value int) *SplayNode {
	return &SplayNode{Value: value, Left: nil, Right: nil, Parent: nil, Size: 1}
}

func (root *SplayNode) Update() {
	root.Size = 1
	if root.Left != nil {
		root.Size += root.Left.Size
	}
	if root.Right != nil {
		root.Size += root.Right.Size
	}
}

func (root *SplayNode) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *SplayNode) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *SplayNode) SetChild(child *SplayNode, direction bool) {
	if direction {
		root.Right = child
	} else {
		root.Left = child
	}
	if child != nil {
		child.Parent = root
	}
}

func (root *SplayNode) Child(direction bool) *SplayNode {
	if direction {
		return root.Right
	}
	return root.Left
}

type Splay struct {
	Root *SplayNode
}

func New() *Splay {
	return &Splay{Root: nil}
}

func LeftRotate(root *SplayNode) *SplayNode {
	right := root.Right
	root.SetChild(right.Left, true)
	right.SetChild(root, false)
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *SplayNode) *SplayNode {
	left := root.Left
	root.SetChild(left.Right, false)
	left.SetChild(root, true)
	root.Update()
	left.Update()
	return left
}

// Rotate root to its parent
// After this operation, parent will be the child of root
func RotateToParent(root *SplayNode) {
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
func SplayRotate(root, target *SplayNode) {
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

func Find(root *SplayNode, value int) *SplayNode {
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

func Kth(root *SplayNode, k uint32) *SplayNode {
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

func Insert(root *SplayNode, value int) *SplayNode {
	if root == nil {
		return NewNode(value)
	} else {
		superRoot := NewNode(int(rune(0)))
		superRoot.SetChild(root, true)

		for p := root; p != nil; {
			p.Size += 1
			if value < p.Value {
				if p.Left == nil {
					p.SetChild(NewNode(value), false)
					SplayRotate(p.Left, root)
					break
				} else {
					p = p.Left
				}
			} else {
				if p.Right == nil {
					p.SetChild(NewNode(value), true)
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

func Delete(root *SplayNode, value int) *SplayNode {
	if root == nil {
		return nil
	}
	superRoot := NewNode(int(rune(0)))
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

func (thisTree *Splay) Insert(value int) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func (thisTree *Splay) Delete(value int) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *Splay) Contains(value int) bool {
	return Find(thisTree.Root, value) != nil
}

func (thisTree *Splay) Kth(k uint32) (int, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return int(rune(0)), errors.New("k is out of range")
	}
	return result.Value, nil
}

func (thisTree *Splay) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *Splay) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *Splay) Clear() {
	thisTree.Root = nil
}

func (thisTree *Splay) Rank(value int) uint32 {
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

func Prev(root *SplayNode, value int) *SplayNode {
	var result *SplayNode
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

func (thisTree *Splay) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("no prev value")
	}
	return prev.Value, nil
}

func Next(root *SplayNode, value int) *SplayNode {
	var result *SplayNode
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

func (thisTree *Splay) Next(value int) (int, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return int(rune(0)), errors.New("no next value")
	}
	return next.Value, nil
}

func Read(istream *bufio.Reader) (int, error) {
	res, sign := int(0), 1
	readed := false
	c, err := istream.ReadByte()
	for ; err == nil && (c < '0' || c > '9'); c, err = istream.ReadByte() {
		if c == '-' {
			sign = -1
		}
	}
	for ; err == nil && c >= '0' && c <= '9'; c, err = istream.ReadByte() {
		readed = true
		res = res<<3 + res<<1 + int(c-'0')
	}
	if !readed {
		return 0, err
	}
	return res * int(sign), err
}

func ReadWithPanic(gin *bufio.Reader) int {
	value, err := Read(gin)
	if err != nil {
		panic(err)
	}
	return value
}

func main() {
	tree := New()
	gin := bufio.NewReader(os.Stdin)
	n := ReadWithPanic(gin)
	for i := 0; i < n; i++ {
		opt := ReadWithPanic(gin)
		value := ReadWithPanic(gin)
		switch opt {
		case 1:
			tree.Insert(value)
		case 2:
			tree.Delete(value)
		case 3:
			fmt.Println(tree.Rank(value))
		case 4:
			{
				kth, err := tree.Kth(uint32(value))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 5:
			{
				kth, err := tree.Prev(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		case 6:
			{
				kth, err := tree.Next(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(kth)
			}
		}
	}
}
