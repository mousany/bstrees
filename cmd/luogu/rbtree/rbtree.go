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

func (thisNode *RBNode) Leaf() bool {
	return thisNode.Left == nil && thisNode.Right == nil
}

func (thisNode *RBNode) Child(direction bool) *RBNode {
	if direction {
		return thisNode.Right
	} else {
		return thisNode.Left
	}
}

func (thisNode *RBNode) SetChild(direction bool, child *RBNode) {
	if direction {
		thisNode.Right = child
	} else {
		thisNode.Left = child
	}
	// thisNode.Update()
}

func IsRed(root *RBNode) bool {
	return root != nil && root.Red()
}

func IsBlack(root *RBNode) bool {
	return root == nil || root.Black()
}

type RBTree struct {
	Root *RBNode
}

func New() *RBTree {
	return &RBTree{Root: nil}
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
	return int(rune(0)), errors.New("index out of range")
}

func SingleRotate(root *RBNode, direction bool) *RBNode {
	save := root.Child(!direction)
	root.SetChild(!direction, save.Child(direction))
	save.SetChild(direction, root)
	root.Update()
	save.Update()
	root.Color = Red
	save.Color = Black
	return save
}

func DoubleRotate(root *RBNode, direction bool) *RBNode {
	root.SetChild(!direction, SingleRotate(root.Child(!direction), !direction))
	return SingleRotate(root, direction)
}

// https://archive.ph/EJTsz, Eternally Confuzzled's Blog
func (thisTree *RBTree) Insert(value int) {
	if thisTree.Root == nil {
		thisTree.Root = NewNode(value)
	} else {
		superRoot := NewNode(int(rune(0))) // Head in Eternally Confuzzled's paper
		superRoot.Right = thisTree.Root

		var child *RBNode = thisTree.Root        // Q in Eternally Confuzzled's paper
		var parent *RBNode = nil                 // P in Eternally Confuzzled's paper
		var grandParent *RBNode = nil            // G in Eternally Confuzzled's paper
		var greatGrandParent *RBNode = superRoot // int in Eternally Confuzzled's paper

		var direction bool = false
		var lastDirection bool = false

		// Search down
		for ok := false; !ok; {
			if child == nil {
				// Insert new node at the bottom
				child = NewNode(value)
				parent.SetChild(direction, child)
				ok = true
			} else {
				// Update size
				child.Size += 1
				if IsRed(child.Left) && IsRed(child.Right) {
					// Color flip
					child.Color = Red
					child.Left.Color = Black
					child.Right.Color = Black
				}
			}

			if IsRed(child) && IsRed(parent) {
				// Fix red violation
				direction2 := greatGrandParent.Right == grandParent
				if child == parent.Child(lastDirection) {
					greatGrandParent.SetChild(direction2, SingleRotate(grandParent, !lastDirection))
				} else {
					greatGrandParent.SetChild(direction2, DoubleRotate(grandParent, !lastDirection))
					if !ok {
						greatGrandParent.Child(direction2).Size += 1
					}
				}
			}

			lastDirection = direction
			direction = child.Value < value
			if grandParent != nil {
				greatGrandParent = grandParent
			}

			grandParent = parent
			parent = child
			child = child.Child(direction)
		}

		// Update root
		thisTree.Root = superRoot.Right
	}

	thisTree.Root.Color = Black
}

func (thisTree *RBTree) Contains(value int) bool {
	for root := thisTree.Root; root != nil; {
		if root.Value == value {
			return true
		} else if root.Value < value {
			root = root.Right
		} else {
			root = root.Left
		}
	}
	return false
}

func (thisTree *RBTree) Delete(value int) {
	if thisTree.Root != nil && thisTree.Contains(value) {
		superRoot := NewNode(int(rune(0))) // Head in Eternally Confuzzled's paper
		superRoot.Right = thisTree.Root

		var child *RBNode = superRoot // Q in Eternally Confuzzled's paper
		var parent *RBNode = nil      // P in Eternally Confuzzled's paper
		var grandParent *RBNode = nil // G in Eternally Confuzzled's paper
		var target *RBNode = nil      // F in Eternally Confuzzled's paper
		direction := true

		// Search and push a red down
		for child.Child(direction) != nil {
			lastDirection := direction

			grandParent = parent
			parent = child
			child = child.Child(direction)
			direction = child.Value < value

			// Update size
			child.Size -= 1

			// Save the target node
			if child.Value == value {
				target = child
			}

			// Push the red node down
			if !IsRed(child) && !IsRed(child.Child(direction)) {
				if IsRed(child.Child(!direction)) {
					parent.SetChild(lastDirection, SingleRotate(child, direction))
					parent = parent.Child(lastDirection)
					child.Size -= 1
					parent.Update()
				} else if !IsRed(child.Child(!direction)) {
					sibling := parent.Child(!lastDirection)
					if sibling != nil {
						if !IsRed(sibling.Child(!lastDirection)) && !IsRed(sibling.Child(lastDirection)) {
							// Color flip
							parent.Color = Black
							sibling.Color = Red
							child.Color = Red
						} else {
							direction2 := grandParent.Right == parent
							if IsRed(sibling.Child(lastDirection)) {
								grandParent.SetChild(direction2, DoubleRotate(parent, lastDirection))
							} else if IsRed(sibling.Child(!lastDirection)) {
								grandParent.SetChild(direction2, SingleRotate(parent, lastDirection))
							}

							// // Update Size
							// parent.Update()
							// grandParent.Child(direction2).Update()

							// Ensure correct coloring
							child.Color = Red
							grandParent.Child(direction2).Color = Red
							grandParent.Child(direction2).Left.Color = Black
							grandParent.Child(direction2).Right.Color = Black
						}
					}
				}
			}
		}

		// Replace and remove the target node
		if target != nil {
			target.Value = child.Value
			parent.SetChild(parent.Right == child, child.Child(child.Left == nil))
		}

		// Update root and make it black
		thisTree.Root = superRoot.Right
		if thisTree.Root != nil {
			thisTree.Root.Color = Black
		}
	}
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
