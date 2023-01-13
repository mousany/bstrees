package treap

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

type TreapTree[T node.Ordered] struct {
	*tree.BaseTree[T]
}

func New[T node.Ordered]() *TreapTree[T] {
	tree := &TreapTree[T]{tree.New[T]()}
	tree.SetRoot((*treapTreeNode[T])(nil))
	return tree
}

func insert[T node.Ordered](root *treapTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return newTreapTreeNode(value)
	}
	if root.Value() <= value {
		root.SetRight(insert(root.Right().(*treapTreeNode[T]), value))
		if root.Right().(*treapTreeNode[T]).Weight() < root.Weight() {
			root = tree.SingleRotate(false, node.Noded[T](root)).(*treapTreeNode[T])
		}
	} else {
		root.SetLeft(insert(root.Left().(*treapTreeNode[T]), value))
		if root.Left().(*treapTreeNode[T]).Weight() < root.Weight() {
			root = tree.SingleRotate(true, node.Noded[T](root)).(*treapTreeNode[T])
		}
	}
	root.Update()
	return root
}

func (tree *TreapTree[T]) Insert(value T) {
	tree.SetRoot(insert(tree.Root().(*treapTreeNode[T]), value))
}

func delete[T node.Ordered](root *treapTreeNode[T], value T) node.Noded[T] {
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
		if root.Left().(*treapTreeNode[T]).Weight() < root.Right().(*treapTreeNode[T]).Weight() {
			root = tree.SingleRotate(true, node.Noded[T](root)).(*treapTreeNode[T])
			root.SetRight(delete(root.Right().(*treapTreeNode[T]), value))
		} else {
			root = tree.SingleRotate(false, node.Noded[T](root)).(*treapTreeNode[T])
			root.SetLeft(delete(root.Left().(*treapTreeNode[T]), value))
		}
	} else if root.Value() < value {
		root.SetRight(delete(root.Right().(*treapTreeNode[T]), value))
	} else {
		root.SetLeft(delete(root.Left().(*treapTreeNode[T]), value))
	}
	root.Update()
	return root
}

func (tree *TreapTree[T]) Delete(value T) {
	tree.SetRoot(delete(tree.Root().(*treapTreeNode[T]), value))
}
