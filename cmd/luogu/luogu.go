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

func Reorient(grandpa, father, me *RBNode) (*RBNode, *RBNode) {
	if grandpa.Left == father {
		grandpa.Color = Red
		father.Color = Black
		if father.Right == me {
			father = LeftRotate(father)
		}
		_ = RightRotate(grandpa)
	} else {
		grandpa.Color = Red
		father.Color = Black
		if father.Left == me {
			father = RightRotate(father)
		}
		_ = LeftRotate(grandpa)
	}
	return father, me
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
	grandpa_ptr, father_ptr, me_ptr := (**RBNode)(nil), (**RBNode)(nil), &root
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
	*me_ptr = NewNode(value)
	if grandpa_ptr != nil && father_ptr != nil {
		// fmt.Println((*grandpa_ptr).Value, (*grandpa_ptr).Color, (*father_ptr).Value, (*father_ptr).Color, (*me_ptr).Value, (*me_ptr).Color)
		if (*father_ptr).Red() && (*me_ptr).Red() {
			father, _ := Reorient(*grandpa_ptr, *father_ptr, *me_ptr)
			*grandpa_ptr = father
		}
	}
	if root.Red() {
		root.Color = Black
	}
	return root
}

func (tree *RBTree) Insert(value int) {
	tree.Root = Insert(tree.Root, value)
}

func Delete(root *RBNode, value int) *RBNode {
	grandpa_ptr, father_ptr, me_ptr := (**RBNode)(nil), (**RBNode)(nil), &root
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
			(*me_ptr).Color = Black
		} else if (*me_ptr).Right == nil {
			*me_ptr = (*me_ptr).Left
			(*me_ptr).Color = Black
		} else {
			// find the min of right subtree
			min, _ := Kth((*me_ptr).Right, 1) // guaranteed to be not nil
			(*me_ptr).Value = min
			(*me_ptr).Right = Delete((*me_ptr).Right, min)
		}
	}
	return root
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
