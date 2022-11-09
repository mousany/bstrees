package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type TreapNode struct {
	Value  int
	Left   *TreapNode
	Right  *TreapNode
	Weight uint32
	Size   uint32
}

func NewNode(value int) *TreapNode {
	return &TreapNode{Value: value, Left: nil, Right: nil, Weight: rand.Uint32(), Size: 1}
}

func (thisTree *TreapNode) Update() {
	thisTree.Size = 1
	if thisTree.Left != nil {
		thisTree.Size += thisTree.Left.Size
	}
	if thisTree.Right != nil {
		thisTree.Size += thisTree.Right.Size
	}
}

type VanillaTreap struct {
	Root *TreapNode
}

func New() *VanillaTreap {
	return &VanillaTreap{Root: nil}
}

func LeftRotate(root *TreapNode) *TreapNode {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *TreapNode) *TreapNode {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Kth(root *TreapNode, k uint32) *TreapNode {
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

func Insert(root *TreapNode, value int) *TreapNode {
	if root == nil {
		return NewNode(value)
	}
	if root.Value <= value {
		root.Right = Insert(root.Right, value)
		if root.Right.Weight < root.Weight {
			root = LeftRotate(root)
		}
	} else {
		root.Left = Insert(root.Left, value)
		if root.Left.Weight < root.Weight {
			root = RightRotate(root)
		}
	}
	root.Update()
	return root
}

func (thisTree *VanillaTreap) Insert(value int) {
	thisTree.Root = Insert(thisTree.Root, value)
}

func Delete(root *TreapNode, value int) *TreapNode {
	if root == nil {
		return nil
	}
	if root.Value == value {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		if root.Left.Weight < root.Right.Weight {
			root = RightRotate(root)
			root.Right = Delete(root.Right, value)
		} else {
			root = LeftRotate(root)
			root.Left = Delete(root.Left, value)
		}
	} else if root.Value < value {
		root.Right = Delete(root.Right, value)
	} else {
		root.Left = Delete(root.Left, value)
	}
	root.Update()
	return root
}

func (thisTree *VanillaTreap) Delete(value int) {
	thisTree.Root = Delete(thisTree.Root, value)
}

func (thisTree *VanillaTreap) Kth(k uint32) (int, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return int(0), errors.New("k is out of range")
	}
	return result.Value, nil
}

func (thisTree *VanillaTreap) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *VanillaTreap) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *VanillaTreap) Clear() {
	thisTree.Root = nil
}

func Rank(root *TreapNode, value int) uint32 {
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

func (thisTree *VanillaTreap) Rank(value int) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev(root *TreapNode, value int) *TreapNode {
	var result *TreapNode = nil
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

func (thisTree *VanillaTreap) Prev(value int) (int, error) {
	result := Prev(thisTree.Root, value)
	if result == nil {
		return int(0), errors.New("no prev value")
	}
	return result.Value, nil
}

func Next(root *TreapNode, value int) *TreapNode {
	var result *TreapNode = nil
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

func (thisTree *VanillaTreap) Next(value int) (int, error) {
	result := Next(thisTree.Root, value)
	if result == nil {
		return int(0), errors.New("no next value")
	}
	return result.Value, nil
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
