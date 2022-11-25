package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type SplayNode struct {
	Value  int
	Left   *SplayNode
	Right  *SplayNode
	Parent *SplayNode
	Size   uint32
	Rec    uint32 // This field is Splay only
	// Because Splay operation will scatter nodes with the same value
	// While traditional BST search mechanics is too slow on Splay
}

func NewNode(value int) *SplayNode {
	return &SplayNode{
		Value:  value,
		Left:   nil,
		Right:  nil,
		Parent: nil,
		Size:   1,
		Rec:    1,
	}
}

func (root *SplayNode) Update() {
	root.Size = root.Rec
	if root.Left != nil {
		root.Size += root.Left.Size
	}
	if root.Right != nil {
		root.Size += root.Right.Size
	}
}

func (root *SplayNode) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *SplayNode) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *SplayNode) SetChild(child *SplayNode, direction bool) {
	if direction {
		root.Right = child
	} else {
		root.Left = child
	}
	if child != nil {
		child.Parent = root
	}
}

func (root *SplayNode) Child(direction bool) *SplayNode {
	if direction {
		return root.Right
	}
	return root.Left
}

type Splay struct {
	superRoot *SplayNode
}

func (thisTree *Splay) Root() *SplayNode {
	return thisTree.superRoot.Right
}

func (thisTree *Splay) SetRoot(root *SplayNode) {
	thisTree.superRoot.SetChild(root, true)
}

func New() *Splay {
	return &Splay{
		superRoot: NewNode(int(rune(0))),
	}
}

func LeftRotate(root *SplayNode) *SplayNode {
	right := root.Right
	root.SetChild(right.Left, true)
	right.SetChild(root, false)
	root.Update()
	right.Update()
	return right
}

func RightRotate(root *SplayNode) *SplayNode {
	left := root.Left
	root.SetChild(left.Right, false)
	left.SetChild(root, true)
	root.Update()
	left.Update()
	return left
}

// Rotate root to its parent
// After this operation, parent will be the child of root
func RotateToParent(root *SplayNode) {
	grandParent := root.Parent.Parent
	if root == root.Parent.Left {
		// root is left child
		root = RightRotate(root.Parent)
	} else {
		// root is right child
		root = LeftRotate(root.Parent)
	}
	if grandParent != nil {
		if grandParent.Left == root.Parent {
			grandParent.SetChild(root, false)
			grandParent.Update()
		} else {
			grandParent.SetChild(root, true)
			grandParent.Update()
		}
	}
}

// Rotate root to target
// After this operation, target will be the child of root
func SplayRotate(root, target *SplayNode) {
	targetParent := target.Parent
	for root.Parent != targetParent {
		parent := root.Parent
		grandParent := parent.Parent
		direction := root == parent.Left
		grandDirection := parent == grandParent.Left
		if parent == target {
			// root is the child of target
			RotateToParent(root)
		} else if direction == grandDirection {
			// zig-zig
			RotateToParent(parent)
			RotateToParent(root)
		} else {
			// zig-zag
			RotateToParent(root)
			RotateToParent(root)
		}
	}
}

func Find(root *SplayNode, value int) *SplayNode {
	for p := root; p != nil; {
		if p.Value == value {
			return p
		} else if value < p.Value {
			p = p.Left
		} else {
			p = p.Right
		}
	}
	return nil
}

func Kth(root *SplayNode, k uint32) *SplayNode {
	for p := root; p != nil; {
		leftSize := uint32(0)
		if p.Left != nil {
			leftSize = p.Left.Size
		}
		if leftSize < k && leftSize+p.Rec >= k {
			// SplayRotate(p, root)
			return p
		} else if leftSize+p.Rec < k {
			k -= leftSize + p.Rec
			p = p.Right
		} else {
			p = p.Left
		}
	}
	return nil
}

