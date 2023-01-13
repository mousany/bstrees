package avltree

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

type AVLTree[T node.Ordered] struct {
	*tree.BaseTree[T]
}

func New[T node.Ordered]() *AVLTree[T] {
	tree := &AVLTree[T]{tree.New[T]()}
	tree.SetRoot((*avlTreeNode[T])(nil))
	return tree
}

func insert[T node.Ordered](root *avlTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return newAvlTreeNode(value)
	}
	if value < root.Value() {
		root.SetLeft(insert(root.Left().(*avlTreeNode[T]), value))
	} else {
		root.SetRight(insert(root.Right().(*avlTreeNode[T]), value))
	}
	root.Update()
	return balance(root)
}

func (tree *AVLTree[T]) Insert(value T) {
	tree.SetRoot(insert(tree.Root().(*avlTreeNode[T]), value))
}

func delete[T node.Ordered](root *avlTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return nil
	}
	if value < root.Value() {
		root.SetLeft(delete(root.Left().(*avlTreeNode[T]), value))
	} else if root.Value() < value {
		root.SetRight(delete(root.Right().(*avlTreeNode[T]), value))
	} else {
		if root.Left().IsNil() {
			return root.Right().(*avlTreeNode[T])
		} else if root.Right().IsNil() {
			return root.Left().(*avlTreeNode[T])
		} else {
			minNode := tree.Kth(root.Right(), 1) // root.Right is not nil, so this will not fail
			root.SetValue(minNode.Value())
			root.SetRight(delete(root.Right().(*avlTreeNode[T]), minNode.Value()))
		}
	}
	root.Update()
	return balance(root)
}

func (tree *AVLTree[T]) Delete(value T) {
	tree.SetRoot(delete(tree.Root().(*avlTreeNode[T]), value))
}