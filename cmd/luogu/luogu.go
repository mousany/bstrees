package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type RotatelessTreap struct {
	Root *TreapNode
}

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

func (this *TreapNode) Update() {
	this.Size = 1
	if this.Left != nil {
		this.Size += this.Left.Size
	}
	if this.Right != nil {
		this.Size += this.Right.Size
	}
}

func New() *RotatelessTreap {
	return &RotatelessTreap{Root: nil}
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

func Split(node *TreapNode, key int) (*TreapNode, *TreapNode) {
	if node == nil {
		return nil, nil
	}
	if node.Value <= key {
		left, right := Split(node.Right, key)
		node.Right = left
		node.Update()
		return node, right
	} else {
		left, right := Split(node.Left, key)
		node.Left = right
		node.Update()
		return left, node
	}
}

func Kth(node *TreapNode, k uint32) *TreapNode {
	for node != nil {
		leftSize := uint32(0)
		if node.Left != nil {
			leftSize = node.Left.Size
		}
		if leftSize+1 == k {
			return node
		} else if leftSize+1 < k {
			k -= leftSize + 1
			node = node.Right
		} else {
			node = node.Left
		}
	}
	return nil
}

func (this *RotatelessTreap) Insert(value int) {
	left, right := Split(this.Root, value)
	this.Root = Merge(Merge(left, NewNode(value)), right)
}

func (this *RotatelessTreap) Delete(value int) {
	left, right := Split(this.Root, value)
	left, mid := Split(left, value-1)
	if mid != nil {
		mid = Merge(mid.Left, mid.Right)
	}
	this.Root = Merge(Merge(left, mid), right)
}

func (this *RotatelessTreap) Rank(value int) uint32 {
	left, right := Split(this.Root, value-1)
	defer func() {
		this.Root = Merge(left, right)
	}()
	if left == nil {
		return 1
	}
	return left.Size + 1
}

func (this *RotatelessTreap) Kth(k uint32) *int {
	node := Kth(this.Root, k)
	if node == nil {
		return nil
	}
	return &node.Value
}

func (this *RotatelessTreap) Size() uint32 {
	if this.Root == nil {
		return 0
	}
	return this.Root.Size
}

func (this *RotatelessTreap) Empty() bool {
	return this.Root == nil
}

func (this *RotatelessTreap) Clear() {
	this.Root = nil
}

func (this *RotatelessTreap) Prev(value int) *int {
	left, right := Split(this.Root, value-1)
	defer func() {
		this.Root = Merge(left, right)
	}()
	result := Kth(left, left.Size)
	if result == nil {
		return nil
	}
	return &result.Value
}

func (this *RotatelessTreap) Next(value int) *int {
	left, right := Split(this.Root, value)
	defer func() {
		this.Root = Merge(left, right)
	}()
	result := Kth(right, 1)
	if result == nil {
		return nil
	}
	return &result.Value
}

func read(istream *bufio.Reader) (int, error) {
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
	n, err := read(gin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		opt, err := read(gin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value, err := read(gin)
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
			fmt.Println(*tree.Kth(uint32(value)))
		case 5:
			fmt.Println(*tree.Prev(value))
		case 6:
			fmt.Println(*tree.Next(value))
		}
	}

}
