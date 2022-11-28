package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type AndersonNode struct {
	Value int
	Left  *AndersonNode
	Right *AndersonNode
	Size  uint32
	Level uint32
}

func NewNode(value int, level uint32) *AndersonNode {
	return &AndersonNode{
		Value: value,
		Left:  nil,
		Right: nil,
		Size:  uint32(1),
		Level: level,
	}
}

func (root *AndersonNode) Update() {
	root.Size = uint32(1)
	if root.Left != nil {
		root.Size += root.Left.Size
	}
	if root.Right != nil {
		root.Size += root.Right.Size
	}
}

func (root *AndersonNode) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *AndersonNode) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *AndersonNode) SetChild(child *AndersonNode, direction bool) {
	if direction {
		root.Right = child
	} else {
		root.Left = child
	}
}

func (root *AndersonNode) Child(direction bool) *AndersonNode {
	if direction {
		return root.Right
	}
	return root.Left
}

type Anderson struct {
	Root *AndersonNode
}

func New() Anderson {
	return Anderson{Root: nil}
}

func LeftRotate(root *AndersonNode) *AndersonNode {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *AndersonNode) *AndersonNode {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Skew(root *AndersonNode) *AndersonNode {
	// Print(root)
	if root.Left == nil || root.Left.Level != root.Level {
		return root
	}
	return RightRotate(root)
}

func Split(root *AndersonNode) *AndersonNode {
	if root.Right == nil || root.Right.Right == nil || root.Right.Right.Level != root.Level {
		return root
	}
	root = LeftRotate(root)
	root.Level += 1
	return root
}

func Insert(root *AndersonNode, value int) *AndersonNode {
	if root == nil {
		return NewNode(value, 1)
	}
	if value < root.Value {
		root.Left = Insert(root.Left, value)
	} else {
		root.Right = Insert(root.Right, value)
	}
	root.Update()
	root = Skew(root)
	root = Split(root)
	return root
}

func Delete(root *AndersonNode, value int) *AndersonNode {
	if root == nil {
		return nil
	}
	if value < root.Value {
		root.Left = Delete(root.Left, value)
	} else if value > root.Value {
		root.Right = Delete(root.Right, value)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		} else {
			minNode := Kth(root.Right, 1)
			root.Value = minNode.Value
			root.Right = Delete(root.Right, minNode.Value)
		}
	}
	root.Update()
	if (root.Left != nil && root.Left.Level < root.Level-1) ||
		(root.Right != nil && root.Right.Level < root.Level-1) {
		root.Level -= 1
		if root.Right != nil && root.Right.Level > root.Level {
			root.Right.Level = root.Level
		}
		root = Skew(root)
		root = Split(root)
	}
	return root
}

func Kth(root *AndersonNode, k uint32) *AndersonNode {
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

func (tree *Anderson) Insert(value int) {
	tree.Root = Insert(tree.Root, value)
}

func (tree *Anderson) Delete(value int) {
	tree.Root = Delete(tree.Root, value)
}

func (tree *Anderson) Kth(k uint32) (int, error) {
	root := Kth(tree.Root, k)
	if root == nil {
		return int(rune(0)), errors.New("k is out of range")
	}
	return root.Value, nil
}

func (tree *Anderson) Size() uint32 {
	if tree.Root == nil {
		return 0
	}
	return tree.Root.Size
}

func (tree *Anderson) Empty() bool {
	return tree.Root == nil
}

func (tree *Anderson) Clear() {
	tree.Root = nil
}

func Find(root *AndersonNode, value int) *AndersonNode {
	for root != nil {
		if value < root.Value {
			root = root.Left
		} else if root.Value < value {
			root = root.Right
		} else {
			return root
		}
	}
	return nil
}

func (tree *Anderson) Contains(value int) bool {
	return Find(tree.Root, value) != nil
}

func Rank(root *AndersonNode, value int) uint32 {
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

func (thisTree *Anderson) Rank(value int) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev(root *AndersonNode, value int) *AndersonNode {
	var prev *AndersonNode = nil
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

func (thisTree *Anderson) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("no prev")
	}
	return prev.Value, nil
}

func Next(root *AndersonNode, value int) *AndersonNode {
	var next *AndersonNode = nil
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

func (thisTree *Anderson) Next(value int) (int, error) {
	prev := Next(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("no next")
	}
	return prev.Value, nil
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
