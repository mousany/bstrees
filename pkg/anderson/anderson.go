package anderson

import (
	"bstrees/internal/node"
	"bstrees/internal/tree"
)

type AndersonTree[T node.Ordered] struct {
	*tree.BaseTree[T]
}

func New[T node.Ordered]() *AndersonTree[T] {
	tr := &AndersonTree[T]{
		BaseTree: tree.New[T](),
	}
	tr.SetRoot((*andersonTreeNode[T])(nil))
	return tr
}

func insert[T node.Ordered](root *andersonTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return newAndersonTreeNode(value, 1)
	}
	if value < root.Value() {
		root.SetLeft(insert(root.Left().(*andersonTreeNode[T]), value))
	} else {
		root.SetRight(insert(root.Right().(*andersonTreeNode[T]), value))
	}
	root.Update()
	root = skew(root).(*andersonTreeNode[T])
	root = split(root).(*andersonTreeNode[T])
	return root
}

func (tr *AndersonTree[T]) Insert(value T) {
	tr.SetRoot(insert(tr.Root().(*andersonTreeNode[T]), value))
}

func delete[T node.Ordered](root *andersonTreeNode[T], value T) node.Noded[T] {
	if root == nil {
		return nil
	}
	if value < root.Value() {
		root.SetLeft(delete(root.Left().(*andersonTreeNode[T]), value))
	} else if value > root.Value() {
		root.SetRight(delete(root.Right().(*andersonTreeNode[T]), value))
	} else {
		if root.Left().IsNil() {
			return root.Right()
		} else if root.Right().IsNil() {
			return root.Left()
		} else {
			minNode := tree.Kth(root.Right(), 1)
			root.SetValue(minNode.Value())
			root.SetRight(delete(root.Right().(*andersonTreeNode[T]), minNode.Value()))
		}
	}
	root.Update()
	if (!root.Left().IsNil() && root.Left().(*andersonTreeNode[T]).Level() < root.Level()-1) ||
		(!root.Right().IsNil() && root.Right().(*andersonTreeNode[T]).Level() < root.Level()-1) {
		root.SetLevel(root.Level() - 1)
		if root.Right().IsNil() && root.Right().(*andersonTreeNode[T]).Level() > root.Level() {
			root.Right().(*andersonTreeNode[T]).SetLevel(
				root.Level(),
			)
		}
		root = skew(root).(*andersonTreeNode[T])
		root = split(root).(*andersonTreeNode[T])
	}
	return root
}

func (tr *AndersonTree[T]) Delete(value T) {
	tr.SetRoot(delete(tr.Root().(*andersonTreeNode[T]), value))
}
