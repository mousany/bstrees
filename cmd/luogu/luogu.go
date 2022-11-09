package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func Min(args ...int32) int32 {
	min := args[0]
	for _, arg := range args {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Max(args ...int32) int32 {
	max := args[0]
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}

type AVLNode struct {
	Value  int
	Left   *AVLNode
	Right  *AVLNode
	Height int32  // Height of the node
	Size   uint32 // Size of subtree, unnecessary if you don'int need kth element
}

func NewNode(value int) *AVLNode {
	return &AVLNode{Value: value, Left: nil, Right: nil, Height: 0, Size: 1}
}

func (thisNode *AVLNode) Update() {
	thisNode.Height = 0
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Height = Max(thisNode.Height, thisNode.Left.Height+1)
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Height = Max(thisNode.Height, thisNode.Right.Height+1)
		thisNode.Size += thisNode.Right.Size
	}
}

type AVL struct {
	Root *AVLNode
}

func New() *AVL {
	return &AVL{Root: nil}
}

func LeftRotate(root *AVLNode) *AVLNode {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *AVLNode) *AVLNode {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Balance(root *AVLNode) *AVLNode {
	leftHeight := int32(-1)
	if root.Left != nil {
		leftHeight = root.Left.Height
	}
	rightHeight := int32(-1)
	if root.Right != nil {
		rightHeight = root.Right.Height
	}
	if leftHeight > rightHeight+1 {
		left := root.Left
		leftLeftHeight := int32(-1)
		if left.Left != nil {
			leftLeftHeight = left.Left.Height
		}
		leftRightHeight := int32(-1)
		if left.Right != nil {
			leftRightHeight = left.Right.Height
		}
		if leftLeftHeight < leftRightHeight {
			root.Left = LeftRotate(left)
		}
		ret := RightRotate(root)
		return ret
	} else if rightHeight > leftHeight+1 {
		right := root.Right
		rightLeftHeight := int32(-1)
		if right.Left != nil {
			rightLeftHeight = right.Left.Height
		}
		rightRightHeight := int32(-1)
		if right.Right != nil {
			rightRightHeight = right.Right.Height
		}
		if rightRightHeight < rightLeftHeight {
			root.Right = RightRotate(right)
		}
		return LeftRotate(root)
	}
	return root
}

func Kth(root *AVLNode, k uint32) (int, error) {
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

func Insert(root *AVLNode, value int) *AVLNode {
	if root == nil {
		return NewNode(value)
	}
	if value < root.Value {
		root.Left = Insert(root.Left, value)
	} else {
		root.Right = Insert(root.Right, value)
	}
	root.Update()
	return Balance(root)
}

func (thisTree *AVL) Insert(value int) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func Delete(root *AVLNode, value int) *AVLNode {
	if root == nil {
		return nil
	}
	if value < root.Value {
		root.Left = Delete(root.Left, value)
	} else if root.Value < value {
		root.Right = Delete(root.Right, value)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		} else {
			min, _ := Kth(root.Right, 1) // root.Right is not nil, so this will not fail
			root.Value = min
			root.Right = Delete(root.Right, min)
		}
	}
	root.Update()
	return Balance(root)
}

func (thisTree *AVL) Delete(value int) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *AVL) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *AVL) Height() int32 {
	if thisTree.Root == nil {
		return -1
	}
	return thisTree.Root.Height
}

func (thisTree *AVL) Kth(k uint32) (int, error) {
	return Kth(thisTree.Root, k)
}

func (thisTree *AVL) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *AVL) Clear() {
	thisTree.Root = nil
}

func Rank(root *AVLNode, value int) uint32 {
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

func (thisTree *AVL) Rank(value int) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev(root *AVLNode, value int) *AVLNode {
	var result *AVLNode = nil
	for root != nil {
		if root.Value < value {
			result = root
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return result
}

func (thisTree *AVL) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("no previous value")
	}
	return prev.Value, nil
}

func Next(root *AVLNode, value int) *AVLNode {
	var result *AVLNode = nil
	for root != nil {
		if root.Value > value {
			result = root
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return result
}

func (thisTree *AVL) Next(value int) (int, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return int(rune(0)), errors.New("no next value")
	}
	return next.Value, nil
}

func Read(istream *bufio.Reader) (int, error) {
	res, sign := 0, 1
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

func main() {
	tree := New()
	gin := bufio.NewReader(os.Stdin)
	n, err := Read(gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		opt, err := Read(gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value, err := Read(gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
