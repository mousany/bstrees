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

type RBColor bool

const (
	Red   RBColor = true
	Black RBColor = false
)

type rbTreeNode struct {
	*BaseTreeNode
	color RBColor
}

func newRBTreeNode(value int) *rbTreeNode {
	n := &rbTreeNode{newBaseTreeNode(value), Red}
	n.SetLeft((*rbTreeNode)(nil))
	n.SetRight((*rbTreeNode)(nil))
	return n
}

func (n *rbTreeNode) IsNil() bool {
	return n == nil
}

func (n *rbTreeNode) Color() RBColor {
	return n.color
}

func (n *rbTreeNode) SetColor(color RBColor) {
	n.color = color
}

func (n *rbTreeNode) IsRed() bool {
	return n != nil && n.color == Red
}

func (n *rbTreeNode) IsBlack() bool {
	return n == nil || n.color == Black
}

type RBTree struct {
	*BaseTree
}

func New() *RBTree {
	tree := &RBTree{newBaseTree()}
	tree.SetRoot((*rbTreeNode)(nil))
	return tree
}

// https://archive.ph/EJTsz, Eternally Confuzzled's Blog
func Insert(root *rbTreeNode, value int) *rbTreeNode {
	if root == nil {
		root = newRBTreeNode(value)
		root.SetColor(Black)
		return root
	}
	superRoot := newRBTreeNode(int(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.SetRight(root)

	var child *rbTreeNode = root                 // Q in Eternally Confuzzled's paper
	var parent *rbTreeNode = nil                 // P in Eternally Confuzzled's paper
	var grandParent *rbTreeNode = nil            // G in Eternally Confuzzled's paper
	var greatGrandParent *rbTreeNode = superRoot // int in Eternally Confuzzled's paper

	var direction bool = false
	var lastDirection bool = false

	// Search down
	for ok := false; !ok; {
		if child == nil {
			// Insert new node at the bottom
			child = newRBTreeNode(value)
			parent.SetChild(direction, child)
			ok = true
		} else {
			// Update size
			child.SetSize(child.Size() + 1)
			if child.Left().(*rbTreeNode).IsRed() && child.Right().(*rbTreeNode).IsRed() {
				// Color flip
				child.SetColor(Red)
				child.Left().(*rbTreeNode).SetColor(Black)
				child.Right().(*rbTreeNode).SetColor(Black)
			}
		}

		if child.IsRed() && parent.IsRed() {
			// Fix red violation
			direction2 := greatGrandParent.Right().(*rbTreeNode) == grandParent
			if child == parent.Child(lastDirection).(*rbTreeNode) {
				greatGrandParent.SetChild(direction2, SingleRotate(grandParent, !lastDirection))
				// When performing a single rotation to grandparent, child is not affected.
				// So when grandparent(old) and parent(old) is updated, there are all +1ed.
			} else {
				greatGrandParent.SetChild(direction2, DoubleRotate(grandParent, !lastDirection))
				if !ok {
					// When performing a double rotation to grandparent, child is affected.
					// So we need to update child(now grandParent)'s size. But there is no need we insert is done.
					greatGrandParent.Child(direction2).SetSize(greatGrandParent.Child(direction2).Size() + 1)
				}
			}
		}

		lastDirection = direction
		direction = child.Value() < value
		if grandParent != nil {
			greatGrandParent = grandParent
		}

		grandParent = parent
		parent = child
		child = child.Child(direction).(*rbTreeNode)
	}

	// Update root
	root = superRoot.Right().(*rbTreeNode)

	root.SetColor(Black)
	return root
}

func (tree *RBTree) Insert(value int) {
	tree.SetRoot(Insert(tree.Root().(*rbTreeNode), value))
}

func Delete(root *rbTreeNode, value int) *rbTreeNode {
	if root == nil || Find(Noded(root), value).IsNil() {
		return root
	}
	superRoot := newRBTreeNode(int(rune(0))) // Head in Eternally Confuzzled's paper
	superRoot.SetRight(root)

	var child *rbTreeNode = superRoot // Q in Eternally Confuzzled's paper
	var parent *rbTreeNode = nil      // P in Eternally Confuzzled's paper
	var grandParent *rbTreeNode = nil // G in Eternally Confuzzled's paper
	var target *rbTreeNode = nil      // F in Eternally Confuzzled's paper
	direction := true

	// Search and push a red down
	for !child.Child(direction).IsNil() {
		lastDirection := direction

		grandParent = parent
		parent = child
		child = child.Child(direction).(*rbTreeNode)
		direction = child.Value() < value

		// Update size
		child.SetSize(child.Size() - 1)

		// Save the target node
		if child.Value() == value {
			target = child
		}

		// Push the red node down
		if !child.IsRed() && !child.Child(direction).(*rbTreeNode).IsRed() {
			if child.Child(!direction).(*rbTreeNode).IsRed() {
				parent.SetChild(lastDirection, SingleRotate(child, direction))
				parent = parent.Child(lastDirection).(*rbTreeNode)

				// When performing a single rotation to child, child is affected.
				// So we need to update child and sibling(now parent)'s size.
				child.SetSize(child.Size() - 1)
				parent.Update()
			} else if !child.Child(!direction).(*rbTreeNode).IsRed() {
				sibling := parent.Child(!lastDirection).(*rbTreeNode)
				if sibling != nil {
					if !sibling.Child(!lastDirection).(*rbTreeNode).IsRed() && !sibling.Child(lastDirection).(*rbTreeNode).IsRed() {
						// Color flip
						parent.SetColor(Black)
						sibling.SetColor(Red)
						child.SetColor(Red)
					} else {
						direction2 := grandParent.Right().(*rbTreeNode) == parent
						if sibling.Child(lastDirection).(*rbTreeNode).IsRed() {
							grandParent.SetChild(direction2, DoubleRotate(parent, lastDirection))
						} else if sibling.Child(!lastDirection).(*rbTreeNode).IsRed() {
							grandParent.SetChild(direction2, SingleRotate(parent, lastDirection))
						}

						// When performing a rotation to parent, child is not affected.
						// So all nodes on the path are -1ed.

						// // Update Size
						// parent.Update()
						// grandParent.Child(direction2).Update()

						// Ensure correct coloring
						child.SetColor(Red)
						grandParent.Child(direction2).(*rbTreeNode).SetColor(Red)
						grandParent.Child(direction2).Left().(*rbTreeNode).SetColor(Black)
						grandParent.Child(direction2).Right().(*rbTreeNode).SetColor(Black)
					}
				}
			}
		}
	}

	// Replace and remove the target node
	if target != nil {
		target.SetValue(child.Value())
		parent.SetChild(parent.Right().(*rbTreeNode) == child, child.Child(child.Left().(*rbTreeNode) == nil))
	}

	// Update root and make it black
	root = superRoot.Right().(*rbTreeNode)
	if root != nil {
		root.SetColor(Black)
	}
	return root
}

func (tree *RBTree) Delete(value int) {
	tree.SetRoot(Delete(tree.Root().(*rbTreeNode), value))
}

func SingleRotate(root *rbTreeNode, direction bool) Noded {
	save := root.Child(!direction).(*rbTreeNode)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	root.SetColor(Red)
	save.SetColor(Black)
	return save
}

func DoubleRotate(root *rbTreeNode, direction bool) Noded {
	root.SetChild(!direction, SingleRotate(root.Child(!direction).(*rbTreeNode), !direction))
	return SingleRotate(root, direction)
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
