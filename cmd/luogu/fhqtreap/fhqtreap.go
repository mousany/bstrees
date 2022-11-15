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
	Weight uint32 // Random weight
	Size   uint32 // Size of subtree, unnecessary if you don'int need kth element
}

func NewNode(value int) *TreapNode {
	return &TreapNode{Value: value, Left: nil, Right: nil, Weight: rand.Uint32(), Size: 1}
}

func (thisNode *TreapNode) Update() {
	thisNode.Size = 1
	if thisNode.Left != nil {
		thisNode.Size += thisNode.Left.Size
	}
	if thisNode.Right != nil {
		thisNode.Size += thisNode.Right.Size
	}
}

type FHQTreap struct {
	Root *TreapNode
}

func New() *FHQTreap {
	return &FHQTreap{Root: nil}
}

func Merge(left *TreapNode, right *TreapNode) *TreapNode {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.Weight < right.Weight {
		left.Right = Merge(left.Right, right)
		left.Update()
		return left
	} else {
		right.Left = Merge(left, right.Left)
		right.Update()
		return right
	}
}

func Split(root *TreapNode, key int) (*TreapNode, *TreapNode) {
	if root == nil {
		return nil, nil
	}
	if root.Value <= key {
		left, right := Split(root.Right, key)
		root.Right = left
		root.Update()
		return root, right
	} else {
		left, right := Split(root.Left, key)
		root.Left = right
		root.Update()
		return left, root
	}
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

func (thisTree *FHQTreap) Insert(value int) {
	left, right := Split(thisTree.Root, value)
	thisTree.Root = Merge(Merge(left, NewNode(value)), right)
}

func (thisTree *FHQTreap) Delete(value int) {
	left, right := Split(thisTree.Root, value)
	left, mid := Split(left, value-1)
	if mid != nil {
		mid = Merge(mid.Left, mid.Right)
	}
	thisTree.Root = Merge(Merge(left, mid), right)
}

func (thisTree *FHQTreap) Rank(value int) uint32 {
	left, right := Split(thisTree.Root, value-1)
	defer func() {
		thisTree.Root = Merge(left, right)
	}()
	if left == nil {
		return 1
	}
	return left.Size + 1
}

func (thisTree *FHQTreap) Kth(k uint32) (int, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return int(0), errors.New("out of range")
	}
	return result.Value, nil
}

func (thisTree *FHQTreap) Size() uint32 {
	if thisTree.Root == nil {
		return 0
	}
	return thisTree.Root.Size
}

func (thisTree *FHQTreap) Empty() bool {
	return thisTree.Root == nil
}

func (thisTree *FHQTreap) Clear() {
	thisTree.Root = nil
}

func (thisTree *FHQTreap) Prev(value int) (int, error) {
	left, right := Split(thisTree.Root, value-1)
	defer func() {
		thisTree.Root = Merge(left, right)
	}()
	result := Kth(left, left.Size)
	if result == nil {
		return int(0), errors.New("no previous value")
	}
	return result.Value, nil
}

func (thisTree *FHQTreap) Next(value int) (int, error) {
	left, right := Split(thisTree.Root, value)
	defer func() {
		thisTree.Root = Merge(left, right)
	}()
	result := Kth(right, 1)
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
	ans, last := 0, 0
	tree := New()
	gin := bufio.NewReader(os.Stdin)
	n, err := Read(gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m, err := Read(gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		x, err := Read(gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tree.Insert(x)
	}
	for i := 0; i < m; i++ {
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
		value ^= last
		switch opt {
		case 1:
			tree.Insert(value)
		case 2:
			tree.Delete(value)
		case 3:
			{
				rank := tree.Rank(value)
				ans ^= int(rank)
				last = int(rank)
			}
		case 4:
			{
				kth, err := tree.Kth(uint32(value))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ans ^= kth
				last = kth
			}
		case 5:
			{
				prev, err := tree.Prev(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ans ^= prev
				last = prev
			}
		case 6:
			{
				next, err := tree.Next(value)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ans ^= next
				last = next
			}
		}
	}
	fmt.Println(ans)
}
