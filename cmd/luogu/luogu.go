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

func IsNil(root Noded) bool {
	return root == nil || root.IsNil()
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
	if tree.root.IsNil() {
		return 0
	}
	return tree.root.Size()
}

func (tree *BaseTree) Empty() bool {
	return tree.root.IsNil()
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
	return !IsNil(Find(tree.root, value))
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
	if IsNil(result) {
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
	if IsNil(prev) {
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
	if IsNil(next) {
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
	return SingleRotate(false, root)
}

func RightRotate(root Noded) Noded {
	return SingleRotate(true, root)
}

type splayTreeNode struct {
	*BaseTreeNode

	parent *splayTreeNode
	rec    uint // This field is Splay only
	// Because Splay operation will scatter nodes with the same value
	// While traditional BST search mechanics is too slow on Splay
}

func newSplayTreeNode(value int) *splayTreeNode {
	n := &splayTreeNode{
		BaseTreeNode: newBaseTreeNode(value),
		parent:       (*splayTreeNode)(nil),
		rec:          1,
	}
	n.SetLeft((*splayTreeNode)(nil))
	n.SetRight((*splayTreeNode)(nil))
	return n
}

func (n *splayTreeNode) IsNil() bool {
	return n == nil
}

func (n *splayTreeNode) Parent() *splayTreeNode {
	return n.parent
}

func (n *splayTreeNode) SetParent(parent *splayTreeNode) {
	n.parent = parent
}

func (n *splayTreeNode) Rec() uint {
	return n.rec
}

func (n *splayTreeNode) SetRec(rec uint) {
	n.rec = rec
}

func (n *splayTreeNode) SetLeft(left Noded) {
	n.BaseTreeNode.SetLeft(left)
	if !left.IsNil() {
		left.(*splayTreeNode).SetParent(n)
	}
}

func (n *splayTreeNode) SetRight(right Noded) {
	n.BaseTreeNode.SetRight(right)
	if !right.IsNil() {
		right.(*splayTreeNode).SetParent(n)
	}
}

func (n *splayTreeNode) SetChild(right bool, child Noded) {
	if right {
		n.SetRight(child)
	} else {
		n.SetLeft(child)
	}
}

func (n *splayTreeNode) Update() {
	n.BaseTreeNode.Update()
	n.SetSize(n.Size() + n.Rec() - 1)
}

type SplayTree struct {
	*BaseTree // In splay tree, root is not the root of the tree, but the right child of superRoot
}

func New() *SplayTree {
	tr := &SplayTree{newBaseTree()}
	tr.BaseTree.SetRoot(newSplayTreeNode(int(rune(0))))
	return tr
}

func (tr *SplayTree) Root() Noded {
	return tr.BaseTree.Root().Right()
}

func (tr *SplayTree) SetRoot(root Noded) {
	tr.BaseTree.Root().SetRight(root)
}

func insert(root *splayTreeNode, value int) Noded {
	if root == nil {
		return newSplayTreeNode(value)
	} else {
		superRoot := root.Parent()
		for p := root; p != nil; {
			p.SetSize(p.Size() + 1)
			if value == p.Value() {
				p.SetRec(p.Rec() + 1)
				splayRotate(p, root)
				break
			} else if value < p.Value() {
				if p.Left().IsNil() {
					p.SetLeft(newSplayTreeNode(value))
					splayRotate(p.Left().(*splayTreeNode), root)
					break
				} else {
					p = p.Left().(*splayTreeNode)
				}
			} else {
				if p.Right().IsNil() {
					p.SetRight(newSplayTreeNode(value))
					splayRotate(p.Right().(*splayTreeNode), root)
					break
				} else {
					p = p.Right().(*splayTreeNode)
				}
			}
		}
		return superRoot.Right()
	}
}

func (tr *SplayTree) Insert(value int) {
	tr.SetRoot(insert(tr.Root().(*splayTreeNode), value))
}

func Delete(root *splayTreeNode, value int) Noded {
	if root == nil {
		return (*splayTreeNode)(nil)
	}
	superRoot := root.Parent()
	p := Find(Noded(root), value).(*splayTreeNode)
	if p == nil {
		return root
	}
	splayRotate(p, root)
	if p.Rec() > 1 {
		p.SetRec(p.Rec() - 1)
		p.SetSize(p.Size() - 1)
	} else {
		if p.Left().IsNil() && p.Right().IsNil() {
			superRoot.SetRight((*splayTreeNode)(nil))
		} else if p.Left().IsNil() {
			superRoot.SetRight(p.Right())
		} else if p.Right().IsNil() {
			superRoot.SetRight(p.Left())
		} else {
			maxLeft := p.Left()
			for !maxLeft.Right().IsNil() {
				maxLeft.SetSize(maxLeft.Size() - 1)
				maxLeft = maxLeft.Right()
			}
			splayRotate(maxLeft.(*splayTreeNode), superRoot.Right().(*splayTreeNode))
			maxLeft.SetRight(p.Right())
			superRoot.SetRight(maxLeft)
			superRoot.Right().Update()
		}
	}

	return superRoot.Right()
}

func (tr *SplayTree) Delete(value int) {
	tr.SetRoot(Delete(tr.Root().(*splayTreeNode), value))
}

func kth(root *splayTreeNode, k uint) Noded {
	for p := root; p != nil; {
		leftSize := uint(0)
		if !p.Left().IsNil() {
			leftSize = p.Left().Size()
		}
		if leftSize < k && leftSize+p.Rec() >= k {
			splayRotate(p, root)
			return p
		} else if leftSize+p.Rec() < k {
			k -= leftSize + p.Rec()
			p = p.Right().(*splayTreeNode)
		} else {
			p = p.Left().(*splayTreeNode)
		}
	}
	return (*splayTreeNode)(nil)
}

func (tr *SplayTree) Kth(k uint) (int, error) {
	result := kth(tr.Root().(*splayTreeNode), k)
	if result.IsNil() {
		return int(rune(0)), errors.New("k is out of range")
	}
	return result.Value(), nil
}

func (tr *SplayTree) Rank(value int) uint {
	p := Find(tr.Root(), value)
	if IsNil(p) {
		prev := Prev(tr.Root(), value).(*splayTreeNode)
		if prev != nil {
			splayRotate(prev, tr.Root().(*splayTreeNode))
			if !prev.Left().IsNil() {
				return prev.Left().Size() + prev.Rec() + 1
			}
			return prev.Rec() + 1
		}
		return 1
	}
	splayRotate(p.(*splayTreeNode), tr.Root().(*splayTreeNode))
	if !p.Left().IsNil() {
		return p.Left().Size() + 1
	}
	return 1
}

func (tr *SplayTree) Size() uint {
	if tr.Root().IsNil() {
		return 0
	}
	return tr.Root().Size()
}

func (tr *SplayTree) Empty() bool {
	return tr.Root().IsNil()
}

func (tr *SplayTree) Clear() {
	tr.SetRoot((*splayTreeNode)(nil))
}

func (tr *SplayTree) Contains(value int) bool {
	return IsNil(Find(tr.Root(), value))
}

func (tr *SplayTree) Prev(value int) (int, error) {
	result := Prev(tr.Root(), value)
	if IsNil(result) {
		return int(rune(0)), errors.New("no previous value")
	}
	return result.Value(), nil
}

func (tr *SplayTree) Next(value int) (int, error) {
	result := Next(tr.Root(), value)
	if IsNil(result) {
		return int(rune(0)), errors.New("no next value")
	}
	return result.Value(), nil
}

func (tr *SplayTree) String() string {
	return String(tr.Root())
}

// Rotate root to its parent
// After this operation, parent will be the child of root
func rotateToParent(root *splayTreeNode) {
	grandParent := root.Parent().Parent()
	parentDirection := root == root.Parent().Left().(*splayTreeNode)
	root = SingleRotate(parentDirection, Noded(root.Parent())).(*splayTreeNode)
	if grandParent != nil {
		grandParentDirection := root.Parent() == grandParent.Left().(*splayTreeNode)
		grandParent.SetChild(!grandParentDirection, Noded(root))
	}
}

// Rotate root to target
// After this operation, target will be the child of root
func splayRotate(root, target *splayTreeNode) {
	targetParent := target.Parent()
	for root.Parent() != targetParent {
		parent := root.Parent()
		grandParent := parent.Parent()
		direction := root == parent.Left().(*splayTreeNode)
		grandDirection := parent == grandParent.Left().(*splayTreeNode)
		if parent == target {
			// root is the child of target
			rotateToParent(root)
		} else if direction == grandDirection {
			// zig-zig
			rotateToParent(parent)
			rotateToParent(root)
		} else {
			// zig-zag
			rotateToParent(root)
			rotateToParent(root)
		}
	}
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
