package scapegoat

import (
	"bstrees/pkg/errors"
	"bstrees/pkg/scapegoat/node"
	"bstrees/pkg/trait/number"
)

type ScapeGoat[T number.Number] struct {
	Root  *node.ScapeGoatNode[T]
	Alpha float64
}

func New[T number.Number](alpha float64) *ScapeGoat[T] {
	return &ScapeGoat[T]{
		Root:  nil,
		Alpha: alpha,
	}
}

func ToSlice[T number.Number](root *node.ScapeGoatNode[T]) []*node.ScapeGoatNode[T] {
	if root == nil {
		return []*node.ScapeGoatNode[T]{}
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

func FromSlice[T number.Number](slice []*node.ScapeGoatNode[T]) *node.ScapeGoatNode[T] {
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

func (thisTree *ScapeGoat[T]) ToSlice() []T {
	slice := ToSlice(thisTree.Root)
	result := make([]T, len(slice))
	for i, node := range slice {
		result[i] = node.Value
	}
	return result
}

func (thisTree *ScapeGoat[T]) FromSlice(slice []T) {
	nodes := make([]*node.ScapeGoatNode[T], len(slice))
	for i, value := range slice {
		nodes[i] = node.New(value)
	}
	thisTree.Root = FromSlice(nodes)
}

func Reconstruct[T number.Number](root *node.ScapeGoatNode[T]) *node.ScapeGoatNode[T] {
	return FromSlice(ToSlice(root))
}

func Imbalance[T number.Number](root *node.ScapeGoatNode[T], alpha float64) bool {
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

func Insert[T number.Number](root *node.ScapeGoatNode[T], value T, alpha float64) *node.ScapeGoatNode[T] {
	if root == nil {
		return node.New(value)
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

func (thisTree *ScapeGoat[T]) Insert(value T) {
	thisTree.Root = Insert(thisTree.Root, value, thisTree.Alpha)
}

func Rank[T number.Number](root *node.ScapeGoatNode[T], value T) uint32 {
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

func (thisTree *ScapeGoat[T]) Rank(value T) uint32 {
	return Rank(thisTree.Root, value)
}

func Kth[T number.Number](root *node.ScapeGoatNode[T], k uint32) *node.ScapeGoatNode[T] {
	var result *node.ScapeGoatNode[T] = nil
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

func (thisTree *ScapeGoat[T]) Kth(k uint32) (T, error) {
	result := Kth(thisTree.Root, k)
	if result == nil {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.Value, nil
}

func Find[T number.Number](root *node.ScapeGoatNode[T], value T) *node.ScapeGoatNode[T] {
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

func (thisTree *ScapeGoat[T]) Find(value T) bool {
	return Find(thisTree.Root, value) != nil
}

func Delete[T number.Number](root *node.ScapeGoatNode[T], value T) *node.ScapeGoatNode[T] {
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

func (thisTree *ScapeGoat[T]) Delete(value T) {
	target := Find(thisTree.Root, value)
	if target != nil {
		Delete(thisTree.Root, value)
	}
}

func (thisTree *ScapeGoat[T]) Clear() {
	thisTree.Root = nil
}

func (thisTree *ScapeGoat[T]) Size() uint32 {
	return thisTree.Root.Size
}

func (thisTree *ScapeGoat[T]) Empty() bool {
	return thisTree.Root == nil
}

func Prev[T number.Number](root *node.ScapeGoatNode[T], value T) *node.ScapeGoatNode[T] {
	return Kth(root, Rank(root, value)-1)
}

func (thisTree *ScapeGoat[T]) Prev(value T) (T, error) {
	prev := Prev(thisTree.Root, value)
	if prev == nil {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value, nil
}

func Next[T number.Number](root *node.ScapeGoatNode[T], value T) *node.ScapeGoatNode[T] {
	return Kth(root, Rank(root, value+1))
}

func (thisTree *ScapeGoat[T]) Next(value T) (T, error) {
	next := Next(thisTree.Root, value)
	if next == nil {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value, nil
}

// func Print[T number.Number](root *node.ScapeGoatNode[T]) string {
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

// func (thisTree *ScapeGoat[T]) Print() {
// 	fmt.Println(Print(thisTree.Root))
// }
