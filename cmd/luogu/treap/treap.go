package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
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

func newBaseTree() *BaseTree {
	return &BaseTree{}
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

func SingleRotate(direction bool, root Noded) Noded {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	return save
}

func LeftRotate(root Noded) Noded {
	return SingleRotate(true, root)
}

func RightRotate(root Noded) Noded {
	return SingleRotate(false, root)
}

type treapTreeNode struct {
	*BaseTreeNode
	weight uint // Random weight
}

func newTreapTreeNode(value int) *treapTreeNode {
	n := &treapTreeNode{BaseTreeNode: newBaseTreeNode(value), weight: uint(rand.Uint32())}
	n.SetLeft((*treapTreeNode)(nil))
	n.SetRight((*treapTreeNode)(nil))
	return n
}

func (n *treapTreeNode) Weight() uint {
	return n.weight
}

func (n *treapTreeNode) IsNil() bool {
	return n == nil
}

type TreapTree struct {
	*BaseTree
}

func New() *TreapTree {
	tree := &TreapTree{newBaseTree()}
	tree.SetRoot((*treapTreeNode)(nil))
	return tree
}

func insert(root *treapTreeNode, value int) Noded {
	if root == nil {
		return newTreapTreeNode(value)
	}
	if root.Value() <= value {
		root.SetRight(insert(root.Right().(*treapTreeNode), value))
		if root.Right().(*treapTreeNode).Weight() < root.Weight() {
			root = SingleRotate(false, Noded(root)).(*treapTreeNode)
		}
	} else {
		root.SetLeft(insert(root.Left().(*treapTreeNode), value))
		if root.Left().(*treapTreeNode).Weight() < root.Weight() {
			root = SingleRotate(true, Noded(root)).(*treapTreeNode)
		}
	}
	root.Update()
	return root
}

func (tree *TreapTree) Insert(value int) {
	tree.SetRoot(insert(tree.Root().(*treapTreeNode), value))
}

func delete(root *treapTreeNode, value int) Noded {
	if root == nil {
		return nil
	}
	if root.Value() == value {
		if root.Left().IsNil() {
			return root.Right()
		}
		if root.Right().IsNil() {
			return root.Left()
		}
		if root.Left().(*treapTreeNode).Weight() < root.Right().(*treapTreeNode).Weight() {
			root = SingleRotate(true, Noded(root)).(*treapTreeNode)
			root.SetRight(delete(root.Right().(*treapTreeNode), value))
		} else {
			root = SingleRotate(false, Noded(root)).(*treapTreeNode)
			root.SetLeft(delete(root.Left().(*treapTreeNode), value))
		}
	} else if root.Value() < value {
		root.SetRight(delete(root.Right().(*treapTreeNode), value))
	} else {
		root.SetLeft(delete(root.Left().(*treapTreeNode), value))
	}
	root.Update()
	return root
}

func (tree *TreapTree) Delete(value int) {
	tree.SetRoot(delete(tree.Root().(*treapTreeNode), value))
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
