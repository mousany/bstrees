package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type RBColor bool

const (
	Red   RBColor = true
	Black RBColor = false
)

type RBNode struct {
	Value int
	Left  *RBNode
	Right *RBNode
	Color RBColor
	Size  uint32 // Size of subtree, unnecessary if you don'int need kth element
	// Father *RBNode // Not necessary, but easier to implement
}

func NewNode(value int) *RBNode {
	return &RBNode{Value: value, Left: nil, Right: nil, Color: Red, Size: 1}
}

func (thisNode *RBNode) Update() {
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Size += thisNode.Right.Size
	}
}

func (thisNode *RBNode) Red() bool {
	return thisNode.Color == Red
}

func (thisNode *RBNode) Black() bool {
	return thisNode.Color == Black
}

func (thisNode *RBNode) Full() bool {
	return thisNode.Left != nil && thisNode.Right != nil
}

type RBTree struct {
	Root *RBNode
}

func New() *RBTree {
	return &RBTree{Root: nil}
}

func LeftRotate(root *RBNode) *RBNode {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *RBNode) *RBNode {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Kth(root *RBNode, k uint32) (int, error) {
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
	return int(rune(0)), errors.New("k is out of range")
}

func Reorient(charles, william, louis *RBNode) *RBNode {
	if william == charles.Left {
		charles.Color = Red
		if william.Right == louis {
			charles.Left = LeftRotate(william)
		}
		charles.Left.Color = Black
		return RightRotate(charles)
	} else {
		charles.Color = Red
		if william.Left == louis {
			charles.Right = RightRotate(william)
		}
		charles.Right.Color = Black
		return LeftRotate(charles)
	}
}

func FlipColor(root *RBNode) {
	if root != nil {
		root.Color = RBColor(!root.Color)
		if root.Left != nil {
			root.Left.Color = RBColor(!root.Left.Color)
		}
		if root.Right != nil {
			root.Right.Color = RBColor(!root.Right.Color)
		}
	}
}

func Insert(root *RBNode, value int) *RBNode {
	var header *RBNode = NewNode(int(rune(0)))
	header.Right = root

	var elizabeth, charles *RBNode
	var william *RBNode = header
	var louis *RBNode = root
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
		header.Right = NewNode(value)
		header.Right.Color = Black
		return header.Right
	}
	// fmt.Println(william, charles, elizabeth)
	louis = NewNode(value)
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
		header.Right.Color = Black
	}
	return header.Right
}

func (tree *RBTree) Insert(value int) {
	tree.Root = Insert(tree.Root, value)
}

func Delete(root *RBNode, value int) *RBNode {
	var header *RBNode = NewNode(int(rune(0)))
	header.Right = root

	var elizabeth, charles *RBNode
	var william *RBNode = header
	var louis *RBNode = root
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
	return header.Right
}

func (tree *RBTree) Delete(value int) {
	tree.Root = Delete(tree.Root, value)
}

func (thisTree *RBTree) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *RBTree) Kth(k uint32) (int, error) {
	return Kth(thisTree.Root, k)
}

func (thisTree *RBTree) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *RBTree) Clear() {
	thisTree.Root = nil
}

func Rank(root *RBNode, value int) uint32 {
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

func (thisTree *RBTree) Rank(value int) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev(root *RBNode, value int) *RBNode {
	var prev *RBNode = nil
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

func (thisTree *RBTree) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("No previous value")
	}
	return prev.Value, nil
}

func Next(root *RBNode, value int) *RBNode {
	var next *RBNode = nil
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

func (thisTree *RBTree) Next(value int) (int, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return int(rune(0)), errors.New("No next value")
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
