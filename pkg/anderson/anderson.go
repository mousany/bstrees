package anderson

import (
	"bstrees/pkg/anderson/node"
	"bstrees/pkg/errors"
	"bstrees/pkg/trait/ordered"
)

type Anderson[T ordered.Ordered] struct {
	Root *node.AndersonNode[T]
}

func New[T ordered.Ordered]() Anderson[T] {
	return Anderson[T]{Root: nil}
}

func LeftRotate[T ordered.Ordered](root *node.AndersonNode[T]) *node.AndersonNode[T] {
	right := root.Right
	root.Right = right.Left
	right.Left = root
	root.Update()
	right.Update()
	return right
}

func RightRotate[T ordered.Ordered](root *node.AndersonNode[T]) *node.AndersonNode[T] {
	left := root.Left
	root.Left = left.Right
	left.Right = root
	root.Update()
	left.Update()
	return left
}

func Skew[T ordered.Ordered](root *node.AndersonNode[T]) *node.AndersonNode[T] {
	// Print(root)
	if root.Left == nil || root.Left.Level != root.Level {
		return root
	}
	return RightRotate(root)
}

func Split[T ordered.Ordered](root *node.AndersonNode[T]) *node.AndersonNode[T] {
	if root.Right == nil || root.Right.Right == nil || root.Right.Right.Level != root.Level {
		return root
	}
	root = LeftRotate(root)
	root.Level += 1
	return root
}

func Insert[T ordered.Ordered](root *node.AndersonNode[T], value T) *node.AndersonNode[T] {
	if root == nil {
		return node.New(value, 1)
	}
	if value < root.Value {
		root.Left = Insert(root.Left, value)
	} else {
		root.Right = Insert(root.Right, value)
	}
	root.Update()
	root = Skew(root)
	root = Split(root)
	return root
}

func Delete[T ordered.Ordered](root *node.AndersonNode[T], value T) *node.AndersonNode[T] {
	if root == nil {
		return nil
	}
	if value < root.Value {
		root.Left = Delete(root.Left, value)
	} else if value > root.Value {
		root.Right = Delete(root.Right, value)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		} else {
			minNode := Kth(root.Right, 1)
			root.Value = minNode.Value
			root.Right = Delete(root.Right, minNode.Value)
		}
	}
	root.Update()
	if (root.Left != nil && root.Left.Level < root.Level-1) ||
		(root.Right != nil && root.Right.Level < root.Level-1) {
		root.Level -= 1
		if root.Right != nil && root.Right.Level > root.Level {
			root.Right.Level = root.Level
		}
		root = Skew(root)
		root = Split(root)
	}
	return root
}

func Kth[T ordered.Ordered](root *node.AndersonNode[T], k uint32) *node.AndersonNode[T] {
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

func (tree *Anderson[T]) Insert(value T) {
	tree.Root = Insert(tree.Root, value)
}

func (tree *Anderson[T]) Delete(value T) {
	tree.Root = Delete(tree.Root, value)
}

func (tree *Anderson[T]) Kth(k uint32) (T, error) {
	root := Kth(tree.Root, k)
	if root == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return root.Value, nil
}

func (tree *Anderson[T]) Size() uint32 {
	if tree.Root == nil {
		return 0
	}
	return tree.Root.Size
}

func (tree *Anderson[T]) Empty() bool {
	return tree.Root == nil
}

func (tree *Anderson[T]) Clear() {
	tree.Root = nil
}

func Find[T ordered.Ordered](root *node.AndersonNode[T], value T) *node.AndersonNode[T] {
	for root != nil {
		if value < root.Value {
			root = root.Left
		} else if root.Value < value {
			root = root.Right
		} else {
			return root
		}
	}
	return nil
}

func (tree *Anderson[T]) Contains(value T) bool {
	return Find(tree.Root, value) != nil
}

func Rank[T ordered.Ordered](root *node.AndersonNode[T], value T) uint32 {
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

func (thisTree *Anderson[T]) Rank(value T) uint32 {
	return Rank(thisTree.Root, value)
}

func Prev[T ordered.Ordered](root *node.AndersonNode[T], value T) *node.AndersonNode[T] {
	var prev *node.AndersonNode[T] = nil
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

func (thisTree *Anderson[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T ordered.Ordered](root *node.AndersonNode[T], value T) *node.AndersonNode[T] {
	var next *node.AndersonNode[T] = nil
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

func (thisTree *Anderson[T]) Next(value T) (T, error) {
	prev := Next(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return prev.Value, nil
}

// func Print[T ordered.Ordered](root *node.AndersonNode[T]) string {
// 	if root == nil {
// 		return "null"
// 	}
// 	var buffer bytes.Buffer
// 	buffer.WriteString(fmt.Sprintf("[%v", root.Value))
// 	// if root.Red() {
// 	// 	buffer.WriteString("1, ")
// 	// } else {
// 	// 	buffer.WriteString("0, ")
// 	// }
// 	buffer.WriteString(fmt.Sprintf(".%v, ", root.Size))
// 	// buffer.WriteString(fmt.Sprintf("%v, ", root.Size))
// 	buffer.WriteString(fmt.Sprintf("%v, ", Print(root.Left)))
// 	buffer.WriteString(fmt.Sprintf("%v]", Print(root.Right)))
// 	return buffer.String()
// }

// func (thisTree *Anderson[T]) Print() {
// 	fmt.Println(Print(thisTree.Root))
// }
