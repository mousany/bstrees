package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func Min(args ...int) int {
	min := args[0]
	for _, arg := range args {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Max(args ...int) int {
	max := args[0]
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}

func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

type Valued interface {
	Value() int
	SetValue(int)
}

type Sized interface {
	Size() uint
	SetSize(uint)
}

type Binary interface {
	Left() Noded
	SetLeft(Noded)
	Right() Noded
	SetRight(Noded)
	Child(bool) Noded
	SetChild(bool, Noded)
}

type Noded interface {
	Valued
	Sized
	Binary
	Update()
	IsNil() bool
}

type BaseTreeNode struct {
	value int
	size  uint
	left  Noded
	right Noded
}

func newBaseTreeNode(value int) *BaseTreeNode {
	return &BaseTreeNode{value, 1, nil, nil}
}

func (n *BaseTreeNode) Value() int {
	return n.value
}

func (n *BaseTreeNode) SetValue(value int) {
	n.value = value
}

func (n *BaseTreeNode) Size() uint {
	return n.size
}

func (n *BaseTreeNode) SetSize(size uint) {
	n.size = size
}

func (n *BaseTreeNode) Left() Noded {
	return n.left
}

func (n *BaseTreeNode) SetLeft(left Noded) {
	n.left = left
}

func (n *BaseTreeNode) Right() Noded {
	return n.right
}

func (n *BaseTreeNode) SetRight(right Noded) {
	n.right = right
}

func (n *BaseTreeNode) Child(right bool) Noded {
	if right {
		return n.right
	}
	return n.left
}

func (n *BaseTreeNode) SetChild(right bool, child Noded) {
	if right {
		n.right = child
	} else {
		n.left = child
	}
}

func (n *BaseTreeNode) IsNil() bool {
	return n == nil
}

func (n *BaseTreeNode) Update() {
	n.size = 1
	if !n.left.IsNil() {
		n.size += n.left.Size()
	}
	if !n.right.IsNil() {
		n.size += n.right.Size()
	}
}

type Rooted interface {
	Root() Noded
	SetRoot(Noded)
	Size() uint
	Empty() bool
	Clear()
}

type Insertable interface {
	Insert(int)
}

type Deleteable interface {
	Delete(int)
}

type Searchable interface {
	Contains(int) bool
	Kth(uint) (int, error)
	Rank(int) uint
	Prev(int) (int, error)
	Next(int) (int, error)
}

type Treed interface {
	Rooted
	Insertable
	Deleteable
	Searchable
}

type BaseTree struct {
	root Noded
}

func (tree *BaseTree) Root() Noded {
	return tree.root
}

func (tree *BaseTree) SetRoot(root Noded) {
	tree.root = root
}

func (tree *BaseTree) Size() uint {
	return tree.root.Size()
}

func (tree *BaseTree) Empty() bool {
	return tree.root == nil
}

func (tree *BaseTree) Clear() {
	tree.root = nil
}

func Find(root Noded, value int) Noded {
	for !root.IsNil() {
		if root.Value() == value {
			return root
		} else if root.Value() < value {
			root = root.Right()
		} else {
			root = root.Left()
		}
	}
	return nil
}

func (tree *BaseTree) Contains(value int) bool {
	return !Find(tree.root, value).IsNil()
}

func Kth(root Noded, k uint) Noded {
	for !root.IsNil() {
		leftSize := uint(0)
		if !root.Left().IsNil() {
			leftSize = root.Left().Size()
		}
		if leftSize+1 == k {
			return root
		} else if leftSize+1 < k {
			k -= leftSize + 1
			root = root.Right()
		} else {
			root = root.Left()
		}
	}
	return nil
}

func (tree *BaseTree) Kth(k uint) (int, error) {
	result := Kth(tree.root, k)
	if result.IsNil() {
		return int(rune(0)), errors.New("k is out of range")
	}
	return result.Value(), nil
}

func Rank(root Noded, value int) uint {
	rank := uint(0)
	for !root.IsNil() {
		if root.Value() < value {
			rank += 1
			if !root.Left().IsNil() {
				rank += root.Left().Size()
			}
			root = root.Right()
		} else {
			root = root.Left()
		}
	}
	return rank + 1
}

func (tree *BaseTree) Rank(value int) uint {
	return Rank(tree.root, value)
}

func Prev(root Noded, value int) Noded {
	var result Noded = nil
	for !root.IsNil() {
		if root.Value() < value {
			result = root
			root = root.Right()
		} else {
			root = root.Left()
		}
	}
	return result
}

func (tree *BaseTree) Prev(value int) (int, error) {
	prev := Prev(tree.root, value)
	if prev.IsNil() {
		return int(rune(0)), errors.New("no previous value")
	}
	return prev.Value(), nil
}

func Next(root Noded, value int) Noded {
	var result Noded = nil
	for !root.IsNil() {
		if root.Value() > value {
			result = root
			root = root.Left()
		} else {
			root = root.Right()
		}
	}
	return result
}

func (tree *BaseTree) Next(value int) (int, error) {
	next := Next(tree.root, value)
	if next.IsNil() {
		return int(rune(0)), errors.New("no next value")
	}
	return next.Value(), nil
}

func String(root Noded) string {
	if root.IsNil() {
		return "null"
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[%v, ", root.Value()))
	// buffer.WriteString(fmt.Sprintf("%v, ", root.Size()))
	buffer.WriteString(fmt.Sprintf("%v, ", String(root.Left())))
	buffer.WriteString(fmt.Sprintf("%v]", String(root.Right())))
	return buffer.String()
}

func (tree *BaseTree) String() string {
	return String(tree.root)
}

type avlTreeNode struct {
	*BaseTreeNode     // Embedded node
	height        int // Height of the node
}

func newAvlTreeNode(value int) *avlTreeNode {
	n := &avlTreeNode{newBaseTreeNode(value), 0}
	n.SetLeft((*avlTreeNode)(nil))
	n.SetRight((*avlTreeNode)(nil))
	return n
}

func (n *avlTreeNode) Height() int {
	if n == nil {
		return -1
	}
	return n.height
}

func (n *avlTreeNode) SetHeight(height int) {
	n.height = height
}

func (n *avlTreeNode) IsNil() bool {
	return n == nil
}

func (n *avlTreeNode) Update() {
	n.BaseTreeNode.Update()
	n.height = 0
	if !n.Left().IsNil() {
		n.height = Max(n.height, n.Left().(*avlTreeNode).Height()+1)
	}
	if !n.Right().IsNil() {
		n.height = Max(n.height, n.Right().(*avlTreeNode).Height()+1)
	}
}

type AVLTree struct {
	BaseTree
}

func New() *AVLTree {
	tree := &AVLTree{BaseTree{}}
	tree.SetRoot((*avlTreeNode)(nil))
	return tree
}

func Insert(root *avlTreeNode, value int) Noded {
	if root == nil {
		return newAvlTreeNode(value)
	}
	if value < root.Value() {
		root.SetLeft(Insert(root.Left().(*avlTreeNode), value))
	} else {
		root.SetRight(Insert(root.Right().(*avlTreeNode), value))
	}
	root.Update()
	return Balance(root)
}

func (tree *AVLTree) Insert(value int) {
	tree.SetRoot(Insert(tree.Root().(*avlTreeNode), value))
}

func Delete(root *avlTreeNode, value int) Noded {
	if root == nil {
		return nil
	}
	if value < root.Value() {
		root.SetLeft(Delete(root.Left().(*avlTreeNode), value))
	} else if root.Value() < value {
		root.SetRight(Delete(root.Right().(*avlTreeNode), value))
	} else {
		if root.Left().IsNil() {
			return root.Right().(*avlTreeNode)
		} else if root.Right().IsNil() {
			return root.Left().(*avlTreeNode)
		} else {
			minNode := Kth(root.Right(), 1).(*avlTreeNode) // root.Right is not nil, so this will not fail
			root.SetValue(minNode.Value())
			root.SetRight(Delete(root.Right().(*avlTreeNode), minNode.Value()))
		}
	}
	root.Update()
	return Balance(root)
}

func (tree *AVLTree) Delete(value int) {
	tree.SetRoot(Delete(tree.Root().(*avlTreeNode), value))
}

func SingleRotate(direction bool, root *avlTreeNode) Noded {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	return save
}

func Balance(root *avlTreeNode) Noded {
	leftHeight := root.Left().(*avlTreeNode).Height()
	rightHeight := root.Right().(*avlTreeNode).Height()
	if Abs(leftHeight-rightHeight) > 1 {
		grandFatherDirection := leftHeight < rightHeight+1
		father := root.Child(grandFatherDirection).(*avlTreeNode)
		fatherLeftHeight := father.Left().(*avlTreeNode).Height()
		fatherRightHeight := father.Right().(*avlTreeNode).Height()
		fatherDirection := fatherLeftHeight < fatherRightHeight+1
		if grandFatherDirection != fatherDirection {
			root.SetChild(grandFatherDirection, SingleRotate(!fatherDirection, father))
		}
		return SingleRotate(!grandFatherDirection, root)
	}
	return root
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
				kth, err := tree.Kth(uint(value))
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