func Insert(root *SplayNode, value int) *SplayNode {
	if root == nil {
		return NewNode(value)
	} else {
		superRoot := root.Parent

		for p := root; p != nil; {
			p.Size += 1
			if value == p.Value {
				p.Rec += 1
				SplayRotate(p, root)
				break
			} else if value < p.Value {
				if p.Left == nil {
					p.SetChild(NewNode(value), false)
					SplayRotate(p.Left, root)
					break
				} else {
					p = p.Left
				}
			} else {
				if p.Right == nil {
					p.SetChild(NewNode(value), true)
					SplayRotate(p.Right, root)
					break
				} else {
					p = p.Right
				}
			}
		}

		return superRoot.Right
	}
}

func Delete(root *SplayNode, value int) *SplayNode {
	if root == nil {
		return nil
	}
	superRoot := root.Parent
	p := Find(root, value)
	if p == nil {
		return root
	}
	SplayRotate(p, root)
	if p.Rec > 1 {
		p.Rec -= 1
		p.Size -= 1
	} else {
		if p.Left == nil && p.Right == nil {
			superRoot.SetChild(nil, true)
		} else if p.Left == nil {
			superRoot.SetChild(p.Right, true)
		} else if p.Right == nil {
			superRoot.SetChild(p.Left, true)
		} else {
			maxLeft := p.Left
			for maxLeft.Right != nil {
				maxLeft.Size -= 1
				maxLeft = maxLeft.Right
			}
			SplayRotate(maxLeft, superRoot.Right)
			maxLeft.SetChild(p.Right, true)
			superRoot.SetChild(maxLeft, true)
			superRoot.Right.Update()
		}
	}

	return superRoot.Right
}

func (thisTree *Splay) Insert(value int) {
	thisTree.SetRoot(Insert(thisTree.Root(), value))
}

func (thisTree *Splay) Delete(value int) {
	thisTree.SetRoot(Delete(thisTree.Root(), value))
}

func (thisTree *Splay) Contains(value int) bool {
	return Find(thisTree.Root(), value) != nil
}

func (thisTree *Splay) Kth(k uint32) (int, error) {
	result := Kth(thisTree.Root(), k)
	if result == nil {
		return int(rune(0)), errors.New("k is out of range")
	}
	return result.Value, nil
}

func (thisTree *Splay) Size() uint32 {
	if thisTree.Root() == nil {
		return 0
	}
	return thisTree.Root().Size
}

func (thisTree *Splay) Empty() bool {
	return thisTree.Root() == nil
}

func (thisTree *Splay) Clear() {
	thisTree.SetRoot(nil)
}

func Rank(root *SplayNode, value int) uint32 {
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

func (thisTree *Splay) Rank(value int) uint32 {
	// return Rank(thisTree.Root, value)
	p := Find(thisTree.Root(), value)
	if p == nil {
		prev := Prev(thisTree.Root(), value)
		if prev != nil {
			SplayRotate(prev, thisTree.Root())
			if prev.Left != nil {
				return prev.Left.Size + prev.Rec + 1
			}
			return prev.Rec + 1
		}
		return 1
	}
	SplayRotate(p, thisTree.Root())
	if p.Left != nil {
		return p.Left.Size + 1
	}
	return 1
}

func Prev(root *SplayNode, value int) *SplayNode {
	var result *SplayNode
	for p := root; p != nil; {
		if value > p.Value {
			result = p
			p = p.Right
		} else {
			p = p.Left
		}
	}
	return result
}

func (thisTree *Splay) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root(), value)
	if prev == nil {
		return int(rune(0)), errors.New("no prev value")
	}
	return prev.Value, nil
}

func Next(root *SplayNode, value int) *SplayNode {
	var result *SplayNode
	for p := root; p != nil; {
		if value < p.Value {
			result = p
			p = p.Left
		} else {
			p = p.Right
		}
	}
	return result
}

func (thisTree *Splay) Next(value int) (int, error) {
	next := Next(thisTree.Root(), value)
	if next == nil {
		return int(rune(0)), errors.New("no next value")
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
