package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type NodeState bool

const (
	Inactive NodeState = false
	Active   NodeState = true
)

type ScapeGoatNode struct {
	Value  int
	Left   *ScapeGoatNode
	Right  *ScapeGoatNode
	State  NodeState
	Size   uint32 // Number of active nodes in the subtree
	Weight uint32 // Number of nodes in the subtree
}

func NewNode(value int) *ScapeGoatNode {
	return &ScapeGoatNode{
		Value:  value,
		Left:   nil,
		Right:  nil,
		State:  Active,
		Size:   1,
		Weight: 1,
	}
}

func (root *ScapeGoatNode) Leaf() bool {
	return root.Left == nil && root.Right == nil
}

func (root *ScapeGoatNode) Full() bool {
	return root.Left != nil && root.Right != nil
}

func (root *ScapeGoatNode) Inactive() bool {
	return root.State == Inactive
}

func (root *ScapeGoatNode) Active() bool {
	return root.State == Active
}

func (root *ScapeGoatNode) Update() {
	if root.Active() {
		root.Size = 1
	} else {
		root.Size = 0
	}
	root.Weight = 1
	if root.Left != nil {
		root.Size += root.Left.Size
		root.Weight += root.Left.Weight
	}
	if root.Right != nil {
		root.Size += root.Right.Size
		root.Weight += root.Right.Weight
	}
}

func (root *ScapeGoatNode) Deactivate() {
	root.State = Inactive
	root.Update()
}

type ScapeGoat struct {
	Root  *ScapeGoatNode
	Alpha float64
}

func New(alpha float64) *ScapeGoat {
	return &ScapeGoat{
		Root:  nil,
		Alpha: alpha,
	}
}

func ToSlice(root *ScapeGoatNode) []*ScapeGoatNode {
	if root == nil {
		return []*ScapeGoatNode{}
	}
	if root.Active() {
		defer func() {
			root.Left = nil
			root.Right = nil
			root.Size = 1
			root.Weight = 1
		}()
		return append(append(ToSlice(root.Left), root), ToSlice(root.Right)...)
	} else {
		return append(ToSlice(root.Left), ToSlice(root.Right)...)
	}
}

func FromSlice(slice []*ScapeGoatNode) *ScapeGoatNode {
	if len(slice) == 0 {
		return nil
	}
	mid := len(slice) / 2
	root := slice[mid]
	root.Left = FromSlice(slice[:mid])
	root.Right = FromSlice(slice[mid+1:])
	root.Update()
	return root
}

func (thisTree *ScapeGoat) ToSlice() []int {
	slice := ToSlice(thisTree.Root)
	result := make([]int, len(slice))
	for i, root := range slice {
		result[i] = root.Value
	}
	return result
}

func (thisTree *ScapeGoat) FromSlice(slice []int) {
	nodes := make([]*ScapeGoatNode, len(slice))
	for i, value := range slice {
		nodes[i] = NewNode(value)
	}
	thisTree.Root = FromSlice(nodes)
}

func Reconstruct(root *ScapeGoatNode) *ScapeGoatNode {
	return FromSlice(ToSlice(root))
}

func Imbalance(root *ScapeGoatNode, alpha float64) bool {
	if root == nil {
		return false
	}
	if root.Left != nil && root.Left.Weight > uint32(alpha*float64(root.Weight)) {
		return true
	}
	if root.Right != nil && root.Right.Weight > uint32(alpha*float64(root.Weight)) {
		return true
	}
	return false
}

func Insert(root *ScapeGoatNode, value int, alpha float64) *ScapeGoatNode {
	if root == nil {
		return NewNode(value)
	}
	if value < root.Value {
		root.Left = Insert(root.Left, value, alpha)
	} else {
		root.Right = Insert(root.Right, value, alpha)
	}
	root.Update()
	if Imbalance(root, alpha) {
		return Reconstruct(root)
	}
	return root
}

func (thisTree *ScapeGoat) Insert(value int) {
	thisTree.Root = Insert(thisTree.Root, value, thisTree.Alpha)
}

func Rank(root *ScapeGoatNode, value int) uint32 {
	result := uint32(0)
	for root != nil {
		if root.Value >= value {
			root = root.Left
		} else {
			if root.Left != nil {
				result += root.Left.Size
			}
			if root.Active() {
				result += 1
			}
			root = root.Right
		}
	}
	return result + 1
}

func (thisTree *ScapeGoat) Rank(value int) uint32 {
	return Rank(thisTree.Root, value)
}

func Kth(root *ScapeGoatNode, k uint32) *ScapeGoatNode {
	var result *ScapeGoatNode = nil
	for root != nil {
		leftSize := uint32(0)
		if root.Left != nil {
			leftSize = root.Left.Size
		}
		if root.Active() && leftSize+1 == k {
			result = root
			break
		} else if leftSize >= k {
			root = root.Left
		} else {
			k -= leftSize
			if root.Active() {
				k -= 1
			}
			root = root.Right
		}
	}
	return result
}

func (thisTree *ScapeGoat) Kth(k uint32) (int, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return int(rune(0)), errors.New("k is out of range")
	}
	return result.Value, nil
}

func Find(root *ScapeGoatNode, value int) *ScapeGoatNode {
	if root == nil {
		return nil
	}
	if root.Value == value {
		if root.Active() {
			return root
		} else {
			if result := Find(root.Left, value); result != nil {
				return result
			} else if result := Find(root.Right, value); result != nil {
				return result
			}
			return nil
		}
	} else if root.Value > value {
		return Find(root.Left, value)
	} else {
		return Find(root.Right, value)
	}
}

func (thisTree *ScapeGoat) Find(value int) bool {
	return Find(thisTree.Root, value) != nil
}

func Delete(root *ScapeGoatNode, value int) *ScapeGoatNode {
	if root == nil {
		return nil
	}
	if root.Value == value {
		if root.Active() {
			root.Deactivate()
			return root
		} else {
			if result := Delete(root.Left, value); result != nil {
				root.Size -= 1
				return result
			} else if result := Delete(root.Right, value); result != nil {
				root.Size -= 1
				return result
			}
			return nil
		}
	} else if root.Value > value {
		if result := Delete(root.Left, value); result != nil {
			root.Size -= 1
			return result
		}
	} else {
		if result := Delete(root.Right, value); result != nil {
			root.Size -= 1
			return result
		}
	}
	return nil
}

func (thisTree *ScapeGoat) Delete(value int) {
	target := Find(thisTree.Root, value)
	if target != nil {
		Delete(thisTree.Root, value)
	}
}

func (thisTree *ScapeGoat) Clear() {
	thisTree.Root = nil
}

func (thisTree *ScapeGoat) Size() uint32 {
	return thisTree.Root.Size
}

func (thisTree *ScapeGoat) Empty() bool {
	return thisTree.Root == nil
}

func Prev(root *ScapeGoatNode, value int) *ScapeGoatNode {
	return Kth(root, Rank(root, value)-1)
}

func (thisTree *ScapeGoat) Prev(value int) (int, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return int(rune(0)), errors.New("no previous value")
	}
	return prev.Value, nil
}

func Next(root *ScapeGoatNode, value int) *ScapeGoatNode {
	return Kth(root, Rank(root, value+1))
}

func (thisTree *ScapeGoat) Next(value int) (int, error) {
	next := Next(thisTree.Root, value)
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
	tree := New(0.7)
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
