package tree

import (
	"bstrees/internal/node"
	"bstrees/pkg/errors"
	"bytes"
	"fmt"
	"reflect"
)

type Rooted[T node.Ordered] interface {
	Root() node.Noded[T]
	SetRoot(node.Noded[T])
	Size() uint
	Empty() bool
	Clear()
}

type Insertable[T node.Ordered] interface {
	Insert(T)
}

type Deleteable[T node.Ordered] interface {
	Delete(T)
}

type Searchable[T node.Ordered] interface {
	Contains(T) bool
	Kth(uint) (T, error)
	Rank(T) uint
	Prev(T) (T, error)
	Next(T) (T, error)
}

type Treed[T node.Ordered] interface {
	Rooted[T]
	Insertable[T]
	Deleteable[T]
	Searchable[T]
}

type BaseTree[T node.Ordered] struct {
	root node.Noded[T]
}

func IsNil[T node.Ordered](root node.Noded[T]) bool {
	return root == nil || root.IsNil()
}

func New[T node.Ordered]() *BaseTree[T] {
	return &BaseTree[T]{}
}

func (tree *BaseTree[T]) Root() node.Noded[T] {
	return tree.root
}

func (tree *BaseTree[T]) SetRoot(root node.Noded[T]) {
	tree.root = root
}

func (tree *BaseTree[T]) Size() uint {
	if tree.root.IsNil() {
		return 0
	}
	return tree.root.Size()
}

func (tree *BaseTree[T]) Empty() bool {
	return tree.root.IsNil()
}

func (tree *BaseTree[T]) Clear() {
	tree.root = reflect.Zero(reflect.TypeOf(tree.root)).Interface().(node.Noded[T])
}

func Find[T node.Ordered](root node.Noded[T], value T) node.Noded[T] {
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

func (tree *BaseTree[T]) Contains(value T) bool {
	return !Find(tree.root, value).IsNil()
}

func Kth[T node.Ordered](root node.Noded[T], k uint) node.Noded[T] {
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

func (tree *BaseTree[T]) Kth(k uint) (T, error) {
	result := Kth(tree.root, k)
	if IsNil(result) {
		return T(rune(0)), errors.ErrOutOfRange
	}
	return result.Value(), nil
}

func Rank[T node.Ordered](root node.Noded[T], value T) uint {
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

func (tree *BaseTree[T]) Rank(value T) uint {
	return Rank(tree.root, value)
}

func Prev[T node.Ordered](root node.Noded[T], value T) node.Noded[T] {
	var result node.Noded[T] = nil
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

func (tree *BaseTree[T]) Prev(value T) (T, error) {
	prev := Prev(tree.root, value)
	if IsNil(prev) {
		return T(rune(0)), errors.ErrNoPrevValue
	}
	return prev.Value(), nil
}

func Next[T node.Ordered](root node.Noded[T], value T) node.Noded[T] {
	var result node.Noded[T] = nil
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

func (tree *BaseTree[T]) Next(value T) (T, error) {
	next := Next(tree.root, value)
	if IsNil(next) {
		return T(rune(0)), errors.ErrNoNextValue
	}
	return next.Value(), nil
}

func String[T node.Ordered](root node.Noded[T]) string {
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

func (tree *BaseTree[T]) String() string {
	return String(tree.root)
}

func CheckBSTProperty[T node.Ordered](root node.Noded[T]) error {
	if root.IsNil() {
		return nil
	}
	left, right := root.Left(), root.Right()
	if !left.IsNil() && left.Value() > root.Value() {
		return errors.ErrViolatedBSTProperty
	}
	if !right.IsNil() && right.Value() < root.Value() {
		return errors.ErrViolatedBSTProperty
	}
	leftSize := uint(0)
	if !left.IsNil() {
		leftSize = left.Size()
	}
	rightSize := uint(0)
	if !right.IsNil() {
		rightSize = right.Size()
	}
	if leftSize+rightSize+1 != root.Size() {
		return errors.ErrViolatedBSTProperty
	}

	if err := CheckBSTProperty(left); err != nil {
		return err
	}
	if err := CheckBSTProperty(right); err != nil {
		return err
	}
	return nil
}
