package treap

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

type TreapTree[T node.Ordered] struct {
	*tree.BaseTree[T]
}

func New[T node.Ordered]() *TreapTree[T] {
	tr := &TreapTree[T]{tree.New[T]()}
	tr.SetRoot((*TreapTreeNode[T])(nil))
	return tr
}

func insert[T node.Ordered](root *TreapTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return NewTreapTreeNode(value)
	}
	if root.Value() <= value {
		root.SetRight(insert(root.Right().(*TreapTreeNode[T]), value))
		if root.Right().(*TreapTreeNode[T]).Weight() < root.Weight() {
			root = tree.SingleRotate(false, node.Noded[T](root)).(*TreapTreeNode[T])
		}
	} else {
		root.SetLeft(insert(root.Left().(*TreapTreeNode[T]), value))
		if root.Left().(*TreapTreeNode[T]).Weight() < root.Weight() {
			root = tree.SingleRotate(true, node.Noded[T](root)).(*TreapTreeNode[T])
		}
	}
	root.Update()
	return root
}

func (tr *TreapTree[T]) Insert(value T) {
	tr.SetRoot(insert(tr.Root().(*TreapTreeNode[T]), value))
}

func delete[T node.Ordered](root *TreapTreeNode[T], value T) node.Noded[T] {
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
		if root.Left().(*TreapTreeNode[T]).Weight() < root.Right().(*TreapTreeNode[T]).Weight() {
			root = tree.SingleRotate(true, node.Noded[T](root)).(*TreapTreeNode[T])
			root.SetRight(delete(root.Right().(*TreapTreeNode[T]), value))
		} else {
			root = tree.SingleRotate(false, node.Noded[T](root)).(*TreapTreeNode[T])
			root.SetLeft(delete(root.Left().(*TreapTreeNode[T]), value))
		}
	} else if root.Value() < value {
		root.SetRight(delete(root.Right().(*TreapTreeNode[T]), value))
	} else {
		root.SetLeft(delete(root.Left().(*TreapTreeNode[T]), value))
	}
	root.Update()
	return root
}

func (tr *TreapTree[T]) Delete(value T) {
	tr.SetRoot(delete(tr.Root().(*TreapTreeNode[T]), value))
}

func fromSlice[T node.Ordered](slice []T) *TreapTreeNode[T] {
	if len(slice) == 0 {
		return nil
	}
	mid := len(slice) / 2
	root := NewTreapTreeNode(slice[mid])
	root.SetLeft(fromSlice(slice[:mid]))
	root.SetRight(fromSlice(slice[mid+1:]))
	root.Update()
	return root
}

func (tr *TreapTree[T]) FromSlice(slice []T) {
	tr.SetRoot(fromSlice(slice))
}
