package fhq

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
	"bstrees/pkg/errors"
	"bstrees/pkg/treap"
)

type FHQTree[T node.Number] struct {
	*tree.BaseTree[T]
}

func New[T node.Number]() *FHQTree[T] {
	tr := &FHQTree[T]{tree.New[T]()}
	tr.SetRoot((*treap.TreapTreeNode[T])(nil))
	return tr
}

func (tr *FHQTree[T]) Insert(value T) {
	left, right := split(tr.Root().(*treap.TreapTreeNode[T]), value)
	tr.SetRoot(merge(merge(left, treap.NewTreapTreeNode(value)), right))
}

func (tr *FHQTree[T]) Delete(value T) {
	left, right := split(tr.Root().(*treap.TreapTreeNode[T]), value)
	left, mid := split(left, value-1)
	if mid != nil {
		mid = merge(mid.Left().(*treap.TreapTreeNode[T]), mid.Right().(*treap.TreapTreeNode[T]))
	}
	tr.SetRoot(merge(merge(left, mid), right))
}

func (tr *FHQTree[T]) Rank(value T) uint {
	left, right := split(tr.Root().(*treap.TreapTreeNode[T]), value-1)
	defer func() {
		tr.SetRoot(merge(left, right))
	}()
	if left == nil {
		return 1
	}
	return left.Size() + 1
}

func (tr *FHQTree[T]) Prev(value T) (T, error) {
	left, right := split(tr.Root().(*treap.TreapTreeNode[T]), value-1)
	defer func() {
		tr.SetRoot(merge(left, right))
	}()
	result := tree.Kth(node.Noded[T](left), left.Size())
	if result == nil {
		return T(0), errors.ErrNoPrevValue
	}
	return result.Value(), nil
}

func (tr *FHQTree[T]) Next(value T) (T, error) {
	left, right := split(tr.Root().(*treap.TreapTreeNode[T]), value)
	defer func() {
		tr.SetRoot(merge(left, right))
	}()
	result := tree.Kth(node.Noded[T](right), 1)
	if result == nil {
		return T(0), errors.ErrNoNextValue
	}
	return result.Value(), nil
}

func fromSlice[T node.Ordered](slice []T) *treap.TreapTreeNode[T] {
	if len(slice) == 0 {
		return nil
	}
	mid := len(slice) / 2
	root := treap.NewTreapTreeNode(slice[mid])
	root.SetLeft(fromSlice(slice[:mid]))
	root.SetRight(fromSlice(slice[mid+1:]))
	root.Update()
	return root
}

func (tr *FHQTree[T]) FromSlice(slice []T) {
	tr.SetRoot(fromSlice(slice))
}
